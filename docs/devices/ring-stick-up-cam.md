---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ring-stick-up-cam"
name: "Ring Stick Up Cam Battery/Wired"
manufacturer: "Ring LLC (Amazon)"
brand: "Ring"
model: "Stick Up Cam"
model_aliases: ["Ring Stick Up Cam Battery", "Ring Stick Up Cam Wired", "Ring Stick Up Cam Elite", "Ring Stick Up Cam 3rd Gen", "8SC1S9-WEN0"]
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
  auth_detail: "Unofficial Ring API. OAuth2 token from oauth.ring.com/oauth/token with mandatory 2FA."
  base_url_template: "https://api.ring.com/clients_api"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "battery"
  mounting: "wall"
  indoor_outdoor: "both"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://ring.com/products/stick-up-cam"
  api_docs: ""
  developer_portal: ""
  support: "https://support.ring.com"
  community_forum: "https://community.ring.com"
  image_url: ""
  fcc_id: "2AEUO-0120"

# --- TAGS ---
tags: ["cloud-only", "amazon", "no-local-api", "motion-detection", "two-way-audio", "1080p", "indoor-outdoor", "battery", "wired-option"]
---

# Ring Stick Up Cam Battery/Wired

## What It Is

> The Ring Stick Up Cam is a versatile indoor/outdoor WiFi security camera available in battery-powered and wired variants. It provides 1080p HD video, two-way audio, motion detection with customizable motion zones, and color night vision. The battery version runs on a removable rechargeable battery pack (or optional solar panel), while the wired version plugs into a standard outlet. Like all Ring cameras, it is entirely cloud-dependent with no local streaming capability.

## How Haus Discovers It

1. **OUI match** -- Ring devices use MAC prefixes `5C:47:5E`, `34:3E:A4`, `0C:47:C9`, `90:48:9A`, `CC:9E:A2`, `A0:23:9F`
2. **Hostname pattern** -- DHCP hostname typically starts with `ring-` or `Ring-`
3. **No port probe** -- No local ports are open; all communication is cloud-mediated
4. **Cloud enrichment** -- Device identification via unofficial Ring API after authentication

## Pairing / Authentication

> Same as all Ring devices: unofficial OAuth2 flow through `oauth.ring.com` with mandatory two-factor authentication. No official developer API exists. See [Ring Indoor Cam](ring-indoor-cam.md) for detailed auth flow.

## API Reference

> No official API. Uses the same unofficial reverse-engineered Ring API as all Ring cameras. See [Ring Indoor Cam](ring-indoor-cam.md) for endpoint details. The battery variant adds battery level reporting via `GET /doorbots/{id}/health` which returns `battery_life` as a percentage string.

## AI Capabilities

> AI integration is not planned due to lack of public API. If implemented, the AI could potentially report motion events, battery level (battery variant), and connection status.

## Quirks & Notes

- **Battery vs Wired** -- The battery version conserves power by only recording on motion events; the wired version can do continuous monitoring
- **Solar panel option** -- Ring Solar Panel (sold separately) can keep the battery topped up outdoors
- **Color night vision** -- Uses ambient light for color night vision; falls back to IR in complete darkness
- **Same cloud limitations as all Ring** -- No RTSP, no ONVIF, no local API
- **Weather resistance** -- IPX5 rated for outdoor use (splash-proof, not submersible)
- **Elite variant** -- The "Elite" version uses Power over Ethernet (PoE) instead of WiFi, but is still cloud-only

## Similar Devices

> - [Ring Indoor Cam](ring-indoor-cam.md) -- Indoor-only, more affordable
> - [Ring Floodlight Cam](ring-floodlight-cam.md) -- Outdoor with integrated floodlights
> - [Arlo Pro 5](arlo-pro-5.md) -- Competing wireless camera with optional RTSP
> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- Battery camera with native RTSP/ONVIF
