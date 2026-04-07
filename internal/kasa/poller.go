package kasa

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// pollInterval is how often the poller queries all known devices.
// 10 seconds -- Mother says that's often enough to check on things.
const pollInterval = 10 * time.Second

// Broadcaster is the interface for pushing events to connected WebSocket clients.
// George Michael is building the actual hub; we just need to know it can broadcast.
type Broadcaster interface {
	BroadcastGlobal(event interface{})
}

// BroadcastEvent is the envelope sent over the WebSocket hub.
type BroadcastEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Poller periodically queries each known Kasa device and broadcasts the
// aggregated state via the WebSocket hub.
type Poller struct {
	deviceIPs []string
	hub       Broadcaster
	devices   atomic.Value // []Device
	stop      chan struct{}
	mu        sync.Mutex
	running   bool
}

// NewPoller creates a Poller for the given device IPs. It does not start
// polling until Start is called. Hub may be nil if WebSocket isn't ready yet.
func NewPoller(deviceIPs []string, hub Broadcaster) *Poller {
	return &Poller{
		deviceIPs: deviceIPs,
		hub:       hub,
	}
}

// Start begins polling each device every 10 seconds. Safe to call from multiple
// goroutines; a second call while already running is a no-op.
func (p *Poller) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running {
		return
	}
	p.stop = make(chan struct{})
	p.running = true
	go p.loop()
	log.Println("[kasa] Poller started. I'll keep checking on them every 10 seconds.")
}

// Stop halts the poller. Safe to call when already stopped.
func (p *Poller) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running {
		return
	}
	close(p.stop)
	p.running = false
	log.Println("[kasa] Poller stopped. The devices are on their own now.")
}

// UpdateDeviceIPs replaces the set of IPs the poller monitors. The change
// takes effect on the next poll cycle.
func (p *Poller) UpdateDeviceIPs(ips []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.deviceIPs = make([]string, len(ips))
	copy(p.deviceIPs, ips)
}

// GetDevices returns the last cached list of devices, or nil if no successful
// poll has occurred yet.
func (p *Poller) GetDevices() []Device {
	v := p.devices.Load()
	if v == nil {
		return nil
	}
	devices, _ := v.([]Device)
	return devices
}

// Refresh triggers an immediate poll and broadcast. Call this after sending
// a command so the UI updates instantly instead of waiting for the next cycle.
func (p *Poller) Refresh() {
	go p.fetch()
}

func (p *Poller) loop() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Fetch immediately on start.
	p.fetch()

	for {
		select {
		case <-ticker.C:
			p.fetch()
		case <-p.stop:
			return
		}
	}
}

func (p *Poller) fetch() {
	p.mu.Lock()
	ips := make([]string, len(p.deviceIPs))
	copy(ips, p.deviceIPs)
	p.mu.Unlock()

	devices := make([]Device, 0, len(ips))
	for _, ip := range ips {
		dev, err := QueryDevice(ip)
		if err != nil {
			log.Printf("[kasa] Poller: query %s failed: %v -- I'm sure it's fine.", ip, err)
			continue
		}
		devices = append(devices, *dev)
	}

	p.devices.Store(devices)

	if p.hub != nil {
		p.hub.BroadcastGlobal(BroadcastEvent{
			Type:    "kasa:state",
			Payload: devices,
		})
	}
}
