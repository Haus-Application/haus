---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "twinkly-smart-lights"
name: "Twinkly Smart LED Lights"
manufacturer: "Ledworks (Twinkly)"
brand: "Twinkly"
model: "TWS600STP"
model_aliases: ["TWS400STP", "TWS250STP", "TWS100STP", "TWI190STP", "TWW210STP", "Twinkly Strings", "Twinkly Curtain", "Twinkly Icicle", "Twinkly Flex", "Twinkly Dots"]
device_type: "twinkly_lights"
category: "lighting"
product_line: "Twinkly"
release_year: 2019
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["D8:F1:5B", "24:46:E4"]
  mdns_services: ["_http._tcp"]
  mdns_txt_keys: ["fwversion", "uuid"]
  default_ports: [80]
  signature_ports: [80]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["Twinkly_.*", "twinkly.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/xled/v1/gestalt"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "product_code"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "twinkly"
  polling_interval_sec: 5
  websocket_event: "twinkly:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M6"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "scenes"]

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 80
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "session_cookie"
  auth_detail: "Login via POST /xled/v1/login with challenge; server returns authentication_token; token passed as X-Auth-Token header; token expires after ~14400 seconds and must be refreshed via /xled/v1/verify"
  base_url_template: "http://{ip}/xled/v1"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "strip"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "both"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.twinkly.com/"
  api_docs: ""
  developer_portal: ""
  support: "https://www.twinkly.com/support/"
  community_forum: ""
  image_url: ""
  fcc_id: "2AYXD-TWS600STP"

# --- TAGS ---
tags: ["wifi", "bluetooth_le", "led_string", "decorative", "addressable", "local_rest_api", "led_mapping", "outdoor"]
---

# Twinkly Smart LED Lights

## What It Is

Twinkly is a line of smart LED string lights, curtains, icicle lights, and other decorative lighting products from Italian company Ledworks. What sets Twinkly apart is their LED mapping technology — using a smartphone camera, the Twinkly app photographs the physical layout of each LED, creating a precise coordinate map that enables spatially-aware effects and animations. Each LED is individually addressable with RGB (or RGBW depending on model) color. The lights connect via WiFi and expose a local REST API on port 80 that has been thoroughly documented by the community, making them a viable candidate for Haus integration. Products include Strings, Curtain, Icicle, Flex (bendable tube), Dots (individual LEDs on wire), and Line (rigid strips).

## How Haus Discovers It

1. **mDNS Discovery**: Twinkly devices advertise `_http._tcp.local.` via mDNS. The TXT records include `fwversion` (firmware version) and `uuid` (unique device ID). The instance name typically follows the pattern `Twinkly_{HEXID}`.
2. **OUI Match**: MAC addresses beginning with `D8:F1:5B` or `24:46:E4` are associated with Twinkly/Ledworks devices.
3. **HTTP Fingerprint**: An unauthenticated GET request to `http://{ip}/xled/v1/gestalt` returns device information including `product_code`, `uuid`, `hw_id`, `flash_size`, `led_type`, `led_version`, `product_name`, `device_name`, `number_of_led`, and firmware details. This endpoint is the definitive identification method.
4. **Hostname Pattern**: Twinkly devices typically use hostnames matching `Twinkly_XXXXXX` where `XXXXXX` is a hex identifier.

## Pairing / Authentication

Twinkly uses a challenge-response authentication flow for its local API:

1. **Login Request**: `POST /xled/v1/login` with body `{"challenge": "<32-byte-random-hex>"}`.
2. **Login Response**: The device returns an `authentication_token` (a base64-encoded token) and a `challenge-response` value.
3. **Verify**: `POST /xled/v1/verify` with the `challenge-response` from step 2 and the `X-Auth-Token` header set to the authentication token.
4. **Token Usage**: All subsequent requests include the `X-Auth-Token` header.
5. **Token Expiry**: The authentication token expires after approximately 4 hours (14400 seconds). Haus should re-authenticate proactively before expiry.

Note: The challenge-response mechanism is relatively weak — the shared secret is derived from the hardware ID. Community implementations have reverse-engineered the token generation. For Haus purposes, the standard login flow works reliably.

## API Reference

### Base URL

`http://{ip}/xled/v1`

### Authentication

All endpoints except `/xled/v1/gestalt` and `/xled/v1/login` require the `X-Auth-Token` header.

### Device Information

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/gestalt` | No | Full device info (model, LED count, firmware, hw_id, etc.) |
| GET | `/fw/version` | Yes | Firmware version details |
| GET | `/device_name` | Yes | Device name |
| PUT | `/device_name` | Yes | Set device name (body: `{"name": "My Twinkly"}`) |
| GET | `/timer` | Yes | Timer/schedule info |

### Power & Mode

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| GET | `/led/mode` | — | Get current mode |
| POST | `/led/mode` | `{"mode": "movie"}` | Set mode |

**Modes:**
- `off` — LEDs off
- `demo` — Built-in demo effect
- `movie` — Play uploaded movie/animation
- `rt` — Real-time control mode (for streaming)
- `color` — Solid color mode
- `effect` — Built-in effect from firmware

### Brightness

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| GET | `/led/out/brightness` | — | Get brightness |
| POST | `/led/out/brightness` | `{"value": 80, "type": "A"}` | Set brightness (0-100, type "A"=absolute) |

### Color (Solid Color Mode)

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| GET | `/led/color` | — | Get current solid color |
| POST | `/led/color` | `{"red": 255, "green": 0, "blue": 128}` | Set solid color (RGB 0-255) |
| POST | `/led/mode` | `{"mode": "color"}` | Must set mode to "color" first |

### Movies (Animations)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/movies` | List uploaded movies |
| GET | `/movies/current` | Get currently playing movie |
| POST | `/movies/current` | Set current movie by ID |
| POST | `/led/movie/config` | Configure movie parameters (frame delay, LED count, frames) |
| POST | `/led/movie/full` | Upload movie data (binary frame data) |

### Real-Time Control (Streaming)

For real-time per-LED color control:

1. Set mode to real-time: `POST /led/mode` with `{"mode": "rt"}`
2. Send UDP packets to device IP on port 7777.

**UDP Frame Format (Version 3):**

| Offset | Length | Field |
|--------|--------|-------|
| 0 | 1 | Version byte (0x03) |
| 1 | 8 | Authentication token (first 8 bytes of base64-decoded token) |
| 9 | 1 | Number of LEDs in this packet |
| 10+ | 3 per LED | R (1 byte) + G (1 byte) + B (1 byte) per LED |

For RGBW models, use 4 bytes per LED (R + G + B + W).

Maximum packet size accommodates approximately 150 RGB LEDs per packet. For larger strings, send multiple packets with appropriate LED offset handling.

### LED Layout / Mapping

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/led/layout/full` | Get full LED coordinate map (x, y, z positions) |
| POST | `/led/layout/full` | Upload LED coordinate map |

The layout data contains coordinates for each LED, as generated by the Twinkly app's camera mapping feature. Coordinates are normalized floats (0.0 to 1.0) in 2D or 3D space.

### Network Configuration

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/network/status` | WiFi connection status |
| GET | `/network/scan` | Scan for WiFi networks |

## AI Capabilities

When the AI concierge is chatting with Twinkly lights, it can:
- Turn lights on/off
- Set brightness
- Set a solid color (RGB)
- Switch between modes (off, color, movie, effect)
- List and activate uploaded movies/animations
- Report device status (current mode, brightness, LED count, firmware version)

## Quirks & Notes

- **Authentication Token Expiry**: Tokens expire after approximately 4 hours. Haus must implement automatic re-authentication. The `/xled/v1/verify` endpoint can extend the session.
- **Mode Switching**: The device operates in distinct modes (off, color, movie, rt, effect, demo). Setting a color requires first switching to "color" mode. Sending real-time data requires "rt" mode. Forgetting to switch modes is a common source of confusion.
- **LED Mapping Data**: The camera-mapped LED coordinates are stored on the device and retrievable via the API. This spatial data is unique to Twinkly and enables position-aware effects. Haus could use this for spatial UI representations.
- **Real-Time Streaming Timeout**: Real-time mode (`rt`) requires continuous UDP packets. If no data is sent for approximately 2.5 seconds, the device reverts to its previous mode.
- **Firmware Variations**: API endpoints and behaviors vary across firmware versions. The `/gestalt` endpoint returns the firmware version; Haus should maintain compatibility tables.
- **RGBW vs RGB Models**: Some models (e.g., Strings AWW+) support warm white in addition to RGB. The `led_profile` field in `/gestalt` indicates whether the device is "RGB" or "RGBW". The real-time streaming format changes accordingly (3 vs 4 bytes per LED).
- **Outdoor Models**: Some Twinkly products are rated IP44 for outdoor use. The API is identical regardless of indoor/outdoor rating.
- **Initial WiFi Setup**: First-time setup requires the Twinkly mobile app to provision WiFi credentials via BLE. Haus cannot provision unconfigured devices.
- **Community Documentation**: The API has no official documentation from Ledworks. All API knowledge comes from community reverse-engineering, primarily from the xled Python library (https://github.com/scrool/xled) and Home Assistant integration.

## Similar Devices

- **govee-rgbic-led-strip** — Addressable WiFi LED strip with local UDP API, simpler protocol
- **nanoleaf-shapes** — WiFi light panels with local REST API, similar local-first paradigm
- **lifx-a19-color** — WiFi bulb with local UDP protocol, no auth required
