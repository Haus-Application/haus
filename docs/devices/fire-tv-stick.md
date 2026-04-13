---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "fire-tv-stick"
name: "Amazon Fire TV Stick"
manufacturer: "Amazon.com, Inc."
brand: "Amazon"
model: "Fire TV Stick 4K Max (2nd Gen)"
model_aliases: ["Fire TV Stick 4K", "Fire TV Stick Lite", "Fire TV Stick (3rd Gen)", "Fire TV Cube (3rd Gen)", "Fire TV Stick 4K Max", "Fire TV Omni QLED", "Fire TV Omni", "Fire TV 2-Series", "Fire TV 4-Series"]
device_type: "streaming_device"
category: "media"
product_line: "Fire TV"
release_year: 2023
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi", "ethernet", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "00:FC:8B"        # Amazon Technologies Inc.
    - "0C:47:C9"        # Amazon Technologies Inc.
    - "10:CE:A9"        # Amazon Technologies Inc.
    - "18:74:2E"        # Amazon Technologies Inc.
    - "24:4C:E3"        # Amazon Technologies Inc.
    - "34:D2:70"        # Amazon Technologies Inc.
    - "38:F7:3D"        # Amazon Technologies Inc.
    - "40:A2:DB"        # Amazon Technologies Inc.
    - "44:65:0D"        # Amazon Technologies Inc.
    - "50:DC:E7"        # Amazon Technologies Inc.
    - "5C:41:5A"        # Amazon Technologies Inc.
    - "68:37:E9"        # Amazon Technologies Inc.
    - "68:54:FD"        # Amazon Technologies Inc.
    - "6C:56:97"        # Amazon Technologies Inc.
    - "74:75:48"        # Amazon Technologies Inc.
    - "74:C2:46"        # Amazon Technologies Inc.
    - "78:E1:03"        # Amazon Technologies Inc.
    - "84:D6:D0"        # Amazon Technologies Inc.
    - "A0:02:DC"        # Amazon Technologies Inc.
    - "AC:63:BE"        # Amazon Technologies Inc.
    - "B4:7C:9C"        # Amazon Technologies Inc.
    - "B4:A5:AC"        # Amazon Technologies Inc.
    - "CC:F7:35"        # Amazon Technologies Inc.
    - "F0:27:2D"        # Amazon Technologies Inc.
    - "F0:D2:F1"        # Amazon Technologies Inc.
    - "F0:F0:A4"        # Amazon Technologies Inc.
    - "FC:65:DE"        # Amazon Technologies Inc.
    - "FC:A1:83"        # Amazon Technologies Inc.
  mdns_services:
    - "_amzn-wplay._tcp"      # Amazon Wireless Display
  mdns_txt_keys: []
  default_ports: [5555, 7000, 8008, 8443, 2870]
  signature_ports: [5555]     # ADB port (when enabled)
  ssdp_search_target: "urn:dial-multiscreen-org:service:dial:1"
  ssdp_server_string: ""
  hostname_patterns:
    - "^amazon-"
    - "^Fire.*TV"
    - "^AFTMM"
    - "^AFTN"
    - "^AFTS"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8008
    path: "/setup/eureka_info"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "detected_only"
  integration_key: "fire_tv"
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "none"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "media_playback"
  - "input_select"

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 5555
  transport: "TCP"
  encoding: "binary"
  auth_method: "app_pairing"
  auth_detail: "ADB (Android Debug Bridge) on port 5555 when developer mode is enabled. First ADB connection requires on-screen RSA key approval. Not intended for consumer automation."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "usb"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.amazon.com/fire-tv-stick-4k-max/dp/B0BXM23V6T"
  api_docs: ""
  developer_portal: "https://developer.amazon.com/apps-and-games/fire-tv"
  support: "https://www.amazon.com/gp/help/customer/display.html?nodeId=201348270"
  community_forum: "https://developer.amazon.com/support"
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["fire_tv", "alexa", "amazon", "adb", "cloud_only", "streaming", "4k", "hdr", "hdmi", "voice_remote", "matter_controller"]
---

# Amazon Fire TV Stick

## What It Is

The Amazon Fire TV Stick is a streaming media player manufactured by Amazon. The lineup includes the Fire TV Stick Lite, Fire TV Stick (3rd Gen), Fire TV Stick 4K, Fire TV Stick 4K Max, Fire TV Cube, and Fire TV branded televisions (Omni series, 2-Series, 4-Series). Fire TV devices run Fire OS, a fork of Android, and are deeply integrated with Amazon's Alexa voice assistant and Amazon Prime Video. Unlike Roku (which has an excellent local ECP API) or LG (which has SSAP), Fire TV has no official local control API. The primary control path is through the Alexa cloud ecosystem. For advanced users, Android Debug Bridge (ADB) over WiFi can be enabled for limited local control, but this requires enabling developer mode and is not a consumer-friendly integration path. Haus can detect Fire TV devices on the network but cannot meaningfully control them locally.

## How Haus Discovers It

1. **OUI Match** -- Amazon Technologies Inc. has many registered MAC prefixes (e.g., `F0:F0:A4`, `FC:65:DE`, `68:37:E9`, `44:65:0D`). However, Amazon makes many devices (Echo, Ring, Kindle, Eero, Fire tablets), so a MAC match to Amazon does not confirm a Fire TV specifically.

2. **SSDP Discovery** -- Fire TV devices respond to DIAL (Discovery and Launch) protocol searches with `urn:dial-multiscreen-org:service:dial:1`. This identifies the device as DIAL-capable but does not uniquely identify it as Fire TV (Chromecast and other devices also support DIAL).

3. **mDNS Discovery** -- Fire TV devices may advertise `_amzn-wplay._tcp.local.` for Amazon Wireless Display (screen mirroring). This is more specific to Amazon devices but does not distinguish Fire TV from Fire tablets.

4. **Hostname Pattern** -- Fire TV devices often have hostnames starting with their hardware model code: `AFTMM` (Fire TV Stick 4K Max), `AFTN` (Fire TV Stick), `AFTS` (Fire TV Stick Lite), etc. These internal codenames in the hostname are a good identifier.

5. **Port Probe** -- If ADB debugging is enabled (developer mode), port 5555 will be open. This is not a reliable discovery signal since most users do not enable developer mode. Port 8008 may respond with DIAL service information.

6. **Negative Identification** -- If a device matches an Amazon MAC prefix but does NOT have ports associated with Echo (port 55443) or Ring (port 443 with Ring-specific certificates), and DOES have DIAL on port 8008 or ADB on port 5555, it is likely a Fire TV.

## Pairing / Authentication

### ADB (Developer Mode Required)

ADB is the only local control mechanism for Fire TV. It requires manual setup:

1. On the Fire TV: **Settings > My Fire TV > Developer Options > ADB Debugging > On**
   - If "Developer Options" is not visible: Go to **Settings > My Fire TV > About** and click "Serial Number" 7 times to enable developer mode.

2. Also enable: **Settings > My Fire TV > Developer Options > Apps from Unknown Sources > On** (required for some ADB commands)

3. Connect via ADB from the Haus hub:
   ```
   adb connect {ip}:5555
   ```

4. The Fire TV displays an on-screen prompt: **"Allow USB debugging? The RSA key fingerprint is: XX:XX:XX:..."** The user must select "Allow" (optionally checking "Always allow from this computer").

5. Once approved, ADB commands can be sent to the device.

### Alexa Cloud API

For cloud-based control, the Alexa Smart Home Skill API or Alexa Voice Service can be used, but this requires:
- Amazon developer account
- OAuth2 with Login with Amazon
- Alexa skill certification
- Internet connectivity

This is a cloud-only path and not suitable for Haus's local-first approach.

## API Reference

### ADB Commands (Developer Mode Only)

ADB provides shell-level access to the Fire OS (Android) system. Useful commands for media control:

**Send Key Events (simulated remote):**
```bash
adb shell input keyevent {keycode}
```

**Common Android Key Codes for Fire TV:**

| Key Code | Key Name | Description |
|----------|----------|-------------|
| 3 | KEYCODE_HOME | Home |
| 4 | KEYCODE_BACK | Back |
| 19 | KEYCODE_DPAD_UP | Up |
| 20 | KEYCODE_DPAD_DOWN | Down |
| 21 | KEYCODE_DPAD_LEFT | Left |
| 22 | KEYCODE_DPAD_RIGHT | Right |
| 23 | KEYCODE_DPAD_CENTER | Select/OK |
| 24 | KEYCODE_VOLUME_UP | Volume Up |
| 25 | KEYCODE_VOLUME_DOWN | Volume Down |
| 26 | KEYCODE_POWER | Power/Sleep |
| 82 | KEYCODE_MENU | Menu |
| 85 | KEYCODE_MEDIA_PLAY_PAUSE | Play/Pause |
| 86 | KEYCODE_MEDIA_STOP | Stop |
| 87 | KEYCODE_MEDIA_NEXT | Next |
| 88 | KEYCODE_MEDIA_PREVIOUS | Previous |
| 89 | KEYCODE_MEDIA_REWIND | Rewind |
| 90 | KEYCODE_MEDIA_FAST_FORWARD | Fast Forward |
| 164 | KEYCODE_MUTE | Mute |

**Launch App by Package Name:**
```bash
adb shell am start -n {package}/{activity}
```

Common Fire TV app packages:
- `com.netflix.ninja/.MainActivity` -- Netflix
- `com.google.android.youtube.tv/com.google.android.apps.youtube.tv.activity.ShellActivity` -- YouTube
- `com.amazon.avod/.HomeActivity` -- Prime Video
- `com.disney.disneyplus/com.bamtechmedia.dominguez.main.MainActivity` -- Disney+
- `com.hulu.plus/.MainActivity` -- Hulu

**Get Current Activity (foreground app):**
```bash
adb shell dumpsys window windows | grep -E 'mCurrentFocus|mFocusedApp'
```

**Screen On/Off:**
```bash
adb shell input keyevent 26  # toggle screen
```

**Get Device Properties:**
```bash
adb shell getprop ro.product.model     # e.g., "AFTMM"
adb shell getprop ro.build.version.sdk  # Android SDK version
adb shell getprop ro.product.brand      # "Amazon"
```

**Install APK:**
```bash
adb install path/to/app.apk
```

**Screenshot:**
```bash
adb shell screencap -p /sdcard/screenshot.png
adb pull /sdcard/screenshot.png
```

### DIAL Protocol (Limited)

Fire TV supports DIAL (Discovery and Launch) for launching specific apps:

```
POST http://{ip}:8008/apps/{app_name}
```

DIAL app names: `Netflix`, `YouTube`, `AmazonInstantVideo`, `Disney+`

This can launch apps but cannot control playback or query state.

## AI Capabilities

Haus marks Fire TV as `detected_only`. The AI concierge can:

- **Identify the device** on the network -- "I see an Amazon Fire TV Stick 4K Max at 192.168.1.50"
- **Report network presence** -- online/offline status via ping
- **Explain limitations** -- "Fire TV doesn't have a local control API. I can see it on your network but can't control it without Alexa cloud integration or ADB developer mode."

The AI cannot control the Fire TV without ADB being manually enabled by the user. If ADB is enabled, limited control is theoretically possible but is not a supported Haus integration path.

## Quirks & Notes

- **No Local API:** This is the fundamental limitation. Amazon has not provided a local control API for Fire TV. All first-party control goes through the Alexa cloud. This is in stark contrast to Roku's completely open ECP, LG's SSAP, and Sony's BRAVIA API.
- **ADB Is a Developer Tool:** ADB was designed for app development and debugging, not consumer automation. It requires manual enablement, on-screen RSA key approval, and can be unstable. Fire OS updates occasionally reset developer mode settings.
- **ADB Security Risks:** ADB over WiFi has no encryption. Any device on the network can connect if the RSA key was previously approved. ADB provides shell access to the entire OS, which is far more access than needed for media control.
- **Amazon MAC Sprawl:** Amazon Technologies has a massive number of OUI registrations covering Echo speakers, Ring cameras, Eero routers, Kindle readers, Fire tablets, Fire TV, and even some AWS hardware. MAC-based identification is only the first step; additional probing is essential to distinguish device type.
- **Matter Controller:** Fire TV Stick 4K Max (2nd Gen, 2023) and Fire TV Cube (3rd Gen) include a Thread border router and can serve as Matter controllers. This is potentially useful for Haus's Matter integration but does not help with controlling the Fire TV itself.
- **Fire TV Cube Special Case:** The Fire TV Cube (3rd Gen) has built-in IR blaster and HDMI-CEC control, allowing it to control other devices (TV power, volume, input). It also has a built-in speaker for Alexa without a TV. The Cube is more of a hub than a simple streaming stick.
- **Alexa Routines:** Users can create Alexa routines that include Fire TV actions (launch apps, play content). These are cloud-only and not accessible via local API.
- **Fire TV Omni TVs:** Amazon's own branded TVs (Omni series) have additional features like ambient mode and hands-free Alexa. They still lack a local control API beyond ADB.
- **HDMI-CEC:** Fire TV Sticks support HDMI-CEC for basic TV control (power on/off, volume). This is a one-way control from Fire TV to TV, not useful for controlling the Fire TV itself.
- **Wireless Display:** Fire TV supports screen mirroring via Miracast (`_amzn-wplay._tcp`). This is a display protocol, not a control protocol.

## Similar Devices

- **roku-tv-streaming-stick** -- Roku with excellent open local ECP API (recommended alternative for local control)
- **lg-oled-tv-webos** -- LG TVs with local WebSocket SSAP API
- **samsung-smart-tv-tizen** -- Samsung TVs with local WebSocket API
- **sony-bravia-google-tv** -- Sony TVs with local REST API
