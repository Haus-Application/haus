---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "unifi-protect-g4-bullet"
name: "UniFi Protect G4 Bullet"
manufacturer: "Ubiquiti Inc."
brand: "UniFi"
model: "UVC-G4-Bullet"
model_aliases: ["UVC-G4-BULLET", "UniFi Video G4 Bullet", "G4 Bullet", "UVC-G4-Bullet-3"]
device_type: "unifi_camera"
category: "security"
product_line: "UniFi Protect"
release_year: 2020
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: false
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["ethernet"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["FC:EC:DA", "24:5A:4C", "78:8A:20", "80:2A:A8", "68:D7:9A", "74:AC:B9", "E0:63:DA", "18:E8:29", "B4:FB:E4"]
  mdns_services: ["_ubnt._tcp"]
  mdns_txt_keys: ["model", "firmware"]
  default_ports: [7443, 7447, 7080, 554, 80]
  signature_ports: [7443, 7447]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^UVC-G4.*", "^UBNT-.*", "^UniFi.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 7080
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "UniFi Video"
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "unifi"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "password"
  ai_chattable: false
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: ["camera_stream", "camera_snapshot", "motion"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 7443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "session_cookie"
  auth_detail: "UniFi Protect API on the UDM/NVR (not the camera). Login via POST to /api/auth/login with local UniFi account credentials. Returns a session cookie and CSRF token. RTSP streams available directly from camera on port 554 once enabled in Protect settings."
  base_url_template: "https://{udm_ip}:7443/proxy/protect/api"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "poe"
  mounting: "wall"
  indoor_outdoor: "outdoor"
  wireless_radios: []

# --- LINKS ---
links:
  product_page: "https://store.ui.com/us/en/collections/unifi-cameras-702/products/uvc-g4-bullet"
  api_docs: ""
  developer_portal: ""
  support: "https://help.ui.com"
  community_forum: "https://community.ui.com"
  image_url: ""
  fcc_id: "SWX-G4BULLET"

# --- TAGS ---
tags: ["rtsp", "poe", "local-first", "unifi-protect", "4mp", "professional", "no-subscription", "nvr-required", "self-signed-tls", "weatherproof"]
---

# UniFi Protect G4 Bullet

## What It Is

> The UniFi Protect G4 Bullet is a professional-grade PoE (Power over Ethernet) security camera from Ubiquiti's UniFi Protect line. It provides 4MP (2688x1512) video at 24fps, 802.3af PoE power, IP67 weather resistance, IR night vision (up to 25m), and a built-in microphone. It requires a UniFi Protect NVR (such as the UniFi Dream Machine Pro, UDM-SE, or Cloud Key Gen2 Plus) for management and recording -- the camera cannot operate standalone. Once adopted by the Protect controller, RTSP streams are available directly from the camera. The entire system is fully local with no cloud dependency and no subscription fees.

## How Haus Discovers It

1. **OUI match** -- Ubiquiti devices use numerous MAC prefixes including `FC:EC:DA`, `24:5A:4C`, `78:8A:20`, `80:2A:A8`, `68:D7:9A`, `74:AC:B9`, `E0:63:DA`, `18:E8:29`, and `B4:FB:E4`
2. **mDNS** -- Ubiquiti devices advertise `_ubnt._tcp` with TXT records containing model and firmware info
3. **Port probe** -- Protect controller on port 7443 (Protect API) and 7447 (Protect WebSocket). Camera itself on port 554 (RTSP) and 80 (adoption)
4. **HTTP fingerprint** -- Camera management interface on port 7080 returns "UniFi Video" server header
5. **RTSP probe** -- Once RTSP is enabled in Protect settings, stream available on camera's port 554

## Pairing / Authentication

> UniFi Protect cameras are "adopted" by a Protect controller:
>
> 1. Connect camera to PoE switch or PoE port on UDM
> 2. Camera powers up and appears in UniFi Protect controller's "Devices" tab as "Pending Adoption"
> 3. Click "Adopt" in the Protect UI
> 4. Camera downloads firmware and joins the Protect system
> 5. Enable RTSP: In Protect UI, go to Camera Settings > Advanced > enable RTSP streams
>
> **API authentication:** UniFi Protect API uses session-based authentication:
> ```
> POST https://{udm_ip}/api/auth/login
> Content-Type: application/json
>
> {"username": "admin", "password": "{password}"}
> ```
> Returns `Set-Cookie` header with session token and `X-CSRF-Token` header for subsequent requests.

## API Reference

> ### RTSP Streams (directly from camera)
>
> RTSP must be enabled per-camera in UniFi Protect settings. Once enabled:
>
> **High quality:**
> ```
> rtsp://{camera_ip}:554/s0
> ```
>
> **Medium quality:**
> ```
> rtsp://{camera_ip}:554/s1
> ```
>
> **Low quality:**
> ```
> rtsp://{camera_ip}:554/s2
> ```
>
> - **Port:** 554 (on camera IP, not controller IP)
> - **Auth:** RTSP uses the Protect-generated credentials (displayed in Protect UI when RTSP is enabled)
> - **Video codec:** H.264
> - **Resolution:** 2688x1512 (s0), 1280x720 (s1), 640x360 (s2)
> - **Framerate:** Up to 24fps
>
> ### UniFi Protect API (on controller/UDM)
>
> Base URL: `https://{udm_ip}/proxy/protect/api` (on UDM) or `https://{ck_ip}:7443/api` (on Cloud Key)
>
> **List cameras:**
> ```
> GET /cameras
> Cookie: {session_cookie}
> X-CSRF-Token: {csrf_token}
> ```
>
> **Get camera details:**
> ```
> GET /cameras/{camera_id}
> ```
>
> **Snapshot:**
> ```
> GET /cameras/{camera_id}/snapshot?ts={timestamp}
> ```
> Returns JPEG snapshot. Omit `ts` for latest.
>
> **Motion events:**
> ```
> GET /events?type=motion&camera={camera_id}&start={epoch_ms}&end={epoch_ms}
> ```
>
> **Live WebSocket events:**
> ```
> WSS wss://{udm_ip}/proxy/protect/ws/updates
> Cookie: {session_cookie}
> ```
> Real-time motion and camera state updates via WebSocket.
>
> ### go2rtc Integration
>
> ```yaml
> streams:
>   unifi_front: "rtsp://{rtsp_user}:{rtsp_pass}@{camera_ip}:554/s0"
> ```

## AI Capabilities

> With RTSP and Protect API snapshots, AI integration follows the standard pattern. The Protect API's WebSocket event stream provides real-time motion notifications that can trigger AI analysis. Planned for M5 cameras milestone.

## Quirks & Notes

- **NVR/UDM required** -- The camera cannot function without a UniFi Protect controller. It will not serve RTSP or any other stream until adopted
- **RTSP not enabled by default** -- Must be explicitly enabled per-camera in Protect settings; this is an intentional Ubiquiti design choice
- **Self-signed TLS** -- The Protect API uses HTTPS with a self-signed certificate; clients must accept or skip TLS verification
- **RTSP credentials auto-generated** -- When RTSP is enabled, Protect generates a unique username/password per camera. These are displayed in the Protect UI and cannot be changed
- **No ONVIF** -- Despite being a professional camera, UniFi Protect cameras do not support ONVIF. Ubiquiti uses their proprietary protocol
- **PoE only** -- No WiFi option; must be wired with Ethernet. Requires 802.3af PoE (max 12.95W)
- **Protect API is undocumented** -- Ubiquiti does not officially document the Protect API. Community reverse-engineering provides the endpoints listed above
- **No subscription** -- All recording, detection, and features are included with the hardware cost. UniFi Protect has no recurring fees
- **Continuous recording** -- All footage is recorded 24/7 to the NVR's storage; no motion-only option needed to save bandwidth

## Similar Devices

> - [Amcrest IP Camera](amcrest-ip-camera.md) -- Professional PoE with RTSP + ONVIF + documented API
> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- WiFi/battery with native RTSP + ONVIF (more accessible)
> - [Nest Cam](nest-cam.md) -- Cloud-only alternative (opposite approach to local-first)
> - [Tapo C200](tapo-c200.md) -- Budget indoor with ONVIF/RTSP
