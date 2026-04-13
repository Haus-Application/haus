---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ring-indoor-cam"
name: "Ring Indoor Cam (2nd Gen)"
manufacturer: "Ring LLC (Amazon)"
brand: "Ring"
model: "Indoor Cam 2nd Gen"
model_aliases: ["Ring Indoor Cam Gen 2", "Ring Indoor Camera 2nd Generation", "B0B6GLQMHJ"]
device_type: "ring_camera"
category: "security"
product_line: "Ring"
release_year: 2023
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
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["5C:47:5E", "34:3E:A4", "0C:47:C9", "90:48:9A", "CC:9E:A2", "A0:23:9F"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^ring-.*", "^Ring-.*"]
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
capabilities: ["camera_stream", "camera_snapshot", "motion"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Ring has no official public API. Unofficial reverse-engineered API uses oauth2 token from Ring account. Two-factor auth required. Token endpoint at oauth.ring.com/oauth/token."
  base_url_template: "https://api.ring.com/clients_api"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "usb"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://ring.com/products/ring-indoor-cam-2nd-gen"
  api_docs: ""
  developer_portal: ""
  support: "https://support.ring.com"
  community_forum: "https://community.ring.com"
  image_url: ""
  fcc_id: "2AEUO-0128"

# --- TAGS ---
tags: ["cloud-only", "amazon", "no-local-api", "motion-detection", "two-way-audio", "1080p", "indoor"]
---

# Ring Indoor Cam (2nd Gen)

## What It Is

> The Ring Indoor Cam (2nd Gen) is an affordable WiFi security camera designed for indoor monitoring. It provides 1080p HD video, two-way audio, motion detection, and night vision. Part of Amazon's Ring ecosystem, it integrates with Alexa and requires the Ring app for setup and live viewing. It is entirely cloud-dependent -- there is no local API, no RTSP stream, and no ONVIF support. All video processing and storage happens on Ring's servers, with optional Ring Protect subscription for video recording.

## How Haus Discovers It

1. **OUI match** -- Ring devices use MAC prefixes including `5C:47:5E`, `34:3E:A4`, `0C:47:C9`, `90:48:9A`, `CC:9E:A2`, and `A0:23:9F` (assigned to Ring LLC and Amazon Technologies)
2. **Hostname pattern** -- Ring devices often register DHCP hostnames starting with `ring-` or `Ring-`
3. **No port probe** -- Ring cameras have no open local ports; they communicate exclusively with Ring cloud servers
4. **Cloud enrichment** -- Device name and type would be retrieved from the Ring API (unofficial) after OAuth authentication

## Pairing / Authentication

> Ring has no official public API or developer program. Integration would require using the unofficial reverse-engineered Ring API:
>
> 1. User provides Ring account email and password
> 2. Two-factor authentication code is required (SMS or email)
> 3. Token exchange at `https://oauth.ring.com/oauth/token`
> 4. Access token used for subsequent API calls
> 5. Tokens must be refreshed periodically
>
> **Risk:** Amazon actively patches against unofficial API access. This integration could break at any time.

## API Reference

> Ring does not publish an official API. An unofficial, community-reverse-engineered API exists but is undocumented and unsupported by Amazon. Key known endpoints:
>
> - `GET /ring_devices` -- list all Ring devices on account
> - `GET /dings/active` -- active motion/doorbell events
> - `GET /doorbots/{id}/history` -- event history
> - `POST /doorbots/{id}/live_view` -- initiate live stream (returns SIP/WebRTC session)
>
> Live streaming uses a proprietary SIP-over-WebSocket protocol, making local stream extraction extremely difficult.

## AI Capabilities

> AI integration is not planned for Ring devices due to the lack of a public API and cloud-only architecture. If implemented in the future, the AI could potentially report motion events and connection status.

## Quirks & Notes

- **No local API whatsoever** -- Ring cameras are fully cloud-locked. No RTSP, no ONVIF, no local HTTP server
- **Unofficial API instability** -- Amazon regularly changes the Ring API, breaking third-party integrations
- **Two-factor auth required** -- Ring enforces 2FA on all accounts, complicating automated authentication
- **SIP-based streaming** -- Live video uses a proprietary SIP protocol over WebSocket, not standard WebRTC or RTSP
- **Ring Protect subscription** -- Video recording requires a paid subscription ($3.99/mo per camera or $12.99/mo for all cameras)
- **Privacy concerns** -- All video passes through Amazon's cloud; no local storage option
- **End-to-end encryption** -- Available as opt-in feature since 2022, but breaks cloud-based video features

## Similar Devices

> - [Ring Stick Up Cam](ring-stick-up-cam.md) -- Indoor/outdoor version with battery option
> - [Ring Floodlight Cam](ring-floodlight-cam.md) -- Outdoor with integrated floodlights
> - [Blink Outdoor](blink-outdoor.md) -- Another Amazon cloud-only camera, even more locked down
> - [Wyze Cam v3](wyze-cam-v3.md) -- Affordable alternative with optional RTSP firmware
