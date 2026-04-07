// Package hue provides a Philips Hue API v2 client, bridge discovery,
// and a background poller for Haus. I've been studying the Hue bridge
// since before it was called a bridge. Mother has one in every room.
package hue

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BridgeInfo holds the minimal information needed to connect to a Hue bridge.
type BridgeInfo struct {
	ID string `json:"id"`
	IP string `json:"internalipaddress"`
}

// DiscoverBridges queries the Philips cloud discovery endpoint to find Hue
// bridges on the local network. The endpoint returns bridge IPs without
// requiring mDNS -- which is nice because mDNS can be... temperamental.
func DiscoverBridges(timeout time.Duration) ([]BridgeInfo, error) {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get("https://discovery.meethue.com/")
	if err != nil {
		return nil, fmt.Errorf("hue: discover bridges: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hue: discover bridges status %d: %s", resp.StatusCode, body)
	}

	var bridges []BridgeInfo
	if err := json.NewDecoder(resp.Body).Decode(&bridges); err != nil {
		return nil, fmt.Errorf("hue: decode discovery response: %w", err)
	}
	return bridges, nil
}
