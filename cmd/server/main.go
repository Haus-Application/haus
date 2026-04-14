package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coalson/haus/internal/ai"
	"github.com/coalson/haus/internal/api"
	"github.com/coalson/haus/internal/db"
	"github.com/coalson/haus/internal/discovery"
	"github.com/coalson/haus/internal/hue"
	"github.com/coalson/haus/internal/kasa"
	"github.com/coalson/haus/internal/kb"
	"github.com/coalson/haus/internal/ws"
	"github.com/gorilla/websocket"
)

// defaultAPIKey is set at build time via -ldflags for release builds.
var defaultAPIKey string

func main() {
	log.Println("[haus] Mother is waking up... initializing Haus server.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "haus.db"
	}

	// Open database
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("[haus] Failed to open database: %v. Mother won't be happy about this.", err)
	}
	defer database.Close()

	// Create scanner (auto-detects subnet)
	scanner := discovery.NewScanner(database)

	// -----------------------------------------------------------------------
	// 1. WebSocket hub -- George Michael built this and I'm very proud of him.
	// -----------------------------------------------------------------------
	hub := ws.NewHub()
	go hub.Run()
	log.Println("[haus] WebSocket hub is running. George Michael did good work.")

	// -----------------------------------------------------------------------
	// 2. Warm-resume pollers from persisted state. Zero network I/O at boot —
	//    discovery happens when the user clicks Scan (POST /api/scan). Any
	//    previously-discovered devices still in the DB keep their pollers alive
	//    across restarts. Mother prefers we knock first.
	// -----------------------------------------------------------------------
	var (
		kasaPoller *kasa.Poller
		hueClient  *hue.Client
		huePoller  *hue.Poller
		started    []string
	)

	if ips, err := db.LoadKasaIPs(database); err != nil {
		log.Printf("[haus] Could not load Kasa devices from DB: %v", err)
	} else if len(ips) > 0 {
		kasaPoller = kasa.NewPoller(ips, &broadcasterAdapter{hub: hub})
		kasaPoller.Start()
		started = append(started, fmt.Sprintf("Kasa: polling %d known device(s)", len(ips)))
	}

	hueConfig, err := db.LoadHueConfig(database)
	if err == nil && hueConfig != nil {
		hueClient = hue.NewClient(hueConfig.BridgeIP, hueConfig.Username)
		huePoller = hue.NewPoller(hueClient, &broadcasterAdapter{hub: hub})
		huePoller.Start()
		started = append(started, fmt.Sprintf("Hue: connected to bridge at %s", hueConfig.BridgeIP))
	}

	// -----------------------------------------------------------------------
	// 4. AI concierge -- GOB's magic trick engine
	// -----------------------------------------------------------------------
	kasaFuncs := buildKasaFuncs(kasaPoller)
	hueFuncs := buildHueFuncs(hueClient, huePoller)
	var concierge *ai.Concierge
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		apiKey = defaultAPIKey
	}
	if apiKey != "" {
		concierge = ai.NewConcierge(apiKey, kasaFuncs, hueFuncs)
		concierge.HTTPQuery = buildHTTPQuery(database)
		concierge.CameraSnapshot = api.CaptureSnapshotBase64
		concierge.JellyFishQuery = buildJellyFishQuery()
		log.Println("[haus] AI concierge is ready. GOB says the magic is real this time.")
	} else {
		log.Println("[haus] ANTHROPIC_API_KEY not set -- AI concierge disabled. GOB can't perform without props.")
	}

	// -----------------------------------------------------------------------
	// 5. API server and routes
	// -----------------------------------------------------------------------
	// Load the device knowledge base (100 markdown files with YAML frontmatter).
	// Used by the probe system to inject rich API docs and by the validation
	// dashboard. Non-fatal if missing — we'll just fall back to docs/api/.
	kbDir := resolveRuntimePath("docs/devices")
	catalog, err := kb.Load(kbDir)
	if err != nil {
		log.Printf("[haus] WARNING: failed to load device knowledge base from %s: %v", kbDir, err)
	}

	server := &api.Server{
		DB:         database,
		Scanner:    scanner,
		KasaPoller: kasaPoller,
		HueClient:  hueClient,
		HuePoller:  huePoller,
		Concierge:  concierge,
		KB:         catalog,
		Hub:        hub,
		APIKey:     apiKey,
		ValidationDir: resolveRuntimePath("validation"),

		// Google Nest SDM credentials -- George Sr. keeps these locked up tight.
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleProjectID:    os.Getenv("GOOGLE_PROJECT_ID"),
	}

	mux := http.NewServeMux()

	// Scan routes
	mux.HandleFunc("GET /api/devices", server.HandleDevices)
	mux.HandleFunc("GET /api/devices/{ip}/probe", server.HandleProbeDevice)
	mux.HandleFunc("POST /api/devices/{ip}/jellyfish", server.HandleJellyfishCommand)
	mux.HandleFunc("POST /api/devices/{ip}/auth", server.HandleDeviceAuth)
	mux.HandleFunc("POST /api/scan", server.HandleStartScan)
	mux.HandleFunc("GET /api/scan/stream", server.HandleScanStream)
	mux.HandleFunc("GET /api/scan/results", server.HandleScanResults)

	// Kasa routes
	mux.HandleFunc("GET /api/kasa/devices", server.HandleKasaDevices)
	mux.HandleFunc("PUT /api/kasa/devices/{ip}/state", server.HandleKasaSetState)
	mux.HandleFunc("PUT /api/kasa/devices/{ip}/brightness", server.HandleKasaSetBrightness)
	mux.HandleFunc("PUT /api/kasa/devices/{ip}/fan-speed", server.HandleKasaSetFanSpeed)

	// Hue routes
	mux.HandleFunc("GET /api/hue/discover", server.HandleHueDiscover)
	mux.HandleFunc("POST /api/hue/pair", server.HandleHuePair)
	mux.HandleFunc("GET /api/hue/status", server.HandleHueStatus)
	mux.HandleFunc("DELETE /api/hue/disconnect", server.HandleHueDisconnect)
	mux.HandleFunc("GET /api/hue/lights", server.HandleHueLights)
	mux.HandleFunc("PUT /api/hue/lights/{id}", server.HandleHueSetLight)
	mux.HandleFunc("GET /api/hue/rooms", server.HandleHueRooms)
	mux.HandleFunc("PUT /api/hue/rooms/{id}", server.HandleHueSetRoom)
	mux.HandleFunc("GET /api/hue/scenes", server.HandleHueScenes)
	mux.HandleFunc("POST /api/hue/scenes/{id}/activate", server.HandleHueActivateScene)

	// Google Nest OAuth routes -- the front door, not the back.
	mux.HandleFunc("GET /api/google/auth", server.HandleGoogleAuthStart)
	mux.HandleFunc("GET /api/google/callback", server.HandleGoogleAuthCallback)
	mux.HandleFunc("GET /api/google/status", server.HandleGoogleStatus)
	mux.HandleFunc("DELETE /api/google/disconnect", server.HandleGoogleDisconnect)
	mux.HandleFunc("GET /api/google/devices", server.HandleGoogleDevices)
	// Camera routes (go2rtc proxy)
	mux.HandleFunc("GET /api/cameras", server.HandleCameraList)
	mux.HandleFunc("POST /api/cameras/{id}/webrtc", server.HandleCameraWebRTC)
	mux.HandleFunc("GET /api/cameras/{id}/stream", server.HandleCameraStream)
	mux.HandleFunc("GET /api/cameras/{id}/snapshot", server.HandleCameraSnapshot)

	mux.HandleFunc("POST /api/google/camera/{deviceID}/stream", server.HandleNestCameraStream)
	mux.HandleFunc("POST /api/google/camera/{deviceID}/extend", server.HandleNestCameraExtend)
	mux.HandleFunc("POST /api/google/camera/{deviceID}/stop", server.HandleNestCameraStop)

	// Chat route
	mux.HandleFunc("POST /api/chat", server.HandleChat)
	mux.HandleFunc("POST /api/chat/device", server.HandleDeviceChat)

	// Validation routes -- Buster & GOB grading the 100-device knowledge base.
	mux.HandleFunc("GET /api/validation/summary", server.HandleValidationSummary)
	mux.HandleFunc("GET /api/validation/devices", server.HandleValidationDeviceList)
	mux.HandleFunc("GET /api/validation/devices/{slug}", server.HandleValidationDeviceReport)
	mux.HandleFunc("POST /api/validation/run", server.HandleValidationRun)

	// WebSocket
	mux.HandleFunc("/api/ws", hub.HandleWebSocket)

	// Resolve runtime asset paths (frontend + docs) relative to either CWD (dev)
	// or the binary's own directory (release). Mother always said: know where you live.
	api.DocsBaseDir = resolveRuntimePath("docs/api")

	// Serve frontend (SPA) if dist directory exists
	frontendDir := resolveRuntimePath("frontend/dist")
	if info, err := os.Stat(frontendDir); err == nil && info.IsDir() {
		log.Printf("[haus] Serving frontend from %s", frontendDir)
		mux.Handle("/", spaHandler(frontendDir))
	} else {
		log.Println("[haus] WARNING: frontend/dist not found. Maeby hasn't built it yet. API-only mode.")
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok","message":"Haus API is running. Frontend not built yet."}`))
		})
	}

	// -----------------------------------------------------------------------
	// 6. Enrich Google/Nest device names from SDM API + go2rtc
	// -----------------------------------------------------------------------
	go enrichNestDeviceNames(database, server)

	// -----------------------------------------------------------------------
	// Startup summary -- only mention integrations that actually started.
	// -----------------------------------------------------------------------
	if len(started) == 0 {
		log.Println("[haus] No devices tracked yet. Click Scan in the UI when you're ready.")
	} else {
		for _, line := range started {
			log.Printf("[haus] %s", line)
		}
	}
	log.Printf("[haus] Mother's house is coming alive! Listening on :%s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("[haus] Server failed: %v. I PANICKED!", err)
	}
}

// broadcasterAdapter wraps the WebSocket Hub to satisfy the Broadcaster
// interfaces defined in the kasa and hue packages. Those interfaces expect
// BroadcastGlobal(event interface{}), but the Hub's method takes a concrete
// ws.BroadcastEvent. I had to study both signatures very carefully. Mother
// said that was "obsessive," but I call it "thorough."
type broadcasterAdapter struct {
	hub *ws.Hub
}

func (b *broadcasterAdapter) BroadcastGlobal(event interface{}) {
	// The kasa and hue pollers send their own BroadcastEvent structs.
	// We need to re-wrap them into ws.BroadcastEvent. Both packages
	// define BroadcastEvent with Type string + Payload interface{},
	// but Go doesn't do structural typing for structs, only interfaces.
	type typedEvent struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}

	switch e := event.(type) {
	case kasa.BroadcastEvent:
		b.hub.BroadcastGlobal(ws.BroadcastEvent{
			Type:    e.Type,
			Payload: e.Payload,
		})
	case hue.BroadcastEvent:
		b.hub.BroadcastGlobal(ws.BroadcastEvent{
			Type:    e.Type,
			Payload: e.Payload,
		})
	default:
		log.Printf("[ws] broadcasterAdapter: unknown event type %T -- I don't know what this is.", event)
	}
}

// buildKasaFuncs creates the closure struct that bridges the kasa package's
// concrete types to the ai package's generic KasaDeviceInfo. If the poller
// is nil (no devices found), we return nil and the concierge gracefully
// degrades. I know all of this because I've read every line of the Kasa
// XOR protocol spec. Twice.
// buildHTTPQuery creates a closure that makes authenticated HTTPS requests
// to devices using stored credentials from the DB.
// enrichNestDeviceNames queries go2rtc for camera streams and the Google SDM
// API for device names, then updates DB entries so Google devices show proper
// names like "Living Room Camera" instead of "Google .89".
func enrichNestDeviceNames(database *sql.DB, server *api.Server) {
	// Give go2rtc and Google API a moment to be ready
	time.Sleep(3 * time.Second)

	// Get Google Nest devices from SDM API
	client, err := server.GetGoogleClient()
	if err != nil {
		log.Printf("[enrich] Google not connected, skipping camera enrichment: %v", err)
		return
	}

	devices, err := client.ListDevices()
	if err != nil {
		log.Printf("[enrich] Failed to list Nest devices: %v", err)
		return
	}

	// Get all devices from DB that are Google/Nest (manufacturer = "Google")
	allDevices, err := db.LoadAllDevices(database)
	if err != nil {
		return
	}

	// Match Nest SDM devices to DB devices by name matching or just update
	// all unnamed Google devices with the Nest names
	googleDevices := []db.DeviceRow{}
	for _, d := range allDevices {
		if d.Manufacturer == "Google" && (d.Name == "" || strings.HasPrefix(d.Name, "Google .") || strings.HasPrefix(d.Name, "Device .")) {
			googleDevices = append(googleDevices, d)
		}
	}

	// Build list of Nest cameras/thermostats with their names
	nestNames := []struct{ Name, Type string }{}
	for _, dev := range devices {
		displayName := ""
		if len(dev.ParentRelations) > 0 {
			displayName = dev.ParentRelations[0].DisplayName
		}
		deviceType := strings.Split(dev.Type, ".")[len(strings.Split(dev.Type, "."))-1]
		if displayName != "" {
			nestNames = append(nestNames, struct{ Name, Type string }{displayName, deviceType})
		}
	}

	// Assign Nest names to unnamed Google devices round-robin
	// (best effort — we can't match IPs to SDM devices)
	updated := 0
	for i, gd := range googleDevices {
		if i < len(nestNames) {
			typeSuffix := strings.ToLower(nestNames[i].Type)
			if typeSuffix == "thermostat" { typeSuffix = "Thermostat" }
			if typeSuffix == "camera" { typeSuffix = "Camera" }
			if typeSuffix == "display" { typeSuffix = "Display" }
			if typeSuffix == "doorbell" { typeSuffix = "Doorbell" }
			newName := fmt.Sprintf("%s %s", nestNames[i].Name, typeSuffix)
			category := "smart_home"
			deviceType := "nest_" + strings.ToLower(nestNames[i].Type)
			if nestNames[i].Type == "CAMERA" || nestNames[i].Type == "DISPLAY" {
				deviceType = "nest_camera"
			}

			db.UpsertDevice(database, gd.IP, gd.MAC, gd.Hostname, newName,
				gd.Manufacturer, gd.Model, deviceType, category,
				gd.Protocols, gd.Services, gd.OpenPorts, gd.Metadata)
			log.Printf("[enrich] %s → %s (%s)", gd.IP, newName, nestNames[i].Type)
			updated++
		}
	}

	if updated > 0 {
		log.Printf("[enrich] Updated %d Google device names from Nest SDM API", updated)
	}
}

// buildJellyFishQuery creates a function that queries JellyFish controllers
// via their WebSocket API.
func buildJellyFishQuery() ai.JellyFishQueryFunc {
	return func(ip string, command map[string]interface{}) (string, error) {
		dialer := websocket.Dialer{HandshakeTimeout: 3 * time.Second}
		conn, _, err := dialer.Dial(fmt.Sprintf("ws://%s:9000/", ip), nil)
		if err != nil {
			return "", fmt.Errorf("connect failed: %w", err)
		}
		defer conn.Close()

		conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
		if err := conn.WriteJSON(command); err != nil {
			return "", fmt.Errorf("write failed: %w", err)
		}

		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("read failed: %w", err)
		}

		return string(msg), nil
	}
}

func buildHTTPQuery(database *sql.DB) ai.DeviceHTTPQuery {
	return func(ip, path string) (string, error) {
		cred, err := db.LoadDeviceCredential(database, ip)
		if err != nil || cred == nil {
			return "", fmt.Errorf("no stored credentials for %s", ip)
		}

		tr := &http.Transport{
			TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
			ForceAttemptHTTP2: false,
		}
		tlsClient := &http.Client{Timeout: 10 * time.Second, Transport: tr}

		// Step 1: Authenticate and extract session token from response body
		authReq, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/auth?login", ip), nil)
		authReq.SetBasicAuth(cred.Username, cred.Password)
		authResp, err := tlsClient.Do(authReq)
		if err != nil {
			return "", fmt.Errorf("auth failed: %w", err)
		}
		authBody, _ := io.ReadAll(authResp.Body)
		authResp.Body.Close()

		if authResp.StatusCode != 200 {
			return "", fmt.Errorf("authentication failed (HTTP %d): %s", authResp.StatusCode, string(authBody))
		}

		// Extract session from response body (PVS returns {"session": "..."})
		// AND from Set-Cookie header — try both
		var session string
		var authJSON struct{ Session string `json:"session"` }
		if json.Unmarshal(authBody, &authJSON) == nil && authJSON.Session != "" {
			session = strings.TrimSpace(authJSON.Session)
		}
		if session == "" {
			for _, cookie := range authResp.Cookies() {
				if cookie.Name == "session" {
					session = cookie.Value
					break
				}
			}
		}
		if session == "" {
			return "", fmt.Errorf("authentication succeeded but no session token returned")
		}

		// Step 2: Make the data request with the session cookie
		dataReq, _ := http.NewRequest("GET", fmt.Sprintf("https://%s%s", ip, path), nil)
		dataReq.AddCookie(&http.Cookie{Name: "session", Value: session})
		dataResp, err := tlsClient.Do(dataReq)
		if err != nil {
			return "", fmt.Errorf("request failed: %w", err)
		}
		defer dataResp.Body.Close()

		body, _ := io.ReadAll(io.LimitReader(dataResp.Body, 32768))

		if dataResp.StatusCode == 403 || dataResp.StatusCode == 400 {
			return "", fmt.Errorf("access denied (HTTP %d): %s", dataResp.StatusCode, string(body))
		}

		return string(body), nil
	}
}

func buildKasaFuncs(poller *kasa.Poller) *ai.KasaFuncs {
	if poller == nil {
		return nil
	}
	return &ai.KasaFuncs{
		ListDevices: func() ([]ai.KasaDeviceInfo, error) {
			devices := poller.GetDevices()
			var infos []ai.KasaDeviceInfo
			for _, d := range devices {
				infos = append(infos, ai.KasaDeviceInfo{
					IP:         d.IP,
					Alias:      d.Alias,
					Model:      d.Model,
					DeviceType: d.DeviceType,
					On:         d.On,
					Brightness: d.Brightness,
					FanSpeed:   d.FanSpeed,
				})
			}
			return infos, nil
		},
		QueryDevice: func(ip string) (*ai.KasaDeviceInfo, error) {
			dev, err := kasa.QueryDevice(ip)
			if err != nil {
				return nil, err
			}
			return &ai.KasaDeviceInfo{
				IP: dev.IP, Alias: dev.Alias, Model: dev.Model,
				DeviceType: dev.DeviceType, On: dev.On,
				Brightness: dev.Brightness, FanSpeed: dev.FanSpeed,
			}, nil
		},
		SetState: func(ip string, on bool) error {
			return kasa.SetState(ip, on)
		},
		SetBrightness: func(ip string, brightness int) error {
			return kasa.SetBrightness(ip, brightness)
		},
		SetFanSpeed: func(ip string, speed int) error {
			return kasa.SetFanSpeed(ip, speed)
		},
	}
}

// buildHueFuncs creates the closure struct that bridges the hue package's
// concrete types to the ai package's generic Hue info types. The Hue API v2
// uses CIE xy color space, which I find very soothing to think about. If
// client or poller is nil, we return nil.
func buildHueFuncs(client *hue.Client, poller *hue.Poller) *ai.HueFuncs {
	if client == nil || poller == nil {
		return nil
	}
	return &ai.HueFuncs{
		ListLights: func() ([]ai.HueLightInfo, error) {
			lights := poller.GetLights()
			var infos []ai.HueLightInfo
			for _, l := range lights {
				infos = append(infos, ai.HueLightInfo{
					ID:         l.ID,
					Name:       l.Name,
					On:         l.On,
					Brightness: l.Brightness,
					RoomName:   l.RoomName,
					Reachable:  true, // Hue API v2 doesn't expose reachable directly on the light resource
				})
			}
			return infos, nil
		},
		ListRooms: func() ([]ai.HueRoomInfo, error) {
			rooms := poller.GetRooms()
			var infos []ai.HueRoomInfo
			for _, r := range rooms {
				anyOn := false
				for _, l := range r.Lights {
					if l.On {
						anyOn = true
						break
					}
				}
				infos = append(infos, ai.HueRoomInfo{
					ID:             r.ID,
					Name:           r.Name,
					GroupedLightID: r.GroupedLightID,
					LightCount:     len(r.Lights),
					AnyOn:          anyOn,
				})
			}
			return infos, nil
		},
		ListScenes: func() ([]ai.HueSceneInfo, error) {
			scenes := poller.GetScenes()
			var infos []ai.HueSceneInfo
			for _, s := range scenes {
				infos = append(infos, ai.HueSceneInfo{
					ID:       s.ID,
					Name:     s.Name,
					RoomName: s.RoomName,
				})
			}
			return infos, nil
		},
		ToggleLight: func(lightID string, on bool) error {
			return client.SetLightState(lightID, &on, nil, nil)
		},
		SetBrightness: func(lightID string, brightness float64) error {
			return client.SetLightState(lightID, nil, &brightness, nil)
		},
		SetColor: func(lightID string, xy [2]float64) error {
			return client.SetLightState(lightID, nil, nil, &xy)
		},
		SetRoomState: func(groupedLightID string, on *bool, brightness *float64) error {
			return client.SetGroupedLightState(groupedLightID, on, brightness)
		},
		ActivateScene: func(sceneID string) error {
			return client.ActivateScene(sceneID)
		},
	}
}

// resolveRuntimePath locates a runtime asset directory. Tries CWD first (dev
// mode — running from source tree), then falls back to a path relative to the
// binary's own directory (release mode — ./run.sh or binary invoked from
// elsewhere). Returns the input unchanged if neither location works; callers
// that use os.Stat will then surface the "not found" branch.
func resolveRuntimePath(relative string) string {
	if info, err := os.Stat(relative); err == nil && info.IsDir() {
		abs, err := filepath.Abs(relative)
		if err == nil {
			return abs
		}
		return relative
	}
	exe, err := os.Executable()
	if err != nil {
		return relative
	}
	return filepath.Join(filepath.Dir(exe), relative)
}

// spaHandler serves static files from the given directory, falling back
// to index.html for unmatched routes (SPA behavior).
func spaHandler(dir string) http.Handler {
	fs := http.Dir(dir)
	fileServer := http.FileServer(fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := filepath.Clean(r.URL.Path)
		if path == "/" {
			path = "/index.html"
		}

		// Check if the file exists
		fullPath := filepath.Join(dir, path)
		if _, err := os.Stat(fullPath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for unmatched routes
		indexPath := filepath.Join(dir, "index.html")
		if _, err := os.Stat(indexPath); err != nil {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, indexPath)
	})
}
