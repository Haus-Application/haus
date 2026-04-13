---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "wiz-connected-bulb"
name: "WiZ Connected Smart Bulb"
manufacturer: "Signify"
brand: "WiZ"
model: "A19 Full Color"
model_aliases: ["WiZ A60", "WiZ A19", "WiZ BR30", "WiZ GU10", "WiZ Lightstrip", "WiZ Filament"]
device_type: "wiz_bulb"
category: "lighting"
product_line: "WiZ Connected"
release_year: 2020
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
  mac_prefixes: ["D8:0D:17", "A8:BB:50", "44:4F:8E"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [38899]
  signature_ports: [38899]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^wiz_[a-f0-9]+$", "^WiZLight_[a-f0-9]+$"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "wiz"
  polling_interval_sec: 15
  websocket_event: "wiz:state"
  setup_type: "none"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp", "scenes"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 38899
  transport: "UDP"
  encoding: "JSON"
  auth_method: "none"
  auth_detail: "No authentication required. Any device on the local network can send UDP commands."
  base_url_template: "udp://{ip}:38899"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "bulb"
  power_source: "mains"
  mounting: "ceiling"
  indoor_outdoor: "both"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.wizconnected.com/en-us/consumer/products/"
  api_docs: "https://github.com/sbidy/pywizlight"
  developer_portal: ""
  support: "https://www.wizconnected.com/en-us/support/"
  community_forum: "https://github.com/home-assistant/core/tree/dev/homeassistant/components/wiz"
  image_url: ""
  fcc_id: "2AOKM-A60FCWF"

# --- TAGS ---
tags: ["wifi_bulb", "udp", "no_auth", "signify", "hubless", "local_control"]
---

# WiZ Connected Smart Bulb

## What It Is

WiZ Connected is Signify's (the parent company of Philips Hue) budget-friendly smart lighting line. Unlike Hue, which requires a Zigbee bridge, WiZ bulbs connect directly to WiFi and are controlled via a local UDP protocol with no authentication. They are available in a wide range of form factors (A19, A21, BR30, GU10, lightstrips, filament bulbs, outdoor fixtures) and support full RGB color, tunable white, or dimmable white depending on the model. WiZ bulbs are typically priced at $10-20 per bulb, roughly one-third the cost of Philips Hue equivalents.

## How Haus Discovers It

1. **OUI match**: WiZ bulbs use WiFi chipsets with MAC prefixes registered to Signify or their chipset suppliers (D8:0D:17, A8:BB:50, 44:4F:8E). However, OUI matching alone is insufficient since Signify makes many products.
2. **Hostname pattern**: WiZ bulbs typically register DHCP hostnames matching `wiz_*` or `WiZLight_*` followed by a hex suffix.
3. **UDP registration broadcast**: WiZ bulbs periodically send UDP broadcast packets on port 38899 advertising their presence. Listening on this port is the most reliable discovery method.
4. **UDP probe**: Send a `{"method": "registration", "params": {"phoneMac": "AAAAAAAAAAAA", "register": false, "phoneIp": "{haus_ip}", "id": "1"}}` UDP packet to port 38899. WiZ bulbs respond with their module info and firmware version.
5. **getPilot probe**: Send `{"method": "getPilot", "params": {}}` to port 38899. The response includes the current light state, confirming the device is a WiZ light.

## Pairing / Authentication

WiZ bulbs require no authentication for local control. Any device on the local network can send UDP commands to port 38899 and the bulb will respond and execute them. This is the simplest integration path possible but also means there is zero security for local API access.

Initial setup of the bulb onto WiFi must be done through the WiZ app (the bulb creates an AP for provisioning), but once on the network, no pairing with Haus is needed.

## API Reference

WiZ bulbs communicate via JSON-encoded UDP datagrams on port 38899. Each request is a JSON object with `method` and `params` fields. Responses include `result` or `error`.

### Methods

| Method | Direction | Description |
|--------|-----------|-------------|
| `getPilot` | Request | Get current light state |
| `setPilot` | Request | Set light state |
| `registration` | Request | Register for push notifications |
| `pulse` | Request | Flash the bulb |
| `getSystemConfig` | Request | Get firmware, MAC, module info |
| `getModelConfig` | Request | Get hardware capabilities |
| `getWifiConfig` | Request | Get WiFi SSID/RSSI |
| `reboot` | Request | Reboot the bulb |
| `syncPilot` | Push | Bulb pushes state change (after registration) |
| `firstBeat` | Push | Bulb announces presence on network |

### getPilot Response

```json
{
  "method": "getPilot",
  "env": "pro",
  "result": {
    "mac": "d80d17xxxxxx",
    "rssi": -55,
    "src": "",
    "state": true,
    "sceneId": 0,
    "r": 255,
    "g": 128,
    "b": 0,
    "c": 0,
    "w": 0,
    "dimming": 80
  }
}
```

### setPilot Examples

**Turn on with RGB color:**
```json
{
  "method": "setPilot",
  "params": {
    "r": 255,
    "g": 0,
    "b": 128,
    "dimming": 80
  }
}
```

**Set color temperature (warm white):**
```json
{
  "method": "setPilot",
  "params": {
    "temp": 2700,
    "dimming": 100
  }
}
```

**Activate a built-in scene:**
```json
{
  "method": "setPilot",
  "params": {
    "sceneId": 4
  }
}
```

**Turn off:**
```json
{
  "method": "setPilot",
  "params": {
    "state": false
  }
}
```

### Built-in Scene IDs

| ID | Scene | ID | Scene |
|----|-------|----|-------|
| 1 | Ocean | 17 | Party |
| 2 | Romance | 18 | Fall |
| 3 | Sunset | 19 | Spring |
| 4 | Party | 20 | Summer |
| 5 | Fireplace | 21 | Deep Dive |
| 6 | Cozy | 22 | Jungle |
| 7 | Forest | 23 | Mojito |
| 8 | Pastel Colors | 24 | Club |
| 9 | Wake Up | 25 | Christmas |
| 10 | Bedtime | 26 | Halloween |
| 11 | Warm White | 27 | Candlelight |
| 12 | Daylight | 28 | Golden White |
| 13 | Cool White | 29 | Pulse |
| 14 | Night Light | 30 | Steampunk |
| 15 | Focus | 31 | Diwali |
| 16 | Relax | 32 | White |

### Registration for Push Updates

To receive real-time state changes without polling, register with the bulb:

```json
{
  "method": "registration",
  "params": {
    "phoneMac": "AABBCCDDEEFF",
    "register": true,
    "phoneIp": "192.168.1.100",
    "id": "1"
  }
}
```

After registration, the bulb sends `syncPilot` UDP packets to the registered IP whenever its state changes.

### Color Model

WiZ uses a 5-channel color model:
- **r, g, b**: RGB channels (0-255)
- **c**: Cold white LED channel (0-255)
- **w**: Warm white LED channel (0-255)
- **temp**: Color temperature in Kelvin (2200-6500)
- **dimming**: Brightness percentage (10-100)

When setting `temp`, do not include `r`, `g`, `b`, `c`, `w`. When setting RGB, do not include `temp`. The channels are mutually exclusive modes.

## AI Capabilities

Not yet planned, but when implemented the AI could set colors and scenes by name, adjust brightness, set color temperatures for different times of day, and activate the built-in dynamic scenes.

## Quirks & Notes

- **No authentication**: Any device on the LAN can control WiZ bulbs. This is a security concern on shared networks. Haus should note this in the UI.
- **UDP unreliability**: UDP is connectionless with no delivery guarantee. Haus should retry commands and verify with a `getPilot` follow-up if the command is critical.
- **Minimum dimming**: The minimum dimming value is 10, not 0. Setting `dimming: 0` is not valid; use `state: false` to turn off.
- **Registration expiry**: Push notification registrations expire and must be periodically renewed (roughly every 30 seconds to 2 minutes, send a keepalive registration).
- **Same parent as Hue**: WiZ and Philips Hue are both Signify products. Some WiZ OUI prefixes overlap with other Signify products. Hostname and UDP probe are more reliable than OUI alone.
- **WiFi only**: WiZ bulbs have no Zigbee or Thread radio. They connect directly to the home WiFi router. Large WiZ installations (20+ bulbs) can strain consumer WiFi routers.
- **Matter support**: Newer WiZ firmware versions advertise Matter-over-WiFi support, though implementation is ongoing. The UDP protocol remains the more reliable integration path.
- **5-channel vs 3-channel**: Not all WiZ models have all 5 LED channels. Basic white bulbs only have warm/cool white. `getModelConfig` reveals the actual capabilities.

## Similar Devices

- [philips-hue-bridge](philips-hue-bridge.md) — Same parent company (Signify), but hub-based with Zigbee
- [kasa-smart-bulb-kl130](kasa-smart-bulb-kl130.md) — Another hubless WiFi bulb with local API
- [cync-smart-bulb](cync-smart-bulb.md) — WiFi bulb but cloud-only (no local API)
