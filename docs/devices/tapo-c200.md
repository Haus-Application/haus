---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "tapo-c200"
name: "TP-Link Tapo C200"
manufacturer: "TP-Link Technologies Co., Ltd."
brand: "TP-Link Tapo"
model: "Tapo C200"
model_aliases: ["TC60", "Tapo C210", "Tapo C200 V2", "Tapo C220"]
device_type: "tapo_camera"
category: "security"
product_line: "Tapo"
release_year: 2020
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["5C:E9:31", "98:25:4A", "B0:A7:B9", "30:DE:4B", "E8:48:B8", "A8:42:A1", "68:FF:7B", "1C:3B:F3"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [554, 2020, 80, 443]
  signature_ports: [2020, 554]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^tapo.*", "^Tapo.*", "^C200.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 443
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "TP-Link"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "tapo"
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
  auth_detail: "RTSP and ONVIF use the camera's credentials set in the Tapo app (Account > Advanced Settings > Camera Account). ONVIF on port 2020. RTSP on port 554. Tapo cloud API uses a separate encrypted JSON protocol on port 443."
  base_url_template: "rtsp://{ip}:554/stream1"
  tls: false
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
  product_page: "https://www.tapo.com/product/smart-camera/tapo-c200/"
  api_docs: ""
  developer_portal: ""
  support: "https://www.tapo.com/support/"
  community_forum: "https://community.tp-link.com"
  image_url: ""
  fcc_id: "TE7C200"

# --- TAGS ---
tags: ["rtsp", "onvif", "local-first", "budget", "pan-tilt", "1080p", "night-vision", "microsd", "no-subscription", "tp-link", "motion-detection", "two-way-audio"]
---

# TP-Link Tapo C200

## What It Is

> The TP-Link Tapo C200 is an affordable ($25-35) indoor WiFi pan/tilt security camera with 1080p video, 360-degree horizontal and 114-degree vertical rotation, night vision (up to 9m), motion detection, two-way audio, and microSD recording. The standout feature for home automation is native ONVIF and RTSP support -- no firmware hacking, no hub required. After setting a camera account in the Tapo app, RTSP and ONVIF streams are accessible locally on the network. Combined with its low price, it is one of the most accessible cameras for local integration with Haus.

## How Haus Discovers It

1. **OUI match** -- TP-Link Tapo devices use MAC prefixes including `5C:E9:31`, `98:25:4A`, `B0:A7:B9`, `30:DE:4B`, `E8:48:B8`, `A8:42:A1`, `68:FF:7B`, and `1C:3B:F3` (TP-Link OUIs)
2. **ONVIF discovery** -- Standard ONVIF WS-Discovery probe (UDP multicast to 239.255.255.250:3702). Tapo cameras respond with ONVIF device service URL on port 2020
3. **Port probe** -- TCP check on port 2020 (ONVIF) and 554 (RTSP)
4. **RTSP probe** -- RTSP `OPTIONS` to `rtsp://{ip}:554/stream1` confirms stream availability
5. **HTTP fingerprint** -- HTTPS on port 443 may contain "TP-Link" branding

## Pairing / Authentication

> Initial setup requires the Tapo app:
>
> 1. Create TP-Link account and add camera via Tapo app
> 2. Camera connects to WiFi during setup
> 3. Set camera account: In Tapo app, go to Camera Settings > Advanced Settings > Camera Account
> 4. Create a username and password for the camera account -- this is used for RTSP/ONVIF
>
> **Important:** The "camera account" (for RTSP/ONVIF) is separate from the TP-Link cloud account. You must explicitly create it in Advanced Settings. Without this, RTSP/ONVIF authentication will fail.

## API Reference

> ### RTSP Streams
>
> **Main stream (high quality):**
> ```
> rtsp://{username}:{password}@{ip}:554/stream1
> ```
>
> **Sub stream (low bandwidth):**
> ```
> rtsp://{username}:{password}@{ip}:554/stream2
> ```
>
> - **Port:** 554
> - **Auth:** Digest authentication (camera account credentials)
> - **Video codec:** H.264
> - **Main resolution:** 1920x1080 (1080p)
> - **Sub resolution:** 640x360
> - **Framerate:** Up to 15fps
>
> ### ONVIF
>
> - **ONVIF port:** 2020 (non-standard; most cameras use 80 or 8000)
> - **Device service URL:** `http://{ip}:2020/onvif/device_service`
> - **Media service URL:** `http://{ip}:2020/onvif/media_service`
> - **PTZ service URL:** `http://{ip}:2020/onvif/ptz_service`
> - **Profiles:** Profile S (streaming)
> - **PTZ:** Pan/tilt control available via ONVIF PTZ service
> - **Events:** Motion detection events via ONVIF pull-point subscription
>
> ### Snapshot
>
> Via ONVIF GetSnapshotUri or directly:
> ```
> GET http://{username}:{password}@{ip}:80/snapshot.cgi
> ```
> Availability varies by firmware version.
>
> ### go2rtc Integration
>
> ```yaml
> streams:
>   tapo_living_room: "rtsp://{user}:{pass}@{ip}:554/stream1"
> ```
>
> Or via ONVIF auto-discovery:
> ```yaml
> streams:
>   tapo_living_room: "onvif://{user}:{pass}@{ip}:2020"
> ```

## AI Capabilities

> With native RTSP and ONVIF, AI integration follows the standard camera pattern: go2rtc captures snapshots, Claude vision analyzes scenes. PTZ control could allow AI-directed camera movement. Not planned until M5 cameras milestone.

## Quirks & Notes

- **ONVIF on port 2020** -- Unlike most cameras that serve ONVIF on port 80, Tapo uses port 2020. This is important for ONVIF discovery and must be accounted for in Haus's discovery engine
- **Camera account is separate** -- The RTSP/ONVIF credentials are NOT the TP-Link cloud account credentials. Users must explicitly create a "camera account" in Advanced Settings. This is the most common setup mistake
- **Pan/tilt via ONVIF** -- The C200's pan/tilt motor is controllable via ONVIF PTZ service, enabling AI-directed camera movement
- **Privacy mode** -- The camera can physically rotate to face a wall ("privacy mode") via the Tapo app or ONVIF PTZ preset
- **Night vision** -- IR LEDs provide up to 9m night vision range; no color night vision (no spotlight)
- **MicroSD recording** -- Supports up to 512GB microSD for local continuous recording
- **No subscription** -- All features work without any subscription or cloud dependency
- **Firmware updates** -- Tapo regularly updates firmware; ONVIF/RTSP support has been consistent across recent firmware versions
- **Tapo C210/C220** -- The C210 (2K) and C220 (2K with AI detection) are newer models with the same protocol support but higher resolution
- **Cloud protocol** -- The Tapo cloud API (port 443) uses an encrypted JSON protocol with RSA key exchange; this is separate from the RTSP/ONVIF local protocols

## Similar Devices

> - [Wyze Cam v3](wyze-cam-v3.md) -- Similar price, requires RTSP firmware flash
> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- Battery-powered with native RTSP/ONVIF
> - [Amcrest IP Camera](amcrest-ip-camera.md) -- More professional, same local protocol support
> - [Eufy Security Cam 2C](eufy-security-cam-2c.md) -- Battery camera, RTSP via HomeBase
