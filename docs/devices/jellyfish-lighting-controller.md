---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "jellyfish-lighting-controller"
name: "JellyFish Lighting Controller"
manufacturer: "JellyFish Lighting"
brand: "JellyFish"
model: "JellyFish Controller V2"
model_aliases: []
device_type: "jellyfish_controller"
category: "lighting"
product_line: "JellyFish"
release_year: 2020
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "ethernet"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["D8:80:39", "DC:4F:22"]
  mdns_services: ["_jellyfishV2._tcp"]
  mdns_txt_keys: []
  default_ports: [9000, 8080]
  signature_ports: [9000]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^jellyfish.*", "^jf-.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8080
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: "JellyFish"
    server_header: ""
    body_contains: "JellyFish"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "jellyfish"
  polling_interval_sec: 0
  websocket_event: "jellyfish:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: ["on_off", "scenes"]

# --- PROTOCOL ---
protocol:
  type: "websocket_json"
  port: 9000
  transport: "WebSocket"
  encoding: "JSON"
  auth_method: "none"
  auth_detail: "No authentication required. Direct local WebSocket connection."
  base_url_template: "ws://{ip}:9000/"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "hardwired"
  mounting: "wall"
  indoor_outdoor: "outdoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://jellyfishlighting.com"
  api_docs: ""
  developer_portal: ""
  support: "https://jellyfishlighting.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["outdoor-lighting", "led", "addressable", "zones", "patterns", "websocket", "local-control", "no-auth"]
---

# JellyFish Lighting Controller

## What It Is

> JellyFish Lighting is a permanent outdoor LED lighting system professionally installed along rooflines, soffits, and architectural features. The controller manages individually addressable LED strips organized into zones, each with a configurable pixel count. Users select from a library of color patterns (holidays, accents, animations) that can be applied per-zone. The controller exposes a fully asynchronous WebSocket API on port 9000 for local control with no authentication required, and a firmware update web UI on port 8080.

## How Haus Discovers It

1. **mDNS** -- Advertises as `_jellyfishV2._tcp` on the local network
2. **Port probe** -- TCP check on port 9000 confirms WebSocket availability
3. **WebSocket probe** -- Connect to `ws://{ip}:9000/` and send `{"cmd":"toCtlrGet","get":[["ctlrName"]]}` to retrieve the controller name
4. **Zone enumeration** -- Query zones and pattern list to populate device capabilities

## Pairing / Authentication

> No pairing or authentication is required. The WebSocket API on port 9000 is open to any device on the local network. Connect directly and start sending commands.

## API Reference

### Protocol: WebSocket JSON

The controller exposes a fully asynchronous WebSocket API. Commands are sent as JSON objects and responses arrive independently. The connection endpoint is `ws://{controller_ip}:9000/`.

- **Handshake timeout:** 5 seconds
- **Message format:** JSON objects
- **Command types:** `toCtlrGet` for queries, `toCtlrSet` for control
- **Response format:** Responses echo the requested data type as the top-level key

### Get Controller Name

**Send:**
```json
{"cmd": "toCtlrGet", "get": [["ctlrName"]]}
```

**Response:**
```json
{"ctlrName": "Happy lites"}
```

### Get Zones

**Send:**
```json
{"cmd": "toCtlrGet", "get": [["zones"]]}
```

**Response:**
```json
{
  "zones": {
    "Zone": {"numPixels": 300},
    "Zone1": {"numPixels": 150}
  }
}
```

Each zone represents a physical LED strip segment with a pixel count.

### Get Pattern List

**Send:**
```json
{"cmd": "toCtlrGet", "get": [["patternFileList"]]}
```

**Response:**
```json
{
  "patternFileList": [
    {"folders": "Easter", "name": "Easter Colors", "readOnly": true},
    {"folders": "Accent", "name": "All Lights Warm White 3000K", "readOnly": false},
    {"folders": "Holiday", "name": "Christmas", "readOnly": true}
  ]
}
```

Pattern file paths are constructed as `{folders}/{name}` (e.g., `"Accent/All Lights Warm White 3000K"`).

- `readOnly: true` -- built-in system patterns that cannot be modified
- `readOnly: false` -- user-created or editable patterns

### Get Zone State

Query the current run state of one or more zones:

**Send:**
```json
{"cmd": "toCtlrGet", "get": [["runPattern", "Zone1", "Zone"]]}
```

**Response (one per zone):**
```json
{
  "runPattern": {
    "state": 1,
    "zoneName": ["Zone1"],
    "file": "Accent/White",
    "id": ""
  }
}
```

- `state: 1` = playing
- `state: 0` = off
- `file` = currently active pattern file path
- `zoneName` = zones this state applies to

### Play a Pattern (Turn On)

**Send:**
```json
{
  "cmd": "toCtlrSet",
  "runPattern": {
    "state": 1,
    "zoneName": ["Zone1", "Zone"],
    "file": "Accent/All Lights Warm White 3000K",
    "id": "",
    "data": ""
  }
}
```

- `state: 1` = start playing
- `zoneName` = array of zone names to activate
- `file` = pattern file path (`folder/name`)
- `id` and `data` = pass empty strings

### Turn Off

**Send:**
```json
{
  "cmd": "toCtlrSet",
  "runPattern": {
    "state": 0,
    "zoneName": ["Zone1", "Zone"],
    "file": "",
    "id": "",
    "data": ""
  }
}
```

- `state: 0` = stop playing / turn off
- `file` = pass empty string when turning off

### Firmware Update UI

The controller serves a web UI on **port 8080** for firmware updates. This is an HTML page using jQuery for checking and applying firmware updates. It is not part of the control API but is useful for device identification via HTTP fingerprinting.

## AI Capabilities

> When chatting with a JellyFish controller, the AI can:
> - **Query current zone states** -- on/off, active pattern per zone
> - **Turn zones on/off** -- activate or deactivate specific zones
> - **Play specific patterns** -- set any available pattern on any combination of zones
> - **List available patterns and zones** -- enumerate all pattern files and zone configurations
> - **Real-time queries** -- uses a direct WebSocket connection for instant responses
>
> The AI speaks as the device: "I have 2 zones running. Zone is playing 'Accent/All Lights Warm White 3000K' with 300 pixels."

## Quirks & Notes

- **No polling needed** -- the WebSocket API is stateless and query-on-demand; there is no need for periodic polling
- **Asynchronous responses** -- responses arrive independently from requests; the client must match responses to requests by data type key
- **Zone names are arbitrary** -- zones are named during installation and can be renamed; do not hardcode zone names
- **Pattern paths** -- always construct pattern file paths as `{folders}/{name}` from the `patternFileList` response
- **Professional installation** -- JellyFish systems are professionally installed; the controller is typically mounted outdoors in a weatherproof enclosure near the electrical panel
- **Multiple controllers** -- a home may have multiple JellyFish controllers (e.g., front and back of house); each has its own IP and WebSocket endpoint
- **Firmware updates** -- managed via the web UI on port 8080, not via the WebSocket API

## Similar Devices

> JellyFish is a unique product category. No directly comparable devices in the Haus knowledge base. Conceptually similar to addressable LED controllers (WLED, Govee) but with a proprietary protocol and professional installation model.
