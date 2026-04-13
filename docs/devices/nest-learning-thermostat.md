---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "nest-learning-thermostat"
name: "Google Nest Learning Thermostat (4th Gen)"
manufacturer: "Google LLC"
brand: "Google Nest"
model: "T4000ES"
model_aliases: ["T4001ES", "T4000EF", "A0044"]
device_type: "nest_thermostat"
category: "climate"
product_line: "Nest"
release_year: 2024
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
  protocols_spoken: ["wifi", "bluetooth", "thread", "matter"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "18:B4:30"        # Google (Nest Labs Inc.)
    - "64:16:66"        # Google (Nest devices)
    - "7C:10:15"        # Google LLC
    - "B0:09:DA"        # Google LLC
    - "F8:0F:F9"        # Google LLC
    - "18:7F:88"        # Google LLC
    - "48:D6:D5"        # Google LLC
  mdns_services: []     # Nest thermostats do not advertise mDNS services
  mdns_txt_keys: []
  default_ports: []     # No open ports -- cloud-only device
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Nest-[A-Z0-9]+"
    - "^Google-Nest"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP interface

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "nest"
  polling_interval_sec: 60
  websocket_event: "nest:state"
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
  auth_detail: "OAuth2 via nestservices.google.com partner connection flow. Access token in Authorization: Bearer header. Tokens expire after 1 hour; refresh tokens do not expire unless revoked."
  base_url_template: "https://smartdevicemanagement.googleapis.com/v1/enterprises/{project_id}"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "thermostat"
  power_source: "hardwired"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le", "thread"]

# --- LINKS ---
links:
  product_page: "https://store.google.com/product/nest_learning_thermostat_4th_gen"
  api_docs: "https://developers.google.com/nest/device-access/api"
  developer_portal: "https://console.nest.google.com/device-access"
  support: "https://support.google.com/googlenest/answer/9241211"
  community_forum: "https://www.googlenestcommunity.com/"
  image_url: ""
  fcc_id: "A4RA0044"

# --- TAGS ---
tags: ["cloud_only", "sdm_api", "oauth2", "thermostat", "learning", "thread", "matter", "google", "nest"]
---

# Google Nest Learning Thermostat (4th Gen)

## What It Is

The Google Nest Learning Thermostat (4th generation, model T4000ES) is Google's flagship smart thermostat, released in 2024. It features a large 2.7" borderless LCD display, a redesigned polished stainless steel housing, and integrates temperature, humidity, and occupancy sensing. The 4th gen is the first Nest thermostat to support Matter and Thread protocols, meaning it can participate in a Thread mesh network and be controlled by any Matter-compatible platform. It retains the hallmark "learning" capability -- it observes your manual temperature adjustments over the first week or two and builds a schedule automatically. The thermostat connects to the home network over WiFi (2.4 GHz) and communicates with Google's cloud servers; there is no local control API. All programmatic control goes through the Google Smart Device Management (SDM) API.

## How Haus Discovers It

1. **OUI Match** -- During a network scan, devices with MAC prefixes `18:B4:30`, `64:16:66`, `7C:10:15`, `B0:09:DA`, `F8:0F:F9`, `18:7F:88`, or `48:D6:D5` are flagged as Google/Nest devices. This does not distinguish thermostats from cameras or displays -- it just identifies the manufacturer.
2. **No Local Probe** -- Nest thermostats have no open ports on the local network. Port scanning returns nothing. Haus skips local HTTP fingerprinting for devices identified by Google MAC prefixes.
3. **SDM API Enrichment** -- Once the user completes the OAuth2 flow, Haus queries `GET /v1/enterprises/{project_id}/devices` and filters for devices with type `sdm.devices.types.THERMOSTAT`. The response includes the device's custom name (from `sdm.devices.traits.Info.customName`), ambient temperature, humidity, and HVAC mode.
4. **Name Matching** -- Haus matches SDM devices to locally discovered Google MAC devices and enriches their names. A device initially seen as "Google .145" becomes "Living Room Thermostat" once matched.

## Pairing / Authentication

Nest thermostats use Google's partner connection OAuth2 flow. This is NOT standard Google OAuth -- the authorization URL is `nestservices.google.com`, not `accounts.google.com`.

### Prerequisites

1. A **Google Cloud project** with the Smart Device Management API enabled.
2. A **Device Access project** created at `https://console.nest.google.com/device-access` ($5 one-time registration fee per Google account).
3. An **OAuth2 Web Application client** configured in Google Cloud Console with `client_id` and `client_secret`.
4. The user's Google account must be linked to their Nest devices via the Google Home app.

### Environment Variables

```
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-...
GOOGLE_PROJECT_ID=bf59ddff-ff13-4b97-8123-4aca84e315dd
```

### OAuth2 Flow

#### Step 1: Redirect to Authorization

```
GET https://nestservices.google.com/partnerconnections/{project_id}/auth
  ?redirect_uri={redirect_uri}
  &access_type=offline
  &prompt=consent
  &client_id={client_id}
  &response_type=code
  &scope=https://www.googleapis.com/auth/sdm.service
```

**Critical details:**
- The URL MUST be `nestservices.google.com`, NOT `accounts.google.com`. Using the wrong URL will result in a valid OAuth flow that lacks SDM permissions.
- `access_type=offline` is required to receive a refresh token.
- `prompt=consent` ensures a refresh token is always returned (even if the user previously authorized).
- The `{project_id}` in the URL is the Device Access project ID, NOT the Google Cloud project ID.

#### Step 2: Exchange Authorization Code for Tokens

```
POST https://www.googleapis.com/oauth2/v4/token
Content-Type: application/x-www-form-urlencoded

client_id={client_id}&client_secret={client_secret}&code={code}&grant_type=authorization_code&redirect_uri={redirect_uri}
```

Response:
```json
{
  "access_token": "ya29.a0...",
  "expires_in": 3599,
  "refresh_token": "1//0d...",
  "scope": "https://www.googleapis.com/auth/sdm.service",
  "token_type": "Bearer"
}
```

#### Step 3: Refresh Access Token

Access tokens expire after 1 hour. Refresh tokens do NOT expire unless explicitly revoked by the user.

```
POST https://www.googleapis.com/oauth2/v4/token
Content-Type: application/x-www-form-urlencoded

client_id={client_id}&client_secret={client_secret}&refresh_token={token}&grant_type=refresh_token
```

**Important:** Google's refresh response does NOT always include a new refresh token. Haus MUST preserve the original refresh token and only replace it if the response explicitly includes a new one.

### Haus Auth Endpoints

- `GET /api/google/auth` -- Redirects the user to the `nestservices.google.com` authorization URL.
- `GET /api/google/callback` -- Handles the OAuth callback, exchanges the code for tokens, and stores them.
- `GET /api/google/status` -- Returns `{"connected": true}` if valid tokens exist.

## API Reference

All requests go to the Smart Device Management (SDM) API.

**Base URL:** `https://smartdevicemanagement.googleapis.com/v1`

**Auth Header:** `Authorization: Bearer {access_token}`

### List Devices

```
GET /v1/enterprises/{project_id}/devices
```

Returns all devices associated with the user's Google account. Thermostats have type `sdm.devices.types.THERMOSTAT`.

### Get Single Device

```
GET /v1/enterprises/{project_id}/devices/{device_id}
```

### Execute Command

```
POST /v1/enterprises/{project_id}/devices/{device_id}:executeCommand
Content-Type: application/json

{
  "command": "sdm.devices.commands.{CommandName}",
  "params": { ... }
}
```

### Thermostat Commands

| Command | Params | Description |
|---------|--------|-------------|
| `ThermostatMode.SetMode` | `{"mode": "HEAT"}` | Set mode. Values: `HEAT`, `COOL`, `HEATCOOL`, `OFF` |
| `ThermostatTemperatureSetpoint.SetHeat` | `{"heatCelsius": 22.0}` | Set heat setpoint (HEAT mode) |
| `ThermostatTemperatureSetpoint.SetCool` | `{"coolCelsius": 24.0}` | Set cool setpoint (COOL mode) |
| `ThermostatTemperatureSetpoint.SetRange` | `{"heatCelsius": 20.0, "coolCelsius": 24.0}` | Set heat/cool range (HEATCOOL mode) |

**Notes:**
- SetHeat only works when mode is `HEAT` or `HEATCOOL`.
- SetCool only works when mode is `COOL` or `HEATCOOL`.
- SetRange only works when mode is `HEATCOOL`.
- All temperatures are in Celsius. Convert from Fahrenheit before sending.
- Temperature values are rounded to the nearest 0.5 degrees Celsius by the API.

### Thermostat Traits

| Trait | Fields | Description |
|-------|--------|-------------|
| `sdm.devices.traits.Info` | `customName` | User-assigned name in Google Home app |
| `sdm.devices.traits.Temperature` | `ambientTemperatureCelsius` | Current ambient temperature reading |
| `sdm.devices.traits.Humidity` | `ambientHumidityPercent` | Current relative humidity (integer, 0-100) |
| `sdm.devices.traits.ThermostatMode` | `mode`, `availableModes` | Current mode and list of supported modes |
| `sdm.devices.traits.ThermostatTemperatureSetpoint` | `heatCelsius`, `coolCelsius` | Active setpoint(s) depending on mode |
| `sdm.devices.traits.ThermostatHvac` | `status` | Current HVAC activity: `HEATING`, `COOLING`, or `OFF` |
| `sdm.devices.traits.Connectivity` | `status` | `ONLINE` or `OFFLINE` |
| `sdm.devices.traits.Settings` | `temperatureScale` | `FAHRENHEIT` or `CELSIUS` (user preference) |

### Example Device Response

```json
{
  "name": "enterprises/{project_id}/devices/{device_id}",
  "type": "sdm.devices.types.THERMOSTAT",
  "traits": {
    "sdm.devices.traits.Info": {
      "customName": "Living Room"
    },
    "sdm.devices.traits.Temperature": {
      "ambientTemperatureCelsius": 22.5
    },
    "sdm.devices.traits.Humidity": {
      "ambientHumidityPercent": 45
    },
    "sdm.devices.traits.ThermostatMode": {
      "mode": "HEAT",
      "availableModes": ["HEAT", "COOL", "HEATCOOL", "OFF"]
    },
    "sdm.devices.traits.ThermostatTemperatureSetpoint": {
      "heatCelsius": 21.0
    },
    "sdm.devices.traits.ThermostatHvac": {
      "status": "HEATING"
    },
    "sdm.devices.traits.Connectivity": {
      "status": "ONLINE"
    },
    "sdm.devices.traits.Settings": {
      "temperatureScale": "FAHRENHEIT"
    }
  },
  "parentRelations": [
    {
      "parent": "enterprises/{project_id}/structures/{structure_id}/rooms/{room_id}",
      "displayName": "Living Room"
    }
  ]
}
```

### Pub/Sub Events

The SDM API supports Google Cloud Pub/Sub for real-time event notifications. When a thermostat's state changes (mode, temperature, connectivity), an event is pushed to a configured Pub/Sub topic.

**Event structure:**
```json
{
  "eventId": "uuid",
  "timestamp": "2026-04-12T10:00:00Z",
  "resourceUpdate": {
    "name": "enterprises/{project_id}/devices/{device_id}",
    "traits": {
      "sdm.devices.traits.ThermostatHvac": {
        "status": "OFF"
      }
    }
  },
  "userId": "obfuscated-user-id"
}
```

## AI Capabilities

When the AI concierge "chats as" a Nest thermostat, it can:

- **Query temperature and humidity** -- real-time readings via the SDM API ("It's 72.5 F and 45% humidity in here right now.")
- **Report thermostat mode** -- current mode (HEAT, COOL, HEATCOOL, OFF) and active HVAC status (HEATING, COOLING, OFF)
- **Set temperature** -- adjust heat or cool setpoints via SDM commands ("I've set the heat to 70 F for you.")
- **Change mode** -- switch between HEAT, COOL, HEATCOOL, and OFF modes
- **Report setpoints** -- current target temperatures for the active mode

The AI speaks in first person as the thermostat, providing a natural conversational interface to climate control.

## Quirks & Notes

- **Cloud-only:** There is absolutely no local API. If internet goes down, the thermostat continues running its schedule autonomously but cannot be controlled remotely.
- **Temperature rounding:** The SDM API rounds all temperature setpoints to the nearest 0.5 degrees Celsius. Setting 22.3 C will result in 22.5 C.
- **Mode constraints:** You must be in the correct mode before setting a setpoint. Attempting `SetHeat` while in `COOL` mode returns an error. Always `SetMode` first if needed.
- **No learning trait:** The SDM API does not expose the learning/scheduling features. You cannot read or modify the learned schedule via the API -- only set manual temperature overrides.
- **Rate limiting:** The SDM API has a quota of 10 queries per minute per device and 100 queries per minute per project. Haus uses 60-second polling to stay well within limits.
- **Sandbox limitations:** The Device Access sandbox tier is limited to 25 users across 5 structures. Google's Commercial Development program (unlimited users) is currently paused for new applicants.
- **4th gen Matter support:** While the 4th gen supports Matter over Thread, the Matter integration provides basic thermostat control. The SDM API remains the richer integration path with full trait access.
- **Temperature scale:** The API always returns Celsius regardless of the user's display preference. The `Settings.temperatureScale` trait indicates what the physical display shows, not what the API returns.
- **Refresh token preservation:** Google's token refresh endpoint does NOT always return a new refresh token. Haus must store and reuse the original refresh token indefinitely.
- **$5 developer fee:** Each Google account that registers a Device Access project pays a one-time $5 fee. In the commercial Haus product, this fee would be absorbed by Haus's own developer account.

## Similar Devices

- **nest-thermostat-2020** -- Google Nest Thermostat (budget model), same SDM API but fewer hardware features
- **ecobee-smart-thermostat-premium** -- Competing smart thermostat with both cloud API and local HomeKit control
- **honeywell-home-t9** -- Honeywell/Resideo competing thermostat with cloud API and room sensors
