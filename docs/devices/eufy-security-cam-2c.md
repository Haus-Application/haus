---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "eufy-security-cam-2c"
name: "Eufy Security Cam 2C"
manufacturer: "Eufy (Anker Innovations)"
brand: "Eufy Security"
model: "eufyCam 2C"
model_aliases: ["eufyCam 2C Pro", "T8113", "T8114", "E220", "eufyCam 2C Wireless"]
device_type: "eufy_camera"
category: "security"
product_line: "Eufy Security"
release_year: 2020
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["98:F4:AB", "78:8C:B5", "AC:CF:85", "8C:85:80"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [554, 8080]
  signature_ports: [554]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^eufy.*", "^homebase.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8080
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "eufy"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "eufy"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "password"
  ai_chattable: false
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: ["camera_stream", "camera_snapshot", "motion", "battery_level"]

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 554
  transport: "TCP"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "RTSP streams enabled through HomeBase settings in Eufy Security app. Once enabled, RTSP is accessible without authentication on the HomeBase IP. Cloud API uses Eufy account credentials."
  base_url_template: "rtsp://{homebase_ip}:554/{camera_serial}"
  tls: false
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
  product_page: "https://www.eufy.com/eufycam-2c"
  api_docs: ""
  developer_portal: ""
  support: "https://www.eufy.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AOKB-T8113"

# --- TAGS ---
tags: ["rtsp-via-homebase", "battery", "local-storage", "no-subscription", "1080p", "indoor-outdoor", "homebase-hub", "person-detection"]
---

# Eufy Security Cam 2C

## What It Is

> The Eufy Security Cam 2C (eufyCam 2C) is a wireless battery-powered security camera that communicates through a HomeBase 2 hub. It provides 1080p video, 135-degree field of view, on-device AI person detection (no cloud processing required), two-way audio, and IP67 weather resistance. The standout feature for home automation is RTSP streaming through the HomeBase, enabling fully local video access without any cloud dependency or subscription. Eufy's "no monthly fee" positioning means local storage on the HomeBase's built-in 16GB eMMC (or external HDD) handles recording.

## How Haus Discovers It

1. **OUI match** -- Eufy HomeBase uses MAC prefixes including `98:F4:AB`, `78:8C:B5`, `AC:CF:85`, and `8C:85:80` (Anker Innovations OUIs). Note: the HomeBase appears on the network, cameras connect to the HomeBase over a proprietary wireless link
2. **Port probe** -- TCP check on port 554 (RTSP) on the HomeBase IP when RTSP is enabled
3. **RTSP probe** -- Attempt RTSP `OPTIONS` to `rtsp://{homebase_ip}:554/` to confirm streaming availability
4. **HTTP probe** -- Check port 8080 on HomeBase for web interface

## Pairing / Authentication

> Setup requires the Eufy Security app:
>
> 1. Create Eufy account and add HomeBase 2 via app
> 2. Pair camera to HomeBase by pressing sync button on camera
> 3. Camera connects to HomeBase over proprietary 2.4GHz wireless link
> 4. Enable RTSP in Eufy Security app: Device Settings > Storage > NAS (RTSP)
> 5. Each camera gets an RTSP URL displayed in the app
>
> **RTSP authentication:** Once enabled, RTSP streams from the HomeBase typically do not require additional authentication. The stream URL includes the camera serial number as the path.

## API Reference

> ### RTSP Stream (via HomeBase)
>
> **Stream URL:**
> ```
> rtsp://{homebase_ip}:554/{camera_serial_number}
> ```
>
> - **Port:** 554 (standard RTSP port on HomeBase)
> - **Path:** Camera serial number (displayed in Eufy app)
> - **Auth:** None required once RTSP is enabled
> - **Video codec:** H.264
> - **Resolution:** 1080p (1920x1080)
> - **Framerate:** Up to 15fps (battery optimization)
>
> ### go2rtc Integration
>
> ```yaml
> streams:
>   eufy_front_door: "rtsp://{homebase_ip}:554/{serial}"
> ```
>
> ### Unofficial Local API (eufy-security-ws)
>
> The community project `eufy-security-ws` provides a WebSocket API wrapper for local HomeBase control:
> - Push notification events for motion/person detection
> - Arm/disarm modes
> - Camera property queries
> - Lives at github.com/bropat/eufy-security-ws

## AI Capabilities

> With RTSP streaming enabled, AI integration follows the standard camera pattern: snapshots via go2rtc, Claude vision analysis for scene description. On-device person detection events can be surfaced through the eufy-security-ws integration. Not planned until M5 cameras milestone.

## Quirks & Notes

- **HomeBase required** -- The camera itself does not appear on the WiFi network; the HomeBase is the network-visible device that serves RTSP streams
- **RTSP must be explicitly enabled** -- Not on by default; must be toggled per-camera in the Eufy app under NAS/RTSP settings
- **Battery impact** -- RTSP streaming significantly increases battery drain; the camera cannot stream continuously on battery. RTSP sessions wake the camera from sleep
- **On-device AI** -- Person detection runs on the camera's local processor, not in the cloud. This means no subscription needed for smart detection
- **180-day battery life** -- Advertised on normal use; RTSP access reduces this dramatically
- **HomeBase storage** -- 16GB built-in eMMC storage; supports external USB HDD for expanded local recording
- **Privacy controversy** -- Eufy faced criticism in 2022 for uploading thumbnails to cloud servers despite "local only" marketing. They've since addressed this but trust remains an issue
- **No ONVIF** -- Despite having RTSP, there is no ONVIF support for standardized discovery and configuration
- **HomeBase 3 (S380)** -- Newer HomeBase 3 supports more cameras and has expanded local storage, but RTSP availability varies by firmware

## Similar Devices

> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- Battery camera with native RTSP + ONVIF (no hub required)
> - [Arlo Pro 5](arlo-pro-5.md) -- Battery camera with optional RTSP via SmartHub
> - [Wyze Cam v3](wyze-cam-v3.md) -- Budget wired camera with RTSP firmware
> - [Blink Outdoor](blink-outdoor.md) -- Amazon's battery camera (no local access at all)
