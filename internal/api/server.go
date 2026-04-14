package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/coalson/haus/internal/ai"
	"github.com/coalson/haus/internal/db"
	"github.com/coalson/haus/internal/discovery"
	"github.com/coalson/haus/internal/hue"
	"github.com/coalson/haus/internal/kasa"
	"github.com/coalson/haus/internal/kb"
	"github.com/coalson/haus/internal/ws"
)

// Server holds the dependencies for HTTP handlers.
type Server struct {
	DB         *sql.DB
	Scanner    *discovery.Scanner
	KasaPoller *kasa.Poller // may be nil until devices are discovered
	HueClient  *hue.Client  // may be nil until paired
	HuePoller  *hue.Poller  // may be nil until paired
	Concierge  *ai.Concierge // may be nil if no API key configured

	// KB is the device knowledge base, loaded from docs/devices/ at startup.
	// May be nil if the KB files are missing — handlers should guard.
	KB *kb.Catalog

	// Hub is the WebSocket hub for broadcasting live events. May be nil.
	Hub *ws.Hub

	// ValidationDir is the output directory for validation reports.
	// main.go sets this via resolveRuntimePath. Default "validation".
	ValidationDir string

	// validation run state (one run at a time)
	validationMu      sync.Mutex
	validationRunning bool
	validationJobID   string

	// APIKey is the Anthropic key for running validations server-side.
	APIKey string

	// Google Nest SDM OAuth config -- keep these close to the chest.
	GoogleClientID     string
	GoogleClientSecret string
	GoogleProjectID    string
}

// NewServer creates a new API server with the given dependencies.
func NewServer(db *sql.DB, scanner *discovery.Scanner) *Server {
	return &Server{
		DB:      db,
		Scanner: scanner,
	}
}

// writeJSON sends a JSON response with the given status code.
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[api] Failed to encode JSON response: %v", err)
	}
}

// writeError sends a JSON error response.
func (s *Server) writeError(w http.ResponseWriter, status int, msg string) {
	s.writeJSON(w, status, map[string]string{"error": msg})
}

// HandleDevices returns all persisted devices from the database.
// GET /api/devices
func (s *Server) HandleDevices(w http.ResponseWriter, r *http.Request) {
	rows, err := db.LoadAllDevices(s.DB)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to load devices")
		return
	}

	devices := make([]map[string]interface{}, 0, len(rows))
	for _, d := range rows {
		dev := map[string]interface{}{
			"ip":           d.IP,
			"mac":          d.MAC,
			"hostname":     d.Hostname,
			"name":         d.Name,
			"manufacturer": d.Manufacturer,
			"model":        d.Model,
			"device_type":  d.DeviceType,
			"category":     d.Category,
		}
		// Parse JSON arrays/objects back to proper types
		var protocols []string
		json.Unmarshal([]byte(d.Protocols), &protocols)
		dev["protocols"] = protocols

		var services []string
		json.Unmarshal([]byte(d.Services), &services)
		dev["services"] = services

		var openPorts []int
		json.Unmarshal([]byte(d.OpenPorts), &openPorts)
		dev["open_ports"] = openPorts

		var metadata map[string]string
		json.Unmarshal([]byte(d.Metadata), &metadata)
		dev["metadata"] = metadata

		devices = append(devices, dev)
	}

	s.writeJSON(w, http.StatusOK, devices)
}
