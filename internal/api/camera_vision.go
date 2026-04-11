package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// HandleCameraSnapshot captures a frame from go2rtc and returns it as JPEG.
// GET /api/cameras/{id}/snapshot
func (s *Server) HandleCameraSnapshot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.writeError(w, http.StatusBadRequest, "stream id required")
		return
	}

	snapshot, err := captureSnapshot(id)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "failed to capture snapshot")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(snapshot)
}

// captureSnapshot grabs a JPEG frame from go2rtc.
func captureSnapshot(streamID string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://localhost:1984/api/frame.jpeg?src=%s", streamID))
	if err != nil {
		return nil, fmt.Errorf("go2rtc snapshot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("go2rtc returned %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024)) // 5MB max
	if err != nil {
		return nil, fmt.Errorf("reading snapshot: %w", err)
	}

	log.Printf("[camera-vision] Captured %dKB snapshot from stream %s", len(data)/1024, streamID)
	return data, nil
}

// CaptureSnapshotBase64 grabs a snapshot and returns it as base64 for the AI.
func CaptureSnapshotBase64(streamID string) (string, error) {
	data, err := captureSnapshot(streamID)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
