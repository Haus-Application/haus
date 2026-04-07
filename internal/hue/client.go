package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client communicates with the Hue API v2 on a local bridge.
// The bridge uses self-signed TLS certs, so we skip verification.
// George Sr. would be concerned, but it's a local network, so it's fine.
// ... right?
type Client struct {
	bridgeIP string
	username string
	http     *http.Client
}

// NewClient creates a Client for the given bridge IP and application username
// (API key). The HTTP client skips TLS verification because Hue bridges use
// self-signed certificates.
func NewClient(bridgeIP, username string) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // Hue bridge uses a self-signed cert
	}
	return &Client{
		bridgeIP: bridgeIP,
		username: username,
		http:     &http.Client{Timeout: 10 * time.Second, Transport: transport},
	}
}

// BridgeIP returns the IP address of the connected bridge.
func (c *Client) BridgeIP() string {
	return c.bridgeIP
}

// baseURL returns the Hue API v2 base URL for this bridge.
func (c *Client) baseURL() string {
	return "https://" + c.bridgeIP + "/clip/v2"
}

// do performs an authenticated request to the Hue API.
func (c *Client) do(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("hue: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL()+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("hue: build request: %w", err)
	}
	req.Header.Set("hue-application-key", c.username)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hue: request %s %s: %w", method, path, err)
	}
	return resp, nil
}

// hueResponse is the standard Hue API v2 envelope.
type hueResponse struct {
	Errors []struct {
		Description string `json:"description"`
	} `json:"errors"`
	Data json.RawMessage `json:"data"`
}

// decodeResponse reads and decodes a Hue API v2 envelope, returning the raw
// data payload. It returns an error if the envelope contains API errors.
func decodeResponse(resp *http.Response) (json.RawMessage, error) {
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hue: response status %d: %s", resp.StatusCode, body)
	}
	var env hueResponse
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("hue: decode response: %w", err)
	}
	if len(env.Errors) > 0 {
		return nil, fmt.Errorf("hue: API error: %s", env.Errors[0].Description)
	}
	return env.Data, nil
}

// rawLight is the on-the-wire format returned by the Hue API for a light resource.
type rawLight struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	On struct {
		On bool `json:"on"`
	} `json:"on"`
	Dimming struct {
		Brightness float64 `json:"brightness"`
	} `json:"dimming"`
	Color struct {
		XY struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"xy"`
	} `json:"color"`
	ColorTemperature struct {
		Mirek *int `json:"mirek"`
	} `json:"color_temperature"`
	Owner struct {
		RID   string `json:"rid"`
		RType string `json:"rtype"`
	} `json:"owner"`
}

// Light represents a single Hue light with its current state.
type Light struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	On         bool       `json:"on"`
	Brightness float64    `json:"brightness"` // 0-100
	ColorXY    [2]float64 `json:"color_xy,omitempty"`
	ColorTemp  int        `json:"color_temp,omitempty"` // mirek
	RoomID     string     `json:"room_id"`
	RoomName   string     `json:"room_name"`
}

// Room represents a Hue room with its member lights.
type Room struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Lights         []Light `json:"lights"`
	GroupedLightID string  `json:"grouped_light_id"`
}

// Scene represents a saved Hue scene.
type Scene struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RoomID   string `json:"room_id"`
	RoomName string `json:"room_name"`
}

// rawRoom is the on-the-wire format for a room resource.
type rawRoom struct {
	ID       string `json:"id"`
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Children []struct {
		RID   string `json:"rid"`
		RType string `json:"rtype"`
	} `json:"children"`
	Services []struct {
		RID   string `json:"rid"`
		RType string `json:"rtype"`
	} `json:"services"`
}

// rawScene is the on-the-wire format for a scene resource.
type rawScene struct {
	ID       string `json:"id"`
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Group struct {
		RID   string `json:"rid"`
		RType string `json:"rtype"`
	} `json:"group"`
}

// ListLights returns all lights registered on the bridge.
//
// GET /clip/v2/resource/light
func (c *Client) ListLights() ([]Light, error) {
	resp, err := c.do("GET", "/resource/light", nil)
	if err != nil {
		return nil, err
	}
	data, err := decodeResponse(resp)
	if err != nil {
		return nil, err
	}

	var raw []rawLight
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("hue: unmarshal lights: %w", err)
	}

	lights := make([]Light, 0, len(raw))
	for _, r := range raw {
		l := Light{
			ID:         r.ID,
			Name:       r.Metadata.Name,
			On:         r.On.On,
			Brightness: r.Dimming.Brightness,
			ColorXY:    [2]float64{r.Color.XY.X, r.Color.XY.Y},
		}
		if r.ColorTemperature.Mirek != nil {
			l.ColorTemp = *r.ColorTemperature.Mirek
		}
		lights = append(lights, l)
	}
	return lights, nil
}

// ListRooms returns all rooms, each populated with the lights that belong to it.
//
// GET /clip/v2/resource/room
func (c *Client) ListRooms() ([]Room, error) {
	resp, err := c.do("GET", "/resource/room", nil)
	if err != nil {
		return nil, err
	}
	data, err := decodeResponse(resp)
	if err != nil {
		return nil, err
	}

	var raw []rawRoom
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("hue: unmarshal rooms: %w", err)
	}

	// We need to fetch device resources to map room children (devices) to lights.
	// Room children reference device IDs, but lights are separate resources owned by devices.
	allLights, err := c.ListLights()
	if err != nil {
		// Non-fatal: return rooms without light details.
		allLights = nil
	}

	// Build a lookup of light by owner device RID so we can match room children.
	// Room children are "device" type, lights have an owner pointing to their device.
	lightByID := make(map[string]*Light, len(allLights))
	for i := range allLights {
		lightByID[allLights[i].ID] = &allLights[i]
	}

	rooms := make([]Room, 0, len(raw))
	for _, r := range raw {
		room := Room{
			ID:   r.ID,
			Name: r.Metadata.Name,
		}
		// Find the grouped_light service for this room.
		for _, svc := range r.Services {
			if svc.RType == "grouped_light" {
				room.GroupedLightID = svc.RID
				break
			}
		}
		// Collect member lights.
		for _, child := range r.Children {
			if child.RType == "light" {
				if l, ok := lightByID[child.RID]; ok {
					lCopy := *l
					lCopy.RoomID = room.ID
					lCopy.RoomName = room.Name
					room.Lights = append(room.Lights, lCopy)
				}
			}
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

// ListScenes returns all scenes registered on the bridge.
//
// GET /clip/v2/resource/scene
func (c *Client) ListScenes() ([]Scene, error) {
	// Fetch rooms first so we can annotate scenes with room names.
	allRooms, _ := c.ListRooms()
	roomNameByID := make(map[string]string, len(allRooms))
	for _, r := range allRooms {
		roomNameByID[r.ID] = r.Name
	}

	resp, err := c.do("GET", "/resource/scene", nil)
	if err != nil {
		return nil, err
	}
	data, err := decodeResponse(resp)
	if err != nil {
		return nil, err
	}

	var raw []rawScene
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("hue: unmarshal scenes: %w", err)
	}

	scenes := make([]Scene, 0, len(raw))
	for _, r := range raw {
		roomID := ""
		if r.Group.RType == "room" {
			roomID = r.Group.RID
		}
		scenes = append(scenes, Scene{
			ID:       r.ID,
			Name:     r.Metadata.Name,
			RoomID:   roomID,
			RoomName: roomNameByID[roomID],
		})
	}
	return scenes, nil
}

// SetLightState updates a light's on/off, brightness, and/or color.
// Pass nil for any parameter you do not want to change.
//
// PUT /clip/v2/resource/light/<id>
func (c *Client) SetLightState(id string, on *bool, brightness *float64, colorXY *[2]float64) error {
	body := make(map[string]interface{})
	if on != nil {
		body["on"] = map[string]bool{"on": *on}
	}
	if brightness != nil {
		body["dimming"] = map[string]float64{"brightness": *brightness}
	}
	if colorXY != nil {
		body["color"] = map[string]interface{}{
			"xy": map[string]float64{"x": colorXY[0], "y": colorXY[1]},
		}
	}

	resp, err := c.do("PUT", "/resource/light/"+id, body)
	if err != nil {
		return err
	}
	_, err = decodeResponse(resp)
	return err
}

// SetGroupedLightState updates the on/off and brightness for all lights in a
// room via its grouped_light resource.
//
// PUT /clip/v2/resource/grouped_light/<id>
func (c *Client) SetGroupedLightState(id string, on *bool, brightness *float64) error {
	body := make(map[string]interface{})
	if on != nil {
		body["on"] = map[string]bool{"on": *on}
	}
	if brightness != nil {
		body["dimming"] = map[string]float64{"brightness": *brightness}
	}

	resp, err := c.do("PUT", "/resource/grouped_light/"+id, body)
	if err != nil {
		return err
	}
	_, err = decodeResponse(resp)
	return err
}

// ActivateScene recalls a scene by ID, making the lights adopt its saved state.
// This is basically a magic trick, but unlike GOB, it actually works every time.
//
// PUT /clip/v2/resource/scene/<id>
func (c *Client) ActivateScene(id string) error {
	body := map[string]interface{}{
		"recall": map[string]string{"action": "active"},
	}
	resp, err := c.do("PUT", "/resource/scene/"+id, body)
	if err != nil {
		return err
	}
	_, err = decodeResponse(resp)
	return err
}

// pairResponse is the response from the Hue bridge pairing endpoint.
type pairResponse struct {
	Success *struct {
		Username string `json:"username"`
	} `json:"success"`
	Error *struct {
		Description string `json:"description"`
	} `json:"error"`
}

// Pair initiates the bridge link button pairing flow. The user must press the
// physical link button on the bridge before calling this. Returns the username
// (application key) that should be stored for future requests.
// This is the most nerve-wracking 30 seconds of my life, every time.
func Pair(bridgeIP string) (string, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // Hue bridge uses a self-signed cert
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	body, _ := json.Marshal(map[string]interface{}{
		"devicetype":        "haus#app",
		"generateclientkey": true,
	})
	resp, err := client.Post("https://"+bridgeIP+"/api", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("hue: pair request: %w", err)
	}
	defer resp.Body.Close()

	var results []pairResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return "", fmt.Errorf("hue: decode pair response: %w", err)
	}
	if len(results) == 0 {
		return "", fmt.Errorf("hue: empty pair response")
	}
	r := results[0]
	if r.Error != nil {
		return "", fmt.Errorf("hue: pair error: %s", r.Error.Description)
	}
	if r.Success == nil || r.Success.Username == "" {
		return "", fmt.Errorf("hue: pair did not return a username")
	}
	return r.Success.Username, nil
}
