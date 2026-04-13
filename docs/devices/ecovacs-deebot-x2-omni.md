---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ecovacs-deebot-x2-omni"
name: "Ecovacs Deebot X2 Omni"
manufacturer: "Ecovacs Robotics"
brand: "Ecovacs"
model: "Deebot X2 Omni"
model_aliases: ["DEX86", "DEEBOT X2 OMNI", "X2 Omni"]
device_type: "robot_vacuum"
category: "smart_home"
product_line: "Deebot"
release_year: 2023
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["C8:E2:65", "70:2C:1F"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^DEEBOT.*", "^ECOVACS.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "ecovacs"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "password"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["on_off"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "password"
  auth_detail: "Cloud-based XMPP/MQTT protocol. Authentication via Ecovacs account credentials. The robot communicates with Ecovacs cloud servers (portal-*.ecouser.net). Community library 'ecovacs-deebot.js' and 'deebot_client' reverse-engineer the protocol. REST API at portal.ecouser.net for auth, MQTT for commands."
  base_url_template: "https://portal-{country}.ecouser.net"
  tls: true
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
  product_page: "https://www.ecovacs.com/us/deebot-x2-omni"
  api_docs: ""
  developer_portal: ""
  support: "https://www.ecovacs.com/us/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AYDJDEX86"

# --- TAGS ---
tags: ["vacuum", "robot-vacuum", "mop", "cloud-only", "xmpp", "mqtt", "lidar", "square-design", "auto-empty", "auto-wash-mop", "auto-refill", "hot-water-wash", "yiko-voice-assistant"]
---

# Ecovacs Deebot X2 Omni

## What It Is

> The Ecovacs Deebot X2 Omni is a premium robot vacuum and mop featuring a distinctive square design optimized for edge and corner cleaning. It offers 8000Pa suction, LiDAR navigation with AI-powered AIVI 3D 2.0 obstacle avoidance (using a structured light sensor), a OZMO Turbo 2.0 dual-spinning mop system, and a built-in voice assistant called YIKO. The all-in-one OMNI station handles auto-emptying, hot-water mop washing (up to 70C/158F), auto-refilling, and hot-air drying. All control is through the Ecovacs Home app and cloud infrastructure. There is no official local API, and the device relies entirely on Ecovacs cloud servers for remote operation.

## How Haus Discovers It

1. **OUI match** -- Ecovacs MAC prefixes: `C8:E2:65`, `70:2C:1F`
2. **Hostname pattern** -- DHCP hostname may contain `DEEBOT` or `ECOVACS`
3. **No local ports** -- The robot does not expose any local network services; all communication routes through Ecovacs cloud

## Pairing / Authentication

### Cloud API (Unofficial)

Authentication uses the Ecovacs cloud portal:

1. **Login** -- `POST https://portal-{country}.ecouser.net/api/users/user.do` with Ecovacs account email and MD5-hashed password
2. **Get auth token** -- Response includes an `accessToken` and `uid`
3. **Get device list** -- `POST https://portal-{country}.ecouser.net/api/users/user.do` with `todo=GetDeviceList`
4. **MQTT connection** -- Connect to the MQTT broker at `mq-{region}.ecouser.net:8883` with credentials derived from the auth token

Country codes: `na` (North America), `eu` (Europe), `cn` (China), etc.

### Haus Auth Flow

`POST /api/devices/{ip}/auth` with Ecovacs account email and password. Haus authenticates against the Ecovacs cloud portal and stores the access token.

## API Reference

### Cloud MQTT Commands

Once connected to the Ecovacs MQTT broker, commands are sent as JSON payloads:

**Start cleaning:**
```json
{
  "header": {
    "pri": 1,
    "ts": "1712000000000",
    "tzm": -300,
    "ver": "0.0.50"
  },
  "body": {
    "data": {
      "act": "start",
      "type": "auto"
    }
  }
}
```

**Common commands:**

| Command | act | type | Description |
|---------|-----|------|-------------|
| Auto clean | `start` | `auto` | Start full auto cleaning |
| Stop | `stop` | -- | Stop current mission |
| Return to dock | `start` | `charge` | Return to charging station |
| Spot clean | `start` | `spotArea` | Clean specific areas |
| Custom area | `start` | `customArea` | Clean custom-defined rectangle |

**State response fields:**

| Field | Description |
|-------|-------------|
| `status` | Current state (idle, cleaning, returning, charging) |
| `battery` | Battery percentage (0-100) |
| `cleanedArea` | Cleaned area in square meters |
| `cleanedTime` | Cleaning time in seconds |
| `waterLevel` | Mop water flow (1-4) |
| `vacuumPower` | Suction level (1-4) |
| `mopMode` | Mop spinning mode |

### REST API (Auth Only)

Base URL: `https://portal-{country}.ecouser.net`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/users/user.do` (todo=EcoVacs_Login) | POST | Authenticate |
| `/api/users/user.do` (todo=GetDeviceList) | POST | List devices |
| `/api/users/user.do` (todo=GetAuthCode) | POST | Get MQTT auth code |

## AI Capabilities

> AI integration is not planned for V1 due to cloud-only architecture. If implemented:
> - Start/stop/dock cleaning missions
> - Clean specific rooms by name from saved maps
> - Report battery level, cleaning area, and duration
> - Control suction power and mop water level
> - Use YIKO voice commands (relayed through cloud)

## Quirks & Notes

- **Square design** -- Unlike most round robot vacuums, the X2 Omni is square-shaped to improve edge and corner cleaning coverage by up to 99.5%
- **Cloud-only** -- No local API whatsoever; if Ecovacs cloud servers go down, the robot can only run using physical buttons
- **YIKO voice assistant** -- Built-in voice assistant allows direct voice commands to the robot without a phone; uses wake word "OK YIKO"
- **8000Pa suction** -- Among the highest suction power in consumer robot vacuums
- **AIVI 3D 2.0** -- Uses structured light sensor (not just camera) for 3D obstacle detection and avoidance
- **Hot water mop washing** -- The OMNI station washes mop pads at up to 70C (158F) for better hygiene
- **Protocol evolution** -- Ecovacs has migrated from XMPP (older models) to MQTT (newer models including X2); the X2 uses MQTT exclusively
- **Regional cloud servers** -- The robot is locked to the region configured during setup; changing regions requires a factory reset
- **Community libraries** -- `ecovacs-deebot.js` (Node.js) and `bumper` (self-hosted cloud replacement) provide unofficial protocol documentation
- **Bumper** -- An open-source project that acts as a local Ecovacs cloud replacement, but compatibility with the X2 Omni's newer MQTT protocol is limited
- **Map data** -- Supports up to 3 saved maps; map data is stored in the cloud and synced to the robot

## Similar Devices

> - [iRobot Roomba j7+](irobot-roomba-j7-plus.md) -- Competing robot vacuum with local MQTT hack
> - [Roborock S8 Pro Ultra](roborock-s8-pro-ultra.md) -- Competing robot vacuum with Valetudo local control option
> - [Shark AI Robot Vacuum](shark-ai-robot-vacuum.md) -- Budget competitor, cloud-only
