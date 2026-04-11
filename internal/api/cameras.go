package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var go2rtcClient = &http.Client{Timeout: 10 * time.Second}
const go2rtcURL = "http://localhost:1984"

// HandleCameraList returns camera streams from go2rtc.
// GET /api/cameras
func (s *Server) HandleCameraList(w http.ResponseWriter, r *http.Request) {
	resp, err := go2rtcClient.Get(go2rtcURL + "/api/streams")
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "go2rtc not reachable")
		return
	}
	defer resp.Body.Close()

	var raw map[string]json.RawMessage
	json.NewDecoder(resp.Body).Decode(&raw)

	type stream struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	streams := make([]stream, 0)
	for id := range raw {
		streams = append(streams, stream{ID: id, Name: id})
	}
	s.writeJSON(w, http.StatusOK, streams)
}

// HandleCameraWebRTC proxies a WebRTC SDP exchange to go2rtc.
// POST /api/cameras/{id}/webrtc
func (s *Server) HandleCameraWebRTC(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.writeError(w, http.StatusBadRequest, "stream id required")
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 64*1024))
	if err != nil || len(body) == 0 {
		s.writeError(w, http.StatusBadRequest, "request body required")
		return
	}

	upstream := fmt.Sprintf("%s/api/webrtc?src=%s", go2rtcURL, id)
	req, _ := http.NewRequestWithContext(r.Context(), http.MethodPost, upstream, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := go2rtcClient.Do(req)
	if err != nil {
		log.Printf("[cameras] go2rtc unreachable: %v", err)
		s.writeError(w, http.StatusBadGateway, "go2rtc unreachable")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[cameras] go2rtc returned %d for stream %q", resp.StatusCode, id)
		s.writeError(w, http.StatusBadGateway, "go2rtc error")
		return
	}

	var answer json.RawMessage
	json.NewDecoder(resp.Body).Decode(&answer)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answer)
}

// HandleCameraStream proxies go2rtc's MP4/MSE stream for non-WebRTC browsers.
// GET /api/cameras/{id}/stream
func (s *Server) HandleCameraStream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.writeError(w, http.StatusBadRequest, "stream id required")
		return
	}

	upstream := fmt.Sprintf("%s/api/stream.mp4?src=%s", go2rtcURL, id)
	req, _ := http.NewRequestWithContext(r.Context(), http.MethodGet, upstream, nil)

	resp, err := go2rtcClient.Do(req)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "go2rtc unreachable")
		return
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}
