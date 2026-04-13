---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "rachio-3-smart-sprinkler"
name: "Rachio 3 Smart Sprinkler Controller"
manufacturer: "Rachio Inc."
brand: "Rachio"
model: "3"
model_aliases: ["Rachio 3", "8ZULW-C", "16ZULW-C", "Rachio 3e"]
device_type: "sprinkler_controller"
category: "smart_home"
product_line: "Rachio"
release_year: 2018
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["B8:27:EB", "9C:9C:1F"]
  mdns_services: ["_rachio._tcp"]
  mdns_txt_keys: ["id", "ver"]
  default_ports: [8080]
  signature_ports: [8080]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^rachio-.*", "^Rachio-.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8080
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "rachio"
  polling_interval_sec: 60
  websocket_event: ""
  setup_type: "api_key"
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
  auth_method: "api_key"
  auth_detail: "Cloud REST API at api.rach.io uses Bearer token from Rachio account (generated at app.rach.io or via API). Local API on port 8080 allows direct zone control without cloud. Cloud API is well-documented at rachio.readme.io."
  base_url_template: "https://api.rach.io/1/public"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "hardwired"
  mounting: "wall"
  indoor_outdoor: "both"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://rachio.com/products/rachio-3/"
  api_docs: "https://rachio.readme.io/docs"
  developer_portal: "https://rachio.readme.io"
  support: "https://support.rachio.com"
  community_forum: "https://community.rachio.com"
  image_url: ""
  fcc_id: "2AESQ-8ZULW"

# --- TAGS ---
tags: ["irrigation", "sprinkler", "outdoor", "local-api", "cloud-api", "documented-api", "zone-control", "weather-intelligence", "water-savings", "epa-watersense"]
---

# Rachio 3 Smart Sprinkler Controller

## What It Is

> The Rachio 3 is a WiFi-connected smart sprinkler controller that replaces a traditional irrigation timer. Available in 8-zone and 16-zone models, it uses Weather Intelligence Plus to automatically adjust watering schedules based on local weather data, soil type, plant type, sun exposure, and slope. The Rachio 3 has one of the best-documented cloud APIs in the smart home industry (at rachio.readme.io), plus an undocumented local REST API on port 8080 that allows zone control without internet. It supports Alexa, Google Assistant, HomeKit (via firmware update), and IFTTT.

## How Haus Discovers It

1. **mDNS** -- Advertises as `_rachio._tcp` on the local network with TXT records containing device ID and firmware version
2. **Hostname pattern** -- DHCP hostname starts with `rachio-` or `Rachio-`
3. **Port probe** -- HTTP on port 8080 for local API
4. **OUI match** -- MAC prefixes `B8:27:EB`, `9C:9C:1F`

## Pairing / Authentication

### Cloud API Key

1. Log in to the Rachio web app at `https://app.rach.io`
2. Navigate to Account Settings
3. Generate an API key (or use the one displayed)
4. The API key is a Bearer token used in all cloud API requests

Alternatively, use the API to generate a token:
```
POST https://api.rach.io/1/public/person/info
Authorization: Bearer {api_key}
```

### Local API

The local API on port 8080 does not require authentication. Any device on the local network can send commands.

### Haus Auth Flow

`POST /api/devices/{ip}/auth` with the Rachio API key. Haus uses the cloud API for full functionality and falls back to the local API for basic zone control when internet is unavailable.

## API Reference

### Cloud REST API

Base URL: `https://api.rach.io/1/public`
Authorization: `Bearer {api_key}` header on all requests.

**Get person info (includes device IDs):**
```
GET /person/info
Authorization: Bearer {api_key}
```

**Response:**
```json
{
  "id": "person-uuid",
  "username": "user@example.com",
  "fullName": "John Doe",
  "devices": [
    {
      "id": "device-uuid",
      "name": "Backyard Rachio",
      "model": "GENERATION3_16ZONE",
      "serialNumber": "VR0123456",
      "status": "ONLINE",
      "zones": [
        {
          "id": "zone-uuid",
          "zoneNumber": 1,
          "name": "Front Lawn",
          "enabled": true,
          "availableWater": 0.17,
          "rootZoneDepth": 6,
          "efficiency": 0.7
        }
      ]
    }
  ]
}
```

**Get device info:**
```
GET /device/{device_id}
Authorization: Bearer {api_key}
```

**Start a zone:**
```
PUT /zone/start
Authorization: Bearer {api_key}
Content-Type: application/json

{
  "id": "zone-uuid",
  "duration": 600
}
```

Duration is in seconds. Maximum 10800 (3 hours).

**Start multiple zones:**
```
PUT /zone/start_multiple
Authorization: Bearer {api_key}
Content-Type: application/json

{
  "zones": [
    {"id": "zone-1-uuid", "duration": 600, "sortOrder": 0},
    {"id": "zone-2-uuid", "duration": 300, "sortOrder": 1}
  ]
}
```

**Stop watering:**
```
PUT /device/stop_water
Authorization: Bearer {api_key}
Content-Type: application/json

{
  "id": "device-uuid"
}
```

**Put device in standby (rain delay):**
```
PUT /device/rain_delay
Authorization: Bearer {api_key}
Content-Type: application/json

{
  "id": "device-uuid",
  "duration": 86400
}
```

Duration in seconds. Use `0` to cancel a rain delay.

**Get current schedule:**
```
GET /device/{device_id}/current_schedule
Authorization: Bearer {api_key}
```

**Response (when running):**
```json
{
  "deviceId": "device-uuid",
  "scheduleId": "schedule-uuid",
  "status": "PROCESSING",
  "zoneId": "zone-uuid",
  "zoneName": "Front Lawn",
  "zoneNumber": 1,
  "startDate": 1712000000000,
  "duration": 600,
  "remaining": 420
}
```

**Webhook events** (configured via API):

| Event Type | Description |
|-----------|-------------|
| `DEVICE_STATUS_EVENT` | Device online/offline |
| `ZONE_STATUS_EVENT` | Zone started/stopped/completed |
| `SCHEDULE_STATUS_EVENT` | Schedule started/stopped |
| `RAIN_DELAY_EVENT` | Rain delay activated/deactivated |
| `WEATHER_INTELLIGENCE_EVENT` | Weather skip triggered |

### Local REST API

Base URL: `http://{ip}:8080`
No authentication required.

**Start a zone (local):**
```
PUT http://{ip}:8080/zone/start
Content-Type: application/json

{
  "zoneNumber": 1,
  "duration": 600
}
```

**Stop all zones (local):**
```
PUT http://{ip}:8080/device/stop
```

**Get device status (local):**
```
GET http://{ip}:8080/device/status
```

The local API supports basic zone start/stop but does not provide full schedule management, weather intelligence data, or zone configuration.

## AI Capabilities

> AI integration planned. When available:
> - Start/stop individual zones by name
> - Run custom watering sequences
> - Report current watering status (which zone, time remaining)
> - Set rain delays
> - Report weather intelligence skips
> - Report zone moisture levels and next scheduled watering

## Quirks & Notes

- **Cloud API is excellent** -- Rachio has one of the best-documented smart home APIs at rachio.readme.io; rate limited to 1700 calls per day
- **Local API is undocumented** -- The local API on port 8080 exists but is not officially supported; it may change without notice in firmware updates
- **No local schedule management** -- The local API only supports starting/stopping zones; schedules, weather intelligence, and zone configuration require the cloud API
- **Weather Intelligence Plus** -- Automatically skips or adjusts watering based on observed and forecasted rain, wind, freeze, and saturation levels
- **HomeKit support** -- Rachio 3 gained HomeKit support via firmware update; exposed as a sprinkler accessory
- **Zone count** -- Available in 8-zone (8ZULW) and 16-zone (16ZULW) models; supports master valve/pump relay
- **Flow meter support** -- Optional wireless flow meter detects high flow (broken pipe/head) and low flow (clogged nozzle)
- **EPA WaterSense certified** -- Certified to save water through weather-based smart scheduling
- **Webhook support** -- Cloud API supports webhook subscriptions for real-time event notifications
- **Rachio 3e** -- A lower-cost variant with fewer features (no Weather Intelligence Plus, no flow meter support)
- **Wiring** -- Replaces existing sprinkler controller; uses standard 24V AC sprinkler valve wiring
- **Outdoor-rated** -- Weather-resistant for outdoor installation, but Hunter recommends a weatherproof enclosure

## Similar Devices

> - [Orbit B-Hyve Smart Timer](orbit-b-hyve-timer.md) -- Budget WiFi/BLE competitor
> - [Ring Smart Lighting Pathlight](ring-smart-lighting-pathlight.md) -- Ring outdoor ecosystem device
