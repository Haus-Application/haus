---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "eve-motion"
name: "Eve Motion"
manufacturer: "Eve Systems GmbH"
brand: "Eve"
model: "20EBN8101"
model_aliases: ["Eve Motion 2nd Gen", "1EM109901000", "20EBN8101"]
device_type: "eve_motion_sensor"
category: "security"
product_line: "Eve"
release_year: 2022
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: false
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["thread", "matter", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []            # Thread device; MAC is on the Thread mesh, not directly visible on WiFi/Ethernet
  mdns_services:
    - "_matterc._udp"         # Matter commissioning (when uncommissioned)
    - "_matter._tcp"           # Matter operational (when commissioned)
  mdns_txt_keys: ["DN", "VP", "D", "CM", "DT", "PH"]  # Matter mDNS TXT keys
  default_ports: [5540]       # Matter default port
  signature_ports: [5540]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []         # Matter uses UDP/CASE sessions, not HTTP

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "matter"
  polling_interval_sec: 0     # Matter subscriptions provide real-time updates
  websocket_event: "matter:state"
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["motion", "battery_level"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 5540
  transport: "UDP"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Matter commissioning via BLE or Thread. Initial setup uses a setup code (QR code or 11-digit manual pairing code on the device). Commissioning establishes CASE (Certificate Authenticated Session Establishment) sessions using NIST P-256 ECDSA certificates. All subsequent communication is encrypted and authenticated via Matter fabric credentials."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "sensor"
  power_source: "battery"
  mounting: "wall"
  indoor_outdoor: "both"
  wireless_radios: ["thread", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.evehome.com/en/eve-motion"
  api_docs: ""
  developer_portal: "https://csa-iot.org/developer-resource/specifications/"
  support: "https://www.evehome.com/en/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AHPQ20EBN8101"

# --- TAGS ---
tags: ["thread", "matter", "homekit", "no_cloud", "battery", "pir", "motion_sensor", "ipv6", "eve"]
---

# Eve Motion

## What It Is

The Eve Motion is a wireless passive infrared (PIR) motion sensor built on Thread and Matter, designed for both indoor and outdoor use (IPX3 water resistance). It detects movement within a 120-degree field of view at up to 9 meters range and was one of the first consumer sensors to ship with native Thread and Matter support. The device operates on two AA batteries with an estimated battery life of up to 2.5 years, communicates over Thread (802.15.4 mesh networking), and requires absolutely no cloud account or internet connection to function. It was originally HomeKit-only but received a firmware update enabling Matter support. Eve Systems is a pioneer in the "no cloud, no data collection" philosophy -- all data stays local.

## How Haus Discovers It

The Eve Motion communicates over Thread, which means it is reachable via IPv6 on the local network through a Thread Border Router (such as an Apple HomePod Mini, Apple TV 4K, or a dedicated Thread border router):

1. **Thread Border Router** -- Haus must first have connectivity to the Thread mesh network via a Thread Border Router that advertises routes to the Thread network on the LAN.
2. **Matter mDNS** -- When uncommissioned, the Eve Motion advertises `_matterc._udp` (Matter commissionable) via mDNS through the Thread Border Router's DNS-SD proxy. TXT records include `VP` (vendor/product), `D` (discriminator), `CM` (commissioning mode), and `DT` (device type).
3. **Matter device type** -- The Matter device type identifier for occupancy sensors is `0x0107`. The vendor ID for Eve Systems is `0x130A` (4874 decimal).
4. **Once commissioned** -- The device advertises `_matter._tcp` as an operational Matter node and is reachable on Matter port 5540 via its Thread IPv6 address.

## Pairing / Authentication

### Matter Commissioning

1. **Obtain setup code** -- The Eve Motion has a QR code on the device body and a manual 11-digit pairing code (format: XXXX-XXX-XXXX) printed on the packaging and device.
2. **Open commissioning window** -- For a new (factory reset) device, the commissioning window is open by default. For a previously paired device, factory reset by pressing and holding the button for 10 seconds until the LED flashes amber.
3. **BLE or Thread commissioning** -- The Matter controller (Haus) connects via BLE for the initial PASE (Passcode-Authenticated Session Establishment) using the setup code. This establishes trust.
4. **Thread network credentials** -- The controller provisions Thread network credentials to the device so it can join the Thread mesh.
5. **CASE session** -- After joining Thread, the device establishes a CASE (Certificate Authenticated Session Establishment) session with the controller using Matter operational certificates. This creates an encrypted, authenticated channel.
6. **Fabric enrollment** -- The device is enrolled in the Haus Matter fabric and receives an operational node ID.

### Multi-Admin (Multi-Fabric)

Matter supports multi-admin, meaning the Eve Motion can be commissioned to multiple controllers simultaneously (e.g., Haus and Apple Home). Each controller maintains its own fabric and CASE session.

## API Reference

The Eve Motion is controlled via the Matter protocol. It implements the following Matter clusters:

### Matter Clusters

| Cluster | ID | Description |
|---------|-----|-------------|
| Occupancy Sensing | `0x0406` | Motion detection state |
| Power Source | `0x002F` | Battery level and status |
| Descriptor | `0x001D` | Device type and cluster list |
| Basic Information | `0x0028` | Vendor name, product name, serial number, firmware version |
| General Commissioning | `0x0030` | Commissioning state management |
| Network Commissioning | `0x0031` | Thread network credentials |
| Thread Network Diagnostics | `0x0035` | Thread mesh diagnostics |

### Occupancy Sensing Cluster (0x0406)

| Attribute | ID | Type | Description |
|-----------|----|------|-------------|
| `Occupancy` | `0x0000` | bitmap8 | Bit 0: `1` = motion detected, `0` = no motion |
| `OccupancySensorType` | `0x0001` | enum8 | `0` = PIR, `1` = Ultrasonic, `2` = PIRAndUltrasonic, `3` = PhysicalContact |
| `PIROccupiedToUnoccupiedDelay` | `0x0010` | uint16 | Seconds before sensor reports "no motion" after last detection |

### Power Source Cluster (0x002F)

| Attribute | ID | Type | Description |
|-----------|----|------|-------------|
| `Status` | `0x0000` | enum8 | Power source status |
| `BatVoltage` | `0x000B` | uint32 | Battery voltage in mV |
| `BatPercentRemaining` | `0x000C` | uint8 | Battery percentage (0-200, divide by 2) |
| `BatChargeLevel` | `0x000E` | enum8 | `0` = OK, `1` = Warning, `2` = Critical |

### Matter Subscriptions

Matter supports attribute subscriptions for real-time updates:

```
Subscribe to Occupancy attribute:
  Node ID: {operational_node_id}
  Endpoint: 1
  Cluster: 0x0406
  Attribute: 0x0000
  MinInterval: 0 seconds
  MaxInterval: 300 seconds
```

The device sends unsolicited reports on state changes (motion detected / no motion) and periodic reports at the max interval.

## AI Capabilities

When Matter integration is complete (M11), the AI concierge could:

- Report motion status for rooms with Eve Motion sensors
- Provide motion event history and occupancy patterns
- Alert when motion is detected in monitored areas (e.g., while away)
- Report battery levels and warn when replacement is needed
- Combine motion data with other sensors for room occupancy intelligence

## Quirks & Notes

- **Thread mesh networking** -- The Eve Motion is a Thread Sleepy End Device (SED), meaning it sleeps most of the time to conserve battery and wakes briefly to poll its Thread parent for messages. This means commands sent to it may have slight latency (typically under 1 second) until the next poll interval.
- **No cloud, ever** -- Eve Systems has a strong privacy stance. There is no cloud account, no cloud API, no telemetry. Everything is local. This aligns perfectly with Haus's local-first philosophy.
- **IPX3 water resistance** -- Rated for outdoor use. Can handle rain splashing from any direction, but should not be submerged or exposed to pressure washing.
- **Thread Border Router required** -- The sensor cannot communicate directly with WiFi/Ethernet devices. A Thread Border Router must be present on the network to bridge between Thread (802.15.4) and IP (WiFi/Ethernet). Apple HomePod Mini, Apple TV 4K (2021+), and some other devices act as Thread Border Routers.
- **Firmware updates via BLE** -- Eve provides firmware updates through the Eve app (iOS) over Bluetooth Low Energy. The Matter firmware update was a significant one that added Matter support to originally HomeKit-only devices.
- **AA batteries** -- Uses two standard AA batteries, which is more convenient and cost-effective than coin cells. Expected life is approximately 2.5 years.
- **120-degree FOV** -- Narrower than some competing sensors (Aqara P1 is 170 degrees). Placement matters more for coverage.
- **PIR-only** -- No illuminance/light level sensor. If you need light level for conditional automations, you will need a separate sensor or a different product.

## Similar Devices

- [eve-door-window](eve-door-window.md) -- Thread/Matter contact sensor from the same Eve platform
- [aqara-motion-sensor-p1](aqara-motion-sensor-p1.md) -- Zigbee motion sensor with light level sensor (requires hub)
- [aqara-hub-m2](aqara-hub-m2.md) -- Zigbee hub that can bridge Aqara motion sensors
