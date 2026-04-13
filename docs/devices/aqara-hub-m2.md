---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "aqara-hub-m2"
name: "Aqara Hub M2"
manufacturer: "Lumi United Technology Co., Ltd."
brand: "Aqara"
model: "HM2-G01"
model_aliases: ["ZHWG12LM", "Aqara Hub M2", "Aqara Smart Hub M2"]
device_type: "aqara_hub"
category: "smart_home"
product_line: "Aqara"
release_year: 2021
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
  protocols_spoken: ["wifi", "zigbee", "ethernet"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "54:EF:44"              # Lumi United Technology (Aqara)
    - "04:CF:8C"              # Xiaomi / Lumi
    - "7C:49:EB"              # Xiaomi Communications
    - "28:6C:07"              # Xiaomi Communications
    - "78:11:DC"              # Xiaomi Communications
  mdns_services:
    - "_hap._tcp"             # HomeKit Accessory Protocol
    - "_aqara._tcp"           # Aqara-specific service (reported by some firmware versions)
  mdns_txt_keys: ["md", "pv", "id", "c#", "s#", "ff", "ci", "sf", "sh"]  # Standard HomeKit TXT keys
  default_ports: [80, 443, 4443, 9898]
  signature_ports: [4443, 9898]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^lumi-gateway-.*", "^Aqara-Hub-.*", "^AqaraHub-.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 4443
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "aqara"
  polling_interval_sec: 10
  websocket_event: "aqara:state"
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp", "scenes", "groups", "motion", "temperature", "humidity"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 4443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "The Aqara Hub M2 supports HomeKit (HAP over IP) as its primary local API. HAP uses mDNS discovery and an 8-digit setup code for initial pairing, then establishes encrypted sessions via SRP-6a and ChaCha20-Poly1305. The Aqara Developer API (cloud) uses OAuth2. The legacy Xiaomi/Lumi UDP protocol on port 9898 is available on older firmware but deprecated."
  base_url_template: ""
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "usb"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "zigbee", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.aqara.com/us/hub-m2.html"
  api_docs: "https://developer.aqara.com/cloud/api-introduction"
  developer_portal: "https://developer.aqara.com/"
  support: "https://www.aqara.com/us/support.html"
  community_forum: "https://community.aqara.com/"
  image_url: ""
  fcc_id: "2AKIT-HM2G01"

# --- TAGS ---
tags: ["zigbee_hub", "homekit", "matter_ready", "ir_blaster", "wifi", "ethernet", "aqara", "xiaomi", "gateway"]
---

# Aqara Hub M2

## What It Is

The Aqara Hub M2 is a smart home gateway that bridges Aqara's extensive line of Zigbee 3.0 sensors and devices to your WiFi/Ethernet network. It is the central coordinator for Aqara's ecosystem -- all those door sensors, motion sensors, temperature sensors, vibration sensors, and smart plugs communicate through this hub. The M2 connects via WiFi (2.4 GHz) or Ethernet and supports Apple HomeKit natively, with Matter support added via firmware update. It also has a built-in infrared (IR) blaster for controlling legacy IR devices like TVs, air conditioners, and fans. The hub features a small built-in speaker for alarm sounds and a night light ring. Setup requires the Aqara Home app initially, but once configured, the hub can operate locally via HomeKit without cloud dependency.

## How Haus Discovers It

1. **OUI match** -- The hub's WiFi or Ethernet MAC address matches Lumi/Xiaomi OUI prefixes: `54:EF:44`, `04:CF:8C`, `7C:49:EB`, `28:6C:07`, or `78:11:DC`.
2. **mDNS** -- The hub advertises `_hap._tcp` (HomeKit Accessory Protocol) via mDNS. The TXT record contains standard HomeKit fields including `md` (model name, e.g., "Aqara Hub M2"), `ci` (category identifier, 2 = bridge), and `sf` (status flags -- `sf=1` means unpaired, `sf=0` means paired).
3. **Hostname pattern** -- DHCP hostname typically matches `lumi-gateway-*` or `Aqara-Hub-*`.
4. **Port probe** -- Port 4443 (HTTPS) and legacy port 9898 (UDP) are signature ports. Port 4443 serves the hub's local HTTPS interface.
5. **HomeKit category** -- The mDNS TXT record `ci=2` identifies this as a HomeKit bridge, narrowing device type.

## Pairing / Authentication

### HomeKit (Primary Local Path)

1. **mDNS discovery** -- Find the hub via `_hap._tcp` mDNS service.
2. **HomeKit setup code** -- The 8-digit setup code is printed on the bottom of the hub and on the included card (format: XXX-XX-XXX).
3. **SRP-6a pairing** -- HomeKit uses Secure Remote Password protocol (SRP-6a) for initial pairing, establishing a shared secret.
4. **Session encryption** -- All subsequent communication uses ChaCha20-Poly1305 AEAD encryption over TCP.
5. **Accessory database** -- Once paired, the hub exposes all its child Zigbee devices as HomeKit accessories in a single HAP accessory database.

### Aqara Developer API (Cloud Path)

1. **Register as developer** -- Create an account at developer.aqara.com.
2. **OAuth2 authorization** -- Redirect user to Aqara OAuth endpoint, receive authorization code.
3. **Token exchange** -- Exchange code for access token and refresh token.
4. **API calls** -- Use bearer token for cloud API calls to query/control devices.

### Legacy Xiaomi UDP Protocol (Port 9898, Deprecated)

Older firmware versions supported a local UDP multicast protocol on port 9898 using the Xiaomi Mi Home Gateway protocol. This used a pre-shared key (obtained from Mi Home app developer options) and AES-128-CBC encryption. This protocol is deprecated on newer firmware and may not be available on the M2.

## API Reference

### HomeKit Accessory Protocol (HAP) -- Local

The HomeKit protocol is the primary local integration path. After pairing, the hub exposes a HAP accessory database with the following service types for child devices:

| HAP Service | UUID | Description |
|-------------|------|-------------|
| `ContactSensor` | `0x0080` | Door/window sensor open/close state |
| `MotionSensor` | `0x0085` | PIR motion detection |
| `TemperatureSensor` | `0x008A` | Temperature reading (Celsius) |
| `HumiditySensor` | `0x0082` | Relative humidity percentage |
| `LightSensor` | `0x0084` | Ambient light level (lux) |
| `Lightbulb` | `0x0043` | Hub night light control |
| `SecuritySystem` | `0x007E` | Alarm mode (home/away/night/off) |
| `Switch` | `0x0049` | Smart plug on/off |

HAP characteristics for each service follow Apple's HomeKit Accessory Protocol Specification (non-commercial version available for open-source projects).

### Aqara Developer Cloud API

Base URL: `https://open-{region}.aqara.com/v3.0/open`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/device/query` | POST | List all devices on account |
| `/resource/query` | POST | Query current state of device resources |
| `/resource/write` | POST | Write value to device resource |
| `/scene/query` | POST | List configured scenes |
| `/scene/run` | POST | Execute a scene |
| `/ifttt/query` | POST | List automations |

Regions: `cn` (China), `usa` (US), `ger` (Europe), `kr` (Korea), `ru` (Russia), `sg` (Singapore).

### IR Blaster

The hub's built-in IR blaster can learn and transmit IR codes. Via the Aqara app, you can:

- Select from a database of known IR codes for TVs, ACs, fans
- Learn custom IR codes from existing remotes
- Create virtual buttons that transmit stored IR sequences

IR control is available via the cloud API as virtual device resources but not directly via HomeKit.

## AI Capabilities

AI integration for the Aqara Hub M2 is planned for M11. When implemented, the AI concierge could:

- Report status of all connected Zigbee sensors (doors open/closed, motion detected, temperature/humidity)
- Control the hub's night light (on/off, brightness, color)
- Trigger scenes configured in the Aqara ecosystem
- Set/query the security alarm mode
- Control IR devices through the hub's IR blaster
- Provide environmental summaries ("The bedroom is 72 degrees and 45% humidity, all doors closed")

## Quirks & Notes

- **Matter support via firmware update** -- Aqara added Matter bridge functionality to the M2 via firmware update (version 4.0.4 or later). When Matter is enabled, the hub can expose child devices as Matter accessories to other Matter controllers. This is a key integration path for Haus at M11.
- **HomeKit is the best local API** -- The HAP protocol is well-documented (Apple's specification is public for non-commercial use) and provides real-time event notifications. Libraries like `hc` (Go) implement HAP controller functionality.
- **Dual network** -- Can connect via both WiFi and Ethernet simultaneously. Ethernet is recommended for reliability.
- **IR blaster limitations** -- The IR blaster has a fixed position and limited range/angle. It works best with line-of-sight to the target device within about 8 meters.
- **Zigbee child device limit** -- The M2 supports up to 128 Zigbee child devices. In practice, having more than 40-50 devices can cause performance degradation.
- **2.4 GHz WiFi only** -- Does not support 5 GHz WiFi networks.
- **Legacy UDP protocol may be disabled** -- Newer firmware versions (3.x+) may have the legacy Xiaomi UDP port 9898 protocol disabled by default. Do not rely on it.
- **Cloud required for initial setup** -- The Aqara Home app (which requires cloud login) is needed for initial hub setup, WiFi configuration, and Zigbee device pairing. After setup, local operation via HomeKit works without internet.
- **Built-in speaker** -- Produces alarm sounds, doorbell tones, and can be used as a basic alert system. Volume is modest.

## Similar Devices

- [aqara-door-window-sensor](aqara-door-window-sensor.md) -- Zigbee contact sensor that pairs to this hub
- [aqara-motion-sensor-p1](aqara-motion-sensor-p1.md) -- Zigbee motion sensor that pairs to this hub
- [smartthings-hub-v3](smartthings-hub-v3.md) -- Competing multi-protocol hub with cloud API
- [ikea-dirigera-hub](ikea-dirigera-hub.md) -- Similar Zigbee + Thread + Matter hub with local API
- [philips-hue-bridge](philips-hue-bridge.md) -- Zigbee bridge with mature local REST API
