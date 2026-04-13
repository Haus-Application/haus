---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ikea-fyrtur-smart-blinds"
name: "IKEA FYRTUR Smart Blinds"
manufacturer: "IKEA"
brand: "IKEA"
model: "FYRTUR"
model_aliases: ["KADRILJ", "PRAKTLYSING", "TREDANSEN", "E1757", "E1926"]
device_type: "smart_shade"
category: "smart_home"
product_line: "IKEA Smart Home"
release_year: 2019
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: false
  cloud_api: false
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["zigbee"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "ikea"
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "Zigbee"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Zigbee device that pairs with IKEA DIRIGERA hub or TRADFRI gateway. The DIRIGERA hub exposes a local REST API over HTTPS (port 8443). Can also pair directly with any Zigbee 3.0 coordinator (e.g., Zigbee2MQTT, deCONZ)."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "shade"
  power_source: "battery"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee"]

# --- LINKS ---
links:
  product_page: "https://www.ikea.com/us/en/p/fyrtur-blackout-roller-blind-smart-wireless-battery-operated-gray-90417462/"
  api_docs: ""
  developer_portal: ""
  support: "https://www.ikea.com/us/en/customer-service/"
  community_forum: ""
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["shades", "blinds", "window-covering", "zigbee", "battery", "ikea", "dirigera", "tradfri", "zigbee2mqtt", "affordable", "solar-panel-optional"]
---

# IKEA FYRTUR Smart Blinds

## What It Is

> The IKEA FYRTUR is a battery-powered motorized blackout roller blind with built-in Zigbee connectivity. It is one of the most affordable smart blinds available, with prices starting around $100-$180 depending on size. The FYRTUR offers full position control (open, close, and any intermediate position), is powered by an internal rechargeable battery (charged via USB-C, with an optional IKEA solar charging panel), and communicates over Zigbee. It can be controlled via the IKEA DIRIGERA hub (or older TRADFRI gateway), or paired directly with any Zigbee 3.0 coordinator like Zigbee2MQTT or deCONZ. KADRILJ is the sheer/light-filtering variant; PRAKTLYSING and TREDANSEN are newer generation models with similar Zigbee functionality.

## How Haus Discovers It

> FYRTUR blinds are Zigbee devices and are not directly discoverable on the IP network. Haus discovers them indirectly:
>
> 1. **Discover the DIRIGERA hub** -- The IKEA DIRIGERA hub is discovered via mDNS (`_ihsp._tcp`) on the local network
> 2. **Query the hub** -- Once authenticated with the DIRIGERA hub, Haus queries for all paired devices
> 3. **Blind identification** -- Blinds appear with device type `blinds` and model identifiers `E1757` (FYRTUR) or `E1926` (KADRILJ)
> 4. **Alternative: Zigbee coordinator** -- If using Zigbee2MQTT, the blind appears as a Zigbee device with model ID `E1757` and manufacturer `IKEA of Sweden`

## Pairing / Authentication

### Pairing to DIRIGERA Hub

1. Open the IKEA Home Smart app
2. Navigate to devices and tap "Add device"
3. Reset the blind by pressing the pairing button on the motor module (4 quick presses of the reset button)
4. The blind LED blinks, and the hub discovers it via Zigbee
5. Assign to a room and configure

### Pairing to Zigbee2MQTT

1. Enable pairing mode on the Zigbee coordinator
2. Press the FYRTUR reset button 4 times quickly
3. The blind joins the Zigbee network and is exposed via MQTT

### DIRIGERA Hub to Haus

See [IKEA DIRIGERA Hub](ikea-dirigera-hub.md) for the hub authentication flow (OAuth2-based local API on port 8443).

## API Reference

### Via DIRIGERA Hub (HTTPS REST)

Base URL: `https://{hub_ip}:8443/v1`

**Get blind status:**
```
GET /devices/{device_id}
Authorization: Bearer {token}
```

**Response:**
```json
{
  "id": "abcd-1234-efgh",
  "type": "blinds",
  "deviceType": "blinds",
  "attributes": {
    "blindsCurrentLevel": 75,
    "blindsTargetLevel": 75,
    "blindsState": "stopped",
    "batteryPercentage": 85,
    "customName": "Bedroom Blind",
    "model": "FYRTUR blackout roller blind",
    "manufacturer": "IKEA of Sweden",
    "firmwareVersion": "24.4.5"
  }
}
```

**Set blind position:**
```
PATCH /devices/{device_id}
Authorization: Bearer {token}
Content-Type: application/json

[{
  "attributes": {
    "blindsTargetLevel": 50
  }
}]
```

**Level values:**
- `0` = fully closed (blind down)
- `100` = fully open (blind up)

### Via Zigbee2MQTT

MQTT topic: `zigbee2mqtt/{friendly_name}`

**Set position:**
```json
{"position": 50}
```

**Commands:**
```json
{"state": "OPEN"}
{"state": "CLOSE"}
{"state": "STOP"}
```

**State payload:**
```json
{
  "position": 75,
  "battery": 85,
  "linkquality": 120,
  "update_available": false
}
```

### Zigbee Cluster Details

| Cluster | ID | Description |
|---------|----|-------------|
| Window Covering | `0x0102` | Position control (lift percentage) |
| Power Configuration | `0x0001` | Battery voltage and percentage |
| Groups | `0x0004` | Group membership |
| OTA Upgrade | `0x0019` | Firmware update |

**Attributes (Window Covering cluster):**

| Attribute | ID | Description |
|-----------|----|-------------|
| Current Position Lift Percentage | `0x0008` | Current position (0-100) |
| Installed Open Limit | `0x0010` | Open limit |
| Installed Closed Limit | `0x0011` | Closed limit |

## AI Capabilities

> AI integration planned via DIRIGERA hub or Zigbee coordinator. When available:
> - Open/close blinds by room name
> - Set specific position percentages
> - Report battery level
> - Coordinate with lighting scenes

## Quirks & Notes

- **Battery life** -- Internal rechargeable lithium-ion battery lasts approximately 6 months to 1 year depending on usage; charge via USB-C cable or optional IKEA solar charging panel
- **USB-C charging** -- The charging port is on the motor module at one end of the headrail; the blind must be removed from the bracket to charge (unless using the solar panel)
- **Solar panel option** -- IKEA offers a solar charging panel (SOLSTRALE) that attaches to the window and keeps the battery topped up
- **Position inversion** -- Some integrations report inverted position (0=open, 100=closed); check the `windowCoveringType` attribute for the orientation convention
- **Slow movement** -- FYRTUR blinds move relatively slowly compared to premium shades; full travel takes approximately 20-30 seconds depending on blind length
- **Limited sizes** -- Available in fixed sizes only (not custom-cut); sizes range from 23" to 48" wide
- **Zigbee signal** -- As a battery-powered Zigbee device, FYRTUR acts as an end device (not a router), meaning it does not relay messages for other Zigbee devices
- **Group control** -- Supports Zigbee groups for synchronized multi-blind control
- **TRADFRI vs DIRIGERA** -- The older TRADFRI gateway uses CoAP protocol; the newer DIRIGERA hub uses HTTPS REST. Both work with FYRTUR, but DIRIGERA is recommended for new setups
- **Firmware updates** -- OTA updates are delivered through the IKEA hub; Zigbee2MQTT can also apply OTA updates

## Similar Devices

> - [Lutron Serena Shades](lutron-serena-shades.md) -- Premium alternative via Lutron bridge and LEAP protocol
> - [Hunter Douglas PowerView](hunter-douglas-powerview.md) -- Premium alternative with local REST API
> - [IKEA DIRIGERA Hub](ikea-dirigera-hub.md) -- Required hub for IKEA smart home ecosystem
