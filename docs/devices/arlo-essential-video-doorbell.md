---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "arlo-essential-video-doorbell"
name: "Arlo Essential Video Doorbell"
manufacturer: "Arlo Technologies, Inc."
brand: "Arlo"
model: "AVD2001"
model_aliases: ["AVD2001B", "AVD2001-100NAS"]
device_type: "doorbell_camera"
category: "security"
product_line: "Arlo Essential"
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
  mac_prefixes:
    - "9C:F6:DD"        # Arlo Technologies
    - "84:B5:D8"        # Arlo Technologies
    - "E4:4E:2D"        # Arlo Technologies (newer)
  mdns_services: []     # Arlo doorbells do not advertise mDNS services
  mdns_txt_keys: []
  default_ports: []     # No open ports -- cloud-only device
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Arlo"
    - "^arlo-doorbell"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP services

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "arlo"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "camera_stream"
  - "doorbell"
  - "motion"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Arlo uses a proprietary authentication flow. Login via POST to https://ocapi-app.arlo.com/api/auth with email/password, returns a token. 2FA is required. Subsequent API calls via POST to https://myapi.arlo.com/hmsweb/users/devices with Auth header. Real-time events use SSE or MQTT via Arlo's cloud."
  base_url_template: "https://myapi.arlo.com/hmsweb"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "battery"
  mounting: "door"
  indoor_outdoor: "outdoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.arlo.com/doorbell/AVD2001.html"
  api_docs: ""
  developer_portal: ""
  support: "https://www.arlo.com/en-us/support/"
  community_forum: "https://community.arlo.com/"
  image_url: ""
  fcc_id: "2AEOS-AVD2001"

# --- TAGS ---
tags: ["cloud_only", "doorbell", "battery", "head_to_toe", "no_official_api", "2fa_required"]
---

# Arlo Essential Video Doorbell

## What It Is

The Arlo Essential Video Doorbell is a Wi-Fi connected wireless video doorbell manufactured by Arlo Technologies. Its headline feature is a 180-degree diagonal field of view with a head-to-toe aspect ratio (1:1 square at 1536x1536 resolution), allowing users to see visitors from head to feet and packages on the ground without a separate package camera. It features two-way audio, HDR video, night vision, and on-device motion detection. The device connects directly to Wi-Fi (2.4GHz only) without requiring an Arlo base station, making setup simpler than Arlo's pro-grade cameras. It runs on a rechargeable battery or can be hardwired to existing doorbell wiring.

## How Haus Discovers It

1. **OUI Match** -- During network scan, devices with MAC prefixes `9C:F6:DD`, `84:B5:D8`, or `E4:4E:2D` are flagged as Arlo devices.
2. **Hostname Pattern** -- Arlo doorbells may appear with hostnames matching `Arlo*` in DHCP tables, though hostname reporting varies.
3. **No Local Probing** -- Arlo Essential devices expose no local API, no open ports, and no mDNS/SSDP services. They connect outbound to Arlo's cloud only. Functional integration requires cloud API access.

## Pairing / Authentication

Arlo does not provide an official public API or developer program. Third-party integrations use the unofficial Arlo API.

### Authentication Flow

1. **Initial Login:** `POST https://ocapi-app.arlo.com/api/auth` with:
   ```json
   {
     "email": "{email}",
     "password": "{password}",
     "EnvSource": "prod",
     "language": "en"
   }
   ```
2. **2FA Challenge:** Arlo mandates 2FA. The auth response indicates 2FA is required and provides a `factorId`. Arlo supports SMS, email, or TOTP authenticator app for the second factor.
3. **2FA Verification:** `POST https://ocapi-app.arlo.com/api/auth/verify` with the 2FA code and `factorId`.
4. **Token Response:** Returns an authentication token and a `userId` used for subsequent API calls.
5. **Validated Token:** After 2FA verification, retrieve the final validated token from `POST https://ocapi-app.arlo.com/api/auth/validate`.

### Security Notes

- Arlo aggressively enforces 2FA on all accounts since 2023 -- no option to disable it.
- The unofficial API changes periodically as Arlo migrates infrastructure.
- Arlo has moved between multiple backend architectures (MQTT, SSE, WebSocket) over the years, making the API a moving target.
- Session tokens expire and must be re-authenticated periodically.

## API Reference

All endpoints use `Authorization: {token}` header (no "Bearer" prefix in Arlo's API).

### List Devices

```
GET https://myapi.arlo.com/hmsweb/users/devices
```

Returns all Arlo devices on the account. Doorbell entries include:
- `deviceId` -- unique device identifier
- `deviceName` -- user-assigned name
- `deviceType` -- device type string (e.g., `doorbell`)
- `properties` -- device properties including battery level, signal strength
- `state` -- current device state

### Get Library (Recordings)

```
POST https://myapi.arlo.com/hmsweb/users/library
Content-Type: application/json

{
  "dateFrom": "20240101",
  "dateTo": "20240102"
}
```

Returns recorded events for the date range including motion, doorbell press, and smart detection events.

### Start Stream

```
POST https://myapi.arlo.com/hmsweb/users/devices/startStream
Content-Type: application/json

{
  "from": "{userId}_{web}",
  "to": "{deviceId}",
  "action": "set",
  "resource": "cameras/{deviceId}",
  "publishResponse": true,
  "properties": {
    "activityState": "startUserStream",
    "cameraId": "{deviceId}"
  }
}
```

Returns an RTSP or HLS stream URL for live viewing. Stream URLs are temporary and expire.

### Real-Time Events

Arlo uses an SSE (Server-Sent Events) endpoint for real-time notifications:

```
GET https://myapi.arlo.com/hmsweb/client/subscribe
```

Events include doorbell press (`doorbellAlert`), motion detection, and device state changes. The event payload varies by device type.

## AI Capabilities

AI integration is not currently planned for Arlo devices due to the unofficial API dependency. If implemented, the AI concierge could:

- Report doorbell press and motion events
- Show event thumbnails from the recording library
- Report device battery level and connectivity status
- Start live streams on demand

## Quirks & Notes

- **Head-to-Toe View:** The 180-degree FOV captures a 1:1 square image that shows the full doorstep area from head to toe. This eliminates the "package blind spot" common in landscape-oriented doorbells but requires UI adjustments for the non-standard aspect ratio.
- **No Base Station Required:** Unlike Arlo Pro and Ultra cameras, the Essential Doorbell connects directly to Wi-Fi. This simplifies setup but means no local storage option -- all recordings go through Arlo's cloud.
- **2.4GHz Only:** The Essential Doorbell only supports 2.4GHz Wi-Fi, which can be limiting in congested wireless environments but provides better range than 5GHz.
- **Subscription Tiers:** Without a subscription, you get live view, two-way audio, and basic motion alerts. Arlo Secure ($7.99/month single device or $17.99/month unlimited) adds 30-day cloud recording, activity zones, smart notifications (person, vehicle, package, animal), and 4K recording where supported.
- **API Instability:** Arlo has migrated their backend infrastructure multiple times. The API has shifted from a WebSocket-based event system to MQTT to SSE. Third-party integrations (including Home Assistant's) frequently need updates.
- **No Local Storage on Essential:** The Essential line does not support USB storage or microSD. Only Arlo Pro/Ultra models with a SmartHub support local recording.
- **Battery Life:** Arlo advertises up to 6 months on a single charge, though real-world usage with frequent motion events typically yields 2-4 months.
- **Silent Mode:** When the internal chime is disabled, the doorbell still sends push notifications but does not make an audible sound. An Arlo Chime accessory can be paired for indoor audible alerts.

## Similar Devices

- **ring-video-doorbell-4** -- Amazon's competing cloud-based video doorbell
- **nest-doorbell-battery** -- Google's cloud-based video doorbell with SDM API
- **eufy-video-doorbell-dual** -- Local storage alternative with dual cameras
