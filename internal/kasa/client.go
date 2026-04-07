// Package kasa communicates with TP-Link Kasa smart switches, dimmers, and fans
// over the local network using the Kasa proprietary TCP protocol (port 9999).
// The protocol uses XOR encryption with key 171 -- no authentication required.
// I memorized the entire protocol spec. Mother said that was "a lot," but I
// think she meant "a lot of talent."
package kasa

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

// dialTimeout is the per-device TCP connect timeout.
const dialTimeout = 1 * time.Second

// Device represents a single Kasa smart switch, dimmer, or fan and its current state.
type Device struct {
	IP         string `json:"ip"`
	Alias      string `json:"alias"`
	Model      string `json:"model"`
	DeviceType string `json:"device_type"` // "dimmer", "switch", or "fan"
	On         bool   `json:"on"`
	Brightness int    `json:"brightness"` // 0-100, only meaningful for dimmers/fans
	FanSpeed   int    `json:"fan_speed"`  // 0-4, only meaningful for fan devices
}

// Encrypt XOR-encrypts a JSON string for the Kasa protocol.
// The result includes a 4-byte big-endian length prefix followed by the
// encrypted payload. Key 171 -- same as in discovery. I could do this in my sleep.
func Encrypt(plaintext string) []byte {
	key := byte(171)
	n := len(plaintext)
	result := make([]byte, 4+n)
	result[0] = byte(n >> 24)
	result[1] = byte(n >> 16)
	result[2] = byte(n >> 8)
	result[3] = byte(n)
	for i := 0; i < n; i++ {
		key = key ^ plaintext[i]
		result[4+i] = key
	}
	return result
}

// Decrypt XOR-decrypts a Kasa protocol response (including the 4-byte length
// prefix). Returns an empty string if data is shorter than the header.
func Decrypt(data []byte) string {
	if len(data) < 4 {
		return ""
	}
	key := byte(171)
	payload := data[4:]
	result := make([]byte, len(payload))
	for i, b := range payload {
		result[i] = b ^ key
		key = b
	}
	return string(result)
}

// sendCommand opens a TCP connection to ip:9999, sends the encrypted command,
// and returns the decrypted response JSON string.
func sendCommand(ip, command string) (string, error) {
	conn, err := net.DialTimeout("tcp", ip+":9999", dialTimeout)
	if err != nil {
		return "", fmt.Errorf("kasa: dial %s: %w", ip, err)
	}
	defer conn.Close()

	// 2-second read/write deadline per operation.
	if err := conn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return "", fmt.Errorf("kasa: set deadline: %w", err)
	}

	if _, err := conn.Write(Encrypt(command)); err != nil {
		return "", fmt.Errorf("kasa: write to %s: %w", ip, err)
	}

	// Read the 4-byte big-endian length prefix.
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return "", fmt.Errorf("kasa: read length from %s: %w", ip, err)
	}

	// Read exactly length bytes of encrypted payload.
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return "", fmt.Errorf("kasa: read payload from %s: %w", ip, err)
	}

	// Prepend a synthetic 4-byte header so Decrypt can strip it.
	buf := make([]byte, 4+len(payload))
	buf[0] = byte(length >> 24)
	buf[1] = byte(length >> 16)
	buf[2] = byte(length >> 8)
	buf[3] = byte(length)
	copy(buf[4:], payload)
	return Decrypt(buf), nil
}

// sysInfoResponse is the on-the-wire shape returned by get_sysinfo.
type sysInfoResponse struct {
	System struct {
		GetSysInfo struct {
			Alias      string `json:"alias"`
			Model      string `json:"model"`
			RelayState int    `json:"relay_state"` // 0 or 1
			Brightness int    `json:"brightness"`  // 0-100; present on HS220 dimmers
		} `json:"get_sysinfo"`
	} `json:"system"`
}

// QueryDevice sends a get_sysinfo request to the device at ip and returns its
// current state. This is the most reliable part of my day.
func QueryDevice(ip string) (*Device, error) {
	const cmd = `{"system":{"get_sysinfo":{}}}`
	raw, err := sendCommand(ip, cmd)
	if err != nil {
		return nil, err
	}

	var resp sysInfoResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("kasa: parse sysinfo from %s: %w", ip, err)
	}

	info := resp.System.GetSysInfo
	dt := "switch"
	if strings.Contains(info.Model, "HS220") {
		dt = "dimmer"
	}
	// Fan devices are identified by alias -- they use a dimmer under the hood
	// but expose speed levels instead of raw brightness.
	if strings.Contains(strings.ToLower(info.Alias), "fan") {
		dt = "fan"
	}

	// Compute fan speed from brightness for fan devices.
	fanSpeed := 0
	if dt == "fan" && info.Brightness > 0 {
		switch {
		case info.Brightness <= 25:
			fanSpeed = 1
		case info.Brightness <= 50:
			fanSpeed = 2
		case info.Brightness <= 75:
			fanSpeed = 3
		default:
			fanSpeed = 4
		}
	}

	return &Device{
		IP:         ip,
		Alias:      info.Alias,
		Model:      info.Model,
		DeviceType: dt,
		On:         info.RelayState == 1,
		Brightness: info.Brightness,
		FanSpeed:   fanSpeed,
	}, nil
}

// SetState turns the device at ip on (on=true) or off (on=false).
func SetState(ip string, on bool) error {
	state := 0
	if on {
		state = 1
	}
	cmd := fmt.Sprintf(`{"system":{"set_relay_state":{"state":%d}}}`, state)
	_, err := sendCommand(ip, cmd)
	return err
}

// SetBrightness sets the dimmer level (0-100) on the device at ip.
// Only meaningful for HS220 dimmers. Mother has opinions about brightness levels.
func SetBrightness(ip string, level int) error {
	cmd := fmt.Sprintf(`{"smartlife.iot.dimmer":{"set_brightness":{"brightness":%d}}}`, level)
	_, err := sendCommand(ip, cmd)
	return err
}

// SetFanSpeed sets the fan speed (1-4) on a dimmer-based fan device.
// Speed levels map to brightness: 1->25, 2->50, 3->75, 4->100.
func SetFanSpeed(ip string, speed int) error {
	brightness := speed * 25
	return SetBrightness(ip, brightness)
}

// DiscoverDevices scans all .1-.254 hosts on the given subnet (e.g.
// "192.168.1") and returns every device that responds on port 9999 within
// timeout. Discovery is concurrent; a semaphore limits parallel dials to 50.
func DiscoverDevices(subnet string, timeout time.Duration) ([]Device, error) {
	const maxParallel = 50

	type result struct {
		dev *Device
		err error
	}

	results := make(chan result, 254)
	sem := make(chan struct{}, maxParallel)
	var wg sync.WaitGroup

	deadline := time.Now().Add(timeout)

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s.%d", subnet, i)
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// Don't bother if we've already hit the overall deadline.
			if time.Now().After(deadline) {
				return
			}

			// Probe port 9999 first with a short dial timeout.
			conn, err := net.DialTimeout("tcp", ip+":9999", dialTimeout)
			if err != nil {
				return // not a Kasa device
			}
			conn.Close()

			dev, err := QueryDevice(ip)
			results <- result{dev: dev, err: err}
		}(ip)
	}

	// Close results once all goroutines finish.
	go func() {
		wg.Wait()
		close(results)
	}()

	var devices []Device
	for r := range results {
		if r.err == nil && r.dev != nil {
			devices = append(devices, *r.dev)
		}
	}
	return devices, nil
}
