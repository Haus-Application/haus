package discovery

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// kasaEncrypt applies the TP-Link Kasa XOR encryption to a plaintext command.
// The first 4 bytes are the big-endian length, followed by XOR-encrypted payload.
// I practiced this algorithm until Mother said I could stop. She never said stop.
func kasaEncrypt(plaintext string) []byte {
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

// kasaDecrypt reverses the XOR encryption on a Kasa response.
func kasaDecrypt(data []byte) string {
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

// kasaSysInfo represents the parts of a Kasa sysinfo response we care about.
type kasaSysInfo struct {
	System struct {
		GetSysinfo struct {
			Alias      string `json:"alias"`
			Model      string `json:"model"`
			DevName    string `json:"dev_name"`
			RelayState int    `json:"relay_state"`
			Brightness int    `json:"brightness"`
		} `json:"get_sysinfo"`
	} `json:"system"`
}

// probeKasa queries devices with port 9999 open using the Kasa XOR protocol.
// This is where it gets really exciting -- the XOR key is 171, which is
// 0xAB, which is the same as... well, it's just a nice number.
func probeKasa(session *ScanSession) {
	query := `{"system":{"get_sysinfo":{}}}`
	encrypted := kasaEncrypt(query)
	probed := 0

	for _, device := range session.Devices {
		if !hasPort(device, 9999) {
			continue
		}

		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:9999", device.IP), 2*time.Second)
		if err != nil {
			continue
		}

		conn.SetDeadline(time.Now().Add(3 * time.Second))
		_, err = conn.Write(encrypted)
		if err != nil {
			conn.Close()
			continue
		}

		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		conn.Close()
		if err != nil || n == 0 {
			continue
		}

		decrypted := kasaDecrypt(buf[:n])
		var info kasaSysInfo
		if err := json.Unmarshal([]byte(decrypted), &info); err != nil {
			continue
		}

		sysinfo := info.System.GetSysinfo

		if sysinfo.Alias != "" {
			device.Name = sysinfo.Alias
		}
		if sysinfo.Model != "" {
			device.Model = sysinfo.Model
		}

		// Classify by model
		model := strings.ToUpper(sysinfo.Model)
		switch {
		case strings.HasPrefix(model, "HS220"):
			device.DeviceType = "dimmer"
		case strings.HasPrefix(model, "HS200"):
			device.DeviceType = "switch"
		case strings.Contains(strings.ToLower(sysinfo.Alias), "fan"):
			device.DeviceType = "fan"
		default:
			device.DeviceType = "switch"
		}

		// Add Kasa protocol
		if !hasProtocol(device, "kasa") {
			device.Protocols = append(device.Protocols, "kasa")
		}

		device.Category = CategoryLighting

		// Store extra metadata
		if device.Metadata == nil {
			device.Metadata = make(map[string]string)
		}
		device.Metadata["relay_state"] = fmt.Sprintf("%d", sysinfo.RelayState)
		if sysinfo.Brightness > 0 {
			device.Metadata["brightness"] = fmt.Sprintf("%d", sysinfo.Brightness)
		}

		emitDevice(session, device)
		probed++
	}

	log.Printf("[kasa] Probed %d Kasa devices. The XOR protocol never lets me down.", probed)
}

// hasPort checks if a device has a specific port in its OpenPorts list.
func hasPort(device *Device, port int) bool {
	for _, p := range device.OpenPorts {
		if p == port {
			return true
		}
	}
	return false
}

// hasProtocol checks if a device already has a protocol in its list.
func hasProtocol(device *Device, protocol string) bool {
	for _, p := range device.Protocols {
		if p == protocol {
			return true
		}
	}
	return false
}
