---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "apple-homepod-mini"
name: "Apple HomePod mini"
manufacturer: "Apple Inc."
brand: "Apple"
model: "MY5G2LL/A"
model_aliases: ["HomePod mini", "A2374"]
device_type: "airplay_speaker"
category: "media"
product_line: "HomePod"
release_year: 2020
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["wifi", "bluetooth", "thread"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "F0:B3:EC"        # Apple, Inc.
    - "28:6A:BA"        # Apple, Inc.
    - "3C:E0:72"        # Apple, Inc.
    - "AC:BC:32"        # Apple, Inc.
    - "DC:A4:CA"        # Apple, Inc.
    - "A8:8C:3E"        # Apple, Inc.
    - "F0:D4:F6"        # Apple, Inc.
    - "8C:85:90"        # Apple, Inc.
    - "3C:06:30"        # Apple, Inc.
    - "14:98:77"        # Apple, Inc.
    - "70:3C:69"        # Apple, Inc.
    - "E0:5F:45"        # Apple, Inc.
  mdns_services:
    - "_airplay._tcp"
    - "_raop._tcp"          # Remote Audio Output Protocol (AirPlay audio)
    - "_companion-link._tcp"
    - "_homekit._tcp"
    - "_hap._tcp"           # HomeKit Accessory Protocol
    - "_meshcop._udp"       # Thread border router
  mdns_txt_keys:
    - "model"           # "AudioAccessory5,1" for HomePod mini
    - "features"        # AirPlay feature bitmask
    - "flags"           # AirPlay flags
    - "pk"              # public key
    - "pi"              # pairing identity
    - "srcvers"         # AirPlay source version
    - "vv"              # protocol version
    - "sf"              # status flags
    - "ci"              # category identifier
  default_ports: [7000, 7100, 49152]
  signature_ports: [7000]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^HomePod"
    - ".*HomePod.*"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "detected_only"
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
  port: 7000
  transport: "TCP"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "AirPlay 2 uses encrypted communication with device pairing. No documented public API for local control exists. The HomePod is controlled through the Apple Home app, Siri, or AirPlay streaming from Apple devices."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "speaker"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le", "thread"]

# --- LINKS ---
links:
  product_page: "https://www.apple.com/homepod-mini/"
  api_docs: ""
  developer_portal: "https://developer.apple.com/homekit/"
  support: "https://support.apple.com/homepod"
  community_forum: "https://discussions.apple.com/community/homepod"
  image_url: ""
  fcc_id: "BCG-E3526"

# --- TAGS ---
tags: ["airplay2", "homekit-hub", "thread-border-router", "matter", "siri", "smart-speaker", "apple", "no-local-api"]
---

# Apple HomePod mini

## What It Is

The Apple HomePod mini is a compact smart speaker manufactured by Apple Inc. It features a full-range driver, two passive radiators, Siri voice assistant, and serves as a HomeKit/Matter smart home hub and Thread border router. It is the smallest and most affordable speaker in Apple's lineup and the most common Apple smart home hub device. The HomePod mini supports AirPlay 2 for audio streaming from Apple devices and multi-room audio with other AirPlay 2 speakers. It connects via Wi-Fi 4 (802.11n) and includes Bluetooth 5.0, Thread, and Ultra Wideband (U1) radios. From a local control perspective, the HomePod mini is largely opaque -- it advertises several mDNS services but does not expose a documented public API for third-party local control.

## How Haus Discovers It

1. **OUI Match** -- MAC address begins with an Apple OUI prefix (`F0:B3:EC`, `28:6A:BA`, `3C:E0:72`, `AC:BC:32`, `DC:A4:CA`, `A8:8C:3E`, etc.). Apple uses hundreds of OUI blocks across iPhones, iPads, Macs, Apple TVs, HomePods, and accessories, so OUI alone is not definitive.
2. **mDNS Discovery** -- The HomePod mini advertises multiple mDNS services:
   - `_airplay._tcp` -- AirPlay 2 receiver. TXT record `model` = "AudioAccessory5,1" identifies it as a HomePod mini.
   - `_raop._tcp` -- Remote Audio Output Protocol (AirPlay audio).
   - `_companion-link._tcp` -- Apple companion link protocol for device pairing.
   - `_homekit._tcp` / `_hap._tcp` -- HomeKit Accessory Protocol advertisement.
   - `_meshcop._udp` -- Thread border router (Thread Mesh Commissioning Protocol).
3. **Model Identification** -- The `model` key in the AirPlay mDNS TXT record definitively identifies the device: `AudioAccessory5,1` = HomePod mini, `AudioAccessory6,1` = HomePod (2nd Gen).
4. **Hostname** -- Typically appears as "HomePod-{hex}" or the user-assigned name with spaces replaced.

## Pairing / Authentication

There is no pairing or authentication flow that Haus can initiate for controlling the HomePod mini.

The device is set up through the Apple Home app on an iPhone/iPad. It requires an Apple ID and iCloud account. AirPlay streaming is initiated from Apple devices (iPhone, iPad, Mac) and uses encrypted peer-to-peer pairing.

The HomeKit Accessory Protocol (HAP) that the HomePod speaks for smart home hub functionality requires Apple's MFi (Made for iPhone) program licensing for certified accessory development. It is not available for third-party control apps to use for controlling the HomePod itself.

## API Reference

There is no documented public local API for controlling the HomePod mini.

### AirPlay 2 Protocol

AirPlay 2 is Apple's proprietary wireless streaming protocol. Key technical details:

- **Port 7000** -- Primary AirPlay service port (RTSP-based control)
- **Port 7100** -- AirPlay event/callback port
- **Protocol** -- Based on RTSP (Real-Time Streaming Protocol) with Apple-specific extensions
- **Encryption** -- Uses FairPlay encryption for DRM-protected content and pair-verify for authentication
- **Codec** -- ALAC (Apple Lossless) or AAC for audio
- **Buffered streaming** -- AirPlay 2 uses buffered (non-real-time) streaming for multi-room synchronization, unlike AirPlay 1 which was real-time

### Open-Source AirPlay Implementations

Community projects have reverse-engineered portions of the AirPlay protocol:
- **shairport-sync** -- Open-source AirPlay audio receiver (receives audio, cannot send to HomePod)
- **RPiPlay** -- AirPlay mirroring receiver for Raspberry Pi

These projects demonstrate that the AirPlay protocol can be reverse-engineered, but sending audio TO a HomePod requires Apple's FairPlay DRM implementation, which is not publicly available.

### HomeKit / HAP

The HomePod advertises `_hap._tcp` as a HomeKit accessory bridge. Through HomeKit, the HomePod controls other smart home devices but does not expose itself as a controllable device. The HAP protocol uses:

- **Port 49152+** -- Dynamic port for HAP connections
- **Encryption** -- SRP (Secure Remote Password) for pairing, ChaCha20-Poly1305 for session encryption
- **Data format** -- TLV8 (Type-Length-Value 8-bit) encoding

### Thread Border Router

The HomePod mini advertises `_meshcop._udp` as a Thread border router. It can commission and communicate with Thread devices (including Matter-over-Thread) within the Apple Home ecosystem. This Thread functionality is managed through Apple's Home app and is not accessible to third-party controllers.

## AI Capabilities

Not applicable for local control. The HomePod mini is detected on the network but cannot be queried or controlled by Haus.

Haus will display the HomePod mini as a detected device and note its roles (AirPlay 2 receiver, HomeKit hub, Thread border router) based on mDNS service advertisements.

## Quirks & Notes

- **Apple ecosystem lock-in** -- The HomePod mini is deeply integrated into the Apple ecosystem. It requires an iPhone/iPad for setup, an Apple ID for operation, and works best with Apple devices for audio streaming. There is no path to third-party local control.
- **Thread border router** -- The HomePod mini is one of the most common Thread border routers in homes, alongside the Apple TV 4K and Google Nest Hub 2nd Gen. While Haus cannot control the HomePod, its presence as a Thread border router benefits Matter devices on the network.
- **AirPlay model identification** -- The `model` TXT record in `_airplay._tcp` mDNS is the definitive identifier:
  - `AudioAccessory1,1` = HomePod (1st Gen)
  - `AudioAccessory5,1` = HomePod mini
  - `AudioAccessory6,1` = HomePod (2nd Gen, 2023)
- **Apple OUI sprawl** -- Apple uses more MAC prefix blocks than almost any other manufacturer. iPhones, iPads, Macs, Apple Watches, AirPods, Apple TVs, and HomePods all use Apple OUIs. Additionally, Apple devices use MAC address randomization for Wi-Fi scanning, though the HomePod uses its real MAC when connected.
- **Ultra Wideband (U1)** -- The HomePod mini includes a U1 chip for spatial awareness (hand-off audio to/from iPhone). Not accessible via network APIs.
- **Temperature and humidity sensor** -- The HomePod mini includes a temperature and humidity sensor (added via software update in 2023). Data is accessible through the Apple Home app but not via any local API.
- **Detected only** -- Integration status is "detected_only" because while Haus can reliably identify the HomePod mini on the network via mDNS, there is no feasible path to control or query it locally. Unlike the Echo (not_feasible), the HomePod at least advertises itself, making detection reliable.

## Similar Devices

- **[apple-tv-4k](apple-tv-4k.md)** -- Apple's media streamer with AirPlay, HomeKit hub, and more promising integration via companion link protocol
- **[sonos-era-100](sonos-era-100.md)** -- Competing smart speaker with full local API
- **[google-nest-mini](google-nest-mini.md)** -- Competing smart speaker with local Cast API (read-only)
- **[amazon-echo-5th-gen](amazon-echo-5th-gen.md)** -- Competing smart speaker (cloud-only)
