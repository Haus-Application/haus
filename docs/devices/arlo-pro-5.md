---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "arlo-pro-5"
name: "Arlo Pro 5"
manufacturer: "Arlo Technologies Inc."
brand: "Arlo"
model: "Pro 5 2K"
model_aliases: ["Arlo Pro 5S 2K", "VMC4060P", "VMC4260P", "Arlo Pro 5S"]
device_type: "arlo_camera"
category: "security"
product_line: "Arlo"
release_year: 2023
discontinued: false
price_range: "$$"

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
  mac_prefixes: ["9C:53:22", "9C:B7:0D", "20:3D:BD", "70:EE:50", "74:DA:38"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^arlo.*", "^Arlo.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "arlo"
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["camera_stream", "camera_snapshot", "motion", "battery_level"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Arlo cloud API. Authentication via https://my.arlo.com. Two-factor auth required. RTSP available on certain firmware versions via direct WiFi connection."
  base_url_template: "https://myapi.arlo.com/hmsweb"
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
  product_page: "https://www.arlo.com/cameras/pro/arlo-pro-5.html"
  api_docs: "https://developer.arlo.com"
  developer_portal: "https://developer.arlo.com"
  support: "https://www.arlo.com/support"
  community_forum: "https://community.arlo.com"
  image_url: ""
  fcc_id: "2AEOS-VMC4060P"

# --- TAGS ---
tags: ["cloud-primary", "battery", "2k-resolution", "color-night-vision", "spotlight", "two-way-audio", "rtsp-possible", "wifi-direct"]
---

# Arlo Pro 5

## What It Is

> The Arlo Pro 5 (also marketed as Pro 5S 2K) is a premium wireless security camera with 2K HDR video, integrated spotlight, color night vision, two-way audio, and a rechargeable battery. It connects directly to WiFi (no hub required, though the Arlo SmartHub is optional for extended range). Arlo's primary interface is cloud-based through the Arlo app, but the camera supports RTSP streaming on certain firmware versions when connected to an Arlo SmartHub, providing a potential path for local integration.

## How Haus Discovers It

1. **OUI match** -- Arlo cameras use MAC prefixes including `9C:53:22`, `9C:B7:0D`, `20:3D:BD`, `70:EE:50`, and `74:DA:38` (Netgear/Arlo OUIs)
2. **Hostname pattern** -- Arlo devices may register hostnames starting with `arlo` or `Arlo`
3. **No local port probe** -- Direct WiFi Arlo cameras have no open local ports for streaming
4. **SmartHub RTSP** -- If connected via SmartHub, RTSP may be available on the hub's IP address

## Pairing / Authentication

> Arlo cameras require cloud account setup:
>
> 1. Create Arlo account at my.arlo.com or through the Arlo app
> 2. Add camera via QR code scan in the Arlo app
> 3. Camera connects to WiFi directly or via Arlo SmartHub
> 4. For RTSP access: requires SmartHub, enable RTSP in SmartHub settings
> 5. API authentication uses Arlo cloud OAuth with mandatory 2FA
>
> **RTSP on SmartHub:** Navigate to SmartHub settings in the Arlo app, enable RTSP for each camera. This generates per-camera RTSP URLs.

## API Reference

> Arlo has a developer portal at developer.arlo.com that provides limited API access. The primary unofficial API:
>
> ### Authentication
> ```
> POST https://ocapi-app.arlo.com/api/auth
> Content-Type: application/json
>
> {"email": "...", "password": "...", "EnvSource": "prod"}
> ```
> Returns a token and requires 2FA verification via `POST /api/auth/verify`.
>
> ### List Devices
> ```
> GET https://myapi.arlo.com/hmsweb/users/devices
> Authorization: {token}
> ```
>
> ### RTSP Stream (SmartHub only)
> When RTSP is enabled on the SmartHub, streams are available at:
> ```
> rtsp://{smarthub_ip}:554/{camera_serial_number}
> ```
> The RTSP URL format and port may vary by SmartHub firmware version.
>
> ### Start/Stop Stream
> ```
> POST https://myapi.arlo.com/hmsweb/users/devices/startStream
> {"from": "{user_id}", "to": "{base_station_id}", "action": "set", "resource": "cameras/{camera_id}", "publishResponse": true, "properties": {"activityState": "startUserStream"}}
> ```

## AI Capabilities

> AI integration is not planned for V1. Future integration could leverage RTSP streams (via SmartHub) for local snapshot and vision analysis, similar to Nest camera AI capabilities.

## Quirks & Notes

- **RTSP requires SmartHub** -- RTSP streaming is only available when the camera is paired to an Arlo SmartHub (VMB4540 or VMB5000); direct WiFi cameras have no local streaming
- **Firmware-dependent RTSP** -- RTSP availability depends on SmartHub firmware version; Arlo has been known to remove RTSP in firmware updates
- **Battery life** -- 3-6 months typical battery life depending on activity level; motion-activated recording to conserve power
- **Arlo Secure subscription** -- Cloud recording, AI detection features, and some advanced features require Arlo Secure plan ($7.99/mo per camera or $17.99/mo for all cameras)
- **2FA mandatory** -- Arlo enforces two-factor authentication, complicating automated API access
- **WiFi 6 support** -- Supports 802.11ax dual-band for better connectivity
- **Magnetic mount** -- Uses a strong magnetic mount for easy positioning and adjustment
- **Developer program** -- Arlo has an official developer portal at developer.arlo.com, but access is limited and primarily aimed at partners

## Similar Devices

> - [Ring Stick Up Cam](ring-stick-up-cam.md) -- Competing battery camera, cloud-only
> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- Battery camera with native RTSP/ONVIF (better local support)
> - [Eufy Security Cam 2C](eufy-security-cam-2c.md) -- Battery camera with local RTSP via HomeBase
> - [Wyze Cam v3](wyze-cam-v3.md) -- Budget alternative with RTSP firmware option
