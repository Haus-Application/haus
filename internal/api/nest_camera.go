package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coalson/haus/internal/db"
	"github.com/coalson/haus/internal/nest"
)

// HandleNestCameraStream generates a live stream URL for a Nest camera.
// POST /api/google/camera/{deviceID}/stream
// The deviceID is the last segment of the SDM device name (after "devices/").
func (s *Server) HandleNestCameraStream(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("deviceID")
	if deviceID == "" {
		s.writeError(w, http.StatusBadRequest, "deviceID required")
		return
	}

	client, err := s.GetGoogleClient()
	if err != nil {
		s.writeError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	// Try RTSP first (works for Nest app cameras)
	streamURL, token, err := client.GetCameraLiveStreamURL(deviceID)
	if err != nil {
		log.Printf("[nest-camera] RTSP stream failed for %s: %v — trying WebRTC", deviceID, err)

		// Try WebRTC (for Google Home app cameras)
		// Frontend will need to handle SDP exchange
		var body struct {
			OfferSDP string `json:"offer_sdp"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		if body.OfferSDP != "" {
			answerSDP, mediaSessionID, expiresAt, webrtcErr := client.GenerateWebRtcStream(deviceID, body.OfferSDP)
			if webrtcErr != nil {
				log.Printf("[nest-camera] WebRTC also failed: %v", webrtcErr)
				s.writeError(w, http.StatusBadGateway, "Could not start camera stream. Camera may not support streaming via SDM API.")
				return
			}

			s.writeJSON(w, http.StatusOK, map[string]interface{}{
				"type":             "webrtc",
				"answer_sdp":       answerSDP,
				"media_session_id": mediaSessionID,
				"expires_at":       expiresAt,
			})
			return
		}

		// No SDP offer — tell frontend to use WebRTC
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"type":    "webrtc_required",
			"message": "This camera requires WebRTC. Send an offer_sdp in the request body.",
		})
		return
	}

	log.Printf("[nest-camera] RTSP stream started for %s", deviceID)
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"type":            "rtsp",
		"stream_url":      streamURL,
		"extension_token": token,
		"expires_in":      300, // 5 minutes
	})
}

// HandleNestCameraExtend extends an active camera stream.
// POST /api/google/camera/{deviceID}/extend
func (s *Server) HandleNestCameraExtend(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("deviceID")
	var body struct {
		Token string `json:"token"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	client, err := s.GetGoogleClient()
	if err != nil {
		s.writeError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	newToken, err := client.ExtendCameraStream(deviceID, body.Token)
	if err != nil {
		s.writeError(w, http.StatusBadGateway, "Failed to extend stream")
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"extension_token": newToken,
		"expires_in":      300,
	})
}

// HandleNestCameraStop stops an active camera stream.
// POST /api/google/camera/{deviceID}/stop
func (s *Server) HandleNestCameraStop(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("deviceID")
	var body struct {
		Token string `json:"token"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	client, err := s.GetGoogleClient()
	if err != nil {
		s.writeError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	client.StopCameraStream(deviceID, body.Token)
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// GetGoogleClient creates a Nest SDK client with valid tokens, refreshing if needed.
func (s *Server) GetGoogleClient() (*nest.Client, error) {
	tokens, err := db.LoadGoogleTokens(s.DB)
	if err != nil {
		return nil, fmt.Errorf("Google not connected — sign in first")
	}

	accessToken := tokens.AccessToken

	// Refresh if expired
	if time.Now().After(tokens.ExpiresAt) {
		newTokens, err := nest.RefreshToken(nest.OAuthConfig{
			ClientID:     s.GoogleClientID,
			ClientSecret: s.GoogleClientSecret,
			ProjectID:    s.GoogleProjectID,
		}, tokens.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh Google token: %w", err)
		}

		accessToken = newTokens.AccessToken
		expiresAt := time.Now().Add(time.Duration(newTokens.ExpiresIn) * time.Second)
		refreshToken := tokens.RefreshToken
		if newTokens.RefreshToken != "" {
			refreshToken = newTokens.RefreshToken
		}
		db.SaveGoogleTokens(s.DB, accessToken, refreshToken, expiresAt)
	}

	return nest.NewClient(s.GoogleProjectID, accessToken), nil
}
