---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "shark-ai-robot-vacuum"
name: "Shark AI Robot Vacuum"
manufacturer: "SharkNinja Operating LLC"
brand: "Shark"
model: "AI Robot Vacuum"
model_aliases: ["RV2502AE", "AV2501S", "AV2501AE", "Shark AI VacMop", "Shark Matrix", "Shark Detect Pro"]
device_type: "robot_vacuum"
category: "smart_home"
product_line: "Shark AI"
release_year: 2022
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi", "bluetooth_le"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^Shark.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "not_feasible"
  integration_key: ""
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities: ["on_off"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "HTTPS"
  encoding: "binary"
  auth_method: "oauth2"
  auth_detail: "Entirely cloud-based. SharkClean app communicates with Shark cloud servers. No known local protocol. No community reverse-engineering efforts have produced a usable API. The robot communicates outbound to Shark cloud infrastructure only."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "battery"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.sharkclean.com/robot-vacuums/"
  api_docs: ""
  developer_portal: ""
  support: "https://www.sharkclean.com/support/"
  community_forum: ""
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: ["vacuum", "robot-vacuum", "cloud-only", "no-local-api", "no-api", "closed-ecosystem", "not-feasible"]
---

# Shark AI Robot Vacuum

## What It Is

> The Shark AI Robot Vacuum line (including the Matrix, Detect Pro, and VacMop variants) is a series of WiFi-connected robot vacuums from SharkNinja. Models feature LiDAR or camera-based navigation, AI-powered object detection, and optional self-emptying bases. They are controlled exclusively through the SharkClean app and are entirely cloud-dependent. SharkNinja does not provide any API (public or private), does not expose any local network services, and has no known protocol that has been successfully reverse-engineered by the community. This makes Shark robot vacuums effectively not feasible for third-party smart home integration.

## How Haus Discovers It

1. **Hostname pattern** -- DHCP hostname may contain `Shark`
2. **No further identification** -- No known MAC OUI prefixes specific to Shark robot vacuums, no mDNS services, no open ports
3. **Network traffic analysis** -- The robot communicates outbound over HTTPS to Shark cloud servers; no inbound connections are possible

## Pairing / Authentication

> Pairing is handled exclusively through the SharkClean mobile app via BLE. There is no way to authenticate or control the device outside of the official app.

## API Reference

> No API exists. SharkNinja does not offer:
> - A public developer API
> - A documented cloud API
> - Any local network API
> - Works with Alexa/Google integration only, but these integrations go through Shark's cloud servers
>
> No community projects have successfully reverse-engineered a usable control protocol. The SharkClean app communicates with Shark's backend over HTTPS with certificate pinning, making traffic interception difficult.

## AI Capabilities

> Not applicable. Integration is not feasible.

## Quirks & Notes

- **Completely closed ecosystem** -- SharkNinja has shown no interest in opening their platform to third-party developers
- **No Home Assistant integration** -- Unlike most smart home devices, there is no community-maintained Home Assistant integration for Shark vacuums
- **Certificate pinning** -- The SharkClean app uses certificate pinning, preventing MITM-based API discovery
- **Alexa/Google only** -- Voice control through Alexa and Google Assistant is supported, but this routes through Shark's cloud servers and does not expose any usable API endpoints to third parties
- **Model confusion** -- SharkNinja frequently refreshes model numbers and names (AI Robot Vacuum, Matrix, Detect Pro, VacMop, PowerDetect) while maintaining the same closed cloud architecture
- **Not recommended** -- For users wanting local smart home integration with robot vacuums, Roborock (with Valetudo) or iRobot (with local MQTT) are significantly better choices
- **Matter support** -- No current or announced Matter support
- **BLE only for setup** -- Bluetooth LE is used only during initial WiFi provisioning; no ongoing BLE control

## Similar Devices

> - [iRobot Roomba j7+](irobot-roomba-j7-plus.md) -- Competing robot vacuum with local MQTT hack
> - [Roborock S8 Pro Ultra](roborock-s8-pro-ultra.md) -- Competing robot vacuum with Valetudo local control option
> - [Ecovacs Deebot X2 Omni](ecovacs-deebot-x2-omni.md) -- Competing cloud-based robot vacuum
