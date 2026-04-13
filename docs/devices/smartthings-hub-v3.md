---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "smartthings-hub-v3"
name: "Samsung SmartThings Hub (v3)"
manufacturer: "Samsung Electronics Co., Ltd."
brand: "SmartThings"
model: "IM6001-V3P"
model_aliases: ["IM6001-V3P01", "SmartThings Hub v3", "SmartThings Hub 2018", "GP-U999SJVLGDA", "STH-ETH-300"]
device_type: "smartthings_hub"
category: "smart_home"
product_line: "SmartThings"
release_year: 2018
discontinued: true
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["zigbee", "zwave", "wifi", "bluetooth", "ethernet"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "D0:52:A8"              # SmartThings (Samsung)
    - "24:FD:5B"              # SmartThings
    - "28:6D:97"              # Samsung Electronics
    - "34:14:B5"              # Samsung Electronics
    - "8C:F5:A3"              # Samsung SmartThings
  mdns_services: []           # SmartThings Hub does not advertise mDNS services
  mdns_txt_keys: []
  default_ports: [39500, 8080]
  signature_ports: [39500]
  ssdp_search_target: "urn:SmartThingsCommunity:device:Hub"
  ssdp_server_string: ""
  hostname_patterns: ["^SmartThings-.*", "^Hub-.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 39500
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "SmartThings"
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "smartthings"
  polling_interval_sec: 30
  websocket_event: "smartthings:state"
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp", "motion", "temperature", "humidity", "lock_unlock", "fan_speed", "thermostat", "battery_level"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "SmartThings API uses OAuth2 with Personal Access Tokens (PATs) or full OAuth2 authorization code flow. PATs can be generated at my.smartthings.com/advanced/pat. The API base URL is https://api.smartthings.com/v1. Bearer token in Authorization header."
  base_url_template: "https://api.smartthings.com/v1"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee", "zwave", "bluetooth_le", "wifi"]

# --- LINKS ---
links:
  product_page: "https://www.smartthings.com/"
  api_docs: "https://developer.smartthings.com/docs/api/public"
  developer_portal: "https://developer.smartthings.com/"
  support: "https://support.smartthings.com/"
  community_forum: "https://community.smartthings.com/"
  image_url: ""
  fcc_id: "2ADZK-IM6001V3P"

# --- TAGS ---
tags: ["zigbee_hub", "zwave_hub", "cloud_api", "oauth2", "samsung", "multi_protocol", "matter_ready", "discontinued_hub"]
---

# Samsung SmartThings Hub (v3)

## What It Is

The Samsung SmartThings Hub (v3) is a multi-protocol smart home hub that bridges Zigbee, Z-Wave, WiFi, and Bluetooth devices into the SmartThings cloud ecosystem. Released in 2018, it was the primary hardware hub for Samsung's SmartThings platform before Samsung shifted toward embedding SmartThings functionality directly into Samsung TVs, refrigerators, and the Aeotec Smart Home Hub (a rebadged successor). The v3 hub connects to your network via Ethernet or WiFi, communicates with local devices over Zigbee 3.0 and Z-Wave Plus, and relays everything through the SmartThings cloud. While the hub itself is discontinued, the SmartThings cloud API remains fully supported and is one of the most well-documented smart home APIs available. The hub received a Matter controller firmware update in late 2023, allowing it to commission and control Matter devices.

## How Haus Discovers It

1. **OUI match** -- SmartThings hubs use MAC prefixes `D0:52:A8`, `24:FD:5B`, `8C:F5:A3`, and Samsung Electronics prefixes `28:6D:97`, `34:14:B5`.
2. **SSDP** -- The hub responds to UPnP/SSDP discovery with search target `urn:SmartThingsCommunity:device:Hub`.
3. **Port probe** -- Port 39500 is the hub's local HTTP callback server (used for SmartApp subscriptions). The `Server` header typically contains "SmartThings".
4. **Hostname pattern** -- DHCP hostname often matches `SmartThings-*` or `Hub-*`.
5. **Cloud enrichment** -- After OAuth2 authentication, the SmartThings API provides full device details including hub firmware version, connected devices, and capabilities.

## Pairing / Authentication

### SmartThings Cloud API (Primary Integration Path)

1. **Personal Access Token (simplest)** -- Generate a PAT at `https://my.smartthings.com/advanced/pat`. Select scopes for devices, locations, and scenes. The PAT is a long-lived bearer token.
2. **OAuth2 Authorization Code Flow (production)** -- For multi-user deployment:
   - Register an app at `https://developer.smartthings.com/`
   - Redirect user to `https://api.smartthings.com/oauth/authorize?client_id={id}&response_type=code&redirect_uri={uri}&scope=r:devices:* x:devices:* r:scenes:* x:scenes:*`
   - User logs in with Samsung account and authorizes
   - Exchange authorization code for access token at `https://api.smartthings.com/oauth/token`
   - Access tokens expire after 24 hours; refresh tokens last 30 days

### Local (Limited)

The SmartThings hub does not expose a general-purpose local API. Port 39500 is a callback server for local SmartApps (Edge drivers), not a general REST API. Samsung has been transitioning to "Edge" drivers that run locally on the hub, but control still routes through the cloud API for third-party integrations.

## API Reference

The SmartThings API is a well-documented cloud REST API.

Base URL: `https://api.smartthings.com/v1`

### Core Endpoints

| Path | Method | Description |
|------|--------|-------------|
| `/devices` | GET | List all devices |
| `/devices/{deviceId}` | GET | Get device details |
| `/devices/{deviceId}/status` | GET | Get full device status (all capabilities) |
| `/devices/{deviceId}/components/{componentId}/capabilities/{capabilityId}/status` | GET | Get specific capability status |
| `/devices/{deviceId}/commands` | POST | Send commands to device |
| `/locations` | GET | List all locations |
| `/locations/{locationId}/rooms` | GET | List rooms in a location |
| `/scenes` | GET | List all scenes |
| `/scenes/{sceneId}/execute` | POST | Execute a scene |
| `/subscriptions` | POST | Subscribe to device events (webhook) |
| `/rules` | GET | List automation rules |

### Device Status Response

```json
{
  "components": {
    "main": {
      "switch": {
        "switch": {
          "value": "on",
          "timestamp": "2024-01-01T12:00:00.000Z"
        }
      },
      "switchLevel": {
        "level": {
          "value": 75,
          "unit": "%",
          "timestamp": "2024-01-01T12:00:00.000Z"
        }
      },
      "colorTemperature": {
        "colorTemperature": {
          "value": 3500,
          "unit": "K"
        }
      }
    }
  }
}
```

### Send Command

```json
POST /devices/{deviceId}/commands
{
  "commands": [
    {
      "component": "main",
      "capability": "switch",
      "command": "on"
    },
    {
      "component": "main",
      "capability": "switchLevel",
      "command": "setLevel",
      "arguments": [75]
    }
  ]
}
```

### Capability Model

SmartThings uses a capability-based abstraction. Common capabilities:

| Capability | Commands | Attributes |
|-----------|----------|------------|
| `switch` | `on`, `off` | `switch` (on/off) |
| `switchLevel` | `setLevel(level)` | `level` (0-100) |
| `colorControl` | `setColor(hue, saturation)` | `hue`, `saturation` |
| `colorTemperature` | `setColorTemperature(temp)` | `colorTemperature` (K) |
| `motionSensor` | (none) | `motion` (active/inactive) |
| `contactSensor` | (none) | `contact` (open/closed) |
| `temperatureMeasurement` | (none) | `temperature` (F or C) |
| `relativeHumidityMeasurement` | (none) | `humidity` (%) |
| `battery` | (none) | `battery` (%) |
| `lock` | `lock`, `unlock` | `lock` (locked/unlocked) |
| `thermostatMode` | `setThermostatMode(mode)` | `thermostatMode` |
| `thermostatCoolingSetpoint` | `setCoolingSetpoint(temp)` | `coolingSetpoint` |
| `thermostatHeatingSetpoint` | `setHeatingSetpoint(temp)` | `heatingSetpoint` |

### Subscriptions (Webhooks)

SmartThings supports webhook subscriptions for real-time event delivery:

```json
POST /subscriptions
{
  "sourceType": "DEVICE",
  "device": {
    "deviceId": "{deviceId}",
    "componentId": "main",
    "capability": "motionSensor",
    "attribute": "motion",
    "stateChangeOnly": true
  },
  "subscription": {
    "subscriptionName": "haus-motion",
    "callbackUrl": "https://your-haus-hub.example.com/webhook/smartthings"
  }
}
```

### Rate Limits

- **Global:** 250 requests per minute per token
- **Per device:** 20 requests per minute per device
- **Subscriptions:** Maximum 200 subscriptions per installed app

## AI Capabilities

When integrated, the AI concierge could leverage the SmartThings API to:

- Control any SmartThings-connected device (lights, locks, thermostats, fans)
- Query sensor states (motion, contact, temperature, humidity)
- Execute SmartThings scenes
- Provide unified device status across Zigbee, Z-Wave, and WiFi devices connected to the hub
- Report battery levels for all battery-powered sensors

## Quirks & Notes

- **Cloud-dependent** -- Despite having local Zigbee/Z-Wave radios, the SmartThings Hub v3 routes all third-party API access through the cloud. If Samsung's servers go down, third-party integrations stop working (though local Edge drivers continue to function).
- **Discontinued but supported** -- Samsung stopped selling the v3 hub but continues to support it with firmware updates. The Aeotec Smart Home Hub is the spiritual successor with identical functionality.
- **Matter controller** -- Firmware update 000.052.00009 (late 2023) added Matter controller support. The hub can commission Thread and WiFi Matter devices.
- **Edge drivers replace Groovy** -- SmartThings migrated from cloud-based Groovy SmartApps to locally-executed Lua-based Edge drivers in 2022-2023. Edge drivers run on the hub itself and provide faster local execution.
- **Port 39500** -- This is the hub's local callback server used by Edge drivers and local SmartApps. It is not a general-purpose API endpoint. Do not attempt to use it for device control.
- **Samsung account required** -- All SmartThings API access requires a Samsung account. There is no way to use the hub without a Samsung account.
- **Webhook requires public URL** -- Subscription webhooks require a publicly accessible HTTPS endpoint. For Haus running locally, this means either polling or setting up a relay/tunnel.
- **Z-Wave inclusion** -- Z-Wave device pairing requires initiating inclusion mode via the SmartThings app or API, then activating the device's pairing button. Z-Wave devices cannot be migrated between hubs without re-pairing.

## Similar Devices

- [aqara-hub-m2](aqara-hub-m2.md) -- Competing Zigbee hub with local HomeKit API
- [ikea-dirigera-hub](ikea-dirigera-hub.md) -- IKEA's Zigbee/Thread hub with local REST API
- [philips-hue-bridge](philips-hue-bridge.md) -- Zigbee bridge focused on lighting with excellent local API
