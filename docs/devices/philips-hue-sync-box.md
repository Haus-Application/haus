---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "philips-hue-sync-box"
name: "Philips Hue Play HDMI Sync Box"
manufacturer: "Signify Netherlands B.V."
brand: "Philips Hue"
model: "HSB1"
model_aliases: ["929002275802", "555227"]
device_type: "hue_sync_box"
category: "media"
product_line: "Hue"
release_year: 2019
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["ethernet", "wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "00:17:88"        # Signify / Philips Lighting BV
    - "EC:B5:FA"        # Signify Netherlands B.V.
  mdns_services:
    - "_huesync._tcp"
  mdns_txt_keys:
    - "uniqueid"        # unique device identifier
    - "name"            # user-assigned device name
    - "devicetype"      # "HSB1"
  default_ports: [443]
  signature_ports: [443]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^HueSyncBox-"
    - "^HSB1-"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 443
    path: "/api/v1"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"device\""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "hue_sync"
  polling_interval_sec: 5
  websocket_event: "hue_sync:state"
  setup_type: "app_pairing"
  ai_chattable: true
  haus_milestone: "M7"

# --- CAPABILITIES ---
capabilities:
  - "on_off"
  - "input_select"
  - "media_playback"
  - "brightness"
  - "scenes"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "Registration via POST /api/v1/registrations with appName, instanceName. User must press and hold the physical button on the Sync Box for ~3 seconds during registration. Returns an access token used as Bearer token in Authorization header."
  base_url_template: "https://{ip}/api/v1"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.philips-hue.com/en-us/p/hue-play-hdmi-sync-box/046677555221"
  api_docs: "https://developers.meethue.com/develop/hue-entertainment/hue-hdmi-sync-box-api/"
  developer_portal: "https://developers.meethue.com/"
  support: "https://www.philips-hue.com/en-us/support"
  community_forum: "https://developers.meethue.com/forum/"
  image_url: ""
  fcc_id: "2ABA6-HSB1"

# --- TAGS ---
tags: ["hdmi", "sync", "entertainment", "media", "local_api", "signify", "passthrough"]
---

# Philips Hue Play HDMI Sync Box

## What It Is

The Philips Hue Play HDMI Sync Box (model HSB1) is a standalone HDMI passthrough device that analyzes the video signal from up to 4 HDMI input sources and synchronizes connected Philips Hue lights to match the on-screen colors in real time. It sits between your media sources (game consoles, streaming sticks, Blu-ray players) and your TV, passing through the video signal while extracting color data to drive ambient lighting effects via the Hue Entertainment API.

Unlike Hue bulbs and lightstrips, the Sync Box has its **own Ethernet/Wi-Fi network connection and its own local HTTPS API**, separate from the Hue Bridge. However, it communicates with the Hue Bridge over the local network to control lights -- it does not directly control Zigbee devices. The Sync Box supports HDMI 2.0b with 4K 60Hz HDR10+ and Dolby Vision passthrough.

Initial setup requires the Hue Sync mobile app and an internet connection to link the device to your Hue account, but once configured, sync operation works locally.

## How Haus Discovers It

1. **OUI Match** -- MAC prefix `00:17:88` or `EC:B5:FA` flags the device as a Signify product.
2. **mDNS Discovery** -- Browse for `_huesync._tcp.local.` services. The Sync Box advertises TXT records including `uniqueid`, `name`, and `devicetype` ("HSB1"). This is distinct from the bridge's `_hue._tcp` service.
3. **Port Probe** -- Port 443 open with self-signed TLS certificate.
4. **HTTP Fingerprint** -- `GET https://{ip}/api/v1` returns device information. Without authentication, a limited response is returned containing the device name and type. The presence of the Sync Box API endpoint structure confirms the device.
5. **Hostname Pattern** -- Sync Boxes typically appear as `HueSyncBox-XXXX` on the network.

## Pairing / Authentication

The Sync Box uses a registration-based authentication model similar to the Hue Bridge but with its own API.

### Registration Flow

1. **Initiate Registration:**
   ```
   POST https://{ip}/api/v1/registrations
   Content-Type: application/json

   {
     "appName": "haus",
     "instanceName": "hub"
   }
   ```

2. **User Action:** The user must **press and hold the physical button** on the Sync Box for approximately **3 seconds** until the LED blinks. This must happen within the registration window (similar to the bridge link-button concept but requires a longer press).

3. **Response:**
   ```json
   {
     "registrationId": "abc123...",
     "accessToken": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
   }
   ```

4. **Subsequent Requests:** Use the access token as a Bearer token:
   ```
   Authorization: Bearer {accessToken}
   ```

### Cloud Setup Requirement

The initial setup of the Sync Box (connecting it to a Hue Bridge, configuring Entertainment areas) requires the **Hue Sync mobile app** and a **Signify cloud account**. This is a one-time requirement. After setup, the local API can be used for all control operations.

## API Reference

**Base URL:** `https://{ip}/api/v1`

All endpoints use HTTPS with self-signed TLS. Authentication via `Authorization: Bearer {token}` header.

### Get Device State

```
GET /api/v1/device
```

**Response:**
```json
{
  "name": "Living Room Sync Box",
  "deviceType": "HSB1",
  "uniqueId": "AABBCCDDEEFF",
  "ipAddress": "192.168.1.100",
  "apiLevel": 7,
  "firmwareVersion": "2.12.0",
  "buildNumber": 1234,
  "lastCheckedUpdate": "2024-01-15T10:30:00Z",
  "updatableBuildNumber": null,
  "updatableFirmwareVersion": null,
  "wifiState": "connected"
}
```

### Get Execution State

```
GET /api/v1/execution
```

**Response:**
```json
{
  "mode": "video",
  "syncActive": true,
  "hdmiActive": true,
  "hdmiSource": "input1",
  "hueTarget": "Entertainment area 1",
  "brightness": 200,
  "lastSyncMode": "video",
  "video": {
    "intensity": "high",
    "backgroundLighting": true
  },
  "game": {
    "intensity": "extreme",
    "backgroundLighting": true
  },
  "music": {
    "intensity": "high"
  }
}
```

### Set Execution State

```
PUT /api/v1/execution
```

```json
{
  "mode": "video",
  "syncActive": true,
  "hdmiSource": "input2",
  "brightness": 180,
  "video": {
    "intensity": "moderate"
  }
}
```

**Mode values:** `"passthrough"` (no sync), `"video"`, `"game"`, `"music"`

**Intensity values:** `"subtle"`, `"moderate"`, `"high"`, `"extreme"`

**Brightness:** 0-200 (integer scale, not percentage)

**HDMI Source:** `"input1"`, `"input2"`, `"input3"`, `"input4"`

### Get HDMI Info

```
GET /api/v1/hdmi
```

**Response:**
```json
{
  "input1": {
    "name": "Xbox",
    "type": "game",
    "status": "linked",
    "lastSyncMode": "game"
  },
  "input2": {
    "name": "Apple TV",
    "type": "video",
    "status": "linked",
    "lastSyncMode": "video"
  },
  "input3": {
    "name": "",
    "type": "generic",
    "status": "unplugged",
    "lastSyncMode": "video"
  },
  "input4": {
    "name": "",
    "type": "generic",
    "status": "unplugged",
    "lastSyncMode": "video"
  },
  "output": {
    "name": "LG TV",
    "type": "generic",
    "status": "linked"
  },
  "contentSpecs": "4K 60Hz HDR10",
  "videoSyncSupported": true,
  "audioSyncSupported": true
}
```

### Set HDMI Input Names

```
PUT /api/v1/hdmi
```

```json
{
  "input1": {"name": "PlayStation 5", "type": "game"},
  "input2": {"name": "Apple TV 4K", "type": "video"}
}
```

**Type values:** `"generic"`, `"video"`, `"game"`, `"music"`, `"xbox"`, `"playstation"`, `"nintendoswitch"`, `"pc"`, `"bluray"`, `"satellite"`, `"streaming"`

### Control Power (Standby)

```
PUT /api/v1/execution
```

```json
{"mode": "powersave"}
```

To wake from standby:
```json
{"mode": "passthrough"}
```

## AI Capabilities

When Haus supports this device, the AI concierge will be able to:

- **Switch HDMI inputs** ("switch to the PlayStation", "show input 2")
- **Start/stop sync** ("start syncing lights", "stop the light show")
- **Change sync mode** ("switch to game mode", "use music mode")
- **Adjust intensity** ("make the sync more subtle", "turn up the intensity")
- **Adjust brightness** of synced lights ("dim the ambient lights a bit")
- **Report current state** ("The Sync Box is on input 1 (Xbox), syncing in game mode at high intensity")
- **Control power** ("put the sync box in standby")

## Quirks & Notes

- **Cloud Required for Initial Setup:** Unlike the bridge, the Sync Box requires the Hue Sync mobile app and a Signify cloud account for first-time configuration (linking to a bridge, setting up entertainment areas). This is a one-time requirement, but it cannot be avoided.
- **Not a Zigbee Controller:** The Sync Box does not directly control Zigbee devices. It communicates with the Hue Bridge over the local network (using the Entertainment API streaming protocol), and the bridge then controls the lights. Both the Sync Box and bridge must be on the same network.
- **HDMI 2.0b Limitation:** The original HSB1 supports HDMI 2.0b (4K 60Hz). It does NOT support HDMI 2.1 features like 4K 120Hz, VRR, or ALLM. Signify released the Hue Play HDMI Sync Box 8K (model HSB2) in late 2023 with HDMI 2.1 support.
- **Brightness Scale:** The brightness value in the API is 0-200 (integer), not 0-100%. This is unique to the Sync Box API and does not match the bridge API's percentage scale.
- **Self-Signed TLS:** Like the bridge, the Sync Box uses self-signed certificates. Haus must skip or pin certificate verification.
- **Entertainment Areas:** The Sync Box syncs lights assigned to a Hue Entertainment area (configured in the Hue Sync app). It uses the bridge's Entertainment streaming API (UDP port 2100) for low-latency color updates.
- **CEC Support:** The Sync Box supports HDMI-CEC for automatic input switching and can trigger TV power on/off.
- **Button Double-Function:** The physical button on the Sync Box toggles sync on/off in normal operation, and is used as the registration button for API pairing (long press).
- **Firmware Updates:** Applied via the Hue Sync mobile app. Updates can add features and change API behavior.
- **Wi-Fi vs Ethernet:** The Sync Box supports both Wi-Fi and Ethernet. Ethernet is recommended for reliable sync operation, especially in environments with Wi-Fi congestion.

## Similar Devices

- **philips-hue-bridge** -- Required for the Sync Box to control lights; they work as a pair
- **philips-hue-lightstrip-plus** -- Popular companion for TV ambient lighting with the Sync Box
- **philips-hue-bulb-a19** -- Can be synced via the Sync Box for room-wide ambient effects
