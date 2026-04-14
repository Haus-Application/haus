package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coalson/haus/internal/db"
	"github.com/coalson/haus/internal/kasa"
	"github.com/coalson/haus/internal/ws"
)

// ensureKasaPoller lazily constructs the Kasa poller the first time we need
// it after a scan has persisted Kasa IPs to the DB. Without this, users had
// to restart the server before toggling switches would work. Returns nil if
// no Kasa devices are known yet; callers must be nil-safe (the poller's own
// methods are too). Guarded by s.kasaMu so two simultaneous handler calls
// can't race to create duplicate pollers.
func (s *Server) ensureKasaPoller() *kasa.Poller {
	s.kasaMu.Lock()
	defer s.kasaMu.Unlock()

	if s.KasaPoller != nil {
		return s.KasaPoller
	}

	ips, err := db.LoadKasaIPs(s.DB)
	if err != nil || len(ips) == 0 {
		return nil
	}

	// Kasa poller needs a broadcaster that can emit BroadcastEvent to the
	// WebSocket hub. Use the same adapter pattern as main.go.
	var broadcaster kasa.Broadcaster
	if s.Hub != nil {
		broadcaster = &kasaHubAdapter{hub: s.Hub}
	}

	p := kasa.NewPoller(ips, broadcaster)
	p.Start()
	s.KasaPoller = p
	log.Printf("[kasa] Lazy-started poller with %d device(s) from DB", len(ips))
	return p
}

// kasaHubAdapter lets the kasa package broadcast via ws.Hub without importing it.
type kasaHubAdapter struct{ hub *ws.Hub }

func (a *kasaHubAdapter) BroadcastGlobal(event interface{}) {
	if a == nil || a.hub == nil {
		return
	}
	// kasa.BroadcastEvent has {Type, Payload}; forward as ws.BroadcastEvent.
	if e, ok := event.(kasa.BroadcastEvent); ok {
		a.hub.BroadcastGlobal(ws.BroadcastEvent{Type: e.Type, Payload: e.Payload})
	}
}

// HandleKasaDevices returns the cached device list from the Kasa poller.
//
// GET /api/kasa/devices
func (s *Server) HandleKasaDevices(w http.ResponseWriter, r *http.Request) {
	poller := s.ensureKasaPoller()
	devices := poller.GetDevices() // nil-safe
	if devices == nil {
		devices = []kasa.Device{}
	}
	s.writeJSON(w, http.StatusOK, devices)
}

// HandleKasaSetState turns a device on or off.
//
// PUT /api/kasa/devices/{ip}/state
// Body: {"on": true}
func (s *Server) HandleKasaSetState(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	if ip == "" {
		s.writeError(w, http.StatusBadRequest, "ip required")
		return
	}

	var body struct {
		On bool `json:"on"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := kasa.SetState(ip, body.On); err != nil {
		log.Printf("[kasa] Set state %s: %v -- the device isn't responding. I'm not panicking.", ip, err)
		s.writeError(w, http.StatusBadGateway, "failed to set device state")
		return
	}
	s.ensureKasaPoller().Refresh()
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// HandleKasaSetBrightness sets the dimmer brightness level on a device.
//
// PUT /api/kasa/devices/{ip}/brightness
// Body: {"brightness": 50}
func (s *Server) HandleKasaSetBrightness(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	if ip == "" {
		s.writeError(w, http.StatusBadRequest, "ip required")
		return
	}

	var body struct {
		Brightness int `json:"brightness"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Brightness < 0 || body.Brightness > 100 {
		s.writeError(w, http.StatusBadRequest, "brightness must be 0-100")
		return
	}

	if err := kasa.SetBrightness(ip, body.Brightness); err != nil {
		log.Printf("[kasa] Set brightness %s: %v", ip, err)
		s.writeError(w, http.StatusBadGateway, "failed to set brightness")
		return
	}
	s.ensureKasaPoller().Refresh()
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// HandleKasaSetFanSpeed sets the fan speed level on a fan device.
//
// PUT /api/kasa/devices/{ip}/fan-speed
// Body: {"speed": 3}
func (s *Server) HandleKasaSetFanSpeed(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	if ip == "" {
		s.writeError(w, http.StatusBadRequest, "ip required")
		return
	}

	var body struct {
		Speed int `json:"speed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Speed < 1 || body.Speed > 4 {
		s.writeError(w, http.StatusBadRequest, "speed must be 1-4")
		return
	}

	if err := kasa.SetFanSpeed(ip, body.Speed); err != nil {
		log.Printf("[kasa] Set fan speed %s: %v", ip, err)
		s.writeError(w, http.StatusBadGateway, "failed to set fan speed")
		return
	}
	s.ensureKasaPoller().Refresh()
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
