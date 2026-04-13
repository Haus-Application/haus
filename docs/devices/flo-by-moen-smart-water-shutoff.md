---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "flo-by-moen-smart-water-shutoff"
name: "Flo by Moen Smart Water Monitor & Shutoff"
manufacturer: "Moen Incorporated (Fortune Brands Innovations)"
brand: "Flo by Moen"
model: "Smart Water Monitor & Shutoff"
model_aliases: ["Flo Smart Water Monitor", "Flo 1.25", "Flo 1.0", "900-001", "900-002", "Flo by Moen 3/4\"", "Flo by Moen 1\"", "Flo by Moen 1-1/4\""]
device_type: "water_shutoff"
category: "smart_home"
product_line: "Flo by Moen"
release_year: 2019
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["44:67:55", "90:97:D5"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^flo-.*", "^Flo-.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "flo"
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
  auth_detail: "Cloud REST API at api.meetflo.com. Authenticate with Flo/Moen account email and password. Returns a bearer token. The device itself has no local API -- all communication goes through Flo cloud servers."
  base_url_template: "https://api-gw.meetflo.com/api/v2"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.moen.com/flo"
  api_docs: ""
  developer_portal: ""
  support: "https://support.meetflo.com"
  community_forum: ""
  image_url: ""
  fcc_id: "2ASYB-FLODVC100"

# --- TAGS ---
tags: ["water", "leak-detection", "shutoff-valve", "flow-monitoring", "cloud-only", "water-pressure", "water-temperature", "insurance-discount", "whole-home"]
---

# Flo by Moen Smart Water Monitor & Shutoff

## What It Is

> The Flo by Moen Smart Water Monitor & Shutoff is a whole-home water monitoring and automatic shutoff device installed on the main water supply line. It continuously monitors water flow rate, water pressure, and water temperature, using machine learning to detect leaks (from small drips to burst pipes) and automatically shut off the water supply to prevent damage. The device runs a daily "MicroLeak Test" (FloProtect Health Test) that pressurizes the plumbing system and detects leaks as small as a drop per minute. Available in 3/4", 1", and 1-1/4" pipe sizes, it installs on the main water line after the shutoff valve. All monitoring and control is through the Flo by Moen cloud service -- there is no local API.

## How Haus Discovers It

1. **OUI match** -- Flo MAC prefixes: `44:67:55`, `90:97:D5`
2. **Hostname pattern** -- DHCP hostname may start with `flo-` or `Flo-`
3. **No local ports** -- The device does not expose any local network services
4. **Cloud enrichment** -- After authenticating with the Flo cloud API, device details including IP address are available

## Pairing / Authentication

### Cloud API Authentication

```
POST https://api-gw.meetflo.com/api/v2/auth/signin
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password"
}
```

**Response:**
```json
{
  "token": "eyJ...",
  "tokenPayload": {
    "user": {
      "user_id": "user-uuid",
      "email": "user@example.com"
    },
    "timestamp": 1712000000
  },
  "tokenExpiration": 86400
}
```

Include in subsequent requests:
```
Authorization: Bearer {token}
```

### Haus Auth Flow

`POST /api/devices/{ip}/auth` with Flo account email and password. Haus authenticates against the Flo cloud and stores the bearer token.

## API Reference

### Cloud REST API

Base URL: `https://api-gw.meetflo.com/api/v2`
Auth: `Authorization: Bearer {token}`

**Get user info and locations:**
```
GET /users/{user_id}?expand=locations
Authorization: Bearer {token}
```

**Response:**
```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "locations": [
    {
      "id": "location-uuid",
      "nickname": "Home",
      "address": "123 Main St",
      "devices": [
        {
          "id": "device-uuid",
          "macAddress": "44:67:55:XX:XX:XX",
          "nickname": "Main Water Line",
          "serialNumber": "FLO-XXXX-XXXX",
          "fwVersion": "7.0.1",
          "deviceType": "flo_device_v2",
          "deviceModel": "flo_device_075_v2",
          "isConnected": true,
          "valve": {
            "lastKnown": "open",
            "target": "open"
          },
          "telemetry": {
            "current": {
              "gpm": 2.5,
              "psi": 58.3,
              "tempF": 65.2,
              "updated": "2024-04-01T12:00:00Z"
            }
          },
          "notifications": {
            "pending": {
              "criticalCount": 0,
              "warningCount": 1,
              "infoCount": 3
            }
          }
        }
      ]
    }
  ]
}
```

**Get device details:**
```
GET /devices/{device_id}
Authorization: Bearer {token}
```

**Close valve (shut off water):**
```
POST /devices/{device_id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "valve": {
    "target": "closed"
  }
}
```

**Open valve (restore water):**
```
POST /devices/{device_id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "valve": {
    "target": "open"
  }
}
```

**Get water consumption data:**
```
GET /water/consumption?macAddress={mac}&startDate=2024-04-01&endDate=2024-04-02&interval=1h
Authorization: Bearer {token}
```

**Response:**
```json
{
  "params": {
    "startDate": "2024-04-01",
    "endDate": "2024-04-02",
    "interval": "1h"
  },
  "aggregations": {
    "sumTotalGallonsConsumed": 85.4
  },
  "items": [
    {
      "time": "2024-04-01T00:00:00Z",
      "gallonsConsumed": 0.2
    },
    {
      "time": "2024-04-01T01:00:00Z",
      "gallonsConsumed": 0.0
    }
  ]
}
```

**Run health test (MicroLeak Test):**
```
POST /devices/{device_id}/healthTest/run
Authorization: Bearer {token}
```

**Get health test results:**
```
GET /devices/{device_id}/healthTest/latest
Authorization: Bearer {token}
```

**Response:**
```json
{
  "status": "passed",
  "roundId": "test-uuid",
  "startDate": "2024-04-01T03:00:00Z",
  "endDate": "2024-04-01T03:05:00Z",
  "leakLossMinGal": 0.0,
  "leakLossMaxGal": 0.0,
  "startPsi": 58.3,
  "endPsi": 58.2,
  "deltaPsi": -0.1
}
```

### Key Telemetry Fields

| Field | Description | Unit |
|-------|-------------|------|
| `gpm` | Current flow rate | Gallons per minute |
| `psi` | Current water pressure | PSI |
| `tempF` | Current water temperature | Fahrenheit |

### Alerts & Notifications

```
GET /alerts?deviceId={device_id}&status=triggered&page=1&size=20
Authorization: Bearer {token}
```

Alert severity levels: `critical` (shutoff triggered), `warning` (abnormal pattern), `info` (informational)

## AI Capabilities

> AI integration planned via cloud API. When available:
> - Report real-time water flow rate, pressure, and temperature
> - Open/close the main water shutoff valve
> - Report daily/weekly/monthly water consumption
> - Report health test results (leak detection)
> - Alert on detected leaks with severity
> - Report valve status (open/closed)
> - The AI speaks as the device: "Water pressure is 58 PSI. I'm seeing 2.5 GPM of flow right now. Last night's health test passed -- no leaks detected."

## Quirks & Notes

- **Cloud-only** -- No local API; all monitoring and control requires the Flo cloud service and internet connectivity
- **Valve operation** -- The motorized ball valve takes approximately 5-10 seconds to fully open or close; the API response returns before the valve has finished moving
- **MicroLeak Test** -- The daily health test (typically 3 AM) briefly closes the valve to pressurize the system; this may cause a brief interruption in water flow
- **Professional installation recommended** -- Installation on the main water line typically requires a licensed plumber; improper installation can void warranty and cause water damage
- **Pipe sizes** -- Available in 3/4", 1", and 1-1/4" sizes; must match the main water supply pipe diameter
- **Power requirement** -- Requires a standard outlet; includes a battery backup for valve operation during power outages
- **Insurance discounts** -- Many insurance companies offer premium discounts (5-15%) for homes with automatic water shutoff systems
- **FloProtect subscription** -- Advanced features (automatic shutoff, detailed alerts, insurance program) require a FloProtect subscription ($5/month or $40/year)
- **Token expiry** -- Bearer tokens expire after 24 hours; re-authenticate when receiving 401 responses
- **Water pressure** -- Flo recommends water pressure between 30-80 PSI; pressures outside this range generate alerts
- **Freeze alerts** -- Alerts when water temperature approaches freezing (below 39F/4C)
- **API undocumented** -- The meetflo.com API is reverse-engineered; Moen does not officially support third-party API access

## Similar Devices

> - [Pentair ScreenLogic](pentair-screenlogic.md) -- Pool water management system (different application, similar home infrastructure)
