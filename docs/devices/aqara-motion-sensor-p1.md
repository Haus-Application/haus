---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "aqara-motion-sensor-p1"
name: "Aqara Motion Sensor P1"
manufacturer: "Lumi United Technology Co., Ltd."
brand: "Aqara"
model: "MS-S02"
model_aliases: ["RTCGQ14LM", "Aqara Motion Sensor P1", "Aqara Motion & Light Sensor P1"]
device_type: "aqara_motion_sensor"
category: "security"
product_line: "Aqara"
release_year: 2022
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
  mdns_services: []
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
capabilities: ["motion"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "Zigbee"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Zigbee 3.0 network join. Paired to a Zigbee coordinator hub. Uses ZCL Occupancy Sensing cluster (0x0406) for motion detection and Illuminance Measurement cluster (0x0400) for light level. Xiaomi proprietary cluster (0xFCC0) for sensitivity and detection interval settings."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "sensor"
  power_source: "battery"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee"]

# --- LINKS ---
links:
  product_page: "https://www.aqara.com/us/motion-sensor-p1.html"
  api_docs: ""
  developer_portal: "https://developer.aqara.com/"
  support: "https://www.aqara.com/us/support.html"
  community_forum: "https://community.aqara.com/"
  image_url: ""
  fcc_id: "2AKIT-MS-S02"

# --- TAGS ---
tags: ["zigbee", "battery", "motion_sensor", "pir", "light_level", "aqara", "requires_hub", "zigbee_3_0", "adjustable_sensitivity"]
---

# Aqara Motion Sensor P1

## What It Is

The Aqara Motion Sensor P1 is a Zigbee 3.0 passive infrared (PIR) motion sensor with an integrated light level (illuminance) sensor. It detects human movement within a 170-degree field of view at distances up to 7 meters and simultaneously reports ambient light levels in lux. The P1 is a significant upgrade over Aqara's original motion sensor -- it offers three adjustable sensitivity levels (low, medium, high), configurable detection intervals (from 2 seconds to 65535 seconds, compared to the original's fixed 60-second re-trigger delay), and Zigbee 3.0 compatibility for broader hub support. It runs on two CR2450 batteries with an estimated 5-year battery life and requires a Zigbee hub to function.

## How Haus Discovers It

Like all Aqara Zigbee sensors, the Motion Sensor P1 has no IP presence and cannot be discovered directly on the network:

1. **Hub discovery** -- Haus discovers the Zigbee hub (Aqara Hub M2, SmartThings, Zigbee2MQTT coordinator) on the local network via mDNS, SSDP, or OUI matching.
2. **Hub API query** -- Haus queries the hub's device list. The P1 appears with Zigbee model identifier `lumi.motion.ac02`.
3. **Cluster inspection** -- The device exposes ZCL Occupancy Sensing cluster (0x0406) for motion events and Illuminance Measurement cluster (0x0400) for light levels. Xiaomi proprietary cluster (0xFCC0) provides sensitivity settings, detection interval, and battery information.
4. **Device type classification** -- Haus maps the model identifier to its internal `aqara_motion_sensor` type.

## Pairing / Authentication

No authentication on the sensor itself. Pairing is performed at the Zigbee network level:

1. **Put hub in pairing mode** -- Enable device discovery on the Zigbee coordinator.
2. **Reset sensor** -- Press and hold the reset button on the bottom of the sensor for 5 seconds until the blue LED blinks rapidly.
3. **Zigbee join** -- The sensor joins the Zigbee network and begins reporting.
4. **Configure settings** -- After pairing, sensitivity (1=low, 2=medium, 3=high) and detection interval (in seconds) can be configured via the proprietary cluster 0xFCC0, attribute 0x010C (sensitivity) and attribute 0x0102 (detection interval).

## API Reference

The sensor has no direct API. Data is accessed through the parent hub.

### Via Aqara Hub (Aqara Developer API)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v3.0/open/device/query` | POST | Query device list including motion sensors |
| `/v3.0/open/resource/query` | POST | Query resource (state) of a specific device |

Resource model:

| Resource ID | Description | Values |
|-------------|-------------|--------|
| `3.1.85` | Motion status | `0` = no motion, `1` = motion detected |
| `0.3.85` | Illuminance (lux) | Integer, 0-83000 |
| `8.0.2008` | Battery voltage (mV) | Integer |
| `0.5.85` | Detection interval (seconds) | Integer, 2-65535 |
| `0.4.85` | Sensitivity | `1` = low, `2` = medium, `3` = high |

### Via Zigbee2MQTT

MQTT payload on `zigbee2mqtt/{friendly_name}`:

```json
{
  "occupancy": true,
  "illuminance": 156,
  "illuminance_lux": 156,
  "motion_sensitivity": "medium",
  "detection_interval": 15,
  "battery": 100,
  "voltage": 3100,
  "temperature": 23
}
```

### Via SmartThings

- Capability: `motionSensor` -- `motion.active` / `motion.inactive`
- Capability: `illuminanceMeasurement` -- `illuminance` (lux)
- Capability: `battery` -- `battery` (percentage)

## AI Capabilities

When integrated, the AI concierge could:

- Report current motion status for any room with a P1 sensor
- Report ambient light levels in lux, with human-readable descriptions ("dark", "dim", "bright")
- Trigger automations based on motion + light level combinations (e.g., "turn on lights when motion detected and light level below 50 lux")
- Provide occupancy history and patterns over time
- Adjust sensitivity and detection interval via hub API

## Quirks & Notes

- **Configurable detection interval** -- The P1's killer feature versus the original Aqara motion sensor. The original had a fixed 60-second re-trigger cooldown. The P1 can be set as low as 2 seconds via the proprietary Zigbee attribute, making it suitable for lighting automations that need fast re-trigger.
- **Sensitivity adjustment** -- Three levels (low, medium, high) configurable via Zigbee attribute 0x010C on cluster 0xFCC0. High sensitivity increases range but also false positive rate from pets.
- **CR2450 batteries** -- Uses two CR2450 coin cells (not CR2032). Battery life is rated at 5 years, which is excellent.
- **170-degree FOV** -- Wide detection angle. The sensor can be wall-mounted or placed on a shelf with the included magnetic mount.
- **Light level sensor** -- Reports illuminance in lux. Useful for conditional automations. Note that the light sensor is on the front face of the sensor and can be affected by sensor placement.
- **Zigbee sleepy end device** -- Wakes on PIR trigger or periodically for heartbeat. Cannot be polled on demand.
- **Xiaomi proprietary attributes** -- Extended features (sensitivity, interval) are on the proprietary cluster 0xFCC0 and require hub firmware or Zigbee2MQTT device definitions that understand these custom attributes.
- **Thread/Matter variant coming** -- Aqara has announced Thread-based motion sensors (FP2 is WiFi-based presence, but future PIR sensors may support Thread). The P1 is Zigbee-only.

## Similar Devices

- [aqara-door-window-sensor](aqara-door-window-sensor.md) -- Contact sensor from the same Aqara Zigbee ecosystem
- [aqara-hub-m2](aqara-hub-m2.md) -- Aqara's WiFi/Zigbee gateway for connecting these sensors
- [eve-motion](eve-motion.md) -- Thread/Matter motion sensor, no hub required
- [smartthings-hub-v3](smartthings-hub-v3.md) -- Multi-protocol hub compatible with Aqara Zigbee sensors
