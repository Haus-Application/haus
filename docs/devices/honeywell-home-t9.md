---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "honeywell-home-t9"
name: "Honeywell Home T9 Smart Thermostat"
manufacturer: "Resideo Technologies Inc."
brand: "Honeywell Home"
model: "RCHT9610WFSW2003/U"
model_aliases: ["T9", "RCHT9610WFSW2003", "TH6320WF2003"]
device_type: "honeywell_thermostat"
category: "climate"
product_line: "Honeywell Home T-Series"
release_year: 2019
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
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "48:A2:E6"        # Resideo Technologies
    - "00:D0:2D"        # Resideo / Honeywell
    - "B0:99:28"        # Honeywell (older production)
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Honeywell"
    - "^T9-"
    - "^TCC-"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "honeywell"
  polling_interval_sec: 120
  websocket_event: "honeywell:state"
  setup_type: "oauth2"
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
  auth_method: "oauth2"
  auth_detail: "OAuth2 authorization code flow via Honeywell Developer Portal. Access tokens expire after 1800 seconds (30 min). Refresh tokens valid for 90 days."
  base_url_template: "https://api.honeywell.com/v2"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "thermostat"
  power_source: "hardwired"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.honeywellhome.com/t9-smart-thermostat-with-sensor/RCHT9610WFSW2003-U/46000.html"
  api_docs: "https://developer.honeywell.com/api-methods"
  developer_portal: "https://developer.honeywell.com/"
  support: "https://www.honeywellhome.com/support/"
  community_forum: ""
  image_url: ""
  fcc_id: "HS9-QUARTZ-01"

# --- TAGS ---
tags: ["cloud_only", "oauth2", "thermostat", "room_sensors", "resideo", "honeywell", "multi_zone"]
---

# Honeywell Home T9 Smart Thermostat

## What It Is

The Honeywell Home T9 is a WiFi-connected smart thermostat manufactured by Resideo Technologies (spun off from Honeywell International in 2018). It features a touchscreen display, support for multi-room temperature sensing via wireless room sensors (sold separately), and "Smart Room" functionality that prioritizes temperature in occupied rooms. The T9 connects via 2.4 GHz WiFi and is controlled through the Honeywell Home app and cloud API. It supports up to 20 wireless room sensors, each with temperature and motion detection. The thermostat is compatible with most 24V HVAC systems including heat pumps, forced air, hydronic, and radiant heating. There is no local API -- all programmatic control requires the Honeywell cloud.

## How Haus Discovers It

1. **OUI Match** -- Devices with MAC prefix `48:A2:E6`, `00:D0:2D`, or `B0:99:28` are flagged as Resideo/Honeywell devices.
2. **No Local Probe** -- Honeywell Home thermostats have no open ports. Haus skips HTTP fingerprinting.
3. **Hostname Pattern** -- May appear as `Honeywell-XXXX` or `TCC-XXXX` in DHCP (TCC = Total Connect Comfort, Honeywell's legacy platform name).
4. **Cloud Enrichment** -- After OAuth2 setup, Haus queries the Honeywell API for registered locations and thermostats.

## Pairing / Authentication

The Honeywell Home API uses standard OAuth2 authorization code flow.

### Developer Registration

1. Create an account at `https://developer.honeywell.com/`.
2. Register an application to obtain `consumer_key` and `consumer_secret`.
3. Set the redirect URI in the developer portal.

### OAuth2 Flow

#### Step 1: Authorization

```
GET https://api.honeywell.com/oauth2/authorize
  ?response_type=code
  &client_id={consumer_key}
  &redirect_uri={redirect_uri}
```

The user logs in with their Honeywell Home account and authorizes the application.

#### Step 2: Token Exchange

```
POST https://api.honeywell.com/oauth2/token
Content-Type: application/x-www-form-urlencoded
Authorization: Basic {base64(consumer_key:consumer_secret)}

grant_type=authorization_code&code={code}&redirect_uri={redirect_uri}
```

**Response:**
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "expires_in": 1799,
  "token_type": "Bearer"
}
```

**Important:** Access tokens expire after only 1800 seconds (30 minutes) -- much shorter than Google or Ecobee. Haus must refresh proactively.

#### Step 3: Refresh Token

```
POST https://api.honeywell.com/oauth2/token
Content-Type: application/x-www-form-urlencoded
Authorization: Basic {base64(consumer_key:consumer_secret)}

grant_type=refresh_token&refresh_token={token}
```

Refresh tokens are valid for approximately 90 days. If the refresh token expires, the user must re-authorize.

## API Reference

**Base URL:** `https://api.honeywell.com/v2`

**Auth Header:** `Authorization: Bearer {access_token}`

### List Locations

```
GET /v2/locations
  ?apikey={consumer_key}
```

Returns all locations (homes) associated with the user's account. Each location contains a list of devices.

**Response (abbreviated):**
```json
[
  {
    "locationID": 123456,
    "name": "My Home",
    "devices": [
      {
        "deviceID": "LCC-00A1B2C3D4E5",
        "name": "Living Room Thermostat",
        "deviceType": 128,
        "deviceModel": "T9-T10",
        "indoorTemperature": 72.0,
        "indoorHumidity": 45,
        "allowedModes": ["Heat", "Cool", "Auto", "Off"],
        "changeableValues": {
          "mode": "Heat",
          "heatSetpoint": 70,
          "coolSetpoint": 76,
          "thermostatSetpointStatus": "NoHold"
        },
        "operationStatus": {
          "mode": "Heat",
          "fanRequest": false,
          "circulationFanRequest": false
        }
      }
    ]
  }
]
```

**Temperature format:** Fahrenheit (integer or float). The API returns and accepts temperatures in the user's configured scale (typically Fahrenheit in the US).

### Get Single Device

```
GET /v2/devices/thermostats/{device_id}
  ?apikey={consumer_key}
  &locationId={location_id}
```

### Set Thermostat State

```
POST /v2/devices/thermostats/{device_id}
  ?apikey={consumer_key}
  &locationId={location_id}
Content-Type: application/json

{
  "mode": "Heat",
  "heatSetpoint": 72,
  "coolSetpoint": 76,
  "thermostatSetpointStatus": "TemporaryHold",
  "autoChangeoverActive": false
}
```

**Mode values:** `Heat`, `Cool`, `Auto`, `Off`

**Hold types:**
- `NoHold` -- Following schedule
- `TemporaryHold` -- Hold until next scheduled period
- `PermanentHold` -- Hold indefinitely
- `HoldUntil` -- Hold until specific time (requires `nextPeriodTime`)

### Set Fan Mode

```
POST /v2/devices/thermostats/{device_id}/fan
  ?apikey={consumer_key}
  &locationId={location_id}
Content-Type: application/json

{
  "mode": "On",
  "allowedModes": ["On", "Auto", "Circulate"]
}
```

### Room Sensors

Room sensor data is included in the thermostat device response when sensors are paired. Each sensor reports temperature and occupancy.

```json
{
  "groups": [
    {
      "id": 1,
      "name": "Smart Room 1",
      "rooms": [
        {
          "id": 1,
          "name": "Bedroom",
          "type": "Bedroom",
          "avgTemperature": 71.5,
          "accessories": [
            {
              "id": 1,
              "name": "Bedroom Sensor",
              "type": "IndoorAirSensor",
              "status": "Alive",
              "temperature": 71.5,
              "motion": true
            }
          ]
        }
      ]
    }
  ]
}
```

## AI Capabilities

When the AI concierge "chats as" a Honeywell Home T9, it can:

- **Query temperature and humidity** -- current readings from the thermostat and all room sensors
- **Report HVAC mode** -- current mode and operational status
- **Set temperature** -- adjust heat/cool setpoints with hold type control
- **Change mode** -- switch between Heat, Cool, Auto, and Off
- **Report room sensor data** -- temperature and occupancy from each paired sensor
- **Report which rooms are occupied** -- using motion data from room sensors

## Quirks & Notes

- **Short-lived tokens:** Access tokens expire after only 30 minutes, much shorter than Google (1 hour) or Ecobee (1 hour). Haus must implement aggressive token refresh.
- **API key in query string:** The Honeywell API requires the `apikey` (consumer_key) as a query parameter on every request, in addition to the Bearer token in the Authorization header. This is unusual.
- **Location-scoped devices:** Devices are scoped to locations. Most API calls require both `device_id` and `location_id`.
- **Temperature scale:** The API returns temperatures in the user's configured scale (Fahrenheit or Celsius) rather than a fixed scale. Haus must detect and handle both.
- **Rate limiting:** The Honeywell API enforces rate limits but does not clearly document them. Community reports suggest approximately 10-20 requests per minute. Haus should poll conservatively (2-minute intervals).
- **Resideo transition:** Honeywell Home products are manufactured by Resideo Technologies, which spun off from Honeywell International in 2018. The API and developer portal remain under the honeywell.com domain.
- **Device model confusion:** The API returns `deviceModel: "T9-T10"` for both T9 and T10 thermostats. Distinguishing between models requires additional heuristics.
- **Refresh token expiry:** Unlike Google (no expiry) or Ecobee (~1 year), Honeywell refresh tokens expire after ~90 days. Haus must handle re-authorization gracefully.
- **Room sensors:** Wireless 915 MHz room sensors (sold separately) pair with the T9 directly. They are not WiFi or Zigbee -- they use a proprietary protocol. Sensor data is only accessible through the thermostat's cloud API.
- **No local API:** Like Nest, there is absolutely no local control available. The thermostat continues running its schedule without internet but cannot be controlled remotely.

## Similar Devices

- **nest-learning-thermostat** -- Google's competing premium thermostat with SDM cloud API
- **nest-thermostat-2020** -- Google's budget thermostat, similar price point
- **ecobee-smart-thermostat-premium** -- Ecobee's competing thermostat with cloud + local HomeKit
