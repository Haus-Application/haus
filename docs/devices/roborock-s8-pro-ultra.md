---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "roborock-s8-pro-ultra"
name: "Roborock S8 Pro Ultra"
manufacturer: "Beijing Roborock Technology Co., Ltd."
brand: "Roborock"
model: "S8 Pro Ultra"
model_aliases: ["S8ProUltra", "roborock.vacuum.a70", "S8PU"]
device_type: "robot_vacuum"
category: "smart_home"
product_line: "Roborock S Series"
release_year: 2023
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["64:90:C1", "78:11:DC", "50:EC:50"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^roborock-.*", "^roborock_.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "roborock"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["on_off"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 58867
  transport: "UDP"
  encoding: "binary"
  auth_method: "api_key"
  auth_detail: "Uses Xiaomi Mi Home / Tuya-derived protocol. Cloud token obtained from Roborock app login. Local communication uses miIO protocol on UDP port 54321 with AES-128-CBC encryption using a device-specific token. Valetudo replaces the cloud firmware for full local REST control."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "battery"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://us.roborock.com/pages/roborock-s8-pro-ultra"
  api_docs: ""
  developer_portal: ""
  support: "https://us.roborock.com/pages/support"
  community_forum: "https://forum.roborock.com"
  image_url: ""
  fcc_id: "2AQMU-S8PU"

# --- TAGS ---
tags: ["vacuum", "robot-vacuum", "mop", "cloud-primary", "valetudo-compatible", "miio-protocol", "lidar", "auto-empty", "auto-wash-mop", "auto-refill"]
---

# Roborock S8 Pro Ultra

## What It Is

> The Roborock S8 Pro Ultra is a premium robot vacuum and mop with a self-maintaining dock that auto-empties the dustbin, auto-washes the mop pad with hot water, auto-refills the clean water tank, and auto-dries the mop. It features dual rubber brushes, LiDAR navigation with 3D structured light obstacle avoidance (Reactive 3D), 6000Pa suction, and a liftable VibraRise 2.0 mopping system that raises the mop pad 5mm when vacuuming carpets. The robot communicates primarily through Roborock's cloud infrastructure, but an alternative firmware project called Valetudo can provide full local control.

## How Haus Discovers It

1. **OUI match** -- Roborock MAC prefixes: `64:90:C1`, `78:11:DC`, `50:EC:50`
2. **Hostname pattern** -- DHCP hostname typically contains `roborock`
3. **miIO discovery** -- Sending a Hello packet (all `0xFF`, 32 bytes) to UDP port 54321 triggers a response containing the device's ID, uptime, and token checksum
4. **Valetudo detection** -- If Valetudo is installed, the robot serves an HTTP REST API on port 80 with an OpenAPI spec at `/api/v2`

## Pairing / Authentication

### Stock Firmware (Cloud)

The Roborock app handles all provisioning via BLE. The robot connects to Roborock's cloud servers. There is no official local API.

### Obtaining the Device Token (miIO Protocol)

For local miIO protocol access on stock firmware:

1. **From the Roborock app** -- Extract the token from the app's local database or via API interception
2. **From provisioning** -- During WiFi setup, the device token is briefly available in plaintext via the miIO Hello response
3. **From rooted device** -- Token stored in `/mnt/data/miio/device.token`

The device token is a 128-bit key used for AES-128-CBC encryption of all miIO protocol messages.

### Valetudo (Full Local Control)

Valetudo is a third-party firmware replacement that:

1. Replaces cloud connectivity with a local HTTP REST API
2. Runs on port 80 with no authentication by default (can be configured)
3. Provides full vacuum control, map data, and configuration
4. Requires rooting the device (varies by model and firmware version)

**Root methods for S8 Pro Ultra:**
- Firmware version dependent
- Community tools like `dustbuilder` generate custom firmware packages
- Rooting typically requires USB access or exploit-based methods
- Roborock has been actively patching root exploits in newer firmware versions

### Haus Auth Flow

For stock firmware: `POST /api/devices/{ip}/auth` with the miIO device token.
For Valetudo: `POST /api/devices/{ip}/pair` auto-detects the Valetudo REST API.

## API Reference

### miIO Protocol (Stock Firmware)

UDP port 54321, AES-128-CBC encrypted JSON-RPC:

**Start cleaning:**
```json
{"method": "app_start", "params": [], "id": 1}
```

**Stop cleaning:**
```json
{"method": "app_stop", "params": [], "id": 2}
```

**Return to dock:**
```json
{"method": "app_charge", "params": [], "id": 3}
```

**Pause:**
```json
{"method": "app_pause", "params": [], "id": 4}
```

**Get status:**
```json
{"method": "get_status", "params": [], "id": 5}
```

**Status response:**
```json
{
  "result": [{
    "msg_ver": 2,
    "msg_seq": 1,
    "state": 8,
    "battery": 100,
    "clean_time": 1200,
    "clean_area": 25500000,
    "error_code": 0,
    "map_present": 1,
    "in_cleaning": 0,
    "in_returning": 0,
    "in_fresh_state": 1,
    "water_box_status": 1,
    "fan_power": 102,
    "mop_mode": 300
  }]
}
```

**State codes:** 1=initiating, 2=sleeping, 3=idle, 5=cleaning, 6=returning, 8=charging, 9=error, 10=paused, 11=spot cleaning, 14=updating, 15=docking

### Valetudo REST API (Rooted Firmware)

Base URL: `http://{ip}/api/v2`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/robot/state` | GET | Current state (status, battery, area) |
| `/robot/capabilities` | GET | List all supported capabilities |
| `/robot/capabilities/BasicControlCapability` | PUT | Start/stop/pause/home commands |
| `/robot/capabilities/FanSpeedControlCapability/preset` | PUT | Set suction level |
| `/robot/capabilities/WaterUsageControlCapability/preset` | PUT | Set mop water flow |
| `/robot/capabilities/MapSegmentationCapability` | PUT | Clean specific rooms |
| `/map/latest` | GET | Current map data (rooms, walls, paths) |

**Start cleaning (Valetudo):**
```
PUT /api/v2/robot/capabilities/BasicControlCapability
Content-Type: application/json
{"action": "start"}
```

**Actions:** `start`, `stop`, `pause`, `home`

## AI Capabilities

> AI integration is not planned for V1. If implemented, particularly with Valetudo:
> - Start/stop/dock cleaning missions
> - Clean specific rooms by name
> - Report battery level, cleaning area, and duration
> - Display live map with robot position
> - Adjust suction power and mop water flow
> - Report consumable status (brush, filter, mop pad)

## Quirks & Notes

- **miIO protocol encryption** -- All local communication on stock firmware is AES-128-CBC encrypted with the device token; messages include a checksum and timestamp
- **Valetudo compatibility varies** -- Roborock actively works to prevent rooting; newer firmware versions may not be rootable; check Valetudo's supported devices list before purchasing
- **S8 Pro Ultra root status** -- As of 2024, rooting requires specific firmware versions; the community maintains a compatibility matrix at valetudo.cloud
- **Dual rubber brushes** -- The S8 Pro Ultra has two rubber main brushes that counter-rotate, reducing hair tangles
- **6000Pa suction** -- Highest suction in the S8 lineup, automatically increases on carpets
- **VibraRise 2.0** -- Mop lifts 5mm when carpet is detected; sonic vibration at 3000 RPM for mopping
- **Empty Wash Fill dock** -- The dock handles dustbin emptying (2.5L bag), mop washing (hot water), water tank refilling, and mop drying (hot air)
- **Map storage** -- Stores up to 4 maps for multi-floor homes
- **miIO Hello packet** -- 32 bytes of 0xFF sent to UDP 54321 triggers a discovery response with device ID and uptime
- **Cloud servers** -- Roborock operates regional cloud servers; the robot connects to the server assigned during initial setup (US, EU, CN, etc.)

## Similar Devices

> - [iRobot Roomba j7+](irobot-roomba-j7-plus.md) -- Competing robot vacuum with local MQTT hack
> - [Ecovacs Deebot X2 Omni](ecovacs-deebot-x2-omni.md) -- Competing cloud-based robot vacuum/mop
> - [Shark AI Robot Vacuum](shark-ai-robot-vacuum.md) -- Budget competitor, cloud-only
