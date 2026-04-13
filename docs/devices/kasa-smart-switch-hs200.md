---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "kasa-smart-switch-hs200"
name: "Kasa Smart Wi-Fi Light Switch"
manufacturer: "TP-Link"
brand: "Kasa Smart"
model: "HS200"
model_aliases: ["HS200(US)"]
device_type: "kasa_switch"
category: "smart_home"
product_line: "Kasa"
release_year: 2016
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
  hostname_patterns: ["^HS200.*"]
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
  form_factor: "switch"
  power_source: "hardwired"
  mounting: "in_wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.kasasmart.com/us/products/smart-switches/kasa-smart-wi-fi-light-switch-hs200"
  api_docs: ""
  developer_portal: ""
  support: "https://www.kasasmart.com/us/support"
  community_forum: ""
  image_url: ""
  fcc_id: "TE7HS200"

# --- TAGS ---
tags: ["switch", "in-wall", "wifi", "no-hub", "local-control", "xor-protocol"]
---

# Kasa Smart Wi-Fi Light Switch (HS200)

## What It Is

> The TP-Link Kasa HS200 is a standard single-pole in-wall Wi-Fi light switch that replaces a traditional toggle or paddle switch. It provides simple on/off control of hardwired lighting loads with no hub required. Setup is done through the Kasa app (cloud-provisioned Wi-Fi credentials), but once configured, the device operates entirely over the local network using TP-Link's proprietary XOR-encrypted TCP protocol on port 9999. It requires a neutral wire and works with most standard light fixtures.

## How Haus Discovers It

1. **OUI match** -- MAC address begins with a known TP-Link prefix (50:C7:BF, B0:BE:76, 60:A4:B7, 1C:3B:F3, 5C:A6:E6, 98:DA:C4, etc.)
2. **Port scan** -- TCP connect scan on port 9999 across the local subnet
3. **Protocol probe** -- Send XOR-encrypted `{"system":{"get_sysinfo":{}}}` and validate decrypted JSON response
4. **Model identification** -- Parse `model` field from `get_sysinfo` response (e.g., `"HS200(US)"`)
5. **Type classification** -- Model prefix `HS200` maps to `kasa_switch` device type with `on_off` capability

## Pairing / Authentication

> No pairing or authentication is required for local control. The XOR encryption on port 9999 uses a static key (`0xAB`) and provides obfuscation, not security. Any device on the local network can control Kasa switches.
>
> **Initial Wi-Fi setup** requires the Kasa mobile app and a TP-Link cloud account to provision the device onto the Wi-Fi network. After setup, cloud connectivity is optional -- all control works locally.

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
| `alias` | string | Device display name (e.g., "Kitchen Lights") |
| `model` | string | Hardware model (e.g., "HS200(US)") |
| `relay_state` | int | 0 = off, 1 = on |
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

## AI Capabilities

> When chatting with an HS200 switch, the AI can:
> - **Query real-time state** via XOR protocol -- on/off status, uptime
> - **Toggle power** on/off instantly via `set_relay_state`
> - **Optimistic UI** -- controls update instantly in the UI, poller confirms within 10 seconds
>
> The AI speaks as the device: "I'm currently on. I've been running for 3 hours and 42 minutes."

## Quirks & Notes

- **Neutral wire required** -- the HS200 requires a neutral wire in the switch box, which older homes may lack
- **No dimming** -- this is a relay-only switch; for dimming, use the HS220
- **3-way switching** -- for 3-way setups, use the HS210 instead, or wire the HS200 as the primary switch with a compatible companion
- **Cloud setup only** -- initial provisioning requires the Kasa app and TP-Link cloud account, but after setup the device works fully local
- **XOR key is universal** -- all Kasa devices use the same static key `0xAB` (171); this is obfuscation, not security
- **1-second timeout** -- TCP connections should use a 1-second connect timeout; unresponsive devices are marked offline
- **Firmware updates** -- delivered via Kasa cloud; may occasionally change response fields but the core XOR protocol has been stable since 2016

## Similar Devices

- **[kasa-smart-dimmer-hs220](kasa-smart-dimmer-hs220.md)** -- same form factor with dimming and fan speed support
- **[kasa-smart-plug-hs103](kasa-smart-plug-hs103.md)** -- plug-in variant with the same on/off protocol
- **[kasa-smart-power-strip-kp303](kasa-smart-power-strip-kp303.md)** -- multi-outlet strip with per-outlet control
