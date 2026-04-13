---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ring-floodlight-cam"
name: "Ring Floodlight Cam Wired Pro"
manufacturer: "Ring LLC (Amazon)"
brand: "Ring"
model: "Floodlight Cam Wired Pro"
model_aliases: ["Ring Floodlight Cam Plus", "Ring Floodlight Cam Wired Plus", "B08F6GPQQ7", "Ring Floodlight Camera"]
device_type: "ring_camera"
category: "security"
product_line: "Ring"
release_year: 2021
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
capabilities: ["camera_stream", "camera_snapshot", "motion", "on_off"]

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
  power_source: "hardwired"
  mounting: "wall"
  indoor_outdoor: "outdoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://ring.com/products/floodlight-cam-wired-pro"
  api_docs: ""
  developer_portal: ""
  support: "https://support.ring.com"
  community_forum: "https://community.ring.com"
  image_url: ""
  fcc_id: "2AEUO-0114"

# --- TAGS ---
tags: ["cloud-only", "amazon", "no-local-api", "motion-detection", "two-way-audio", "1080p", "outdoor", "floodlight", "siren", "3d-motion", "bird-eye-view"]
---

# Ring Floodlight Cam Wired Pro

## What It Is

> The Ring Floodlight Cam Wired Pro is a premium outdoor security camera with integrated dual LED floodlights (2000 lumens) and a 110dB siren. It provides 1080p HDR video, two-way audio with noise cancellation, 3D motion detection with Bird's Eye View (radar-based overhead motion map), and color night vision. Hardwired to an exterior junction box (replaces an existing outdoor light), it combines security camera, motion-activated floodlights, and audible alarm in one device. Like all Ring products, it is entirely cloud-dependent.

## How Haus Discovers It

1. **OUI match** -- Ring MAC prefixes: `5C:47:5E`, `34:3E:A4`, `0C:47:C9`, `90:48:9A`, `CC:9E:A2`, `A0:23:9F`
2. **Hostname pattern** -- DHCP hostname typically starts with `ring-` or `Ring-`
3. **No port probe** -- No local ports; all communication via Ring cloud
4. **Cloud enrichment** -- Device type, name, and capabilities retrieved via unofficial Ring API

## Pairing / Authentication

> Same unofficial OAuth2 flow as all Ring devices. See [Ring Indoor Cam](ring-indoor-cam.md) for detailed auth flow.

## API Reference

> No official API. Uses the same unofficial reverse-engineered Ring API. Additional endpoints for the Floodlight Cam:
>
> - `PUT /doorbots/{id}/floodlight_light_on` -- Turn floodlights on
> - `PUT /doorbots/{id}/floodlight_light_off` -- Turn floodlights off
> - `PUT /doorbots/{id}/siren_on` -- Activate 110dB siren
> - `PUT /doorbots/{id}/siren_off` -- Deactivate siren
>
> The `on_off` capability maps to floodlight and siren control, not the camera itself (camera is always active when powered).

## AI Capabilities

> AI integration is not planned due to lack of public API. If implemented, the AI could potentially control the floodlights (on/off), trigger the siren, and report motion events with Bird's Eye View location data.

## Quirks & Notes

- **Hardwired installation** -- Requires an existing outdoor electrical junction box; replaces a standard outdoor light fixture
- **3D Motion Detection** -- Uses radar in addition to the camera sensor to track motion position and direction; generates a Bird's Eye View motion map
- **Dual floodlights** -- Two adjustable LED panels, 2000 lumens total, motion-activated or manually controlled
- **110dB siren** -- Remotely triggered or automated via motion rules
- **Pre-roll** -- 4 seconds of pre-roll video capture (wired power enables always-on recording to buffer)
- **Same cloud limitations** -- No RTSP, no ONVIF, no local control; floodlight and siren control also requires cloud
- **WiFi 6 support** -- Supports 802.11ax (WiFi 6) for better throughput and range
- **Power requirements** -- Requires 16-24V AC or standard 120V/240V AC junction box depending on model

## Similar Devices

> - [Ring Indoor Cam](ring-indoor-cam.md) -- Indoor-only, budget option
> - [Ring Stick Up Cam](ring-stick-up-cam.md) -- Indoor/outdoor, battery or wired
> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- Outdoor camera with local RTSP (no floodlight)
> - [UniFi Protect G4 Bullet](unifi-protect-g4-bullet.md) -- Professional outdoor camera with local API
