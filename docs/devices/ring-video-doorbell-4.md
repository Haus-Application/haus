---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ring-video-doorbell-4"
name: "Ring Video Doorbell 4"
manufacturer: "Ring LLC (Amazon)"
brand: "Ring"
model: "B08JNR77QY"
model_aliases: ["5AT3T5", "8VR1S1-0EN0"]
device_type: "doorbell_camera"
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
  mac_prefixes:
    - "2C:AA:8E"        # Ring LLC
    - "50:27:C7"        # Ring LLC
    - "44:73:D6"        # Ring LLC
    - "34:3E:A4"        # Ring LLC
    - "5C:47:5E"        # Ring LLC
    - "0C:96:E6"        # Ring LLC
  mdns_services: []     # Ring doorbells do not advertise mDNS services
  mdns_txt_keys: []
  default_ports: []     # No open ports -- cloud-only device
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Ring-[0-9a-f]+"
    - "^ring-doorbell"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP services

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "ring"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "camera_stream"
  - "camera_snapshot"
  - "doorbell"
  - "motion"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Ring uses a proprietary OAuth-like flow. Initial auth via POST to https://oauth.ring.com/oauth/token with username/password/2FA, returns access_token and refresh_token. All API calls use Authorization: Bearer {token} header. Tokens expire and must be refreshed."
  base_url_template: "https://api.ring.com/clients_api"
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
  product_page: "https://ring.com/products/video-doorbell-4"
  api_docs: ""
  developer_portal: ""
  support: "https://support.ring.com/"
  community_forum: "https://community.ring.com/"
  image_url: ""
  fcc_id: "2AEUP-1090"

# --- TAGS ---
tags: ["cloud_only", "amazon", "doorbell", "pre_roll", "battery", "no_local_api"]
---

# Ring Video Doorbell 4

## What It Is

The Ring Video Doorbell 4 is a Wi-Fi connected video doorbell manufactured by Ring LLC, an Amazon subsidiary. It features 1080p HD video, two-way audio, Pre-Roll Video Preview (4 seconds of black-and-white video captured before the motion event using a low-power secondary sensor), advanced motion detection with customizable motion zones, and infrared night vision. It runs on a rechargeable battery pack or can be hardwired to existing doorbell wiring (8-24V AC) for trickle charging. The device is entirely cloud-dependent -- all video processing, storage (requires Ring Protect subscription), and live streaming route through Ring's cloud infrastructure.

## How Haus Discovers It

1. **OUI Match** -- During network scan, devices with MAC prefixes `2C:AA:8E`, `50:27:C7`, `44:73:D6`, `34:3E:A4`, `5C:47:5E`, or `0C:96:E6` are flagged as Ring devices.
2. **Hostname Pattern** -- Ring doorbells typically appear in DHCP tables with hostnames matching `Ring-*` or `ring-doorbell-*`.
3. **No Local Probing** -- Ring devices expose no local API, no open ports, and no mDNS/SSDP services. Discovery is limited to network-level identification (MAC + hostname). Functional integration requires cloud API authentication.

## Pairing / Authentication

Ring does not offer an official public API or developer program. All third-party integrations rely on the unofficial Ring API, which Ring has periodically attempted to restrict.

### Authentication Flow

1. **Initial Login:** `POST https://oauth.ring.com/oauth/token` with body:
   ```json
   {
     "client_id": "ring_official_android",
     "grant_type": "password",
     "password": "{password}",
     "username": "{email}",
     "scope": "client"
   }
   ```
2. **2FA Required:** Ring enforces two-factor authentication. If 2FA is enabled, the initial request returns HTTP 412 with a `tsv_state` indicating 2FA is pending. The user receives an SMS/email code.
3. **2FA Verification:** Re-send the same POST with the `X-Time-Based-OTP-Token` header set to the 2FA code.
4. **Token Response:** On success, returns `access_token` and `refresh_token`. The access token is short-lived (minutes); the refresh token lasts longer.
5. **Token Refresh:** `POST https://oauth.ring.com/oauth/token` with `grant_type: "refresh_token"`.
6. **Hardware ID:** A unique `hardware_id` (UUID) must be generated per client and included in all auth requests. This identifies the "device" to Ring's system.

### Security Notes

- Ring has no official developer program or public API documentation.
- The unofficial API is reverse-engineered and subject to change without notice.
- Ring has added rate limiting and device verification to slow automated access.
- Two-factor authentication is mandatory on all Ring accounts.

## API Reference

All endpoints use `Authorization: Bearer {access_token}` header.

### List Devices

```
GET https://api.ring.com/clients_api/ring_devices
```

Returns all Ring devices on the account including doorbells, cameras, chimes, and alarm devices. Each doorbell includes `id`, `description`, `firmware_version`, `kind` (device type), `address`, `latitude`, `longitude`, and battery status.

### Get Doorbell Health

```
GET https://api.ring.com/clients_api/doorbots/{id}/health
```

Returns Wi-Fi signal strength (RSSI), firmware version, network name, and connection quality metrics.

### Ding Events (Doorbell Press / Motion)

```
GET https://api.ring.com/clients_api/dings/active
```

Returns currently active ding events. Each event includes `kind` ("ding" for button press, "motion" for motion detected), `doorbot_id`, `state`, and SIP connection details for live streaming.

### Event History

```
GET https://api.ring.com/clients_api/doorbots/{id}/history?limit=20
```

Returns recent events with timestamps, event types, and recording availability.

### Live Stream

Ring uses SIP (Session Initiation Protocol) over WebRTC for live video streaming. When a ding event is active, the API provides SIP credentials and ICE servers for establishing a peer-to-peer (via TURN relay) video connection. The SIP negotiation flow is complex and involves Ring's proprietary SIP proxy infrastructure.

### Snapshot

```
POST https://api.ring.com/clients_api/snapshots/timestamps
```

Requests a new snapshot from the device. The snapshot is captured asynchronously and can be retrieved via a separate endpoint.

```
GET https://api.ring.com/clients_api/snapshots/image/{id}
```

Returns the most recent snapshot as a JPEG image.

## AI Capabilities

AI integration is not currently planned for Ring devices due to the cloud-only, unofficial API dependency. If implemented, the AI concierge could:

- Report doorbell press events and motion alerts
- Show recent snapshots
- List recent event history
- Report device health (battery level, Wi-Fi signal)

## Quirks & Notes

- **No Official API:** Ring does not provide a sanctioned developer API. All third-party access relies on reverse-engineered endpoints that Ring periodically changes or restricts. Home Assistant's Ring integration has been broken and fixed multiple times over the years.
- **Pre-Roll Video:** The Doorbell 4's headline feature is Pre-Roll Video Preview -- 4 seconds of black-and-white footage captured before a motion event triggers. This uses a secondary low-power sensor that records continuously in a rolling buffer. The pre-roll is stitched onto the beginning of the recorded event in the cloud.
- **Subscription Required:** Video recording, snapshot history, and person detection require a Ring Protect subscription ($3.99/month Basic or $12.99/month Plus as of 2024). Without a subscription, you only get real-time notifications and live view.
- **Battery Life:** Typical battery life is 6-12 months depending on activity level and temperature. Cold weather significantly reduces battery performance.
- **SIP Streaming Complexity:** Live video uses SIP-over-WebRTC, which is significantly more complex than standard RTSP. This makes local streaming integration challenging.
- **Amazon Sidewalk:** Ring devices participate in Amazon Sidewalk, a shared neighborhood network using BLE and 900MHz ISM band. This can be disabled in the Ring app but is enabled by default.
- **Dual-Band WiFi:** The Doorbell 4 supports both 2.4GHz and 5GHz Wi-Fi (802.11 b/g/n/ac). Previous models were 2.4GHz only.
- **Hardwire Bypass:** If hardwired, the device charges the battery but does not run on wired power directly. If the battery dies while hardwired, it may still go offline briefly.

## Similar Devices

- **nest-doorbell-battery** -- Google's competing battery video doorbell with similar cloud-dependent model
- **arlo-essential-video-doorbell** -- Arlo's cloud-based video doorbell
- **eufy-video-doorbell-dual** -- Eufy's locally-stored alternative with dual cameras
