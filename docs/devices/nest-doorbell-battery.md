---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "nest-doorbell-battery"
name: "Google Nest Doorbell (Battery)"
manufacturer: "Google LLC"
brand: "Google Nest"
model: "GA01318-US"
model_aliases: ["GWX3T", "G6AUD", "GA02076-US", "GA01318"]
device_type: "doorbell_camera"
category: "security"
product_line: "Nest"
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
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "18:B4:30"        # Google (Nest Labs Inc.)
    - "64:16:66"        # Google (Nest devices)
    - "7C:10:15"        # Google LLC
    - "B0:09:DA"        # Google LLC
    - "F8:0F:F9"        # Google LLC
    - "18:7F:88"        # Google LLC
    - "48:D6:D5"        # Google LLC
  mdns_services: []     # Nest doorbells do not advertise mDNS services
  mdns_txt_keys: []
  default_ports: []     # No open ports -- cloud-only device
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Nest-[A-Z0-9]+"
    - "^Google-Nest"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP services

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "nest"
  polling_interval_sec: 30
  websocket_event: "nest:state"
  setup_type: "oauth2"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities:
  - "camera_stream"
  - "camera_snapshot"
  - "doorbell"
  - "motion"
  - "battery_level"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Google OAuth 2.0 via Smart Device Management (SDM) API. Requires Google Cloud project, OAuth consent screen, and Device Access Console registration ($5 fee). Authorization via standard Google OAuth flow with scope https://www.googleapis.com/auth/sdm.service. Access tokens refreshed via https://oauth2.googleapis.com/token."
  base_url_template: "https://smartdevicemanagement.googleapis.com/v1"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "battery"
  mounting: "door"
  indoor_outdoor: "outdoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://store.google.com/product/nest_doorbell_battery"
  api_docs: "https://developers.google.com/nest/device-access/api"
  developer_portal: "https://console.nest.google.com/device-access"
  support: "https://support.google.com/googlenest/"
  community_forum: "https://www.googlenestcommunity.com/"
  image_url: ""
  fcc_id: "A4RGX3T"

# --- TAGS ---
tags: ["cloud_only", "google", "nest", "sdm_api", "doorbell", "battery", "on_device_ml", "person_detection"]
---

# Google Nest Doorbell (Battery)

## What It Is

The Google Nest Doorbell (Battery) is a wireless video doorbell with on-device machine learning for person, package, animal, and vehicle detection. It captures 960x1280 resolution video (3:4 portrait aspect ratio optimized for doorways), supports HDR and night vision, and features two-way audio with noise cancellation. The device runs on an internal rechargeable battery (6.3Ah) and can optionally be wired to existing doorbell wiring (8-24V AC) for continuous charging. All video processing and cloud features are managed through Google's cloud infrastructure, with local on-device ML providing intelligent alerts even without a Nest Aware subscription. Integration is via the official Google Smart Device Management (SDM) API, the same API used by Nest thermostats and cameras.

## How Haus Discovers It

1. **OUI Match** -- During network scan, devices with MAC prefixes associated with Google/Nest (`18:B4:30`, `64:16:66`, `7C:10:15`, etc.) are flagged as potential Nest devices. These OUIs are shared across all Google/Nest products.
2. **Hostname Pattern** -- Nest doorbells typically appear with hostnames matching `Nest-*` or `Google-Nest-*` in DHCP tables.
3. **No Local Probing** -- Like all Nest devices, the doorbell exposes no local API, no open ports, and no mDNS/SSDP services. Functional integration requires SDM API authentication.
4. **SDM API Discovery** -- Once OAuth is established, `GET /v1/enterprises/{project-id}/devices` lists all devices. Nest Doorbells are identified by the presence of `sdm.devices.traits.CameraLiveStream` and `sdm.devices.traits.DoorbellChime` traits.

## Pairing / Authentication

Nest devices use Google's Smart Device Management (SDM) API, which requires a one-time developer setup and standard Google OAuth 2.0 for user authorization.

### Prerequisites

1. **Google Cloud Project** -- Create a project in the Google Cloud Console with the SDM API enabled.
2. **Device Access Console** -- Register at `console.nest.google.com/device-access` (one-time $5 fee). Create a Device Access project and note the `project-id`.
3. **OAuth Consent Screen** -- Configure in Google Cloud Console with the `https://www.googleapis.com/auth/sdm.service` scope.
4. **OAuth Client ID** -- Create OAuth 2.0 credentials (Web application type) in Google Cloud Console.

### OAuth 2.0 Flow

1. **Authorization URL:** Redirect user to:
   ```
   https://nestservices.google.com/partnerconnections/{project-id}/auth
     ?redirect_uri={redirect_uri}
     &access_type=offline
     &prompt=consent
     &client_id={client_id}
     &response_type=code
     &scope=https://www.googleapis.com/auth/sdm.service
   ```
   Note: The authorization URL uses `nestservices.google.com`, NOT `accounts.google.com`. This is specific to the Device Access program and presents the user with a Nest-specific consent screen showing which devices to share.

2. **Token Exchange:** Exchange the authorization code:
   ```
   POST https://oauth2.googleapis.com/token
   Content-Type: application/x-www-form-urlencoded

   code={auth_code}&client_id={client_id}&client_secret={client_secret}&redirect_uri={redirect_uri}&grant_type=authorization_code
   ```

3. **Token Refresh:** Access tokens expire after 1 hour. Refresh via:
   ```
   POST https://oauth2.googleapis.com/token
   Content-Type: application/x-www-form-urlencoded

   client_id={client_id}&client_secret={client_secret}&refresh_token={refresh_token}&grant_type=refresh_token
   ```

### Security Notes

- OAuth refresh tokens do not expire but can be revoked by the user.
- The $5 Device Access fee is per-developer, not per-user.
- Google enforces rate limits: 10 queries per minute (QPM) per device for most endpoints, 5 QPM for camera streams.

## API Reference

All endpoints use `Authorization: Bearer {access_token}` header.

**Base URL:** `https://smartdevicemanagement.googleapis.com/v1`

### List Devices

```
GET /v1/enterprises/{project-id}/devices
```

Returns all authorized Nest devices. Each device includes:
- `name` -- resource name (`enterprises/{project-id}/devices/{device-id}`)
- `type` -- device type (e.g., `sdm.devices.types.DOORBELL`)
- `traits` -- map of device traits and their current values
- `parentRelations` -- room/structure assignments

### Doorbell-Specific Traits

**`sdm.devices.traits.DoorbellChime`** -- Indicates this device is a doorbell. Doorbell press events are delivered via Pub/Sub (see Events section).

**`sdm.devices.traits.CameraLiveStream`** -- Live stream capability:
- `maxVideoResolution` -- maximum stream resolution
- `videoCodecs` -- supported codecs (typically `["H264"]`)
- `audioCodecs` -- supported audio codecs

**`sdm.devices.traits.CameraEventImage`** -- Snapshot capability for event-triggered images.

**`sdm.devices.traits.CameraPerson`** -- Person detection events (requires Nest Aware for cloud, but on-device detection works without subscription for alerts).

**`sdm.devices.traits.CameraMotion`** -- Motion detection events.

### Generate Live Stream

```
POST /v1/enterprises/{project-id}/devices/{device-id}:executeCommand
Content-Type: application/json

{
  "command": "sdm.devices.commands.CameraLiveStream.GenerateRtspStream"
}
```

**Response:**
```json
{
  "results": {
    "streamUrls": {
      "rtspUrl": "rtsps://somehost:443/stream?auth=token..."
    },
    "streamExtensionToken": "...",
    "streamToken": "...",
    "expiresAt": "2024-01-01T00:00:00Z"
  }
}
```

Streams are RTSPS (RTSP over TLS) and expire after 5 minutes. Use `ExtendRtspStream` command with the `streamExtensionToken` to extend before expiry.

### Get Event Image

```
POST /v1/enterprises/{project-id}/devices/{device-id}:executeCommand
Content-Type: application/json

{
  "command": "sdm.devices.commands.CameraEventImage.GenerateImage",
  "params": {
    "eventId": "{event-id-from-pubsub}"
  }
}
```

Returns a URL and token to download the event snapshot image.

### Pub/Sub Events

Nest device events are delivered via Google Cloud Pub/Sub. Create a subscription to the topic `projects/sdm-prod/topics/enterprise-{project-id}`.

**Doorbell Press Event:**
```json
{
  "eventId": "...",
  "timestamp": "2024-01-01T00:00:00Z",
  "resourceUpdate": {
    "name": "enterprises/{project-id}/devices/{device-id}",
    "events": {
      "sdm.devices.events.DoorbellChime.Chime": {
        "eventSessionId": "...",
        "eventId": "..."
      }
    }
  }
}
```

**Motion Event:**
```json
{
  "resourceUpdate": {
    "events": {
      "sdm.devices.events.CameraMotion.Motion": {
        "eventSessionId": "...",
        "eventId": "..."
      }
    }
  }
}
```

**Person Event:**
```json
{
  "resourceUpdate": {
    "events": {
      "sdm.devices.events.CameraPerson.Person": {
        "eventSessionId": "...",
        "eventId": "..."
      }
    }
  }
}
```

## AI Capabilities

When the AI concierge integrates with the Nest Doorbell, it can:

- **Report doorbell press events** in real-time via Pub/Sub
- **Show event snapshots** from doorbell press, motion, and person detection events
- **Start live streams** and provide the RTSP URL for viewing
- **Report battery level** and charging status
- **List recent events** by type (doorbell, motion, person, package, animal, vehicle)
- **Provide context** -- "Someone pressed the doorbell at 3:42 PM. A person was detected."

## Quirks & Notes

- **Portrait Video:** The Nest Doorbell records in 3:4 portrait orientation (960x1280), unlike most security cameras that use landscape. This is optimized for seeing people head-to-toe at the door but can look odd in traditional camera UIs.
- **On-Device ML:** The doorbell runs ML models locally on a dedicated TPU chip for person, package, animal, and vehicle detection. These detections trigger smart alerts even without a Nest Aware subscription, but event history and continuous recording require a subscription.
- **Battery Life:** Google advertises approximately 2.5 months per charge with typical use. Heavy traffic areas or frequent live viewing will drain faster. Below 15C (59F) battery life decreases significantly.
- **Wired Mode Differences:** When hardwired, the doorbell enables 24/7 continuous recording (with Nest Aware) and a higher quality stream. On battery, it only records events.
- **SDM API Limitations:** The SDM API does not expose all features available in the Google Home app. Notable missing capabilities include: device settings changes, event clip playback (only snapshots), chime settings, and quiet time configuration.
- **Rate Limits:** Google enforces strict rate limits on the SDM API. Camera stream generation is limited to 5 requests per minute per device. Exceeding limits returns HTTP 429.
- **RTSPS Streams:** Streams use RTSPS (TLS-encrypted RTSP) and are short-lived (5 minutes). You must call `ExtendRtspStream` before expiry to keep the stream alive. Maximum stream duration with extensions is not documented but is approximately 5 minutes per extension.
- **Pub/Sub Costs:** Google Cloud Pub/Sub has its own pricing. For typical home use, costs are negligible (well within the free tier), but the infrastructure setup is more complex than simple polling.
- **Shared OUI Prefixes:** Google uses the same MAC OUI prefixes across all Nest/Google Home products (thermostats, cameras, speakers, Chromecast, Wi-Fi routers). MAC alone cannot distinguish a doorbell from a Nest Hub -- SDM API device listing is required.
- **2nd Gen (2022):** Google released a wired-only 2nd gen Nest Doorbell in 2022 (model G9AJC). It has a similar form factor but different internals. The SDM API treats both generations the same way.

## Similar Devices

- **ring-video-doorbell-4** -- Amazon's competing battery video doorbell (cloud-only, unofficial API)
- **arlo-essential-video-doorbell** -- Arlo's cloud-based video doorbell
- **eufy-video-doorbell-dual** -- Eufy's locally-stored alternative
- **nest-learning-thermostat** -- Same SDM API integration, different device type
