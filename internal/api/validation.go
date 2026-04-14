package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coalson/haus/internal/validation"
	"github.com/coalson/haus/internal/ws"
	"github.com/google/uuid"
)

// deviceListEntry is the lean per-device summary for the validation grid.
type deviceListEntry struct {
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	Category          string `json:"category"`
	IntegrationKey    string `json:"integration_key"`
	IntegrationStatus string `json:"integration_status"`
	TotalPct          *int   `json:"total_pct"` // nil if not yet validated
	Ran               bool   `json:"ran"`
}

// HandleValidationSummary returns the aggregate summary JSON.
// GET /api/validation/summary
func (s *Server) HandleValidationSummary(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.ValidationDir, "summary.json")
	data, err := os.ReadFile(path)
	if err != nil {
		s.writeJSON(w, http.StatusNotFound, map[string]string{"status": "no_runs_yet"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

// HandleValidationDeviceList returns all catalog devices with latest score.
// GET /api/validation/devices
func (s *Server) HandleValidationDeviceList(w http.ResponseWriter, r *http.Request) {
	if s.KB == nil {
		s.writeJSON(w, http.StatusOK, []deviceListEntry{})
		return
	}
	out := make([]deviceListEntry, 0, len(s.KB.All))
	for _, d := range s.KB.All {
		entry := deviceListEntry{
			Slug:              d.Slug,
			Name:              d.Name,
			Category:          d.Category,
			IntegrationKey:    d.IntegrationKey,
			IntegrationStatus: d.IntegrationStatus,
			Ran:               false,
		}
		// Peek the device report for its score, if any.
		reportPath := filepath.Join(s.ValidationDir, "devices", d.Slug+".json")
		if data, err := os.ReadFile(reportPath); err == nil {
			var rep validation.DeviceReport
			if json.Unmarshal(data, &rep) == nil {
				pct := rep.Score.TotalPct
				entry.TotalPct = &pct
				entry.Ran = true
			}
		}
		out = append(out, entry)
	}
	s.writeJSON(w, http.StatusOK, out)
}

// HandleValidationDeviceReport returns one device's full report.
// GET /api/validation/devices/{slug}
func (s *Server) HandleValidationDeviceReport(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" || strings.ContainsAny(slug, "/\\.") {
		s.writeError(w, http.StatusBadRequest, "invalid slug")
		return
	}
	path := filepath.Join(s.ValidationDir, "devices", slug+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		s.writeJSON(w, http.StatusNotFound, map[string]string{"status": "not_run", "slug": slug})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

type validationRunRequest struct {
	Only string `json:"only,omitempty"`
}

type validationRunResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HandleValidationRun kicks off a validation run in the background.
// POST /api/validation/run  body: {"only":"optional-slug"}
// Returns 202 with {job_id}. Progress streams via WebSocket event
// "validation:progress". Returns 409 if a run is already in progress.
func (s *Server) HandleValidationRun(w http.ResponseWriter, r *http.Request) {
	if s.KB == nil {
		s.writeError(w, http.StatusServiceUnavailable, "knowledge base not loaded")
		return
	}
	if s.APIKey == "" {
		s.writeError(w, http.StatusServiceUnavailable, "no Anthropic API key configured")
		return
	}

	var req validationRunRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	s.validationMu.Lock()
	if s.validationRunning {
		current := s.validationJobID
		s.validationMu.Unlock()
		s.writeJSON(w, http.StatusConflict, validationRunResponse{JobID: current, Status: "already_running"})
		return
	}
	jobID := uuid.New().String()
	s.validationRunning = true
	s.validationJobID = jobID
	s.validationMu.Unlock()

	go s.runValidation(context.Background(), jobID, req.Only)

	s.writeJSON(w, http.StatusAccepted, validationRunResponse{JobID: jobID, Status: "running"})
}

// runValidation is the goroutine body for a validation run. It emits
// `validation:progress` events through the WebSocket hub (if present) and
// clears the running flag when done.
func (s *Server) runValidation(ctx context.Context, jobID, only string) {
	defer func() {
		if rec := recover(); rec != nil {
			s.broadcastValidation(map[string]any{
				"job_id": jobID,
				"status": "failed",
				"error":  fmt.Sprintf("panic: %v", rec),
			})
		}
		s.validationMu.Lock()
		s.validationRunning = false
		s.validationJobID = ""
		s.validationMu.Unlock()
	}()

	runner := validation.NewRunner(s.KB, s.APIKey, "", s.ValidationDir, 8)

	if only != "" {
		dev, ok := s.KB.BySlug[only]
		if !ok {
			s.broadcastValidation(map[string]any{"job_id": jobID, "status": "failed", "error": "unknown slug " + only})
			return
		}
		s.broadcastValidation(map[string]any{"job_id": jobID, "status": "running", "slug": only, "done": 0, "total": 1})
		rep, err := runner.RunOne(ctx, dev.Slug)
		if err != nil {
			s.broadcastValidation(map[string]any{"job_id": jobID, "status": "failed", "slug": only, "error": err.Error()})
			return
		}
		s.broadcastValidation(map[string]any{"job_id": jobID, "status": "done", "slug": only, "done": 1, "total": 1, "score": rep.Score.TotalPct})
		return
	}

	total := len(s.KB.All)
	s.broadcastValidation(map[string]any{"job_id": jobID, "status": "running", "done": 0, "total": total})

	progress := func(p validation.Progress) {
		s.broadcastValidation(map[string]any{
			"job_id": jobID,
			"slug":   p.Slug,
			"done":   p.Done,
			"total":  p.Total,
			"score":  p.Score,
			"status": p.Status,
			"error":  p.Error,
		})
	}

	sum, err := runner.RunAll(ctx, progress)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.broadcastValidation(map[string]any{"job_id": jobID, "status": "failed", "error": err.Error()})
		return
	}
	s.broadcastValidation(map[string]any{
		"job_id":    jobID,
		"status":    "done",
		"done":      total,
		"total":     total,
		"avg_score": sum.AvgScore,
		"ran_at":    time.Now().UTC(),
	})
}

func (s *Server) broadcastValidation(payload map[string]any) {
	if s.Hub == nil {
		return
	}
	s.Hub.BroadcastGlobal(ws.BroadcastEvent{
		Type:    "validation:progress",
		Payload: payload,
	})
}
