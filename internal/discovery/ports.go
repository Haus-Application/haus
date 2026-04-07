package discovery

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// targetPorts is the list of ports we check on each device.
// I've studied every one of these port numbers. Port 9999 is Kasa,
// 8008 is Google Cast, 554 is RTSP -- I could go on. Mother says
// I shouldn't, but I could.
var targetPorts = []int{22, 80, 443, 554, 5455, 8008, 8009, 8080, 8443, 9999}

// scanPorts performs a concurrent TCP connect scan across all devices
// for the target port list. Uses a semaphore to limit concurrency to 100.
func scanPorts(session *ScanSession) {
	sem := make(chan struct{}, 100)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, device := range session.Devices {
		for _, port := range targetPorts {
			wg.Add(1)
			sem <- struct{}{}
			go func(dev *Device, p int) {
				defer wg.Done()
				defer func() { <-sem }()

				addr := fmt.Sprintf("%s:%d", dev.IP, p)
				conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
				if err != nil {
					return
				}
				conn.Close()

				mu.Lock()
				dev.OpenPorts = append(dev.OpenPorts, p)
				mu.Unlock()
			}(device, port)
		}
	}

	wg.Wait()

	// Emit device events for any devices that got new ports
	portsFound := 0
	for _, device := range session.Devices {
		if len(device.OpenPorts) > 0 {
			emitDevice(session, device)
			portsFound += len(device.OpenPorts)
		}
	}

	log.Printf("[ports] Found %d open ports across all devices. I knocked on every door!", portsFound)
}
