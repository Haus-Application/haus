package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coalson/haus/internal/db"
	"github.com/coalson/haus/internal/kasa"
	"github.com/gorilla/websocket"
)

// loadAPIDocs reads the markdown documentation for a device integration type.
func loadAPIDocs(integration string) string {
	if integration == "" || integration == "unknown" || integration == "generic" {
		return ""
	}
	data, err := os.ReadFile(fmt.Sprintf("docs/api/%s.md", integration))
	if err != nil {
		return ""
	}
	return string(data)
}

// WebFingerprint captures what we found when probing a specific port on a device.
// Mother always said you can tell a lot about a device by its HTTP headers.
type WebFingerprint struct {
	Port       int    `json:"port"`
	URL        string `json:"url"`
	Title      string `json:"title,omitempty"`
	Server     string `json:"server,omitempty"`
	HasLogin   bool   `json:"has_login"`
	StatusCode int    `json:"status_code"`
	Body       string `json:"body,omitempty"` // first 200 chars for identification
}

// DeviceProbeResult is the real-time capability report for a device.
type DeviceProbeResult struct {
	IP           string                 `json:"ip"`
	Name         string                 `json:"name"`
	Reachable    bool                   `json:"reachable"`
	Integration  string                 `json:"integration"`  // "kasa", "hue", "cast", "generic", "unknown", "jellyfish", "yamaha", "sunpower", "brilliant"
	Status       string                 `json:"status"`       // "connected", "needs_pairing", "unreachable", "read_only", "discovered", "needs_auth", "offline"
	Capabilities []string               `json:"capabilities"` // ["on_off", "brightness", "fan_speed", "color", "scenes"]
	State        map[string]interface{} `json:"state"`        // current live state
	Actions      []DeviceAction         `json:"actions"`      // available actions for the UI
	SetupNeeded  *SetupStep             `json:"setup_needed,omitempty"` // if pairing/auth is required
	Fingerprints []WebFingerprint       `json:"fingerprints,omitempty"` // web fingerprints from port probing
	API          *APIInfo               `json:"api,omitempty"`          // API documentation
	APIDocs      string                 `json:"api_docs,omitempty"`     // full markdown documentation
}

// APIInfo describes the device's API protocol and available endpoints.
type APIInfo struct {
	Protocol    string       `json:"protocol"`    // "tcp_xor", "https_rest", "websocket_json", "http_rest"
	Port        int          `json:"port"`
	Description string       `json:"description"` // human-readable summary
	AuthMethod  string       `json:"auth_method,omitempty"` // "none", "link_button", "basic_auth", "api_key"
	DocURL      string       `json:"doc_url,omitempty"`     // link to manufacturer API docs
	Endpoints   []APIEndpoint `json:"endpoints,omitempty"`
}

// APIEndpoint describes a single API operation.
type APIEndpoint struct {
	Method      string `json:"method"`      // "TCP", "GET", "PUT", "POST", "WS_SEND"
	Path        string `json:"path"`        // endpoint path or command
	Description string `json:"description"`
	Example     string `json:"example,omitempty"` // example payload
}

// DeviceAction describes a control the UI can render.
type DeviceAction struct {
	ID          string                 `json:"id"`          // "toggle", "brightness", "fan_speed", "pair"
	Label       string                 `json:"label"`       // "Power", "Brightness", "Fan Speed"
	Type        string                 `json:"type"`        // "toggle", "slider", "buttons", "button"
	Value       interface{}            `json:"value"`       // current value
	Min         int                    `json:"min,omitempty"`
	Max         int                    `json:"max,omitempty"`
	Options     []ActionOption         `json:"options,omitempty"` // for button groups
}

// ActionOption is a single option in a button group.
type ActionOption struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// SetupStep describes what the user needs to do to connect.
type SetupStep struct {
	Type        string `json:"type"`        // "link_button", "auth", "info"
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`      // API endpoint to call
	ActionLabel string `json:"action_label"` // button text
}

// HandleProbeDevice does a real-time probe of a device to discover its
// capabilities, current state, and any setup needed.
// GET /api/devices/{ip}/probe
func (s *Server) HandleProbeDevice(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	if ip == "" {
		s.writeError(w, http.StatusBadRequest, "ip required")
		return
	}

	// Load stored info about this device
	rows, err := db.LoadAllDevices(s.DB)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to load device")
		return
	}

	var stored *db.DeviceRow
	for _, d := range rows {
		if d.IP == ip {
			stored = &d
			break
		}
	}

	result := &DeviceProbeResult{
		IP:    ip,
		State: make(map[string]interface{}),
	}

	if stored != nil {
		result.Name = stored.Name
	}

	// Determine integration based on stored data
	var protocols []string
	if stored != nil {
		json.Unmarshal([]byte(stored.Protocols), &protocols)
	}

	hasProtocol := func(p string) bool {
		for _, proto := range protocols {
			if proto == p {
				return true
			}
		}
		return false
	}

	var services []string
	if stored != nil {
		json.Unmarshal([]byte(stored.Services), &services)
	}

	hasService := func(s string) bool {
		for _, svc := range services {
			if svc == s {
				return true
			}
		}
		return false
	}

	deviceType := ""
	if stored != nil {
		deviceType = stored.DeviceType
	}

	manufacturer := ""
	if stored != nil {
		manufacturer = stored.Manufacturer
	}

	switch {
	case hasProtocol("kasa") || deviceType == "dimmer" || deviceType == "switch" || deviceType == "fan":
		probeKasa(result, ip)
	case deviceType == "hue_bridge":
		s.probeHue(result, ip)
	case hasProtocol("cast"):
		probeCast(result, ip)
	case deviceType == "jellyfish" || hasService("_jellyfishV2._tcp"):
		probeJellyFish(result, ip)
	case deviceType == "solar_gateway" || manufacturer == "SunPower" || hasService("_pvs6._tcp"):
		s.probeSunPower(result, ip)
	case deviceType == "av_receiver" || manufacturer == "Yamaha":
		probeYamaha(result, ip)
	case manufacturer == "Brilliant" || hasService("_brilliant._tcp"):
		probeBrilliant(result, ip)
	default:
		s.probeGeneric(result, ip, stored)
	}

	// Load API documentation markdown if available
	result.APIDocs = loadAPIDocs(result.Integration)

	log.Printf("[probe] %s (%s): integration=%s status=%s caps=%v",
		ip, result.Name, result.Integration, result.Status, result.Capabilities)

	s.writeJSON(w, http.StatusOK, result)
}

func probeKasa(result *DeviceProbeResult, ip string) {
	result.Integration = "kasa"

	dev, err := kasa.QueryDevice(ip)
	if err != nil {
		result.Reachable = false
		result.Status = "unreachable"
		return
	}

	result.Reachable = true
	result.Status = "connected"
	result.Name = dev.Alias

	result.State["on"] = dev.On
	result.State["brightness"] = dev.Brightness
	result.State["fan_speed"] = dev.FanSpeed
	result.State["model"] = dev.Model
	result.State["device_type"] = dev.DeviceType

	result.API = &APIInfo{
		Protocol:    "tcp_xor",
		Port:        9999,
		Description: "TP-Link Kasa local protocol. XOR-encrypted JSON over TCP. No authentication required.",
		AuthMethod:  "none",
		Endpoints: []APIEndpoint{
			{Method: "TCP", Path: `{"system":{"get_sysinfo":{}}}`, Description: "Query device state (on/off, brightness, model)"},
			{Method: "TCP", Path: `{"system":{"set_relay_state":{"state":1}}}`, Description: "Turn device on (state=1) or off (state=0)"},
		},
	}

	if dev.DeviceType == "dimmer" {
		result.API.Endpoints = append(result.API.Endpoints,
			APIEndpoint{Method: "TCP", Path: `{"smartlife.iot.dimmer":{"set_brightness":{"brightness":N}}}`, Description: "Set brightness 0-100%"},
		)
	}

	// Always has on/off
	result.Capabilities = []string{"on_off"}
	result.Actions = []DeviceAction{
		{ID: "toggle", Label: "Power", Type: "toggle", Value: dev.On},
	}

	switch dev.DeviceType {
	case "dimmer":
		result.Capabilities = append(result.Capabilities, "brightness")
		result.Actions = append(result.Actions, DeviceAction{
			ID: "brightness", Label: "Brightness", Type: "slider",
			Value: dev.Brightness, Min: 0, Max: 100,
		})
	case "fan":
		result.Capabilities = append(result.Capabilities, "fan_speed")
		result.Actions = append(result.Actions, DeviceAction{
			ID: "fan_speed", Label: "Fan Speed", Type: "buttons",
			Value: dev.FanSpeed,
			Options: []ActionOption{
				{Value: 1, Label: "Low"},
				{Value: 2, Label: "Med"},
				{Value: 3, Label: "High"},
				{Value: 4, Label: "Max"},
			},
		})
	}
}

func (s *Server) probeHue(result *DeviceProbeResult, ip string) {
	result.Integration = "hue"
	result.API = &APIInfo{
		Protocol:    "https_rest",
		Port:        443,
		Description: "Philips Hue API v2 (CLIP). HTTPS with self-signed certificate. Requires link-button pairing to obtain API key.",
		AuthMethod:  "link_button",
		DocURL:      "https://developers.meethue.com/develop/hue-api-v2/",
		Endpoints: []APIEndpoint{
			{Method: "GET", Path: "/clip/v2/resource/light", Description: "List all lights with on/off, brightness, color state"},
			{Method: "PUT", Path: "/clip/v2/resource/light/{id}", Description: "Control a light (on/off, brightness 0-100, color XY)", Example: `{"on":{"on":true},"dimming":{"brightness":75}}`},
			{Method: "GET", Path: "/clip/v2/resource/room", Description: "List rooms with grouped lights"},
			{Method: "PUT", Path: "/clip/v2/resource/grouped_light/{id}", Description: "Control all lights in a room"},
			{Method: "GET", Path: "/clip/v2/resource/scene", Description: "List saved scenes"},
			{Method: "PUT", Path: "/clip/v2/resource/scene/{id}", Description: "Activate a scene", Example: `{"recall":{"action":"active"}}`},
		},
	}

	// Check if we're already paired
	cfg, err := db.LoadHueConfig(s.DB)
	if err == nil && cfg != nil && cfg.BridgeIP == ip {
		// Already paired — show full capabilities
		result.Reachable = true
		result.Status = "connected"
		result.Capabilities = []string{"lights", "rooms", "scenes", "color", "brightness"}

		// Get live state from poller or direct query
		if s.HuePoller != nil {
			lights := s.HuePoller.GetLights()
			rooms := s.HuePoller.GetRooms()
			scenes := s.HuePoller.GetScenes()
			result.State["light_count"] = len(lights)
			result.State["room_count"] = len(rooms)
			result.State["scene_count"] = len(scenes)

			lightsOn := 0
			for _, l := range lights {
				if l.On {
					lightsOn++
				}
			}
			result.State["lights_on"] = lightsOn
		}
		return
	}

	// Not paired — check if bridge is reachable
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://%s/api/0/config", ip))
	if err != nil {
		result.Reachable = false
		result.Status = "unreachable"
		return
	}
	defer resp.Body.Close()

	result.Reachable = true
	result.Status = "needs_pairing"
	result.Capabilities = []string{"lights", "rooms", "scenes", "color", "brightness"}
	result.SetupNeeded = &SetupStep{
		Type:        "link_button",
		Title:       "Pair with Hue Bridge",
		Description: "Press the link button on your Hue bridge, then click Pair below.",
		Action:      fmt.Sprintf("/api/hue/pair"),
		ActionLabel: "Pair",
	}
}

func probeCast(result *DeviceProbeResult, ip string) {
	result.Integration = "cast"
	result.API = &APIInfo{
		Protocol:    "http_rest",
		Port:        8008,
		Description: "Google Cast local device info API. HTTP on port 8008. Read-only — control requires the Cast SDK.",
		AuthMethod:  "none",
		DocURL:      "https://developers.google.com/cast",
		Endpoints: []APIEndpoint{
			{Method: "GET", Path: "/setup/eureka_info", Description: "Device info (name, model, build info)"},
			{Method: "GET", Path: "/setup/configured_networks", Description: "WiFi network configuration"},
		},
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s:8008/setup/eureka_info", ip))
	if err != nil {
		result.Reachable = false
		result.Status = "unreachable"
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info struct {
		Name string `json:"name"`
	}
	json.Unmarshal(body, &info)

	result.Reachable = true
	result.Status = "read_only"
	result.Name = info.Name
	result.Capabilities = []string{"media_info"}
	result.State["name"] = info.Name
}

// probeGeneric is the aggressive device fingerprinter. It doesn't just knock on the
// door — it peers through every window, reads every nameplate, and catalogs
// everything it finds. Mother would say I'm being thorough. I say I'm being Buster.
func (s *Server) probeGeneric(result *DeviceProbeResult, ip string, stored *db.DeviceRow) {
	result.Integration = "unknown"

	// Gather stored data: open ports, services, metadata
	var storedPorts []int
	var services []string
	if stored != nil {
		json.Unmarshal([]byte(stored.OpenPorts), &storedPorts)
		json.Unmarshal([]byte(stored.Services), &services)
		result.State["manufacturer"] = stored.Manufacturer
		result.State["model"] = stored.Model
		result.State["category"] = stored.Category
	}

	// Build the complete port list: stored ports + common probe ports.
	// Like Mother's guest list — you invite everyone and see who shows up.
	commonPorts := []int{80, 443, 8080, 8443, 9000, 9999, 5455, 1883, 8008}
	portSet := make(map[int]bool)
	for _, p := range storedPorts {
		portSet[p] = true
	}
	for _, p := range commonPorts {
		portSet[p] = true
	}

	// Fingerprint every port — this is the part where I get excited
	var fingerprints []WebFingerprint
	for port := range portSet {
		fp := fingerprintPort(ip, port)
		if fp != nil {
			fingerprints = append(fingerprints, *fp)
		}
	}

	// Helper closures for signature matching
	hasService := func(svc string) bool {
		for _, s := range services {
			if strings.Contains(s, svc) {
				return true
			}
		}
		return false
	}

	titleContains := func(substr string) bool {
		for _, fp := range fingerprints {
			if strings.Contains(strings.ToLower(fp.Title), strings.ToLower(substr)) {
				return true
			}
		}
		return false
	}

	serverContains := func(substr string) bool {
		for _, fp := range fingerprints {
			if strings.Contains(fp.Server, substr) {
				return true
			}
		}
		return false
	}

	portOpen := func(p int) bool {
		for _, fp := range fingerprints {
			if fp.Port == p {
				return true
			}
		}
		return false
	}

	// Stash fingerprints on result for downstream use
	result.Fingerprints = fingerprints

	// Check known device signatures — like identifying a Bluth by their
	// particular brand of dysfunction
	if hasService("_jellyfishV2._tcp") || (titleContains("Software Update") && portOpen(9000)) {
		probeJellyFish(result, ip)
		return
	}
	if serverContains("Network_Module") || serverContains("RX-V") {
		probeYamaha(result, ip)
		return
	}
	if hasService("_pvs6._tcp") {
		s.probeSunPower(result, ip)
		return
	}
	if hasService("_brilliant._tcp") {
		probeBrilliant(result, ip)
		return
	}

	// No known signature matched. Did anything respond at all?
	if len(fingerprints) > 0 {
		result.Reachable = true
		result.Status = "discovered"

		// If any fingerprint found a login form, this device wants credentials
		hasAuth := false
		for _, fp := range fingerprints {
			if fp.HasLogin {
				hasAuth = true
				break
			}
		}
		if hasAuth {
			result.Status = "needs_auth"
			result.SetupNeeded = &SetupStep{
				Type:        "auth",
				Title:       "Authentication Required",
				Description: "This device has a web interface that requires login credentials.",
				Action:      fmt.Sprintf("/api/devices/%s/auth", ip),
				ActionLabel: "Configure",
			}
		}

		// Populate stored metadata into state
		var allPorts []int
		for _, fp := range fingerprints {
			allPorts = append(allPorts, fp.Port)
		}
		result.State["open_ports"] = allPorts
		return
	}

	// Nothing responded at all — the device has gone to the banana stand in the sky
	result.Reachable = false
	result.Status = "offline"
}

// probeJellyFish identifies and probes JellyFish LED controllers.
// These are the lights Mother would never let me install on the model home,
// but they have a WebSocket API on port 9000 and a firmware UI on 8080.
func probeJellyFish(result *DeviceProbeResult, ip string) {
	result.Integration = "jellyfish"
	result.Capabilities = []string{"zones", "patterns", "on_off"}
	result.SetupNeeded = nil
	result.Name = "JellyFish Lights"
	result.API = &APIInfo{
		Protocol:    "websocket_json",
		Port:        9000,
		Description: "JellyFish Lighting async WebSocket API. JSON commands over ws://host:9000/. No authentication required.",
		AuthMethod:  "none",
		Endpoints: []APIEndpoint{
			{Method: "WS_SEND", Path: `{"cmd":"toCtlrGet","get":[["zones"]]}`, Description: "Get all lighting zones"},
			{Method: "WS_SEND", Path: `{"cmd":"toCtlrGet","get":[["patternFileList"]]}`, Description: "List available patterns"},
			{Method: "WS_SEND", Path: `{"cmd":"toCtlrGet","get":[["runPattern","Zone1"]]}`, Description: "Get current state of a zone"},
			{Method: "WS_SEND", Path: `{"cmd":"toCtlrSet","runPattern":{"state":1,"zoneName":["Zone1"],"file":"Accent/White"}}`, Description: "Play a pattern on zones (state=1=on, 0=off)"},
			{Method: "WS_SEND", Path: `{"cmd":"toCtlrGet","get":[["ctlrName"]]}`, Description: "Get controller name"},
		},
	}

	// Connect to WebSocket API on port 9000
	dialer := websocket.Dialer{HandshakeTimeout: 3 * time.Second}
	wsURL := fmt.Sprintf("ws://%s:9000/", ip)
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		result.Reachable = false
		result.Status = "offline"
		log.Printf("[probe] JellyFish at %s: WebSocket connect failed: %v", ip, err)
		return
	}
	defer conn.Close()

	result.Reachable = true
	result.Status = "connected"

	// Query zones
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	conn.WriteJSON(map[string]any{"cmd": "toCtlrGet", "get": [][]string{{"zones"}}})
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, zonesMsg, err := conn.ReadMessage()
	if err == nil {
		var zonesResp map[string]interface{}
		if json.Unmarshal(zonesMsg, &zonesResp) == nil {
			if zones, ok := zonesResp["zones"]; ok {
				result.State["zones"] = zones
				if zm, ok := zones.(map[string]interface{}); ok {
					var zoneNames []string
					for name := range zm {
						zoneNames = append(zoneNames, name)
					}
					result.State["zone_names"] = zoneNames
					result.State["zone_count"] = len(zoneNames)
				}
			}
		}
	}

	// Query controller name
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	conn.WriteJSON(map[string]any{"cmd": "toCtlrGet", "get": [][]string{{"ctlrName"}}})
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, nameMsg, err := conn.ReadMessage()
	if err == nil {
		var nameResp map[string]interface{}
		if json.Unmarshal(nameMsg, &nameResp) == nil {
			if name, ok := nameResp["ctlrName"].(string); ok && name != "" {
				result.Name = name
			}
		}
	}

	// Query patterns
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	conn.WriteJSON(map[string]any{"cmd": "toCtlrGet", "get": [][]string{{"patternFileList"}}})
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, patternsMsg, err := conn.ReadMessage()
	if err == nil {
		var patternsResp map[string]interface{}
		if json.Unmarshal(patternsMsg, &patternsResp) == nil {
			if patterns, ok := patternsResp["patternFileList"]; ok {
				if pl, ok := patterns.([]interface{}); ok {
					type patternInfo struct {
						Path   string `json:"path"`
						Name   string `json:"name"`
						Folder string `json:"folder"`
					}
					var patternList []patternInfo
					for _, p := range pl {
						if pm, ok := p.(map[string]interface{}); ok {
							name, _ := pm["name"].(string)
							folder, _ := pm["folders"].(string)
							if name == "" {
								continue
							}
							path := name
							if folder != "" {
								path = folder + "/" + name
							}
							patternList = append(patternList, patternInfo{Path: path, Name: name, Folder: folder})
						}
					}
					result.State["pattern_count"] = len(patternList)
					result.State["patterns"] = patternList
				}
			}
		}
	}

	// Add actions
	result.Actions = []DeviceAction{
		{ID: "toggle", Label: "Power", Type: "toggle", Value: false},
	}

	// Add fingerprint for firmware UI on 8080
	fp := fingerprintPort(ip, 8080)
	if fp != nil {
		result.Fingerprints = append(result.Fingerprints, *fp)
	}
}

// probeYamaha queries Yamaha AV receivers via their Extended Control API.
// The Yamaha API is surprisingly well-documented, unlike the instructions
// Mother left for the microwave.
func probeYamaha(result *DeviceProbeResult, ip string) {
	result.Integration = "yamaha"
	result.API = &APIInfo{
		Protocol:    "http_rest",
		Port:        80,
		Description: "Yamaha MusicCast / Extended Control API. HTTP REST on port 80. No authentication required.",
		AuthMethod:  "none",
		DocURL:      "https://github.com/rsc-dev/pyamaha",
		Endpoints: []APIEndpoint{
			{Method: "GET", Path: "/YamahaExtendedControl/v1/system/getDeviceInfo", Description: "Device info (model, firmware, serial)"},
			{Method: "GET", Path: "/YamahaExtendedControl/v1/main/getStatus", Description: "Current state (power, volume, input, mute)"},
			{Method: "GET", Path: "/YamahaExtendedControl/v1/main/setPower?power=on", Description: "Power on/off (on, standby, toggle)"},
			{Method: "GET", Path: "/YamahaExtendedControl/v1/main/setVolume?volume=N", Description: "Set volume (0-161)"},
			{Method: "GET", Path: "/YamahaExtendedControl/v1/main/setInput?input=hdmi1", Description: "Switch input source"},
			{Method: "GET", Path: "/YamahaExtendedControl/v1/main/setMute?enable=true", Description: "Mute/unmute"},
		},
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s/YamahaExtendedControl/v1/system/getDeviceInfo", ip))
	if err != nil {
		result.Reachable = false
		result.Status = "unreachable"
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		result.Reachable = false
		result.Status = "unreachable"
		return
	}

	var info struct {
		ResponseCode int     `json:"response_code"`
		ModelName    string  `json:"model_name"`
		Destination  string  `json:"destination"`
		SystemVer    float64 `json:"system_version"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		result.Reachable = true
		result.Status = "discovered"
		result.State["raw_response"] = string(body)
		return
	}

	result.Reachable = true
	result.Status = "connected"
	result.Capabilities = []string{"power", "volume", "input_select", "mute"}
	result.State["model"] = info.ModelName
	result.State["system_version"] = fmt.Sprintf("%.1f", info.SystemVer)
	result.State["destination"] = info.Destination
	result.Name = info.ModelName

	// Query current power/volume state
	statusResp, err := client.Get(fmt.Sprintf("http://%s/YamahaExtendedControl/v1/main/getStatus", ip))
	if err == nil {
		defer statusResp.Body.Close()
		statusBody, _ := io.ReadAll(io.LimitReader(statusResp.Body, 4096))
		var status struct {
			Power  string `json:"power"`
			Volume int    `json:"volume"`
			Mute   bool   `json:"mute"`
			Input  string `json:"input"`
		}
		if json.Unmarshal(statusBody, &status) == nil {
			result.State["power"] = status.Power
			result.State["volume"] = status.Volume
			result.State["mute"] = status.Mute
			result.State["input"] = status.Input

			isOn := status.Power == "on"
			result.Actions = []DeviceAction{
				{ID: "toggle", Label: "Power", Type: "toggle", Value: isOn},
			}
		}
	}
}

// probeSunPower probes SunPower PVS solar monitoring systems.
// They use self-signed certs because apparently solar companies and the
// SEC have the same attitude toward proper certificates.
func (s *Server) probeSunPower(result *DeviceProbeResult, ip string) {
	result.Integration = "sunpower"
	result.Capabilities = []string{"solar_production", "solar_consumption", "grid_power", "battery"}
	result.API = &APIInfo{
		Protocol:    "https_rest",
		Port:        443,
		Description: "SunPower PVS local API. HTTPS with self-signed certificate. Requires Basic Auth (username: ssm_owner) to obtain a session cookie.",
		AuthMethod:  "basic_auth",
		Endpoints: []APIEndpoint{
			{Method: "GET", Path: "/auth?login", Description: "Authenticate with Basic Auth (ssm_owner:password), returns session cookie"},
			{Method: "GET", Path: "/vars?match=livedata&fmt=obj", Description: "Live solar production, house consumption, grid, battery data"},
			{Method: "GET", Path: "/vars?match=sys/devices&fmt=obj", Description: "All inverter/panel data with serial numbers, temps, production per panel"},
			{Method: "GET", Path: "/vars?match=sys&fmt=obj", Description: "Complete system data including all devices and live readings"},
		},
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// First check if the PVS is reachable at all
	resp, err := client.Get(fmt.Sprintf("https://%s/cgi-bin/dl_cgi?Command=DeviceList", ip))
	if err != nil {
		result.Reachable = false
		result.Status = "offline"
		return
	}
	defer resp.Body.Close()

	result.Reachable = true
	result.Name = "SunPower PVS"

	// Check for stored credentials
	cred, credErr := db.LoadDeviceCredential(s.DB, ip)
	if credErr == nil && cred != nil && cred.Password != "" {
		// Auth and extract session from JSON body (cookie jar doesn't work with self-signed TLS)
		authReq, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/auth?login", ip), nil)
		authReq.SetBasicAuth("ssm_owner", cred.Password)
		authResp, err := client.Do(authReq)
		if err == nil {
			authBody, _ := io.ReadAll(authResp.Body)
			authResp.Body.Close()
			if authResp.StatusCode == 200 {
				var authJSON struct{ Session string `json:"session"` }
				session := ""
				if json.Unmarshal(authBody, &authJSON) == nil {
					session = strings.TrimSpace(authJSON.Session)
				}
				if session != "" {
					varsReq, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/vars?match=sys&fmt=obj", ip), nil)
					varsReq.AddCookie(&http.Cookie{Name: "session", Value: session})
					varsResp, err := client.Do(varsReq)
					if err == nil {
						defer varsResp.Body.Close()
						if varsResp.StatusCode == 200 {
							body, _ := io.ReadAll(io.LimitReader(varsResp.Body, 65536))
							var data map[string]string
							if json.Unmarshal(body, &data) == nil {
								result.Status = "connected"
								result.SetupNeeded = nil
								parseSunPowerData(result, data)
								return
							}
						}
					}
				}
			}
		}
	}

	// No stored credentials or they failed — ask for password
	result.Status = "needs_auth"
	result.SetupNeeded = &SetupStep{
		Type:        "password",
		Title:       "Connect to SunPower PVS",
		Description: "Enter the installer password for your SunPower system. The username is always 'ssm_owner'. You can find this password in the SunPower app or on the sticker inside your PVS unit.",
		Action:      fmt.Sprintf("/api/devices/%s/auth", ip),
		ActionLabel: "Connect",
	}

	// Try unauthenticated access as fallback
	varsResp, err := client.Get(fmt.Sprintf("https://%s/vars?match=livedata&fmt=obj", ip))
	if err != nil {
		return
	}
	defer varsResp.Body.Close()

	if varsResp.StatusCode == 200 {
		result.Status = "connected"
		result.SetupNeeded = nil

		body, _ := io.ReadAll(io.LimitReader(varsResp.Body, 8192))
		var data map[string]string
		if json.Unmarshal(body, &data) == nil {
			parseSunPowerData(result, data)
		}
		return
	}

	// Got some other status — still worth reporting
	result.Status = "discovered"
	result.State["http_status"] = resp.StatusCode
}

// parseSunPowerData extracts solar metrics from the PVS vars response.
func parseSunPowerData(result *DeviceProbeResult, data map[string]string) {
	parseFloat := func(keys ...string) float64 {
		for _, key := range keys {
			if v, ok := data[key]; ok && v != "nan" && v != "" {
				var f float64
				fmt.Sscanf(v, "%f", &f)
				return f
			}
		}
		return 0
	}

	// Try both key formats — PVS firmware versions differ
	production := parseFloat("/sys/livedata/pv_p", "livedata.production.p_3phsum_kw")
	grid := parseFloat("/sys/livedata/net_p", "livedata.grid.p_3phsum_kw")
	consumption := parseFloat("/sys/livedata/site_load_p", "livedata.consumption.p_3phsum_kw")
	batterySOC := parseFloat("/sys/livedata/soc", "livedata.battery.soc_pct")
	batteryW := parseFloat("/sys/livedata/ess_p", "livedata.battery.p_3phsum_kw")
	lifetimeKWh := parseFloat("/sys/livedata/pv_en", "livedata.production.e_3phsum_kwh")

	result.State["production_kw"] = fmt.Sprintf("%.1f", production)
	result.State["grid_kw"] = fmt.Sprintf("%.1f", grid)
	result.State["consumption_kw"] = fmt.Sprintf("%.1f", consumption)
	result.State["battery_soc"] = fmt.Sprintf("%.0f", batterySOC)
	result.State["battery_kw"] = fmt.Sprintf("%.1f", batteryW)
	result.State["lifetime_kwh"] = fmt.Sprintf("%.0f", lifetimeKWh)
	result.State["exporting"] = grid < 0

	// Count inverters (each is a panel with micro-inverter)
	panelCount := 0
	for key := range data {
		if strings.HasPrefix(key, "/sys/devices/inverter/") && strings.HasSuffix(key, "/sn") {
			panelCount++
		}
	}
	if panelCount > 0 {
		result.State["panel_count"] = panelCount
	}

	log.Printf("[sunpower] Production: %.1fkW, Consumption: %.1fkW, Grid: %.1fkW, Lifetime: %.0fkWh, Panels: %d",
		production, consumption, grid, lifetimeKWh, panelCount)
}

// probeBrilliant checks for Brilliant smart home controllers.
// Their API is about as well-documented as the Bluth family's
// offshore accounts, so we just see what ports are open.
func probeBrilliant(result *DeviceProbeResult, ip string) {
	result.Integration = "brilliant"
	result.Reachable = true
	result.Status = "discovered"
	result.Capabilities = []string{"switch", "dimmer"}

	// Brilliant uses port 5455 for its control protocol
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:5455", ip), 2*time.Second)
	if err != nil {
		result.State["control_port"] = false
		log.Printf("[probe] Brilliant at %s: control port 5455 not reachable: %v", ip, err)
	} else {
		conn.Close()
		result.State["control_port"] = true
	}
}

// fingerprintPort attempts HTTP and HTTPS connections to a given port,
// capturing headers, title, login forms, and a body snippet.
// It's like reading someone's mail, but for network services.
func fingerprintPort(ip string, port int) *WebFingerprint {
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		// Don't follow redirects — we want to see the raw response
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Try HTTP first, then HTTPS
	schemes := []string{"http", "https"}
	// For well-known TLS ports, try HTTPS first
	if port == 443 || port == 8443 {
		schemes = []string{"https", "http"}
	}

	for _, scheme := range schemes {
		url := fmt.Sprintf("%s://%s:%d/", scheme, ip, port)
		resp, err := client.Get(url)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		resp.Body.Close()

		bodyStr := string(body)

		fp := &WebFingerprint{
			Port:       port,
			URL:        url,
			Server:     resp.Header.Get("Server"),
			StatusCode: resp.StatusCode,
			Title:      extractTitle(bodyStr),
			HasLogin:   hasLoginForm(bodyStr),
		}

		// Store first 200 chars of body for identification
		if len(bodyStr) > 200 {
			fp.Body = bodyStr[:200]
		} else {
			fp.Body = bodyStr
		}

		return fp
	}

	return nil
}

// extractTitle pulls the content of the first <title> tag from an HTML snippet.
func extractTitle(html string) string {
	lower := strings.ToLower(html)
	start := strings.Index(lower, "<title>")
	if start == -1 {
		return ""
	}
	start += len("<title>")
	end := strings.Index(lower[start:], "</title>")
	if end == -1 {
		// No closing tag — take what we have up to 100 chars
		remaining := html[start:]
		if len(remaining) > 100 {
			return remaining[:100]
		}
		return remaining
	}
	return html[start : start+end]
}

// hasLoginForm checks if an HTML body contains indicators of an authentication form.
func hasLoginForm(html string) bool {
	lower := strings.ToLower(html)
	indicators := []string{"password", "login", "sign in", "signin", "authenticate", "log in"}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			return true
		}
	}
	return false
}
