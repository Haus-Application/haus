---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "orbit-b-hyve-timer"
name: "Orbit B-Hyve Smart Timer"
manufacturer: "Orbit Irrigation Products LLC"
brand: "Orbit"
model: "B-Hyve"
model_aliases: ["B-Hyve XR", "B-Hyve XD", "57946", "57950", "57985", "21004"]
device_type: "sprinkler_controller"
category: "smart_home"
product_line: "B-Hyve"
release_year: 2017
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi", "bluetooth_le"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^bhyve.*", "^orbit.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "bhyve"
  polling_interval_sec: 60
  websocket_event: ""
  setup_type: "password"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["on_off"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "password"
  auth_detail: "Cloud REST API and WebSocket API at api.orbitbhyve.com. Authenticate with Orbit account email/password to receive an Orbit-Session-Token. WebSocket at wss://api.orbitbhyve.com/v1/events for real-time events and commands."
  base_url_template: "https://api.orbitbhyve.com/v1"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "hardwired"
  mounting: "wall"
  indoor_outdoor: "both"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://bhyve.orbitirigation.com"
  api_docs: ""
  developer_portal: ""
  support: "https://support.orbitonline.com"
  community_forum: ""
  image_url: ""
  fcc_id: "YF6-57946"

# --- TAGS ---
tags: ["irrigation", "sprinkler", "outdoor", "cloud-api", "websocket", "bluetooth", "weather-adjustment", "budget", "epa-watersense"]
---

# Orbit B-Hyve Smart Timer

## What It Is

> The Orbit B-Hyve is a WiFi and Bluetooth-connected smart sprinkler controller available in several form factors: indoor wall-mount controllers (4, 8, 12, or 16 zones), outdoor-rated controllers, and hose-end timers. It uses WeatherSense technology to adjust watering schedules based on local weather conditions, soil type, and plant type. The B-Hyve line is positioned as a budget-friendly alternative to Rachio, with similar smart watering features at a lower price point. All models are controlled via the B-Hyve app through Orbit's cloud infrastructure. BLE is available for local setup and control when within range, but there is no documented local IP API.

## How Haus Discovers It

1. **Hostname pattern** -- DHCP hostname may contain `bhyve` or `orbit`
2. **Cloud enrichment** -- After authentication with Orbit cloud, device list includes B-Hyve devices and their local IP addresses
3. **BLE discovery** -- B-Hyve devices advertise over Bluetooth LE when nearby (not usable for IP-based integration)

## Pairing / Authentication

### Cloud API

1. Create or log in to an Orbit B-Hyve account
2. Authenticate via the API:

```
POST https://api.orbitbhyve.com/v1/session
Content-Type: application/json

{
  "session": {
    "email": "user@example.com",
    "password": "password"
  }
}
```

**Response:**
```json
{
  "orbit_session_token": "session-token-string",
  "user_id": "user-uuid",
  "user_name": "John Doe",
  "devices": [...]
}
```

Include `orbit-session-token: {token}` in subsequent request headers.

### BLE (Local Setup)

BLE is used for initial WiFi provisioning and can be used for direct zone control when within Bluetooth range via the B-Hyve app. BLE communication protocol is proprietary and not documented.

### Haus Auth Flow

`POST /api/devices/{ip}/auth` with Orbit account email and password. Haus authenticates against the Orbit cloud and stores the session token.

## API Reference

### Cloud REST API

Base URL: `https://api.orbitbhyve.com/v1`
Auth header: `orbit-session-token: {token}`

**Get devices:**
```
GET /devices
orbit-session-token: {token}
```

**Response:**
```json
[
  {
    "id": "device-uuid",
    "name": "Front Yard Timer",
    "type": "sprinkler_timer",
    "hardware_version": "WT25G2",
    "firmware_version": "0050",
    "is_connected": true,
    "status": {
      "watering_status": "auto",
      "run_mode": "off",
      "rain_delay": 0
    },
    "zones": [
      {
        "station": 1,
        "name": "Front Lawn",
        "enabled": true,
        "smart_watering_enabled": true,
        "num_sprinklers": 4
      }
    ],
    "address": {
      "line_1": "123 Main St",
      "city": "Denver",
      "state": "CO",
      "country": "US"
    }
  }
]
```

**Start a zone:**
```
POST /devices/{device_id}/zones/{zone_number}/start
orbit-session-token: {token}
Content-Type: application/json

{
  "duration": 600
}
```

**Stop watering:**
```
POST /devices/{device_id}/stop
orbit-session-token: {token}
```

**Set rain delay:**
```
PUT /devices/{device_id}
orbit-session-token: {token}
Content-Type: application/json

{
  "rain_delay": 24
}
```

Rain delay in hours.

### WebSocket API (Real-Time)

Connect to `wss://api.orbitbhyve.com/v1/events` for real-time device events:

**Connection:**
```
GET wss://api.orbitbhyve.com/v1/events
orbit-session-token: {token}
```

**Subscribe to device events:**
```json
{
  "event": "app_connection",
  "orbit_session_token": "token",
  "subscribe_device_id": "device-uuid"
}
```

**Event types received:**

| Event | Description |
|-------|-------------|
| `watering_in_progress_notification` | Zone is currently watering |
| `watering_complete` | Zone watering finished |
| `device_connected` | Device came online |
| `device_disconnected` | Device went offline |
| `rain_delay` | Rain delay activated |
| `flow_sensor_state_changed` | Flow sensor triggered (if equipped) |

**Start zone via WebSocket:**
```json
{
  "event": "change_mode",
  "mode": "manual",
  "device_id": "device-uuid",
  "stations": [
    {"station": 1, "run_time": 10}
  ]
}
```

Run time in minutes (not seconds).

## AI Capabilities

> AI integration planned via cloud API. When available:
> - Start/stop individual zones by name
> - Report current watering status
> - Set rain delays
> - Report weather-based skip events
> - Report device online/offline status

## Quirks & Notes

- **No local IP API** -- Unlike Rachio, there is no local HTTP/REST API; all control requires the cloud API or BLE
- **BLE for local control** -- The B-Hyve app can control the timer via BLE when within Bluetooth range (approximately 30 feet); this is the only local control option
- **Cloud API is unofficial** -- The REST and WebSocket APIs at api.orbitbhyve.com are reverse-engineered; Orbit does not officially document them for third-party use
- **Session tokens expire** -- The orbit-session-token expires periodically; re-authenticate when receiving 401 responses
- **WebSocket for real-time** -- The WebSocket API provides real-time events and can also send commands; it is often more reliable than the REST API for zone control
- **Model range** -- B-Hyve spans from simple hose-end timers (1-2 zones, battery/BLE-only) to full 16-zone indoor/outdoor controllers; only WiFi models have cloud connectivity
- **B-Hyve XR** -- The XR model adds extended-range WiFi mesh networking between multiple B-Hyve devices
- **EPA WaterSense** -- Most B-Hyve controllers are EPA WaterSense certified
- **WeatherSense** -- Automatic watering adjustments based on local weather data; available on WiFi models only
- **Budget option** -- B-Hyve controllers are typically 40-60% less expensive than Rachio equivalents
- **Hardwired installation** -- Indoor models wire directly to existing sprinkler valve wiring (24V AC); outdoor models include weatherproof enclosure

## Similar Devices

> - [Rachio 3 Smart Sprinkler](rachio-3-smart-sprinkler.md) -- Premium WiFi sprinkler controller with documented local and cloud APIs
> - [Ring Smart Lighting Pathlight](ring-smart-lighting-pathlight.md) -- Ring outdoor ecosystem device
