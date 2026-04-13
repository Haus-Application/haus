---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "aqara-door-window-sensor"
name: "Aqara Door and Window Sensor"
manufacturer: "Lumi United Technology Co., Ltd."
brand: "Aqara"
model: "MCCGQ11LM"
model_aliases: ["MCCGQ12LM", "MCCGQ14LM", "Aqara Door & Window Sensor", "Aqara Contact Sensor"]
device_type: "aqara_contact_sensor"
category: "security"
product_line: "Aqara"
release_year: 2018
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["zigbee"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []            # Zigbee device, no WiFi/Ethernet MAC
  mdns_services: []           # Zigbee-only, no mDNS
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "aqara"
  polling_interval_sec: 0     # Event-driven via Zigbee reports
  websocket_event: ""
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["motion"]      # open/close detection modeled as motion

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "Zigbee"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Zigbee 3.0 network join. Device must be paired to a Zigbee coordinator (Aqara Hub, SmartThings Hub, or Zigbee2MQTT). Uses Zigbee Cluster Library (ZCL) IAS Zone cluster (0x0500) for open/close status and Xiaomi proprietary cluster (0xFCC0) for extended attributes."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "sensor"
  power_source: "battery"
  mounting: "door"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee"]

# --- LINKS ---
links:
  product_page: "https://www.aqara.com/us/door-and-window-sensor.html"
  api_docs: ""
  developer_portal: "https://developer.aqara.com/"
  support: "https://www.aqara.com/us/support.html"
  community_forum: "https://community.aqara.com/"
  image_url: ""
  fcc_id: "2AKIT-MCCGQ11LM"

# --- TAGS ---
tags: ["zigbee", "battery", "contact_sensor", "door_window", "xiaomi", "aqara", "requires_hub", "zigbee_3_0"]
---

# Aqara Door and Window Sensor

## What It Is

The Aqara Door and Window Sensor is a compact, battery-powered Zigbee 3.0 contact sensor that detects whether a door or window is open or closed. It consists of two pieces -- a main body containing a reed switch and electronics, and a small magnet that attaches to the moving part (door/window). When the magnet separates from the sensor body, the device reports an "open" state. It runs on a single CR1632 coin cell battery with roughly two years of battery life, and it communicates exclusively over Zigbee, requiring a Zigbee hub (Aqara Hub M1S/M2/E1, SmartThings, or a Zigbee2MQTT coordinator) to function. There is no WiFi, no Bluetooth, and no IP connectivity on the sensor itself.

## How Haus Discovers It

Haus cannot discover this sensor directly on the network because it is a Zigbee end device with no IP presence. Discovery follows this path:

1. **Hub discovery** -- Haus first discovers the Zigbee hub (Aqara Hub M2, SmartThings, etc.) via mDNS, SSDP, or OUI matching on the local network.
2. **Hub API query** -- Once authenticated with the hub, Haus queries the hub's device list. The Aqara Door and Window Sensor appears as a child device with model identifier `lumi.sensor_magnet.aq2` (or `lumi.sensor_magnet` for the original model).
3. **Zigbee cluster inspection** -- The device reports ZCL IAS Zone cluster (0x0500) for open/close state. Xiaomi's proprietary cluster (0xFCC0) provides battery voltage and temperature.
4. **Device type classification** -- Haus maps the Zigbee model identifier to its internal `aqara_contact_sensor` type.

## Pairing / Authentication

The sensor itself has no authentication. Pairing occurs at the Zigbee network level:

1. **Put hub in pairing mode** -- In the Aqara Home app (or whatever hub manages the Zigbee network), enter device pairing mode.
2. **Reset sensor** -- Press and hold the small reset button on the side of the sensor with a pin for 5 seconds until the blue LED flashes three times.
3. **Zigbee join** -- The sensor sends a Zigbee association request and joins the network. The hub confirms the pairing.
4. **Haus reads from hub** -- Haus then reads the sensor's state through the hub's API. No direct authentication with the sensor is needed.

## API Reference

The sensor has no API of its own. Its data is accessed through the hub it is paired with. Depending on the hub:

### Via Aqara Hub (Aqara Developer API)

The Aqara Developer Cloud API exposes sensor data:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v3.0/open/device/query` | POST | Query device list including contact sensors |
| `/v3.0/open/resource/query` | POST | Query resource (state) of a specific device |

Resource model for the contact sensor:

| Resource ID | Description | Values |
|-------------|-------------|--------|
| `3.1.85` | Contact status | `0` = closed, `1` = open |
| `8.0.2008` | Battery voltage (mV) | Integer |

### Via SmartThings Hub

SmartThings REST API exposes the sensor via the standard capability model:

- Capability: `contactSensor` -- `contact.open` / `contact.closed`
- Capability: `battery` -- `battery` (percentage)

### Via Zigbee2MQTT

The sensor exposes MQTT topics:

- `zigbee2mqtt/{friendly_name}` -- JSON payload with `contact` (boolean), `battery` (percentage), `voltage` (mV), `temperature` (Celsius)

## AI Capabilities

AI integration is minimal for a binary contact sensor. When implemented, the AI concierge could:

- Report whether a specific door or window is currently open or closed
- Provide history of open/close events with timestamps
- Alert when a door has been left open for a configurable duration
- Include door/window state in security status summaries

## Quirks & Notes

- **Zigbee sleepy end device** -- This sensor spends most of its time asleep to conserve battery. It wakes on state change (open/close) or every ~50 minutes for a heartbeat report. You cannot poll it on demand.
- **CR1632 battery** -- Not the most common coin cell. Battery life is roughly 2 years under normal use. The sensor reports battery voltage via Xiaomi's proprietary cluster.
- **Temperature reporting** -- The sensor has an internal temperature sensor (used for battery voltage compensation) that also reports ambient temperature. Accuracy is rough (plus or minus 2 degrees Celsius) since it is inside a small plastic enclosure.
- **Xiaomi proprietary attributes** -- Aqara/Xiaomi Zigbee devices use proprietary cluster 0xFCC0 and custom attributes that are not part of the standard ZCL spec. Hub firmware or Zigbee2MQTT must understand these.
- **Zigbee range** -- As a battery-powered end device, it does not act as a Zigbee router. Range is typically 10-20 meters indoors depending on walls and interference.
- **Multiple revisions** -- The MCCGQ11LM is the most common model. The MCCGQ14LM is a newer revision with identical functionality but slightly updated hardware. Both use the same Zigbee clusters.
- **No direct IP integration** -- Haus must always go through a hub to access this sensor. There is no way to communicate with it over WiFi or Ethernet.

## Similar Devices

- [aqara-motion-sensor-p1](aqara-motion-sensor-p1.md) -- PIR motion sensor from the same Aqara Zigbee ecosystem
- [aqara-hub-m2](aqara-hub-m2.md) -- Aqara's WiFi/Zigbee gateway that connects these sensors to the network
- [eve-door-window](eve-door-window.md) -- Thread/Matter contact sensor that does not require a proprietary hub
- [smartthings-hub-v3](smartthings-hub-v3.md) -- Multi-protocol hub that can pair Aqara Zigbee sensors
