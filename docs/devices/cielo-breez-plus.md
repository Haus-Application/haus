---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "cielo-breez-plus"
name: "Cielo Breez Plus"
manufacturer: "Cielo WiGle Inc."
brand: "Cielo"
model: "CBP-01"
model_aliases: ["Cielo Breez Plus", "CBP01", "Breez Plus"]
device_type: "cielo_ac_controller"
category: "climate"
product_line: "Cielo Breez"
release_year: 2020
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
    - "^cielo"
    - "^Cielo"
    - "^espressif"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "cielo"
  polling_interval_sec: 120
  websocket_event: "cielo:state"
  setup_type: "api_key"
  ai_chattable: true
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities:
  - "thermostat"
  - "temperature"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "Cloud API via api.cielowigle.com. API key authentication. Developer access available through Cielo developer portal."
  base_url_template: "https://api.cielowigle.com/1"
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
  product_page: "https://www.cielowigle.com/products/cielo-breez-plus"
  api_docs: "https://developer.cielowigle.com/"
  developer_portal: "https://developer.cielowigle.com/"
  support: "https://www.cielowigle.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AU5A-CBP01"

# --- TAGS ---
tags: ["cloud_only", "ir_blaster", "mini_split", "ac_controller", "cielo", "hvac", "universal_remote"]
---

# Cielo Breez Plus

## What It Is

The Cielo Breez Plus is a WiFi-enabled infrared (IR) blaster designed to make any IR-controlled air conditioner, mini-split, heat pump, or window unit "smart." Manufactured by Cielo WiGle Inc. (a Pakistani-American company based in Redmond, WA), the device sits on a wall or shelf with line-of-sight to the AC unit's IR receiver and sends IR commands that mimic the original remote control. It connects to the home network via 2.4 GHz WiFi and is controlled through the Cielo Home app and cloud API. The Breez Plus includes a built-in temperature and humidity sensor for ambient monitoring. It supports over 500 AC brands with a database of IR codes. The device is powered via micro-USB and includes scheduling, geofencing, comfy mode (automatic temperature range control), and weekly usage reports.

## How Haus Discovers It

1. **OUI Match** -- Cielo devices use Espressif (ESP-based) WiFi modules. MAC prefixes like `CC:50:E3`, `AC:67:B2` identify Espressif hardware, but this is a broad match shared with many IoT devices.
2. **Hostname Pattern** -- May appear as `cielo-XXXX` or `espressif` in DHCP.
3. **No Local Probe** -- No open ports on the local network. Communication goes through Cielo's cloud servers.
4. **Cloud Enrichment** -- After API setup, Haus queries the Cielo API for registered devices and their current state.

**Note:** Like Mysa, positive network identification is difficult due to generic Espressif OUIs. Cloud API integration or user confirmation is needed.

## Pairing / Authentication

### Cielo Developer API

Cielo offers a developer API through their developer portal at `https://developer.cielowigle.com/`.

#### Registration

1. Create a Cielo Home account (via the mobile app).
2. Apply for API access at the developer portal.
3. Receive an API key for authentication.

#### Authentication

```
Authorization: Bearer {api_key}
```

All API requests include the API key as a Bearer token.

### Device Setup

1. The physical device is set up via the Cielo Home mobile app (WiFi provisioning via Bluetooth or AP mode).
2. The user selects their AC brand and model from the IR code database.
3. The app tests IR codes to confirm correct operation.
4. Once paired, the device is accessible via both the app and the API.

## API Reference

**Base URL:** `https://api.cielowigle.com/1` (version may vary)

**Auth Header:** `Authorization: Bearer {api_key}`

### List Devices

```
GET /devices
```

Returns all devices registered to the account.

**Response (abbreviated):**
```json
{
  "devices": [
    {
      "deviceId": "abc123",
      "macAddress": "CC:50:E3:XX:XX:XX",
      "name": "Bedroom AC",
      "brand": "Mitsubishi",
      "model": "MSZ-FH12NA",
      "isOnline": true,
      "lastSeen": "2026-04-12T10:00:00Z",
      "ambientTemperature": 75.0,
      "ambientHumidity": 48,
      "deviceState": {
        "power": "on",
        "mode": "cool",
        "targetTemperature": 72,
        "fanSpeed": "auto",
        "swingMode": "auto"
      }
    }
  ]
}
```

### Control Device

```
POST /devices/{deviceId}/control
Content-Type: application/json

{
  "power": "on",
  "mode": "cool",
  "targetTemperature": 72,
  "fanSpeed": "high",
  "swingMode": "auto"
}
```

**Power values:** `on`, `off`

**Mode values:** `cool`, `heat`, `auto`, `dry`, `fan` (availability depends on the AC unit)

**Fan speed values:** `auto`, `low`, `medium`, `high`, `turbo` (availability depends on the AC unit)

**Swing mode values:** `auto`, `position1` through `position6`, `swing`, `off` (availability depends on the AC unit)

### Get Device State

```
GET /devices/{deviceId}/state
```

Returns current device state including power, mode, temperature, fan, and swing settings, plus ambient temperature and humidity from the built-in sensor.

### Get Usage Data

```
GET /devices/{deviceId}/usage
  ?startDate=2026-04-01
  &endDate=2026-04-12
```

Returns daily usage statistics including hours of operation per mode.

## AI Capabilities

When the AI concierge "chats as" a Cielo Breez Plus, it can:

- **Query ambient temperature and humidity** -- readings from the built-in sensor
- **Report AC state** -- power, mode, target temperature, fan speed, swing position
- **Turn AC on/off** -- power control
- **Set temperature** -- adjust target temperature
- **Change mode** -- cool, heat, auto, dry, fan
- **Adjust fan speed** -- auto, low, medium, high, turbo
- **Report usage** -- operating hours and energy estimates

## Quirks & Notes

- **IR line-of-sight required:** The device must have direct line-of-sight to the AC unit's IR receiver. If the IR signal is blocked, commands won't reach the AC. The Breez Plus has a ~30-foot IR range with a 120-degree coverage angle.
- **No feedback from AC:** IR communication is one-way. The Cielo device sends commands but cannot verify whether the AC actually received and executed them. If the user changes settings with the physical remote, the Cielo's state becomes out of sync until the next sensor reading provides indirect feedback (e.g., temperature changes confirm the AC is running).
- **AC brand compatibility:** The device supports 500+ brands via an IR code database. Some obscure or very old AC models may not be supported. Testing IR codes during setup is critical.
- **USB powered:** The device requires constant USB power (5V/1A). It does not have a battery. A USB power adapter is included.
- **Temperature sensor accuracy:** The built-in sensor reads ambient temperature near the device (typically mounted on the wall near the AC). This may differ from the room's average temperature, especially if mounted near the AC's output airflow.
- **Cloud dependency:** All control goes through Cielo's cloud servers. If internet goes down, the Cielo device cannot be controlled via the app or API. The user must use the original AC remote.
- **Comfy Mode:** Cielo's "Comfy Mode" automatically adjusts the AC to maintain a temperature within a user-defined range. This runs in the cloud, not on the device.
- **ESP-based hardware:** Like Mysa, uses Espressif WiFi modules, making network-level identification unreliable.
- **Multiple Breez models:** Cielo offers Breez Lite (basic), Breez Plus (with humidity), Breez Max (with display), and Breez Eco (budget). The API is similar across models but capabilities vary.

## Similar Devices

- **sensibo-air** -- Competing WiFi IR blaster for AC control, with documented REST API
- **mysa-smart-thermostat** -- Mysa also makes an AC controller model (MY300WMN), similar concept
- **honeywell-home-t9** -- Traditional thermostat (different approach to climate control)
