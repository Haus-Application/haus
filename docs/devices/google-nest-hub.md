---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "google-nest-hub"
name: "Google Nest Hub (2nd Gen)"
manufacturer: "Google LLC"
brand: "Google Nest"
model: "GXCA6"
model_aliases: ["Nest Hub", "Nest Hub 2nd Gen", "Google Home Hub", "H1A", "GXCA6"]
device_type: "cast_display"
category: "media"
product_line: "Nest"
release_year: 2021
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
  protocols_spoken: ["wifi", "bluetooth", "thread"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "F4:F5:D8"        # Google, Inc.
    - "30:FD:38"        # Google, Inc.
    - "54:60:09"        # Google, Inc.
    - "A4:77:33"        # Google, Inc.
    - "48:D6:D5"        # Google, Inc.
    - "20:DF:B9"        # Google, Inc.
    - "F8:0F:F9"        # Google, Inc.
    - "E4:F0:42"        # Google, Inc.
    - "1C:F2:9A"        # Google, Inc.
    - "44:07:0B"        # Google, Inc.
  mdns_services:
    - "_googlecast._tcp"
    - "_meshcop._udp"       # Thread border router (Nest Hub 2nd Gen)
  mdns_txt_keys:
    - "id"              # unique device ID
    - "cd"              # device certificate hash
    - "rm"              # room name
    - "ve"              # version
    - "md"              # model description (e.g., "Google Nest Hub")
    - "ic"              # icon path
    - "fn"              # friendly name
    - "ca"              # capabilities bitmask
    - "st"              # state
    - "bs"              # setup state
    - "nf"              # flags
    - "rs"              # receiver status
  default_ports: [8008, 8009, 8443]
  signature_ports: [8008, 8009]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Google-Nest-Hub"
    - "^[0-9a-f]{32}$"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8008
    path: "/setup/eureka_info?params=name,build_info,detail,device_info,opt_in"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"model_name\""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "read_only"
  integration_key: "cast"
  polling_interval_sec: 30
  websocket_event: "cast:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities:
  - "media_playback"
  - "volume"

# --- PROTOCOL ---
protocol:
  type: "protobuf_tls"
  port: 8009
  transport: "TLS"
  encoding: "Protobuf"
  auth_method: "none"
  auth_detail: "Port 8008 HTTP info API requires no authentication. Port 8009 Cast control protocol uses TLS with device certificate. Cast SDK required for control."
  base_url_template: "http://{ip}:8008"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "display"
  power_source: "mains"
  mounting: "tabletop"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le", "thread"]

# --- LINKS ---
links:
  product_page: "https://store.google.com/us/product/google_nest_hub_2nd_gen"
  api_docs: "https://developers.google.com/cast"
  developer_portal: "https://developers.google.com/cast"
  support: "https://support.google.com/googlenest/"
  community_forum: "https://support.google.com/googlenest/community"
  image_url: ""
  fcc_id: "A4RGXCA6"

# --- TAGS ---
tags: ["cast", "display", "google-assistant", "thread-border-router", "matter", "smart-display", "wifi"]
---

# Google Nest Hub (2nd Gen)

## What It Is

The Google Nest Hub (2nd Gen, model GXCA6) is a 7-inch smart display manufactured by Google LLC. It combines a touchscreen display with a full-range speaker, Google Assistant, and the Google Cast protocol. It serves as a media controller, photo frame, smart home dashboard, and -- critically for Haus -- a Thread border router. The 2nd Gen model added a Soli radar chip for sleep tracking and gesture control, plus Thread radio support for Matter-over-Thread device commissioning. It connects to the home network via Wi-Fi and exposes a read-only HTTP info API on port 8008 and the Cast control protocol on port 8009 (TLS/Protobuf). The original Google Home Hub and 1st Gen Nest Hub are functionally similar for Cast purposes but lack Thread support.

## How Haus Discovers It

1. **OUI Match** -- MAC address begins with a Google OUI prefix (`F4:F5:D8`, `30:FD:38`, `54:60:09`, `A4:77:33`, `48:D6:D5`, `20:DF:B9`, `F8:0F:F9`, `E4:F0:42`, `1C:F2:9A`, `44:07:0B`). Note that Google uses many OUI blocks across all hardware products, so OUI alone is not definitive.
2. **mDNS Discovery** -- Browse for `_googlecast._tcp.local.` services. The Nest Hub advertises rich TXT records including `fn` (friendly name, e.g., "Living Room Hub"), `md` (model description, e.g., "Google Nest Hub"), `id` (unique device ID), and `ca` (capabilities bitmask).
3. **Thread Border Router** -- The 2nd Gen Nest Hub also advertises `_meshcop._udp.local.` as a Thread border router, enabling Matter-over-Thread device commissioning.
4. **HTTP Fingerprint** -- `GET http://{ip}:8008/setup/eureka_info?params=name,build_info,detail,device_info,opt_in` returns JSON with device name, model, build info, and locale. The `device_info.model_name` field identifies it as "Google Nest Hub".
5. **Port Probe** -- Ports 8008 (HTTP info), 8009 (Cast TLS), and 8443 (HTTPS info) are open.

## Pairing / Authentication

The HTTP info API on port 8008 requires no authentication. Any device on the local network can query device information.

The Cast control protocol on port 8009 uses TLS with the device's self-signed certificate. The Google Cast SDK handles the TLS handshake and protocol negotiation. No user credentials or pairing flow is required for basic Cast functionality, though launching apps may require the sender to be registered with the Google Cast SDK Developer Console.

Initial device setup requires the Google Home app (iOS/Android) and a Google account.

## API Reference

### HTTP Info API (Port 8008)

**Base URL:** `http://{ip}:8008`

No authentication required. Read-only.

#### Device Info (eureka_info)

```
GET /setup/eureka_info?params=name,build_info,detail,device_info,opt_in
```

**Response:**
```json
{
  "name": "Living Room Hub",
  "build_info": {
    "cast_build_revision": "1.56.313396",
    "cast_control_version": 1
  },
  "device_info": {
    "model_name": "Google Nest Hub",
    "manufacturer": "Google Inc.",
    "product_name": "Google Nest Hub"
  },
  "detail": {
    "locale": {
      "display_string": "English (United States)"
    }
  },
  "opt_in": {
    "stats": true,
    "crash": true
  }
}
```

#### Configured Networks

```
GET /setup/configured_networks
```

Returns Wi-Fi network configuration (SSID, security type).

#### Supported App IDs

```
GET /setup/supported_app_ids
```

Returns list of supported Cast application IDs.

#### Device Offer

```
GET /setup/offer
```

Returns device capabilities and supported features.

### Cast Control Protocol (Port 8009)

The Cast protocol on port 8009 uses Protocol Buffers over TLS. This is a binary protocol that requires the Google Cast SDK or a compatible open-source implementation (e.g., go-chromecast, pychromecast).

**Protocol overview:**
1. TLS connection to port 8009 (accept self-signed certificate)
2. Exchange Cast Channel messages (Protobuf-encoded `CastMessage`)
3. Namespaces define message types:
   - `urn:x-cast:com.google.cast.tp.connection` -- connection management
   - `urn:x-cast:com.google.cast.tp.heartbeat` -- keep-alive PING/PONG
   - `urn:x-cast:com.google.cast.receiver` -- receiver status, app launch/stop
   - `urn:x-cast:com.google.cast.media` -- media control (play, pause, seek, volume)

**Key operations:**
- **Connect:** Send `CONNECT` message to `receiver-0` transport
- **Get Status:** Send `GET_STATUS` on receiver namespace
- **Set Volume:** Send `SET_VOLUME` with `{ "volume": { "level": 0.5 } }` on receiver namespace
- **Launch App:** Send `LAUNCH` with `{ "appId": "CC1AD845" }` (default media receiver)
- **Media Control:** Send `PLAY`, `PAUSE`, `SEEK`, `STOP` on media namespace to active session

### HTTPS Info API (Port 8443)

```
GET https://{ip}:8443/setup/eureka_info
```

Same as port 8008 but over HTTPS with self-signed certificate.

## AI Capabilities

When chatting with a Nest Hub via Haus, the AI can:
- **Report device identity** -- name, model, software version from eureka_info
- **Report network status** -- online/offline, responding on port 8008
- **Describe Thread role** -- whether the device is functioning as a Thread border router (2nd Gen only)

Cast control (play/pause/volume) requires Cast SDK integration and is planned for a future milestone.

For Nest Hub Max models that include a camera: if the user has connected their Google account via OAuth, the AI can access camera feeds and describe what the camera sees using vision AI.

## Quirks & Notes

- **Thread border router** -- The 2nd Gen Nest Hub is one of the most common Thread border routers in homes. It advertises `_meshcop._udp` for Thread commissioning. This is essential for Matter-over-Thread devices.
- **Google OUI overlap** -- Google uses dozens of OUI prefixes across Pixel phones, Chromecasts, Nest devices, and enterprise hardware. OUI alone cannot distinguish a Nest Hub from a Chromecast or Pixel phone. The mDNS `md` TXT record or eureka_info `model_name` is required for definitive identification.
- **No local control API** -- Unlike Hue or Kasa, there is no documented local REST API for controlling the Nest Hub. The Cast SDK on port 8009 is the only control path, and it requires Protobuf serialization.
- **Sleep tracking radar** -- The 2nd Gen Soli radar chip is not accessible via any network API. It reports to Google Health via cloud only.
- **Firmware auto-updates** -- Google pushes OTA firmware updates automatically. No user or developer control over update timing.
- **Cast protocol complexity** -- The Cast protocol is substantially more complex than simple REST APIs. It uses multiple "namespaces" (channels) multiplexed over a single TLS connection, Protobuf encoding, and a request-response model with async status updates. Open-source libraries like `go-chromecast` and `pychromecast` abstract this complexity.
- **Hub Max camera** -- The Nest Hub Max (model H2A) includes a built-in Nest camera. It advertises the same `_googlecast._tcp` service but has additional camera capabilities accessible through the Google Nest/SDM API (cloud, OAuth2 required).

## Similar Devices

- **[google-nest-mini](google-nest-mini.md)** -- Speaker-only Cast device with the same protocol but no display
- **[chromecast-google-tv](chromecast-google-tv.md)** -- HDMI Cast device with the same protocol
- **[apple-homepod-mini](apple-homepod-mini.md)** -- Competing smart speaker/hub ecosystem (HomeKit vs Google Home)
- **[amazon-echo-show-15](amazon-echo-show-15.md)** -- Competing smart display (Alexa vs Google Assistant)
