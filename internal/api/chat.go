package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"

	"github.com/coalson/haus/internal/ai"
)

// chatRequest is the JSON body for POST /api/chat.
type chatRequest struct {
	Message string                   `json:"message"`
	History []anthropic.MessageParam `json:"history"`
}

// chatResponse is the JSON body returned from POST /api/chat.
type chatResponse struct {
	Text      string                   `json:"text"`
	ToolCalls []ai.ToolCallResult      `json:"tool_calls,omitempty"`
	Messages  []anthropic.MessageParam `json:"messages"`
}

// HandleChat processes a natural language message through the AI concierge.
// POST /api/chat
//
// The Concierge field must be set on the Server struct for this to work.
// If it's nil, the endpoint returns 503 — the show can't go on without the
// magician.
func (s *Server) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}

	if s.Concierge == nil {
		s.writeError(w, http.StatusServiceUnavailable, "AI concierge is not configured")
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Message == "" {
		s.writeError(w, http.StatusBadRequest, "message is required")
		return
	}

	log.Printf("[chat] received: %q", req.Message)

	result, err := s.Concierge.Chat(r.Context(), req.Message, req.History)
	if err != nil {
		log.Printf("[chat] error: %v", err)
		s.writeError(w, http.StatusInternalServerError, "failed to process message")
		return
	}

	s.writeJSON(w, http.StatusOK, chatResponse{
		Text:      result.Text,
		ToolCalls: result.ToolCalls,
		Messages:  result.Messages,
	})
}

// deviceChatRequest is the JSON body for POST /api/chat/device.
type deviceChatRequest struct {
	Device  ai.DeviceContext         `json:"device"`
	Message string                   `json:"message"`
	History []anthropic.MessageParam `json:"history"`
}

// HandleDeviceChat processes a message scoped to a single device.
// POST /api/chat/device
func (s *Server) HandleDeviceChat(w http.ResponseWriter, r *http.Request) {
	if s.Concierge == nil {
		s.writeError(w, http.StatusServiceUnavailable, "AI concierge is not configured")
		return
	}

	var req deviceChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Message == "" {
		s.writeError(w, http.StatusBadRequest, "message is required")
		return
	}
	if req.Device.IP == "" {
		s.writeError(w, http.StatusBadRequest, "device is required")
		return
	}

	log.Printf("[device-chat] %s (%s): %q", req.Device.Name, req.Device.IP, req.Message)

	result, err := s.Concierge.DeviceChat(r.Context(), req.Device, req.Message, req.History)
	if err != nil {
		log.Printf("[device-chat] error: %v", err)
		s.writeError(w, http.StatusInternalServerError, "failed to process message")
		return
	}

	s.writeJSON(w, http.StatusOK, chatResponse{
		Text:      result.Text,
		ToolCalls: result.ToolCalls,
		Messages:  result.Messages,
	})
}
