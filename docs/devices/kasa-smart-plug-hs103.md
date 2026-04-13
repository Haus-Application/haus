---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "kasa-smart-plug-hs103"
name: "Kasa Smart Wi-Fi Plug Mini"
manufacturer: "TP-Link"
brand: "Kasa Smart"
model: "HS103"
model_aliases: ["HS103(US)", "HS105", "HS105(US)", "EP10", "EP10(US)"]
device_type: "kasa_plug"
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
  hostname_patterns: ["^HS103.*", "^HS105.*", "^EP10.*"]
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
capabilities: ["on_off"]

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
  form_factor: "plug"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.kasasmart.com/us/products/smart-plugs/kasa-smart-wifi-plug-mini"
  api_docs: ""
  developer_portal: ""
  support: "https://www.kasasmart.com/us/support"
  community_forum: ""
  image_url: ""
  fcc_id: "TE7HS103"

# --- TAGS ---
tags: ["plug", "mini", "wifi", "no-hub", "local-control", "xor-protocol"]
---

# Kasa Smart Wi-Fi Plug Mini (HS103)

## What It Is

> The TP-Link Kasa HS103 is a compact Wi-Fi smart plug that plugs into a standard outlet and provides on/off control for any connected load. Its "mini" form factor is designed to not block the adjacent outlet. No hub is required. Variants include the HS105 (slightly different hardware revision) and EP10 (newer generation with identical protocol). For energy monitoring, the HS110 and EP25 models add real-time power measurement. Like all Kasa devices, it communicates locally using the XOR-encrypted TCP protocol on port 9999.

## How Haus Discovers It

1. **OUI match** -- MAC address begins with a known TP-Link prefix (50:C7:BF, B0:BE:76, 60:A4:B7, 1C:3B:F3, etc.)
2. **Port scan** -- TCP connect scan on port 9999 across the local subnet
3. **Protocol probe** -- Send XOR-encrypted `{"system":{"get_sysinfo":{}}}` and validate decrypted JSON response
4. **Model identification** -- Parse `model` field from `get_sysinfo` response (e.g., `"HS103(US)"`, `"HS105(US)"`, `"EP10(US)"`)
5. **Type classification** -- Model prefix `HS103`, `HS105`, or `EP10` maps to `kasa_plug` device type with `on_off` capability

## Pairing / Authentication

> No pairing or authentication is required for local control. The XOR encryption on port 9999 uses a static key (`0xAB`) and provides obfuscation, not security.
>
> **Initial Wi-Fi setup** requires the Kasa mobile app and a TP-Link cloud account. After setup, cloud connectivity is optional.

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
| `alias` | string | Device display name (e.g., "Desk Lamp") |
| `model` | string | Hardware model (e.g., "HS103(US)") |
| `relay_state` | int | 0 = off, 1 = on |
| `dev_name` | string | Device description |
| `mac` | string | MAC address |
| `deviceId` | string | Unique device identifier |
| `on_time` | int | Seconds since last power on |
| `feature` | string | Feature flags (e.g., "TIM" for timer support) |

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

### Energy Monitoring (HS110/EP25 only)

The HS110 and EP25 variants include energy monitoring. These commands return errors on the HS103/HS105/EP10.

**Get real-time power:**
```json
{"emeter":{"get_realtime":{}}}
```

**Response:**
```json
{
  "emeter": {
    "get_realtime": {
      "current_ma": 432,
      "voltage_mv": 121345,
      "power_mw": 52340,
      "total_wh": 12450,
      "err_code": 0
    }
  }
}
```

## AI Capabilities

> When chatting with a Kasa smart plug, the AI can:
> - **Query real-time state** via XOR protocol -- on/off status, uptime
> - **Toggle power** on/off instantly via `set_relay_state`
> - **Read energy data** on HS110/EP25 models -- current, voltage, power, lifetime usage
> - **Optimistic UI** -- controls update instantly in the UI, poller confirms within 10 seconds
>
> The AI speaks as the device: "I'm currently on. The desk lamp has been running for 45 minutes."

## Quirks & Notes

- **Mini form factor** -- designed not to block the adjacent outlet, but some outlet configurations may still have clearance issues
- **No dimming** -- the HS103/HS105/EP10 are relay-only; they cannot dim connected loads
- **Energy monitoring** -- only available on HS110 (older, bulkier) and EP25 (newer mini form); the HS103/HS105/EP10 do not support the `emeter` namespace
- **EP10 is the successor** -- the EP10 is TP-Link's newer mini plug with identical protocol; Kasa and Tapo product lines are converging
- **15A max load** -- rated for 15 amps / 1800 watts; do not exceed with resistive or inductive loads
- **Away mode** -- the Kasa app supports "Away Mode" that randomly toggles the plug; this is cloud-only and does not affect local protocol behavior
- **XOR key is universal** -- all Kasa devices use the same static key `0xAB` (171)

## Similar Devices

- **[kasa-smart-switch-hs200](kasa-smart-switch-hs200.md)** -- in-wall switch variant
- **[kasa-smart-dimmer-hs220](kasa-smart-dimmer-hs220.md)** -- in-wall dimmer variant
- **[kasa-smart-power-strip-kp303](kasa-smart-power-strip-kp303.md)** -- multi-outlet strip with per-outlet control
