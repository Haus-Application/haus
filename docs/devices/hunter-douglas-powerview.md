---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "hunter-douglas-powerview"
name: "Hunter Douglas PowerView Motorization"
manufacturer: "Hunter Douglas"
brand: "Hunter Douglas"
model: "PowerView"
model_aliases: ["PowerView Gen 2", "PowerView Gen 3", "PowerView Automation"]
device_type: "smart_shade"
category: "smart_home"
product_line: "PowerView"
release_year: 2015
discontinued: false
price_range: "$$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "bluetooth_le"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["00:26:74"]
  mdns_services: ["_powerview._tcp"]
  mdns_txt_keys: []
  default_ports: [80]
  signature_ports: [80]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^PowerView-Hub.*", "^PVHUB.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/api/fwversion"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "firmware"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "powerview"
  polling_interval_sec: 15
  websocket_event: ""
  setup_type: "none"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness"]

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 80
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "none"
  auth_detail: "No authentication required for Gen 2 hub. Gen 3 hub uses HTTPS on port 443 with a local API key generated during setup. Gen 2 REST API is completely open on port 80."
  base_url_template: "http://{ip}/api"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.hunterdouglas.com/powerview-motorization"
  api_docs: ""
  developer_portal: ""
  support: "https://www.hunterdouglas.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["shades", "blinds", "window-covering", "local-api", "rest-api", "no-auth", "hub", "powerview", "gen2", "gen3", "position-control", "tilt-control", "scenes"]
---

# Hunter Douglas PowerView Motorization

## What It Is

> Hunter Douglas PowerView is a motorization system for Hunter Douglas window treatments (shades, blinds, shutters, sheers). The system consists of motorized shade/blind units and a PowerView Hub that connects them to the home WiFi network. The hub serves a local REST API on port 80 (Gen 2) or port 443 (Gen 3) that provides full control over shade positions, tilt angles, scenes, rooms, and schedules. The Gen 2 hub requires no authentication whatsoever -- a wide-open JSON REST API. The Gen 3 hub, released in 2023, adds Bluetooth LE, Thread, and Matter support alongside a slightly modified API that requires a local API key. PowerView is one of the most integration-friendly motorized shade systems due to its documented local API.

## How Haus Discovers It

1. **mDNS** -- The PowerView Hub advertises as `_powerview._tcp` on the local network
2. **Port probe** -- HTTP on port 80 (Gen 2) or HTTPS on port 443 (Gen 3)
3. **HTTP fingerprint** -- `GET /api/fwversion` returns firmware version JSON, confirming a PowerView Hub
4. **Hostname pattern** -- DHCP hostname typically starts with `PowerView-Hub` or `PVHUB`
5. **OUI match** -- Hunter Douglas MAC prefix: `00:26:74`

## Pairing / Authentication

### Gen 2 Hub

No authentication required. The REST API on port 80 is completely open to any device on the local network. Connect and send commands immediately.

### Gen 3 Hub

The Gen 3 hub uses HTTPS on port 443 with a locally-generated API key:

1. The API key is generated during initial hub setup via the PowerView app
2. Include the key as a header: `X-API-Key: {key}`
3. The hub uses a self-signed TLS certificate

### Haus Auth Flow

For Gen 2: `POST /api/devices/{ip}/pair` auto-detects the PowerView Hub and begins polling immediately (no auth needed).
For Gen 3: `POST /api/devices/{ip}/auth` with the API key from the PowerView app.

## API Reference

### Gen 2 REST API

Base URL: `http://{hub_ip}/api`

#### Shades

**List all shades:**
```
GET /api/shades
```

**Response:**
```json
{
  "shadeData": [
    {
      "id": 12345,
      "name": "TGl2aW5nIFJvb20=",
      "roomId": 100,
      "groupId": 200,
      "type": 42,
      "batteryStrength": 178,
      "batteryStatus": 3,
      "positions": {
        "posKind1": 1,
        "position1": 32767
      },
      "firmware": {
        "revision": 2,
        "subRevision": 1,
        "build": 916
      }
    }
  ]
}
```

**Note:** Shade names are Base64-encoded.

**Get single shade:**
```
GET /api/shades/{shade_id}
```

**Set shade position:**
```
PUT /api/shades/{shade_id}
Content-Type: application/json

{
  "shade": {
    "positions": {
      "posKind1": 1,
      "position1": 32767
    }
  }
}
```

**Position values:**
- `position1` range: `0` (closed) to `65535` (fully open)
- `posKind1`: `1` = primary position (lift), `2` = secondary position (vane/tilt), `3` = secondary with vane tilt
- Some shades support both lift and tilt (e.g., Silhouette), requiring both `position1`/`posKind1` and `position2`/`posKind2`

**Jog shade (short movement for identification):**
```
PUT /api/shades/{shade_id}
Content-Type: application/json

{
  "shade": {
    "motion": "jog"
  }
}
```

**Calibrate shade:**
```
PUT /api/shades/{shade_id}
Content-Type: application/json

{
  "shade": {
    "motion": "calibrate"
  }
}
```

#### Rooms

**List rooms:**
```
GET /api/rooms
```

**Response:**
```json
{
  "roomData": [
    {
      "id": 100,
      "name": "TGl2aW5nIFJvb20=",
      "order": 0,
      "colorId": 1,
      "iconId": 1
    }
  ]
}
```

#### Scenes

**List scenes:**
```
GET /api/scenes
```

**Activate scene:**
```
GET /api/scenes?sceneId={scene_id}
```

#### Scene Collections (groups of scenes)

**List scene collections:**
```
GET /api/scenecollections
```

**Activate scene collection:**
```
GET /api/scenecollections?sceneCollectionId={id}
```

### Gen 3 API Differences

The Gen 3 hub modifies the API:

- HTTPS on port 443 (self-signed certificate)
- Requires `X-API-Key` header
- Shade names are returned in plaintext (no Base64 encoding)
- Position values use percentage (0-100) instead of 0-65535
- New endpoint structure: `/home/shades/{id}` instead of `/api/shades/{id}`
- Supports SSE (Server-Sent Events) on `/home/shades/events` for real-time position updates

### Shade Types

| Type ID | Description |
|---------|-------------|
| 4 | Roman shade |
| 5 | Top-down/bottom-up (dual position) |
| 6 | Duette |
| 18 | Pirouette |
| 23 | Silhouette |
| 42 | Roller shade |
| 44 | Twist roller |
| 47 | Top-down/bottom-up roller |
| 49 | AC roller shade |
| 51 | Venetian blind |
| 62 | Venetian blind (tilt only) |
| 66 | Duette DuoLite top-down/bottom-up |
| 69 | Curtain (track) |
| 70 | Duette with PowerView |

### Battery Status

| Value | Meaning |
|-------|---------|
| 1 | Low |
| 2 | Medium |
| 3 | Full |
| 4 | Plugged in (AC wired) |

## AI Capabilities

> AI integration planned. When available:
> - Open/close shades by room name
> - Set specific position percentages and tilt angles
> - Activate scenes (e.g., "Movie mode" preset)
> - Report battery levels across all shades
> - Identify shades by jogging
> - Coordinate with lighting for ambiance control

## Quirks & Notes

- **Base64 shade names (Gen 2)** -- All shade and room names in the Gen 2 API are Base64-encoded; decode with standard Base64 before display
- **Position range (Gen 2)** -- Gen 2 uses 0-65535 for positions, not 0-100; convert with `percentage = position1 / 655.35`
- **No authentication (Gen 2)** -- The Gen 2 hub has zero authentication on port 80; anyone on the local network can control all shades
- **Self-signed TLS (Gen 3)** -- The Gen 3 hub's HTTPS certificate is self-signed; TLS verification must be disabled
- **Polling required (Gen 2)** -- Gen 2 has no push mechanism; poll `/api/shades` for position updates. Gen 3 adds Server-Sent Events
- **Battery-powered shades** -- Most PowerView shades run on batteries (D-cell or rechargeable battery wand); wired options available for some shade types
- **Rate limiting** -- The hub has limited processing power; avoid sending more than ~1 request per second
- **Shade movement time** -- Shades take several seconds to reach target position; the position in the API updates only after movement completes
- **Dual-position shades** -- Top-down/bottom-up shades and Silhouettes use both `position1` (lift) and `position2` (tilt/secondary), each with their own `posKind`
- **Scene activation via GET** -- Scenes are activated with a GET request (not POST), which is unusual for a REST API
- **Thread/Matter (Gen 3)** -- Gen 3 hub supports Thread and Matter, enabling direct control from Matter controllers without the REST API

## Similar Devices

> - [Lutron Serena Shades](lutron-serena-shades.md) -- Premium competitor via Lutron bridge/LEAP protocol
> - [IKEA FYRTUR Smart Blinds](ikea-fyrtur-smart-blinds.md) -- Budget Zigbee alternative
