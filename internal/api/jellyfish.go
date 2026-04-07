package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// jellyfishRequest is the body for POST /api/devices/{ip}/jellyfish
type jellyfishRequest struct {
	Action  string   `json:"action"`  // "on", "off", "pattern"
	Zones   []string `json:"zones"`   // zone names
	Pattern string   `json:"pattern"` // pattern file path (for "on"/"pattern")
}

// HandleJellyfishCommand sends a command to a JellyFish controller via WebSocket.
// POST /api/devices/{ip}/jellyfish
func (s *Server) HandleJellyfishCommand(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	if ip == "" {
		s.writeError(w, http.StatusBadRequest, "ip required")
		return
	}

	var req jellyfishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Connect to JellyFish WebSocket
	dialer := websocket.Dialer{HandshakeTimeout: 3 * time.Second}
	conn, _, err := dialer.Dial(fmt.Sprintf("ws://%s:9000/", ip), nil)
	if err != nil {
		log.Printf("[jellyfish] Failed to connect to %s: %v", ip, err)
		s.writeError(w, http.StatusBadGateway, "failed to connect to JellyFish controller")
		return
	}
	defer conn.Close()

	var msg map[string]interface{}

	switch req.Action {
	case "on":
		zones := req.Zones
		if len(zones) == 0 {
			zones = []string{"Zone1", "Zone"} // default all zones
		}
		pattern := req.Pattern
		if pattern == "" {
			pattern = "Accent/White" // default pattern
		}
		msg = map[string]interface{}{
			"cmd": "toCtlrSet",
			"runPattern": map[string]interface{}{
				"state":    1,
				"zoneName": zones,
				"file":     pattern,
				"id":       "",
				"data":     "",
			},
		}
	case "off":
		zones := req.Zones
		if len(zones) == 0 {
			zones = []string{"Zone1", "Zone"}
		}
		msg = map[string]interface{}{
			"cmd": "toCtlrSet",
			"runPattern": map[string]interface{}{
				"state":    0,
				"zoneName": zones,
				"file":     "",
				"id":       "",
				"data":     "",
			},
		}
	case "pattern":
		zones := req.Zones
		if len(zones) == 0 {
			zones = []string{"Zone1", "Zone"}
		}
		msg = map[string]interface{}{
			"cmd": "toCtlrSet",
			"runPattern": map[string]interface{}{
				"state":    1,
				"zoneName": zones,
				"file":     req.Pattern,
				"id":       "",
				"data":     "",
			},
		}
	default:
		s.writeError(w, http.StatusBadRequest, "action must be 'on', 'off', or 'pattern'")
		return
	}

	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("[jellyfish] Failed to send command to %s: %v", ip, err)
		s.writeError(w, http.StatusBadGateway, "failed to send command")
		return
	}

	log.Printf("[jellyfish] Sent %s command to %s", req.Action, ip)
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"ok":     true,
		"action": req.Action,
	})
}
