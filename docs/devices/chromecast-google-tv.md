---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "chromecast-google-tv"
name: "Chromecast with Google TV"
manufacturer: "Google LLC"
brand: "Google"
model: "GZRNL"
model_aliases: ["Chromecast with Google TV", "Chromecast with Google TV (4K)", "Chromecast with Google TV (HD)", "GZRNL", "G9N9N", "GA01919-US", "GA03131-US"]
device_type: "cast_media_player"
category: "media"
product_line: "Chromecast"
release_year: 2020
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
  protocols_spoken: ["wifi", "bluetooth", "hdmi"]

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
    - "md"              # "Chromecast" or "Chromecast with Google TV"
    - "ic"
    - "fn"              # friendly name
    - "ca"              # capabilities bitmask
    - "st"
    - "bs"
    - "nf"
    - "rs"
  default_ports: [8008, 8009, 8443]
  signature_ports: [8008, 8009]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Chromecast"
    - "^[0-9a-f]{32}$"
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
  - "input_select"

# --- PROTOCOL ---
protocol:
  type: "protobuf_tls"
  port: 8009
  transport: "TLS"
  encoding: "Protobuf"
  auth_method: "none"
  auth_detail: "Port 8008 HTTP info API requires no authentication. Port 8009 Cast control protocol uses TLS with device certificate. Same as all Google Cast devices."
  base_url_template: "http://{ip}:8008"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "usb"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://store.google.com/us/product/chromecast_google_tv"
  api_docs: "https://developers.google.com/cast"
  developer_portal: "https://developers.google.com/cast"
  support: "https://support.google.com/chromecast/"
  community_forum: "https://support.google.com/chromecast/community"
  image_url: ""
  fcc_id: "A4RGZRNL"

# --- TAGS ---
tags: ["cast", "media-player", "google-tv", "hdmi", "wifi", "streaming-dongle", "4k-hdr"]
---

# Chromecast with Google TV

## What It Is

The Chromecast with Google TV is an HDMI streaming dongle manufactured by Google LLC. It plugs into a TV's HDMI port and runs Google TV (Android TV-based), providing a full app-based streaming experience with a dedicated remote control. Available in 4K (model GZRNL, released 2020) and HD (model G9N9N, released 2022) variants, it supports the Google Cast protocol for receiving streamed content from phones and computers, as well as standalone app navigation via the included Bluetooth remote. It connects to the home network via Wi-Fi and exposes the standard Google Cast HTTP info API on port 8008 and Cast control protocol on port 8009. Unlike the original Chromecast (cast-only dongle), the Google TV version runs full Android TV apps independently.

## How Haus Discovers It

Discovery follows the standard Google Cast flow. See **[google-nest-hub](google-nest-hub.md)** for the complete discovery process.

1. **OUI Match** -- Google MAC prefix. Not definitive alone.
2. **mDNS Discovery** -- `_googlecast._tcp.local.` with TXT record `md` containing "Chromecast" and `fn` containing the user-assigned friendly name.
3. **HTTP Fingerprint** -- `GET http://{ip}:8008/setup/eureka_info` returns `device_info.model_name` = "Chromecast" or "Chromecast with Google TV".
4. **Port Probe** -- Ports 8008 (HTTP info) and 8009 (Cast TLS) are open.

### Distinguishing from Other Cast Devices

The `md` (model description) mDNS TXT record and `device_info.model_name` in eureka_info differentiate the Chromecast from Nest Hub, Nest Mini, and third-party Cast devices:
- `md` = "Chromecast" -- Chromecast with Google TV (4K)
- `md` = "Chromecast HD" -- Chromecast with Google TV (HD)
- `md` = "Google Nest Hub" -- Nest Hub
- `md` = "Google Nest Mini" -- Nest Mini

The Chromecast with Google TV does NOT advertise `_meshcop._udp` (no Thread radio) or `_homekit._tcp`.

## Pairing / Authentication

No authentication required for the HTTP info API on port 8008. See **[google-nest-hub](google-nest-hub.md)** for full Cast protocol details.

Initial device setup requires the Google Home app and a Google account.

## API Reference

The Chromecast with Google TV uses the identical Cast protocol as all Google Cast devices. See **[google-nest-hub](google-nest-hub.md)** for the complete API reference covering:

- HTTP info API on port 8008 (eureka_info, configured_networks, supported_app_ids)
- Cast control protocol on port 8009 (Protobuf/TLS)
- Cast namespaces (connection, heartbeat, receiver, media)
- Volume control, media playback, app launch

### Google TV / Android TV ADB

Since the Chromecast with Google TV runs Android TV, it potentially supports ADB (Android Debug Bridge) when developer mode is enabled:

- **Port 5555** -- ADB over network (must be manually enabled in Settings > System > Developer options)
- **Capabilities** -- App install/launch, input injection, screen capture, shell commands
- **Not default** -- ADB is disabled by default and requires the user to enable developer options

This is not a reliable integration path (requires user action to enable, security implications) but is documented for completeness.

## AI Capabilities

When chatting with a Chromecast via Haus, the AI can:
- **Report device identity** -- name, model, software version from eureka_info
- **Report network status** -- online/offline, responding on port 8008

Cast control (play/pause/volume) requires Cast SDK integration and is planned for a future milestone.

## Quirks & Notes

- **Google TV vs Chrome OS** -- Despite the "Chromecast" name, the Google TV version runs a full Android TV OS, not Chrome OS. It has a Play Store, runs apps independently, and has a dedicated remote. The Cast protocol is just one input method.
- **No Thread radio** -- Unlike the Nest Hub 2nd Gen, the Chromecast with Google TV does not include a Thread radio and cannot serve as a Matter border router.
- **HDMI-CEC** -- The Chromecast supports HDMI-CEC for TV power control and input switching. When CEC is active, powering on the Chromecast can wake the TV and switch to its HDMI input.
- **USB-C power** -- Powered via USB-C. Some TVs provide enough power through their USB port, but Google recommends the included power adapter for reliable operation.
- **Ambient mode** -- When idle, the Chromecast displays ambient content (photos, art, weather). This state is visible via the Cast protocol as the "backdrop" app.
- **Same Cast protocol** -- From a network perspective, the Chromecast with Google TV is identical to the original Chromecast, Nest Hub, and any other Cast device. The eureka_info and Cast control protocol are the same.
- **All Cast quirks apply** -- See [google-nest-hub](google-nest-hub.md) for general Cast protocol notes (Google OUI overlap, firmware auto-updates, Cast protocol complexity).

## Similar Devices

- **[google-nest-hub](google-nest-hub.md)** -- Smart display with Cast protocol plus Thread border router
- **[google-nest-mini](google-nest-mini.md)** -- Speaker-only Cast device
- **[apple-tv-4k](apple-tv-4k.md)** -- Competing HDMI media streamer with AirPlay and companion protocol
- **[amazon-echo-show-15](amazon-echo-show-15.md)** -- Competing media platform with Fire TV (cloud-only)
