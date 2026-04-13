---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "philips-hue-bulb-a19"
name: "Philips Hue White and Color Ambiance A19"
manufacturer: "Signify Netherlands B.V."
brand: "Philips Hue"
model: "9290024688"
model_aliases: ["548727", "563585", "9290022169", "9290012573", "LCA001", "LCA003", "LCA007"]
device_type: "hue_light"
category: "lighting"
product_line: "Hue"
release_year: 2020
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
  protocols_spoken: ["zigbee", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
# Zigbee device -- no direct IP network presence.
# Controlled exclusively via the Hue Bridge API.
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
  status: "supported"
  integration_key: "hue"
  polling_interval_sec: 5
  websocket_event: "hue:state"
  setup_type: "link_button"
  ai_chattable: true
  haus_milestone: "M3"

# --- CAPABILITIES ---
capabilities:
  - "on_off"
  - "brightness"
  - "color"
  - "color_temp"

# --- PROTOCOL ---
# No direct protocol -- accessed via bridge CLIP v2 API.
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "link_button"
  auth_detail: "Accessed indirectly via Hue Bridge API. Bridge handles Zigbee communication."
  base_url_template: "https://{bridge_ip}/clip/v2/resource/light/{id}"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "bulb"
  power_source: "mains"
  mounting: "ceiling"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.philips-hue.com/en-us/p/hue-white-and-color-ambiance-1-pack-e26/046677548728"
  api_docs: "https://developers.meethue.com/develop/hue-api-v2/"
  developer_portal: "https://developers.meethue.com/"
  support: "https://www.philips-hue.com/en-us/support"
  community_forum: "https://developers.meethue.com/forum/"
  image_url: ""
  fcc_id: "2ABA6-LCA001"

# --- TAGS ---
tags: ["zigbee", "color", "ambiance", "e26", "a19", "bluetooth_provisioning", "via_bridge"]
---

# Philips Hue White and Color Ambiance A19

## What It Is

The Philips Hue White and Color Ambiance A19 is a standard E26-base LED smart bulb that produces up to 1100 lumens (75W equivalent) of tunable white light (2000K-6500K) and 16 million colors via its RGB+WW LED array. It communicates over Zigbee 3.0 (Zigbee Light Link profile) and requires a Philips Hue Bridge for full functionality. Bluetooth LE is included for limited direct control from the Hue app (up to 10 bulbs, no scenes/automations). This is the most common bulb in the Hue ecosystem and the one most users encounter first.

The bulb has gone through several hardware generations, with model identifiers including LCA001 (Gen 3, Bluetooth-enabled), LCA003 (higher lumen), and LCA007 (latest revision with richer color gamut). Retail SKUs include 548727 and 563585 depending on region and packaging.

## How Haus Discovers It

This device has **no direct network presence** -- it communicates exclusively over Zigbee via the Hue Bridge.

1. **Bridge Discovery** -- Haus first discovers and pairs with the Hue Bridge (see `philips-hue-bridge`).
2. **Device Enumeration** -- `GET /clip/v2/resource/device` returns all paired Zigbee devices. Each device entry includes its `product_data.model_id` (e.g., "LCA001") and `product_data.product_name`.
3. **Light Resource** -- Each bulb exposes a `light` service accessible via `GET /clip/v2/resource/light`. The light resource contains the bulb's current state and capabilities.
4. **Capability Detection** -- Haus reads the light resource to determine what the bulb supports:
   - `color` object present = full color support
   - `color_temperature` object present = tunable white support
   - `dimming` object present = dimmable
   - `on` object = power control

## Pairing / Authentication

The bulb itself requires no separate pairing with Haus. It is paired to the Hue Bridge via the Hue app or the bridge's Zigbee touchlink process.

### How a Bulb Joins the Bridge

1. Power on the bulb (factory-fresh or after factory reset).
2. In the Hue app, tap "Add light" -- the bridge initiates a Zigbee network scan.
3. The bulb joins the bridge's Zigbee network and is assigned a device ID.
4. From Haus's perspective, the bulb simply appears in the bridge's device/light listing.

### Factory Reset

A Hue bulb can be factory reset by:
- Using the Hue app's "Delete" function
- Toggling power off-on in a specific pattern (5 cycles of off 5s / on 8s)
- Using a Hue Dimmer Switch held close to the bulb (hold all 4 buttons for 5 seconds)

## API Reference

All control is via the Hue Bridge CLIP v2 API. See `philips-hue-bridge` for full endpoint documentation.

### Get Light State

```
GET /clip/v2/resource/light/{light_id}
```

**Response fields for an A19 color bulb:**
```json
{
  "id": "a1b2c3d4-...",
  "type": "light",
  "metadata": {
    "name": "Living Room Lamp",
    "archetype": "sultan_bulb"
  },
  "on": {"on": true},
  "dimming": {
    "brightness": 75.0,
    "min_dim_level": 0.2
  },
  "color": {
    "xy": {"x": 0.4578, "y": 0.4101},
    "gamut": {
      "red":   {"x": 0.6915, "y": 0.3083},
      "green": {"x": 0.1700, "y": 0.7000},
      "blue":  {"x": 0.1532, "y": 0.0475}
    },
    "gamut_type": "C"
  },
  "color_temperature": {
    "mirek": 370,
    "mirek_valid": true,
    "mirek_schema": {"mirek_minimum": 153, "mirek_maximum": 500}
  },
  "dynamics": {
    "status": "none",
    "speed": 0.0
  }
}
```

### Control Light

```
PUT /clip/v2/resource/light/{light_id}
```

```json
{
  "on": {"on": true},
  "dimming": {"brightness": 50.0},
  "color": {"xy": {"x": 0.675, "y": 0.322}},
  "dynamics": {"duration": 500}
}
```

### Color Gamut

The A19 Gen 3+ (LCA001 and later) uses **Gamut C**, the widest color gamut in the Hue lineup:

| Point | CIE x | CIE y |
|-------|--------|--------|
| Red | 0.6915 | 0.3083 |
| Green | 0.1700 | 0.7000 |
| Blue | 0.1532 | 0.0475 |

Older A19 models may use Gamut A or B with a narrower range. The gamut is reported in the `color.gamut` field.

## AI Capabilities

When the AI interacts with this bulb (via the bridge), it can:

- **Turn on/off** by name ("turn off the bedroom lamp")
- **Set brightness** as percentage ("dim the hallway to 30%")
- **Set color** by name or description ("make the living room lamp red", "set a warm sunset glow")
- **Set color temperature** ("set to daylight", "make it warm and cozy")
- **Transition smoothly** using dynamics duration
- **Report current state** ("The kitchen light is on at 80%, warm white")

## Quirks & Notes

- **Minimum Brightness:** The bulb cannot dim below ~0.2% (`min_dim_level`). Setting brightness to 0 does not turn the bulb off -- you must set `on.on` to false.
- **Color vs Color Temp:** Setting `color.xy` puts the bulb in color mode; setting `color_temperature.mirek` puts it in white mode. These are mutually exclusive -- the last one set wins.
- **Bluetooth Limitation:** When controlled via Bluetooth directly (no bridge), features are limited to on/off, brightness, color, and color temperature. No scenes, automations, or room grouping. Maximum 10 bulbs.
- **Power-On Behavior:** The bulb defaults to "last state" on power restore, but this can be configured via the bridge to always turn on, always stay off, or return to a specific state.
- **Zigbee Repeater:** Like all mains-powered Zigbee devices, the A19 acts as a Zigbee router/repeater, extending the mesh network range.
- **Max Distance from Bridge:** Zigbee signal typically reaches 30-50 feet indoors. The mesh network means bulbs relay signals to each other.
- **Wattage:** 9W actual power consumption for 1100 lumens output (LCA001).
- **Lifespan:** Rated for 25,000 hours (approximately 22 years at 3 hours/day).

## Similar Devices

- **philips-hue-bridge** -- Required hub for full functionality
- **philips-hue-lightstrip-plus** -- LED strip with similar color capabilities but different form factor
- **philips-hue-sync-box** -- HDMI sync device that can coordinate with this bulb for entertainment
