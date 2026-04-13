---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "lifx-a19-color"
name: "LIFX A19 Color"
manufacturer: "Feit Electric (LIFX)"
brand: "LIFX"
model: "L3A19LC08"
model_aliases: ["A19", "LIFX Color", "LIFX A19 E26", "L3A19LC09", "LHA19E26UC10P"]
device_type: "lifx_bulb"
category: "lighting"
product_line: "LIFX"
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
  local_only_capable: true
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["D0:73:D5", "D0:73:D5:xx", "50:76:AF"]
  mdns_services: ["_lifx._tcp"]
  mdns_txt_keys: ["id", "md", "pv"]
  default_ports: [56700]
  signature_ports: [56700]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["LIFX.*", "lifx[a-f0-9]+"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "lifx"
  polling_interval_sec: 5
  websocket_event: "lifx:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 56700
  transport: "UDP"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "LAN protocol requires no authentication; broadcast discovery on UDP 56700"
  base_url_template: "udp://{ip}:56700"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "bulb"
  power_source: "mains"
  mounting: "ceiling"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.lifx.com/products/color"
  api_docs: "https://lan.developer.lifx.com/"
  developer_portal: "https://api.developer.lifx.com/"
  support: "https://support.lifx.com/"
  community_forum: "https://community.lifx.com/"
  image_url: ""
  fcc_id: "2AEMI-LHA19E26UC10P"

# --- TAGS ---
tags: ["wifi", "no_hub", "color_bulb", "local_lan_protocol", "udp", "rgbw"]
---

# LIFX A19 Color

## What It Is

The LIFX A19 Color is a WiFi-connected smart bulb that requires no hub — it connects directly to your home WiFi network and is controllable via both a cloud HTTP API and a local LAN protocol. Acquired by Feit Electric in 2022, LIFX has been a pioneer in hubless smart lighting. The A19 form factor fits standard E26 sockets and supports full RGBW color (1600 万 colors), tunable white from 1500K to 9000K, and up to 1100 lumens of brightness. It is one of the best locally-controllable WiFi bulbs on the market, which makes it a strong candidate for Haus integration.

## How Haus Discovers It

1. **mDNS Discovery**: The bulb advertises `_lifx._tcp.local.` via multicast DNS. Each bulb publishes a service instance with TXT records including `id` (MAC-based unique ID), `md` (model description, e.g. "LIFX A19"), and `pv` (protocol version, typically "2").
2. **OUI Match**: MAC addresses beginning with `D0:73:D5` or `50:76:AF` are associated with LIFX devices.
3. **UDP Port Probe**: Port 56700/UDP is the LIFX LAN protocol port. Sending a GetService (type 2) message and receiving a StateService (type 3) response confirms the device.
4. **LAN Protocol Identification**: The GetVersion (type 32) message returns the product ID and vendor ID, allowing Haus to identify the exact LIFX product. LIFX vendor ID is 1; the A19 Color has product IDs including 22, 27, 29, 43, 49, and others depending on hardware revision.

## Pairing / Authentication

No pairing or authentication is required for the LIFX LAN protocol. Any device on the same local network can send UDP messages to port 56700 and control the bulb. Initial WiFi setup is done via the LIFX mobile app, which provisions the bulb's WiFi credentials over a temporary SoftAP or BLE connection. Once connected to the home network, the bulb is immediately controllable locally.

For the cloud HTTP API (used as fallback or for remote access), an OAuth2 bearer token is obtained from `https://cloud.lifx.com/`.

## API Reference

### LIFX LAN Protocol (Primary — Local)

The LIFX LAN protocol operates over UDP on port 56700. Messages use a binary framing format with a 36-byte header followed by type-specific payloads.

#### Message Header Structure (36 bytes)

| Offset | Length | Field | Description |
|--------|--------|-------|-------------|
| 0 | 2 | size | Total message size in bytes (little-endian) |
| 2 | 2 | protocol + flags | Bits 0-11: protocol (1024), Bits 12-15: flags |
| 4 | 4 | source | Unique client identifier |
| 8 | 8 | target | Target device MAC (6 bytes + 2 padding), or all-zeros for broadcast |
| 16 | 6 | reserved | Reserved bytes |
| 22 | 1 | res_required | Response required flag |
| 23 | 1 | ack_required | Acknowledgement required flag |
| 24 | 8 | reserved | Reserved |
| 32 | 2 | type | Message type number |
| 34 | 2 | reserved | Reserved |

#### Key Message Types

| Type | Name | Direction | Description |
|------|------|-----------|-------------|
| 2 | GetService | Client -> Device | Discover devices on network (broadcast) |
| 3 | StateService | Device -> Client | Response with service port |
| 14 | GetHostFirmware | Client -> Device | Get firmware version |
| 15 | StateHostFirmware | Device -> Client | Returns firmware build/version |
| 20 | GetPower | Client -> Device | Get device power level |
| 21 | StatePower | Device -> Client | Returns power level (0=off, 65535=on) |
| 22 | SetPower | Client -> Device | Set power (0 or 65535) |
| 32 | GetVersion | Client -> Device | Get hardware version info |
| 33 | StateVersion | Device -> Client | Returns vendor, product, hw_version |
| 101 | Get (Light) | Client -> Device | Get light state (HSBK + power + label) |
| 102 | SetColor | Client -> Device | Set HSBK color with duration |
| 107 | State (Light) | Device -> Client | Returns full light state |
| 117 | SetLightPower | Client -> Device | Set light power with duration |

#### HSBK Color Model

LIFX uses HSBK (Hue, Saturation, Brightness, Kelvin) with 16-bit unsigned values:

- **Hue**: 0-65535 maps to 0-360 degrees
- **Saturation**: 0-65535 (0 = white, 65535 = fully saturated)
- **Brightness**: 0-65535
- **Kelvin**: 1500-9000 (only meaningful when saturation is 0 for white light)

#### SetColor Payload (13 bytes)

| Offset | Length | Field | Description |
|--------|--------|-------|-------------|
| 0 | 1 | reserved | Reserved |
| 1 | 2 | hue | Hue (0-65535) |
| 3 | 2 | saturation | Saturation (0-65535) |
| 5 | 2 | brightness | Brightness (0-65535) |
| 7 | 2 | kelvin | Color temperature (1500-9000) |
| 9 | 4 | duration | Transition time in milliseconds |

### LIFX Cloud HTTP API (Fallback)

Base URL: `https://api.lifx.com/v1/`

Authentication: `Authorization: Bearer {token}`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/lights/all` | List all lights |
| PUT | `/lights/{selector}/state` | Set light state |
| POST | `/lights/{selector}/toggle` | Toggle power |
| POST | `/lights/{selector}/effects/breathe` | Breathe effect |
| POST | `/lights/{selector}/effects/pulse` | Pulse effect |

## AI Capabilities

When the AI concierge is chatting with a LIFX bulb, it can:
- Turn the light on/off with optional transition duration
- Set brightness as a percentage
- Set any RGB color or named color (converted to HSBK)
- Set color temperature in Kelvin (1500K-9000K)
- Report current state (color, brightness, power, connectivity)
- Apply smooth transitions between colors

## Quirks & Notes

- **No Authentication on LAN**: Any device on the network can control any LIFX bulb. This is a feature for Haus (no pairing needed) but a consideration for network security.
- **UDP Reliability**: The LAN protocol uses UDP, so messages can be lost. Haus should implement retry logic with acknowledgement requests for critical commands.
- **Rate Limiting**: LIFX recommends no more than 20 messages per second per bulb on the LAN protocol. The cloud API is rate-limited to approximately 120 requests per minute.
- **WiFi Congestion**: Each LIFX bulb is a full WiFi client. Networks with many LIFX bulbs (10+) may experience WiFi congestion. This is a common complaint versus hub-based systems.
- **Firmware Variations**: Older firmware versions may not support all message types. Always check firmware version via GetHostFirmware before using newer features.
- **Initial Setup Requires App**: The first-time WiFi provisioning must be done through the LIFX mobile app. Haus cannot provision a brand-new bulb, only discover and control already-configured ones.
- **Multizone Products**: Some LIFX products (Beam, Strip) use extended multizone messages. The A19 is single-zone only.

## Similar Devices

- **lifx-mini-color** — Smaller form factor LIFX bulb, same LAN protocol
- **lifx-clean** — LIFX bulb with HEV (germicidal) light, same protocol with additional HEV messages
- **nanoleaf-essentials-a19** — Thread/Matter bulb, different protocol entirely
- **wyze-bulb-color** — WiFi bulb but cloud-only, no local control
