---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "yamaha-rx-v-receiver"
name: "Yamaha RX-V AV Receiver"
manufacturer: "Yamaha Corporation"
brand: "Yamaha"
model: "RX-V6A"
model_aliases: ["RX-V4A", "RX-A2A", "RX-A4A", "RX-A6A", "RX-A8A"]
device_type: "av_receiver"
category: "media"
product_line: "RX-V / Aventage"
release_year: 2020
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "ethernet", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["00:A0:DE", "04:09:86", "10:C6:FC", "28:87:BA", "34:54:3C", "40:B0:76", "58:C4:C4", "A8:5B:78", "B0:99:28"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [80]
  signature_ports: [80]
  ssdp_search_target: "urn:schemas-upnp-org:device:MediaRenderer:1"
  ssdp_server_string: ""
  hostname_patterns: ["^RX-V.*", "^RX-A.*", "^yamaha.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/YamahaExtendedControl/v1/system/getDeviceInfo"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "response_code"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "yamaha"
  polling_interval_sec: 10
  websocket_event: "yamaha:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: ["on_off", "volume", "input_select", "media_playback"]

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 80
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "none"
  auth_detail: "No authentication required. All endpoints are open on the local network."
  base_url_template: "http://{ip}/YamahaExtendedControl/v1"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "receiver"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://usa.yamaha.com/products/audio_visual/av_receivers_amps/"
  api_docs: "https://github.com/rsc-dev/pyamaha"
  developer_portal: ""
  support: "https://usa.yamaha.com/support/"
  community_forum: ""
  image_url: ""
  fcc_id: "A3L-RXV6A"

# --- TAGS ---
tags: ["av-receiver", "musiccast", "multi-zone", "surround-sound", "hdmi", "airplay", "bluetooth", "local-control", "no-auth"]
---

# Yamaha RX-V AV Receiver

## What It Is

> Yamaha RX-V and Aventage (RX-A) series AV receivers are multi-channel surround sound receivers with network streaming capabilities via Yamaha's MusicCast platform. They expose a local HTTP REST API on port 80 called "Extended Control" (also known as the MusicCast API) that provides full control over power, volume, input selection, mute, sound programs, and multi-zone audio -- all via simple GET requests with no authentication required. Models with MusicCast / Network Module are supported.

## How Haus Discovers It

1. **SSDP** -- Yamaha receivers advertise via UPnP as `urn:schemas-upnp-org:device:MediaRenderer:1`
2. **HTTP fingerprint** -- `GET /YamahaExtendedControl/v1/system/getDeviceInfo` on port 80 returns JSON with `response_code: 0`
3. **Model identification** -- Parse `model_name` from `getDeviceInfo` response (e.g., `"RX-V6A"`)
4. **OUI match** -- MAC address begins with a known Yamaha prefix (00:A0:DE, 04:09:86, etc.)

## Pairing / Authentication

> No pairing or authentication is required. All MusicCast Extended Control API endpoints are open to any device on the local network. All commands are HTTP GET requests.

## API Reference

### Base URL

```
http://{receiver_ip}/YamahaExtendedControl/v1
```

### System Info

```
GET /system/getDeviceInfo
```

**Response:**
```json
{
  "response_code": 0,
  "model_name": "RX-V6A",
  "destination": "U",
  "device_id": "04098643797E",
  "system_version": 1.70,
  "api_version": 2.15,
  "serial_number": "Y726554RT"
}
```

### Get Main Zone Status

```
GET /main/getStatus
```

**Response:**
```json
{
  "response_code": 0,
  "power": "standby",
  "volume": 45,
  "mute": false,
  "input": "hdmi1",
  "sound_program": "straight"
}
```

- `power`: `"on"` or `"standby"`
- `volume`: 0-161 (0 = -80dB, 161 = +16.5dB)
- `input`: `"hdmi1"`-`"hdmi7"`, `"av1"`-`"av7"`, `"audio1"`-`"audio4"`, `"tuner"`, `"bluetooth"`, `"airplay"`, `"spotify"`, `"net_radio"`
- `sound_program`: active DSP program name

### Power Control

```
GET /main/setPower?power=on
GET /main/setPower?power=standby
GET /main/setPower?power=toggle
```

### Volume Control

```
GET /main/setVolume?volume=50
GET /main/setVolume?volume=up&step=1
GET /main/setVolume?volume=down&step=1
```

Volume range is 0-161 (maps to -80dB to +16.5dB in 0.5dB steps).

### Mute

```
GET /main/setMute?enable=true
GET /main/setMute?enable=false
```

### Input Selection

```
GET /main/setInput?input=hdmi1
GET /main/setInput?input=bluetooth
GET /main/setInput?input=airplay
```

### Sound Program

```
GET /main/setSoundProgram?program=straight
GET /main/setSoundProgram?program=surr_decoder
GET /main/getSoundProgramList
```

### Multi-Zone Control (Zone 2, Zone 3)

Replace `main` with `zone2` or `zone3`:

```
GET /zone2/getStatus
GET /zone2/setPower?power=on
GET /zone2/setVolume?volume=30
GET /zone2/setInput?input=audio1
```

Not all models support Zone 3. Zone 2 typically has limited input options (no HDMI passthrough).

### Network Standby

```
GET /system/getNetworkStandby
GET /system/setNetworkStandby?standby=on
```

Network standby must be enabled for the receiver to accept commands while in standby mode.

### Response Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Initializing |
| 2 | Internal Error |
| 3 | Invalid Request |
| 4 | Invalid Parameter |
| 5 | Guarded (operation blocked) |

## AI Capabilities

> When chatting with a Yamaha receiver, the AI can:
> - **Query current state** -- power, volume, input, mute, sound program
> - **Get device info** -- model name, firmware version, serial number
> - **Toggle power** on/off/standby
> - **Set volume** to absolute value or relative up/down
> - **Change input** -- HDMI, Bluetooth, AirPlay, etc.
> - **Set sound program** -- Straight, Surround Decoder, etc.
> - **Multi-zone control** -- independent power, volume, input per zone
>
> All via the MusicCast HTTP REST API on port 80. No authentication required.
> The AI speaks as the device: "I'm in standby mode. Volume is at 45, input is HDMI 1."

## Quirks & Notes

- **All GET requests** -- every command in the Extended Control API is a GET request with query string parameters; there are no POST/PUT endpoints
- **Volume scale** -- volume 0-161 maps to -80dB to +16.5dB in 0.5dB steps; typical listening volume is 40-80
- **Network standby** -- must be enabled (`setNetworkStandby?standby=on`) for power-on commands to work when the receiver is in standby; this setting persists across power cycles
- **No auth** -- the API is completely open on the local network; anyone on the LAN can control the receiver
- **Response code 5 (Guarded)** -- returned when an operation is temporarily blocked (e.g., changing input during initialization); retry after a brief delay
- **Multi-zone limitations** -- Zone 2 typically supports analog inputs only (no HDMI); Zone 3 availability varies by model
- **HDMI CEC** -- the receiver also supports HDMI CEC for power/input sync with connected displays, but this is outside the network API

## Similar Devices

> Other MusicCast-enabled Yamaha products (soundbars, wireless speakers) use the same Extended Control API but with a reduced feature set. Denon/Marantz HEOS receivers have a similar HTTP-based API but use different endpoints.
