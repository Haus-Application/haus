---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ring-smart-lighting-pathlight"
name: "Ring Smart Lighting Pathlight"
manufacturer: "Ring LLC (Amazon)"
brand: "Ring"
model: "Smart Lighting Pathlight"
model_aliases: ["5AT1S6-BEN0", "Ring Pathlight", "Ring Solar Pathlight"]
device_type: "outdoor_light"
category: "lighting"
product_line: "Ring Smart Lighting"
release_year: 2019
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
  protocols_spoken: ["proprietary_rf"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "ring"
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "motion"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Ring Smart Lighting devices connect to a Ring Bridge via proprietary sub-GHz RF protocol. The Ring Bridge connects to Ring cloud via WiFi. Control is via the unofficial Ring cloud API. No local API available. Uses the same auth flow as other Ring devices (oauth.ring.com/oauth/token with 2FA)."
  base_url_template: "https://api.ring.com/clients_api"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "outdoor_fixture"
  power_source: "battery"
  mounting: "outdoor"
  indoor_outdoor: "outdoor"
  wireless_radios: []

# --- LINKS ---
links:
  product_page: "https://ring.com/products/smart-lighting-pathlight"
  api_docs: ""
  developer_portal: ""
  support: "https://support.ring.com"
  community_forum: "https://community.ring.com"
  image_url: ""
  fcc_id: "2AEUO-0128"

# --- TAGS ---
tags: ["outdoor-lighting", "pathlight", "ring", "amazon", "cloud-only", "ring-bridge", "sub-ghz-rf", "battery", "solar-option", "motion-sensor", "no-local-api"]
---

# Ring Smart Lighting Pathlight

## What It Is

> The Ring Smart Lighting Pathlight is a battery-powered outdoor path light with an integrated motion sensor. Part of Ring's Smart Lighting ecosystem, it communicates via a proprietary sub-GHz RF protocol (not WiFi, not Zigbee) to a Ring Bridge, which then connects to Ring's cloud servers via WiFi. The pathlight provides 80 lumens of warm white light (adjustable brightness), motion detection up to approximately 15 feet, and can be grouped with other Ring Smart Lighting devices and Ring cameras for coordinated motion-triggered lighting. Available in battery-powered and solar-powered variants, the pathlight is designed for walkways, driveways, and garden areas.

## How Haus Discovers It

> Ring Smart Lighting Pathlights are not directly discoverable on the IP network because they use proprietary RF, not WiFi:
>
> 1. **Discover the Ring Bridge** -- The Ring Bridge may be identifiable on the network by Ring MAC OUI prefixes (`5C:47:5E`, `34:3E:A4`, `0C:47:C9`) or hostname patterns (`ring-*`)
> 2. **Query Ring cloud** -- After authenticating with Ring's cloud API, query for all devices including Smart Lighting devices connected through the bridge
> 3. **Device identification** -- Pathlights appear in the Ring API device list with a specific device type and associated bridge ID

## Pairing / Authentication

### Pathlight to Ring Bridge

1. The Ring Bridge must be set up first via the Ring app (connects to WiFi)
2. In the Ring app, select "Set Up a Device" > "Smart Lighting"
3. Pull the battery tab on the pathlight (or insert batteries)
4. The Ring app discovers the pathlight via the bridge's RF scanning
5. Assign to a group and configure motion settings

### Ring Cloud API Auth

Uses the same unofficial Ring OAuth2 flow as all Ring devices:

1. `POST https://oauth.ring.com/oauth/token` with email, password, and 2FA code
2. Receive access and refresh tokens
3. Include `Authorization: Bearer {access_token}` in API requests

See [Ring Indoor Cam](ring-indoor-cam.md) for the full authentication flow.

### Haus Auth Flow

`POST /api/devices/{ip}/auth` with Ring account credentials and 2FA code. Haus stores the refresh token for automatic token renewal.

## API Reference

### Ring Cloud API (Unofficial)

Base URL: `https://api.ring.com/clients_api`
Auth: `Authorization: Bearer {access_token}`

**Get all devices (including Smart Lighting):**
```
GET /ring_devices
Authorization: Bearer {token}
```

**Smart Lighting devices appear in the `stickup_cams` array (confusingly, Ring uses this for all devices):**
```json
{
  "stickup_cams": [
    {
      "id": 12345678,
      "description": "Front Walkway",
      "device_id": "abcdef123456",
      "kind": "lpd_v2",
      "firmware_version": "3.8.67",
      "led_status": "on",
      "ring_id": null,
      "location_id": "location-uuid",
      "features": {
        "motions_enabled": true,
        "show_recordings": false
      },
      "owned": true,
      "alerts": {
        "connection": "online"
      }
    }
  ]
}
```

**Turn light on:**
```
PUT /doorbots/{device_id}/floodlight_light_on
Authorization: Bearer {token}
```

**Turn light off:**
```
PUT /doorbots/{device_id}/floodlight_light_off
Authorization: Bearer {token}
```

**Set brightness (0-100):**
```
PUT /doorbots/{device_id}/settings
Authorization: Bearer {token}
Content-Type: application/json

{
  "light_settings": {
    "brightness": 75
  }
}
```

**Set motion sensitivity:**
```
PUT /doorbots/{device_id}/settings
Authorization: Bearer {token}
Content-Type: application/json

{
  "motion_settings": {
    "motion_detection_enabled": true,
    "motion_sensitivity": "mid"
  }
}
```

Sensitivity values: `low`, `mid`, `high`, `highest`

**Get motion events:**
```
GET /doorbots/{device_id}/history?limit=20
Authorization: Bearer {token}
```

## AI Capabilities

> AI integration is not planned for V1 due to cloud-only API. If implemented:
> - Turn pathlights on/off by group name
> - Adjust brightness levels
> - Report motion events and activity history
> - Configure motion sensitivity zones
> - Coordinate with Ring cameras for motion-triggered actions

## Quirks & Notes

- **Ring Bridge required** -- Pathlights cannot operate independently; the Ring Bridge (sold separately or bundled) provides the WiFi-to-RF gateway
- **Proprietary sub-GHz RF** -- Not WiFi, Zigbee, or Z-Wave; uses Ring's proprietary sub-GHz frequency (around 900 MHz in the US) for longer range and better outdoor penetration
- **No local API** -- All control routes through Ring cloud servers; no local control is possible even through the bridge
- **Linked motion groups** -- Multiple pathlights and Ring cameras can be linked in a "Light Group" so that motion detected on any device triggers all lights in the group
- **Battery or solar** -- Standard model uses 4x D-cell batteries (approx. 1 year life); solar variant includes an integrated solar panel that keeps batteries topped up
- **80 lumens** -- Relatively dim; designed for ambient pathway lighting, not security floodlighting
- **Motion range** -- PIR motion sensor detects movement up to approximately 15 feet at a 120-degree angle
- **Ring Bridge capacity** -- A single Ring Bridge supports up to 50 Smart Lighting devices
- **RF range** -- Sub-GHz RF range is approximately 250 feet line-of-sight from the bridge; walls and obstacles reduce range
- **Unofficial API** -- The Ring API is reverse-engineered and undocumented; Amazon may change or restrict it at any time
- **2FA mandatory** -- Ring enforces two-factor authentication on all accounts, which complicates automated API access (requires manual 2FA code entry during initial setup)
- **Weather resistant** -- IP65 rated for outdoor use in rain, snow, and temperature extremes

## Similar Devices

> - [Ring Floodlight Cam](ring-floodlight-cam.md) -- Ring outdoor camera with integrated floodlights
> - [Rachio 3 Smart Sprinkler](rachio-3-smart-sprinkler.md) -- Another outdoor smart home device with better API access
> - [Orbit B-Hyve Smart Timer](orbit-b-hyve-timer.md) -- Outdoor irrigation with cloud API
