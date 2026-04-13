---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "nest-thermostat-2020"
name: "Google Nest Thermostat (2020)"
manufacturer: "Google LLC"
brand: "Google Nest"
model: "G4CVZ"
model_aliases: ["GA01334-US", "GA02081-US", "GA02083-US", "T3017US", "A0063"]
device_type: "nest_thermostat"
category: "climate"
product_line: "Nest"
release_year: 2020
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
    - "18:B4:30"        # Google (Nest Labs Inc.)
    - "64:16:66"        # Google (Nest devices)
    - "7C:10:15"        # Google LLC
    - "B0:09:DA"        # Google LLC
    - "F8:0F:F9"        # Google LLC
    - "18:7F:88"        # Google LLC
    - "48:D6:D5"        # Google LLC
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Nest-[A-Z0-9]+"
    - "^Google-Nest"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

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
  auth_detail: "OAuth2 via nestservices.google.com partner connection flow. See nest-learning-thermostat.md for full details."
  base_url_template: "https://smartdevicemanagement.googleapis.com/v1/enterprises/{project_id}"
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
  product_page: "https://store.google.com/product/nest_thermostat"
  api_docs: "https://developers.google.com/nest/device-access/api"
  developer_portal: "https://console.nest.google.com/device-access"
  support: "https://support.google.com/googlenest/answer/9241211"
  community_forum: "https://www.googlenestcommunity.com/"
  image_url: ""
  fcc_id: "A4RA0063"

# --- TAGS ---
tags: ["cloud_only", "sdm_api", "oauth2", "thermostat", "budget", "google", "nest"]
---

# Google Nest Thermostat (2020)

## What It Is

The Google Nest Thermostat (2020, model G4CVZ) is Google's budget-friendly smart thermostat. Released in 2020 at roughly half the price of the Nest Learning Thermostat, it features a mirrored glass display with a sideswipe touch strip for navigation and a clean, minimalist design available in multiple colors (Snow, Charcoal, Sand, Fog). Unlike the Learning Thermostat, it does NOT have the self-learning schedule algorithm, far-field temperature sensors, or a metal housing. It does include Google's "Savings Finder" feature which suggests schedule adjustments to save energy, and it supports basic scheduling through the Google Home app. It connects via 2.4 GHz WiFi and Bluetooth (for initial setup). Like all Nest devices, it is cloud-only -- no local API exists. All programmatic control goes through the Google Smart Device Management (SDM) API, using the same OAuth2 flow and endpoints as the Learning Thermostat.

## How Haus Discovers It

Discovery is identical to the Nest Learning Thermostat. See **nest-learning-thermostat.md** for the full discovery flow.

1. **OUI Match** -- Google/Nest MAC prefixes (`18:B4:30`, `64:16:66`, `7C:10:15`, etc.) identify this as a Google device.
2. **No Local Probe** -- No open ports on the local network. Haus skips HTTP fingerprinting.
3. **SDM API Enrichment** -- After OAuth2 setup, Haus queries the SDM API for `sdm.devices.types.THERMOSTAT` devices and matches them to locally-discovered Google MAC addresses.

**Note:** The SDM API does not distinguish between the 2020 Nest Thermostat and the Learning Thermostat at the device type level. Both report as `sdm.devices.types.THERMOSTAT` with the same trait set. Haus cannot programmatically determine which physical model the user has without additional heuristics.

## Pairing / Authentication

Identical to the Nest Learning Thermostat. See **nest-learning-thermostat.md** for the complete OAuth2 flow, including:

- Prerequisites (Google Cloud project, Device Access project, $5 fee)
- The critical `nestservices.google.com` authorization URL (NOT `accounts.google.com`)
- Token exchange and refresh token preservation
- Haus auth endpoints (`/api/google/auth`, `/api/google/callback`, `/api/google/status`)

## API Reference

The 2020 Nest Thermostat uses the exact same SDM API as the Learning Thermostat. See **nest-learning-thermostat.md** for complete API documentation including:

- List/Get device endpoints
- All thermostat commands (SetMode, SetHeat, SetCool, SetRange)
- All thermostat traits (Temperature, Humidity, ThermostatMode, ThermostatTemperatureSetpoint, ThermostatHvac)
- Pub/Sub event structure

### Trait Differences from Learning Thermostat

The SDM API exposes the same traits for both models. The key hardware differences (no learning, no far-field sensor) are not reflected in the API -- they affect the device's autonomous behavior, not the available API surface.

| Feature | Learning Thermostat (4th Gen) | Nest Thermostat (2020) |
|---------|-------------------------------|------------------------|
| Learning schedule | Yes (on-device) | No (Savings Finder only) |
| Far-field temp sensor | Yes | No |
| Matter/Thread | Yes | No |
| Available modes via API | HEAT, COOL, HEATCOOL, OFF | HEAT, COOL, HEATCOOL, OFF |
| Temperature trait | Same | Same |
| Humidity trait | Same | Same |
| HVAC status trait | Same | Same |
| Setpoint commands | Same | Same |

## AI Capabilities

Identical to the Nest Learning Thermostat. The AI concierge can:

- Query temperature and humidity in real time
- Report and change thermostat mode
- Set heat/cool setpoints
- Report current HVAC activity (HEATING, COOLING, OFF)

See **nest-learning-thermostat.md** for full AI capability details.

## Quirks & Notes

- **No learning:** Despite being a "Nest" thermostat, this model does not learn your schedule. It offers "Savings Finder" -- periodic suggestions in the Google Home app to adjust your schedule for energy savings, which the user must manually accept.
- **No Thread/Matter:** Unlike the 4th gen Learning Thermostat, the 2020 model has WiFi and Bluetooth only. No Thread radio, no Matter support.
- **C-wire requirement:** The 2020 model includes a "Y-wire" jumper plate for systems without a common (C) wire. This is a hardware compatibility feature not exposed via the API.
- **Sideswipe control:** The physical interface is a capacitive touch strip on the right side of the device, not a rotating bezel like the Learning Thermostat. This is irrelevant to API integration.
- **Same API quotas:** 10 queries/minute/device, 100 queries/minute/project. Same sandbox limitations (25 users, 5 structures).
- **Cloud-only:** Same limitation as all Nest devices -- no local control whatsoever.
- **SDM indistinguishable:** The API does not expose the hardware model. Both the 2020 and Learning thermostats appear as `sdm.devices.types.THERMOSTAT` with identical trait sets.

## Similar Devices

- **nest-learning-thermostat** -- Google's premium learning thermostat (4th gen), same API, more hardware features
- **ecobee-smart-thermostat-premium** -- Competing mid-range thermostat with cloud + HomeKit local control
- **honeywell-home-t9** -- Competing budget-friendly smart thermostat with room sensors
