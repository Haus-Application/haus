package hue

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

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

// Poller periodically fetches light, room, and scene state from the Hue bridge
// and broadcasts changes via the WebSocket hub.
type Poller struct {
	client  *Client
	hub     Broadcaster
	lights  atomic.Value // []Light
	rooms   atomic.Value // []Room
	scenes  atomic.Value // []Scene
	stop    chan struct{}
	mu      sync.Mutex
	running bool
}

// NewPoller creates a Poller. It does not start polling until Start is called.
// Hub may be nil if WebSocket isn't ready yet.
func NewPoller(client *Client, hub Broadcaster) *Poller {
	return &Poller{
		client: client,
		hub:    hub,
	}
}

// Start begins polling the Hue bridge every 5 seconds. Safe to call from
// multiple goroutines; a second call while running is a no-op.
func (p *Poller) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running {
		return
	}
	p.stop = make(chan struct{})
	p.running = true
	go p.loop()
	log.Println("[hue] Poller started. Checking Mother's lights every 5 seconds.")
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
	log.Println("[hue] Poller stopped.")
}

// IsRunning reports whether the poller is currently active.
func (p *Poller) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

// GetLights returns the last cached list of lights, or nil.
func (p *Poller) GetLights() []Light {
	v := p.lights.Load()
	if v == nil {
		return nil
	}
	lights, _ := v.([]Light)
	return lights
}

// GetRooms returns the last cached list of rooms, or nil.
func (p *Poller) GetRooms() []Room {
	v := p.rooms.Load()
	if v == nil {
		return nil
	}
	rooms, _ := v.([]Room)
	return rooms
}

// GetScenes returns the last cached list of scenes, or nil.
func (p *Poller) GetScenes() []Scene {
	v := p.scenes.Load()
	if v == nil {
		return nil
	}
	scenes, _ := v.([]Scene)
	return scenes
}

func (p *Poller) loop() {
	ticker := time.NewTicker(5 * time.Second)
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
	lights, err := p.client.ListLights()
	if err != nil {
		log.Printf("[hue] Poller: list lights: %v", err)
		return
	}
	p.lights.Store(lights)

	rooms, err := p.client.ListRooms()
	if err != nil {
		log.Printf("[hue] Poller: list rooms: %v", err)
		return
	}
	p.rooms.Store(rooms)

	scenes, err := p.client.ListScenes()
	if err != nil {
		log.Printf("[hue] Poller: list scenes: %v", err)
		return
	}
	p.scenes.Store(scenes)

	if p.hub != nil {
		p.hub.BroadcastGlobal(BroadcastEvent{
			Type: "hue:state",
			Payload: map[string]interface{}{
				"lights": lights,
				"rooms":  rooms,
				"scenes": scenes,
			},
		})
	}
}
