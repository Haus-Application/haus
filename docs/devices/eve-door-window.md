---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "eve-door-window"
name: "Eve Door & Window"
manufacturer: "Eve Systems GmbH"
brand: "Eve"
model: "20EBN1101"
model_aliases: ["Eve Door & Window 2nd Gen", "1ED109901000", "20EBN1101"]
device_type: "eve_contact_sensor"
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
  mac_prefixes: []            # Thread device; MAC on Thread mesh, not WiFi/Ethernet
  mdns_services:
    - "_matterc._udp"         # Matter commissioning (uncommissioned)
    - "_matter._tcp"           # Matter operational (commissioned)
  mdns_txt_keys: ["DN", "VP", "D", "CM", "DT", "PH"]
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
  polling_interval_sec: 0     # Matter subscriptions for real-time updates
  websocket_event: "matter:state"
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["motion", "battery_level"]   # open/close modeled as motion

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 5540
  transport: "UDP"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Matter commissioning via BLE or Thread. Setup code (QR or manual 11-digit) printed on device. PASE for initial pairing, CASE for operational sessions. All communication encrypted with Matter fabric credentials."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "sensor"
  power_source: "battery"
  mounting: "door"
  indoor_outdoor: "indoor"
  wireless_radios: ["thread", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.evehome.com/en/eve-door-window"
  api_docs: ""
  developer_portal: "https://csa-iot.org/developer-resource/specifications/"
  support: "https://www.evehome.com/en/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AHPQ20EBN1101"

# --- TAGS ---
tags: ["thread", "matter", "homekit", "no_cloud", "battery", "contact_sensor", "door_window", "ipv6", "eve"]
---

# Eve Door & Window

## What It Is

The Eve Door & Window is a Thread/Matter contact sensor that detects whether a door, window, cabinet, or drawer is open or closed. Like the Eve Motion, it is built on Eve Systems' privacy-first, no-cloud platform. The sensor consists of two parts -- a main body with a magnetometer and a separate magnet -- that detect proximity changes when a door or window opens and closes. It communicates over Thread (802.15.4 mesh), supports Matter natively (via firmware update for older units, natively on newer production runs), and requires no cloud account, no internet connection, and no proprietary hub. It runs on a single CR2032 coin cell battery with approximately one year of battery life. The sensor was originally HomeKit-only but gained Matter compatibility through a firmware update.

## How Haus Discovers It

Discovery follows the same Thread/Matter path as other Eve devices:

1. **Thread Border Router** -- A Thread Border Router (Apple HomePod Mini, Apple TV 4K, or dedicated TBR) must be present on the network to bridge between the Thread mesh and the IP network.
2. **Matter mDNS** -- When uncommissioned, the sensor advertises `_matterc._udp` via the Thread Border Router's DNS-SD proxy. TXT records include vendor/product ID (`VP`), discriminator (`D`), and device type (`DT`).
3. **Matter device type** -- The Matter device type for contact sensors is `0x0015` (Contact Sensor). Eve's vendor ID is `0x130A` (4874).
4. **Commissioned state** -- Once commissioned, the device appears as an operational Matter node advertising `_matter._tcp` and is reachable on port 5540 via its Thread IPv6 address.

## Pairing / Authentication

### Matter Commissioning

1. **Setup code** -- QR code on the device body and manual 11-digit pairing code on packaging.
2. **Factory reset** -- Press and hold the button on the device for 10 seconds until the LED flashes amber to open a new commissioning window.
3. **BLE commissioning** -- The Matter controller connects via BLE for PASE (Passcode-Authenticated Session Establishment) using the setup code.
4. **Thread provisioning** -- The controller sends Thread network credentials so the device joins the Thread mesh.
5. **CASE session** -- Encrypted CASE session established using Matter operational certificates.
6. **Multi-admin** -- Supports simultaneous commissioning to multiple Matter fabrics (e.g., Haus and Apple Home).

## API Reference

The Eve Door & Window implements the following Matter clusters:

### Matter Clusters

| Cluster | ID | Description |
|---------|-----|-------------|
| Boolean State | `0x0045` | Open/closed state (primary contact detection) |
| Power Source | `0x002F` | Battery level and status |
| Descriptor | `0x001D` | Device type and cluster list |
| Basic Information | `0x0028` | Vendor name, product name, serial, firmware |
| General Commissioning | `0x0030` | Commissioning state |
| Network Commissioning | `0x0031` | Thread network credentials |
| Thread Network Diagnostics | `0x0035` | Thread mesh diagnostics |

### Boolean State Cluster (0x0045)

| Attribute | ID | Type | Description |
|-----------|----|------|-------------|
| `StateValue` | `0x0000` | boolean | `true` = closed (contact), `false` = open (no contact) |

Note: The Boolean State cluster reports the physical state of the magnetic reed switch. `true` means the magnet is in proximity (door closed), `false` means the magnet is away (door open). Some implementations invert this, so verify during integration testing.

The cluster also generates a `StateChange` event when the state transitions, which is delivered via Matter subscriptions.

### Power Source Cluster (0x002F)

| Attribute | ID | Type | Description |
|-----------|----|------|-------------|
| `Status` | `0x0000` | enum8 | Power source status |
| `BatVoltage` | `0x000B` | uint32 | Battery voltage in mV |
| `BatPercentRemaining` | `0x000C` | uint8 | Battery percentage (0-200, divide by 2) |
| `BatChargeLevel` | `0x000E` | enum8 | `0` = OK, `1` = Warning, `2` = Critical |
| `BatReplacementDescription` | `0x0013` | string | "CR2032" |

### Matter Subscriptions

```
Subscribe to StateValue:
  Node ID: {operational_node_id}
  Endpoint: 1
  Cluster: 0x0045
  Attribute: 0x0000
  MinInterval: 0 seconds
  MaxInterval: 3600 seconds
```

The device sends immediate unsolicited reports on state changes (open/close events) and periodic keep-alive reports at the max interval.

## AI Capabilities

When Matter integration is complete (M11), the AI concierge could:

- Report whether specific doors and windows are currently open or closed
- Track open/close history with timestamps
- Alert when a door or window has been left open for a configurable duration
- Provide security summaries ("All doors and windows are closed" or "The back door has been open for 15 minutes")
- Combine with motion sensor data for occupancy and security intelligence
- Report battery levels and warn when CR2032 replacement is needed

## Quirks & Notes

- **Thread Sleepy End Device** -- Like the Eve Motion, this is a Thread SED. It sleeps most of the time and wakes briefly to poll its Thread parent. State change reports are sent immediately upon detection (interrupt-driven), but responses to read requests may have slight latency.
- **No cloud, period** -- Zero cloud dependency. No Eve account required. No telemetry. All data stays on the local Thread network. Eve is one of the very few IoT manufacturers that truly commits to this.
- **CR2032 battery** -- Single coin cell. Battery life is approximately 1 year, which is shorter than the Eve Motion's AA battery life. The trade-off is a much smaller form factor.
- **Magnetic sensitivity** -- The sensor uses a magnetometer rather than a simple reed switch. This provides more reliable detection but can be affected by strong magnetic fields from nearby appliances or speakers.
- **Thread Border Router required** -- Cannot function without a Thread Border Router on the network. If no TBR is present, the sensor has no way to communicate.
- **Firmware update via BLE** -- Updates delivered through the Eve app (iOS) over Bluetooth. The Matter enablement update was a significant change for originally HomeKit-only units.
- **Indoor only** -- Unlike the Eve Motion, the Door & Window sensor is not water-resistant. Indoor use only.
- **Matter Boolean State vs. Occupancy** -- The Eve Door & Window uses the Boolean State cluster (0x0045), not the Occupancy Sensing cluster (0x0406) used by the Eve Motion. Despite both being mapped to the Haus "motion" capability, they use different Matter clusters and should be handled accordingly in the integration code.
- **Last open/close duration** -- The Eve app tracks duration data (how long a door was open, last open time). This is implemented as Eve-specific custom attributes, not standard Matter attributes. Haus would need to compute these from the raw state change events.

## Similar Devices

- [eve-motion](eve-motion.md) -- Thread/Matter motion sensor from the same Eve platform
- [aqara-door-window-sensor](aqara-door-window-sensor.md) -- Zigbee contact sensor (requires hub, lower cost)
- [aqara-hub-m2](aqara-hub-m2.md) -- Zigbee hub for connecting Aqara contact sensors
