package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/coalson/haus/internal/db"
)

type deviceAuthRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// HandleDeviceAuth attempts to authenticate with a device using provided credentials.
// POST /api/devices/{ip}/auth
func (s *Server) HandleDeviceAuth(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	if ip == "" {
		s.writeError(w, http.StatusBadRequest, "ip required")
		return
	}

	var req deviceAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Try SunPower PVS auth: GET /auth?login with Basic auth
	username := req.Username
	if username == "" {
		username = "ssm_owner" // SunPower default
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	authURL := fmt.Sprintf("https://%s/auth?login", ip)
	authReq, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to create auth request")
		return
	}
	authReq.SetBasicAuth(username, req.Password)

	resp, err := client.Do(authReq)
	if err != nil {
		log.Printf("[device-auth] Failed to connect to %s: %v", ip, err)
		s.writeError(w, http.StatusBadGateway, "Could not reach device. Check the IP address.")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Printf("[device-auth] Auth failed for %s: HTTP %d — %s", ip, resp.StatusCode, string(body))
		s.writeError(w, http.StatusUnauthorized, "Authentication failed. Check your password.")
		return
	}

	// Look for session cookie
	var sessionToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session" {
			sessionToken = cookie.Value
			break
		}
	}

	if sessionToken == "" {
		log.Printf("[device-auth] No session cookie from %s", ip)
		s.writeError(w, http.StatusInternalServerError, "Authentication succeeded but no session was returned.")
		return
	}

	log.Printf("[device-auth] Successfully authenticated with %s", ip)

	// Persist credentials so the probe can use them next time
	if err := db.SaveDeviceCredential(s.DB, ip, "sunpower", username, req.Password, sessionToken); err != nil {
		log.Printf("[device-auth] Failed to save credentials for %s: %v", ip, err)
	}

	// Now try to fetch live data with the session
	varsURL := fmt.Sprintf("https://%s/vars?match=livedata&fmt=obj", ip)
	varsReq, _ := http.NewRequest("GET", varsURL, nil)
	varsReq.AddCookie(&http.Cookie{Name: "session", Value: sessionToken})

	varsResp, err := client.Do(varsReq)
	if err != nil {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"ok":      true,
			"session": sessionToken,
			"message": "Authenticated successfully, but could not fetch data.",
		})
		return
	}
	defer varsResp.Body.Close()
	varsBody, _ := io.ReadAll(varsResp.Body)

	var data map[string]interface{}
	json.Unmarshal(varsBody, &data)

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"session": sessionToken,
		"message": "Connected to SunPower PVS!",
		"data":    data,
	})
}
