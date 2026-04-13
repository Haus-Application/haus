---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "sonos-era-100"
name: "Sonos Era 100"
manufacturer: "Sonos, Inc."
brand: "Sonos"
model: "S36"
model_aliases: ["Sonos One", "Sonos One SL", "S13", "S18", "S38"]
device_type: "sonos_speaker"
category: "media"
product_line: "Sonos"
release_year: 2023
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "bluetooth", "ethernet"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "00:0E:58"        # Sonos, Inc. (primary OUI)
    - "5C:AA:FD"        # Sonos, Inc. (newer devices)
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
    - "info"            # device info string
    - "vers"            # protocol version
    - "protovers"       # Sonos protocol version
    - "hhsn"            # household serial number
  default_ports: [1400, 1443, 3400, 3401, 3500]
  signature_ports: [1400]
  ssdp_search_target: "urn:schemas-upnp-org:device:ZonePlayer:1"
  ssdp_server_string: "Linux UPnP/1.0 Sonos/70.x (ZPS36)"
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
  - port: 1400
    path: "/status"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: ""
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

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 1400
  transport: "HTTP"
  encoding: "XML"
  auth_method: "none"
  auth_detail: "UPnP SOAP API requires no authentication on the local network. Newer Sonos REST API (port 1443) requires OAuth2 via Sonos Developer Portal."
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
  product_page: "https://www.sonos.com/en-us/shop/era-100"
  api_docs: "https://developer.sonos.com/reference/"
  developer_portal: "https://developer.sonos.com/"
  support: "https://support.sonos.com/"
  community_forum: "https://en.community.sonos.com/"
  image_url: ""
  fcc_id: "E8THS36"

# --- TAGS ---
tags: ["speaker", "wifi", "upnp", "sonos", "airplay2", "multi-room", "local-api", "soap"]
---

# Sonos Era 100

## What It Is

The Sonos Era 100 is a compact smart speaker from Sonos, Inc. that replaced the Sonos One in the product lineup. It supports Wi-Fi 6 and Bluetooth 5.0, AirPlay 2, and Sonos's proprietary multi-room audio system. The Era 100 features a single tweeter and mid-woofer in a sealed acoustic architecture. It connects to the home network via Wi-Fi (or optional Ethernet adapter) and exposes a local UPnP/SOAP API on port 1400 that allows media control, volume adjustment, and group management without cloud dependency. The speaker can be stereo-paired with a second Era 100 or used as surround channels with a Sonos soundbar. Earlier models in this lineage include the Sonos One (S13/S18) and Sonos One SL (S38), which share the same local API and protocol.

## How Haus Discovers It

1. **OUI Match** -- During network scan, any device with MAC prefix `00:0E:58`, `5C:AA:FD`, `54:2A:1B`, `78:28:CA`, `48:A6:B8`, `94:9F:3E`, `B8:E9:37`, `34:7E:5C`, `F0:F6:C1`, or `7C:B2:7D` is flagged as a Sonos device.
2. **mDNS Discovery** -- Browse for `_sonos._tcp.local.` services. Each speaker advertises its household ID and protocol version in TXT records.
3. **SSDP Discovery** -- The speaker responds to UPnP M-SEARCH with search target `urn:schemas-upnp-org:device:ZonePlayer:1`. The response includes the device description XML URL at `/xml/device_description.xml`.
4. **HTTP Fingerprint** -- `GET http://{ip}:1400/xml/device_description.xml` returns UPnP device description XML containing model name, model number, serial number, software version, and room name. The `Server` header contains `Linux UPnP/1.0 Sonos/`.
5. **Model Identification** -- Parse `<modelNumber>` from the device description XML. `S36` maps to Era 100, `S13` to Sonos One, `S18` to Sonos One (Gen 2), `S38` to Sonos One SL.

## Pairing / Authentication

No authentication is required for the local UPnP/SOAP API on port 1400. Any device on the local network can discover and control Sonos speakers.

Initial setup requires the Sonos app (iOS/Android) to provision the speaker onto the Wi-Fi network and register it with a Sonos household. After setup, all local control works without cloud connectivity.

The newer Sonos Developer REST API (cloud-based, available on port 1443 locally or via `api.ws.sonos.com`) requires OAuth2 authentication through the Sonos Developer Portal. This API provides additional capabilities but is not required for basic playback and volume control.

## API Reference

### UPnP/SOAP API (Port 1400)

All Sonos speakers expose a UPnP control API on port 1400 over HTTP. Commands are sent as SOAP requests with appropriate `SOAPAction` headers.

**Base URL:** `http://{ip}:1400`

### Device Description

```
GET /xml/device_description.xml
```

Returns UPnP device description including model, serial number, software version, room name, and available services.

### Get Current Transport State

```
POST /MediaRenderer/AVTransport/Control
SOAPAction: "urn:schemas-upnp-org:service:AVTransport:1#GetTransportInfo"
Content-Type: text/xml; charset="utf-8"

<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
  s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:GetTransportInfo xmlns:u="urn:schemas-upnp-org:service:AVTransport:1">
      <InstanceID>0</InstanceID>
    </u:GetTransportInfo>
  </s:Body>
</s:Envelope>
```

**Response includes:**
- `CurrentTransportState` -- "PLAYING", "PAUSED_PLAYBACK", "STOPPED", "TRANSITIONING"
- `CurrentTransportStatus` -- "OK" or error status
- `CurrentSpeed` -- "1"

### Play / Pause / Stop

```
POST /MediaRenderer/AVTransport/Control
SOAPAction: "urn:schemas-upnp-org:service:AVTransport:1#Play"
```

```xml
<u:Play xmlns:u="urn:schemas-upnp-org:service:AVTransport:1">
  <InstanceID>0</InstanceID>
  <Speed>1</Speed>
</u:Play>
```

Replace `Play` with `Pause` or `Stop` for other transport controls. `Next` and `Previous` advance tracks.

### Get Volume

```
POST /MediaRenderer/RenderingControl/Control
SOAPAction: "urn:schemas-upnp-org:service:RenderingControl:1#GetVolume"
```

```xml
<u:GetVolume xmlns:u="urn:schemas-upnp-org:service:RenderingControl:1">
  <InstanceID>0</InstanceID>
  <Channel>Master</Channel>
</u:GetVolume>
```

**Response:** `<CurrentVolume>` -- integer 0-100.

### Set Volume

```
POST /MediaRenderer/RenderingControl/Control
SOAPAction: "urn:schemas-upnp-org:service:RenderingControl:1#SetVolume"
```

```xml
<u:SetVolume xmlns:u="urn:schemas-upnp-org:service:RenderingControl:1">
  <InstanceID>0</InstanceID>
  <Channel>Master</Channel>
  <DesiredVolume>35</DesiredVolume>
</u:SetVolume>
```

### Get/Set Mute

```
SOAPAction: "urn:schemas-upnp-org:service:RenderingControl:1#GetMute"
SOAPAction: "urn:schemas-upnp-org:service:RenderingControl:1#SetMute"
```

`<DesiredMute>1</DesiredMute>` to mute, `0` to unmute.

### Get Current Track Info

```
SOAPAction: "urn:schemas-upnp-org:service:AVTransport:1#GetPositionInfo"
```

**Response includes:**
- `Track` -- track number in queue
- `TrackDuration` -- "H:MM:SS"
- `TrackMetaData` -- DIDL-Lite XML containing title, artist, album, album art URI
- `TrackURI` -- source URI of current track
- `RelTime` -- current position "H:MM:SS"

### Zone Group Topology

```
GET /status/topology
```

Returns XML describing all speakers in the household, their group assignments, and coordinator roles. This is the key endpoint for understanding multi-room state.

### Subscribe to Events (UPnP Eventing)

```
SUBSCRIBE /MediaRenderer/AVTransport/Event HTTP/1.1
CALLBACK: <http://{haus_ip}:{port}/notify>
NT: upnp:event
TIMEOUT: Second-3600
```

Sonos sends HTTP NOTIFY callbacks when transport state changes. Haus would need to run a local HTTP server to receive these callbacks.

## AI Capabilities

AI integration is planned for a future milestone. When implemented, the AI concierge will be able to:
- Report current playback state (playing, paused, stopped) and track information
- Control playback (play, pause, skip, previous)
- Adjust volume (0-100) and mute/unmute
- Report group/zone membership

## Quirks & Notes

- **Port 1400 is the key** -- This is the signature port for all Sonos speakers. If port 1400 is open and responds to a GET on `/xml/device_description.xml` with Sonos-specific XML, it is definitively a Sonos device.
- **SOAP is verbose** -- The UPnP SOAP API works well but requires careful XML construction. The DIDL-Lite metadata format for track info is particularly complex to parse.
- **Multi-room grouping** -- Sonos speakers form "zones" with one coordinator per group. Commands sent to any group member are forwarded to the coordinator. The `/status/topology` endpoint reveals the full group structure.
- **AirPlay 2** -- Era 100 supports AirPlay 2, which means it also advertises `_airplay._tcp` via mDNS. This is a separate protocol from the Sonos SOAP API.
- **Sonos S2 app required** -- Era 100 requires the S2 version of the Sonos app. Older S1-only devices (Play:5 Gen 1, Connect, etc.) use a slightly different API version but the same SOAP protocol.
- **SonosNet mesh** -- When at least one Sonos device is wired via Ethernet, the speakers form a proprietary 5 GHz mesh network called SonosNet. This reduces Wi-Fi congestion.
- **New REST API** -- Sonos introduced a cloud-based REST API at `api.ws.sonos.com` and a local WebSocket API on port 1443. This is cleaner than SOAP but requires OAuth2 and developer registration. For Haus, the unauthenticated SOAP API on port 1400 is preferred.
- **Firmware updates** -- Delivered via Sonos cloud. Major updates have occasionally broken third-party integrations, particularly the controversial 2024 app redesign.

## Similar Devices

- **[sonos-beam-arc](sonos-beam-arc.md)** -- Sonos soundbar with the same SOAP API plus HDMI input
- **[apple-homepod-mini](apple-homepod-mini.md)** -- Competing smart speaker with AirPlay but no local API
- **[google-nest-mini](google-nest-mini.md)** -- Competing smart speaker with Cast protocol
