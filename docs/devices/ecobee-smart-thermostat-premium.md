---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ecobee-smart-thermostat-premium"
name: "Ecobee Smart Thermostat Premium"
manufacturer: "Ecobee (Generac Holdings)"
brand: "Ecobee"
model: "EB-STATE5P-01"
model_aliases: ["EB-STATE5P-02", "Smart Thermostat Premium"]
device_type: "ecobee_thermostat"
category: "climate"
product_line: "Ecobee"
release_year: 2022
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true       # HomeKit/local accessory protocol
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false    # HomeKit works locally
  local_only_capable: true    # via HomeKit
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "44:61:32"        # Ecobee Inc.
    - "3C:2E:FF"        # Ecobee Inc.
  mdns_services:
    - "_hap._tcp"       # HomeKit Accessory Protocol
  mdns_txt_keys:
    - "md"              # model name
    - "id"              # device ID
    - "sf"              # status flags (1 = not paired, 0 = paired)
    - "ci"              # category identifier (9 = thermostat)
  default_ports: [80, 443]
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^ecobee[_-]"
    - "^Ecobee"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []       # HomeKit uses encrypted sessions, not plain HTTP

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "ecobee"
  polling_interval_sec: 180
  websocket_event: "ecobee:state"
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
  auth_detail: "OAuth2 with PIN-based or standard authorization code flow via api.ecobee.com. Access tokens expire after 1 hour; refresh tokens expire after ~1 year."
  base_url_template: "https://api.ecobee.com/1"
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
  product_page: "https://www.ecobee.com/en-us/smart-thermostats/smart-thermostat-premium/"
  api_docs: "https://www.ecobee.com/home/developer/api/introduction/index.shtml"
  developer_portal: "https://www.ecobee.com/home/developer/loginDeveloper.jsp"
  support: "https://support.ecobee.com/"
  community_forum: "https://support.ecobee.com/s/community"
  image_url: ""
  fcc_id: "WR2EB-STATE5P"

# --- TAGS ---
tags: ["cloud_api", "homekit", "oauth2", "thermostat", "alexa_built_in", "ecobee", "hybrid", "pin_auth"]
---

# Ecobee Smart Thermostat Premium

## What It Is

The Ecobee Smart Thermostat Premium is Ecobee's flagship thermostat, featuring a large glass touchscreen display, built-in Amazon Alexa speaker and microphone, air quality monitoring (VOC sensor), and a built-in temperature/occupancy sensor with an included wireless SmartSensor for multi-room temperature averaging. Ecobee was acquired by Generac Holdings in 2021. The thermostat connects via 2.4 GHz WiFi and supports both the Ecobee cloud API and Apple HomeKit for local control. The HomeKit integration means it can function on a local network without internet for basic thermostat operations. The cloud API at `api.ecobee.com` provides richer functionality including scheduling, vacation holds, weather integration, and energy reports.

## How Haus Discovers It

1. **OUI Match** -- Devices with MAC prefix `44:61:32` or `3C:2E:FF` are flagged as Ecobee devices.
2. **mDNS Discovery** -- The thermostat advertises `_hap._tcp.local.` (HomeKit Accessory Protocol). TXT records include `ci=9` (category: thermostat) and `md` (model name). The `sf` flag indicates pairing status (1 = unpaired, 0 = paired).
3. **Hostname Pattern** -- Typically appears as `ecobee-XXXX` or similar in DHCP.
4. **Cloud Enrichment** -- After OAuth2 setup, Haus queries the Ecobee API for registered thermostats and enriches locally discovered devices with names and capabilities.

## Pairing / Authentication

Ecobee supports two OAuth2 flows: a standard authorization code flow and a unique PIN-based flow designed for devices without a browser.

### Developer Registration

1. Create a developer account at `https://www.ecobee.com/home/developer/loginDeveloper.jsp`.
2. Register an application to obtain an `api_key` (also called `client_id`). No client secret is used for the PIN flow.

### PIN-Based Authorization Flow

This is Ecobee's signature auth flow, designed for hub devices.

#### Step 1: Request a PIN

```
GET https://api.ecobee.com/authorize
  ?response_type=ecobeePin
  &client_id={api_key}
  &scope=smartWrite
```

**Response:**
```json
{
  "ecobeePin": "abc-defg",
  "code": "authorization_code_here",
  "scope": "smartWrite",
  "expires_in": 9,
  "interval": 30
}
```

- `ecobeePin` -- Display this to the user. They enter it in their Ecobee app or web portal under "My Apps".
- `code` -- Use this to poll for token exchange.
- `expires_in` -- Minutes until the PIN expires (typically 9 minutes).
- `interval` -- Minimum seconds between poll attempts.

#### Step 2: User Authorizes

The user navigates to their Ecobee web portal (ecobee.com > My Account > My Apps > Add Application) and enters the displayed PIN.

#### Step 3: Exchange Code for Tokens

Poll this endpoint after the user confirms:

```
POST https://api.ecobee.com/token
Content-Type: application/x-www-form-urlencoded

grant_type=ecobeePin&code={code}&client_id={api_key}
```

**Response:**
```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 3599,
  "refresh_token": "...",
  "scope": "smartWrite"
}
```

If the user hasn't authorized yet, the response will be an error with `authorization_pending`. Continue polling at the specified interval.

#### Step 4: Refresh Token

```
POST https://api.ecobee.com/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&refresh_token={token}&client_id={api_key}
```

**Important:** Ecobee refresh tokens expire after approximately 1 year. Unlike Google, every refresh response DOES include a new refresh token -- always store the latest one.

### Scope Options

| Scope | Access |
|-------|--------|
| `smartRead` | Read-only access to thermostat data |
| `smartWrite` | Read/write access to thermostat settings and state |

### Standard Authorization Code Flow

Ecobee also supports a standard OAuth2 authorization code flow for web applications:

```
GET https://api.ecobee.com/authorize
  ?response_type=code
  &client_id={api_key}
  &redirect_uri={redirect_uri}
  &scope=smartWrite
```

## API Reference

**Base URL:** `https://api.ecobee.com/1`

**Auth Header:** `Authorization: Bearer {access_token}`

**Content-Type:** `application/json` (most endpoints also accept `text/json`)

### Get Thermostats

```
GET /1/thermostat
  ?json={"selection":{"selectionType":"registered","selectionMatch":"","includeRuntime":true,"includeSettings":true,"includeSensors":true}}
```

The Ecobee API uses a `selection` object to specify which data to include. The selection is passed as a JSON-encoded query parameter.

**Key selection include flags:**
- `includeRuntime` -- current temperature, humidity, HVAC mode
- `includeSettings` -- setpoints, schedule, fan settings
- `includeSensors` -- remote sensor data (temperature, occupancy)
- `includeWeather` -- local weather forecast
- `includeEquipmentStatus` -- what's currently running (compressor1, fan, auxHeat1, etc.)
- `includeProgram` -- full schedule and comfort profiles

**Response (abbreviated):**
```json
{
  "thermostatList": [
    {
      "identifier": "123456789012",
      "name": "Living Room",
      "thermostatRev": "2026041210",
      "isRegistered": true,
      "modelNumber": "nikeSmart",
      "runtime": {
        "actualTemperature": 725,
        "actualHumidity": 45,
        "desiredHeat": 700,
        "desiredCool": 760,
        "desiredFanMode": "auto"
      },
      "settings": {
        "hvacMode": "heat",
        "heatStages": 1,
        "coolStages": 1,
        "fanMinOnTime": 0
      },
      "remoteSensors": [
        {
          "id": "rs:100",
          "name": "Living Room",
          "type": "thermostat",
          "inUse": true,
          "capability": [
            {"id": "1", "type": "temperature", "value": "725"},
            {"id": "2", "type": "occupancy", "value": "true"}
          ]
        }
      ]
    }
  ]
}
```

**Temperature format:** Ecobee uses integer values in tenths of a degree Fahrenheit. `725` = 72.5 F. Divide by 10 for display. Convert to Celsius: `(value / 10 - 32) * 5 / 9`.

### Set Thermostat State

```
POST /1/thermostat
Content-Type: application/json

{
  "selection": {
    "selectionType": "thermostats",
    "selectionMatch": "123456789012"
  },
  "functions": [
    {
      "type": "setHold",
      "params": {
        "holdType": "nextTransition",
        "heatHoldTemp": 700,
        "coolHoldTemp": 760
      }
    }
  ]
}
```

**Hold types:**
- `nextTransition` -- Hold until the next scheduled comfort setting
- `indefinite` -- Hold indefinitely until manually cancelled
- `holdHours` -- Hold for a specified number of hours (requires `holdHours` param)
- `dateTime` -- Hold until a specific date/time

### Set HVAC Mode

```
POST /1/thermostat
Content-Type: application/json

{
  "selection": {
    "selectionType": "thermostats",
    "selectionMatch": "123456789012"
  },
  "thermostat": {
    "settings": {
      "hvacMode": "heat"
    }
  }
}
```

**HVAC modes:** `heat`, `cool`, `auto`, `auxHeatOnly`, `off`

### Resume Program

Cancel any active hold and resume the normal schedule:

```
POST /1/thermostat
Content-Type: application/json

{
  "selection": {
    "selectionType": "thermostats",
    "selectionMatch": "123456789012"
  },
  "functions": [
    {
      "type": "resumeProgram",
      "params": {
        "resumeAll": true
      }
    }
  ]
}
```

## AI Capabilities

When the AI concierge "chats as" an Ecobee thermostat, it can:

- **Query temperature and humidity** -- current readings from the thermostat and all connected SmartSensors
- **Report HVAC mode and status** -- what mode it's in and what equipment is actively running
- **Set temperature holds** -- adjust heat/cool setpoints with configurable hold duration
- **Change HVAC mode** -- switch between heat, cool, auto, and off
- **Report sensor data** -- temperature and occupancy from remote SmartSensors
- **Resume schedule** -- cancel holds and return to the programmed schedule
- **Report air quality** -- VOC readings from the built-in sensor (Premium model only)

## Quirks & Notes

- **Temperature in tenths of Fahrenheit:** The API uses integer values representing tenths of a degree Fahrenheit (e.g., 725 = 72.5 F). This is unique among thermostat APIs and requires conversion for display and for Celsius users.
- **Selection object:** Almost every API call requires a `selection` object specifying which thermostats and which data to include. This is verbose but allows precise control over response payloads.
- **Hold vs program:** Ecobee distinguishes between the programmed schedule ("comfort settings" like Home, Away, Sleep) and temporary holds. Setting a temperature via the API creates a hold that overrides the schedule. The `resumeProgram` function cancels holds.
- **Built-in Alexa:** The Premium model includes a full Alexa speaker. This is independent of the API integration but means the device can respond to voice commands directly.
- **HomeKit local control:** The thermostat supports Apple HomeKit, which operates over the local network without internet. Haus could potentially use the HomeKit Accessory Protocol (HAP) for local control, though this requires HomeKit pairing.
- **Rate limiting:** The Ecobee API allows approximately 3 requests per minute per thermostat. Polling should be infrequent (3-minute intervals recommended).
- **Model numbers in API:** Ecobee uses internal model names in the API. The Premium is `nikeSmart` (or similar codename). Don't rely on model numbers for identification -- use the capabilities and feature flags instead.
- **SmartSensor data:** Remote sensor data is included in the thermostat response when `includeSensors` is true. Sensors report temperature and occupancy but not humidity.
- **Generac acquisition:** Ecobee was acquired by Generac Holdings in 2021. The API and developer portal remain under the ecobee.com domain.

## Similar Devices

- **ecobee-smartsensor** -- Wireless room sensor that pairs with Ecobee thermostats via BLE
- **nest-learning-thermostat** -- Google's competing premium thermostat (cloud-only, SDM API)
- **honeywell-home-t9** -- Honeywell/Resideo competing thermostat with room sensors
