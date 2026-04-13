---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "sengled-smart-bulb-zigbee"
name: "Sengled Smart LED Bulb (Zigbee)"
manufacturer: "Sengled"
brand: "Sengled"
model: "E11-G13"
model_aliases: ["E11-G14", "E11-N1EA", "E12-N1E", "E21-N1EA", "E11-U2E", "E11-U3E"]
device_type: "sengled_zigbee_bulb"
category: "lighting"
product_line: "Sengled Smart"
release_year: 2018
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
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
  status: "detected_only"
  integration_key: ""
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "none"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp"]

# --- PROTOCOL ---
protocol:
  type: ""
  port: 0
  transport: ""
  encoding: ""
  auth_method: ""
  auth_detail: ""
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "bulb"
  power_source: "mains"
  mounting: "ceiling"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee"]

# --- LINKS ---
links:
  product_page: "https://www.sengled.com/products/smart-led-multicolor-a19-bulb"
  api_docs: ""
  developer_portal: ""
  support: "https://support.sengled.com/"
  community_forum: "https://github.com/home-assistant/core/tree/dev/homeassistant/components/sengled"
  image_url: ""
  fcc_id: "2AGWR-E11G13"

# --- TAGS ---
tags: ["zigbee_bulb", "no_local_api", "requires_hub", "budget", "no_zigbee_router"]
---

# Sengled Smart LED Bulb (Zigbee)

## What It Is

Sengled Smart LED Bulbs are budget-friendly Zigbee-only smart bulbs. Unlike most Zigbee bulbs, Sengled intentionally omits the Zigbee router function — these bulbs are Zigbee end devices only, meaning they do not extend the Zigbee mesh network. This design choice was made to improve network stability (bulbs that are turned off at the wall switch do not leave a gap in the mesh). They require a Zigbee hub to function: the Sengled Element Hub (now discontinued), Samsung SmartThings, Amazon Echo Plus/Echo Show (4th gen with Zigbee), Hubitat, or any standard Zigbee 3.0 coordinator. Models range from dimmable white (E11-G13) to full RGB color (E11-N1EA, E21-N1EA). Priced at roughly $7-12 per bulb, they are among the cheapest smart bulbs available.

## How Haus Discovers It

Sengled Zigbee bulbs have **no IP network presence**. They communicate exclusively over Zigbee and do not appear on the WiFi/Ethernet network at all. Haus cannot directly discover these devices through network scanning.

**Indirect detection paths:**
1. **Via a paired Zigbee hub**: If the user has a SmartThings, Hubitat, or similar hub that Haus integrates with, Sengled bulbs will appear in that hub's device list with manufacturer string "Sengled" and model IDs like "E11-G13".
2. **Via Zigbee2MQTT or ZHA**: If the user runs a Zigbee USB coordinator with Zigbee2MQTT or Home Assistant ZHA, Sengled devices are exposed via MQTT or the HA API.
3. **User declaration**: The user tells Haus they have Sengled bulbs, and Haus records them as detected but not directly controllable.

Haus will flag Sengled bulbs as "detected_only" with guidance that they require a supported Zigbee hub.

## Pairing / Authentication

There is no direct pairing with Haus. The bulbs pair with their Zigbee coordinator:

1. Power cycle the bulb (off for 2 seconds, on for 2 seconds) 10 times to enter pairing mode.
2. The bulb blinks to indicate pairing mode.
3. The Zigbee coordinator discovers and pairs the bulb.

If Haus later gains integration with the user's Zigbee coordinator (e.g., SmartThings, Hubitat), the Sengled bulbs become controllable through that hub's API.

## API Reference

Sengled Zigbee bulbs have no direct local API. They speak standard Zigbee Cluster Library (ZCL) commands:

- **On/Off Cluster (0x0006)**: Turn on, turn off, toggle
- **Level Control Cluster (0x0008)**: Set brightness (0-254)
- **Color Control Cluster (0x0300)**: Set hue, saturation, color temperature (on color models)

These clusters are exposed through whatever Zigbee coordinator the bulb is paired with. There is no IP-based API.

### Sengled Cloud API

Sengled has a cloud API accessible via their app, but it is undocumented, authenticated via Sengled account credentials, and provides no local access. It is not a viable integration path for Haus.

## AI Capabilities

Not applicable. Haus cannot directly communicate with Sengled Zigbee bulbs without an intermediary hub integration.

## Quirks & Notes

- **Not a Zigbee router**: This is Sengled's most notable design decision. Unlike Hue, TRADFRI, and most other Zigbee bulbs, Sengled bulbs are end devices only. They do not repeat/route Zigbee traffic. This means turning off a Sengled bulb at the wall switch does not disrupt the Zigbee mesh, but it also means Sengled bulbs do not extend mesh range.
- **Sengled Element Hub discontinued**: Sengled's own hub (the Element Hub) has been discontinued. Sengled now recommends using third-party hubs like SmartThings.
- **WiFi models exist**: Sengled also makes WiFi bulbs (model numbers with "W" suffix, e.g., W11-U2E), which use a cloud-only API. These are different products from the Zigbee models documented here.
- **Matter models**: Sengled has announced Matter-over-Thread bulbs. These would be directly controllable via Matter integration when Haus adds Thread Border Router support.
- **No OUI/MAC**: Since these are Zigbee-only devices with no WiFi or Ethernet radio, there are no MAC prefixes on the IP network to detect.
- **Power-on behavior**: Sengled bulbs default to turning on at full brightness when power is restored, which is the expected behavior for bulbs controlled by wall switches.
- **Model confusion**: The Sengled product line has dozens of model numbers across generations. The E11/E12/E21 prefix generally indicates Zigbee, while W11/W21 indicates WiFi.

## Similar Devices

- [wiz-connected-bulb](wiz-connected-bulb.md) — WiFi bulb with local API (no hub required)
- [kasa-smart-bulb-kl130](kasa-smart-bulb-kl130.md) — WiFi bulb with local API (no hub required)
- [cync-smart-bulb](cync-smart-bulb.md) — WiFi+BLE bulb, also lacking local API
