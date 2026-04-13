---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "kasa-smart-bulb-kl130"
name: "TP-Link Kasa Smart Bulb KL130"
manufacturer: "TP-Link"
brand: "Kasa Smart"
model: "KL130"
model_aliases: ["KL110", "KL120", "KL125", "KL135", "KL50", "KL60", "LB130", "LB120", "LB110"]
device_type: "kasa_bulb"
category: "lighting"
product_line: "Kasa"
release_year: 2017
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["50:C7:BF", "B0:A7:B9", "54:AF:97", "68:FF:7B", "98:DA:C4", "1C:3B:F3", "60:32:B1", "B4:B0:24"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [9999]
  signature_ports: [9999]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^KL\\d{2,3}.*$", "^LB\\d{2,3}.*$"]
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
  haus_milestone: "M3"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp"]

# --- PROTOCOL ---
protocol:
  type: "tcp_xor"
  port: 9999
  transport: "TCP"
  encoding: "XOR-JSON"
  auth_method: "none"
  auth_detail: "No authentication. Commands are XOR-encrypted JSON over TCP on port 9999. The XOR key is 0xAB (171) with autokey cipher."
  base_url_template: "tcp://{ip}:9999"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "bulb"
  power_source: "mains"
  mounting: "ceiling"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.kasasmart.com/us/products/smart-lighting"
  api_docs: "https://github.com/python-kasa/python-kasa"
  developer_portal: ""
  support: "https://www.kasasmart.com/us/support"
  community_forum: "https://github.com/home-assistant/core/tree/dev/homeassistant/components/tplink"
  image_url: ""
  fcc_id: "TE7KL130"

# --- TAGS ---
tags: ["wifi_bulb", "xor_protocol", "no_auth", "hubless", "local_control", "kasa", "tp-link", "color"]
---

# TP-Link Kasa Smart Bulb KL130

## What It Is

The TP-Link Kasa KL130 is a WiFi smart bulb with full RGB color and tunable white support. It uses the exact same XOR-encoded JSON over TCP protocol on port 9999 that all Kasa smart switches, dimmers, and plugs use, making it part of the existing Haus Kasa integration with zero additional protocol work. The Kasa bulb lineup spans several models: KL50/KL60 (filament/soft glow, dimmable only), KL110 (dimmable white), KL120/KL125 (tunable white), and KL130/KL135 (full color). The older "LB" prefix models (LB110, LB120, LB130) are predecessors that speak the same protocol. Kasa bulbs are priced at $10-25 and are widely available.

## How Haus Discovers It

Kasa bulbs are discovered identically to Kasa switches and plugs — they are part of the same integration:

1. **OUI match**: TP-Link MAC prefixes (50:C7:BF, B0:A7:B9, 54:AF:97, 68:FF:7B, 98:DA:C4, 1C:3B:F3, 60:32:B1, B4:B0:24).
2. **Port probe**: TCP port 9999 is the signature port for all Kasa devices.
3. **XOR probe**: Send an XOR-encrypted `{"system":{"get_sysinfo":{}}}` to port 9999. The response includes `model` (e.g., "KL130(US)"), `type` or `mic_type` (e.g., "IOT.SMARTBULB"), and device-specific fields.
4. **Device type detection**: The `get_sysinfo` response for bulbs includes `is_dimmable`, `is_color`, `is_variable_color_temp` fields that identify exact capabilities. The `type` field will be `IOT.SMARTBULB` for bulbs vs `IOT.SMARTPLUGSWITCH` for switches/plugs.
5. **Hostname pattern**: Kasa bulbs register DHCP hostnames matching their model number (e.g., `KL130_*`, `LB130_*`).

Since Haus already has Kasa integration for switches, bulbs are discovered and controlled through the same code path with type-specific command handling.

## Pairing / Authentication

No authentication is required. Same as all Kasa devices — any device on the local network can send XOR-encrypted commands to port 9999.

Initial WiFi provisioning must be done through the Kasa app (the bulb creates a setup AP), but once on the network, no pairing is needed.

## API Reference

Kasa bulbs use the same XOR-over-TCP protocol as all Kasa devices. The encryption is a simple autokey XOR cipher starting with key byte `0xAB` (171). The JSON command structure differs slightly from switches because bulbs use the `smartlife.iot.smartbulb.lightingservice` namespace.

### XOR Encryption

```go
func encrypt(plaintext []byte) []byte {
    key := byte(0xAB)
    result := make([]byte, len(plaintext)+4)
    binary.BigEndian.PutUint32(result[:4], uint32(len(plaintext)))
    for i, b := range plaintext {
        result[i+4] = b ^ key
        key = result[i+4]
    }
    return result
}
```

### Key Commands

**Get system info (identifies device type and capabilities):**
```json
{"system": {"get_sysinfo": {}}}
```

Response includes:
```json
{
  "system": {
    "get_sysinfo": {
      "model": "KL130(US)",
      "type": "IOT.SMARTBULB",
      "alias": "Living Room Lamp",
      "mic_type": "IOT.SMARTBULB",
      "is_dimmable": 1,
      "is_color": 1,
      "is_variable_color_temp": 1,
      "light_state": {
        "on_off": 1,
        "mode": "normal",
        "hue": 240,
        "saturation": 100,
        "color_temp": 0,
        "brightness": 80
      }
    }
  }
}
```

**Get light state:**
```json
{
  "smartlife.iot.smartbulb.lightingservice": {
    "get_light_state": {}
  }
}
```

**Set light state (color via HSB):**
```json
{
  "smartlife.iot.smartbulb.lightingservice": {
    "transition_light_state": {
      "on_off": 1,
      "hue": 240,
      "saturation": 100,
      "brightness": 80,
      "color_temp": 0,
      "transition_period": 1000,
      "mode": "normal"
    }
  }
}
```

**Set color temperature (warm white):**
```json
{
  "smartlife.iot.smartbulb.lightingservice": {
    "transition_light_state": {
      "on_off": 1,
      "color_temp": 2700,
      "brightness": 100,
      "transition_period": 500
    }
  }
}
```

**Turn off:**
```json
{
  "smartlife.iot.smartbulb.lightingservice": {
    "transition_light_state": {
      "on_off": 0
    }
  }
}
```

**Get preferred light states (saved presets):**
```json
{
  "smartlife.iot.smartbulb.lightingservice": {
    "get_preferred_state": {}
  }
}
```

### Color Model

Kasa bulbs use HSB (Hue/Saturation/Brightness):
- **hue**: 0-360 (degrees on the color wheel)
- **saturation**: 0-100 (percentage)
- **brightness**: 0-100 (percentage)
- **color_temp**: 2500-9000 Kelvin (when in white mode; set to 0 for color mode)
- **transition_period**: Transition time in milliseconds

When `color_temp` is non-zero, the bulb is in white mode and `hue`/`saturation` are ignored. Set `color_temp: 0` to use HSB color mode.

### Model Capabilities Matrix

| Model | Dimmable | Color Temp | Full Color | Notes |
|-------|----------|------------|------------|-------|
| KL50 | Yes | No | No | Filament style |
| KL60 | Yes | No | No | Filament style |
| KL110 | Yes | No | No | Basic white |
| KL120 | Yes | Yes (2700-5000K) | No | Tunable white |
| KL125 | Yes | Yes (2500-6500K) | Yes | Color, newer |
| KL130 | Yes | Yes (2500-9000K) | Yes | Full color |
| KL135 | Yes | Yes (2500-6500K) | Yes | Color, newer |
| LB110 | Yes | No | No | Predecessor to KL110 |
| LB120 | Yes | Yes | No | Predecessor to KL120 |
| LB130 | Yes | Yes | Yes | Predecessor to KL130 |

## AI Capabilities

Since Kasa is a supported integration, the AI concierge can control Kasa bulbs:

- "Turn on the living room lamp" (on_off)
- "Set the bedroom light to 50%" (brightness)
- "Make the office light warm white" (color_temp: 2700)
- "Set the accent light to blue" (hue: 240, saturation: 100)
- "Dim the hallway light over 5 seconds" (brightness + transition_period)

The AI can also query state: "Is the kitchen light on?", "What color is the desk lamp?"

## Quirks & Notes

- **Same protocol as switches**: Kasa bulbs use the identical XOR-over-TCP protocol on port 9999 as Kasa switches and plugs. The only difference is the JSON command namespace (`smartlife.iot.smartbulb.lightingservice` vs `system` for switches).
- **No authentication**: Like all Kasa devices, there is no auth. Anyone on the LAN can control the bulb.
- **Transition support**: Kasa bulbs support smooth transitions via `transition_period` (milliseconds). This is useful for sunrise/sunset automations.
- **KLAP protocol on newer firmware**: Some newer Kasa devices have moved to the KLAP protocol (encrypted handshake + AES). As of early 2025, most bulbs still support the classic XOR protocol, but firmware updates may change this. Haus should implement KLAP as a fallback.
- **Tapo rebrand**: TP-Link is migrating the Kasa brand to "Tapo" in some markets. Tapo devices use a different protocol (KLAP or HTTPS) and are NOT compatible with the XOR protocol. Tapo bulbs (L530, L510) are different products.
- **UDP discovery**: Kasa devices can also be discovered via UDP broadcast to port 9999 with the encrypted `get_sysinfo` command. This is faster than TCP port scanning for large networks.
- **Power-on state**: Kasa bulbs can be configured to restore their previous state or turn on at a specific brightness after a power outage. This is set via `smartlife.iot.common.default_softON_behavior`.
- **Max brightness ceiling**: Some KL130 units have a firmware bug where brightness 100 is slightly dimmer than brightness 99 due to a PWM rollover. Setting brightness to 99 is a common workaround.

## Similar Devices

- [wiz-connected-bulb](wiz-connected-bulb.md) — Another hubless WiFi bulb with local API (UDP instead of TCP)
- [cync-smart-bulb](cync-smart-bulb.md) — WiFi bulb but cloud-only (no local API)
- [sengled-smart-bulb-zigbee](sengled-smart-bulb-zigbee.md) — Zigbee bulb, requires hub
