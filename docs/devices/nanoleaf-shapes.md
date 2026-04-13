---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "nanoleaf-shapes"
name: "Nanoleaf Shapes"
manufacturer: "Nanoleaf"
brand: "Nanoleaf"
model: "NL42"
model_aliases: ["NL42", "NL47", "NL52", "Shapes Triangles", "Shapes Mini Triangles", "Shapes Hexagons"]
device_type: "nanoleaf_panel"
category: "lighting"
product_line: "Nanoleaf Shapes"
release_year: 2020
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["00:55:DA", "80:E4:DA"]
  mdns_services: ["_nanoleafapi._tcp"]
  mdns_txt_keys: ["srcvers", "md", "id", "NL-DeviceId"]
  default_ports: [16021]
  signature_ports: [16021]
  ssdp_search_target: "nanoleaf_aurora:light"
  ssdp_server_string: ""
  hostname_patterns: ["Nanoleaf.*", "nanoleaf.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 16021
    path: "/api/v1/new"
    method: "GET"
    expect_status: 401
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "nanoleaf"
  polling_interval_sec: 5
  websocket_event: "nanoleaf:state"
  setup_type: "api_key"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp", "scenes"]

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 16021
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "Authorization token obtained by holding power button for 5-7 seconds then POSTing to /api/v1/new; token passed as path segment: /api/v1/{auth_token}/"
  base_url_template: "http://{ip}:16021/api/v1/{auth_token}"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "panel"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://nanoleaf.me/en-US/products/shapes/"
  api_docs: "https://forum.nanoleaf.me/docs"
  developer_portal: "https://forum.nanoleaf.me/docs"
  support: "https://nanoleaf.me/en-US/support/"
  community_forum: "https://forum.nanoleaf.me/"
  image_url: ""
  fcc_id: "2AWLPNL42"

# --- TAGS ---
tags: ["wifi", "modular_panels", "decorative", "effects", "local_rest_api", "no_hub"]
---

# Nanoleaf Shapes

## What It Is

Nanoleaf Shapes are modular, wall-mounted LED light panels that come in three form factors: Triangles, Mini Triangles, and Hexagons. They snap together magnetically and link via electrical connectors to form custom layouts. Each panel is individually addressable, supporting 16 million colors and tunable white (1200K-6500K). They connect directly to WiFi with no hub required, and expose a full local REST API on port 16021 that makes them an excellent candidate for Haus integration. The panels also support touch gestures, music visualization (via built-in microphone), and a rich scene/effect engine called "Canvas" with community-shared effects.

## How Haus Discovers It

1. **mDNS Discovery**: Nanoleaf devices advertise `_nanoleafapi._tcp.local.` via multicast DNS. The TXT records include `md` (model description, e.g. "NL42"), `id` (unique device ID), `srcvers` (API version), and `NL-DeviceId`.
2. **SSDP Discovery**: Nanoleaf also responds to SSDP M-SEARCH with search target `nanoleaf_aurora:light` (legacy name from the original Aurora product).
3. **OUI Match**: MAC addresses beginning with `00:55:DA` or `80:E4:DA` are associated with Nanoleaf devices.
4. **Port Probe**: HTTP GET to port 16021 at `/api/v1/new` returns 401 Unauthorized if no auth token exists, confirming the device is a Nanoleaf.
5. **Device Info**: Once authenticated, GET `/api/v1/{token}/` returns full device info including model, firmware version, serial number, and panel layout.

## Pairing / Authentication

Nanoleaf uses a physical-confirmation authentication flow:

1. User holds the power button on the Nanoleaf controller for 5-7 seconds until the LEDs begin flashing in a pattern.
2. While the device is in pairing mode (within 30 seconds), Haus sends: `POST http://{ip}:16021/api/v1/new`
3. The device responds with a JSON object containing the auth token: `{"auth_token": "AbCdEfGh1234"}`
4. Haus stores this token and uses it as a path segment in all subsequent API calls: `/api/v1/{auth_token}/...`
5. The token does not expire unless the user performs a factory reset on the device.

## API Reference

### Base URL

`http://{ip}:16021/api/v1/{auth_token}`

### Authentication

All endpoints (except `POST /api/v1/new`) require the auth token as a URL path segment.

### Device Info

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Full device info (name, serial, model, firmware, state, effects, layout) |
| GET | `/state` | Current state (on, brightness, hue, sat, ct, colorMode) |
| GET | `/effects` | Current effect and effect list |

### Power

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| PUT | `/state` | `{"on": {"value": true}}` | Turn on |
| PUT | `/state` | `{"on": {"value": false}}` | Turn off |

### Brightness

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| PUT | `/state` | `{"brightness": {"value": 80}}` | Set brightness (0-100) |
| PUT | `/state` | `{"brightness": {"value": 50, "duration": 3}}` | Set with transition (seconds) |

### Color

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| PUT | `/state` | `{"hue": {"value": 120}, "sat": {"value": 100}}` | Set HSV color |
| PUT | `/state` | `{"ct": {"value": 4000}}` | Set color temperature (1200-6500K) |

### Effects / Scenes

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| GET | `/effects/effectsList` | — | List available effects |
| PUT | `/effects` | `{"select": "Flames"}` | Activate a named effect |
| PUT | `/effects` | `{"write": {...}}` | Write a custom streaming effect |

### Panel Layout

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/panelLayout/layout` | Returns panel positions, orientations, and IDs |
| GET | `/panelLayout/globalOrientation` | Layout rotation |

### Streaming (External Control)

For real-time per-panel color control:

1. `PUT /effects` with body: `{"write": {"command": "display", "animType": "extControl", "extControlVersion": "v2"}}`
2. Device opens a UDP socket (port returned in response, typically 60222).
3. Send UDP packets with panel color data in binary format.

#### Streaming UDP Packet Format (v2)

| Offset | Length | Field |
|--------|--------|-------|
| 0 | 2 | Number of panels (big-endian) |
| 2+ | 8 per panel | Panel ID (2 bytes) + R (1) + G (1) + B (1) + White (1) + Transition time (2, deciseconds) |

### Server-Sent Events

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/events?id=1,2,3,4` | SSE stream: 1=state, 2=layout, 3=effects, 4=touch |

Touch events deliver panel ID and gesture type (tap, double-tap, swipe up/down/left/right).

## AI Capabilities

When the AI concierge is chatting with a Nanoleaf device, it can:
- Turn panels on/off
- Set brightness with optional transition
- Set a solid color across all panels (HSV or color temperature)
- Activate named scenes/effects from the effect list
- Report the current state, active effect, and panel layout
- Describe the physical layout (number of panels, arrangement)

## Quirks & Notes

- **Legacy "Aurora" Naming**: The SSDP search target still uses `nanoleaf_aurora:light` even for Shapes products. The original product was called "Aurora" and some internal naming persists.
- **Auth Token in URL**: The auth token is embedded in the URL path, not in a header. This is unusual and means URLs are sensitive.
- **Panel Layout Coordinates**: Panel positions are returned in centimeters relative to an origin. The coordinate system can be rotated via globalOrientation. This is essential for rendering an accurate visual representation.
- **Streaming Mode Timeout**: External control (UDP streaming) mode times out after 60 seconds of no data. The mode must be re-activated.
- **Firmware API Differences**: Older firmware (pre-8.x) uses API v1 with slightly different endpoint behaviors. Shapes generally ship with firmware that supports all v1 features.
- **Max Panels**: A single controller supports up to 500 panels (though practical WiFi bandwidth limits effective per-panel streaming to around 70-80 panels).
- **Touch Events**: Touch gesture detection requires firmware 8.4+ and is only available on Shapes (not older Canvas/Aurora products).

## Similar Devices

- **nanoleaf-essentials-a19** — Thread/Matter bulb from Nanoleaf, different protocol entirely
- **nanoleaf-canvas** — Older square panels, same API but different form factor
- **nanoleaf-elements** — Wood-look hexagons, same API, warm white only
- **lifx-a19-color** — WiFi bulb with local control, different protocol
- **govee-rgbic-led-strip** — Addressable LED strip with LAN API, different form factor
