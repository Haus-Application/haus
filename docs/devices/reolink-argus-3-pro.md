---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "reolink-argus-3-pro"
name: "Reolink Argus 3 Pro"
manufacturer: "Reolink Innovation Inc."
brand: "Reolink"
model: "Argus 3 Pro"
model_aliases: ["Reolink Argus 3 Pro WiFi", "B09DPNTTLS"]
device_type: "reolink_camera"
category: "security"
product_line: "Reolink"
release_year: 2021
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["EC:71:DB", "B4:6D:83", "B8:A5:89", "9C:8E:CD"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [554, 80, 8000]
  signature_ports: [554, 8000]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^reolink.*", "^Reolink.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: "Reolink"
    server_header: ""
    body_contains: "Reolink"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "reolink"
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
  auth_method: "basic_auth"
  auth_detail: "RTSP and ONVIF use the camera's admin credentials (set during initial setup via Reolink app or web UI). Default username is 'admin'. Web UI on port 80 uses session-based auth."
  base_url_template: "rtsp://{ip}:554/h264Preview_01_main"
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
  product_page: "https://reolink.com/product/argus-3-pro/"
  api_docs: "https://support.reolink.com/hc/en-us/articles/360007010473-How-to-Live-View-Reolink-Cameras-via-VLC-Media-Player"
  developer_portal: ""
  support: "https://support.reolink.com"
  community_forum: "https://community.reolink.com"
  image_url: ""
  fcc_id: "2ATKB-ARGUS3PRO"

# --- TAGS ---
tags: ["rtsp", "onvif", "local-first", "battery", "solar-option", "2k-resolution", "spotlight", "person-vehicle-detection", "microsd", "no-subscription"]
---

# Reolink Argus 3 Pro

## What It Is

> The Reolink Argus 3 Pro is a wireless battery-powered security camera with exceptional local integration support. It provides 2K (2560x1440) video, color night vision via integrated spotlight, person/vehicle detection, two-way audio, and a 122-degree field of view. What makes it standout for Haus is its native RTSP and ONVIF support -- no special firmware, no hub required, no cloud dependency. The camera connects directly to WiFi and exposes standard protocols out of the box. It can run entirely local with microSD recording and no internet connection required. Compatible with the Reolink Solar Panel for indefinite outdoor operation.

## How Haus Discovers It

1. **OUI match** -- Reolink cameras use MAC prefixes including `EC:71:DB`, `B4:6D:83`, `B8:A5:89`, and `9C:8E:CD`
2. **Port probe** -- TCP check on ports 554 (RTSP), 80 (web UI), and 8000 (Reolink proprietary)
3. **ONVIF discovery** -- Standard ONVIF WS-Discovery probe (UDP multicast to 239.255.255.250:3702). Reolink cameras respond with ONVIF device service URL
4. **HTTP fingerprint** -- Web UI on port 80 contains "Reolink" in the page title
5. **RTSP probe** -- RTSP `OPTIONS` request to `rtsp://{ip}:554/` confirms stream availability

## Pairing / Authentication

> Setup via Reolink app or web UI:
>
> 1. Power on camera and connect to its setup WiFi hotspot (or use Reolink app QR code scan)
> 2. Configure WiFi credentials
> 3. Set admin password (this becomes the RTSP/ONVIF password)
> 4. Camera joins home WiFi network directly -- no hub required
>
> **RTSP/ONVIF credentials:** Same admin username and password set during initial setup. Default username is `admin`.

## API Reference

> ### RTSP Streams
>
> Reolink provides two stream profiles:
>
> **Main stream (high quality):**
> ```
> rtsp://admin:{password}@{ip}:554/h264Preview_01_main
> ```
>
> **Sub stream (low bandwidth):**
> ```
> rtsp://admin:{password}@{ip}:554/h264Preview_01_sub
> ```
>
> - **Port:** 554
> - **Auth:** Basic authentication (admin credentials)
> - **Video codec:** H.264 (main), H.264 (sub)
> - **Main resolution:** 2560x1440 (2K)
> - **Sub resolution:** 640x360
> - **Framerate:** Up to 15fps (battery model)
>
> ### ONVIF
>
> - **ONVIF port:** 8000
> - **Device service URL:** `http://{ip}:8000/onvif/device_service`
> - **Media service URL:** `http://{ip}:8000/onvif/media_service`
> - **Profiles:** Profile S (streaming), Profile T (advanced video)
> - **PTZ:** Not applicable (fixed camera)
> - **Events:** Motion detection events via ONVIF pull-point subscription
>
> ### Snapshot
>
> ```
> GET http://admin:{password}@{ip}/cgi-bin/api.cgi?cmd=Snap&channel=0&rs=abc123
> ```
> Returns JPEG snapshot. The `rs` parameter is a random string for cache-busting.
>
> ### Web UI API (JSON-RPC style)
>
> ```
> POST http://{ip}/cgi-bin/api.cgi?cmd=Login
> Content-Type: application/json
>
> [{"cmd": "Login", "action": 0, "param": {"User": {"userName": "admin", "password": "{password}"}}}]
> ```
> Returns a session token used for subsequent API calls (battery level, motion settings, etc.).
>
> ### go2rtc Integration
>
> ```yaml
> streams:
>   reolink_front: "rtsp://admin:{pass}@{ip}:554/h264Preview_01_main"
> ```

## AI Capabilities

> With native RTSP, AI integration follows the standard camera pattern: go2rtc captures snapshots from the RTSP stream, Claude vision analyzes the scene. Person/vehicle detection events from ONVIF can trigger AI notifications. Planned for M5 cameras milestone.

## Quirks & Notes

- **Best-in-class local support** -- Native RTSP + ONVIF with no firmware hacks, no hubs, no cloud required. This is the gold standard for Haus camera integration
- **Battery limitations on streaming** -- RTSP streaming wakes the camera from sleep and drains battery quickly. Continuous RTSP streaming is not practical on battery; use motion-triggered streams
- **Solar panel option** -- The Reolink Solar Panel (sold separately) provides enough power for continuous operation in most climates, making always-on RTSP feasible
- **Spotlight** -- Integrated white LED spotlight activates on motion for color night vision; configurable brightness and schedule
- **MicroSD recording** -- Supports up to 128GB microSD for local recording; accessible via Reolink app or FTP
- **No cloud subscription** -- All features including person/vehicle detection work without any subscription
- **Port 8000** -- Reolink's proprietary port for ONVIF and device communication; some routers/firewalls may block this
- **Dual-band WiFi** -- Supports 2.4GHz and 5GHz WiFi for better connectivity
- **FTP upload** -- Can upload snapshots/clips to an FTP server on motion events

## Similar Devices

> - [Amcrest IP Camera](amcrest-ip-camera.md) -- PoE/WiFi with native RTSP + ONVIF (wired, more professional)
> - [Tapo C200](tapo-c200.md) -- Budget indoor with native ONVIF + RTSP
> - [Eufy Security Cam 2C](eufy-security-cam-2c.md) -- Battery camera with RTSP via HomeBase (less direct)
> - [Wyze Cam v3](wyze-cam-v3.md) -- Budget wired with RTSP firmware (requires flash)
> - [UniFi Protect G4 Bullet](unifi-protect-g4-bullet.md) -- Professional PoE with RTSP
