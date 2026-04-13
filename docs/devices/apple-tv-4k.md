---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "apple-tv-4k"
name: "Apple TV 4K (3rd Gen)"
manufacturer: "Apple Inc."
brand: "Apple"
model: "A2843"
model_aliases: ["Apple TV 4K", "A2737", "A2843", "A2169", "AppleTV11,1", "AppleTV14,1"]
device_type: "airplay_media_player"
category: "media"
product_line: "Apple TV"
release_year: 2022
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["wifi", "ethernet", "bluetooth", "thread", "hdmi"]

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
    - "_appletv-v2._tcp"
    - "_airplay._tcp"
    - "_raop._tcp"
    - "_companion-link._tcp"
    - "_homekit._tcp"
    - "_hap._tcp"
    - "_meshcop._udp"       # Thread border router
    - "_touch-able._tcp"    # Remote app protocol
    - "_mediaremotetv._tcp" # Media Remote protocol
  mdns_txt_keys:
    - "model"           # "AppleTV14,1" for Apple TV 4K 3rd Gen
    - "features"        # AirPlay feature bitmask
    - "flags"
    - "pk"              # public key
    - "pi"              # pairing identity
    - "srcvers"
    - "vv"
    - "sf"
    - "UniqueIdentifier"
    - "Name"            # user-assigned device name
  default_ports: [3689, 7000, 7100, 49152]
  signature_ports: [7000, 3689]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Apple-TV"
    - ".*AppleTV.*"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "apple_tv"
  polling_interval_sec: 10
  websocket_event: "appletv:state"
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "M7"

# --- CAPABILITIES ---
capabilities:
  - "media_playback"
  - "volume"
  - "input_select"

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 49152
  transport: "TCP"
  encoding: "Protobuf"
  auth_method: "none"
  auth_detail: "Companion Link Protocol uses SRP pairing (PIN displayed on TV, entered on controller). After initial pairing, subsequent connections use stored credentials. The open-source pyatv library implements this protocol."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le", "thread"]

# --- LINKS ---
links:
  product_page: "https://www.apple.com/apple-tv-4k/"
  api_docs: ""
  developer_portal: "https://developer.apple.com/tvos/"
  support: "https://support.apple.com/apple-tv"
  community_forum: "https://discussions.apple.com/community/apple-tv"
  image_url: ""
  fcc_id: "BCG-E3528"

# --- TAGS ---
tags: ["airplay2", "homekit-hub", "thread-border-router", "matter-controller", "media-player", "apple", "companion-link", "pyatv", "hdmi"]
---

# Apple TV 4K (3rd Gen)

## What It Is

The Apple TV 4K (3rd Gen, 2022) is a premium streaming media player manufactured by Apple Inc. It connects to a TV via HDMI and runs tvOS, providing access to streaming apps (Apple TV+, Netflix, Disney+, etc.), the Apple Arcade gaming service, and AirPlay 2 audio/video receiving. Beyond media, the Apple TV 4K serves as a HomeKit/Matter smart home hub and Thread border router, making it a critical infrastructure device for Apple's smart home ecosystem. It features an A15 Bionic chip, 64GB or 128GB storage, Wi-Fi 6, Gigabit Ethernet (128GB model), Bluetooth 5.0, Thread, and HDMI 2.1 with 4K HDR/Dolby Vision/Dolby Atmos output. Critically for Haus, the Apple TV exposes several mDNS services and can be controlled locally via the Companion Link Protocol and the Media Remote Protocol, both of which have been reverse-engineered by the open-source pyatv project.

## How Haus Discovers It

1. **OUI Match** -- MAC address begins with an Apple OUI prefix. Not definitive alone due to Apple's large OUI pool.
2. **mDNS Discovery** -- The Apple TV advertises a rich set of mDNS services:
   - `_appletv-v2._tcp` -- Primary Apple TV identification service. This is the definitive service for Apple TV devices.
   - `_airplay._tcp` -- AirPlay 2 receiver. TXT record `model` = "AppleTV14,1" (3rd Gen 4K).
   - `_mediaremotetv._tcp` -- Media Remote Protocol for remote control.
   - `_companion-link._tcp` -- Companion Link Protocol for app-to-device communication.
   - `_touch-able._tcp` -- Legacy Remote app protocol.
   - `_homekit._tcp` / `_hap._tcp` -- HomeKit hub advertisement.
   - `_meshcop._udp` -- Thread border router.
3. **Model Identification** -- The `model` TXT key in `_airplay._tcp` definitively identifies the generation:
   - `AppleTV6,2` = Apple TV 4K (1st Gen, 2017)
   - `AppleTV11,1` = Apple TV 4K (2nd Gen, 2021)
   - `AppleTV14,1` = Apple TV 4K (3rd Gen, 2022, Wi-Fi + Ethernet)
   - `AppleTV14,2` = Apple TV 4K (3rd Gen, 2022, Wi-Fi only)
4. **Port Probe** -- Port 7000 (AirPlay), 3689 (DAAP/companion), and dynamic ports (49152+) for HAP and Media Remote.

## Pairing / Authentication

The Apple TV supports pairing for the Companion Link and Media Remote protocols:

### Companion Link Pairing (pyatv method)

1. Haus initiates a pairing request via the Companion Link Protocol.
2. The Apple TV displays a **4-digit PIN** on the connected TV screen.
3. The user enters this PIN in the Haus interface.
4. SRP (Secure Remote Password) authentication is performed using the PIN.
5. On success, Haus receives a set of credentials (device identifier + authentication key) that are stored for future connections.
6. Subsequent connections use the stored credentials without requiring a PIN.

### AirPlay Pairing

AirPlay uses a similar PIN-based pairing flow when the Apple TV is configured to require verification for AirPlay connections. The pairing credentials are separate from Companion Link credentials.

### MRP (Media Remote Protocol) Pairing

The Media Remote Protocol uses its own pairing flow, also PIN-based. The pyatv library handles all three pairing protocols.

## API Reference

There is no official Apple-documented local API. However, the open-source **pyatv** project (https://github.com/postlund/pyatv) has reverse-engineered the key protocols, and their Go equivalents could be implemented for Haus.

### Companion Link Protocol

The primary control protocol for modern Apple TV (tvOS 15+). Uses Protobuf-encoded messages over an encrypted TCP connection.

**Capabilities:**
- Get current playing app and media information
- Play/Pause/Stop media
- Navigate menus (up, down, left, right, select, menu, home)
- Volume control (when Apple TV controls audio output)
- Launch apps by bundle ID
- Get device state (on/off, idle/playing)
- Power on/off

**Key protocol details:**
- Port: Dynamic (advertised via `_companion-link._tcp` mDNS)
- Encryption: SRP-based session encryption after pairing
- Encoding: Protocol Buffers (Protobuf)
- Keep-alive: Periodic heartbeat messages required

### Media Remote Protocol (MRP)

An older protocol (tvOS 12-14 era) for remote control. Still functional but being superseded by Companion Link.

**Capabilities:**
- Transport control (play, pause, skip, previous, seek)
- Now playing information (title, artist, album, artwork)
- Playback state and position
- Volume control
- Navigation (arrow keys, select, menu)

**Key protocol details:**
- Port: Dynamic (advertised via `_mediaremotetv._tcp` mDNS)
- Encryption: SRP-based session encryption
- Encoding: Protocol Buffers

### AirPlay 2

See **[apple-homepod-mini](apple-homepod-mini.md)** for AirPlay protocol details. The Apple TV functions as an AirPlay 2 receiver for both audio and video.

### pyatv Library Reference

The pyatv Python library (https://pyatv.dev/) provides a complete implementation of all Apple TV protocols:

```
pip install pyatv
atvremote scan                    # Discover Apple TVs
atvremote -s 192.168.1.x pair    # Pair with device
atvremote -s 192.168.1.x playing # Get now playing info
atvremote -s 192.168.1.x pause   # Pause playback
```

For Haus (Go), the pyatv protocol documentation and source code serve as the reference implementation. A Go port of the core Companion Link and MRP protocols would be needed.

## AI Capabilities

AI integration is planned for a future milestone. When implemented, the AI concierge will be able to:
- Report current playback state (app, media title, artist, position)
- Control playback (play, pause, skip, previous, seek)
- Adjust volume
- Report device power state
- Launch apps by name
- Navigate the Apple TV UI via remote control commands

## Quirks & Notes

- **pyatv is the Rosetta Stone** -- The pyatv project by Pierre Stitchlebaut is the most comprehensive reverse-engineering of Apple TV protocols. It implements Companion Link, MRP, AirPlay, DMAP, and pairing for all of them. Any Haus Apple TV integration should reference pyatv extensively.
- **Thread border router** -- The Apple TV 4K (2nd Gen and later) is a Thread border router, joining the HomePod mini as Apple's Thread infrastructure. Having multiple Thread border routers improves Matter device reliability.
- **Matter controller** -- The Apple TV 4K (2nd Gen+) is also a Matter controller and can commission Matter devices through the Apple Home app. This Matter fabric is managed by Apple Home, not accessible to third parties.
- **HDMI CEC** -- The Apple TV supports HDMI-CEC for TV power control and volume passthrough. When HDMI-CEC is active, volume commands sent via Companion Link control the TV/soundbar volume.
- **Multiple protocol eras** -- Apple has iterated through several remote control protocols: DAAP/DMAP (legacy), MRP (tvOS 12-14), Companion Link (tvOS 15+). pyatv supports all of them and auto-negotiates.
- **Sleep mode** -- The Apple TV enters sleep mode when idle. It remains on the network and responds to mDNS/Wake-on-LAN but requires a wake command before control. pyatv handles this transparently.
- **Apple OUI overlap** -- Same limitation as HomePod: Apple uses many OUI blocks across all products. The `_appletv-v2._tcp` mDNS service is the definitive identifier.
- **Go implementation required** -- pyatv is Python. Haus will need a Go implementation of the Companion Link Protocol. The Protobuf message definitions can be extracted from pyatv's source code. The SRP pairing library exists in Go (`github.com/tadglines/go-pkgs/crypto/srp`).

## Similar Devices

- **[apple-homepod-mini](apple-homepod-mini.md)** -- Apple smart speaker, shares HomeKit hub and Thread border router roles
- **[chromecast-google-tv](chromecast-google-tv.md)** -- Competing HDMI media streamer with Cast protocol
- **[amazon-echo-show-15](amazon-echo-show-15.md)** -- Competing media platform with Fire TV (cloud-only)
- **[sonos-beam-arc](sonos-beam-arc.md)** -- Often receives audio from Apple TV via HDMI eARC or AirPlay
