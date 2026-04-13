---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "nanoleaf-essentials-a19"
name: "Nanoleaf Essentials A19 Bulb"
manufacturer: "Nanoleaf"
brand: "Nanoleaf"
model: "NL45"
model_aliases: ["NL45-0800", "Essentials A19", "Essentials Bulb", "NL67"]
device_type: "matter_bulb"
category: "lighting"
product_line: "Nanoleaf Essentials"
release_year: 2020
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: false
  cloud_api: false
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["thread", "matter", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["00:55:DA", "80:E4:DA"]
  mdns_services: ["_matterc._udp", "_matterd._udp"]
  mdns_txt_keys: ["VP", "D", "CM", "DT", "DN", "RI", "PI", "PH", "SII", "SAI"]
  default_ports: [5540]
  signature_ports: [5540]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "matter"
  polling_interval_sec: 0
  websocket_event: "matter:state"
  setup_type: "app_pairing"
  ai_chattable: true
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp"]

# --- PROTOCOL ---
protocol:
  type: "coap"
  port: 5540
  transport: "UDP"
  encoding: "CBOR"
  auth_method: "none"
  auth_detail: "Matter commissioning via BLE or Thread; CASE/PASE session establishment; no traditional API key"
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "bulb"
  power_source: "mains"
  mounting: "ceiling"
  indoor_outdoor: "indoor"
  wireless_radios: ["thread", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://nanoleaf.me/en-US/products/essentials/bulb/"
  api_docs: ""
  developer_portal: "https://csa-iot.org/developer-resource/specifications-download-request/"
  support: "https://nanoleaf.me/en-US/support/"
  community_forum: "https://forum.nanoleaf.me/"
  image_url: ""
  fcc_id: "2AWLPNL45"

# --- TAGS ---
tags: ["thread", "matter", "bluetooth_le", "no_wifi", "no_hub_required_with_border_router", "color_bulb"]
---

# Nanoleaf Essentials A19 Bulb

## What It Is

The Nanoleaf Essentials A19 is a smart bulb that communicates over Thread mesh networking with Matter protocol support. Unlike WiFi-based smart bulbs, it does not connect to your WiFi network at all — instead it joins a Thread mesh network coordinated by a Thread border router (such as an Apple HomePod Mini, Apple TV 4K, Google Nest Hub 2nd gen, or Samsung SmartThings Station). The bulb was one of the first consumer Thread/Matter lighting products on the market. It supports full RGBW color (16M colors), tunable white (2700K-6500K), and 806 lumens brightness. Commissioning is done via Bluetooth LE.

## How Haus Discovers It

Thread/Matter devices are discovered differently from WiFi devices:

1. **Thread Border Router**: Haus must be acting as, or communicating with, a Thread border router to see Thread devices on the mesh network.
2. **mDNS (Matter)**: Matter-commissioned devices advertise `_matterc._udp` (commissionable) or `_matterd._udp` (commissioned) via mDNS. The TXT records include `VP` (Vendor ID + Product ID), `D` (discriminator), `CM` (commissioning mode), `DT` (device type), and `DN` (device name).
3. **BLE Advertising**: Before commissioning, the bulb advertises via BLE with Matter-specific BLE service UUID `0xFFF6`. The advertisement contains the discriminator and vendor/product info.
4. **Matter Fabric Discovery**: Once commissioned onto a Matter fabric, the device is discoverable via the fabric's operational discovery (mDNS with `_matterd._udp` using fabric-compressed ID).
5. **OUI Match**: While Thread devices do not use WiFi MACs on the network, the BLE MAC may match Nanoleaf OUIs (`00:55:DA`, `80:E4:DA`).

## Pairing / Authentication

Matter commissioning uses a multi-step cryptographic pairing flow:

1. **Setup Code**: The bulb ships with a Matter setup code (numeric, printed on the bulb or packaging) and a QR code. The setup code format is `XXXX-XXX-XXXX` (11-digit numeric with dashes).
2. **PASE (Passcode-Authenticated Session Establishment)**: Haus initiates a PASE session over BLE using the setup code as the shared secret (SPAKE2+ protocol).
3. **Thread Network Credentials**: During commissioning, Haus provides the Thread network credentials (obtained from the Thread border router) to the bulb over the PASE-secured BLE connection.
4. **CASE (Certificate-Authenticated Session Establishment)**: After joining the Thread network, the bulb and Haus establish a CASE session using certificates from the Matter fabric's Certificate Authority.
5. **Fabric Enrollment**: The device is now part of the Haus Matter fabric and can be controlled via Matter clusters over the Thread network.

No traditional API key or password is involved. The Matter fabric's PKI handles all ongoing authentication.

## API Reference

The Nanoleaf Essentials A19 has no HTTP REST API. All control is performed via Matter protocol data model (clusters and attributes) over Thread.

### Matter Clusters Supported

| Cluster | Cluster ID | Description |
|---------|-----------|-------------|
| On/Off | 0x0006 | Power state control |
| Level Control | 0x0008 | Brightness control (0-254) |
| Color Control | 0x0300 | HSV color, color temperature (mireds), XY color |
| Descriptor | 0x001D | Device type and cluster lists |
| Identify | 0x0003 | Identify the device (flashing) |
| Groups | 0x0004 | Group membership |
| Scenes | 0x0005 | Scene recall and storage |

### Key Attributes

| Cluster | Attribute | ID | Type | Description |
|---------|-----------|-----|------|-------------|
| On/Off | OnOff | 0x0000 | bool | Current power state |
| Level Control | CurrentLevel | 0x0000 | uint8 | Brightness 0-254 |
| Color Control | CurrentHue | 0x0000 | uint8 | Hue 0-254 |
| Color Control | CurrentSaturation | 0x0001 | uint8 | Saturation 0-254 |
| Color Control | ColorTemperatureMireds | 0x0007 | uint16 | Color temp in mireds (154-370 for this bulb) |
| Color Control | ColorMode | 0x0008 | enum8 | 0=HS, 1=XY, 2=CT |

### Key Commands

| Cluster | Command | Description |
|---------|---------|-------------|
| On/Off | On (0x01) | Turn on |
| On/Off | Off (0x00) | Turn off |
| On/Off | Toggle (0x02) | Toggle power |
| Level Control | MoveToLevel (0x00) | Set brightness with transition time |
| Color Control | MoveToHue (0x00) | Set hue with transition |
| Color Control | MoveToSaturation (0x03) | Set saturation with transition |
| Color Control | MoveToHueAndSaturation (0x06) | Set both with transition |
| Color Control | MoveToColorTemperature (0x0A) | Set color temp in mireds with transition |

## AI Capabilities

When the AI concierge is chatting with this bulb via Matter, it can:
- Turn the bulb on/off
- Set brightness level
- Set color via hue/saturation or color temperature
- Report current state (power, brightness, color, color mode)
- Identify the bulb (cause it to flash)

## Quirks & Notes

- **No WiFi**: This bulb has no WiFi radio. It communicates exclusively over Thread (802.15.4) mesh networking. Haus must have access to a Thread border router or implement one.
- **Thread Border Router Required**: Without a Thread border router on the network, the bulb cannot be reached from the IP network. Apple HomePod Mini, Apple TV 4K (2nd gen+), Google Nest Hub (2nd gen), and some other devices serve as Thread border routers.
- **BLE Range for Commissioning**: Initial commissioning is done over BLE, which has limited range (typically 5-10 meters). The hub must be within BLE range of the bulb during setup.
- **Firmware Updates via Matter OTA**: The bulb supports Matter OTA update cluster for firmware updates.
- **Color Temperature Range**: 2700K-6500K (approximately 154-370 mireds). Attempting to set values outside this range will be clamped.
- **Original Apple HomeKit Thread**: The first hardware revisions (NL45) shipped with Apple HomeKit over Thread before Matter existed. These were later updated via firmware to support Matter. Some very old firmware versions may only support HomeKit.
- **Multi-Admin**: Matter supports multi-admin — the bulb can be commissioned to multiple fabrics simultaneously (e.g., both Haus and Apple Home).
- **Subscription Model**: Matter supports attribute subscriptions, so Haus can receive real-time state change notifications without polling.

## Similar Devices

- **nanoleaf-shapes** — WiFi panel lights from Nanoleaf, REST API, different protocol
- **lifx-a19-color** — WiFi bulb with local UDP protocol, no Thread/Matter
- **eve-energy** — Thread/Matter smart plug, same Matter commissioning flow
- **nanoleaf-essentials-strip** — Thread/Matter LED strip from same product line
