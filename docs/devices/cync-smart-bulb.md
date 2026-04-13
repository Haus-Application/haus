---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "cync-smart-bulb"
name: "GE Cync Smart Bulb"
manufacturer: "GE Lighting (Savant Systems)"
brand: "Cync"
model: "93128983"
model_aliases: ["C by GE", "GE Cync Full Color A19", "GE Cync Soft White", "GE Cync Direct Connect"]
device_type: "cync_bulb"
category: "lighting"
product_line: "Cync"
release_year: 2021
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
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["34:86:5D", "68:27:19", "C8:2E:18"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^Cync[-_][A-Za-z0-9]+$", "^GE[-_]Light[-_][A-Za-z0-9]+$"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "not_feasible"
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
  type: "proprietary"
  port: 0
  transport: "TLS"
  encoding: "binary"
  auth_method: "oauth2"
  auth_detail: "Cync cloud auth via GE account credentials. Proprietary binary protocol over TLS to cloud servers. No documented local API."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "bulb"
  power_source: "mains"
  mounting: "ceiling"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.cyncsmart.com/"
  api_docs: ""
  developer_portal: ""
  support: "https://www.cyncsmart.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AQHD93128983"

# --- TAGS ---
tags: ["wifi_bulb", "cloud_only", "no_local_api", "ble_setup", "not_feasible", "ge_lighting", "savant"]
---

# GE Cync Smart Bulb

## What It Is

GE Cync (formerly "C by GE") is GE Lighting's consumer smart home line, now owned by Savant Systems. Cync bulbs connect via WiFi and Bluetooth Low Energy and are controlled exclusively through the Cync app and cloud service. They are widely available at retail (Home Depot, Walmart, Amazon) and priced aggressively at $8-15 per bulb. The lineup includes dimmable white, tunable white, and full color models across A19, BR30, and other form factors. Despite the attractive price point and wide availability, Cync bulbs have no documented local API and communicate with GE/Savant cloud servers using a proprietary binary protocol, making local integration infeasible.

## How Haus Discovers It

1. **OUI match**: Cync bulbs use WiFi chipsets with MAC prefixes including 34:86:5D, 68:27:19, and C8:2E:18. These OUIs are associated with the WiFi module manufacturers used by GE Lighting.
2. **Hostname pattern**: Cync bulbs may register DHCP hostnames matching `Cync-*` or `GE_Light_*` patterns, though this varies by firmware version.
3. **No open ports**: Cync bulbs do not expose any open TCP or UDP ports for local control. A port scan will return no results, which itself is a signal when combined with OUI matching.
4. **Cloud traffic analysis**: Cync bulbs communicate with `*.gelighting.com` and `*.cyncsmart.com` servers. DNS query monitoring could identify their presence, though Haus does not use this method.

When Haus detects a likely Cync device via OUI + hostname, it will flag the device as "detected but not controllable locally" and recommend alternatives with local API support.

## Pairing / Authentication

There is no pairing with Haus. Cync bulbs pair exclusively through the Cync mobile app:

1. Create a Cync/GE account in the app.
2. The app discovers nearby bulbs via BLE.
3. WiFi credentials are transferred to the bulb over BLE.
4. The bulb connects to WiFi and registers with GE cloud servers.
5. All subsequent control goes through the cloud.

There is no local control path after setup.

## API Reference

There is no usable local API. The cloud protocol has been partially reverse-engineered by community projects:

- **Cloud endpoint**: Cync bulbs connect to GE/Savant cloud servers via TLS.
- **Protocol**: Proprietary binary protocol (not REST, not JSON). Each message has a binary header with command codes.
- **Authentication**: Requires GE account credentials to obtain session tokens.
- **Community efforts**: The `nikshriv/cync_lights` Home Assistant integration reverse-engineers the cloud API, but it is fragile, breaks with firmware updates, and requires the user's GE account password.

**Haus will NOT implement cloud API integration for Cync** because:
1. It requires storing user's third-party cloud credentials (security risk).
2. The protocol is undocumented and changes without notice.
3. It creates a dependency on GE/Savant cloud availability.
4. It violates Haus's local-first design philosophy.

## AI Capabilities

Not applicable. Haus cannot control Cync devices.

## Quirks & Notes

- **Cloud-only by design**: GE/Savant has made no effort to provide a local API. There is no official developer program for Cync.
- **BLE mesh legacy**: Early "C by GE" products used BLE mesh exclusively (no WiFi). Newer "Cync Direct Connect" products use WiFi. The BLE-only models are even less controllable.
- **Matter announced but limited**: GE has announced Matter support for some Cync products, but rollout has been slow and limited to newer hardware revisions. If/when Cync bulbs gain Matter-over-WiFi support, they would become controllable through Haus's Matter integration.
- **Savant acquisition**: GE Lighting was acquired by Savant Systems in 2020. Savant is a luxury home automation company, and it is unclear how Cync will evolve under their ownership.
- **Google Home preferred partner**: Cync has a deep integration with Google Home, including some local execution through the Google Home hub. This does not help Haus.
- **WiFi congestion**: Like all WiFi bulbs, large Cync installations add many devices to the WiFi network, which can strain consumer routers.
- **Recommendation for users**: If Haus detects Cync bulbs, it should suggest alternatives with local API support: WiZ (similar price, local UDP), Kasa (local TCP), or Philips Hue (premium, local REST).

## Similar Devices

- [wiz-connected-bulb](wiz-connected-bulb.md) — Similar price point but with local UDP API (recommended alternative)
- [kasa-smart-bulb-kl130](kasa-smart-bulb-kl130.md) — WiFi bulb with local TCP API (recommended alternative)
- [sengled-smart-bulb-zigbee](sengled-smart-bulb-zigbee.md) — Also lacks direct local API (requires Zigbee hub)
