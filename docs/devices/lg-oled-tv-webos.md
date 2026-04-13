---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "lg-oled-tv-webos"
name: "LG OLED TV (webOS)"
manufacturer: "LG Electronics"
brand: "LG"
model: "OLED55C3PUA"
model_aliases: ["OLED55C4PUA", "OLED65C3PUA", "OLED65C4PUA", "OLED55G3PUA", "OLED77C3PUA", "OLED55B3PUA", "OLED65B3PUA", "OLED48C3PUA"]
device_type: "smart_tv"
category: "media"
product_line: "LG OLED"
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
    - "00:E0:91"        # LG Innotek (common in LG TVs)
    - "10:68:3F"        # LG Electronics
    - "2C:54:CF"        # LG Electronics
    - "30:B4:9E"        # LG Electronics
    - "34:4D:F7"        # LG Electronics
    - "38:8C:50"        # LG Electronics
    - "3C:BD:D8"        # LG Electronics
    - "40:B0:FA"        # LG Electronics
    - "58:A2:B5"        # LG Electronics
    - "64:99:5D"        # LG Electronics
    - "74:40:BE"        # LG Electronics
    - "78:5D:C8"        # LG Electronics
    - "88:C9:D0"        # LG Electronics
    - "A0:39:F7"        # LG Electronics
    - "A8:23:FE"        # LG Electronics
    - "AC:F1:08"        # LG Electronics
    - "B4:E6:2A"        # LG Electronics
    - "C4:36:6C"        # LG Electronics
    - "CC:2D:8C"        # LG Electronics
  mdns_services:
    - "_lgwebostv._tcp"       # primary webOS TV mDNS service
  mdns_txt_keys:
    - "deviceid"              # TV's unique identifier
    - "modelName"             # e.g., "OLED55C3PUA"
    - "modelNum"              # numeric model code
    - "lgTvVersion"           # webOS version
  default_ports: [3000, 3001, 8080, 9998, 18181, 1871]
  signature_ports: [3001]     # SSAP WebSocket over TLS is the strongest signal
  ssdp_search_target: "urn:schemas-upnp-org:device:MediaRenderer:1"
  ssdp_server_string: "WebOS/5.0 UPnP/1.0"
  hostname_patterns:
    - "^LG.*webOS.*TV"
    - "^LGWEBOSTV"
    - "^LGSmartTV"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 3000
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "WebSocket"
    body_contains: ""
    headers: {}
  - port: 1871
    path: "/udap/api/data?target=deviceinfo"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "LG"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "webos"
  polling_interval_sec: 10
  websocket_event: "webos:state"
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
  port: 3001
  transport: "WebSocket"
  encoding: "JSON"
  auth_method: "app_pairing"
  auth_detail: "Connect via WebSocket to wss://{ip}:3001. On first connection, the TV displays a pairing prompt on-screen. User accepts, and the TV returns a client-key token. Store the client-key and include it in the 'client-key' field of the hello/register payload on all subsequent connections."
  base_url_template: "wss://{ip}:3001"
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
  product_page: "https://www.lg.com/us/tvs/lg-oled55c3pua/"
  api_docs: "https://github.com/nicholasgasior/gowon"
  developer_portal: "https://webostv.developer.lge.com/"
  support: "https://www.lg.com/us/support"
  community_forum: "https://github.com/nicholasgasior/gowon/issues"
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["webos", "ssap", "websocket", "local_api", "oled", "4k", "hdr", "hdmi_arc", "airplay", "homekit", "cast"]
---

# LG OLED TV (webOS)

## What It Is

LG OLED TVs running webOS are premium 4K/8K displays manufactured by LG Electronics. They feature self-lit OLED pixels for perfect blacks, Dolby Vision/Atmos, HDMI 2.1, and a full smart TV platform with built-in streaming apps. Most importantly for Haus, they expose a local WebSocket API called SSAP (Simple Service Access Protocol) that allows comprehensive control -- power, volume, input switching, app launching, media transport, and more -- all without any cloud dependency. The webOS platform has been LG's TV OS since 2014, and the SSAP API has remained remarkably stable across generations. Recent models (2020+) also support Apple AirPlay 2, HomeKit, and Google Cast, making them among the most locally controllable TVs on the market.

## How Haus Discovers It

1. **OUI Match** -- During network scan, devices with MAC prefixes registered to LG Electronics (common prefixes include `A8:23:FE`, `CC:2D:8C`, `38:8C:50`, `78:5D:C8`, and many others) are flagged as potential LG devices. LG has a large number of OUI registrations, so MAC match alone is not sufficient for TV identification.

2. **mDNS Discovery** -- Browse for `_lgwebostv._tcp.local.` services. This is the strongest discovery signal. The TXT records include `deviceid` (unique TV identifier), `modelName` (e.g., "OLED55C3PUA"), and `lgTvVersion` (webOS version number). The presence of this service type definitively identifies an LG webOS TV.

3. **SSDP Discovery** -- The TV responds to UPnP M-SEARCH with `urn:schemas-upnp-org:device:MediaRenderer:1`. The response includes a description XML URL containing the model name, manufacturer "LG Electronics", and DLNA capabilities.

4. **Port Probe** -- Attempt WebSocket connection to port 3001 (TLS) or port 3000 (plaintext). Port 3001 with a self-signed TLS certificate accepting WebSocket upgrade is a strong indicator. Port 1871 may also respond with UDAP (Universal Device Access Protocol) XML data on older models.

5. **HTTP Fingerprint** -- `GET http://{ip}:1871/udap/api/data?target=deviceinfo` returns XML with manufacturer and model info on webOS 3.x and earlier. Newer models primarily use the WebSocket endpoint.

## Pairing / Authentication

LG webOS TVs use an on-screen pairing prompt model. No cloud account is required.

### Pairing Flow

1. Haus opens a WebSocket connection to `wss://{ip}:3001` (accepting the self-signed TLS certificate).

2. Haus sends a registration message:
   ```json
   {
     "type": "register",
     "id": "register_0",
     "payload": {
       "pairingType": "PROMPT",
       "manifest": {
         "manifestVersion": 1,
         "appVersion": "1.0.0",
         "signed": {
           "created": "20240101",
           "appId": "com.haus.hub",
           "vendorId": "com.haus",
           "localizedAppNames": {
             "": "Haus Hub"
           },
           "localizedVendorNames": {
             "": "Haus"
           },
           "permissions": [
             "LAUNCH",
             "LAUNCH_WEBAPP",
             "APP_TO_APP",
             "CLOSE",
             "TEST_OPEN",
             "TEST_PROTECTED",
             "CONTROL_AUDIO",
             "CONTROL_DISPLAY",
             "CONTROL_INPUT_JOYSTICK",
             "CONTROL_INPUT_MEDIA_RECORDING",
             "CONTROL_INPUT_MEDIA_PLAYBACK",
             "CONTROL_INPUT_TV",
             "CONTROL_POWER",
             "READ_APP_STATUS",
             "READ_CURRENT_CHANNEL",
             "READ_INPUT_DEVICE_LIST",
             "READ_NETWORK_STATE",
             "READ_RUNNING_APPS",
             "READ_TV_CHANNEL_LIST",
             "WRITE_NOTIFICATION"
           ],
           "serial": "haus001"
         },
         "permissions": [
           "LAUNCH",
           "LAUNCH_WEBAPP",
           "APP_TO_APP",
           "CLOSE",
           "TEST_OPEN",
           "TEST_PROTECTED",
           "CONTROL_AUDIO",
           "CONTROL_DISPLAY",
           "CONTROL_INPUT_JOYSTICK",
           "CONTROL_INPUT_MEDIA_RECORDING",
           "CONTROL_INPUT_MEDIA_PLAYBACK",
           "CONTROL_INPUT_TV",
           "CONTROL_POWER",
           "READ_APP_STATUS",
           "READ_CURRENT_CHANNEL",
           "READ_INPUT_DEVICE_LIST",
           "READ_NETWORK_STATE",
           "READ_RUNNING_APPS",
           "READ_TV_CHANNEL_LIST",
           "WRITE_NOTIFICATION"
         ],
         "signatures": [
           {
             "signatureVersion": 1,
             "signature": ""
           }
         ]
       },
       "client-key": ""
     }
   }
   ```

3. The TV displays an on-screen prompt: **"Allow Haus Hub to control this TV?"** with Accept/Deny buttons.

4. When the user accepts (via TV remote), the TV responds:
   ```json
   {
     "type": "registered",
     "id": "register_0",
     "payload": {
       "client-key": "abcdef1234567890abcdef1234567890"
     }
   }
   ```

5. Haus stores the `client-key`. On subsequent connections, include it in the registration payload to skip the on-screen prompt:
   ```json
   {
     "type": "register",
     "id": "register_0",
     "payload": {
       "client-key": "abcdef1234567890abcdef1234567890"
     }
   }
   ```

6. The TV responds immediately with `"type": "registered"` if the key is valid.

### Security Notes

- The TV uses a self-signed TLS certificate on port 3001. Haus must accept or pin this certificate.
- Client keys do not expire unless the TV is factory reset.
- The TV must be powered on (or in standby with "Quick Start+" enabled) to accept WebSocket connections.
- Wake-on-LAN (WOL) can be used to wake the TV from deep standby if the feature is enabled in settings.

## API Reference

All commands are sent as JSON over the WebSocket connection after registration. Each command has a `type`, `id`, and `uri` field. Responses are correlated by `id`.

### General Command Format

**Request:**
```json
{
  "type": "request",
  "id": "unique_request_id",
  "uri": "ssap://api/endpoint",
  "payload": {}
}
```

**Response:**
```json
{
  "type": "response",
  "id": "unique_request_id",
  "payload": {
    "returnValue": true,
    ...
  }
}
```

### Power

**Turn Off:**
```json
{
  "type": "request",
  "id": "power_off_1",
  "uri": "ssap://system/turnOff"
}
```

**Note:** There is no `turnOn` SSAP command because the WebSocket connection requires the TV to be on. Use Wake-on-LAN (WOL) to power on the TV. Send a magic packet to the TV's MAC address on UDP port 9. The TV must have "Quick Start+" or "Turn on via Wi-Fi" enabled in settings.

### Volume Control

**Get Volume:**
```json
{
  "type": "request",
  "id": "vol_get_1",
  "uri": "ssap://audio/getVolume"
}
```
Response payload: `{"returnValue": true, "volumeStatus": {"volume": 15, "muteStatus": false, "activeStatus": true, "adjustVolume": true}}`

**Set Volume:**
```json
{
  "type": "request",
  "id": "vol_set_1",
  "uri": "ssap://audio/setVolume",
  "payload": {"volume": 25}
}
```

**Volume Up:**
```json
{
  "type": "request",
  "id": "vol_up_1",
  "uri": "ssap://audio/volumeUp"
}
```

**Volume Down:**
```json
{
  "type": "request",
  "id": "vol_down_1",
  "uri": "ssap://audio/volumeDown"
}
```

**Set Mute:**
```json
{
  "type": "request",
  "id": "mute_1",
  "uri": "ssap://audio/setMute",
  "payload": {"mute": true}
}
```

### Input Selection

**Get Input List:**
```json
{
  "type": "request",
  "id": "input_list_1",
  "uri": "ssap://tv/getExternalInputList"
}
```
Response payload includes array of inputs: `{"devices": [{"id": "HDMI_1", "label": "HDMI 1", "connected": true, "appId": "com.webos.app.hdmi1"}, ...]}`

**Switch Input:**
```json
{
  "type": "request",
  "id": "input_switch_1",
  "uri": "ssap://tv/switchInput",
  "payload": {"inputId": "HDMI_1"}
}
```

### App Control

**Get App List:**
```json
{
  "type": "request",
  "id": "app_list_1",
  "uri": "ssap://com.webos.applicationManager/listApps"
}
```
Response: array of installed apps with `id`, `title`, `type`, `icon`, `visible`.

**Launch App:**
```json
{
  "type": "request",
  "id": "app_launch_1",
  "uri": "ssap://system.launcher/launch",
  "payload": {"id": "netflix"}
}
```

Common app IDs:
- `netflix` -- Netflix
- `youtube.leanback.v4` -- YouTube
- `amazon` -- Prime Video
- `com.webos.app.hdmi1` through `com.webos.app.hdmi4` -- HDMI inputs
- `com.webos.app.livetv` -- Live TV tuner
- `com.webos.app.browser` -- Web browser
- `disney+` or `com.disney.disneyplus-prod` -- Disney+
- `com.apple.appletv` -- Apple TV+

**Get Foreground App:**
```json
{
  "type": "request",
  "id": "fg_app_1",
  "uri": "ssap://com.webos.applicationManager/getForegroundAppInfo"
}
```
Response: `{"appId": "netflix", "processId": "...", "windowId": "..."}`

**Close App:**
```json
{
  "type": "request",
  "id": "app_close_1",
  "uri": "ssap://system.launcher/close",
  "payload": {"id": "netflix"}
}
```

### Media Playback

**Play:**
```json
{
  "type": "request",
  "id": "play_1",
  "uri": "ssap://media.controls/play"
}
```

**Pause:**
```json
{
  "type": "request",
  "id": "pause_1",
  "uri": "ssap://media.controls/pause"
}
```

**Stop:**
```json
{
  "type": "request",
  "id": "stop_1",
  "uri": "ssap://media.controls/stop"
}
```

**Rewind:**
```json
{
  "type": "request",
  "id": "rw_1",
  "uri": "ssap://media.controls/rewind"
}
```

**Fast Forward:**
```json
{
  "type": "request",
  "id": "ff_1",
  "uri": "ssap://media.controls/fastForward"
}
```

### TV Channel Control

**Get Current Channel:**
```json
{
  "type": "request",
  "id": "ch_get_1",
  "uri": "ssap://tv/getCurrentChannel"
}
```

**Get Channel List:**
```json
{
  "type": "request",
  "id": "ch_list_1",
  "uri": "ssap://tv/getChannelList"
}
```

**Set Channel:**
```json
{
  "type": "request",
  "id": "ch_set_1",
  "uri": "ssap://tv/openChannel",
  "payload": {"channelId": "3-1"}
}
```

**Channel Up/Down:**
```json
{"type": "request", "id": "ch_up_1", "uri": "ssap://tv/channelUp"}
{"type": "request", "id": "ch_dn_1", "uri": "ssap://tv/channelDown"}
```

### Notification (Toast)

**Show Toast:**
```json
{
  "type": "request",
  "id": "toast_1",
  "uri": "ssap://system.notifications/createToast",
  "payload": {"message": "Hello from Haus!"}
}
```

### System Info

**Get System Info:**
```json
{
  "type": "request",
  "id": "sys_info_1",
  "uri": "ssap://system/getSystemInfo"
}
```
Response: `{"features": {...}, "receiverType": "lg_tv", "modelName": "OLED55C3PUA"}`

**Get Software Info:**
```json
{
  "type": "request",
  "id": "sw_info_1",
  "uri": "ssap://com.webos.service.update/getCurrentSWInformation"
}
```

### Subscriptions (Real-Time Updates)

Any `request` URI can be used with `"type": "subscribe"` instead to receive continuous updates:

```json
{
  "type": "subscribe",
  "id": "vol_sub_1",
  "uri": "ssap://audio/getVolume"
}
```

This sends an initial response and then pushes updates whenever the volume changes. Useful subscriptions:
- `ssap://audio/getVolume` -- volume and mute changes
- `ssap://com.webos.applicationManager/getForegroundAppInfo` -- app switches
- `ssap://tv/getCurrentChannel` -- channel changes

### Mouse / Pointer Input Socket

After registration, request a pointer input socket:

```json
{
  "type": "request",
  "id": "pointer_1",
  "uri": "ssap://com.webos.service.networkinput/getPointerInputSocket"
}
```

Response includes a secondary WebSocket URL for sending remote-control button presses and pointer movements:
```
type:button\nname:ENTER\n
type:button\nname:HOME\n
type:button\nname:BACK\n
type:button\nname:LEFT\n
type:button\nname:RIGHT\n
type:button\nname:UP\n
type:button\nname:DOWN\n
type:move\ndx:10\ndy:5\ndown:0\n
type:click\n
type:scroll\ndx:0\ndy:1\n
```

Button names include: `HOME`, `BACK`, `ENTER`, `LEFT`, `RIGHT`, `UP`, `DOWN`, `RED`, `GREEN`, `YELLOW`, `BLUE`, `MENU`, `EXIT`, `1`-`9`, `0`, `DASH`.

## AI Capabilities

When the AI concierge "chats as" an LG OLED TV, it can:

- **Power off** the TV (power on via Wake-on-LAN)
- **Get and set volume** to specific levels, mute/unmute
- **List HDMI inputs** and switch between them ("switch to HDMI 2", "show me what's on the Xbox input")
- **Launch streaming apps** by name ("open Netflix", "put on YouTube")
- **Control media playback** -- play, pause, rewind, fast forward
- **Change channels** when watching live TV
- **Show notifications** as on-screen toasts ("display a message on the TV")
- **Report current state** -- what app is running, what input is selected, current volume
- **Navigate** using virtual remote button presses (home, back, arrow keys, enter)

The AI speaks in first person as the TV, referencing its OLED display capabilities and current viewing state.

## Quirks & Notes

- **Self-signed TLS on port 3001:** The TV generates its own TLS certificate. Haus must skip TLS verification or pin the cert on first connection. Port 3000 is available for unencrypted WebSocket but may not be enabled on all models.
- **No Power On via API:** The SSAP API cannot turn on the TV because the WebSocket server is only running when the TV is on. Use Wake-on-LAN (magic packet to TV's MAC on UDP port 9) with "Quick Start+" or "Turn on via Wi-Fi" enabled in TV settings. Alternatively, HDMI-CEC can wake the TV from connected devices.
- **Quick Start+ Required:** For WOL and fast WebSocket availability, the TV's "Quick Start+" mode must be enabled. Without it, the TV fully powers down and is unreachable on the network.
- **webOS Version Differences:** The SSAP API is consistent across webOS 3.x through 6.x (2016-2023+ TVs). Older webOS 1.x/2.x models (2014-2015) may lack some endpoints. The mDNS service name `_lgwebostv._tcp` is present on webOS 3.0+.
- **Connection Timeout:** The WebSocket connection may drop after ~10 minutes of inactivity. Implement a ping/pong heartbeat or reconnect logic.
- **Rate Limits:** No formal rate limit, but sending commands faster than ~100ms apart can cause dropped commands. Volume changes in particular should be debounced.
- **Dual-Band WiFi:** LG TVs support both 2.4GHz and 5GHz. The TV must be on the same subnet as the Haus hub for local API access.
- **HDMI-CEC Naming:** Input labels reflect CEC device names when available. HDMI_1 might show as "PlayStation 5" if CEC is active.
- **Screen Off vs Power Off:** `ssap://system/turnOff` fully powers off the TV. For "screen off" (audio continues, useful for music), some models support picture-off mode via settings commands, but this is not a standard SSAP endpoint.
- **Multiple Connections:** The TV accepts multiple simultaneous WebSocket clients, each with their own client-key.
- **Cast and AirPlay:** Models from 2020+ also support Google Cast (port 8008/8009) and AirPlay 2. These are separate control planes from SSAP.

## Similar Devices

- **samsung-smart-tv-tizen** -- Samsung's competing smart TV platform with its own WebSocket API
- **sony-bravia-google-tv** -- Sony's Google TV platform with REST API
- **roku-tv-streaming-stick** -- Roku OS TVs with ECP REST API
