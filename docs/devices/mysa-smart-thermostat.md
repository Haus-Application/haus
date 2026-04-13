---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "mysa-smart-thermostat"
name: "Mysa Smart Thermostat for Electric Baseboard Heaters"
manufacturer: "Empowered Homes Inc."
brand: "Mysa"
model: "MY100WMN"
model_aliases: ["MY100WMN-01", "MY200WMN", "Mysa V2"]
device_type: "mysa_thermostat"
category: "climate"
product_line: "Mysa"
release_year: 2018
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
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "CC:50:E3"        # Espressif Systems (common ESP32 OUI used by Mysa)
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
    - "^mysa"
    - "^Mysa"
    - "^espressif"      # Generic ESP32 hostname (may appear before cloud provisioning)
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "mysa"
  polling_interval_sec: 120
  websocket_event: "mysa:state"
  setup_type: "oauth2"
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
  auth_method: "oauth2"
  auth_detail: "Cloud API via Mysa backend. OAuth2 or proprietary token auth. API is not publicly documented; reverse engineering or partnership required."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "thermostat"
  power_source: "hardwired"
  mounting: "in_wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.mysa.com/products/mysa-baseboard"
  api_docs: ""
  developer_portal: ""
  support: "https://www.mysa.com/pages/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2ATBK-MY100WMN"

# --- TAGS ---
tags: ["cloud_only", "thermostat", "baseboard", "high_voltage", "line_voltage", "240v", "120v", "electric_heat", "mysa"]
---

# Mysa Smart Thermostat for Electric Baseboard Heaters

## What It Is

The Mysa Smart Thermostat is a WiFi-connected line-voltage thermostat designed specifically for electric baseboard heaters, convectors, and in-floor radiant heating systems. Manufactured by Empowered Homes Inc. (a Canadian company based in St. John's, Newfoundland), Mysa fills a niche that most smart thermostats ignore: high-voltage (120V/240V) heating systems that are common in older homes, apartments, and throughout Canada and the northeastern United States. Unlike the Nest or Ecobee which control low-voltage (24V) HVAC systems, the Mysa sits directly in the high-voltage circuit between the breaker panel and the baseboard heater, switching up to 3800W of load at 240V. It features a minimalist white design with a touchscreen interface, connects via 2.4 GHz WiFi, and is controlled through the Mysa app and cloud service. It supports scheduling, geofencing, and energy usage monitoring.

## How Haus Discovers It

1. **OUI Match** -- Mysa thermostats use Espressif (ESP32) WiFi modules. MAC prefixes like `CC:50:E3`, `AC:67:B2`, `24:62:AB`, or `A4:CF:12` identify Espressif-based devices, but this is a very broad match -- many IoT devices use ESP32 chips.
2. **Hostname Pattern** -- May appear as `mysa-XXXX` or `espressif` in DHCP, depending on firmware version and provisioning state.
3. **No Local Probe** -- Mysa has no open ports on the local network. All communication goes through AWS IoT Core (MQTT over TLS).
4. **Cloud Enrichment** -- If/when Haus integrates with Mysa's API, device names and state would be pulled from the cloud.

**Note:** Reliably identifying a Mysa on the network is challenging because it uses generic Espressif OUIs shared by thousands of other IoT products. Haus would need cloud API integration or user confirmation to positively identify the device.

## Pairing / Authentication

Mysa does not have a publicly documented developer API. Integration options:

### Current State

- **No public API documentation.** Mysa has not released a developer portal or public REST API.
- **Third-party integrations** exist via Home Assistant (community-maintained, using reverse-engineered API calls to Mysa's AWS backend) and Apple HomeKit (native support in newer firmware).
- **MQTT backend:** Mysa devices communicate with an AWS IoT Core MQTT broker. The device authenticates with AWS IoT certificates provisioned during factory setup.

### Potential Integration Paths

1. **HomeKit (local):** Newer Mysa firmware supports HomeKit Accessory Protocol (HAP), which would allow local control without cloud dependency. Haus could use HAP for local thermostat control.
2. **Reverse-engineered cloud API:** Community projects (Home Assistant mysa integration) have documented API endpoints. This is fragile and could break with firmware updates.
3. **Partnership:** A formal partnership with Mysa/Empowered Homes would provide stable API access.

### HomeKit Discovery

When HomeKit is enabled, the thermostat advertises `_hap._tcp.local.` with `ci=9` (thermostat category). HAP pairing requires the setup code printed on the device or displayed during setup.

## API Reference

No public API is documented. The following is based on community reverse-engineering efforts.

### Known Cloud Architecture

- **Backend:** AWS IoT Core (MQTT over TLS on port 8883)
- **Mobile app API:** HTTPS REST calls to a Mysa-hosted API (likely `api.mysa.com` or similar AWS API Gateway endpoint)
- **Authentication:** Cognito-based user authentication in the mobile app; device-level authentication via AWS IoT X.509 certificates

### Community-Discovered Endpoints

These endpoints are based on Home Assistant community integrations and may change without notice:

- Device list / status query
- Set target temperature
- Set mode (heat, off)
- Set schedule

**Note:** Haus should NOT ship a production integration based on reverse-engineered endpoints. This information is for planning purposes only.

## AI Capabilities

When the AI concierge "chats as" a Mysa thermostat, it can:

- **Query temperature** -- current ambient temperature reading
- **Report mode** -- heating or off
- **Set target temperature** -- adjust the heat setpoint
- **Report energy usage** -- power consumption data (if available via API)
- **Explain high-voltage context** -- inform the user that this thermostat directly controls baseboard heaters

## Quirks & Notes

- **Line voltage, not low voltage:** This is a critical distinction. Mysa switches 120V or 240V directly, unlike Nest/Ecobee/Honeywell which control 24V HVAC wiring. Installation involves high-voltage wiring and should be done by an electrician.
- **Heating only:** Mysa is designed for electric resistance heating (baseboard heaters, convectors, in-floor radiant). It does NOT support cooling, heat pumps, or forced-air systems.
- **No humidity sensor:** Unlike most smart thermostats, Mysa does not include a humidity sensor. The capabilities list is `[thermostat, temperature]` only.
- **ESP32-based:** The hardware uses an Espressif ESP32 WiFi module, which means generic Espressif OUIs make network identification unreliable.
- **Canadian market focus:** Mysa was designed for the Canadian market where baseboard heating is extremely common. It works equally well in the US with 120V or 240V systems.
- **Multiple product lines:** Mysa makes separate thermostats for baseboard heaters (MY100WMN), in-floor heating (MY200WMN), and air conditioners/mini-splits (MY300WMN). Each has different capabilities.
- **Energy monitoring:** Mysa tracks energy consumption in the app, showing daily/weekly/monthly kWh usage. This data would be valuable for Haus energy dashboards if API access is available.
- **No public API:** The lack of a documented API makes this a challenging integration. HomeKit support provides the best local control path.
- **AWS IoT Core:** The MQTT backend means the device-to-cloud communication is well-secured with X.509 certificates, but it also means there's no simple REST API to tap into.

## Similar Devices

- **ecobee-smart-thermostat-premium** -- Low-voltage smart thermostat (different use case, but same category)
- **nest-learning-thermostat** -- Low-voltage smart thermostat (different use case)
- **cielo-breez-plus** -- Controls AC/mini-splits via IR (Mysa also has an AC model)
- **sensibo-air** -- Controls AC via IR (comparable to Mysa's AC model)
