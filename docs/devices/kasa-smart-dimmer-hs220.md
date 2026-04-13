---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "kasa-smart-dimmer-hs220"
name: "Kasa Smart Wi-Fi Dimmer Switch"
manufacturer: "TP-Link"
brand: "Kasa Smart"
model: "HS220"
model_aliases: ["HS220(US)", "KS220", "KS220(US)"]
device_type: "kasa_dimmer"
category: "smart_home"
product_line: "Kasa"
release_year: 2017
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["50:C7:BF", "B0:BE:76", "60:A4:B7", "1C:3B:F3", "5C:A6:E6", "98:DA:C4", "B0:4E:26", "A8:42:A1", "68:FF:7B", "30:DE:4B", "E8:48:B8", "AC:15:A2"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [9999]
  signature_ports: [9999]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^HS220.*", "^KS220.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "kasa"
  polling_interval_sec: 10
  websocket_event: "kasa:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M4"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "fan_speed"]

# --- PROTOCOL ---
protocol:
  type: "tcp_xor"
  port: 9999
  transport: "TCP"
  encoding: "XOR-JSON"
  auth_method: "none"
  auth_detail: "No authentication. XOR encryption with static key 0xAB (171) provides obfuscation only."
  base_url_template: "tcp://{ip}:9999"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "switch"
  power_source: "hardwired"
  mounting: "in_wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.kasasmart.com/us/products/smart-switches/kasa-smart-wi-fi-dimmer-switch-hs220"
  api_docs: ""
  developer_portal: ""
  support: "https://www.kasasmart.com/us/support"
  community_forum: ""
  image_url: ""
  fcc_id: "TE7HS220"

# --- TAGS ---
tags: ["dimmer", "switch", "in-wall", "wifi", "no-hub", "local-control", "xor-protocol", "fan-capable"]
---

# Kasa Smart Wi-Fi Dimmer Switch (HS220)

## What It Is

> The TP-Link Kasa HS220 is an in-wall Wi-Fi dimmer switch that replaces a standard single-pole light switch. It provides on/off control and smooth brightness dimming (0-100%) for compatible dimmable lighting loads. No hub is required. The KS220 is the newer variant with identical protocol behavior. When the device alias contains "fan", Haus treats it as a fan controller with 4-speed tier mapping. Like all Kasa devices, it communicates over the local network using TP-Link's XOR-encrypted TCP protocol on port 9999.

## How Haus Discovers It

1. **OUI match** -- MAC address begins with a known TP-Link prefix (50:C7:BF, B0:BE:76, 60:A4:B7, 1C:3B:F3, 5C:A6:E6, 98:DA:C4, etc.)
2. **Port scan** -- TCP connect scan on port 9999 across the local subnet
3. **Protocol probe** -- Send XOR-encrypted `{"system":{"get_sysinfo":{}}}` and validate decrypted JSON response
4. **Model identification** -- Parse `model` field from `get_sysinfo` response (e.g., `"HS220(US)"` or `"KS220(US)"`)
5. **Type classification** -- Model prefix `HS220` or `KS220` maps to `kasa_dimmer` device type
6. **Fan detection** -- If `alias` field contains the substring "fan" (case-insensitive), capabilities include `fan_speed` and the device is treated as a fan controller

## Pairing / Authentication

> No pairing or authentication is required for local control. The XOR encryption on port 9999 uses a static key (`0xAB`) and provides obfuscation, not security. Any device on the local network can control Kasa dimmers.
>
> **Initial Wi-Fi setup** requires the Kasa mobile app and a TP-Link cloud account to provision the device onto the Wi-Fi network. After setup, cloud connectivity is optional.

## API Reference

### Protocol: TCP XOR Encryption

All communication uses TCP on port 9999 with XOR-encrypted JSON payloads. Messages are framed with a 4-byte big-endian length prefix followed by the encrypted payload.

**XOR Encryption Algorithm:**

```
Encrypt: key = 171; for each byte in plaintext: encrypted_byte = key XOR byte; key = encrypted_byte; output encrypted_byte
Decrypt: key = 171; for each byte in ciphertext: decrypted_byte = byte XOR key; key = byte; output decrypted_byte
```

**Connection flow:**
1. TCP connect to `{device_ip}:9999` (1-second timeout)
2. Encrypt JSON command using XOR algorithm
3. Prepend 4-byte big-endian uint32 payload length
4. Send length prefix + encrypted payload
5. Read 4-byte response length prefix
6. Read that many bytes of encrypted response
7. Decrypt response using XOR algorithm
8. Parse resulting JSON

### Query Device State

**Request:**
```json
{"system":{"get_sysinfo":{}}}
```

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| `alias` | string | Device display name (e.g., "Living Room Lights") |
| `model` | string | Hardware model (e.g., "HS220(US)") |
| `relay_state` | int | 0 = off, 1 = on |
| `brightness` | int | 0-100, current brightness percentage (dimmers only) |
| `dev_name` | string | Device description |
| `mac` | string | MAC address |
| `deviceId` | string | Unique device identifier |
| `on_time` | int | Seconds since last power on |

### Turn On

**Request:**
```json
{"system":{"set_relay_state":{"state":1}}}
```

### Turn Off

**Request:**
```json
{"system":{"set_relay_state":{"state":0}}}
```

### Set Brightness

**Request:**
```json
{"smartlife.iot.dimmer":{"set_brightness":{"brightness":75}}}
```

- `brightness` -- integer 0-100
- Only works on HS220/KS220 dimmers
- Device must be on (relay_state=1) for brightness changes to take visible effect
- Setting brightness while off will store the value; it applies when the device is turned on

### Fan Speed Control

Fan devices are detected by the `alias` field containing "fan" (case-insensitive). Fan speed is mapped to brightness tiers:

| Speed | Brightness Value | Label |
|-------|-----------------|-------|
| 1 | 25 | Low |
| 2 | 50 | Medium |
| 3 | 75 | High |
| 4 | 100 | Max |

Use the `set_brightness` command with the mapped brightness value. For example, to set a fan to medium speed:

```json
{"smartlife.iot.dimmer":{"set_brightness":{"brightness":50}}}
```

## AI Capabilities

> When chatting with an HS220 dimmer, the AI can:
> - **Query real-time state** via XOR protocol -- on/off, brightness level, uptime
> - **Toggle power** on/off instantly via `set_relay_state`
> - **Set brightness** to any value 0-100% via `set_brightness`
> - **Set fan speed** 1-4 when alias indicates a fan (maps to brightness tiers)
> - **Optimistic UI** -- controls update instantly in the UI, poller confirms within 10 seconds
>
> The AI speaks as the device: "I'm currently off, brightness set to 89%."

## Quirks & Notes

- **Neutral wire required** -- the HS220 requires a neutral wire in the switch box
- **Fan detection is alias-based** -- Haus detects fans by checking if the device `alias` contains "fan"; rename the device alias in the Kasa app to trigger fan mode
- **Minimum brightness** -- some dimmable LED loads may flicker below 10-15%; the HS220 has a built-in minimum brightness setting configurable via the Kasa app
- **Incandescent vs LED** -- rated for 300W LED / 150W CFL / 600W incandescent loads
- **KS220 is protocol-identical** -- the KS220 model is a newer hardware revision of the HS220 with identical firmware protocol behavior
- **Brightness while off** -- setting brightness when relay_state=0 stores the value; it takes effect on next power-on
- **XOR key is universal** -- all Kasa devices use the same static key `0xAB` (171)
- **Firmware updates** -- delivered via Kasa cloud; the `set_brightness` namespace (`smartlife.iot.dimmer`) has been stable across all known firmware versions

## Similar Devices

- **[kasa-smart-switch-hs200](kasa-smart-switch-hs200.md)** -- same form factor, on/off only (no dimming)
- **[kasa-smart-plug-hs103](kasa-smart-plug-hs103.md)** -- plug-in variant with on/off only
- **[kasa-smart-power-strip-kp303](kasa-smart-power-strip-kp303.md)** -- multi-outlet strip with per-outlet control
