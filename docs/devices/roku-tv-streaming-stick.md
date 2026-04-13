---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "roku-tv-streaming-stick"
name: "Roku TV / Streaming Stick"
manufacturer: "Roku, Inc."
brand: "Roku"
model: "Roku Ultra (4802RW)"
model_aliases: ["3820RW", "3810RW", "3811RW", "3920RW", "3930RW", "3940RW", "3941RW", "3942RW", "4800RW", "4801RW", "4802RW", "C11LE-9554RW", "7000X"]
device_type: "streaming_device"
category: "media"
product_line: "Roku"
release_year: 2023
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["wifi", "ethernet", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "00:0D:4B"        # Roku, Inc.
    - "08:05:81"        # Roku, Inc.
    - "10:59:32"        # Roku, Inc.
    - "20:EF:BD"        # Roku, Inc.
    - "28:3B:82"        # Roku, Inc.
    - "2C:E4:10"        # Roku, Inc.
    - "34:1E:58"        # Roku, Inc.
    - "3C:59:1E"        # Roku, Inc.
    - "58:FD:20"        # Roku, Inc.
    - "64:B0:A6"        # Roku, Inc.
    - "84:EA:ED"        # Roku, Inc.
    - "88:DE:A9"        # Roku, Inc.
    - "A8:62:A2"        # Roku, Inc.
    - "AC:3A:7A"        # Roku, Inc.
    - "B0:A7:37"        # Roku, Inc.
    - "B0:EE:7B"        # Roku, Inc.
    - "B8:3E:59"        # Roku, Inc.
    - "BC:D7:D4"        # Roku, Inc.
    - "C8:3A:6B"        # Roku, Inc.
    - "CC:6D:A0"        # Roku, Inc.
    - "D0:4D:C6"        # Roku, Inc.
    - "D4:E2:2F"        # Roku, Inc.
    - "D8:31:34"        # Roku, Inc.
    - "DC:3A:5E"        # Roku, Inc.
  mdns_services:
    - "_roku._tcp"            # primary Roku mDNS service
  mdns_txt_keys:
    - "name"                  # device friendly name
    - "uuid"                  # device unique identifier
  default_ports: [8060, 9080]
  signature_ports: [8060]     # ECP HTTP API -- definitive Roku signal
  ssdp_search_target: "roku:ecp"
  ssdp_server_string: "Roku/9.4 UPnP/1.0 Roku/9.4"
  hostname_patterns:
    - "^Roku"
    - "^roku-"
    - "^RokuPlayer"
    - "^RokuStreamingStick"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8060
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "Roku"
    body_contains: ""
    headers: {}
  - port: 8060
    path: "/query/device-info"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "Roku"
    body_contains: "<device-info>"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "roku"
  polling_interval_sec: 5
  websocket_event: "roku:state"
  setup_type: "none"
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
  port: 8060
  transport: "HTTP"
  encoding: "XML"
  auth_method: "none"
  auth_detail: "No authentication required. ECP is open by default on port 8060 to any device on the local network."
  base_url_template: "http://{ip}:8060"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.roku.com/products/players"
  api_docs: "https://developer.roku.com/docs/developer-program/debugging/external-control-api.md"
  developer_portal: "https://developer.roku.com/"
  support: "https://www.roku.com/en-us/support"
  community_forum: "https://community.roku.com/"
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["roku", "ecp", "rest_api", "no_auth", "local_api", "streaming", "4k", "hdr", "hdmi", "voice_remote", "roku_os"]
---

# Roku TV / Streaming Stick

## What It Is

Roku devices -- including the Roku Streaming Stick, Roku Express, Roku Ultra, Roku Streambar, and Roku-branded TVs from TCL, Hisense, and others -- run Roku OS and are manufactured or licensed by Roku, Inc. They are the most popular streaming platform in the United States. Critically for Haus, Roku exposes an outstanding local HTTP REST API called the External Control Protocol (ECP) on port 8060. ECP requires no authentication whatsoever -- any device on the local network can query device info, list installed apps, launch apps, send remote keypresses, and control media playback. This makes Roku one of the easiest and most reliable devices to integrate locally. ECP has been stable since Roku OS 7.x and works identically across all Roku hardware (sticks, boxes, TVs, soundbars).

## How Haus Discovers It

1. **OUI Match** -- Roku, Inc. has many registered MAC prefixes (e.g., `B0:A7:37`, `D8:31:34`, `AC:3A:7A`, `CC:6D:A0`). A match to a Roku OUI is a strong indicator since Roku only makes streaming devices. For Roku-branded TVs (TCL Roku TV, Hisense Roku TV), the MAC may belong to TCL or Hisense instead.

2. **SSDP Discovery** -- Roku devices respond to UPnP M-SEARCH with the custom search target `roku:ecp`. This is the fastest and most reliable discovery method. The SSDP response includes the ECP URL directly:
   ```
   ST: roku:ecp
   Location: http://192.168.1.100:8060/
   USN: uuid:roku:ecp:XXXXXXXXXXXX
   ```

3. **mDNS Discovery** -- Roku devices advertise `_roku._tcp.local.` with TXT records containing `name` (friendly device name) and `uuid` (unique identifier).

4. **HTTP Fingerprint** -- `GET http://{ip}:8060/query/device-info` returns an XML document with comprehensive device information. The response `Server` header contains "Roku". This endpoint requires no authentication and is the definitive identification.

5. **Port Probe** -- Port 8060 open and responding with a `Server: Roku` header is nearly 100% conclusive.

## Pairing / Authentication

No pairing or authentication is required. ECP is completely open.

The ECP API on port 8060 accepts requests from any device on the local network without any form of authentication -- no API keys, no tokens, no pairing prompts, no passwords. This is by design and has been Roku's approach since the protocol was introduced.

### Security Implications

- Any device on the same LAN can control the Roku.
- Roku has a "Device Access" setting in newer firmware that can restrict ECP to specific companion apps, but this is off by default and does not affect most ECP endpoints.
- There is no TLS option; all communication is plaintext HTTP.

## API Reference

The External Control Protocol (ECP) is a RESTful HTTP API on port 8060. Responses are in XML format. Commands use a mix of GET (queries) and POST (actions).

**Base URL:** `http://{ip}:8060`

### Device Info

```
GET /query/device-info
```

**Response (XML):**
```xml
<device-info>
  <udn>xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx</udn>
  <serial-number>XXXXXXXXXXXX</serial-number>
  <device-id>XXXXXXXXXXXX</device-id>
  <advertising-id>xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx</advertising-id>
  <vendor-name>Roku</vendor-name>
  <model-name>Roku Ultra</model-name>
  <model-number>4802X</model-number>
  <model-region>US</model-region>
  <is-tv>false</is-tv>
  <is-stick>false</is-stick>
  <supports-ethernet>true</supports-ethernet>
  <wifi-mac>AA:BB:CC:DD:EE:FF</wifi-mac>
  <wifi-driver>realtek</wifi-driver>
  <ethernet-mac>AA:BB:CC:DD:EE:00</ethernet-mac>
  <network-type>wifi</network-type>
  <friendly-device-name>Living Room Roku</friendly-device-name>
  <friendly-model-name>Roku Ultra</friendly-model-name>
  <default-device-name>Roku Ultra</default-device-name>
  <user-device-name>Living Room Roku</user-device-name>
  <user-device-location>Living Room</user-device-location>
  <build-number>AAH.00E04209A</build-number>
  <software-version>12.5.0</software-version>
  <software-build>4209</software-build>
  <secure-device>true</secure-device>
  <language>en</language>
  <country>US</country>
  <locale>en_US</locale>
  <time-zone-auto>true</time-zone-auto>
  <time-zone>US/Central</time-zone>
  <time-zone-name>United States/Central</time-zone-name>
  <time-zone-tz>America/Chicago</time-zone-tz>
  <time-zone-offset>-360</time-zone-offset>
  <clock-format>12-hour</clock-format>
  <uptime>345600</uptime>
  <power-mode>PowerOn</power-mode>
  <supports-suspend>true</supports-suspend>
  <supports-find-remote>true</supports-find-remote>
  <find-remote-is-possible>true</find-remote-is-possible>
  <supports-audio-guide>true</supports-audio-guide>
  <supports-rva>true</supports-rva>
  <developer-enabled>false</developer-enabled>
  <keyed-developer-id/>
  <search-enabled>true</search-enabled>
  <search-channels-enabled>true</search-channels-enabled>
  <voice-search-enabled>true</voice-search-enabled>
  <notifications-enabled>true</notifications-enabled>
  <notifications-first-use>true</notifications-first-use>
  <supports-private-listening>true</supports-private-listening>
  <headphones-connected>false</headphones-connected>
  <supports-ecs-textedit>true</supports-ecs-textedit>
  <supports-ecs-microphone>true</supports-ecs-microphone>
  <supports-wake-on-wlan>true</supports-wake-on-wlan>
  <supports-airplay>false</supports-airplay>
  <has-play-on-roku>true</has-play-on-roku>
  <has-mobile-screensaver>false</has-mobile-screensaver>
  <support-url>roku.com/support</support-url>
  <grandcentral-version>6.0.70</grandcentral-version>
  <trc-version>3.0</trc-version>
  <trc-channel-version>4.2</trc-channel-version>
  <davinci-version>2.8.20</davinci-version>
</device-info>
```

Key fields:
- `power-mode` -- `PowerOn` or `Standby` (indicates current power state)
- `is-tv` -- `true` if this is a Roku TV (vs a stick/box)
- `model-name` -- hardware model
- `software-version` -- Roku OS version
- `user-device-name` -- user-configured friendly name
- `supports-wake-on-wlan` -- whether WOL is available

### App List (Installed Channels)

```
GET /query/apps
```

**Response:**
```xml
<apps>
  <app id="12" subtype="ndka" type="appl" version="14.40.14">Netflix</app>
  <app id="837" subtype="ndka" type="appl" version="4.1.1803">YouTube</app>
  <app id="13" subtype="ndka" type="appl" version="16.20.1">Amazon Prime Video</app>
  <app id="291097" subtype="ndka" type="appl" version="2.17.61">Disney+</app>
  <app id="2285" subtype="ndka" type="appl" version="7.30.0">Hulu</app>
  <app id="551012" subtype="ndka" type="appl" version="1.13.49">Apple TV</app>
  <app id="195316" subtype="ndka" type="appl" version="8.12.3">HBO Max</app>
  <app id="46041" subtype="ndka" type="appl" version="6.0.2">Sling TV</app>
  <app id="61322" subtype="ndka" type="appl" version="5.14.0">Tubi</app>
  <app id="593099" subtype="ndka" type="appl" version="2.3.4">Peacock</app>
  <app id="13535" subtype="ndka" type="appl" version="7.8.3">Plex</app>
</apps>
```

### Active App

```
GET /query/active-app
```

**Response:**
```xml
<active-app>
  <app id="12" subtype="ndka" type="appl" version="14.40.14">Netflix</app>
  <screensaver id="55545" type="ssvr" version="2.0.1">Default screensaver</screensaver>
</active-app>
```

If the user is on the Roku home screen: `<app id="" type="" version="">Roku</app>`

### Launch App

```
POST /launch/{app_id}
```

No request body required. Launches the app by its numeric ID.

**Common App IDs:**
| App ID | Name |
|--------|------|
| 12 | Netflix |
| 837 | YouTube |
| 13 | Amazon Prime Video |
| 291097 | Disney+ |
| 2285 | Hulu |
| 551012 | Apple TV+ |
| 195316 | HBO Max / Max |
| 593099 | Peacock |
| 13535 | Plex |
| 46041 | Sling TV |
| 61322 | Tubi |
| 34376 | ESPN |
| 151908 | The Roku Channel |
| 27536 | CBS / Paramount+ |
| 74519 | Spotify |

**Deep Link into App:**
```
POST /launch/{app_id}?contentId={content_id}&mediaType=movie
```

Query parameters allow deep-linking directly to content within an app (e.g., a specific movie or show).

### Install App

```
POST /install/{app_id}
```

Opens the Roku Channel Store page for the specified app and prompts installation.

### Keypress (Remote Control)

**Single Press:**
```
POST /keypress/{key}
```

**Key Down (hold):**
```
POST /keydown/{key}
```

**Key Up (release):**
```
POST /keyup/{key}
```

No request body. The key name is in the URL path.

**Available Keys:**

| Key | Description |
|-----|-------------|
| `Home` | Home button |
| `Rev` | Rewind |
| `Fwd` | Fast Forward |
| `Play` | Play/Pause toggle |
| `Select` | OK/Select |
| `Left` | Left arrow |
| `Right` | Right arrow |
| `Down` | Down arrow |
| `Up` | Up arrow |
| `Back` | Back button |
| `InstantReplay` | Instant Replay (10 sec back) |
| `Info` | Options/Info (*) |
| `Backspace` | Backspace (text entry) |
| `Search` | Search |
| `Enter` | Enter (text entry) |
| `FindRemote` | Find Remote (beeps the remote) |
| `VolumeDown` | Volume Down |
| `VolumeUp` | Volume Up |
| `VolumeMute` | Volume Mute toggle |
| `ChannelUp` | Channel Up (Roku TV) |
| `ChannelDown` | Channel Down (Roku TV) |
| `InputTuner` | TV tuner input (Roku TV) |
| `InputHDMI1` | HDMI 1 (Roku TV) |
| `InputHDMI2` | HDMI 2 (Roku TV) |
| `InputHDMI3` | HDMI 3 (Roku TV) |
| `InputHDMI4` | HDMI 4 (Roku TV) |
| `InputAV1` | AV 1 (Roku TV) |
| `PowerOff` | Power Off (Roku TV) |
| `PowerOn` | Power On (Roku TV, if WOL supported) |

**Text Entry (character by character):**
```
POST /keypress/Lit_{character}
```
Example: To type "hello":
```
POST /keypress/Lit_h
POST /keypress/Lit_e
POST /keypress/Lit_l
POST /keypress/Lit_l
POST /keypress/Lit_o
```

URL-encode special characters: space = `Lit_%20`, @ = `Lit_%40`.

### Input (Touch/Accelerometer)

```
POST /input?touch.0.x={x}&touch.0.y={y}&touch.0.op=press
```

Sends touch/accelerometer events for games and apps that support it.

### Media Player Status

```
GET /query/media-player
```

**Response:**
```xml
<player error="false" state="play">
  <plugin bandwidth="20000000 bps" id="12" name="Netflix"/>
  <format audio="aac_adts" captions="none" container="cenc" drm="widevine" video="h265"/>
  <buffering current="1000" max="1000" target="0"/>
  <new_stream speed="128"/>
  <position milliseconds="345678" runtime="7200000"/>
  <duration milliseconds="7200000"/>
  <is_live>false</is_live>
</player>
```

Key fields:
- `state` -- `play`, `pause`, `stop`, `buffer`, `close`
- `position` -- current playback position in milliseconds
- `duration` -- total content duration in milliseconds
- `format` -- video/audio codec info
- `plugin` -- which app is playing

### TV Input List (Roku TV Only)

```
GET /query/tv-channels
```

Returns the list of broadcast TV channels (for Roku TVs with tuners).

```
GET /query/tv-active-channel
```

Returns the currently tuned channel.

### Search

```
POST /search/browse?keyword={search_term}&type=movie&season=1&show-unavailable=true
```

Searches across all installed channels for content matching the keyword.

### Device Image / App Icon

```
GET /query/icon/{app_id}
```

Returns the app icon as a PNG image. Useful for UI display.

## AI Capabilities

When the AI concierge "chats as" a Roku device, it can:

- **Power on/off** the TV (PowerOn/PowerOff keys for Roku TVs, WOL for sticks/boxes)
- **List all installed apps** with names and IDs
- **Launch any app** by name ("open Netflix", "switch to YouTube")
- **Report what's playing** -- current app, playback state, position, duration
- **Control media playback** -- play, pause, rewind, fast forward, instant replay
- **Adjust volume** -- up, down, mute (Roku TV or Roku Streambar)
- **Switch inputs** on Roku TVs (HDMI 1-4, tuner, AV)
- **Navigate** -- send directional keys, select, back, home
- **Enter text** into search fields character by character
- **Search for content** across all installed channels
- **Find the remote** -- trigger the remote finder beep
- **Report device info** -- model, software version, uptime, power state

The AI speaks in first person as the Roku, aware of its streaming-first identity and installed channel lineup.

## Quirks & Notes

- **No Authentication:** ECP is completely open. This is both a strength (zero-config for Haus) and a potential concern (any LAN device can control the Roku). Roku has added optional "Device Access" settings in newer firmware but these are off by default.
- **XML Responses:** Unlike most modern APIs, ECP returns XML, not JSON. Haus must parse XML for all query responses. POST commands (keypress, launch) return empty 200 responses on success.
- **Roku TV vs Streaming Device:** Roku TVs (`is-tv: true` in device-info) support additional keys like `PowerOff`, `PowerOn`, `ChannelUp`, `ChannelDown`, `InputHDMI1`-`InputHDMI4`, `InputTuner`, and `InputAV1`. Streaming sticks/boxes do not have these keys. Volume keys work on Roku TVs and Roku Streambar but not on basic Roku sticks.
- **TCL/Hisense Roku TVs:** These third-party Roku TVs use TCL or Hisense MAC prefixes, not Roku OUIs. Discovery should rely on SSDP `roku:ecp` or mDNS `_roku._tcp` rather than MAC matching for these devices.
- **Power State:** Roku streaming sticks/boxes do not truly power off; they enter a screensaver or standby. The `power-mode` field in device-info distinguishes `PowerOn` from `Standby`. Roku TVs can be fully powered off.
- **App Launch Timing:** After sending `POST /launch/{id}`, the app may take several seconds to load. Sending keypress commands too soon after launch will be lost. Poll `/query/active-app` to confirm the app is running before sending further commands.
- **Play/Pause Toggle:** The `Play` key toggles between play and pause. There are no separate "play" and "pause" keys. Check `/query/media-player` state to know current playback state before toggling.
- **ECP v2 (Newer Firmware):** Roku OS 11+ introduced ECP v2 with additional endpoints under `/api/v2/`. These include structured input list queries and enhanced device information. The v1 endpoints remain available and stable.
- **Private Listening:** Roku supports "private listening" (audio routed to phone/headphones). The `headphones-connected` field in device-info indicates if this is active.
- **Rate Limits:** No formal rate limit, but keypress commands sent faster than ~50ms apart may be dropped or buffered. For text entry, add a small delay between characters.
- **Screensaver:** When the Roku is idle, a screensaver activates. The active-app query shows the screensaver ID. Any keypress or app launch command will dismiss the screensaver.

## Similar Devices

- **fire-tv-stick** -- Amazon's competing streaming platform with limited local control (ADB)
- **lg-oled-tv-webos** -- LG TVs with WebSocket API
- **samsung-smart-tv-tizen** -- Samsung TVs with WebSocket API
- **sony-bravia-google-tv** -- Sony TVs with REST API
