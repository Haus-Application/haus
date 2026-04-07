package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// startScanRequest is the optional JSON body for POST /api/scan.
type startScanRequest struct {
	Subnet string `json:"subnet"`
}

// HandleStartScan initiates a network scan. Mother said I could scan
// whenever I want, as long as I tell the API about it first.
func (s *Server) HandleStartScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req startScanRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.writeError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}
	}

	session := s.Scanner.StartScan(req.Subnet)
	log.Printf("[api] Scan started: %s. Here we go!", session.ID)

	s.writeJSON(w, http.StatusOK, map[string]string{
		"scan_id": session.ID,
		"status":  "running",
	})
}

// HandleScanStream provides a Server-Sent Events stream of scan progress.
// This is my favorite handler because I get to emit events in real-time,
// which is almost as exciting as when Mother lets me use the remote.
func (s *Server) HandleScanStream(w http.ResponseWriter, r *http.Request) {
	scanID := r.URL.Query().Get("scan_id")
	if scanID == "" {
		s.writeError(w, http.StatusBadRequest, "scan_id is required")
		return
	}

	session := s.Scanner.GetSession(scanID)
	if session == nil {
		s.writeError(w, http.StatusNotFound, "Scan session not found")
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		s.writeError(w, http.StatusInternalServerError, "Streaming not supported")
		return
	}

	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			log.Printf("[api] SSE client disconnected for scan %s", scanID)
			return
		case event, ok := <-session.Events:
			if !ok {
				// Channel closed -- scan is complete
				return
			}
			data, err := json.Marshal(event.Data)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, string(data))
			flusher.Flush()
		}
	}
}

// HandleScanResults returns the full device list from a completed scan.
func (s *Server) HandleScanResults(w http.ResponseWriter, r *http.Request) {
	scanID := r.URL.Query().Get("scan_id")
	if scanID == "" {
		s.writeError(w, http.StatusBadRequest, "scan_id is required")
		return
	}

	session := s.Scanner.GetSession(scanID)
	if session == nil {
		s.writeError(w, http.StatusNotFound, "Scan session not found")
		return
	}

	// Convert device map to slice
	devices := make([]interface{}, 0, len(session.Devices))
	for _, device := range session.Devices {
		devices = append(devices, device)
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"scan_id":      session.ID,
		"status":       session.Status,
		"device_count": len(session.Devices),
		"devices":      devices,
	})
}
