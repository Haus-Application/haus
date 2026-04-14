package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/coalson/haus/internal/db"
)

// HandleGoogleAuthStart initiates the Google Nest OAuth2 flow.
// Builds the authorization URL and redirects the user to Google's consent
// page. You want to get into the Nest? You go through the front door.
// Nobody sneaks in on my watch.
//
// GET /api/google/auth
func (s *Server) HandleGoogleAuthStart(w http.ResponseWriter, r *http.Request) {
	if s.GoogleClientID == "" || s.GoogleProjectID == "" {
		s.writeError(w, http.StatusServiceUnavailable, "Google Nest not configured -- missing client ID or project ID")
		return
	}

	// Build the redirect URI from the incoming request so it works in
	// any environment. The callback endpoint is always /api/google/callback.
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	redirectURI := fmt.Sprintf("%s://%s/api/google/callback", scheme, r.Host)

	// Google Nest partner connections auth endpoint -- not the standard
	// accounts.google.com flow. The SDM API has its own door.
	authURL := fmt.Sprintf(
		"https://nestservices.google.com/partnerconnections/%s/auth?redirect_uri=%s&access_type=offline&prompt=consent&client_id=%s&response_type=code&scope=%s",
		url.PathEscape(s.GoogleProjectID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(s.GoogleClientID),
		url.QueryEscape("https://www.googleapis.com/auth/sdm.service"),
	)

	log.Printf("[google] Redirecting to Google consent page. If you don't come back, I understand.")
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleGoogleAuthCallback handles the OAuth2 callback from Google.
// Exchanges the authorization code for access and refresh tokens, then
// stores them securely in the database. No plaintext tokens lying around
// -- I learned that lesson in Balboa Towers.
//
// GET /api/google/callback?code=xxx
func (s *Server) HandleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errMsg := r.URL.Query().Get("error")
		if errMsg == "" {
			errMsg = "no authorization code received"
		}
		log.Printf("[google] OAuth callback error: %s", errMsg)
		s.writeError(w, http.StatusBadRequest, "authorization failed: "+errMsg)
		return
	}

	// Build the same redirect URI that was used in the auth request.
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	redirectURI := fmt.Sprintf("%s://%s/api/google/callback", scheme, r.Host)

	// Exchange the code for tokens via Google's token endpoint.
	tokenResp, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", url.Values{
		"client_id":     {s.GoogleClientID},
		"client_secret": {s.GoogleClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectURI},
	})
	if err != nil {
		log.Printf("[google] Token exchange failed: %v", err)
		s.writeError(w, http.StatusBadGateway, "failed to exchange authorization code")
		return
	}
	defer tokenResp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(tokenResp.Body, 16384))

	if tokenResp.StatusCode != http.StatusOK {
		log.Printf("[google] Token exchange returned %d: %s", tokenResp.StatusCode, string(body))
		s.writeError(w, http.StatusBadGateway, "Google token exchange failed")
		return
	}

	var tokenData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		log.Printf("[google] Failed to parse token response: %v", err)
		s.writeError(w, http.StatusInternalServerError, "failed to parse token response")
		return
	}

	if tokenData.AccessToken == "" {
		log.Printf("[google] No access token in response: %s", string(body))
		s.writeError(w, http.StatusBadGateway, "no access token received from Google")
		return
	}

	// Calculate expiry time. Default to 1 hour if not specified.
	expiresIn := tokenData.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	if err := db.SaveGoogleTokens(s.DB, tokenData.AccessToken, tokenData.RefreshToken, expiresAt); err != nil {
		log.Printf("[google] Failed to save tokens: %v", err)
		s.writeError(w, http.StatusInternalServerError, "failed to store tokens")
		return
	}

	log.Printf("[google] OAuth tokens stored. The Nest is under our control now.")
	http.Redirect(w, r, "/#google-connected", http.StatusTemporaryRedirect)
}

// HandleGoogleStatus returns whether Google Nest is connected (tokens exist)
// and whether the server is configured to do Google OAuth at all (has the
// client ID + project ID baked in or set via env).
//
// GET /api/google/status
func (s *Server) HandleGoogleStatus(w http.ResponseWriter, r *http.Request) {
	configured := s.GoogleClientID != "" && s.GoogleProjectID != ""

	tokens, err := db.LoadGoogleTokens(s.DB)
	if err != nil || tokens == nil {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"connected":  false,
			"configured": configured,
		})
		return
	}

	// Check if the token has expired
	expired := time.Now().After(tokens.ExpiresAt)

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"connected":     true,
		"configured":    configured,
		"token_expired": expired,
		"has_refresh":   tokens.RefreshToken != "",
	})
}

// HandleGoogleDisconnect removes stored Google OAuth tokens.
// Clean break. No evidence left behind.
//
// DELETE /api/google/disconnect
func (s *Server) HandleGoogleDisconnect(w http.ResponseWriter, r *http.Request) {
	if err := db.DeleteGoogleTokens(s.DB); err != nil {
		log.Printf("[google] Failed to delete tokens: %v", err)
		s.writeError(w, http.StatusInternalServerError, "failed to remove Google tokens")
		return
	}
	log.Printf("[google] Google Nest disconnected. Tokens shredded.")
	s.writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Cached Google devices response — avoids hitting SDM API on every page load.
var cachedGoogleDevicesJSON []byte
var cachedGoogleDevicesTime time.Time
const googleDevicesCacheTTL = 60 * time.Second

// HandleGoogleDevices lists all Nest devices from the Smart Device Management API.
// Caches the response for 60 seconds to avoid repeated slow API calls.
//
// GET /api/google/devices
func (s *Server) HandleGoogleDevices(w http.ResponseWriter, r *http.Request) {
	// Return cached response if fresh
	if len(cachedGoogleDevicesJSON) > 0 && time.Since(cachedGoogleDevicesTime) < googleDevicesCacheTTL {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(cachedGoogleDevicesJSON)
		return
	}
	tokens, err := db.LoadGoogleTokens(s.DB)
	if err != nil || tokens == nil {
		s.writeError(w, http.StatusUnauthorized, "Google Nest not connected -- authorize first via /api/google/auth")
		return
	}

	// Check expiry and refresh if needed
	accessToken := tokens.AccessToken
	if time.Now().After(tokens.ExpiresAt) {
		refreshed, err := s.refreshGoogleToken(tokens.RefreshToken)
		if err != nil {
			log.Printf("[google] Token refresh failed: %v", err)
			s.writeError(w, http.StatusUnauthorized, "token expired and refresh failed -- re-authorize via /api/google/auth")
			return
		}
		accessToken = refreshed
	}

	// Call the SDM API to list devices
	sdmURL := fmt.Sprintf("https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices", s.GoogleProjectID)
	req, _ := http.NewRequest("GET", sdmURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[google] SDM API request failed: %v", err)
		s.writeError(w, http.StatusBadGateway, "failed to reach Google SDM API")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 65536))

	if resp.StatusCode != http.StatusOK {
		log.Printf("[google] SDM API returned %d: %s", resp.StatusCode, string(body))
		s.writeError(w, resp.StatusCode, "SDM API error: "+string(body))
		return
	}

	// Cache and return
	cachedGoogleDevicesJSON = body
	cachedGoogleDevicesTime = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// refreshGoogleToken uses the refresh token to obtain a new access token.
// Returns the new access token or an error. Updates the DB with new tokens.
// This is the kind of thing that keeps me up at night -- token rotation.
func (s *Server) refreshGoogleToken(refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", fmt.Errorf("no refresh token available")
	}

	resp, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", url.Values{
		"client_id":     {s.GoogleClientID},
		"client_secret": {s.GoogleClientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		return "", fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 16384))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("refresh returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return "", fmt.Errorf("failed to parse refresh response: %w", err)
	}

	if tokenData.AccessToken == "" {
		return "", fmt.Errorf("no access token in refresh response")
	}

	expiresIn := tokenData.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Google refresh token grants don't always return a new refresh token.
	// Keep the existing one if not provided.
	if err := db.SaveGoogleTokens(s.DB, tokenData.AccessToken, refreshToken, expiresAt); err != nil {
		log.Printf("[google] Failed to save refreshed tokens: %v", err)
		// Non-fatal -- we still have the new access token for this request
	}

	log.Printf("[google] Access token refreshed. The walls have ears, but we rotated the keys.")
	return tokenData.AccessToken, nil
}
