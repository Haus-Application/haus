---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "samsung-smart-tv-tizen"
name: "Samsung Smart TV (Tizen)"
manufacturer: "Samsung Electronics"
brand: "Samsung"
model: "QN55S95CAFXZA"
model_aliases: ["QN65S95CAFXZA", "UN55TU8000FXZA", "QN65QN90CAFXZA", "UN50TU7000FXZA", "QN55Q80CAFXZA", "QN65QN85CAFXZA", "UN43TU7000FXZA", "QN55S90CAFXZA"]
device_type: "smart_tv"
category: "media"
product_line: "Samsung Smart TV"
release_year: 2023
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "ethernet", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "00:07:AB"        # Samsung Electronics
    - "00:12:FB"        # Samsung Electronics
    - "00:15:99"        # Samsung Electronics
    - "00:16:32"        # Samsung Electronics
    - "00:17:D5"        # Samsung Electronics
    - "00:18:AF"        # Samsung Electronics
    - "00:1C:43"        # Samsung Electronics
    - "00:1E:E2"        # Samsung Electronics
    - "00:21:19"        # Samsung Electronics
    - "00:24:54"        # Samsung Electronics
    - "00:26:37"        # Samsung Electronics
    - "08:D4:6A"        # Samsung Electronics
    - "10:1D:C0"        # Samsung Electronics
    - "14:49:E0"        # Samsung Electronics
    - "28:39:5E"        # Samsung Electronics
    - "30:C7:AE"        # Samsung Electronics
    - "34:14:5F"        # Samsung Electronics
    - "40:16:3B"        # Samsung Electronics
    - "50:01:D9"        # Samsung Electronics
    - "54:BD:79"        # Samsung Electronics
    - "64:B5:C6"        # Samsung Electronics
    - "78:AB:BB"        # Samsung Electronics
    - "8C:71:F8"        # Samsung Electronics
    - "94:8B:C1"        # Samsung Electronics
    - "A8:7C:01"        # Samsung Electronics
    - "BC:14:85"        # Samsung Electronics
    - "C4:73:1E"        # Samsung Electronics
    - "D0:66:7B"        # Samsung Electronics
    - "E4:7C:F9"        # Samsung Electronics
    - "F0:25:B7"        # Samsung Electronics
    - "F4:7B:5E"        # Samsung Electronics
  mdns_services:
    - "_samsung-tv._tcp"      # Samsung TV mDNS (not always present)
  mdns_txt_keys: []
  default_ports: [8001, 8002, 9197, 9110, 9119, 7676, 55000]
  signature_ports: [8002]     # WSS endpoint is the strongest signal
  ssdp_search_target: "urn:samsung.com:device:RemoteControlReceiver:1"
  ssdp_server_string: "SHP, UPnP/1.0, Samsung UPnP SDK/1.0"
  hostname_patterns:
    - "^Samsung.*TV"
    - "^\\[TV\\]"
    - "^SAMSUNG-TV"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8001
    path: "/api/v2/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"device\""
    headers: {}
  - port: 8002
    path: "/api/v2/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"device\""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "samsung_tv"
  polling_interval_sec: 10
  websocket_event: "samsung_tv:state"
  setup_type: "app_pairing"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities:
  - "on_off"
  - "volume"
  - "input_select"
  - "media_playback"

# --- PROTOCOL ---
protocol:
  type: "websocket_json"
  port: 8002
  transport: "WebSocket"
  encoding: "JSON"
  auth_method: "app_pairing"
  auth_detail: "Connect to wss://{ip}:8002/api/v2/channels/samsung.remote.control?name={base64_app_name}. On first connection, TV shows an Allow/Deny prompt. After allowing, a token is returned in the WebSocket response. Include this token in subsequent connections via the 'token' query parameter."
  base_url_template: "wss://{ip}:8002/api/v2/channels/samsung.remote.control"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "display"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.samsung.com/us/televisions-home-theater/tvs/"
  api_docs: "https://developer.samsung.com/smarttv/develop/extension-libraries/smart-view-sdk/receiver-apps.html"
  developer_portal: "https://developer.samsung.com/smarttv"
  support: "https://www.samsung.com/us/support/"
  community_forum: "https://developer.samsung.com/forum"
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["tizen", "websocket", "local_api", "smartthings", "4k", "8k", "qled", "neo_qled", "hdr", "hdmi_arc", "bixby"]
---

# Samsung Smart TV (Tizen)

## What It Is

Samsung Smart TVs run the Tizen operating system and are manufactured by Samsung Electronics. The lineup spans QLED, Neo QLED, OLED, and Crystal UHD panels across various screen sizes. For local control, Samsung TVs expose a WebSocket API on ports 8001 (HTTP) and 8002 (HTTPS/WSS) that supports remote control key simulation, app launching, and device info queries. Samsung also offers the SmartThings cloud API for more structured control, but the local WebSocket API works without internet. Tizen-based Samsung TVs have been shipping since 2015, with the WebSocket remote control API available on 2016+ models.

## How Haus Discovers It

1. **OUI Match** -- Samsung Electronics has an extremely large number of registered MAC prefixes (hundreds of OUIs). A MAC match to Samsung narrows candidates but does not distinguish TVs from phones, tablets, refrigerators, or other Samsung devices. Additional probing is required.

2. **SSDP Discovery** -- The TV responds to UPnP M-SEARCH with search target `urn:samsung.com:device:RemoteControlReceiver:1`. This is highly specific to Samsung TVs and is the best discovery signal. The SSDP response includes the TV model and device description XML URL.

3. **mDNS Discovery** -- Some Samsung TVs advertise `_samsung-tv._tcp.local.` but this is not consistent across all models and firmware versions. SSDP is more reliable.

4. **HTTP Fingerprint** -- `GET http://{ip}:8001/api/v2/` returns a JSON object with device information:
   ```json
   {
     "device": {
       "FrameTVSupport": "false",
       "GamePadSupport": "true",
       "ImeSyncedSupport": "true",
       "OS": "Tizen",
       "TokenAuthSupport": "true",
       "VoiceSupport": "true",
       "countryCode": "US",
       "description": "Samsung DTV RCR",
       "developerIP": "0.0.0.0",
       "developerMode": "0",
       "duid": "uuid:...",
       "firmwareVersion": "Unknown",
       "id": "uuid:...",
       "ip": "192.168.1.100",
       "model": "23_KANTM2E_QLED",
       "modelName": "QN55S95CAFXZA",
       "name": "[TV] Samsung Q95C 55\"",
       "networkType": "wireless",
       "resolution": "3840x2160",
       "smartHubAgreement": "true",
       "type": "Samsung SmartTV",
       "udn": "uuid:...",
       "wifiMac": "AA:BB:CC:DD:EE:FF"
     }
   }
   ```
   The `type: "Samsung SmartTV"` and `OS: "Tizen"` fields definitively confirm a Samsung TV.

5. **Port Probe** -- Ports 8001 and 8002 being open simultaneously, combined with the `/api/v2/` fingerprint, strongly identifies a Samsung TV.

## Pairing / Authentication

Samsung TVs use an on-screen Allow/Deny prompt similar to LG TVs. Newer models (2018+) support token-based authentication.

### Pairing Flow

1. Haus opens a WebSocket connection to:
   ```
   wss://{ip}:8002/api/v2/channels/samsung.remote.control?name={base64_name}
   ```
   Where `{base64_name}` is the base64-encoded app name (e.g., base64 of "Haus Hub" = `SGF1cyBIdWI=`).

2. On first connection, the TV displays an on-screen prompt: **"Haus Hub is requesting permission to connect. Allow?"**

3. When the user selects "Allow" on the TV remote, the WebSocket connection is established and the TV sends a response containing a token:
   ```json
   {
     "data": {
       "clients": [...],
       "id": "...",
       "token": "12345678"
     },
     "event": "ms.channel.connect"
   }
   ```

4. Store the `token` value. On subsequent connections, include it as a query parameter:
   ```
   wss://{ip}:8002/api/v2/channels/samsung.remote.control?name={base64_name}&token={token}
   ```

5. With a valid token, the TV connects immediately without showing a prompt.

### Older Models (Pre-2018)

Older Samsung TVs (2016-2017) use port 8001 without TLS and without token authentication:
```
ws://{ip}:8001/api/v2/channels/samsung.remote.control?name={base64_name}
```
These models still show the Allow/Deny prompt but do not issue tokens. The TV remembers allowed apps by name.

### Security Notes

- Port 8002 uses a self-signed TLS certificate. Haus must accept the certificate.
- Tokens do not expire unless the TV is factory reset or the app is explicitly blocked in TV settings.
- The TV must be on (not in standby) for the initial pairing. After pairing, WOL can be used if enabled.

## API Reference

All commands are sent as JSON messages over the WebSocket connection on the `samsung.remote.control` channel.

### Remote Control Keys

The primary control mechanism is simulating remote control button presses:

```json
{
  "method": "ms.remote.control",
  "params": {
    "Cmd": "Click",
    "DataOfCmd": "KEY_VOLUP",
    "Option": "false",
    "TypeOfRemote": "SendRemoteKey"
  }
}
```

**Cmd values:**
- `Click` -- single press
- `Press` -- key down (hold)
- `Release` -- key up (release after hold)

### Common Key Codes

**Power & Navigation:**
- `KEY_POWER` -- toggle power (also `KEY_POWEROFF`)
- `KEY_HOME` -- home/smart hub
- `KEY_MENU` -- settings menu
- `KEY_RETURN` -- back
- `KEY_EXIT` -- exit app
- `KEY_UP`, `KEY_DOWN`, `KEY_LEFT`, `KEY_RIGHT` -- navigation
- `KEY_ENTER` -- select/confirm
- `KEY_INFO` -- show info overlay

**Volume & Audio:**
- `KEY_VOLUP` -- volume up
- `KEY_VOLDOWN` -- volume down
- `KEY_MUTE` -- toggle mute

**Channel:**
- `KEY_CHUP` -- channel up
- `KEY_CHDOWN` -- channel down
- `KEY_0` through `KEY_9` -- numeric keys
- `KEY_PRECH` -- previous channel
- `KEY_GUIDE` -- TV guide

**Media Transport:**
- `KEY_PLAY` -- play
- `KEY_PAUSE` -- pause
- `KEY_STOP` -- stop
- `KEY_FF` -- fast forward
- `KEY_REWIND` -- rewind
- `KEY_REC` -- record

**Input:**
- `KEY_SOURCE` -- input source selector
- `KEY_HDMI` -- cycle HDMI inputs
- `KEY_HDMI1`, `KEY_HDMI2`, `KEY_HDMI3`, `KEY_HDMI4` -- direct HDMI input

**Color Buttons:**
- `KEY_RED`, `KEY_GREEN`, `KEY_YELLOW`, `KEY_BLUE`

### App Launch

Launch apps using the `ms.channel.emit` method:

```json
{
  "method": "ms.channel.emit",
  "params": {
    "event": "ed.apps.launch",
    "to": "host",
    "data": {
      "appId": "Netflix",
      "action_type": "DEEP_LINK"
    }
  }
}
```

Common app IDs:
- `Netflix` -- Netflix
- `YouTube` -- YouTube
- `AmazonInstantVideo` -- Prime Video
- `DisneyPlus` -- Disney+
- `AppleTV` -- Apple TV+
- `Hulu` -- Hulu
- `11101200001` -- YouTube (numeric alternative)
- `3201907018807` -- Disney+ (numeric alternative)
- `3201512006785` -- Plex

### Get Installed Apps

```json
{
  "method": "ms.channel.emit",
  "params": {
    "event": "ed.installedApp.get",
    "to": "host"
  }
}
```

Response arrives as a separate event with the list of installed applications.

### Device Info (REST)

The REST endpoint is always available without WebSocket:

```
GET http://{ip}:8001/api/v2/
```

Returns device model, name, OS, resolution, network type, and unique identifiers.

### Wake-on-LAN

Samsung TVs support WOL when "Network Standby" is enabled in Power settings. Send a magic packet (6x `0xFF` followed by 16x repetitions of the TV's MAC address) to UDP port 9.

For 2019+ models with SmartThings integration, the SmartThings cloud API can also power on the TV remotely via `POST /devices/{deviceId}/commands` with `{"component":"main","capability":"switch","command":"on"}`.

## AI Capabilities

When the AI concierge "chats as" a Samsung TV, it can:

- **Power on/off** the TV (off via KEY_POWER, on via Wake-on-LAN)
- **Adjust volume** up, down, or mute
- **Switch inputs** -- direct HDMI selection or cycle through sources
- **Launch apps** by name ("open Netflix", "start YouTube")
- **Control media playback** -- play, pause, stop, rewind, fast forward
- **Navigate menus** -- simulate remote control button presses
- **Change channels** for live TV viewing
- **Report device info** -- model name, resolution, current network state

The AI speaks in first person as the TV, aware of its Samsung/Tizen identity and capabilities.

## Quirks & Notes

- **Samsung Has Many OUIs:** Samsung Electronics has hundreds of registered MAC prefixes. MAC-based identification alone is insufficient -- a Samsung MAC could be a phone, tablet, refrigerator, washer, or TV. The SSDP `RemoteControlReceiver` service or the HTTP `/api/v2/` endpoint must be used to confirm a TV.
- **Port 8001 vs 8002:** Port 8001 is plaintext HTTP/WS, port 8002 is TLS/WSS. Both expose the same API. Prefer 8002 for security. Some older models only have 8001. Some newer models may only respond on 8002.
- **Token Format Changes:** The token format and authentication mechanism has changed slightly across model years. 2016-2017 models use name-based tracking. 2018+ models use numeric tokens. 2022+ models may use longer alphanumeric tokens.
- **Frame TV:** Samsung Frame TVs (`"FrameTVSupport": "true"`) have additional art-mode capabilities that can be controlled via the same WebSocket API with additional commands for art display.
- **No Structured State API:** Unlike LG's SSAP which returns structured state (volume level, current input), Samsung's local API is primarily key-simulation. To know the current volume, you can only send KEY_VOLUP/KEY_VOLDOWN -- there is no "get current volume" query via the local API. SmartThings cloud API provides structured state.
- **Connection Keep-Alive:** The WebSocket connection may timeout. Send periodic ping frames to maintain the connection.
- **SmartThings Cloud API:** For structured state and richer control, Samsung offers the SmartThings API at `https://api.smartthings.com/v1/`. This requires OAuth2 and a Samsung account but provides proper state queries, capability-based control, and push notifications. Haus may use this as a supplement to the local API.
- **Power State Detection:** The TV drops off the network entirely when fully powered off (not standby). Detecting power state requires either polling the HTTP endpoint or using SmartThings cloud push events.

## Similar Devices

- **lg-oled-tv-webos** -- LG's competing platform with richer local SSAP API
- **sony-bravia-google-tv** -- Sony's Google TV with REST-based control
- **roku-tv-streaming-stick** -- Roku with excellent local ECP API
