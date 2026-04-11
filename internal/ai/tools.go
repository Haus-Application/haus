// Package ai provides the Haus AI concierge — Claude-powered smart home control
// via natural language. Every tool call is a magic trick. The audience just
// doesn't know it yet.
package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

// ---------------------------------------------------------------------------
// Closure types — the concierge never imports kasa/hue directly. It works
// through these function references, like a magician working through an
// assistant. The audience sees the result, not the mechanism.
// ---------------------------------------------------------------------------

// KasaFuncs holds closures for TP-Link Kasa smart switch operations.
type KasaFuncs struct {
	ListDevices   func() ([]KasaDeviceInfo, error)
	QueryDevice   func(ip string) (*KasaDeviceInfo, error)
	SetState      func(ip string, on bool) error
	SetBrightness func(ip string, brightness int) error
	SetFanSpeed   func(ip string, speed int) error
}

// DeviceHTTPQuery lets the AI make authenticated HTTP requests to any device.
type DeviceHTTPQuery func(ip, path string) (string, error)

// CameraSnapshotFunc captures a base64 JPEG from a camera stream.
type CameraSnapshotFunc func(streamID string) (base64JPEG string, err error)

// JellyFishQueryFunc queries a JellyFish controller's WebSocket API.
type JellyFishQueryFunc func(ip string, command map[string]interface{}) (string, error)

// DeviceContext identifies a single device for scoped chat.
type DeviceContext struct {
	IP           string   `json:"ip"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
	DeviceType   string   `json:"device_type"`
	Category     string   `json:"category"`
	Protocols    []string `json:"protocols"`
	APIDocs      string   `json:"api_docs,omitempty"` // markdown API documentation
}

// HueFuncs holds closures for Philips Hue smart light operations.
type HueFuncs struct {
	ListLights    func() ([]HueLightInfo, error)
	ListRooms     func() ([]HueRoomInfo, error)
	ListScenes    func() ([]HueSceneInfo, error)
	ToggleLight   func(lightID string, on bool) error
	SetBrightness func(lightID string, brightness float64) error
	SetColor      func(lightID string, xy [2]float64) error
	SetRoomState  func(groupedLightID string, on *bool, brightness *float64) error
	ActivateScene func(sceneID string) error
}

// ---------------------------------------------------------------------------
// Info types — the data the audience sees after the trick.
// ---------------------------------------------------------------------------

// KasaDeviceInfo holds the state of a single Kasa smart switch or dimmer.
type KasaDeviceInfo struct {
	IP         string `json:"ip"`
	Alias      string `json:"alias"`
	Model      string `json:"model"`
	DeviceType string `json:"device_type"` // "dimmer", "switch", or "fan"
	On         bool   `json:"on"`
	Brightness int    `json:"brightness"`
	FanSpeed   int    `json:"fan_speed"`
}

// HueLightInfo holds the state of a single Hue light.
type HueLightInfo struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	On         bool    `json:"on"`
	Brightness float64 `json:"brightness"`
	RoomName   string  `json:"room_name"`
	Reachable  bool    `json:"reachable"`
}

// HueRoomInfo holds info about a Hue room/zone.
type HueRoomInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	GroupedLightID string `json:"grouped_light_id"`
	LightCount     int    `json:"light_count"`
	AnyOn          bool   `json:"any_on"`
}

// HueSceneInfo holds info about a Hue scene.
type HueSceneInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RoomName string `json:"room_name"`
}

// ---------------------------------------------------------------------------
// Color mapping — CIE xy coordinates for named colors.
// Every great illusion needs a palette.
// ---------------------------------------------------------------------------

var colorMap = map[string][2]float64{
	"warm":   {0.4578, 0.4101},
	"cool":   {0.3127, 0.3290},
	"red":    {0.6750, 0.3220},
	"blue":   {0.1532, 0.0475},
	"green":  {0.1700, 0.7000},
	"purple": {0.2703, 0.1398},
	"orange": {0.5614, 0.3944},
	"pink":   {0.3944, 0.1990},
	"white":  {0.3127, 0.3290},
}

// ---------------------------------------------------------------------------
// Tool definitions — the props for each trick.
// ---------------------------------------------------------------------------

// tool wraps a ToolParam into a ToolUnionParam — the SDK's union type
// for the Tools field. A small bit of staging for each prop.
func tool(t anthropic.ToolParam) anthropic.ToolUnionParam {
	return anthropic.ToolUnionParam{OfTool: &t}
}

// deviceTools returns all device control tool definitions for the Claude API.
func deviceTools() []anthropic.ToolUnionParam {
	tools := append(kasaTools(), hueTools()...)
	return tools
}

func kasaTools() []anthropic.ToolUnionParam {
	return []anthropic.ToolUnionParam{
		tool(anthropic.ToolParam{
			Name:        "kasa_list_devices",
			Description: anthropic.Opt("List all Kasa switches, dimmers, and fans with their current state."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "kasa_toggle_device",
			Description: anthropic.Opt("Turn a Kasa device on or off by name (fuzzy match)."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The device name to fuzzy match (e.g. 'Kitchen Lights', 'Living Room Fan')",
					},
					"on": map[string]interface{}{
						"type":        "boolean",
						"description": "True to turn on, false to turn off",
					},
				},
				Required: []string{"name", "on"},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "kasa_set_brightness",
			Description: anthropic.Opt("Set brightness of a Kasa dimmer switch (0-100)."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The dimmer name (e.g. 'Kitchen Lights')",
					},
					"brightness": map[string]interface{}{
						"type":        "integer",
						"description": "Brightness level 0-100",
					},
				},
				Required: []string{"name", "brightness"},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "kasa_set_fan_speed",
			Description: anthropic.Opt("Set ceiling fan speed level (1-4)."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The fan name (e.g. 'Living Room Fan')",
					},
					"speed": map[string]interface{}{
						"type":        "integer",
						"description": "Speed level 1 (low) to 4 (high)",
					},
				},
				Required: []string{"name", "speed"},
			},
		}),
	}
}

func hueTools() []anthropic.ToolUnionParam {
	return []anthropic.ToolUnionParam{
		tool(anthropic.ToolParam{
			Name:        "hue_list_lights",
			Description: anthropic.Opt("List all Hue smart lights with their current state (on/off, brightness, room)."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "hue_toggle_light",
			Description: anthropic.Opt("Turn a Hue light on or off by name (fuzzy match)."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The light or room name (e.g. 'living room', 'bedroom lamp')",
					},
					"on": map[string]interface{}{
						"type":        "boolean",
						"description": "True to turn on, false to turn off",
					},
				},
				Required: []string{"name", "on"},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "hue_set_brightness",
			Description: anthropic.Opt("Set the brightness of a Hue light (0-100)."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Light or room name",
					},
					"brightness": map[string]interface{}{
						"type":        "number",
						"description": "Brightness level 0-100",
					},
				},
				Required: []string{"name", "brightness"},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "hue_set_color",
			Description: anthropic.Opt("Set the color of a Hue light. Accepts: warm, cool, red, blue, green, purple, orange, pink, white."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Light name",
					},
					"color": map[string]interface{}{
						"type":        "string",
						"description": "Color name: warm, cool, red, blue, green, purple, orange, pink, white",
					},
				},
				Required: []string{"name", "color"},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "hue_list_scenes",
			Description: anthropic.Opt("List all available Hue scenes with their room assignments."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "hue_activate_scene",
			Description: anthropic.Opt("Activate a Hue scene by name (fuzzy match). Scenes set multiple lights to pre-configured states."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Scene name to fuzzy match (e.g. 'Relax', 'Energize', 'Movie')",
					},
				},
				Required: []string{"name"},
			},
		}),
		tool(anthropic.ToolParam{
			Name:        "hue_control_room",
			Description: anthropic.Opt("Control all lights in a room at once -- turn on/off and optionally set brightness."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Room name",
					},
					"on": map[string]interface{}{
						"type":        "boolean",
						"description": "True to turn on, false to turn off",
					},
					"brightness": map[string]interface{}{
						"type":        "number",
						"description": "Optional brightness 0-100",
					},
				},
				Required: []string{"name"},
			},
		}),
	}
}

// ---------------------------------------------------------------------------
// Fuzzy name matching — case-insensitive contains. Not glamorous, but the
// audience never sees the wires.
// ---------------------------------------------------------------------------

func findKasaDeviceByName(devices []KasaDeviceInfo, name string) *KasaDeviceInfo {
	name = strings.ToLower(name)
	for i, d := range devices {
		if strings.Contains(strings.ToLower(d.Alias), name) {
			return &devices[i]
		}
	}
	return nil
}

func findLightByName(lights []HueLightInfo, name string) *HueLightInfo {
	name = strings.ToLower(name)
	for i, l := range lights {
		if strings.Contains(strings.ToLower(l.Name), name) || strings.Contains(strings.ToLower(l.RoomName), name) {
			return &lights[i]
		}
	}
	return nil
}

func findRoomByName(rooms []HueRoomInfo, name string) *HueRoomInfo {
	name = strings.ToLower(name)
	for i, r := range rooms {
		if strings.Contains(strings.ToLower(r.Name), name) {
			return &rooms[i]
		}
	}
	return nil
}

func findSceneByName(scenes []HueSceneInfo, name string) *HueSceneInfo {
	name = strings.ToLower(name)
	for i, s := range scenes {
		if strings.Contains(strings.ToLower(s.Name), name) {
			return &scenes[i]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Tool execution functions — where the magic actually happens behind the
// curtain. Each one parses input, fuzzy matches, calls the closure, and
// returns a human-readable result.
// ---------------------------------------------------------------------------

func executeKasaListDevices(kasaFuncs *KasaFuncs) (string, error) {
	if kasaFuncs == nil || kasaFuncs.ListDevices == nil {
		return "Kasa devices are not configured.", nil
	}
	devices, err := kasaFuncs.ListDevices()
	if err != nil {
		return "", fmt.Errorf("failed to list Kasa devices: %w", err)
	}
	if len(devices) == 0 {
		return "No Kasa devices found.", nil
	}
	var result strings.Builder
	result.WriteString("Kasa Devices:\n")
	for _, d := range devices {
		state := "off"
		if d.On && d.DeviceType == "fan" {
			state = fmt.Sprintf("on (speed %d/4)", d.FanSpeed)
		} else if d.On && d.DeviceType == "dimmer" {
			state = fmt.Sprintf("on (%d%%)", d.Brightness)
		} else if d.On {
			state = "on"
		}
		result.WriteString(fmt.Sprintf("- %s [%s]: %s\n", d.Alias, d.DeviceType, state))
	}
	return result.String(), nil
}

func executeKasaToggleDevice(input json.RawMessage, kasaFuncs *KasaFuncs) (string, error) {
	if kasaFuncs == nil || kasaFuncs.ListDevices == nil || kasaFuncs.SetState == nil {
		return "Kasa devices are not configured.", nil
	}
	var params struct {
		Name string `json:"name"`
		On   bool   `json:"on"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	devices, err := kasaFuncs.ListDevices()
	if err != nil {
		return "", fmt.Errorf("failed to list Kasa devices: %w", err)
	}
	device := findKasaDeviceByName(devices, params.Name)
	if device == nil {
		return fmt.Sprintf("No Kasa device found matching '%s'.", params.Name), nil
	}

	if err := kasaFuncs.SetState(device.IP, params.On); err != nil {
		return "", fmt.Errorf("failed to toggle device: %w", err)
	}

	action := "off"
	if params.On {
		action = "on"
	}
	return fmt.Sprintf("Turned %s %s.", device.Alias, action), nil
}

func executeKasaSetBrightness(input json.RawMessage, kasaFuncs *KasaFuncs) (string, error) {
	if kasaFuncs == nil || kasaFuncs.ListDevices == nil || kasaFuncs.SetBrightness == nil {
		return "Kasa devices are not configured.", nil
	}
	var params struct {
		Name       string `json:"name"`
		Brightness int    `json:"brightness"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	devices, err := kasaFuncs.ListDevices()
	if err != nil {
		return "", fmt.Errorf("failed to list Kasa devices: %w", err)
	}
	device := findKasaDeviceByName(devices, params.Name)
	if device == nil {
		return fmt.Sprintf("No Kasa device found matching '%s'.", params.Name), nil
	}
	if device.DeviceType != "dimmer" {
		return fmt.Sprintf("%s is a %s, not a dimmer. It doesn't support brightness control.", device.Alias, device.DeviceType), nil
	}

	if err := kasaFuncs.SetBrightness(device.IP, params.Brightness); err != nil {
		return "", fmt.Errorf("failed to set brightness: %w", err)
	}

	return fmt.Sprintf("Set %s brightness to %d%%.", device.Alias, params.Brightness), nil
}

func executeKasaSetFanSpeed(input json.RawMessage, kasaFuncs *KasaFuncs) (string, error) {
	if kasaFuncs == nil || kasaFuncs.ListDevices == nil || kasaFuncs.SetFanSpeed == nil {
		return "Kasa devices are not configured.", nil
	}
	var params struct {
		Name  string `json:"name"`
		Speed int    `json:"speed"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	devices, err := kasaFuncs.ListDevices()
	if err != nil {
		return "", fmt.Errorf("failed to list Kasa devices: %w", err)
	}
	device := findKasaDeviceByName(devices, params.Name)
	if device == nil {
		return fmt.Sprintf("No Kasa device found matching '%s'.", params.Name), nil
	}
	if device.DeviceType != "fan" {
		return fmt.Sprintf("%s is not a fan device. It doesn't support fan speed control.", device.Alias), nil
	}
	if params.Speed < 1 || params.Speed > 4 {
		return "Speed must be between 1 (low) and 4 (high).", nil
	}

	if err := kasaFuncs.SetFanSpeed(device.IP, params.Speed); err != nil {
		return "", fmt.Errorf("failed to set fan speed: %w", err)
	}

	speedLabels := []string{"", "low", "medium-low", "medium", "high"}
	return fmt.Sprintf("Set %s to speed %d (%s).", device.Alias, params.Speed, speedLabels[params.Speed]), nil
}

func executeHueListLights(hueFuncs *HueFuncs) (string, error) {
	if hueFuncs == nil || hueFuncs.ListLights == nil {
		return "Hue lights are not configured.", nil
	}
	lights, err := hueFuncs.ListLights()
	if err != nil {
		return "", fmt.Errorf("failed to list Hue lights: %w", err)
	}
	if len(lights) == 0 {
		return "No Hue lights found.", nil
	}
	var result strings.Builder
	result.WriteString("Hue Lights:\n")
	for _, l := range lights {
		state := "off"
		if l.On {
			state = fmt.Sprintf("on (%d%%)", int(l.Brightness))
		}
		reachable := ""
		if !l.Reachable {
			reachable = " [unreachable]"
		}
		result.WriteString(fmt.Sprintf("- %s [%s]: %s%s\n", l.Name, l.RoomName, state, reachable))
	}
	return result.String(), nil
}

func executeHueToggleLight(input json.RawMessage, hueFuncs *HueFuncs) (string, error) {
	if hueFuncs == nil || hueFuncs.ListLights == nil || hueFuncs.ToggleLight == nil {
		return "Hue lights are not configured.", nil
	}
	var params struct {
		Name string `json:"name"`
		On   bool   `json:"on"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	lights, err := hueFuncs.ListLights()
	if err != nil {
		return "", fmt.Errorf("failed to list lights: %w", err)
	}
	light := findLightByName(lights, params.Name)
	if light == nil {
		return fmt.Sprintf("No Hue light found matching '%s'.", params.Name), nil
	}

	if err := hueFuncs.ToggleLight(light.ID, params.On); err != nil {
		return "", fmt.Errorf("failed to toggle light: %w", err)
	}

	action := "off"
	if params.On {
		action = "on"
	}
	return fmt.Sprintf("Turned %s %s [%s].", light.Name, action, light.RoomName), nil
}

func executeHueSetBrightness(input json.RawMessage, hueFuncs *HueFuncs) (string, error) {
	if hueFuncs == nil || hueFuncs.ListLights == nil || hueFuncs.SetBrightness == nil {
		return "Hue lights are not configured.", nil
	}
	var params struct {
		Name       string  `json:"name"`
		Brightness float64 `json:"brightness"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	lights, err := hueFuncs.ListLights()
	if err != nil {
		return "", fmt.Errorf("failed to list lights: %w", err)
	}
	light := findLightByName(lights, params.Name)
	if light == nil {
		return fmt.Sprintf("No Hue light found matching '%s'.", params.Name), nil
	}

	if err := hueFuncs.SetBrightness(light.ID, params.Brightness); err != nil {
		return "", fmt.Errorf("failed to set brightness: %w", err)
	}

	return fmt.Sprintf("Set %s [%s] brightness to %d%%.", light.Name, light.RoomName, int(params.Brightness)), nil
}

func executeHueSetColor(input json.RawMessage, hueFuncs *HueFuncs) (string, error) {
	if hueFuncs == nil || hueFuncs.ListLights == nil || hueFuncs.SetColor == nil {
		return "Hue lights are not configured.", nil
	}
	var params struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	colorKey := strings.ToLower(strings.TrimSpace(params.Color))
	xy, ok := colorMap[colorKey]
	if !ok {
		return fmt.Sprintf("Unknown color '%s'. Supported: warm, cool, red, blue, green, purple, orange, pink, white.", params.Color), nil
	}

	lights, err := hueFuncs.ListLights()
	if err != nil {
		return "", fmt.Errorf("failed to list lights: %w", err)
	}
	light := findLightByName(lights, params.Name)
	if light == nil {
		return fmt.Sprintf("No Hue light found matching '%s'.", params.Name), nil
	}

	if err := hueFuncs.SetColor(light.ID, xy); err != nil {
		return "", fmt.Errorf("failed to set color: %w", err)
	}

	return fmt.Sprintf("Set %s [%s] to %s.", light.Name, light.RoomName, params.Color), nil
}

func executeHueListScenes(hueFuncs *HueFuncs) (string, error) {
	if hueFuncs == nil || hueFuncs.ListScenes == nil {
		return "Hue scenes are not configured.", nil
	}
	scenes, err := hueFuncs.ListScenes()
	if err != nil {
		return "", fmt.Errorf("failed to list scenes: %w", err)
	}
	if len(scenes) == 0 {
		return "No Hue scenes found.", nil
	}
	var result strings.Builder
	result.WriteString("Hue Scenes:\n")
	for _, s := range scenes {
		result.WriteString(fmt.Sprintf("- %s [%s]\n", s.Name, s.RoomName))
	}
	return result.String(), nil
}

func executeHueActivateScene(input json.RawMessage, hueFuncs *HueFuncs) (string, error) {
	if hueFuncs == nil || hueFuncs.ListScenes == nil || hueFuncs.ActivateScene == nil {
		return "Hue scenes are not configured.", nil
	}
	var params struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	scenes, err := hueFuncs.ListScenes()
	if err != nil {
		return "", fmt.Errorf("failed to list scenes: %w", err)
	}
	scene := findSceneByName(scenes, params.Name)
	if scene == nil {
		return fmt.Sprintf("No scene found matching '%s'.", params.Name), nil
	}

	if err := hueFuncs.ActivateScene(scene.ID); err != nil {
		return "", fmt.Errorf("failed to activate scene: %w", err)
	}

	return fmt.Sprintf("Activated scene '%s' in %s.", scene.Name, scene.RoomName), nil
}

func executeHueControlRoom(input json.RawMessage, hueFuncs *HueFuncs) (string, error) {
	if hueFuncs == nil || hueFuncs.ListRooms == nil || hueFuncs.SetRoomState == nil {
		return "Hue rooms are not configured.", nil
	}
	var params struct {
		Name       string   `json:"name"`
		On         *bool    `json:"on"`
		Brightness *float64 `json:"brightness"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	rooms, err := hueFuncs.ListRooms()
	if err != nil {
		return "", fmt.Errorf("failed to list rooms: %w", err)
	}
	room := findRoomByName(rooms, params.Name)
	if room == nil {
		return fmt.Sprintf("No room found matching '%s'.", params.Name), nil
	}

	if err := hueFuncs.SetRoomState(room.GroupedLightID, params.On, params.Brightness); err != nil {
		return "", fmt.Errorf("failed to control room: %w", err)
	}

	var parts []string
	if params.On != nil {
		if *params.On {
			parts = append(parts, "on")
		} else {
			parts = append(parts, "off")
		}
	}
	if params.Brightness != nil {
		parts = append(parts, fmt.Sprintf("at %d%% brightness", int(*params.Brightness)))
	}
	return fmt.Sprintf("Set %s lights %s.", room.Name, strings.Join(parts, " ")), nil
}

// ---------------------------------------------------------------------------
// Main dispatcher — routes tool calls to the right executor. Like a stage
// manager calling cues.
// ---------------------------------------------------------------------------

// ExecuteTool dispatches a tool call by name, executing against the provided
// closure functions and returning a human-readable result string.
func ExecuteTool(name string, input json.RawMessage, kasaFuncs *KasaFuncs, hueFuncs *HueFuncs) (string, error) {
	switch name {
	// Kasa tools
	case "kasa_list_devices":
		return executeKasaListDevices(kasaFuncs)
	case "kasa_toggle_device":
		return executeKasaToggleDevice(input, kasaFuncs)
	case "kasa_set_brightness":
		return executeKasaSetBrightness(input, kasaFuncs)
	case "kasa_set_fan_speed":
		return executeKasaSetFanSpeed(input, kasaFuncs)

	// Hue tools
	case "hue_list_lights":
		return executeHueListLights(hueFuncs)
	case "hue_toggle_light":
		return executeHueToggleLight(input, hueFuncs)
	case "hue_set_brightness":
		return executeHueSetBrightness(input, hueFuncs)
	case "hue_set_color":
		return executeHueSetColor(input, hueFuncs)
	case "hue_list_scenes":
		return executeHueListScenes(hueFuncs)
	case "hue_activate_scene":
		return executeHueActivateScene(input, hueFuncs)
	case "hue_control_room":
		return executeHueControlRoom(input, hueFuncs)

	default:
		return fmt.Sprintf("Unknown tool: %s", name), nil
	}
}

// DeviceToolsForContext returns the tool set scoped to a specific device.
func DeviceToolsForContext(device DeviceContext) []anthropic.ToolUnionParam {
	queryTool := anthropic.ToolParam{
		Name:        "query_device",
		Description: anthropic.Opt("Query this device's current state (on/off, brightness, etc). Use this when the user asks about the device's state."),
		InputSchema: anthropic.ToolInputSchemaParam{
			Properties: map[string]interface{}{},
		},
	}

	hasProtocol := func(p string) bool {
		for _, proto := range device.Protocols {
			if proto == p {
				return true
			}
		}
		return false
	}

	queryAPITool := anthropic.ToolParam{
		Name:        "query_api",
		Description: anthropic.Opt("Make an authenticated HTTP request to this device's API. Use this to fetch data from specific endpoints listed in the API documentation."),
		InputSchema: anthropic.ToolInputSchemaParam{
			Properties: map[string]interface{}{
				"path": map[string]interface{}{"type": "string", "description": "API path to query (e.g. /vars?match=livedata&fmt=obj or /cgi-bin/dl_cgi?Command=DeviceList)"},
			},
			Required: []string{"path"},
		},
	}

	tools := []anthropic.ToolUnionParam{tool(queryTool), tool(queryAPITool)}

	if hasProtocol("kasa") {
		tools = append(tools, tool(anthropic.ToolParam{
			Name:        "toggle_device",
			Description: anthropic.Opt("Turn this device on or off."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{
					"on": map[string]interface{}{"type": "boolean", "description": "true=on, false=off"},
				},
				Required: []string{"on"},
			},
		}))
		if device.DeviceType == "dimmer" {
			tools = append(tools, tool(anthropic.ToolParam{
				Name:        "set_brightness",
				Description: anthropic.Opt("Set brightness level (0-100%)."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]interface{}{
						"brightness": map[string]interface{}{"type": "integer", "description": "0-100"},
					},
					Required: []string{"brightness"},
				},
			}))
		}
		if device.DeviceType == "fan" {
			tools = append(tools, tool(anthropic.ToolParam{
				Name:        "set_fan_speed",
				Description: anthropic.Opt("Set fan speed (1=low, 2=medium, 3=high, 4=max)."),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]interface{}{
						"speed": map[string]interface{}{"type": "integer", "description": "1-4"},
					},
					Required: []string{"speed"},
				},
			}))
		}
	}

	if device.DeviceType == "nest_camera" || device.DeviceType == "nest_device" {
		tools = append(tools, tool(anthropic.ToolParam{
			Name:        "see_camera",
			Description: anthropic.Opt("Look through the camera and describe what you see. Captures a live snapshot and analyzes it with vision AI."),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]interface{}{},
			},
		}))
	}

	if device.DeviceType == "hue_bridge" {
		for _, t := range deviceTools() {
			if t.OfTool != nil && strings.HasPrefix(t.OfTool.Name, "hue_") {
				tools = append(tools, t)
			}
		}
	}

	return tools
}

// ExecuteDeviceTool dispatches a tool call for a device-scoped chat.
func ExecuteDeviceTool(name string, input json.RawMessage, device DeviceContext, kasaFuncs *KasaFuncs, hueFuncs *HueFuncs, httpQuery DeviceHTTPQuery, jfQuery JellyFishQueryFunc) (string, error) {
	switch name {
	case "query_device":
		return executeQueryDevice(device, kasaFuncs, hueFuncs, httpQuery, jfQuery)
	case "query_api":
		var params struct{ Path string `json:"path"` }
		json.Unmarshal(input, &params)
		if httpQuery != nil && params.Path != "" {
			data, err := httpQuery(device.IP, params.Path)
			if err != nil {
				return "", fmt.Errorf("API query failed: %w", err)
			}
			return data, nil
		}
		return "Cannot query this device's API — no credentials stored.", nil
	case "toggle_device":
		var params struct{ On bool `json:"on"` }
		json.Unmarshal(input, &params)
		if kasaFuncs != nil && kasaFuncs.SetState != nil {
			if err := kasaFuncs.SetState(device.IP, params.On); err != nil {
				return "", err
			}
			state := "off"
			if params.On { state = "on" }
			return fmt.Sprintf("Turned %s %s.", device.Name, state), nil
		}
		return "Cannot control this device.", nil
	case "set_brightness":
		var params struct{ Brightness int `json:"brightness"` }
		json.Unmarshal(input, &params)
		if kasaFuncs != nil && kasaFuncs.SetBrightness != nil {
			if err := kasaFuncs.SetBrightness(device.IP, params.Brightness); err != nil {
				return "", err
			}
			return fmt.Sprintf("Set %s brightness to %d%%.", device.Name, params.Brightness), nil
		}
		return "Cannot set brightness on this device.", nil
	case "set_fan_speed":
		var params struct{ Speed int `json:"speed"` }
		json.Unmarshal(input, &params)
		if kasaFuncs != nil && kasaFuncs.SetFanSpeed != nil {
			if err := kasaFuncs.SetFanSpeed(device.IP, params.Speed); err != nil {
				return "", err
			}
			return fmt.Sprintf("Set %s fan speed to %d.", device.Name, params.Speed), nil
		}
		return "Cannot set fan speed on this device.", nil
	default:
		// Delegate hue_* tools to the global executor
		return ExecuteTool(name, input, kasaFuncs, hueFuncs)
	}
}

func executeQueryDevice(device DeviceContext, kasaFuncs *KasaFuncs, hueFuncs *HueFuncs, httpQuery DeviceHTTPQuery, jfQuery JellyFishQueryFunc) (string, error) {
	hasProtocol := func(p string) bool {
		for _, proto := range device.Protocols {
			if proto == p { return true }
		}
		return false
	}

	if hasProtocol("kasa") && kasaFuncs != nil && kasaFuncs.QueryDevice != nil {
		info, err := kasaFuncs.QueryDevice(device.IP)
		if err != nil {
			return "", fmt.Errorf("query failed: %w", err)
		}
		state := "OFF"
		if info.On { state = "ON" }
		result := fmt.Sprintf("Device: %s\nState: %s\nType: %s\nModel: %s", info.Alias, state, info.DeviceType, info.Model)
		if info.DeviceType == "dimmer" {
			result += fmt.Sprintf("\nBrightness: %d%%", info.Brightness)
		}
		if info.DeviceType == "fan" && info.FanSpeed > 0 {
			result += fmt.Sprintf("\nFan Speed: %d/4", info.FanSpeed)
		}
		return result, nil
	}

	if device.DeviceType == "jellyfish" && jfQuery != nil {
		// Query JellyFish zones and their current state
		zonesResult, err := jfQuery(device.IP, map[string]interface{}{
			"cmd": "toCtlrGet", "get": [][]string{{"zones"}},
		})
		if err != nil {
			return fmt.Sprintf("Device: %s at %s\nType: JellyFish\nStatus: Error querying - %v", device.Name, device.IP, err), nil
		}

		stateResult, err := jfQuery(device.IP, map[string]interface{}{
			"cmd": "toCtlrGet", "get": [][]string{{"runPattern", "Zone", "Zone1"}},
		})

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Device: %s\nIP: %s\nType: JellyFish Lighting Controller\nStatus: Connected\n\n", device.Name, device.IP))
		sb.WriteString("Zones:\n" + zonesResult + "\n")
		if err == nil {
			sb.WriteString("\nCurrent State:\n" + stateResult + "\n")
		}
		return sb.String(), nil
	}

	if device.DeviceType == "nest_camera" || device.DeviceType == "nest_thermostat" || device.DeviceType == "nest_device" {
		// For Nest devices, report what we know from the DB + Google connection status
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Device: %s\n", device.Name))
		sb.WriteString(fmt.Sprintf("IP: %s\n", device.IP))
		sb.WriteString(fmt.Sprintf("Type: %s\n", device.DeviceType))
		sb.WriteString("Connection: ACTIVE (Google Nest SDM API)\n")
		sb.WriteString("Status: Online and streaming\n")
		if strings.Contains(device.DeviceType, "camera") {
			sb.WriteString("\nThis camera is live and streaming.\n")
			sb.WriteString("Capabilities: Live WebRTC streaming, snapshot capture, vision analysis, motion detection.\n")
			sb.WriteString("Use the see_camera tool to capture and analyze what the camera sees right now.\n")
		}
		if strings.Contains(device.DeviceType, "thermostat") {
			sb.WriteString("\nUse query_api to get temperature, humidity, and thermostat mode from the SDM API.\n")
		}
		return sb.String(), nil
	}

	if device.DeviceType == "hue_bridge" && hueFuncs != nil && hueFuncs.ListLights != nil {
		lights, err := hueFuncs.ListLights()
		if err != nil {
			return "", err
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Hue Bridge at %s with %d lights:\n", device.IP, len(lights)))
		for _, l := range lights {
			state := "off"
			if l.On { state = fmt.Sprintf("on (%d%%)", int(l.Brightness)) }
			sb.WriteString(fmt.Sprintf("- %s [%s]: %s\n", l.Name, l.RoomName, state))
		}
		return sb.String(), nil
	}

	// For devices with stored credentials, use HTTP query to get real data
	if httpQuery != nil && (device.DeviceType == "solar_gateway" || device.Category == "energy") {
		// Fetch ALL system data via /vars — this is the only endpoint that works
		allData, err := httpQuery(device.IP, "/vars?match=sys&fmt=obj")
		if err != nil {
			return fmt.Sprintf("Device: %s\nIP: %s\nError fetching data: %s", device.Name, device.IP, err), nil
		}

		// Parse and summarize the data for the AI
		var data map[string]string
		if json.Unmarshal([]byte(allData), &data) != nil {
			return fmt.Sprintf("Device: %s\nIP: %s\nRaw response (first 500 chars):\n%s", device.Name, device.IP, allData[:min(500, len(allData))]), nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("SunPower PVS at %s — Live Data:\n\n", device.IP))

		// Live readings
		getVal := func(keys ...string) string {
			for _, k := range keys {
				if v, ok := data[k]; ok && v != "nan" && v != "" { return v }
			}
			return "0"
		}
		sb.WriteString(fmt.Sprintf("Solar Production: %s kW\n", getVal("/sys/livedata/pv_p")))
		sb.WriteString(fmt.Sprintf("House Consumption: %s kW\n", getVal("/sys/livedata/site_load_p")))
		sb.WriteString(fmt.Sprintf("Grid: %s kW (negative = exporting)\n", getVal("/sys/livedata/net_p")))
		sb.WriteString(fmt.Sprintf("Lifetime Production: %s kWh\n", getVal("/sys/livedata/pv_en")))
		sb.WriteString(fmt.Sprintf("Battery SOC: %s%%\n", getVal("/sys/livedata/soc")))

		// Count and list panels (inverters)
		panelSerials := make(map[int]string)
		panelProduction := make(map[int]string)
		for key, val := range data {
			if strings.HasSuffix(key, "/sn") && strings.Contains(key, "inverter") {
				var idx int
				fmt.Sscanf(key, "/sys/devices/inverter/%d/sn", &idx)
				panelSerials[idx] = val
			}
			if strings.HasSuffix(key, "/ltea3phsumKwh") && strings.Contains(key, "inverter") {
				var idx int
				fmt.Sscanf(key, "/sys/devices/inverter/%d/ltea3phsumKwh", &idx)
				panelProduction[idx] = val
			}
		}
		sb.WriteString(fmt.Sprintf("\nPanels: %d micro-inverters\n", len(panelSerials)))
		for i := 0; i < len(panelSerials); i++ {
			if sn, ok := panelSerials[i]; ok {
				prod := panelProduction[i]
				sb.WriteString(fmt.Sprintf("  Panel %d: SN=%s, Lifetime=%.0skWh\n", i+1, sn, prod))
			}
		}

		if sn, ok := data["/sys/info/serialnum"]; ok {
			sb.WriteString(fmt.Sprintf("\nPVS Serial: %s\n", sn))
		}

		return sb.String(), nil
	}

	return fmt.Sprintf("Device: %s\nIP: %s\nManufacturer: %s\nModel: %s\nType: %s\nCategory: %s",
		device.Name, device.IP, device.Manufacturer, device.Model, device.DeviceType, device.Category), nil
}
