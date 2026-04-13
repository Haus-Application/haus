---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "blink-outdoor"
name: "Blink Outdoor Camera"
manufacturer: "Blink (Amazon)"
brand: "Blink"
model: "Blink Outdoor 4th Gen"
model_aliases: ["Blink Outdoor 4", "Blink Outdoor (3rd Gen)", "B0B1N4LM4J", "Blink XT2"]
device_type: "blink_camera"
category: "security"
product_line: "Blink"
release_year: 2023
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
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["18:74:2E", "F0:F0:A4", "FC:65:DE", "A0:02:DC"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^blink.*", "^Blink.*"]
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
capabilities: ["camera_stream", "camera_snapshot", "motion", "battery_level"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Blink uses a proprietary cloud API with no public documentation. Authentication goes through Amazon/Blink servers. All video processing is cloud-side. No local API of any kind exists."
  base_url_template: "https://rest-{region}.immedia-semi.com"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "battery"
  mounting: "wall"
  indoor_outdoor: "outdoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://blinkforhome.com/products/blink-outdoor-4"
  api_docs: ""
  developer_portal: ""
  support: "https://support.blinkforhome.com"
  community_forum: ""
  image_url: ""
  fcc_id: "2AQTL-BLK73"

# --- TAGS ---
tags: ["cloud-only", "amazon", "no-local-api", "battery", "long-battery-life", "not-feasible", "motion-detection", "1080p", "outdoor", "sync-module"]
---

# Blink Outdoor Camera

## What It Is

> The Blink Outdoor is an affordable, battery-powered wireless security camera from Amazon's Blink subsidiary. It runs on two AA lithium batteries with an advertised 2-year battery life, provides 1080p video, infrared night vision, motion detection, and two-way audio. It requires a Blink Sync Module (a small hub) connected to the router via Ethernet or WiFi. The camera is completely cloud-locked with no local API, no RTSP, no ONVIF, and no known path to local integration. Integration status is "not_feasible."

## How Haus Discovers It

1. **OUI match** -- Blink Sync Module uses MAC prefixes including `18:74:2E`, `F0:F0:A4`, `FC:65:DE`, and `A0:02:DC` (Immedia Semiconductor / Amazon OUIs). Note: the Sync Module appears on the network, not the camera itself (cameras connect to the Sync Module via a proprietary 900MHz protocol, not WiFi)
2. **Hostname pattern** -- Sync Module may register as `blink-*` on the local network
3. **No port probe** -- Sync Module has no open local ports
4. **Identification only** -- Haus can detect the Blink Sync Module on the network and display it as a recognized but unsupported device

## Pairing / Authentication

> Blink devices can only be set up through the Blink app (iOS/Android):
>
> 1. Create Blink/Amazon account
> 2. Plug in Sync Module and add it via app
> 3. Scan camera QR code to pair to Sync Module
> 4. Camera communicates with Sync Module over proprietary 900MHz wireless protocol
>
> There is no way to pair or authenticate with Blink devices outside of the official app.

## API Reference

> **No usable API.** Blink uses Immedia Semiconductor's proprietary cloud infrastructure. While an unofficial Python library (blinkpy) has reverse-engineered some API endpoints, the protocol is designed to prevent third-party access:
>
> - Authentication requires Amazon account credentials with 2FA
> - API endpoints are region-specific (e.g., `rest-prod.immedia-semi.com`)
> - Video clips are stored in Amazon's cloud and delivered via signed URLs that expire quickly
> - Live view sessions use a proprietary protocol routed through Blink's relay servers
> - Camera-to-Sync-Module communication uses a proprietary 900MHz protocol (not WiFi), making local interception impossible
>
> **This device has been classified as not_feasible for Haus integration.**

## AI Capabilities

> No AI integration is possible due to the complete lack of local or documented cloud API access.

## Quirks & Notes

- **Not feasible for integration** -- This is one of the most locked-down consumer cameras on the market. No RTSP, no ONVIF, no local API, no documented cloud API, and the camera itself does not even use WiFi
- **Proprietary 900MHz wireless** -- Cameras connect to the Sync Module via a proprietary 900MHz radio link, not WiFi. The cameras have WiFi radios but only use them during initial setup
- **Sync Module required** -- The Sync Module acts as a bridge between the cameras and the internet. It is the only device that appears on the home network
- **2-year battery life** -- Achieved by extreme power management; camera is in deep sleep except when triggered by the PIR motion sensor
- **Local storage option** -- Blink Sync Module 2 supports a USB drive for local clip storage, but clips are still proprietary format and require the Blink app to view
- **Blink Subscription** -- Cloud storage requires Blink Subscription Plus ($9.99/mo for all cameras) or Basic ($2.99/mo per camera)
- **No continuous recording** -- Battery conservation means record-on-motion only; no 24/7 recording option

## Similar Devices

> - [Ring Indoor Cam](ring-indoor-cam.md) -- Amazon ecosystem, cloud-only but at least uses WiFi for video
> - [Ring Stick Up Cam](ring-stick-up-cam.md) -- Amazon's other battery camera option
> - [Arlo Pro 5](arlo-pro-5.md) -- Battery camera with possible RTSP via SmartHub
> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- Battery camera with native RTSP/ONVIF (the polar opposite of Blink)
