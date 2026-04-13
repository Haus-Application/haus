---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "amcrest-ip-camera"
name: "Amcrest IP Camera (IP8M Series)"
manufacturer: "Amcrest Technologies"
brand: "Amcrest"
model: "IP8M-2496E"
model_aliases: ["IP8M-T2599E", "IP8M-2493E", "IP5M-T1179E", "IP4M-1051", "Amcrest 4K", "Amcrest UltraHD"]
device_type: "amcrest_camera"
category: "security"
product_line: "Amcrest"
release_year: 2021
discontinued: false
price_range: "$$"

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
  mac_prefixes: ["9C:8E:CD", "3C:EF:8C", "E0:50:8B"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [554, 80, 443, 37777]
  signature_ports: [37777, 554]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^amcrest.*", "^Amcrest.*", "^IP8M.*", "^IP5M.*", "^IP4M.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: "Web Service"
    server_header: ""
    body_contains: "Amcrest"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "amcrest"
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
  port: 554
  transport: "TCP"
  encoding: "binary"
  auth_method: "basic_auth"
  auth_detail: "RTSP, ONVIF, and web UI all use the camera's admin credentials. Default username is 'admin'. Password set during initial setup via web UI or Amcrest app. Web UI on port 80 uses Digest authentication. Port 37777 is Amcrest's proprietary binary protocol."
  base_url_template: "rtsp://{ip}:554/cam/realmonitor?channel=1&subtype=0"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "poe"
  mounting: "ceiling"
  indoor_outdoor: "both"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://amcrest.com/ip-cameras.html"
  api_docs: "https://support.amcrest.com/hc/en-us/articles/360048498012-Amcrest-HTTP-API-SDK"
  developer_portal: ""
  support: "https://support.amcrest.com"
  community_forum: "https://amcrest.com/forum/"
  image_url: ""
  fcc_id: "2AN3B-IP8M2496E"

# --- TAGS ---
tags: ["rtsp", "onvif", "local-first", "poe", "wifi", "4k-resolution", "web-ui", "http-api", "professional", "dahua-compatible", "no-subscription", "microsd", "nas"]
---

# Amcrest IP Camera (IP8M Series)

## What It Is

> Amcrest IP cameras (IP8M series) are professional-grade security cameras available in 4K (8MP), 5MP, and 4MP variants with turret, bullet, and dome form factors. They provide native RTSP, ONVIF, and a full HTTP API for local control -- no cloud required. Most models support both PoE (Power over Ethernet) and WiFi. The cameras include a comprehensive web UI for configuration, support microSD recording, NAS/FTP upload, and multiple simultaneous stream profiles. Amcrest cameras are built on the Dahua platform, meaning they also work with Dahua NVRs and software. They represent one of the best options for local-first camera integration in Haus.

## How Haus Discovers It

1. **OUI match** -- Amcrest cameras use MAC prefixes including `9C:8E:CD`, `3C:EF:8C`, and `E0:50:8B` (shared with Dahua, as Amcrest uses Dahua hardware)
2. **Port probe** -- TCP check on port 37777 (Amcrest/Dahua proprietary protocol) is a strong identifier. Also ports 554 (RTSP) and 80 (web UI)
3. **ONVIF discovery** -- ONVIF WS-Discovery probe (UDP multicast 239.255.255.250:3702). Amcrest cameras respond with ONVIF device service URL
4. **HTTP fingerprint** -- Web UI on port 80 contains "Amcrest" or "Web Service" in the title
5. **RTSP probe** -- RTSP `OPTIONS` to `rtsp://{ip}:554/` confirms availability

## Pairing / Authentication

> Setup via web UI or Amcrest app:
>
> 1. Connect camera via Ethernet (PoE) or WiFi (use Amcrest app for WiFi setup)
> 2. Find camera IP via Amcrest IP Config Tool or network scan
> 3. Access web UI at `http://{ip}` (default user: `admin`, password set on first login)
> 4. Configure video settings, motion detection, network, and recording
>
> **Credentials:** Single admin account used for web UI, RTSP, ONVIF, and HTTP API. Some models support additional user accounts with restricted permissions.

## API Reference

> ### RTSP Streams
>
> Amcrest provides main and sub streams per channel:
>
> **Main stream (full resolution):**
> ```
> rtsp://admin:{password}@{ip}:554/cam/realmonitor?channel=1&subtype=0
> ```
>
> **Sub stream (low bandwidth):**
> ```
> rtsp://admin:{password}@{ip}:554/cam/realmonitor?channel=1&subtype=1
> ```
>
> - **Port:** 554
> - **Auth:** Digest authentication (admin credentials)
> - **Video codec:** H.264 or H.265 (configurable)
> - **Main resolution:** Up to 3840x2160 (4K) on IP8M models
> - **Sub resolution:** Configurable, typically 704x480 or 640x360
> - **Framerate:** Up to 20fps (4K) or 30fps (lower resolutions)
>
> ### ONVIF
>
> - **ONVIF port:** 80 (same as web UI)
> - **Device service URL:** `http://{ip}/onvif/device_service`
> - **Media service URL:** `http://{ip}/onvif/media_service`
> - **Events service URL:** `http://{ip}/onvif/event_service`
> - **Profiles:** Profile S (streaming), Profile T (advanced video), Profile G (recording)
> - **PTZ:** Supported on PTZ models via ONVIF PTZ service
> - **Events:** Motion, video loss, tampering via pull-point subscription
>
> ### HTTP API (CGI)
>
> Amcrest publishes a comprehensive HTTP API (Dahua-compatible):
>
> **Snapshot:**
> ```
> GET http://admin:{password}@{ip}/cgi-bin/snapshot.cgi?channel=1
> ```
> Returns JPEG image.
>
> **Get device info:**
> ```
> GET http://admin:{password}@{ip}/cgi-bin/magicBox.cgi?action=getDeviceType
> ```
>
> **Get system info:**
> ```
> GET http://admin:{password}@{ip}/cgi-bin/magicBox.cgi?action=getSystemInfo
> ```
>
> **Motion detection events (long-poll):**
> ```
> GET http://admin:{password}@{ip}/cgi-bin/eventManager.cgi?action=attach&codes=[VideoMotion]
> ```
> Returns chunked HTTP response with events as they occur (Server-Sent Events style).
>
> **PTZ control (PTZ models):**
> ```
> GET http://admin:{password}@{ip}/cgi-bin/ptz.cgi?action=start&channel=1&code=Right&arg1=0&arg2=1&arg3=0
> ```
>
> ### go2rtc Integration
>
> ```yaml
> streams:
>   amcrest_front: "rtsp://admin:{pass}@{ip}:554/cam/realmonitor?channel=1&subtype=0"
> ```

## AI Capabilities

> With native RTSP and HTTP snapshot API, AI integration is straightforward: go2rtc captures snapshots, Claude vision analyzes scenes. The HTTP event long-poll API enables real-time motion detection notifications. Planned for M5 cameras milestone.

## Quirks & Notes

- **Dahua platform** -- Amcrest cameras are rebranded/modified Dahua hardware. The firmware, HTTP API, and ONVIF implementation are Dahua-based. Most Dahua documentation and tools also work with Amcrest
- **Port 37777** -- This is the Amcrest/Dahua proprietary binary protocol port used by their desktop software (SmartPSS, Amcrest Surveillance Pro). Detecting this port is a reliable way to identify Amcrest/Dahua cameras
- **H.265 support** -- Most IP8M models support H.265 encoding for better compression, but H.264 is recommended for broader compatibility with streaming tools
- **Digest auth on RTSP** -- Amcrest uses Digest authentication for RTSP (not Basic). Some older RTSP clients may not support this; go2rtc handles it correctly
- **Web UI requires browser plugins** -- The full web UI (with live view) requires a browser plugin on some browsers; the settings pages work without plugins
- **Dual network** -- Most PoE models also have WiFi; the camera can connect via both simultaneously, though this is not recommended (use one or the other)
- **MicroSD + NAS** -- Supports local microSD recording plus NAS/SMB and FTP upload simultaneously
- **API documentation** -- Amcrest publishes the HTTP API SDK document (PDF), which is one of the most comprehensive camera API docs available for consumer-grade cameras
- **Default port conflicts** -- The web UI on port 80 and ONVIF on port 80 can cause confusion; they share the same port

## Similar Devices

> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- More affordable, battery-powered, also native RTSP/ONVIF
> - [UniFi Protect G4 Bullet](unifi-protect-g4-bullet.md) -- Professional PoE with RTSP + Protect API
> - [Tapo C200](tapo-c200.md) -- Budget indoor with ONVIF/RTSP
> - [Wyze Cam v3](wyze-cam-v3.md) -- Budget option with RTSP firmware
