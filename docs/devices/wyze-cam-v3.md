---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "wyze-cam-v3"
name: "Wyze Cam v3"
manufacturer: "Wyze Labs Inc."
brand: "Wyze"
model: "Wyze Cam v3"
model_aliases: ["WYZEC3", "Wyze Cam v3 Pro", "WYZE-CAKP2JFUS"]
device_type: "wyze_camera"
category: "security"
product_line: "Wyze"
release_year: 2020
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["2C:AA:8E", "7C:78:B2", "D0:3F:27", "A4:DA:22"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [8554]
  signature_ports: [8554]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^wyze.*", "^Wyze.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "wyze"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "password"
  ai_chattable: false
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: ["camera_stream", "camera_snapshot", "motion"]

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 8554
  transport: "TCP"
  encoding: "binary"
  auth_method: "basic_auth"
  auth_detail: "RTSP firmware uses RTSP authentication. Default credentials set during RTSP firmware setup. Stream available on port 8554. Stock firmware has no local API."
  base_url_template: "rtsp://{ip}:8554/live"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "usb"
  mounting: "shelf"
  indoor_outdoor: "both"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.wyze.com/products/wyze-cam"
  api_docs: ""
  developer_portal: ""
  support: "https://support.wyze.com"
  community_forum: "https://forums.wyze.com"
  image_url: ""
  fcc_id: "2AUIU-WYZEC3"

# --- TAGS ---
tags: ["rtsp-firmware", "budget", "1080p", "color-night-vision", "motion-detection", "microsd", "indoor-outdoor", "ip65"]
---

# Wyze Cam v3

## What It Is

> The Wyze Cam v3 is an ultra-affordable ($20-35) indoor/outdoor WiFi security camera with 1080p video, color night vision (Starlight sensor), motion detection, two-way audio, and IP65 weather resistance. Out of the box it is cloud-dependent through the Wyze app, but Wyze offers an official RTSP firmware that enables local RTSP streaming on port 8554. With RTSP firmware installed, the camera becomes one of the most cost-effective cameras for local home automation integration.

## How Haus Discovers It

1. **OUI match** -- Wyze cameras use MAC prefixes including `2C:AA:8E`, `7C:78:B2`, `D0:3F:27`, and `A4:DA:22` (Wyze Labs and partner manufacturer OUIs)
2. **Hostname pattern** -- Wyze devices may register hostnames starting with `wyze` or `Wyze`
3. **Port probe** -- TCP check on port 8554 (only if RTSP firmware is installed)
4. **RTSP probe** -- Attempt RTSP `OPTIONS` request to `rtsp://{ip}:8554/live` to confirm RTSP availability

## Pairing / Authentication

> **Stock firmware:** Setup through Wyze app only. Cloud account required. No local access.
>
> **RTSP firmware installation:**
> 1. Download RTSP firmware from Wyze support (search "Wyze Cam RTSP")
> 2. Flash firmware to camera via microSD card
> 3. Camera reboots with RTSP enabled
> 4. Configure RTSP credentials in Wyze app under camera settings > Advanced Settings > RTSP
> 5. Set username and password for RTSP authentication
>
> **RTSP authentication:** Basic auth with user-configured credentials.
>
> **Important:** The official RTSP firmware is based on an older firmware version and may lack newer features. Community alternatives like wz_mini_hacks provide RTSP on current firmware.

## API Reference

> ### RTSP Stream (RTSP firmware only)
>
> **Stream URL:**
> ```
> rtsp://{username}:{password}@{ip}:8554/live
> ```
>
> - **Port:** 8554
> - **Path:** `/live`
> - **Auth:** Basic authentication (credentials set in Wyze app)
> - **Video codec:** H.264
> - **Audio codec:** G.711 mu-law (optional)
> - **Resolution:** 1080p (1920x1080)
> - **Framerate:** Up to 20fps
>
> ### Alternative: wz_mini_hacks
>
> Community firmware mod that adds RTSP to current Wyze firmware:
> ```
> rtsp://{ip}:8554/unicast
> ```
> No authentication by default. Also supports ONVIF discovery.
>
> ### go2rtc Integration
>
> ```yaml
> streams:
>   wyze_cam: "rtsp://{user}:{pass}@{ip}:8554/live"
> ```

## AI Capabilities

> With RTSP firmware, AI integration follows the same pattern as other RTSP cameras. The AI could capture snapshots via go2rtc and use vision analysis to describe scenes. Not planned until M5 cameras milestone.

## Quirks & Notes

- **RTSP firmware is separate** -- The official RTSP firmware is a dedicated build, not a toggle in stock firmware. Installing it means running older firmware
- **wz_mini_hacks preferred** -- The community wz_mini_hacks project (github.com/gtxaspec/wz_mini_hacks) provides RTSP on current firmware with more features and is widely used in the home automation community
- **MicroSD local recording** -- Supports local recording to microSD card (up to 256GB) independent of cloud
- **Starlight sensor** -- Exceptional low-light performance with f/1.6 aperture; color night vision works in very dim conditions
- **IP65 weather resistance** -- Rated for outdoor use; USB power cable is the weak point for outdoor installations
- **Cam Plus subscription** -- Cloud features like person/pet/vehicle detection require Cam Plus ($1.99/mo per camera)
- **RTSP stream limits** -- Only one concurrent RTSP viewer is supported; multiple viewers may cause dropped frames
- **USB-A power** -- Powered by USB-A cable (included) with a 5V/1A adapter; no PoE option
- **Discontinued RTSP firmware** -- Wyze stopped updating the official RTSP firmware; community firmware is the better long-term path

## Similar Devices

> - [Tapo C200](tapo-c200.md) -- Similar price point with native ONVIF/RTSP (no firmware flash needed)
> - [Eufy Security Cam 2C](eufy-security-cam-2c.md) -- Slightly higher price, RTSP via HomeBase
> - [Amcrest IP Camera](amcrest-ip-camera.md) -- More expensive but professional-grade local streaming
> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- Battery-powered with native RTSP
