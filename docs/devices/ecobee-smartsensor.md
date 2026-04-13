---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ecobee-smartsensor"
name: "Ecobee SmartSensor"
manufacturer: "Ecobee (Generac Holdings)"
brand: "Ecobee"
model: "EB-RSE3PK2-01"
model_aliases: ["EB-RSE3PK1-01", "EB-RSHM2PK-01", "SmartSensor", "Room Sensor"]
device_type: "ecobee_sensor"
category: "climate"
product_line: "Ecobee"
release_year: 2022
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true        # Data accessible via Ecobee thermostat's cloud API
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []        # BLE only -- no WiFi, no MAC on the network
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "detected_only"
  integration_key: "ecobee"
  polling_interval_sec: 0     # Data comes through the Ecobee thermostat API, not directly
  websocket_event: ""
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities:
  - "temperature"
  - "motion"

# --- PROTOCOL ---
protocol:
  type: ""                    # No direct protocol -- BLE peripheral paired to Ecobee thermostat
  port: 0
  transport: ""
  encoding: ""
  auth_method: "none"
  auth_detail: "No direct API. Pairs with Ecobee thermostat via BLE. Data is accessible through the Ecobee thermostat's cloud API (see ecobee-smart-thermostat-premium.md)."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "sensor"
  power_source: "battery"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.ecobee.com/en-us/accessories/smart-temperature-occupancy-sensor/"
  api_docs: "https://www.ecobee.com/home/developer/api/introduction/index.shtml"
  developer_portal: "https://www.ecobee.com/home/developer/loginDeveloper.jsp"
  support: "https://support.ecobee.com/"
  community_forum: "https://support.ecobee.com/s/community"
  image_url: ""
  fcc_id: "WR2EB-RSE3"

# --- TAGS ---
tags: ["ble_only", "sensor", "temperature", "occupancy", "motion", "ecobee", "battery", "peripheral"]
---

# Ecobee SmartSensor

## What It Is

The Ecobee SmartSensor is a wireless room sensor that extends the temperature sensing and occupancy detection of an Ecobee thermostat to additional rooms. It communicates exclusively via Bluetooth Low Energy (BLE) -- it has no WiFi radio and does not appear on the local network. The sensor pairs directly with a compatible Ecobee thermostat (Smart Thermostat Premium, Smart Thermostat Enhanced, SmartThermostat with Voice Control, or Ecobee3 Lite). Once paired, it reports room temperature and occupancy (motion) data to the thermostat, which in turn makes this data available through the Ecobee cloud API. The thermostat uses sensor data to enable "Smart Home/Away" (auto-detecting occupancy) and "Follow Me" (prioritizing temperature in occupied rooms). The sensor is powered by a single CR2477 coin cell battery that lasts approximately 18-24 months.

## How Haus Discovers It

The Ecobee SmartSensor is NOT directly discoverable on the network because it has no WiFi interface.

1. **No Network Presence** -- The sensor has no IP address, no MAC address on the WiFi/Ethernet network, no open ports. It is invisible to network scans.
2. **BLE Advertisement** -- The sensor advertises via BLE for pairing. In theory, Haus could detect BLE advertisements from unpaired sensors, but this requires a BLE radio on the hub and is not part of standard network discovery.
3. **Cloud API Discovery** -- The sensor appears as a `remoteSensor` in the Ecobee thermostat's API response when queried with `includeSensors: true`. This is the only reliable way to detect it.

### Detection via Ecobee API

When Haus queries the Ecobee thermostat API with `includeSensors: true`, the response includes:

```json
{
  "remoteSensors": [
    {
      "id": "rs:100:1",
      "name": "Bedroom",
      "type": "ecobee3_remote_sensor",
      "inUse": true,
      "capability": [
        {"id": "1", "type": "temperature", "value": "715"},
        {"id": "2", "type": "occupancy", "value": "true"}
      ]
    },
    {
      "id": "ei:0",
      "name": "Living Room",
      "type": "thermostat",
      "inUse": true,
      "capability": [
        {"id": "1", "type": "temperature", "value": "725"},
        {"id": "2", "type": "occupancy", "value": "false"}
      ]
    }
  ]
}
```

**Note:** The thermostat itself also appears as a remote sensor (type `thermostat`) with its built-in temperature and occupancy readings. External SmartSensors have type `ecobee3_remote_sensor`.

## Pairing / Authentication

The SmartSensor has no direct authentication. It pairs with the Ecobee thermostat via BLE.

### Pairing Flow

1. In the Ecobee app or on the thermostat's touchscreen, navigate to **Sensors > Add Sensor**.
2. The thermostat enters BLE discovery mode.
3. Pull the plastic battery tab on the SmartSensor (or reinsert the battery if already activated).
4. The sensor advertises via BLE and the thermostat discovers it.
5. The thermostat assigns the sensor to a room and begins receiving data.

### No Direct API

There is NO way to communicate directly with the SmartSensor from Haus. All sensor data flows:

```
SmartSensor --[BLE]--> Ecobee Thermostat --[WiFi/Cloud]--> Ecobee API --[HTTPS]--> Haus
```

Haus accesses sensor data exclusively through the Ecobee thermostat's cloud API. See **ecobee-smart-thermostat-premium.md** for full API documentation.

## API Reference

The SmartSensor has no direct API. All data is accessed through the parent Ecobee thermostat.

### Reading Sensor Data

Query the Ecobee thermostat API with `includeSensors: true`:

```
GET https://api.ecobee.com/1/thermostat
  ?json={"selection":{"selectionType":"registered","selectionMatch":"","includeSensors":true}}
```

### Sensor Data Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Sensor identifier (format: `rs:XXX:Y` for remote sensors) |
| `name` | string | User-assigned room name |
| `type` | string | `ecobee3_remote_sensor` for SmartSensors, `thermostat` for the built-in sensor |
| `inUse` | boolean | Whether the sensor is actively participating in comfort decisions |
| `capability[type=temperature].value` | string | Temperature in tenths of Fahrenheit (e.g., "715" = 71.5 F) |
| `capability[type=occupancy].value` | string | `"true"` or `"false"` -- whether motion was detected recently |

### Temperature Format

Like all Ecobee API values, temperature is in tenths of a degree Fahrenheit. Divide by 10 for the actual temperature. Convert to Celsius: `(value / 10 - 32) * 5 / 9`.

### Occupancy Behavior

The occupancy sensor reports `true` when motion is detected and remains `true` for approximately 30 minutes after the last motion event. It does NOT provide instantaneous motion events -- it provides a "room is occupied" signal with a decay timer.

## AI Capabilities

The SmartSensor is NOT directly chattable (`ai_chattable: false`). Sensor data is exposed through the parent Ecobee thermostat's AI capabilities.

When chatting with the Ecobee thermostat, the AI can report SmartSensor data:
- "The bedroom sensor reads 71.5 F and detects someone in the room."
- "The office sensor shows 68.2 F with no occupancy."

## Quirks & Notes

- **BLE only:** No WiFi radio whatsoever. The sensor is completely invisible to network scanning tools. It exists only in the Ecobee thermostat's BLE mesh.
- **No direct API:** There is no way to talk to the sensor directly. All data goes thermostat -> cloud -> Haus. This means sensor data is subject to the thermostat's reporting interval and the Ecobee API's rate limits.
- **Temperature in tenths of Fahrenheit:** Same quirky format as the Ecobee thermostat API. "715" means 71.5 F.
- **Occupancy decay timer:** The sensor reports "occupied" for ~30 minutes after the last motion event. This is designed for HVAC purposes (you don't want the AC to shut off every time you sit still for a minute) but is too slow for security/presence detection.
- **Battery life:** CR2477 coin cell, 18-24 months typical. The sensor reports battery status via the API (low battery warning).
- **Range:** BLE range is approximately 45 feet from the thermostat. Thick walls, metal structures, and interference can reduce this.
- **Max sensors:** Ecobee thermostats support up to 32 SmartSensors (though practical limits are lower due to BLE bandwidth).
- **"Follow Me" feature:** When enabled, the thermostat averages temperature readings from sensors in occupied rooms only, effectively "following" occupants around the house. This is controlled via the Ecobee app, not the API.
- **Detected only:** Haus marks this as `detected_only` because the sensor has no direct integration -- its data is a subset of the Ecobee thermostat integration. Haus can display sensor readings but cannot control the sensor itself (there's nothing to control).
- **Previous generation:** Earlier Ecobee sensors (EB-RSE3PK1-01) used a proprietary 915 MHz radio instead of BLE. The current SmartSensor uses BLE and is not backward-compatible with Ecobee3 thermostats.

## Similar Devices

- **ecobee-smart-thermostat-premium** -- The parent thermostat required for this sensor to function
- **honeywell-home-t9** -- Honeywell's T9 also supports room sensors (proprietary 915 MHz radio, not BLE)
