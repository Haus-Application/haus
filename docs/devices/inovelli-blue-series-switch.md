---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "inovelli-blue-series-switch"
name: "Inovelli Blue Series 2-1 Switch"
manufacturer: "Inovelli, LLC"
brand: "Inovelli"
model: "VZM31-SN"
model_aliases: ["VZM31-SN-BLU", "Blue Series", "2-in-1 Switch + Dimmer"]
device_type: "zigbee_switch"
category: "lighting"
product_line: "Blue Series"
release_year: 2022
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
  mac_prefixes: []       # Zigbee devices do not have WiFi/Ethernet MACs visible on IP network
  mdns_services: []      # Not an IP device
  mdns_txt_keys: []
  default_ports: []      # No IP ports
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
  integration_key: "matter"
  polling_interval_sec: 0
  websocket_event: "matter:state"
  setup_type: "app_pairing"
  ai_chattable: true
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities:
  - "on_off"
  - "brightness"

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "Zigbee"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Zigbee 3.0 pairing via install code or open network join. Matter commissioning available with Matter firmware via Thread border router or Zigbee bridge."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "switch"
  power_source: "hardwired"
  mounting: "in_wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee"]

# --- LINKS ---
links:
  product_page: "https://inovelli.com/products/blue-series-smart-2-1-switch-on-off-or-dimmer"
  api_docs: "https://inovelli.com/pages/zigbee-technical-docs"
  developer_portal: ""
  support: "https://support.inovelli.com/"
  community_forum: "https://community.inovelli.com/"
  image_url: ""
  fcc_id: "2AW4GVZM31SN"

# --- TAGS ---
tags: ["zigbee", "in_wall", "dimmer", "switch", "led_bar", "matter_capable", "silicon_labs", "no_neutral_optional"]
---

# Inovelli Blue Series 2-1 Switch

## What It Is

The Inovelli Blue Series 2-1 Switch (model VZM31-SN) is a premium in-wall smart switch that functions as either an on/off switch or a dimmer, configurable via software. It communicates via Zigbee 3.0 and requires a Zigbee coordinator (hub) to operate in a smart home context. The standout feature is a fully configurable 16-million-color LED notification bar on the paddle face, which can display colors, effects, and brightness levels for notifications (e.g., red pulsing when a door is open, blue when the washing machine finishes). The switch is built on a Silicon Labs EFR32MG24 chipset and Inovelli has released beta Matter-over-Thread firmware, making it a potential future Matter device. It supports single-pole, 3-way (with dumb or aux switch), and multi-way configurations, and has an optional no-neutral-wire mode for older homes. It is extremely popular in the Home Assistant and Zigbee2MQTT communities.

## How Haus Discovers It

The VZM31-SN is a Zigbee device and is not directly discoverable on the IP network. Haus discovers it through a Zigbee coordinator:

1. **Zigbee Network Join** -- The switch is put into pairing mode by pressing the config button 3x rapidly. A Zigbee coordinator (e.g., a Zigbee USB dongle, Hue Bridge with third-party Zigbee support, or a Thread border router with Matter firmware) detects the device joining.
2. **Zigbee Device Interview** -- Once joined, the coordinator reads the device's Zigbee clusters. The VZM31-SN reports manufacturer "Inovelli" and model "VZM31-SN". It supports the following Zigbee clusters:
   - `0x0000` -- Basic
   - `0x0003` -- Identify
   - `0x0004` -- Groups
   - `0x0005` -- Scenes
   - `0x0006` -- On/Off
   - `0x0008` -- Level Control
   - `0xFC31` -- Inovelli custom cluster (LED bar, config parameters)
3. **Matter Commissioning (Future)** -- With Matter firmware, the switch joins a Thread network and is commissioned via a Matter QR code or numeric pairing code. Haus would discover it through Matter's mDNS-based commissioning advertisement (`_matterc._udp` or `_matterd._udp`).

## Pairing / Authentication

### Zigbee Pairing

1. Pull the air-gap switch (bottom tab) out and push it back in to power-cycle the switch.
2. Press the config button (small button at top of paddle) 3 times quickly. The LED bar will pulse blue, indicating pairing mode.
3. The Zigbee coordinator must have its network open for joining (permit join mode).
4. The switch joins and the LED bar turns solid green briefly to confirm.
5. Factory reset: hold the config button for 20+ seconds until the LED bar turns red.

### Matter Commissioning (Beta Firmware)

1. Flash Matter firmware via OTA from Inovelli or via manual firmware file.
2. The switch advertises on Thread as a Matter-commissionable device.
3. Scan the Matter QR code or enter the 11-digit pairing code.
4. Matter controller (Haus) commissions the device onto the Thread/Matter fabric.

## API Reference

### Zigbee Clusters

The VZM31-SN is controlled via standard Zigbee clusters plus Inovelli's custom cluster.

#### On/Off (Cluster 0x0006)

Standard Zigbee on/off cluster:
- `On` command (0x01) -- turn on
- `Off` command (0x00) -- turn off
- `Toggle` command (0x02) -- toggle state

#### Level Control (Cluster 0x0008)

Standard Zigbee level control for dimming:
- `MoveToLevel` (0x00) -- set brightness 0-254 with transition time
- `Move` (0x01) -- start continuous dimming up/down
- `Step` (0x02) -- step brightness by amount
- `Stop` (0x03) -- stop continuous dimming
- Attribute `CurrentLevel` (0x0000) -- current brightness 0-254

#### Inovelli Custom Cluster (0xFC31)

The custom cluster provides access to over 50 configuration parameters:

| Parameter | ID | Description | Range |
|-----------|-----|-------------|-------|
| Dimming Speed (up) | 1 | Remote dimming speed | 0-127 (seconds) |
| Dimming Speed (down) | 5 | Remote dimming speed down | 0-127 |
| Ramp Rate (on) | 3 | Physical ramp rate on | 0-127 |
| Ramp Rate (off) | 7 | Physical ramp rate off | 0-127 |
| Minimum Level | 9 | Minimum dim level | 1-254 |
| Maximum Level | 10 | Maximum dim level | 2-254 |
| Auto Off Timer | 12 | Auto-off in seconds | 0-32767 |
| Default LED Color | 95 | LED bar color (0-255 hue) | 0-255 |
| Default LED Intensity | 96 | LED bar brightness when on | 0-100 |
| Default LED Intensity (off) | 97 | LED bar brightness when off | 0-100 |
| Switch Mode | 258 | 0=Dimmer, 1=On/Off | 0 or 1 |
| LED Effect Type | 99 | Notification effect type | See below |
| LED Effect Color | 100 | Notification effect color | 0-255 (hue wheel) |

#### LED Notification Bar Commands

The LED bar supports notification effects sent via the custom cluster:

**Effect Types:**
- 0 = Off
- 1 = Solid
- 2 = Fast Blink
- 3 = Slow Blink
- 4 = Pulse
- 5 = Chase
- 6 = Open/Close (rising/falling)
- 7 = Small to Big
- 8 = Aurora
- 9 = Slow Falling
- 10 = Medium Falling
- 11 = Fast Falling
- 12 = Slow Rising
- 13 = Medium Rising
- 14 = Fast Rising
- 15 = Medium Blink
- 16 = Slow Chase
- 17 = Fast Chase
- 18 = Fast Siren
- 19 = Slow Siren

The notification is sent as a composite value encoding color (0-255), brightness (0-10), duration (1-255 seconds, 255=indefinite), and effect type into a single 32-bit parameter on attribute ID 16 of the custom cluster.

### Matter API (Future)

With Matter firmware, the switch exposes standard Matter clusters:
- On/Off cluster (same semantics as Zigbee)
- Level Control cluster (same semantics)
- Custom Inovelli cluster support TBD in Matter firmware

## AI Capabilities

When the AI concierge is chatting with an Inovelli Blue Series switch, it can:

- **Turn the switch on or off**
- **Set dimming level** as a percentage (0-100%)
- **Set LED bar color** for status notifications (e.g., "show red on the LED bar")
- **Activate LED effects** for visual alerts (pulse, chase, blink)
- **Report current state** -- on/off, brightness level, switch mode
- **Change configuration** -- minimum/maximum levels, ramp rates, auto-off timer

## Quirks & Notes

- **No Direct IP Access:** This is a Zigbee device with no WiFi or Ethernet interface. It cannot be discovered or controlled directly over the IP network. All communication goes through a Zigbee coordinator.
- **Silicon Labs EFR32MG24:** The chipset supports both Zigbee 3.0 and Thread/Matter via firmware swap. You cannot run both simultaneously.
- **Matter Firmware is Beta:** As of early 2026, the Matter firmware is functional but still in beta. Some features (particularly LED notifications) may not be fully exposed via Matter clusters yet.
- **No Neutral Required (with limitations):** The switch can work without a neutral wire, but minimum load requirements apply (about 25W for incandescent, higher for LED). A bypass (Inovelli LZW36-BYPASS) may be needed for low-wattage LED loads.
- **Multi-Tap Actions:** The switch supports multi-tap actions (2x, 3x, 4x, 5x tap up/down) that can be bound to automations. These are reported as Zigbee scene events.
- **Power Reporting:** The switch can report instantaneous wattage and energy consumption via the Zigbee Metering cluster (0x0702), though accuracy varies.
- **Air Gap Switch:** The bottom tab physically disconnects power to the load. This is required by electrical code for in-wall switches and also serves as a factory reset mechanism (pull out, hold config button, push back in).
- **OTA Updates:** Inovelli provides firmware updates via Zigbee OTA cluster. Updates can be applied through Zigbee2MQTT, ZHA, or the Inovelli firmware tool.
- **Zigbee2MQTT / ZHA:** The VZM31-SN is fully supported in both Zigbee2MQTT and Home Assistant ZHA with all parameters exposed.

## Similar Devices

- **inovelli-red-series-switch** -- Previous generation Z-Wave version (VZW31-SN), similar features but Z-Wave protocol
- **zooz-zwave-switch-zen76** -- Z-Wave on/off switch, simpler but no LED bar
- **leviton-decora-smart-dimmer** -- WiFi dimmer, different connectivity approach
- **lutron-caseta-dimmer** -- Proprietary Clear Connect protocol, requires Lutron bridge
- **shelly-plus-1** -- WiFi relay, in-wall but no dimming or LED bar
