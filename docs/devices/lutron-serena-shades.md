---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "lutron-serena-shades"
name: "Lutron Serena Smart Shades"
manufacturer: "Lutron Electronics"
brand: "Lutron"
model: "Serena"
model_aliases: ["Serena Shades", "Serena Smart Roller Shade", "Serena Honeycomb Shade", "Serena Wood Blinds"]
device_type: "smart_shade"
category: "smart_home"
product_line: "Serena"
release_year: 2014
discontinued: false
price_range: "$$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["proprietary_rf"]

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
  integration_key: "lutron"
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
  transport: "proprietary_rf"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Serena shades do not connect to IP networks directly. They communicate via Lutron's proprietary Clear Connect RF protocol to a Lutron bridge (Caseta Smart Bridge or RA2 Select). The bridge exposes the LEAP protocol on the local network. See lutron-caseta-bridge.md for the LEAP protocol details."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "shade"
  power_source: "battery"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: []

# --- LINKS ---
links:
  product_page: "https://www.lutron.com/en-US/Products/Pages/SingleSiteBlindsShades/SerenaShadesbyLutron/Overview.aspx"
  api_docs: ""
  developer_portal: "https://developer.lutron.com"
  support: "https://www.lutron.com/en-US/Support/Pages/default.aspx"
  community_forum: ""
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["shades", "blinds", "window-covering", "battery", "clear-connect-rf", "lutron", "caseta", "leap-protocol", "position-control", "no-direct-ip"]
---

# Lutron Serena Smart Shades

## What It Is

> Lutron Serena Smart Shades are premium motorized window coverings available in roller shade, honeycomb shade, and wood blind styles. They operate on battery power (6 D-cell batteries or a rechargeable battery wand, depending on the model) and communicate wirelessly using Lutron's proprietary Clear Connect RF protocol. Serena shades are not directly IP-addressable -- they require a Lutron bridge (Caseta Smart Bridge Pro, RA2 Select, or RadioRA 3) to be controlled from a phone, voice assistant, or smart home system. The bridge handles the RF-to-IP translation and exposes the shades through the LEAP protocol on the local network.

## How Haus Discovers It

> Serena shades are not discoverable on the IP network because they communicate via RF only. Haus discovers them indirectly:
>
> 1. **Discover the Lutron bridge** -- The Caseta Smart Bridge is discovered via mDNS (`_leap._tcp`) on the local network
> 2. **Query the bridge** -- Once paired with the bridge via the LEAP protocol, Haus queries the bridge for all connected devices
> 3. **Shade identification** -- Shades appear as devices with zone type `WindowCovering` in the bridge's device list, with a `Category` of `SheerBlind`, `RollerShade`, or `HoneycombShade`

## Pairing / Authentication

### Shade to Bridge Pairing

Serena shades are paired to the Lutron bridge using the Lutron mobile app:

1. Open the Lutron app and select "Add Device"
2. Select the shade type
3. Press and hold the shade's programming button (on the back of the headrail) until the LED blinks
4. The bridge discovers the shade via Clear Connect RF
5. Assign the shade to a room and configure positioning

### Bridge to Haus Pairing

See [Lutron Caseta Bridge](lutron-caseta-bridge.md) for the LEAP protocol authentication flow. The Caseta Smart Bridge Pro (L-BDGPRO2-WH) is required for third-party integration -- the standard Caseta bridge does not support the LEAP protocol.

## API Reference

> All API access is through the Lutron bridge. The shade is controlled as a zone on the bridge.

### LEAP Protocol (via Bridge)

**Set shade position:**
```json
{
  "CommuniqueType": "CreateRequest",
  "Header": {
    "MessageBodyType": "OneZoneStatus",
    "Url": "/zone/{zone_id}/commandprocessor"
  },
  "Body": {
    "CommandType": "GoToLevel",
    "ZoneStatus": {
      "Zone": {"href": "/zone/{zone_id}"},
      "Level": 75
    }
  }
}
```

**Level values:**
- `0` = fully closed
- `100` = fully open
- Intermediate values for partial positions

**Preset positions:**
```json
{
  "CommandType": "GoToFavoriteLevel"
}
```

Triggers the shade's preset "favorite" position (configured in the Lutron app).

**Get shade status:**
```json
{
  "CommuniqueType": "ReadRequest",
  "Header": {
    "Url": "/zone/{zone_id}/status"
  }
}
```

**Status response:**
```json
{
  "CommuniqueType": "ReadResponse",
  "Body": {
    "ZoneStatus": {
      "Zone": {"href": "/zone/{zone_id}"},
      "Level": 75,
      "StatusAccuracy": "Good"
    }
  }
}
```

### Capability Mapping

| Haus Capability | LEAP Action |
|----------------|-------------|
| `on_off` (open) | `GoToLevel` with `Level: 100` |
| `on_off` (close) | `GoToLevel` with `Level: 0` |
| `brightness` (position) | `GoToLevel` with `Level: 0-100` |

## AI Capabilities

> AI integration planned via Lutron bridge. When available:
> - Open/close shades by room name
> - Set specific position percentages
> - Trigger favorite/preset positions
> - Report current shade position
> - Coordinate shades with scenes (e.g., "Movie mode" closes all living room shades)

## Quirks & Notes

- **Not directly IP-addressable** -- Serena shades use Lutron's Clear Connect RF protocol and require a bridge for any smart home integration
- **Bridge requirement** -- Must use the Caseta Smart Bridge Pro (L-BDGPRO2-WH) or RA2 Select/RadioRA 3 for third-party LEAP protocol access; the standard Caseta bridge does not support LEAP
- **Battery life** -- Typical battery life is 3-5 years with normal use (D-cell batteries) or 1-2 years with the rechargeable battery wand
- **Position accuracy** -- `StatusAccuracy` in the response indicates position confidence; RF communication means the bridge may not always know the exact shade position
- **Quiet operation** -- Serena shades are notably quieter than many competitors due to the motor design
- **Custom sizing** -- Shades are custom-ordered to exact window dimensions from lutron.com or authorized dealers
- **Tilt vs. lift** -- Honeycomb and roller shades support lift (open/close); wood blinds also support tilt (angle adjustment) which maps to a separate zone
- **Group control** -- Multiple shades can be grouped in the Lutron app and controlled as a single zone
- **Clear Connect range** -- RF range is approximately 30 feet through walls; signal repeaters (Pico remotes, other Lutron devices) extend range
- **Price** -- Serena shades are premium-priced ($300-$1000+ per window); the Caseta bridge adds $100-$140

## Similar Devices

> - [IKEA FYRTUR Smart Blinds](ikea-fyrtur-smart-blinds.md) -- Budget Zigbee alternative
> - [Hunter Douglas PowerView](hunter-douglas-powerview.md) -- Premium competitor with local REST API
