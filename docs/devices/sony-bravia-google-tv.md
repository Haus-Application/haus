---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "sony-bravia-google-tv"
name: "Sony Bravia (Google TV)"
manufacturer: "Sony Group Corporation"
brand: "Sony"
model: "XR-55A95L"
model_aliases: ["XR-65A95L", "XR-55A80L", "XR-65A80L", "KD-55X85L", "KD-65X85L", "XR-75X90L", "KD-43X85L", "XR-55X90L", "XR-65X90L", "XR-85X90L"]
device_type: "smart_tv"
category: "media"
product_line: "Sony Bravia"
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
    - "00:13:A9"        # Sony Corporation
    - "00:1A:80"        # Sony Corporation
    - "00:1D:BA"        # Sony Corporation
    - "00:1E:A4"        # Sony Corporation
    - "00:24:BE"        # Sony Corporation
    - "04:5D:4B"        # Sony Corporation
    - "10:4F:A8"        # Sony Corporation
    - "24:21:AB"        # Sony Corporation
    - "30:17:C8"        # Sony Corporation
    - "40:B8:37"        # Sony Corporation
    - "54:42:49"        # Sony Corporation
    - "70:2A:D5"        # Sony Corporation
    - "78:84:3C"        # Sony Corporation
    - "AC:9B:0A"        # Sony Corporation
    - "B4:52:7E"        # Sony Corporation
    - "D8:D4:3C"        # Sony Corporation
    - "FC:0F:E6"        # Sony Corporation
  mdns_services:
    - "_googlecast._tcp"      # Google Cast (built into Google TV)
    - "_sony-bravia._tcp"     # Sony Bravia-specific service (newer models)
  mdns_txt_keys:
    - "fn"                    # friendly name (Cast)
    - "md"                    # model name (Cast)
    - "mn"                    # manufacturer name
  default_ports: [80, 443, 8008, 8009, 10000, 20000]
  signature_ports: [80]       # BRAVIA REST API on port 80
  ssdp_search_target: "urn:schemas-sony-com:service:IRCC:1"
  ssdp_server_string: "Linux/3.x UPnP/1.0 Sony-BRAVIA"
  hostname_patterns:
    - "^BRAVIA"
    - "^Sony.*TV"
    - "^XBR-"
    - "^KD-"
    - "^XR-"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/sony/system"
    method: "POST"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"result\""
    headers: {}
  - port: 80
    path: "/sony/avContent"
    method: "POST"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "bravia"
  polling_interval_sec: 10
  websocket_event: "bravia:state"
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
  type: "http_rest"
  port: 80
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "Pre-Shared Key (PSK) set in TV settings (Settings > Network > Home Network > IP Control > Pre-Shared Key). Send PSK in 'X-Auth-PSK' header with every request. Alternatively, use PIN-based pairing for an auth cookie."
  base_url_template: "http://{ip}/sony"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "display"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://electronics.sony.com/tv-video/televisions"
  api_docs: "https://pro-bravia.sony.net/develop/integrate/ip-control/index.html"
  developer_portal: "https://pro-bravia.sony.net/develop/"
  support: "https://www.sony.com/en/support"
  community_forum: ""
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["bravia", "google_tv", "android_tv", "rest_api", "ip_control", "cast", "ircc", "4k", "8k", "hdr", "oled", "mini_led", "homekit", "airplay"]
---

# Sony Bravia (Google TV)

## What It Is

Sony Bravia TVs are premium displays manufactured by Sony Group Corporation, running Google TV (formerly Android TV) as their smart platform. The lineup includes OLED (A-series), Mini LED (X-series), and Full Array LED models. Sony Bravia TVs expose a comprehensive local REST API called the BRAVIA IP Control API (also known as the Sony Scalar API or JSON-RPC API) on port 80. This API supports power control, volume, input switching, app launching, and media playback. Additionally, since these TVs run Google TV, they also support Google Cast protocol (ports 8008/8009) for media casting. Newer models (2020+) also support Apple AirPlay 2 and HomeKit.

## How Haus Discovers It

1. **OUI Match** -- Devices with MAC prefixes registered to Sony Corporation (e.g., `AC:9B:0A`, `D8:D4:3C`, `FC:0F:E6`) are flagged as potential Sony devices. Like Samsung, Sony makes many products, so additional probing is needed to confirm a TV.

2. **SSDP Discovery** -- The TV responds to UPnP M-SEARCH with search target `urn:schemas-sony-com:service:IRCC:1` (Infrared Compatible Control over IP). This is highly specific to Sony Bravia TVs and is the best discovery signal. The response includes the TV model and service description URL.

3. **mDNS Discovery** -- The TV advertises `_googlecast._tcp.local.` (as a Cast-enabled device). The TXT record `md` contains the model name (e.g., "BRAVIA XR-55A95L") and `mn` contains "Sony". Some newer models also advertise `_sony-bravia._tcp.local.`.

4. **HTTP Fingerprint** -- `POST http://{ip}/sony/system` with a JSON-RPC body returns system information. This is the definitive identification method:
   ```json
   {
     "method": "getSystemInformation",
     "id": 1,
     "params": [],
     "version": "1.0"
   }
   ```
   Response includes `product: "TV"`, `model: "XR-55A95L"`, `name: "BRAVIA"`.

5. **Port Probe** -- Port 80 open with `/sony/system` responding to POST requests, combined with port 8008 (Cast), strongly identifies a Sony Bravia TV.

## Pairing / Authentication

Sony Bravia TVs support two authentication methods for the local API.

### Method 1: Pre-Shared Key (Recommended for Haus)

1. User navigates to: **Settings > Network > Home Network > IP Control > Pre-Shared Key**
2. User sets a PSK string (e.g., "haus123")
3. User also enables: **Settings > Network > Home Network > IP Control > Simple IP Control > On**
4. Haus includes the PSK in every request header:
   ```
   X-Auth-PSK: haus123
   ```

This is the simplest method and requires no handshake or token exchange.

### Method 2: PIN-Based Pairing

1. Haus sends a registration request:
   ```json
   POST http://{ip}/sony/accessControl
   Content-Type: application/json

   {
     "method": "actRegister",
     "id": 1,
     "params": [
       {
         "clientid": "haus:hub:001",
         "nickname": "Haus Hub",
         "level": "private"
       },
       [{"value": "yes", "function": "WOL"}]
     ],
     "version": "1.0"
   }
   ```

2. The TV displays a 4-digit PIN on screen.

3. Haus sends the same request again with the PIN as Basic Auth:
   ```
   Authorization: Basic base64("haus:hub:001" + ":" + "1234")
   ```

4. The TV responds with a cookie that Haus stores for subsequent requests.

### Security Notes

- The PSK method transmits the key in plaintext over HTTP (port 80). The API is local-only.
- PIN-based pairing is more secure but requires user interaction on the TV.
- The TV must have "Remote Start" enabled for Wake-on-LAN functionality.
- IP Control must be explicitly enabled in TV settings (disabled by default on consumer models).

## API Reference

The BRAVIA API uses JSON-RPC 1.0 over HTTP POST. All endpoints are under `http://{ip}/sony/`.

### Service Endpoints

| Endpoint | Purpose |
|----------|---------|
| `/sony/system` | System info, power, LED, network |
| `/sony/avContent` | Inputs, content, playing info |
| `/sony/audio` | Volume, mute, speaker settings |
| `/sony/appControl` | App launch, app list |
| `/sony/IRCC` | IR remote control codes (SOAP/XML) |

### General Request Format

```json
POST http://{ip}/sony/{service}
Content-Type: application/json
X-Auth-PSK: {psk}

{
  "method": "methodName",
  "id": 1,
  "params": [],
  "version": "1.0"
}
```

### Power

**Get Power Status:**
```json
POST /sony/system
{"method": "getPowerStatus", "id": 1, "params": [], "version": "1.0"}
```
Response: `{"result": [{"status": "active"}], "id": 1}` (values: `active`, `standby`)

**Set Power Status:**
```json
POST /sony/system
{"method": "setPowerStatus", "id": 1, "params": [{"status": false}], "version": "1.0"}
```
Setting `status: true` turns on the TV (if Remote Start is enabled). Setting `status: false` puts it in standby.

### Volume

**Get Volume:**
```json
POST /sony/audio
{"method": "getVolumeInformation", "id": 1, "params": [], "version": "1.0"}
```
Response:
```json
{
  "result": [[
    {"target": "speaker", "volume": 25, "mute": false, "maxVolume": 100, "minVolume": 0},
    {"target": "headphone", "volume": 15, "mute": false, "maxVolume": 100, "minVolume": 0}
  ]],
  "id": 1
}
```

**Set Volume:**
```json
POST /sony/audio
{"method": "setAudioVolume", "id": 1, "params": [{"target": "speaker", "volume": "30"}], "version": "1.0"}
```

Relative volume is also supported: `"volume": "+5"` or `"volume": "-5"`.

**Set Mute:**
```json
POST /sony/audio
{"method": "setAudioMute", "id": 1, "params": [{"status": true}], "version": "1.0"}
```

### Input Selection

**Get Input List:**
```json
POST /sony/avContent
{"method": "getCurrentExternalInputsStatus", "id": 1, "params": [], "version": "1.0"}
```
Response includes array of inputs:
```json
{
  "result": [[
    {"uri": "extInput:hdmi?port=1", "title": "HDMI 1", "connection": true, "label": "PlayStation 5", "icon": "meta:hdmi"},
    {"uri": "extInput:hdmi?port=2", "title": "HDMI 2", "connection": true, "label": "", "icon": "meta:hdmi"},
    {"uri": "extInput:hdmi?port=3", "title": "HDMI 3", "connection": false, "label": "", "icon": "meta:hdmi"},
    {"uri": "extInput:hdmi?port=4", "title": "HDMI 4/ARC", "connection": false, "label": "", "icon": "meta:hdmi"}
  ]],
  "id": 1
}
```

**Set Input:**
```json
POST /sony/avContent
{"method": "setPlayContent", "id": 1, "params": [{"uri": "extInput:hdmi?port=1"}], "version": "1.0"}
```

### Content Info

**Get Playing Content Info:**
```json
POST /sony/avContent
{"method": "getPlayingContentInfo", "id": 1, "params": [], "version": "1.0"}
```
Response: `{"result": [{"uri": "extInput:hdmi?port=1", "source": "extInput:hdmi", "title": "HDMI 1"}], "id": 1}`

### App Control

**Get App List:**
```json
POST /sony/appControl
{"method": "getApplicationList", "id": 1, "params": [], "version": "1.0"}
```
Response: array of apps with `title`, `uri`, `icon`.

**Launch App:**
```json
POST /sony/appControl
{"method": "setActiveApp", "id": 1, "params": [{"uri": "com.sony.dtv.com.netflix.ninja/.MainActivity"}], "version": "1.0"}
```

Common app URIs:
- `com.sony.dtv.com.netflix.ninja/.MainActivity` -- Netflix
- `com.sony.dtv.com.google.android.youtube.tv/.MainActivity` -- YouTube
- `com.sony.dtv.com.amazon.amazonvideo.livingroom/.MainActivity` -- Prime Video
- `com.sony.dtv.com.disney.disneyplus/.MainActivity` -- Disney+

### System Information

**Get System Info:**
```json
POST /sony/system
{"method": "getSystemInformation", "id": 1, "params": [], "version": "1.0"}
```
Response:
```json
{
  "result": [{
    "product": "TV",
    "region": "US",
    "language": "en",
    "model": "XR-55A95L",
    "serial": "...",
    "macAddr": "AA:BB:CC:DD:EE:FF",
    "name": "BRAVIA",
    "generation": "2023",
    "area": "US"
  }],
  "id": 1
}
```

### IRCC (IR Remote Control Codes)

For button-press simulation, Sony uses IRCC over SOAP/XML:

```
POST http://{ip}/sony/IRCC
Content-Type: text/xml; charset=UTF-8
X-Auth-PSK: {psk}
SOAPACTION: "urn:schemas-sony-com:service:IRCC:1#X_SendIRCC"

<?xml version="1.0"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:X_SendIRCC xmlns:u="urn:schemas-sony-com:service:IRCC:1">
      <IRCCCode>{code}</IRCCCode>
    </u:X_SendIRCC>
  </s:Body>
</s:Envelope>
```

**Common IRCC Codes:**
- `AAAAAQAAAAEAAAAVAw==` -- Power toggle
- `AAAAAQAAAAEAAAASAw==` -- Volume Up
- `AAAAAQAAAAEAAAATAw==` -- Volume Down
- `AAAAAQAAAAEAAAAUAw==` -- Mute
- `AAAAAQAAAAEAAAAQAw==` -- Channel Up
- `AAAAAQAAAAEAAAARAw==` -- Channel Down
- `AAAAAQAAAAEAAAB0Aw==` -- Input
- `AAAAAgAAAJcAAAAjAw==` -- Home
- `AAAAAgAAAJcAAAA9Aw==` -- Guide
- `AAAAAQAAAAEAAAA0Aw==` -- Return/Back
- `AAAAAgAAAJcAAAANAw==` -- Play
- `AAAAAgAAAJcAAAAYAw==` -- Pause
- `AAAAAgAAAJcAAAAaAw==` -- Stop
- `AAAAAgAAAJcAAAAcAw==` -- Rewind
- `AAAAAgAAAJcAAAAdAw==` -- Fast Forward
- `AAAAAQAAAAEAAAB0Aw==` -- Netflix (direct)

**Get IRCC Code List:**
```json
POST /sony/system
{"method": "getRemoteControllerInfo", "id": 1, "params": [], "version": "1.0"}
```
Returns an array of all supported IRCC codes with their names.

## AI Capabilities

When the AI concierge "chats as" a Sony Bravia TV, it can:

- **Power on/off** the TV (true power on supported via REST API with Remote Start enabled)
- **Get and set volume** to specific numeric levels, mute/unmute
- **List HDMI inputs** with connection status and CEC labels, switch inputs
- **Report current content** -- what input or app is active
- **Launch streaming apps** by name
- **Control media playback** via IRCC codes (play, pause, stop, rewind, fast forward)
- **Query system info** -- model, generation, serial, firmware
- **Send any IRCC remote button** -- full virtual remote capability

The AI speaks in first person as the TV, aware of its Sony Bravia identity and Google TV platform.

## Quirks & Notes

- **IP Control Must Be Enabled:** The BRAVIA REST API is disabled by default on consumer models. Users must go to Settings > Network > Home Network > IP Control and enable "Simple IP Control" and optionally set a Pre-Shared Key. Without this step, port 80 will not respond to `/sony/` endpoints.
- **PSK vs PIN:** Pre-Shared Key is simpler for automation but requires the user to set it via TV menus. PIN pairing is more user-friendly for initial setup but involves a multi-step handshake. Haus should support both, preferring PSK for simplicity.
- **Google TV vs Android TV:** Older Sony TVs (2015-2020) run Android TV; newer ones (2021+) run Google TV. The BRAVIA REST API is the same on both. The difference is in the smart platform UI and app availability, not the control API.
- **Cast Protocol:** Since all Sony Google TV / Android TV models include Chromecast built-in, they also support the Cast protocol on ports 8008/8009. This is a separate control plane from the BRAVIA API and can be used for media casting.
- **AirPlay 2 / HomeKit:** 2020+ models support AirPlay 2 and HomeKit. These are independent from the BRAVIA API and handled by Apple's frameworks.
- **REST API Power On:** Unlike LG and Samsung where the API server only runs when the TV is on, Sony Bravia TVs with "Remote Start" enabled keep a minimal network stack active in standby, allowing `setPowerStatus(true)` to work. This is a significant advantage.
- **IRCC Code Encoding:** IRCC codes are base64-encoded binary values. The `getRemoteControllerInfo` method returns all supported codes, which vary by model. Always query the TV for its supported codes rather than hardcoding.
- **JSON-RPC Version:** The API uses JSON-RPC 1.0 (not 2.0). The `version` field in requests refers to the API version, not the JSON-RPC spec version. Always use `"version": "1.0"`.
- **Rate Limits:** No formal rate limit, but rapid IRCC commands (faster than ~200ms apart) may be dropped. REST API calls are more reliable under rapid use.
- **Dual API Surface:** The REST API (`/sony/*`) provides structured read/write control. IRCC provides button simulation. Use REST for state queries and direct control; use IRCC for commands not available in REST (e.g., specific menu navigation).

## Similar Devices

- **lg-oled-tv-webos** -- LG's competing OLED platform with WebSocket SSAP API
- **samsung-smart-tv-tizen** -- Samsung Tizen platform with WebSocket API
- **roku-tv-streaming-stick** -- Roku with ECP REST API
