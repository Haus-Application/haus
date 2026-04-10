package nest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const baseURL = "https://smartdevicemanagement.googleapis.com/v1"

// Client communicates with the Google Smart Device Management (SDM) API.
// This is like connecting to Mother, but Mother is Google, and she requires
// OAuth2 tokens instead of juice boxes.
type Client struct {
	projectID   string // Device Access project ID (UUID from Google Cloud)
	accessToken string
	http        *http.Client
}

// NewClient creates a Client for the given SDM project ID and OAuth2 access token.
// The project ID is the one you get from the Device Access Console -- not the
// Google Cloud project, that's a different thing entirely, and yes, I've made
// that mistake before and Mother was NOT happy.
func NewClient(projectID, accessToken string) *Client {
	return &Client{
		projectID:   projectID,
		accessToken: accessToken,
		http:        &http.Client{Timeout: 15 * time.Second},
	}
}

// SetAccessToken updates the client's OAuth2 access token. Useful after a
// token refresh -- tokens expire every 3600 seconds, which is exactly how
// long it takes me to stop worrying about the last token expiring.
func (c *Client) SetAccessToken(token string) {
	c.accessToken = token
}

// NestDevice represents a single device from the SDM API. The Name field is
// the fully qualified resource path: enterprises/{project}/devices/{id}.
// The Type tells you what kind of device it is -- thermostat, camera, doorbell,
// or display. Each type has different traits, which is a LOT like how each
// Bluth has different... traits.
type NestDevice struct {
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Traits          map[string]interface{} `json:"traits"`
	ParentRelations []ParentRelation       `json:"parentRelations"`
}

// ParentRelation describes the room/structure a device belongs to.
type ParentRelation struct {
	Parent      string `json:"parent"`
	DisplayName string `json:"displayName"`
}

// listResponse is the envelope for the ListDevices endpoint.
type listResponse struct {
	Devices []NestDevice `json:"devices"`
}

// deviceResponse is the envelope for the GetDevice endpoint.
// (It's just the device itself, no wrapper. Google is inconsistent like that.)

// executeCommandRequest is the request body for ExecuteCommand.
type executeCommandRequest struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

// do performs an authenticated request to the SDM API. Every request gets
// the Bearer token header. If Google rejects us, we return a clear error
// so the caller knows to refresh the token.
func (c *Client) do(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("nest: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	url := baseURL + path
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("nest: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nest: request %s %s: %w", method, path, err)
	}
	return resp, nil
}

// checkResponse reads the response body and returns an error if the status
// code is not 2xx. SDM API errors come back as JSON with an "error" field.
func checkResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nest: read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("nest: response status %d: %s", resp.StatusCode, data)
	}
	return data, nil
}

// enterprisesPrefix returns the path prefix for this project's resources.
func (c *Client) enterprisesPrefix() string {
	return "/enterprises/" + c.projectID
}

// ListDevices returns all devices in the SDM project. This is the first thing
// you call after pairing -- it's like roll call, but for thermostats and cameras.
//
// GET /v1/enterprises/{project}/devices
func (c *Client) ListDevices() ([]NestDevice, error) {
	resp, err := c.do("GET", c.enterprisesPrefix()+"/devices", nil)
	if err != nil {
		return nil, err
	}
	data, err := checkResponse(resp)
	if err != nil {
		return nil, err
	}

	var result listResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("nest: unmarshal devices: %w", err)
	}
	return result.Devices, nil
}

// GetDevice returns a single device by its device ID (not the full resource name,
// just the UUID part at the end).
//
// GET /v1/enterprises/{project}/devices/{id}
func (c *Client) GetDevice(deviceID string) (*NestDevice, error) {
	resp, err := c.do("GET", c.enterprisesPrefix()+"/devices/"+deviceID, nil)
	if err != nil {
		return nil, err
	}
	data, err := checkResponse(resp)
	if err != nil {
		return nil, err
	}

	var device NestDevice
	if err := json.Unmarshal(data, &device); err != nil {
		return nil, fmt.Errorf("nest: unmarshal device: %w", err)
	}
	return &device, nil
}

// ExecuteCommand sends a command to a device. Commands are how you tell a
// thermostat to change mode or a camera to start streaming. The command string
// is the full SDM command name, e.g. "sdm.devices.commands.ThermostatMode.SetMode".
//
// POST /v1/enterprises/{project}/devices/{id}:executeCommand
func (c *Client) ExecuteCommand(deviceID, command string, params map[string]interface{}) error {
	body := executeCommandRequest{
		Command: command,
		Params:  params,
	}
	resp, err := c.do("POST", c.enterprisesPrefix()+"/devices/"+deviceID+":executeCommand", body)
	if err != nil {
		return err
	}
	_, err = checkResponse(resp)
	return err
}

// --- Trait helper functions ---
// These extract typed values from the unstructured traits map. The SDM API
// returns traits as a map of string -> arbitrary JSON, which means we have to
// do a lot of type assertions. It's tedious but necessary, like sorting
// Mother's pill organizer.

// getTrait extracts a trait map by its full trait name.
func getTrait(device NestDevice, traitName string) map[string]interface{} {
	raw, ok := device.Traits[traitName]
	if !ok {
		return nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	return m
}

// GetThermostatMode returns the current thermostat mode (HEAT, COOL, HEATCOOL, OFF).
// Returns empty string if the device doesn't have the ThermostatMode trait.
func GetThermostatMode(device NestDevice) string {
	t := getTrait(device, "sdm.devices.traits.ThermostatMode")
	if t == nil {
		return ""
	}
	mode, _ := t["mode"].(string)
	return mode
}

// GetTemperature returns the ambient temperature in Celsius from the Temperature trait.
// Returns 0 if the trait is not present.
func GetTemperature(device NestDevice) float64 {
	t := getTrait(device, "sdm.devices.traits.Temperature")
	if t == nil {
		return 0
	}
	temp, _ := t["ambientTemperatureCelsius"].(float64)
	return temp
}

// GetHumidity returns the ambient relative humidity percentage from the Humidity trait.
// Returns 0 if the trait is not present.
func GetHumidity(device NestDevice) float64 {
	t := getTrait(device, "sdm.devices.traits.Humidity")
	if t == nil {
		return 0
	}
	humidity, _ := t["ambientHumidityPercent"].(float64)
	return humidity
}

// GetThermostatSetpoints returns the heat and cool setpoint temperatures in Celsius.
// Which values are populated depends on the current mode:
//   - HEAT mode: only heatCelsius is set
//   - COOL mode: only coolCelsius is set
//   - HEATCOOL mode: both are set
func GetThermostatSetpoints(device NestDevice) (heatCelsius, coolCelsius float64) {
	t := getTrait(device, "sdm.devices.traits.ThermostatTemperatureSetpoint")
	if t == nil {
		return 0, 0
	}
	heatCelsius, _ = t["heatCelsius"].(float64)
	coolCelsius, _ = t["coolCelsius"].(float64)
	return heatCelsius, coolCelsius
}

// GetDisplayName returns the best human-readable name for a device. It checks
// the Info trait's customName first, then falls back to the first parent
// relation's display name (which is the room name in the Google Home app).
func GetDisplayName(device NestDevice) string {
	// Try custom name from Info trait first.
	t := getTrait(device, "sdm.devices.traits.Info")
	if t != nil {
		if name, ok := t["customName"].(string); ok && name != "" {
			return name
		}
	}
	// Fall back to parent relation display name (room name).
	if len(device.ParentRelations) > 0 && device.ParentRelations[0].DisplayName != "" {
		return device.ParentRelations[0].DisplayName
	}
	// Last resort: extract device ID from the resource name.
	parts := strings.Split(device.Name, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "Unknown Nest Device"
}

// SetThermostatMode sets the thermostat to the specified mode.
// Valid modes: HEAT, COOL, HEATCOOL, OFF.
func (c *Client) SetThermostatMode(deviceID, mode string) error {
	return c.ExecuteCommand(deviceID, "sdm.devices.commands.ThermostatMode.SetMode", map[string]interface{}{
		"mode": mode,
	})
}

// SetTemperature sets the thermostat target temperature. It uses SetHeat for
// HEAT mode and SetCool for COOL mode. For HEATCOOL, use SetTemperatureRange instead.
// The temperature is in Celsius because that's what the API uses, and honestly
// I've memorized the conversion formula (multiply by 9/5 and add 32) so I can
// do it in my head faster than most people can find a calculator.
func (c *Client) SetTemperature(deviceID string, tempC float64) error {
	// We need to know the current mode to pick the right command.
	device, err := c.GetDevice(deviceID)
	if err != nil {
		return fmt.Errorf("nest: get device for mode check: %w", err)
	}

	mode := GetThermostatMode(*device)
	switch mode {
	case "HEAT":
		return c.ExecuteCommand(deviceID, "sdm.devices.commands.ThermostatTemperatureSetpoint.SetHeat", map[string]interface{}{
			"heatCelsius": tempC,
		})
	case "COOL":
		return c.ExecuteCommand(deviceID, "sdm.devices.commands.ThermostatTemperatureSetpoint.SetCool", map[string]interface{}{
			"coolCelsius": tempC,
		})
	default:
		return fmt.Errorf("nest: cannot set temperature in mode %q (must be HEAT or COOL)", mode)
	}
}

// SetTemperatureRange sets both heat and cool setpoints for HEATCOOL mode.
func (c *Client) SetTemperatureRange(deviceID string, heatCelsius, coolCelsius float64) error {
	return c.ExecuteCommand(deviceID, "sdm.devices.commands.ThermostatTemperatureSetpoint.SetRange", map[string]interface{}{
		"heatCelsius": heatCelsius,
		"coolCelsius": coolCelsius,
	})
}

// GetCameraLiveStreamURL generates an RTSP live stream URL for a camera device.
// This executes the GenerateRtspStream command and returns the stream URL and
// a token that can be used to extend or stop the stream.
//
// Note: The stream URL expires after 5 minutes. You'll need to call
// ExtendCameraStream before it expires if you want to keep watching.
// It's like how I have to keep renewing my library card every month because
// they don't trust me with the microfiche anymore.
func (c *Client) GetCameraLiveStreamURL(deviceID string) (streamURL string, token string, err error) {
	resp, err := c.do("POST", c.enterprisesPrefix()+"/devices/"+deviceID+":executeCommand", executeCommandRequest{
		Command: "sdm.devices.commands.CameraLiveStream.GenerateRtspStream",
		Params:  map[string]interface{}{},
	})
	if err != nil {
		return "", "", err
	}
	data, err := checkResponse(resp)
	if err != nil {
		return "", "", err
	}

	var result struct {
		Results struct {
			StreamURLs struct {
				RTSPURL string `json:"rtspUrl"`
			} `json:"streamUrls"`
			StreamToken          string `json:"streamToken"`
			StreamExtensionToken string `json:"streamExtensionToken"`
			ExpiresAt            string `json:"expiresAt"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", "", fmt.Errorf("nest: unmarshal stream response: %w", err)
	}

	return result.Results.StreamURLs.RTSPURL, result.Results.StreamExtensionToken, nil
}

// ExtendCameraStream extends an active RTSP stream before it expires.
func (c *Client) ExtendCameraStream(deviceID, extensionToken string) (newToken string, err error) {
	resp, err := c.do("POST", c.enterprisesPrefix()+"/devices/"+deviceID+":executeCommand", executeCommandRequest{
		Command: "sdm.devices.commands.CameraLiveStream.ExtendRtspStream",
		Params: map[string]interface{}{
			"streamExtensionToken": extensionToken,
		},
	})
	if err != nil {
		return "", err
	}
	data, err := checkResponse(resp)
	if err != nil {
		return "", err
	}

	var result struct {
		Results struct {
			StreamExtensionToken string `json:"streamExtensionToken"`
			ExpiresAt            string `json:"expiresAt"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("nest: unmarshal extend response: %w", err)
	}

	return result.Results.StreamExtensionToken, nil
}

// StopCameraStream stops an active RTSP stream.
func (c *Client) StopCameraStream(deviceID, extensionToken string) error {
	return c.ExecuteCommand(deviceID, "sdm.devices.commands.CameraLiveStream.StopRtspStream", map[string]interface{}{
		"streamExtensionToken": extensionToken,
	})
}

// DeviceType constants for the SDM API device types.
const (
	TypeThermostat = "sdm.devices.types.THERMOSTAT"
	TypeCamera     = "sdm.devices.types.CAMERA"
	TypeDoorbell   = "sdm.devices.types.DOORBELL"
	TypeDisplay    = "sdm.devices.types.DISPLAY"
)

// IsThermostat returns true if the device is a Nest thermostat.
func IsThermostat(device NestDevice) bool {
	return device.Type == TypeThermostat
}

// IsCamera returns true if the device is a Nest camera.
func IsCamera(device NestDevice) bool {
	return device.Type == TypeCamera
}

// IsDoorbell returns true if the device is a Nest doorbell.
func IsDoorbell(device NestDevice) bool {
	return device.Type == TypeDoorbell
}

// IsDisplay returns true if the device is a Nest Hub display.
func IsDisplay(device NestDevice) bool {
	return device.Type == TypeDisplay
}
