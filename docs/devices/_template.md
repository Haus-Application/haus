---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: ""                        # unique slug, kebab-case (e.g. "philips-hue-bridge")
name: ""                      # display name
manufacturer: ""              # legal entity (e.g. "Signify")
brand: ""                     # consumer brand (e.g. "Philips Hue")
model: ""                     # primary model number
model_aliases: []             # other model strings seen in the wild
device_type: ""               # haus internal type (e.g. "hue_bridge", "kasa_dimmer")
category: ""                  # lighting | media | security | climate | energy | smart_home | network | compute
product_line: ""              # family grouping (e.g. "Hue", "Kasa", "Ring Alarm")
release_year: null            # approximate retail launch
discontinued: false
price_range: ""               # "$", "$$", "$$$", "$$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: ""                    # "local" | "cloud" | "hybrid"
  local_api: false
  cloud_api: false
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: false
  protocols_spoken: []        # ["wifi", "zigbee", "zwave", "thread", "matter", "bluetooth", "ethernet"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []            # OUI prefixes for discovery (e.g. ["00:17:88"])
  mdns_services: []           # mDNS service types (e.g. ["_hue._tcp"])
  mdns_txt_keys: []           # notable TXT record keys
  default_ports: []           # commonly open ports
  signature_ports: []         # ports that strongly ID this device
  ssdp_search_target: ""      # UPnP ST header
  ssdp_server_string: ""      # UPnP Server header
  hostname_patterns: []       # regex patterns for hostnames
  ip_ranges: []               # known static/DHCP ranges if applicable

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 0
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: ""                  # "supported" | "read_only" | "detected_only" | "planned" | "not_feasible"
  integration_key: ""         # matches DeviceProbeResult.Integration (e.g. "hue", "kasa")
  polling_interval_sec: 0
  websocket_event: ""         # e.g. "hue:state"
  setup_type: ""              # "none" | "link_button" | "oauth2" | "password" | "api_key" | "app_pairing"
  ai_chattable: false         # AI can "speak as" this device
  haus_milestone: ""          # which milestone adds support (e.g. "M3")

# --- CAPABILITIES ---
capabilities: []              # ["on_off", "brightness", "color", "color_temp", "scenes", "groups",
                              #  "motion", "temperature", "humidity", "lock_unlock", "garage_open_close",
                              #  "media_playback", "volume", "input_select", "fan_speed", "thermostat",
                              #  "camera_stream", "camera_snapshot", "doorbell", "arm_disarm", "battery_level"]

# --- PROTOCOL ---
protocol:
  type: ""                    # "https_rest" | "http_rest" | "tcp_xor" | "websocket_json" | "protobuf_tls" | "mqtt" | "coap" | "proprietary"
  port: 0
  transport: ""               # TCP | HTTP | HTTPS | WebSocket | UDP | TLS
  encoding: ""                # JSON | Protobuf | XML | XOR-JSON | binary | CBOR
  auth_method: ""             # "none" | "api_key" | "basic_auth" | "oauth2" | "link_button" | "session_cookie"
  auth_detail: ""             # e.g. header name, token endpoint
  base_url_template: ""       # e.g. "https://{ip}/clip/v2"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: ""             # hub | bulb | switch | plug | sensor | panel | receiver | camera | thermostat | controller | gateway | speaker | display | lock | shade | strip | outdoor_fixture
  power_source: ""            # mains | battery | usb | poe | hardwired | solar
  mounting: ""                # shelf | wall | ceiling | in_wall | outdoor | tabletop | door
  indoor_outdoor: ""          # indoor | outdoor | both
  wireless_radios: []         # ["wifi", "zigbee", "zwave", "thread", "bluetooth_le", "matter"]

# --- LINKS ---
links:
  product_page: ""
  api_docs: ""
  developer_portal: ""
  support: ""
  community_forum: ""
  image_url: ""
  fcc_id: ""

# --- TAGS ---
tags: []
---

# {Device Name}

## What It Is

> One-paragraph summary of the device, what it does, and why someone buys it.

## How Haus Discovers It

> Step-by-step: OUI match → mDNS → port probe → HTTP fingerprint → protocol probe.
> Only include steps that apply.

## Pairing / Authentication

> Setup flow. Omit section if no auth required.

## API Reference

> Full endpoint documentation. Omit section for devices with no usable API.

## AI Capabilities

> What the AI concierge can do when "chatting as" this device.

## Quirks & Notes

> Gotchas, edge cases, rate limits, firmware differences.

## Similar Devices

> Cross-references to related entries in the knowledge base.
