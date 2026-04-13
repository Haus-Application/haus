---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "amazon-echo-show-15"
name: "Amazon Echo Show 15"
manufacturer: "Amazon.com, Inc."
brand: "Amazon Echo"
model: "S3JK7Y"
model_aliases: ["Echo Show 15", "Echo Show 15 2nd Gen"]
device_type: "echo_display"
category: "media"
product_line: "Echo"
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
  protocols_spoken: ["wifi", "bluetooth", "zigbee", "matter", "thread"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "FC:65:DE"        # Amazon Technologies Inc.
    - "A0:02:DC"        # Amazon Technologies Inc.
    - "68:54:FD"        # Amazon Technologies Inc.
    - "40:B4:CD"        # Amazon Technologies Inc.
    - "74:C2:46"        # Amazon Technologies Inc.
    - "F0:F0:A4"        # Amazon Technologies Inc.
    - "44:00:49"        # Amazon Technologies Inc.
    - "38:F7:3D"        # Amazon Technologies Inc.
    - "B0:FC:0D"        # Amazon Technologies Inc.
    - "AC:63:BE"        # Amazon Technologies Inc.
    - "14:91:82"        # Amazon Technologies Inc.
    - "18:74:2E"        # Amazon Technologies Inc.
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^amazon-"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "not_feasible"
  integration_key: ""
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: ""
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "media_playback"
  - "volume"
  - "input_select"

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "HTTPS"
  encoding: "Protobuf"
  auth_method: "oauth2"
  auth_detail: "All control flows through Amazon's cloud. No local API. Same cloud-only architecture as all Echo devices."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "display"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le", "zigbee", "thread"]

# --- LINKS ---
links:
  product_page: "https://www.amazon.com/dp/B0BFZVFG6N"
  api_docs: "https://developer.amazon.com/en-US/docs/alexa/alexa-voice-service/api-overview.html"
  developer_portal: "https://developer.amazon.com/alexa"
  support: "https://www.amazon.com/gp/help/customer/display.html?nodeId=201399130"
  community_forum: "https://www.amazonforum.com/s/"
  image_url: ""
  fcc_id: "2AETW-1864"

# --- TAGS ---
tags: ["alexa", "cloud-only", "display", "fire-tv", "zigbee-hub", "thread-border-router", "matter-controller", "voice-assistant", "no-local-api"]
---

# Amazon Echo Show 15

## What It Is

The Amazon Echo Show 15 is a 15.6-inch full HD smart display manufactured by Amazon. It is designed to be wall-mounted (landscape or portrait) or placed on an optional stand, functioning as a family hub with calendar, sticky notes, photo frames, and Fire TV built-in for streaming. It includes an Amazon AZ2 Neural Edge processor, a 5MP camera with visual ID for face recognition, and Alexa voice assistant. The 2nd Gen model (2024) added a built-in smart home hub with Zigbee, Thread, and Matter radios, making it a Thread border router and Matter controller. Like all Echo devices, it is entirely cloud-dependent -- there is no local API for control or state queries.

## How Haus Discovers It

Discovery follows the same limited pattern as the Amazon Echo. See **[amazon-echo-5th-gen](amazon-echo-5th-gen.md)** for the complete discovery process.

1. **OUI Match** -- Amazon Technologies MAC prefix.
2. **Hostname Pattern** -- May appear as `amazon-` prefixed hostname in DHCP.
3. **Classification** -- Haus detects the device on the network by OUI and hostname but cannot integrate with it.

The Echo Show 15 does not advertise any mDNS services or expose HTTP fingerprint endpoints that would allow reliable automated identification.

## Pairing / Authentication

Not applicable for Haus integration. Setup requires the Amazon Alexa app and an Amazon account. All configuration and control flows through Amazon's cloud.

## API Reference

There is no usable local API. See **[amazon-echo-5th-gen](amazon-echo-5th-gen.md)** for a detailed explanation of why local control is not feasible for Amazon Echo devices.

### Fire TV Integration

The Echo Show 15 runs Fire TV OS, which means it can run Fire TV apps and stream content from Amazon Prime Video, Netflix, Disney+, and other services. The Fire TV platform has a remote control API (ADB-based) that some community projects have reverse-engineered, but this requires enabling developer mode on the device and is not a stable or supported integration path.

### Camera

The 5MP camera supports visual ID (face recognition) and can be used for video calls via Alexa. The camera feed is not accessible via any local API -- it is processed on-device for visual ID and streamed through Amazon's cloud for video calls.

## AI Capabilities

Not applicable. The Echo Show 15 cannot be controlled or queried locally.

Haus will display it as a detected network device with a note that it is cloud-dependent and not locally controllable.

## Quirks & Notes

- **Fire TV built-in** -- The Echo Show 15 runs Fire TV OS, making it the only Echo device that doubles as a streaming TV. The ADB (Android Debug Bridge) remote control protocol could theoretically be used for basic control, but requires developer mode, is undocumented by Amazon, and breaks with firmware updates.
- **Wall mount design** -- Unlike other Echo devices, the Show 15 is designed for wall mounting with a custom bracket. It can also be used on an optional tilt stand.
- **Visual ID** -- The camera uses on-device face recognition to personalize the display for different family members. This data stays on-device (Amazon's claim) but is not accessible via any API.
- **Matter/Thread (2nd Gen)** -- The 2nd Gen Echo Show 15 includes Zigbee, Thread, and Matter radios, like the 5th Gen Echo. Same cloud-managed limitation applies.
- **Not feasible for Haus** -- Same rationale as all Echo devices. Amazon's closed ecosystem does not provide local control APIs.

## Similar Devices

- **[amazon-echo-5th-gen](amazon-echo-5th-gen.md)** -- Speaker-only Echo with the same cloud-only limitation
- **[google-nest-hub](google-nest-hub.md)** -- Competing smart display with local Cast API (read-only)
- **[apple-tv-4k](apple-tv-4k.md)** -- Competing media platform with AirPlay and companion protocol
