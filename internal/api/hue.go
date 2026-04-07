package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/coalson/haus/internal/db"
	"github.com/coalson/haus/internal/hue"
)

// hueMu guards dynamic creation of HueClient and HuePoller after pairing.
var hueMu sync.Mutex

// HandleHueDiscover queries the Philips cloud endpoint and returns discovered
// Hue bridges on the local network.
//
// GET /api/hue/discover
func (s *Server) HandleHueDiscover(w http.ResponseWriter, r *http.Request) {
	bridges, err := hue.DiscoverBridges(10e9) // 10s
	if err != nil {
		log.Printf("[hue] Discover: %v -- the bridges are hiding from me.", err)
		s.writeError(w, http.StatusBadGateway, "bridge discovery failed")
		return
	}
	s.writeJSON(w, http.StatusOK, bridges)
}

// HandleHuePair initiates the bridge link-button pairing flow.
// The user must have pressed the physical button before calling this endpoint.
// On success the config is saved to the database and the client+poller are
// activated immediately.
//
// POST /api/hue/pair
// Body: {"bridge_ip":"192.168.x.x"}
func (s *Server) HandleHuePair(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BridgeIP string `json:"bridge_ip"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.BridgeIP == "" {
		s.writeError(w, http.StatusBadRequest, "bridge_ip is required")
		return
	}

	username, err := hue.Pair(body.BridgeIP)
	if err != nil {
		log.Printf("[hue] Pair %s: %v -- did they press the button? They never press the button.", body.BridgeIP, err)
		s.writeError(w, http.StatusBadGateway, "pairing failed: "+err.Error())
		return
	}

	if err := db.SaveHueConfig(s.DB, body.BridgeIP, username, ""); err != nil {
		log.Printf("[hue] Save config: %v", err)
		s.writeError(w, http.StatusInternalServerError, "failed to save Hue config")
		return
	}

	// Activate client and poller dynamically -- no server restart needed.
	hueMu.Lock()
	if s.HuePoller != nil {
		s.HuePoller.Stop()
	}
	s.HueClient = hue.NewClient(body.BridgeIP, username)
	s.HuePoller = hue.NewPoller(s.HueClient, nil) // hub will be wired up by George Michael
	s.HuePoller.Start()
	hueMu.Unlock()

	log.Printf("[hue] Paired with bridge %s. Mother would approve of this connection.", body.BridgeIP)
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"ok":        true,
		"bridge_ip": body.BridgeIP,
	})
}

// HandleHueStatus returns the current Hue integration connection status.
//
// GET /api/hue/status
func (s *Server) HandleHueStatus(w http.ResponseWriter, r *http.Request) {
	cfg, err := db.LoadHueConfig(s.DB)
	if err != nil {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"connected": false,
		})
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"connected": true,
		"bridge_ip": cfg.BridgeIP,
		"bridge_id": cfg.BridgeID,
		"running":   s.HuePoller != nil && s.HuePoller.IsRunning(),
	})
}

// HandleHueDisconnect removes the stored Hue config and stops the poller.
//
// DELETE /api/hue/disconnect
func (s *Server) HandleHueDisconnect(w http.ResponseWriter, r *http.Request) {
	hueMu.Lock()
	if s.HuePoller != nil {
		s.HuePoller.Stop()
		s.HuePoller = nil
	}
	s.HueClient = nil
	hueMu.Unlock()

	if err := db.DeleteHueConfig(s.DB); err != nil {
		log.Printf("[hue] Delete config: %v", err)
		s.writeError(w, http.StatusInternalServerError, "failed to remove Hue config")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// HandleHueLights returns the cached list of lights from the poller.
//
// GET /api/hue/lights
func (s *Server) HandleHueLights(w http.ResponseWriter, r *http.Request) {
	if s.HuePoller == nil {
		s.writeError(w, http.StatusServiceUnavailable, "Hue not configured")
		return
	}
	lights := s.HuePoller.GetLights()
	if lights == nil {
		lights = []hue.Light{}
	}
	s.writeJSON(w, http.StatusOK, lights)
}

// HandleHueSetLight updates a single light's state.
// Body fields are all optional; omit what you don't want to change.
//
// PUT /api/hue/lights/{id}
// Body: {"on":true,"brightness":75.0,"color_xy":[0.3,0.4]}
func (s *Server) HandleHueSetLight(w http.ResponseWriter, r *http.Request) {
	if s.HueClient == nil {
		s.writeError(w, http.StatusServiceUnavailable, "Hue not configured")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		s.writeError(w, http.StatusBadRequest, "light id required")
		return
	}

	var body struct {
		On         *bool       `json:"on"`
		Brightness *float64    `json:"brightness"`
		ColorXY    *[2]float64 `json:"color_xy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := s.HueClient.SetLightState(id, body.On, body.Brightness, body.ColorXY); err != nil {
		log.Printf("[hue] Set light %s: %v", id, err)
		s.writeError(w, http.StatusBadGateway, "failed to update light")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// HandleHueRooms returns the cached list of rooms from the poller.
//
// GET /api/hue/rooms
func (s *Server) HandleHueRooms(w http.ResponseWriter, r *http.Request) {
	if s.HuePoller == nil {
		s.writeError(w, http.StatusServiceUnavailable, "Hue not configured")
		return
	}
	rooms := s.HuePoller.GetRooms()
	if rooms == nil {
		rooms = []hue.Room{}
	}
	// Ensure no nil light slices -- JSON null would crash the frontend,
	// and Maeby would never let me hear the end of it.
	for i := range rooms {
		if rooms[i].Lights == nil {
			rooms[i].Lights = []hue.Light{}
		}
	}
	s.writeJSON(w, http.StatusOK, rooms)
}

// HandleHueSetRoom updates all lights in a room via its grouped_light resource.
//
// PUT /api/hue/rooms/{id}
// Body: {"on":true,"brightness":50.0}
func (s *Server) HandleHueSetRoom(w http.ResponseWriter, r *http.Request) {
	if s.HueClient == nil {
		s.writeError(w, http.StatusServiceUnavailable, "Hue not configured")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		s.writeError(w, http.StatusBadRequest, "room id required")
		return
	}

	var body struct {
		On         *bool    `json:"on"`
		Brightness *float64 `json:"brightness"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// The id in the URL is the room id; we need the grouped_light id.
	// Look it up from the poller cache.
	groupedLightID := ""
	if s.HuePoller != nil {
		for _, room := range s.HuePoller.GetRooms() {
			if room.ID == id {
				groupedLightID = room.GroupedLightID
				break
			}
		}
	}
	if groupedLightID == "" {
		// Fall back to treating id as grouped_light id directly.
		groupedLightID = id
	}

	if err := s.HueClient.SetGroupedLightState(groupedLightID, body.On, body.Brightness); err != nil {
		log.Printf("[hue] Set room %s (grouped_light %s): %v", id, groupedLightID, err)
		s.writeError(w, http.StatusBadGateway, "failed to update room")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// HandleHueScenes returns the cached list of scenes from the poller.
//
// GET /api/hue/scenes
func (s *Server) HandleHueScenes(w http.ResponseWriter, r *http.Request) {
	if s.HuePoller == nil {
		s.writeError(w, http.StatusServiceUnavailable, "Hue not configured")
		return
	}
	scenes := s.HuePoller.GetScenes()
	if scenes == nil {
		scenes = []hue.Scene{}
	}
	s.writeJSON(w, http.StatusOK, scenes)
}

// HandleHueActivateScene recalls the named scene on the bridge.
//
// POST /api/hue/scenes/{id}/activate
func (s *Server) HandleHueActivateScene(w http.ResponseWriter, r *http.Request) {
	if s.HueClient == nil {
		s.writeError(w, http.StatusServiceUnavailable, "Hue not configured")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		s.writeError(w, http.StatusBadRequest, "scene id required")
		return
	}

	if err := s.HueClient.ActivateScene(id); err != nil {
		log.Printf("[hue] Activate scene %s: %v", id, err)
		s.writeError(w, http.StatusBadGateway, "failed to activate scene")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
