---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "irobot-roomba-j7-plus"
name: "iRobot Roomba j7+"
manufacturer: "iRobot (Amazon)"
brand: "iRobot"
model: "j7+"
model_aliases: ["j7 Plus", "Roomba j7+", "j755020", "j755820", "Roomba Combo j7+"]
device_type: "robot_vacuum"
category: "smart_home"
product_line: "Roomba"
release_year: 2021
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi", "bluetooth_le"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["50:14:79", "80:91:33"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [8883]
  signature_ports: [8883]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^iRobot-.*", "^Roomba-.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "irobot"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "password"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["on_off"]

# --- PROTOCOL ---
protocol:
  type: "mqtt"
  port: 8883
  transport: "TLS"
  encoding: "JSON"
  auth_method: "password"
  auth_detail: "Local MQTT on port 8883 with TLS. Username is 'user', password is a BLID-derived token retrieved from the robot via a UDP discovery packet on port 5678. Cloud API uses iRobot cloud MQTT broker at msg.irobot.com."
  base_url_template: "mqtts://{ip}:8883"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "battery"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.irobot.com/en_US/roomba-combo-j7-plus-robot-vacuum-and-mop/j755820.html"
  api_docs: ""
  developer_portal: ""
  support: "https://www.irobot.com/en_US/customer-care.html"
  community_forum: ""
  image_url: ""
  fcc_id: "2AHJ4-ONKYO"

# --- TAGS ---
tags: ["vacuum", "robot-vacuum", "mqtt", "cloud-primary", "local-mqtt-hack", "obstacle-avoidance", "auto-empty", "amazon"]
---

# iRobot Roomba j7+

## What It Is

> The iRobot Roomba j7+ is a WiFi-connected robot vacuum with PrecisionVision Navigation that uses a front-facing camera to identify and avoid obstacles like pet waste, shoes, and cords. The "+" denotes the Clean Base Automatic Dirt Disposal dock that empties the robot's bin automatically. It supports Imprint Smart Mapping for room-by-room cleaning, and can be controlled via the iRobot Home app, Alexa, or Google Assistant. Since Amazon's acquisition of iRobot, it operates entirely through iRobot's cloud infrastructure for remote control and scheduling.

## How Haus Discovers It

1. **OUI match** -- iRobot MAC prefixes: `50:14:79`, `80:91:33`
2. **Hostname pattern** -- DHCP hostname typically starts with `iRobot-` or `Roomba-`
3. **UDP discovery** -- Sending a specific UDP packet to port 5678 on the robot triggers a JSON response containing the robot's BLID (unique identifier), firmware version, SKU, and hostname
4. **Port probe** -- MQTT/TLS on port 8883 indicates local MQTT capability

## Pairing / Authentication

### Cloud API (Official)

The iRobot Home app handles all provisioning. The robot connects to iRobot's cloud MQTT broker (`msg.irobot.com:8883`). There is no official public REST API.

### Local MQTT (Community-Discovered)

The Roomba j7+ runs a local MQTT broker on port 8883 with TLS (self-signed certificate). To authenticate:

1. **Get the BLID** -- Send a UDP broadcast to port 5678. The robot responds with JSON containing its `robotid` (BLID).
2. **Get the password** -- While the robot is in provisioning mode (hold the Home button until it chimes), send a specific UDP packet to port 5678. The robot responds with a password blob.
3. **Connect via MQTT** -- Use BLID as the username and the extracted password as the MQTT password. Connect to `mqtts://{robot_ip}:8883` with TLS verification disabled (self-signed cert).

The password extraction only works during the provisioning window. Once obtained, it persists across reboots and firmware updates.

### Haus Auth Flow

`POST /api/devices/{ip}/pair` triggers the UDP discovery and password extraction sequence. The user must press and hold the Home button on the robot before initiating.

## API Reference

### Local MQTT Topics

Once connected to the local MQTT broker:

**Subscribe to state:**
```
$aws/things/{BLID}/shadow/get
```

**Publish commands:**
```
Topic: cmd
Payload: {"command": "start", "time": {unix_epoch}, "initiator": "localApp"}
```

**Available commands:**

| Command | Description |
|---------|-------------|
| `start` | Start a cleaning mission |
| `stop` | Stop the current mission |
| `pause` | Pause the current mission |
| `resume` | Resume a paused mission |
| `dock` | Send robot to dock |
| `find` | Make the robot play a sound |
| `train` | Start a training run (mapping) |

**State payload (via shadow topic):**
```json
{
  "state": {
    "reported": {
      "cleanMissionStatus": {
        "cycle": "clean",
        "phase": "run",
        "error": 0,
        "mssnM": "none",
        "nMssn": 42,
        "sqft": 850,
        "expireM": 0,
        "rechrgM": 0
      },
      "batPct": 87,
      "bin": {
        "present": true,
        "full": false
      },
      "name": "Rosie",
      "sku": "j755020",
      "softwareVer": "22.29.2",
      "lastCommand": {
        "command": "start",
        "time": 1712000000
      }
    }
  }
}
```

### Cloud MQTT (Unofficial)

The cloud broker at `msg.irobot.com:8883` uses the same MQTT topic structure but requires OAuth tokens obtained through the iRobot Home app's authentication flow. This is undocumented and subject to change.

## AI Capabilities

> AI integration is not planned for V1. If implemented, the AI concierge could:
> - Start/stop/dock cleaning missions
> - Report battery level, bin status, and current cleaning phase
> - Report room-specific cleaning history via Imprint Smart Maps
> - Alert on errors (stuck, cliff detected, bin full)

## Quirks & Notes

- **Local MQTT is unofficial** -- iRobot does not document or support the local MQTT interface; it exists because the robot runs an AWS IoT-compatible MQTT broker internally
- **Password extraction window** -- The provisioning mode for password extraction times out after about 60 seconds; if missed, the user must restart the process
- **Self-signed TLS** -- The local MQTT broker uses a self-signed certificate; TLS verification must be disabled
- **Firmware updates may break local access** -- Amazon/iRobot firmware updates could theoretically disable or change the local MQTT interface at any time
- **AWS IoT shadow format** -- The MQTT state uses AWS IoT Device Shadow format, with `reported` and `desired` state objects
- **Clean Base communication** -- The auto-empty dock communicates with the robot directly; there is no separate network entity for the Clean Base
- **BLE for initial setup** -- Bluetooth LE is used only during initial WiFi provisioning via the iRobot Home app; it is not usable for ongoing control
- **Dorita980** -- The community library `dorita980` (Node.js) and `rest980` provide the best documentation of the local MQTT protocol and password extraction process
- **Amazon cloud migration** -- Following Amazon's acquisition, iRobot may eventually migrate to Amazon's own cloud infrastructure, which could affect both cloud and local API access

## Similar Devices

> - [Roborock S8 Pro Ultra](roborock-s8-pro-ultra.md) -- Competing robot vacuum with Valetudo local control option
> - [Ecovacs Deebot X2 Omni](ecovacs-deebot-x2-omni.md) -- Competing cloud-based robot vacuum
> - [Shark AI Robot Vacuum](shark-ai-robot-vacuum.md) -- Budget competitor, cloud-only
