---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "amazon-echo-5th-gen"
name: "Amazon Echo (5th Gen)"
manufacturer: "Amazon.com, Inc."
brand: "Amazon Echo"
model: "B09ZX7MS5B"
model_aliases: ["Echo 5th Gen", "Echo Dot 5th Gen", "Echo 2022", "C2N6L4"]
device_type: "echo_speaker"
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
  default_ports: [8008, 55443]
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^amazon-"
    - "^echo-"
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

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "HTTPS"
  encoding: "Protobuf"
  auth_method: "oauth2"
  auth_detail: "All control flows through Amazon's cloud Alexa Voice Service (AVS). No documented or stable local API exists. The device communicates with Amazon servers over HTTPS/HTTP2 with proprietary Protobuf payloads."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "speaker"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le", "zigbee", "thread"]

# --- LINKS ---
links:
  product_page: "https://www.amazon.com/dp/B09ZX7MS5B"
  api_docs: "https://developer.amazon.com/en-US/docs/alexa/alexa-voice-service/api-overview.html"
  developer_portal: "https://developer.amazon.com/alexa"
  support: "https://www.amazon.com/gp/help/customer/display.html?nodeId=201399130"
  community_forum: "https://www.amazonforum.com/s/topic/0TO4P000000gSRHWA2/echo-alexa"
  image_url: ""
  fcc_id: "2AETW-1677"

# --- TAGS ---
tags: ["alexa", "cloud-only", "zigbee-hub", "thread-border-router", "matter-controller", "voice-assistant", "no-local-api"]
---

# Amazon Echo (5th Gen)

## What It Is

The Amazon Echo (5th Gen, released 2022) is a spherical smart speaker with Alexa voice assistant, manufactured by Amazon. It features a 3.0" woofer, dual 0.8" tweeters, and a built-in smart home hub with Zigbee, Thread, and Matter radios. The 5th Gen Echo is also an Eero mesh Wi-Fi extender and a Thread border router for Matter devices. Despite its extensive hardware capabilities, the Echo is fundamentally a cloud-dependent device -- all voice processing, smart home control, and media playback routing flows through Amazon's Alexa Voice Service (AVS) cloud. There is no documented, stable, or supported local API for controlling the Echo or querying its state from the local network.

## How Haus Discovers It

1. **OUI Match** -- MAC address begins with an Amazon Technologies prefix (`FC:65:DE`, `A0:02:DC`, `68:54:FD`, `40:B4:CD`, `74:C2:46`, `F0:F0:A4`, `44:00:49`, `38:F7:3D`, `B0:FC:0D`, `AC:63:BE`, `14:91:82`, `18:74:2E`). Amazon uses many OUI blocks across Echo, Fire TV, Kindle, Ring, and Eero products, so OUI alone is not definitive.
2. **Hostname Pattern** -- Echo devices often appear in DHCP with hostnames starting with `amazon-` or containing the device's serial number.
3. **Port Scan** -- Some Echo devices respond on port 8008 (limited HTTP endpoint) and port 55443, but these do not expose useful control APIs.
4. **Classification** -- Without mDNS advertisement or a reliable HTTP fingerprint, Haus classifies Amazon OUI devices as "detected_only" and prompts the user to identify the device type manually if needed.

Haus will detect the Echo on the network but cannot integrate with it for control or state monitoring.

## Pairing / Authentication

Not applicable for Haus integration. The Echo is set up exclusively through the Amazon Alexa app (iOS/Android) with an Amazon account. All control and configuration flows through Amazon's cloud.

## API Reference

There is no usable local API for the Amazon Echo.

### Why Local Control Is Not Feasible

- **Alexa Voice Service (AVS)** -- All voice commands are streamed to Amazon's cloud for processing. Responses and actions are returned from the cloud.
- **Alexa Smart Home API** -- Smart home device control (lights, thermostats, etc.) is routed through Amazon's cloud even when the Echo is controlling local Zigbee/Thread devices.
- **No local REST/SOAP/UPnP API** -- Unlike Sonos (port 1400 SOAP) or Google Cast (port 8008 HTTP), Amazon does not expose a local API for media control, volume, or device state queries.
- **Undocumented endpoints** -- Community reverse-engineering has found some HTTP endpoints on Echo devices, but they are undocumented, change frequently with firmware updates, and require Amazon session cookies that expire.
- **Alexa API (cloud)** -- The Alexa Smart Home Skill API and Alexa Voice Service API exist for cloud-to-cloud integrations but require an Amazon Developer account, OAuth2, and persistent cloud connectivity. This is architecturally incompatible with Haus's local-first design.

### Matter Controller Role

The 5th Gen Echo functions as a Matter controller and Thread border router. It can commission and control Matter devices via the Alexa app. However, this Matter controller functionality is managed entirely through Amazon's cloud -- Haus cannot interact with the Echo's Matter fabric locally.

## AI Capabilities

Not applicable. The Echo cannot be controlled or queried locally, so the AI concierge cannot interact with it.

Haus will display the Echo as a detected network device with a note that it is cloud-dependent and not locally controllable.

## Quirks & Notes

- **Cloud dependency is absolute** -- If internet connectivity is lost, the Echo loses virtually all functionality including voice control, music playback, and smart home control. Some basic Zigbee device control may continue via local processing, but this is limited and undocumented.
- **Matter/Thread hardware** -- The 5th Gen Echo has Zigbee, Thread, and Matter radios. It can serve as a Thread border router and Matter controller. However, these capabilities are managed through Amazon's cloud, not locally.
- **Eero mesh extension** -- The 5th Gen Echo contains Eero mesh Wi-Fi technology and can extend an Eero mesh network. This is configured through the Eero app, not the Alexa app.
- **Amazon OUI sprawl** -- Amazon uses many MAC prefix blocks across Echo, Fire TV, Kindle, Ring, Blink, and Eero products. Distinguishing an Echo from a Fire TV Stick or Ring doorbell by MAC alone is not reliable.
- **Not feasible for Haus** -- Integration status is "not_feasible" because Amazon's closed ecosystem does not provide any stable local API. The cloud-only Alexa API requires persistent internet and Amazon account linking, which conflicts with Haus's local-first, privacy-focused design.
- **Competing ecosystem** -- Users with Echo devices likely also have Alexa routines and integrations. Haus should acknowledge the Echo's presence on the network gracefully and explain why it cannot be integrated.

## Similar Devices

- **[amazon-echo-show-15](amazon-echo-show-15.md)** -- Amazon's smart display with the same cloud-only limitation
- **[google-nest-mini](google-nest-mini.md)** -- Competing smart speaker with local Cast API (read-only)
- **[apple-homepod-mini](apple-homepod-mini.md)** -- Competing smart speaker with AirPlay (limited local control)
- **[sonos-era-100](sonos-era-100.md)** -- Smart speaker with full local API (UPnP/SOAP)
