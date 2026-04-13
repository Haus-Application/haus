---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "brilliant-smart-home-control"
name: "Brilliant Smart Home Control Panel"
manufacturer: "Brilliant"
brand: "Brilliant"
model: "BHA120US-WH"
model_aliases: ["BHA120US-BK", "BHA230US-WH"]
device_type: "smart_panel"
category: "smart_home"
product_line: "Brilliant"
release_year: 2019
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi", "bluetooth", "zigbee"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["B8:B7:F1", "CC:DB:A7"]
  mdns_services: ["_brilliant._tcp"]
  mdns_txt_keys: ["home_id"]
  default_ports: [5455]
  signature_ports: [5455]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^brilliant.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "detected_only"
  integration_key: "brilliant"
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "scenes", "motion"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 5455
  transport: "TCP"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Proprietary binary/protobuf protocol. Encrypted communication. Cloud authentication via Brilliant app required. Local protocol is not publicly documented."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "panel"
  power_source: "hardwired"
  mounting: "in_wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le", "zigbee"]

# --- LINKS ---
links:
  product_page: "https://www.brilliant.tech"
  api_docs: ""
  developer_portal: ""
  support: "https://support.brilliant.tech"
  community_forum: ""
  image_url: ""
  fcc_id: "2AQ4B-BHA120"

# --- TAGS ---
tags: ["panel", "touchscreen", "in-wall", "proprietary", "detected-only", "intercom", "motion-sensor", "dimmer"]
---

# Brilliant Smart Home Control Panel

## What It Is

> The Brilliant Smart Home Control panel is an in-wall touchscreen device that replaces a standard light switch. It features a built-in LCD display, motion sensor, ambient light sensor, and speaker/microphone for intercom functionality. The panel provides built-in light dimming for hardwired loads, scene control, and integrations with third-party platforms (Ring, Sonos, Honeywell, SmartThings, HomeKit). Multiple Brilliant panels in a home share a `home_id` for synchronized control. The device communicates via a proprietary binary protocol on port 5455 that is not publicly documented.

## How Haus Discovers It

1. **mDNS** -- Advertises as `_brilliant._tcp` on the local network
2. **mDNS TXT records** -- Contains `home_id` key linking panels in the same home
3. **Port probe** -- TCP check on port 5455 confirms the proprietary protocol endpoint is listening

Example mDNS record:
```
Instance: 01663ad46e010003271243acea345b26
Port: 5455
TXT: home_id=01919f35f09c00024384dccbd839c147
```

## Pairing / Authentication

> **Not yet achievable for local control.** The Brilliant protocol on port 5455 is proprietary and requires cloud authentication through the Brilliant mobile app. Local control requires reverse engineering the binary protocol or using an alternative integration path.
>
> **Potential future approaches:**
> 1. **Brilliant Cloud API** -- may offer OAuth-based control (requires Brilliant account)
> 2. **Local protocol reverse engineering** -- binary/protobuf protocol on port 5455
> 3. **HomeKit bridge** -- Brilliant supports Apple HomeKit, which could be used as a control path
> 4. **SmartThings integration** -- Brilliant integrates with SmartThings, which has a documented API

## API Reference

> The Brilliant protocol on port 5455 is not publicly documented. Known characteristics:
>
> - Binary/protobuf-based encoding (not JSON)
> - Encrypted communication
> - Cloud authentication required via Brilliant mobile app
> - Local pairing may require initial setup through the Brilliant app
>
> No usable local API is available at this time.

## AI Capabilities

> Direct device control is not yet available due to the proprietary protocol. When a Brilliant panel is detected, the AI can:
> - **Report discovery info** -- IP address, MAC address, mDNS service, port 5455 status, home_id
> - **Explain the device** -- describe what Brilliant panels are and their built-in capabilities
> - **Suggest connection paths** -- recommend using the Brilliant app, HomeKit, or SmartThings as alternative control methods

## Quirks & Notes

- **Detected but not controllable** -- Haus can discover Brilliant panels via mDNS but cannot control them due to the proprietary protocol
- **home_id links panels** -- multiple Brilliant panels in the same home share a `home_id` in their mDNS TXT records; this can be used to group panels
- **Built-in dimmer** -- the panel physically controls hardwired lights via a built-in dimmer; this works regardless of network connectivity
- **HomeKit support** -- Brilliant panels support Apple HomeKit, which may provide a future integration path for Haus
- **Requires neutral wire** -- like most smart switches, the Brilliant panel requires a neutral wire in the switch box
- **1-gang and 2-gang** -- available in single-switch (BHA120) and double-switch (BHA230) configurations
- **Professional installation recommended** -- the panel requires careful wiring and may need a deeper gang box due to its electronics

## Similar Devices

> The Brilliant panel occupies a unique niche as an in-wall touchscreen controller. No directly comparable devices in the Haus knowledge base. Conceptually similar to Crestron or Control4 in-wall panels but at a consumer price point.
