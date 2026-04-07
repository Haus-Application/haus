package discovery

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// eurekaInfo represents the subset of the Google Cast eureka_info response
// that we actually need.
type eurekaInfo struct {
	Name string `json:"name"`
}

// probeCast queries devices with port 8008 open for Google Cast eureka_info.
// It's like asking "who are you?" but over HTTP, which is how I prefer
// to communicate anyway.
func probeCast(session *ScanSession) {
	client := &http.Client{Timeout: 3 * time.Second}
	probed := 0

	for _, device := range session.Devices {
		if !hasPort(device, 8008) {
			continue
		}

		url := fmt.Sprintf("http://%s:8008/setup/eureka_info", device.IP)
		resp, err := client.Get(url)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		var info eurekaInfo
		if err := json.Unmarshal(body, &info); err != nil {
			continue
		}

		if info.Name != "" {
			device.Name = info.Name
		}

		if !hasProtocol(device, "cast") {
			device.Protocols = append(device.Protocols, "cast")
		}

		device.Category = CategoryMedia

		emitDevice(session, device)
		probed++
	}

	log.Printf("[cast] Probed %d Cast devices. They all have such lovely names.", probed)
}
