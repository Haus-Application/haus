---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "google-nest-mini"
name: "Google Nest Mini (2nd Gen)"
manufacturer: "Google LLC"
brand: "Google Nest"
model: "GXCA6"
model_aliases: ["Nest Mini", "Google Home Mini", "H0A", "H2C", "GXCA6"]
device_type: "cast_speaker"
category: "media"
product_line: "Nest"
release_year: 2019
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "F4:F5:D8"        # Google, Inc.
    - "30:FD:38"        # Google, Inc.
    - "54:60:09"        # Google, Inc.
    - "A4:77:33"        # Google, Inc.
    - "48:D6:D5"        # Google, Inc.
    - "20:DF:B9"        # Google, Inc.
    - "F8:0F:F9"        # Google, Inc.
    - "E4:F0:42"        # Google, Inc.
    - "1C:F2:9A"        # Google, Inc.
    - "44:07:0B"        # Google, Inc.
  mdns_services:
    - "_googlecast._tcp"
  mdns_txt_keys:
    - "id"
    - "cd"
    - "rm"
    - "ve"
    - "md"              # "Google Nest Mini"
    - "ic"
    - "fn"              # friendly name
    - "ca"
    - "st"
    - "bs"
    - "nf"
    - "rs"
  default_ports: [8008, 8009, 8443]
  signature_ports: [8008, 8009]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Google-Nest-Mini"
    - "^Google-Home-Mini"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8008
    path: "/setup/eureka_info?params=name,build_info,detail,device_info,opt_in"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"model_name\""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "read_only"
  integration_key: "cast"
  polling_interval_sec: 30
  websocket_event: "cast:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities:
  - "media_playback"
  - "volume"

# --- PROTOCOL ---
protocol:
  type: "protobuf_tls"
  port: 8009
  transport: "TLS"
  encoding: "Protobuf"
  auth_method: "none"
  auth_detail: "Port 8008 HTTP info API requires no authentication. Port 8009 Cast control protocol uses TLS with device certificate."
  base_url_template: "http://{ip}:8008"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "speaker"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://store.google.com/us/product/google_nest_mini"
  api_docs: "https://developers.google.com/cast"
  developer_portal: "https://developers.google.com/cast"
  support: "https://support.google.com/googlenest/"
  community_forum: "https://support.google.com/googlenest/community"
  image_url: ""
  fcc_id: "A4RH2C"

# --- TAGS ---
tags: ["cast", "speaker", "google-assistant", "wifi", "smart-speaker"]
---

# Google Nest Mini (2nd Gen)

## What It Is

The Google Nest Mini (2nd Gen, model H2C) is a compact smart speaker manufactured by Google LLC. It is the successor to the Google Home Mini and features a 40mm driver, three far-field microphones, Google Assistant, and the Google Cast protocol. It is Google's entry-level smart speaker, designed primarily for voice control and casual audio playback. It connects via Wi-Fi and exposes the same read-only HTTP info API on port 8008 and Cast control protocol on port 8009 as all Google Cast devices. Unlike the Nest Hub (2nd Gen), the Nest Mini does not include a Thread radio and cannot serve as a Thread border router.

## How Haus Discovers It

Discovery follows the standard Google Cast flow. See **[google-nest-hub](google-nest-hub.md)** for the complete discovery process.

1. **OUI Match** -- MAC prefix from Google OUI pool. Not definitive alone (shared across all Google hardware).
2. **mDNS Discovery** -- `_googlecast._tcp.local.` with TXT record `md` = "Google Nest Mini" and `fn` containing the user-assigned friendly name.
3. **HTTP Fingerprint** -- `GET http://{ip}:8008/setup/eureka_info` returns `device_info.model_name` = "Google Nest Mini".
4. **Port Probe** -- Ports 8008 and 8009 are open.

The key differentiator from the Nest Hub is the `md` mDNS TXT record and the `model_name` in eureka_info. The Nest Mini does NOT advertise `_meshcop._udp` (no Thread radio).

## Pairing / Authentication

No authentication required for the HTTP info API on port 8008. See **[google-nest-hub](google-nest-hub.md)** for full details on the Cast protocol authentication model.

Initial device setup requires the Google Home app and a Google account.

## API Reference

The Nest Mini uses the identical Cast protocol as all Google Cast devices. See **[google-nest-hub](google-nest-hub.md)** for the complete API reference covering:

- HTTP info API on port 8008 (eureka_info, configured_networks, supported_app_ids)
- Cast control protocol on port 8009 (Protobuf/TLS)
- Cast namespaces (connection, heartbeat, receiver, media)
- Volume control, media playback, app launch

All endpoints and protocol details are identical. The only difference is the device capabilities -- the Nest Mini is audio-only (no display, no camera).

## AI Capabilities

When chatting with a Nest Mini via Haus, the AI can:
- **Report device identity** -- name, model, software version
- **Report network status** -- online/offline, responding on port 8008

Cast control (volume, playback) requires Cast SDK integration and is planned for a future milestone.

## Quirks & Notes

- **No Thread radio** -- Unlike the Nest Hub 2nd Gen, the Nest Mini does not support Thread and cannot serve as a Matter border router. This is the primary functional difference from the Nest Hub for Haus integration purposes.
- **Speaker groups** -- Nest Minis can be grouped with other Cast speakers for multi-room audio via the Google Home app. Group membership is visible in the Cast protocol but requires the cloud API to modify.
- **Ultrasonic presence** -- The Nest Mini uses ultrasonic sensing to detect nearby presence and illuminate touch-sensitive LEDs. This capability is not accessible via any network API.
- **Wall mount** -- The Nest Mini has a built-in wall mount slot, unlike the original Home Mini. This has no API implications but explains why users may have these in unexpected locations.
- **Same OUI pool** -- Google uses the same MAC prefix pool for all hardware. Cannot distinguish a Nest Mini from a Nest Hub or Chromecast by MAC alone.
- **All Cast quirks apply** -- See [google-nest-hub](google-nest-hub.md) for general Cast protocol notes.

## Similar Devices

- **[google-nest-hub](google-nest-hub.md)** -- Smart display with Cast protocol plus Thread border router
- **[chromecast-google-tv](chromecast-google-tv.md)** -- HDMI Cast dongle with the same protocol
- **[amazon-echo-5th-gen](amazon-echo-5th-gen.md)** -- Competing entry-level smart speaker (Alexa vs Google Assistant)
- **[apple-homepod-mini](apple-homepod-mini.md)** -- Competing smart speaker (Siri/HomeKit vs Google Assistant)
