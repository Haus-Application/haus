---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "kasa-smart-power-strip-kp303"
name: "Kasa Smart Wi-Fi Power Strip"
manufacturer: "TP-Link"
brand: "Kasa Smart"
model: "KP303"
model_aliases: ["KP303(US)", "KP400", "KP400(US)", "KP303P5"]
device_type: "kasa_power_strip"
category: "smart_home"
product_line: "Kasa"
release_year: 2019
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
  hostname_patterns: ["^KP303.*", "^KP400.*"]
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
  form_factor: "strip"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.kasasmart.com/us/products/smart-plugs/kasa-smart-wi-fi-power-strip-kp303"
  api_docs: ""
  developer_portal: ""
  support: "https://www.kasasmart.com/us/support"
  community_forum: ""
  image_url: ""
  fcc_id: "TE7KP303"

# --- TAGS ---
tags: ["power-strip", "multi-outlet", "wifi", "no-hub", "local-control", "xor-protocol", "child-devices"]
---

# Kasa Smart Wi-Fi Power Strip (KP303)

## What It Is

> The TP-Link Kasa KP303 is a Wi-Fi smart power strip with 3 individually controllable outlets (plus 2 always-on USB-A ports for charging). Each outlet can be turned on or off independently, making it ideal for entertainment centers or desk setups where you need per-device power control. The KP400 is the outdoor-rated variant with 2 outlets and weather-resistant housing. Both use the same XOR-encrypted TCP protocol on port 9999 with child device addressing to control individual outlets.

## How Haus Discovers It

1. **OUI match** -- MAC address begins with a known TP-Link prefix (50:C7:BF, B0:BE:76, 60:A4:B7, 1C:3B:F3, etc.)
2. **Port scan** -- TCP connect scan on port 9999 across the local subnet
3. **Protocol probe** -- Send XOR-encrypted `{"system":{"get_sysinfo":{}}}` and validate decrypted JSON response
4. **Model identification** -- Parse `model` field from `get_sysinfo` response (e.g., `"KP303(US)"`)
5. **Child detection** -- Presence of `children` array in `get_sysinfo` response indicates a multi-outlet device
6. **Type classification** -- Model prefix `KP303` or `KP400` maps to `kasa_power_strip` device type; each child outlet gets its own virtual device in Haus

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

### Query Device State (with Children)

**Request:**
```json
{"system":{"get_sysinfo":{}}}
```

**Response (abbreviated):**
```json
{
  "system": {
    "get_sysinfo": {
      "model": "KP303(US)",
      "alias": "Entertainment Center",
      "mac": "B0:BE:76:XX:XX:XX",
      "deviceId": "80067...",
      "children": [
        {
          "id": "80067...00",
          "state": 1,
          "alias": "TV"
        },
        {
          "id": "80067...01",
          "state": 0,
          "alias": "Sound Bar"
        },
        {
          "id": "80067...02",
          "state": 1,
          "alias": "Game Console"
        }
      ]
    }
  }
}
```

Each child in the `children` array represents one physical outlet:

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Child device ID (parent deviceId + 2-digit suffix) |
| `state` | int | 0 = off, 1 = on |
| `alias` | string | Outlet display name |

### Control Individual Outlet

To control a specific outlet, wrap the command in a `context` block with the child ID:

**Turn on outlet 2 (index 01):**
```json
{
  "context": {
    "child_ids": ["80067...01"]
  },
  "system": {
    "set_relay_state": {"state": 1}
  }
}
```

**Turn off outlet 2:**
```json
{
  "context": {
    "child_ids": ["80067...01"]
  },
  "system": {
    "set_relay_state": {"state": 0}
  }
}
```

The `child_ids` array can contain multiple IDs to control several outlets simultaneously.

### Control All Outlets

To turn all outlets on or off, omit the `context` block:

**Turn all on:**
```json
{"system":{"set_relay_state":{"state":1}}}
```

**Turn all off:**
```json
{"system":{"set_relay_state":{"state":0}}}
```

## AI Capabilities

> When chatting with a Kasa power strip, the AI can:
> - **Query state of all outlets** -- each outlet's on/off status and alias
> - **Toggle individual outlets** by child ID
> - **Toggle all outlets** at once
> - **Optimistic UI** -- controls update instantly in the UI, poller confirms within 10 seconds
>
> The AI speaks as the device: "I have 3 outlets. TV is on, Sound Bar is off, Game Console is on."

## Quirks & Notes

- **Child device addressing** -- each outlet has a unique `child_id` derived from the parent `deviceId` plus a 2-digit suffix (00, 01, 02); commands must include the `context.child_ids` array to target specific outlets
- **USB ports are not smart** -- the 2 USB-A ports on the KP303 are always powered and cannot be toggled via the protocol
- **KP400 is outdoor** -- the KP400 variant has 2 outlets (not 3) in a weather-resistant IP64 housing; protocol is identical
- **Single IP, multiple outlets** -- the power strip has one IP address and one TCP connection for all outlets; Haus creates virtual sub-devices for each outlet in the device list
- **15A total** -- the KP303 is rated for 15 amps total across all outlets, not per-outlet
- **XOR key is universal** -- all Kasa devices use the same static key `0xAB` (171)
- **Polling returns all children** -- a single `get_sysinfo` call returns the state of all outlets, so polling is efficient

## Similar Devices

- **[kasa-smart-plug-hs103](kasa-smart-plug-hs103.md)** -- single-outlet plug variant
- **[kasa-smart-switch-hs200](kasa-smart-switch-hs200.md)** -- in-wall switch variant
- **[kasa-smart-dimmer-hs220](kasa-smart-dimmer-hs220.md)** -- in-wall dimmer variant
