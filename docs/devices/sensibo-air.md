---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "sensibo-air"
name: "Sensibo Air"
manufacturer: "Sensibo Ltd."
brand: "Sensibo"
model: "SEN-AIR-02"
model_aliases: ["Sensibo Air", "Sensibo Sky V2", "Air Pro"]
device_type: "sensibo_ac_controller"
category: "climate"
product_line: "Sensibo"
release_year: 2021
discontinued: false
price_range: "$"

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
  mac_prefixes:
    - "CC:50:E3"        # Espressif Systems (ESP-based WiFi module)
    - "AC:67:B2"        # Espressif Systems
    - "24:62:AB"        # Espressif Systems
    - "A4:CF:12"        # Espressif Systems
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^sensibo"
    - "^Sensibo"
    - "^espressif"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "sensibo"
  polling_interval_sec: 90
  websocket_event: "sensibo:state"
  setup_type: "api_key"
  ai_chattable: true
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities:
  - "thermostat"
  - "temperature"
  - "humidity"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "API key passed as query parameter: ?apiKey={key}. Key generated from Sensibo dashboard at https://home.sensibo.com/me/api"
  base_url_template: "https://home.sensibo.com/api/v2"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "usb"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://sensibo.com/products/sensibo-air"
  api_docs: "https://sensibo.github.io/"
  developer_portal: "https://home.sensibo.com/me/api"
  support: "https://sensibo.com/pages/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AXUD-SENAIR02"

# --- TAGS ---
tags: ["cloud_only", "ir_blaster", "mini_split", "ac_controller", "sensibo", "documented_api", "rest_api"]
---

# Sensibo Air

## What It Is

The Sensibo Air is a WiFi-enabled infrared (IR) blaster that adds smart connectivity to any IR-controlled air conditioner, mini-split, heat pump, or window unit. Manufactured by Sensibo Ltd. (an Israeli company), the device mounts on the wall with line-of-sight to the AC unit and sends IR commands to replicate the original remote control. It connects via 2.4 GHz WiFi and includes a built-in temperature and humidity sensor. What sets Sensibo apart from competitors like Cielo is its **well-documented public REST API** at `home.sensibo.com/api/v2`, which makes it significantly easier to integrate with third-party systems. The API uses simple API key authentication (no OAuth flow required). Sensibo also supports Apple HomeKit, Google Home, Amazon Alexa, IFTTT, and has native Home Assistant integration.

## How Haus Discovers It

1. **OUI Match** -- Sensibo devices use Espressif (ESP-based) WiFi modules. MAC prefixes `CC:50:E3`, `AC:67:B2`, etc. identify Espressif hardware, but this is shared across many IoT products.
2. **Hostname Pattern** -- May appear as `sensibo-XXXX` or `espressif` in DHCP.
3. **No Local Probe** -- No open ports on the local network. Communication goes through Sensibo's cloud.
4. **Cloud Enrichment** -- After API key setup, Haus queries the Sensibo API for all devices and enriches locally-discovered Espressif MAC devices.

## Pairing / Authentication

Sensibo uses simple API key authentication -- no OAuth2 flow required. This is the simplest auth model of any thermostat/climate device in the Haus knowledge base.

### Getting an API Key

1. Log in to the Sensibo dashboard at `https://home.sensibo.com/me/api`.
2. Click "Generate API Key."
3. Copy the API key (a long alphanumeric string).
4. The key provides full access to all devices on the account.

### Using the API Key

Pass the API key as a query parameter on every request:

```
GET https://home.sensibo.com/api/v2/users/me/pods?apiKey={api_key}
```

There is no Bearer token, no Authorization header -- just the query parameter.

### Device Setup

1. Physical device setup is done via the Sensibo mobile app (WiFi provisioning).
2. The user selects their AC brand/model from the IR code database.
3. The app tests IR codes to confirm the correct code set.
4. Once paired, the device is accessible via the app and the REST API.

## API Reference

**Base URL:** `https://home.sensibo.com/api/v2`

**Auth:** `?apiKey={api_key}` query parameter on all requests.

The Sensibo API is well-documented at `https://sensibo.github.io/`.

### List All Devices (Pods)

```
GET /api/v2/users/me/pods?apiKey={key}
```

**Response:**
```json
{
  "status": "success",
  "result": [
    {
      "id": "ABCDEF12",
      "room": {
        "name": "Bedroom",
        "icon": "bedroom"
      },
      "acState": {
        "on": true,
        "mode": "cool",
        "targetTemperature": 24,
        "temperatureUnit": "C",
        "fanLevel": "auto",
        "swing": "rangeFull"
      },
      "measurements": {
        "temperature": 26.2,
        "humidity": 55,
        "time": {
          "secondsAgo": 42
        }
      },
      "connectionStatus": {
        "isAlive": true
      },
      "remoteCapabilities": {
        "modes": {
          "cool": {
            "temperatures": {
              "C": {"isNative": true, "values": [16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30]},
              "F": {"isNative": false, "values": [61, 63, 64, 66, 68, 70, 72, 73, 75, 77, 79, 81, 82, 84, 86]}
            },
            "fanLevels": ["quiet", "low", "medium", "high", "auto"],
            "swing": ["stopped", "rangeFull"]
          }
        }
      }
    }
  ]
}
```

**Key fields:**
- `id` -- 8-character device ID (Sensibo calls devices "pods")
- `acState` -- Current AC state as last commanded
- `measurements` -- Ambient sensor readings (temperature, humidity, time since last reading)
- `remoteCapabilities` -- What the paired AC remote supports (modes, temperature ranges, fan levels, swing options). This is critical for building UI controls.
- `connectionStatus.isAlive` -- Whether the device is online

### Get Single Device

```
GET /api/v2/pods/{pod_id}?apiKey={key}&fields=acState,measurements,connectionStatus,remoteCapabilities
```

The `fields` query parameter controls which data is returned. This reduces payload size.

### Set AC State

```
POST /api/v2/pods/{pod_id}/acStates?apiKey={key}
Content-Type: application/json

{
  "acState": {
    "on": true,
    "mode": "cool",
    "targetTemperature": 23,
    "temperatureUnit": "C",
    "fanLevel": "auto",
    "swing": "rangeFull"
  }
}
```

**Mode values:** `cool`, `heat`, `fan`, `dry`, `auto` (availability from `remoteCapabilities`)

**Fan level values:** `quiet`, `low`, `medium`, `medium_high`, `high`, `strong`, `auto` (availability from `remoteCapabilities`)

**Swing values:** `stopped`, `rangeFull`, `rangeUp`, `rangeMid`, `rangeDown`, `fixedTop`, `fixedMid`, `fixedBottom` (availability from `remoteCapabilities`)

### Update Single Property

Instead of sending the entire AC state, update a single property:

```
PATCH /api/v2/pods/{pod_id}/acStates/{property}?apiKey={key}
Content-Type: application/json

{
  "newValue": "high"
}
```

Valid properties: `on`, `mode`, `targetTemperature`, `fanLevel`, `swing`

### Turn Off

```
PATCH /api/v2/pods/{pod_id}/acStates/on?apiKey={key}
Content-Type: application/json

{"newValue": false}
```

### Get Historical Measurements

```
GET /api/v2/pods/{pod_id}/measurements?apiKey={key}&days=7
```

Returns temperature and humidity readings over the specified period. Useful for energy dashboards and trend analysis.

### Get Schedules (Timer)

```
GET /api/v2/pods/{pod_id}/schedules?apiKey={key}
```

### Climate React (Smart Mode)

Sensibo's "Climate React" automatically adjusts the AC based on temperature/humidity thresholds:

```
GET /api/v2/pods/{pod_id}/smartmode?apiKey={key}
```

```
PUT /api/v2/pods/{pod_id}/smartmode?apiKey={key}
Content-Type: application/json

{
  "enabled": true,
  "type": "temperature",
  "highTemperatureThreshold": 27,
  "highTemperatureState": {"on": true, "mode": "cool", "targetTemperature": 24, "fanLevel": "auto"},
  "lowTemperatureThreshold": 20,
  "lowTemperatureState": {"on": false}
}
```

## AI Capabilities

When the AI concierge "chats as" a Sensibo Air, it can:

- **Query ambient temperature and humidity** -- real-time sensor readings with age ("It's 26.2 C and 55% humidity in the bedroom, measured 42 seconds ago.")
- **Report AC state** -- power, mode, target temperature, fan level, swing position
- **Turn AC on/off** -- simple power toggle
- **Set target temperature** -- adjust within the AC's supported range
- **Change mode** -- cool, heat, dry, fan, auto (based on AC capabilities)
- **Adjust fan speed** -- quiet through strong, or auto
- **Report what the AC supports** -- uses `remoteCapabilities` to inform the user which modes, fan levels, and swing positions their specific AC unit supports
- **Show temperature trends** -- historical measurements for the past day/week

## Quirks & Notes

- **Best-documented API:** Sensibo has the best-documented public API among IR blaster products. The documentation at `sensibo.github.io` is comprehensive with request/response examples. This makes it significantly easier to integrate than Cielo or other competitors.
- **Simple auth:** API key in query parameter. No OAuth, no PIN flow, no token refresh. Generate a key in the dashboard and start making calls. The downside is that the key provides full account access with no scope restrictions.
- **IR one-way communication:** Like all IR blasters, Sensibo cannot verify whether the AC received and executed a command. The `acState` reflects what was sent, not what the AC is actually doing. If the user changes settings with the physical remote, the state becomes desynchronized.
- **remoteCapabilities:** The API returns what the paired AC remote supports. This is essential for building a correct UI -- different AC models support different modes, temperature ranges, fan levels, and swing positions. Always check capabilities before sending commands.
- **Temperature units:** The API supports both Celsius and Fahrenheit. The `temperatureUnit` field in acState determines which unit the `targetTemperature` is in. Sensor measurements are always in Celsius.
- **USB powered:** Requires constant 5V USB power. Includes a micro-USB cable and power adapter.
- **Pod ID format:** Device IDs are 8-character alphanumeric strings, case-sensitive.
- **ESP-based hardware:** Uses Espressif WiFi module, making MAC-based identification unreliable (shared OUIs with many devices).
- **Rate limiting:** The API has a rate limit of approximately 3 requests per 10 seconds per device. Haus should poll no more frequently than every 90 seconds.
- **HomeKit support:** The Sensibo Air supports HomeKit, which could provide a local control path independent of the cloud API.
- **Multiple Sensibo models:** Sensibo Sky (original), Sensibo Air (current), Sensibo Air Pro (with air quality sensor). The API is the same across models.

## Similar Devices

- **cielo-breez-plus** -- Competing WiFi IR blaster for AC control (less documented API)
- **mysa-smart-thermostat** -- Mysa also makes an AC controller model (different form factor)
- **ecobee-smart-thermostat-premium** -- Traditional smart thermostat (different approach, but same climate category)
