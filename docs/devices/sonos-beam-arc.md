---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "sonos-beam-arc"
name: "Sonos Beam / Arc Soundbar"
manufacturer: "Sonos, Inc."
brand: "Sonos"
model: "S14"
model_aliases: ["S14", "Sonos Beam", "Sonos Beam Gen 2", "S36-BEAM", "S19", "Sonos Arc", "Sonos Arc Ultra", "S23"]
device_type: "sonos_soundbar"
category: "media"
product_line: "Sonos"
release_year: 2020
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "ethernet", "bluetooth", "hdmi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "00:0E:58"        # Sonos, Inc. (primary OUI)
    - "5C:AA:FD"        # Sonos, Inc.
    - "54:2A:1B"        # Sonos, Inc.
    - "78:28:CA"        # Sonos, Inc.
    - "48:A6:B8"        # Sonos, Inc.
    - "94:9F:3E"        # Sonos, Inc.
    - "B8:E9:37"        # Sonos, Inc.
    - "34:7E:5C"        # Sonos, Inc.
    - "F0:F6:C1"        # Sonos, Inc.
    - "7C:B2:7D"        # Sonos, Inc.
  mdns_services:
    - "_sonos._tcp"
  mdns_txt_keys:
    - "info"
    - "vers"
    - "protovers"
    - "hhsn"
  default_ports: [1400, 1443, 3400, 3401, 3500]
  signature_ports: [1400]
  ssdp_search_target: "urn:schemas-upnp-org:device:ZonePlayer:1"
  ssdp_server_string: "Linux UPnP/1.0 Sonos/70.x"
  hostname_patterns:
    - "^Sonos-"
    - "^sonos[0-9a-f]{12}$"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 1400
    path: "/xml/device_description.xml"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "Linux UPnP/1.0 Sonos"
    body_contains: "Sonos"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "sonos"
  polling_interval_sec: 10
  websocket_event: "sonos:state"
  setup_type: "none"
  ai_chattable: false
  haus_milestone: "M6"

# --- CAPABILITIES ---
capabilities:
  - "media_playback"
  - "volume"
  - "input_select"

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 1400
  transport: "HTTP"
  encoding: "XML"
  auth_method: "none"
  auth_detail: "UPnP SOAP API requires no authentication on the local network. Same protocol as all Sonos speakers."
  base_url_template: "http://{ip}:1400"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "speaker"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.sonos.com/en-us/shop/arc"
  api_docs: "https://developer.sonos.com/reference/"
  developer_portal: "https://developer.sonos.com/"
  support: "https://support.sonos.com/"
  community_forum: "https://en.community.sonos.com/"
  image_url: ""
  fcc_id: "E8THS19"

# --- TAGS ---
tags: ["soundbar", "wifi", "upnp", "sonos", "airplay2", "hdmi-earc", "dolby-atmos", "multi-room", "local-api", "soap"]
---

# Sonos Beam / Arc Soundbar

## What It Is

The Sonos Beam and Sonos Arc are premium smart soundbars manufactured by Sonos, Inc. The Beam (model S14, Gen 2 released 2021) is a compact soundbar with HDMI eARC, Dolby Atmos decoding via virtualization, and five Class-D amplifiers. The Arc (model S19, released 2020) is a full-size premium soundbar with 11 high-excursion drivers, true Dolby Atmos with upward-firing speakers, and HDMI eARC. Both connect to the TV via HDMI eARC (or optical with adapter) and join the home network via Wi-Fi or Ethernet. They expose the same UPnP/SOAP local API as all Sonos speakers on port 1400, with the addition of HDMI input source selection. Both support AirPlay 2, can be paired with Sonos surround speakers (Era 100/300) and a Sub, and participate in Sonos multi-room audio groups.

## How Haus Discovers It

Discovery is identical to the Sonos Era 100. See **[sonos-era-100](sonos-era-100.md)** for the complete discovery flow.

1. **OUI Match** -- MAC prefix from Sonos OUI pool (`00:0E:58`, `5C:AA:FD`, etc.)
2. **mDNS Discovery** -- `_sonos._tcp.local.`
3. **SSDP Discovery** -- `urn:schemas-upnp-org:device:ZonePlayer:1`
4. **HTTP Fingerprint** -- `GET http://{ip}:1400/xml/device_description.xml` returns Sonos device description.
5. **Model Identification** -- `<modelNumber>` field: `S14` = Beam Gen 1, `S36-BEAM` or similar = Beam Gen 2, `S19` = Arc, `S23` = Arc Ultra.

## Pairing / Authentication

No authentication required for the local UPnP/SOAP API. See **[sonos-era-100](sonos-era-100.md)** for full pairing details.

## API Reference

The Sonos Beam and Arc use the exact same UPnP/SOAP protocol on port 1400 as all Sonos speakers. See **[sonos-era-100](sonos-era-100.md)** for the complete API reference covering:

- Transport control (play, pause, stop, next, previous)
- Volume control (get/set volume, mute/unmute)
- Track information and position
- Zone group topology
- UPnP event subscriptions

### Additional Soundbar Capabilities

#### Audio Input Selection

The soundbar adds HDMI/TV audio input as a source. When the soundbar is playing TV audio (via HDMI eARC), the `TrackURI` in `GetPositionInfo` will reflect the TV input source rather than a streaming track:

```
x-sonos-htastream:{RINCON_xxxxxxxx}:spdif
```

This URI indicates the soundbar is receiving audio from its HDMI/SPDIF input. The transport state will be "PLAYING" when TV audio is active.

#### Night Mode and Speech Enhancement

```
POST /MediaRenderer/RenderingControl/Control
SOAPAction: "urn:schemas-upnp-org:service:RenderingControl:1#SetEQ"
```

```xml
<u:SetEQ xmlns:u="urn:schemas-upnp-org:service:RenderingControl:1">
  <InstanceID>0</InstanceID>
  <EQType>NightMode</EQType>
  <DesiredValue>1</DesiredValue>
</u:SetEQ>
```

- `NightMode` -- 0 (off) or 1 (on). Reduces bass and dynamic range for late-night viewing.
- `DialogLevel` -- 0 (off) or 1 (on). Enhances speech clarity.

#### Surround and Sub Configuration

The soundbar acts as the group coordinator for home theater setups. Surround speakers and sub are bonded to the soundbar and appear as members in the `/status/topology` response. The soundbar manages surround volume offsets and sub crossover settings.

## AI Capabilities

AI integration is planned for a future milestone. When implemented, the AI concierge will be able to:
- Report current playback state and whether the soundbar is playing TV audio or streaming content
- Control playback and volume
- Toggle night mode and speech enhancement
- Report home theater group configuration (surrounds, sub)

## Quirks & Notes

- **HDMI eARC** -- The soundbar connects to the TV's eARC-capable HDMI port. CEC (Consumer Electronics Control) allows the TV remote to control soundbar volume. Haus sees this as a network device, not an HDMI device.
- **TV input detection** -- When the soundbar is receiving TV audio, the `TrackURI` contains `x-sonos-htastream:`. This is how Haus distinguishes between TV audio and Sonos streaming playback.
- **Dolby Atmos** -- Arc has true Atmos with upward-firing speakers. Beam Gen 2 virtualizes Atmos. This distinction is not visible via the API but affects the user experience.
- **Home theater bonding** -- Surround speakers and sub are "bonded" to the soundbar and cannot be independently controlled for volume in home theater mode. They appear in the topology as group members.
- **Same SOAP API** -- Despite being soundbars with HDMI, the network API is identical to any Sonos speaker. The HDMI/eARC connection is handled entirely in hardware/firmware.
- **All Sonos quirks apply** -- See [sonos-era-100](sonos-era-100.md) for general Sonos protocol notes (SonosNet mesh, firmware updates, SOAP verbosity, etc.).

## Similar Devices

- **[sonos-era-100](sonos-era-100.md)** -- Compact speaker with the same API, can serve as surround channels for this soundbar
- **[apple-tv-4k](apple-tv-4k.md)** -- Often connected to the same TV, controls audio output to this soundbar via AirPlay or eARC
- **[chromecast-google-tv](chromecast-google-tv.md)** -- Streaming device that outputs audio to this soundbar via HDMI/eARC
