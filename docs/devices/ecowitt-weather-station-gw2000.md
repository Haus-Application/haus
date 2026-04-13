---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ecowitt-weather-station-gw2000"
name: "Ecowitt GW2000 WiFi Weather Station Gateway"
manufacturer: "Ecowitt (Fine Offset Electronics Co., Ltd.)"
brand: "Ecowitt"
model: "GW2000"
model_aliases: ["GW2000A", "GW2000B", "GW2000C", "Ecowitt WiFi Gateway"]
device_type: "ecowitt_weather_gateway"
category: "climate"
product_line: "Ecowitt"
release_year: 2022
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "48:3F:DA"              # Espressif (ESP32 inside GW2000)
    - "DC:4F:22"              # Espressif Systems
    - "EC:94:CB"              # Espressif Systems
    - "30:C6:F7"              # Espressif Systems
  mdns_services: []           # GW2000 does not advertise mDNS services
  mdns_txt_keys: []
  default_ports: [80, 49123]
  signature_ports: [49123]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^GW2000.*", "^ecowitt.*", "^espressif.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: "GW2000"
    server_header: ""
    body_contains: "Ecowitt"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "ecowitt"
  polling_interval_sec: 60
  websocket_event: "ecowitt:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["temperature", "humidity"]

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 80
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "none"
  auth_detail: "No authentication required for the local HTTP API. The GW2000's built-in web server on port 80 provides a configuration interface and the /get_livedata_info endpoint returns all sensor readings without any credentials."
  base_url_template: "http://{ip}"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "gateway"
  power_source: "usb"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.ecowitt.com/shop/goodsDetail/107"
  api_docs: "https://osswww.ecowitt.net/uploads/20220407/WN1900%20GW1000,bindest1100%20WiFi%20Gateway%20API.pdf"
  developer_portal: ""
  support: "https://www.ecowitt.com/support"
  community_forum: "https://www.wxforum.net/index.php?board=63.0"
  image_url: ""
  fcc_id: "2AC7Z-GW2000"

# --- TAGS ---
tags: ["weather_station", "wifi", "local_api", "no_auth", "esp32", "temperature", "humidity", "wind", "rain", "uv", "solar_radiation", "barometric_pressure", "hackable", "fine_offset"]
---

# Ecowitt GW2000 WiFi Weather Station Gateway

## What It Is

The Ecowitt GW2000 is a compact WiFi gateway that receives data from Ecowitt's extensive line of outdoor weather sensors (sold separately) and makes that data available over your local network and optionally to cloud services. It is essentially a small ESP32-based box with a 915 MHz (or 868 MHz in EU) RF receiver that listens for wireless transmissions from outdoor sensor arrays -- temperature/humidity probes, wind speed/direction sensors, rain gauges, UV/solar radiation sensors, soil moisture probes, and barometric pressure sensors. The GW2000 has a built-in temperature, humidity, and barometric pressure sensor of its own and serves as the central aggregation point for all sensor data. The device has an excellent local HTTP API that returns all sensor readings as JSON with no authentication required, making it one of the most integration-friendly weather stations available. It is powered by USB-C and connects over 2.4 GHz WiFi. The Ecowitt platform is a rebranded version of Fine Offset Electronics, a Chinese manufacturer whose hardware is also sold under the Ambient Weather, Froggit, Sainlogic, and other brand names.

## How Haus Discovers It

1. **OUI match** -- The GW2000 uses an Espressif ESP32 SoC, so its WiFi MAC address matches Espressif OUI prefixes: `48:3F:DA`, `DC:4F:22`, `EC:94:CB`, `30:C6:F7`. This is a broad match (many IoT devices use ESP32), so additional fingerprinting is needed.
2. **Hostname pattern** -- DHCP hostname is typically `GW2000A` or similar, matching `^GW2000.*`.
3. **Port probe** -- Port 80 (HTTP) hosts the configuration web interface. Port 49123 is used for the device's custom UDP protocol (used by the WS View Plus app for configuration).
4. **HTTP fingerprint** -- A GET request to `http://{ip}/` returns an HTML page with "GW2000" in the title and "Ecowitt" in the body.
5. **API probe** -- A GET request to `http://{ip}/get_livedata_info` returns JSON with all current sensor readings, confirming the device type definitively.

## Pairing / Authentication

No authentication is required. The GW2000's local API is completely open by design.

1. **Initial WiFi setup** -- On first power-up, the GW2000 creates a WiFi AP (typically named `GW2000A-WIFI????`). Connect to it and configure your home WiFi SSID/password via the web interface at `http://192.168.4.1`.
2. **Network join** -- The GW2000 joins your WiFi network and receives a DHCP address.
3. **Immediate access** -- The local HTTP API is available immediately with no pairing, no API keys, no accounts, no tokens. Just HTTP GET requests.
4. **Optional cloud** -- The GW2000 can push data to Ecowitt.net, Weather Underground, WeatherCloud, WOW, and custom HTTP endpoints. This is configured via the web interface but is entirely optional.

## API Reference

### Local HTTP API

The GW2000 exposes a local HTTP API on port 80. All endpoints return JSON. No authentication required.

#### GET /get_livedata_info

Returns all current sensor readings from all connected sensors.

Example response (abbreviated):

```json
{
  "common_list": [
    {
      "id": "0x02",
      "val": "72.5",
      "unit": "℉"
    },
    {
      "id": "0x07",
      "val": "45",
      "unit": "%"
    },
    {
      "id": "0x0A",
      "val": "29.83",
      "unit": "inHg"
    },
    {
      "id": "0x0B",
      "val": "29.79",
      "unit": "inHg"
    }
  ],
  "wh25": [
    {
      "id": "0x01",
      "val": "68.9",
      "unit": "℉"
    },
    {
      "id": "0x02",
      "val": "52",
      "unit": "%"
    }
  ],
  "rain": [
    {
      "id": "0x0D",
      "val": "0.000",
      "unit": "in/Hr"
    },
    {
      "id": "0x0E",
      "val": "0.000",
      "unit": "in"
    },
    {
      "id": "0x10",
      "val": "0.453",
      "unit": "in"
    },
    {
      "id": "0x11",
      "val": "12.567",
      "unit": "in"
    }
  ],
  "wh65": [
    {
      "id": "0x0A",
      "val": "0",
      "unit": ""
    },
    {
      "id": "0x15",
      "val": "0",
      "unit": "W/m²"
    },
    {
      "id": "0x17",
      "val": "0.0",
      "unit": ""
    }
  ]
}
```

#### Sensor ID Reference

| ID | Description | Unit |
|----|-------------|------|
| `0x01` | Indoor temperature | F/C |
| `0x02` | Outdoor temperature | F/C |
| `0x03` | Dew point | F/C |
| `0x04` | Wind chill | F/C |
| `0x05` | Heat index | F/C |
| `0x06` | Indoor humidity | % |
| `0x07` | Outdoor humidity | % |
| `0x08` | Absolute barometric pressure | inHg/hPa |
| `0x09` | Relative barometric pressure | inHg/hPa |
| `0x0A` | Wind direction | degrees |
| `0x0B` | Wind speed | mph/m/s |
| `0x0C` | Wind gust | mph/m/s |
| `0x0D` | Rain rate | in/Hr |
| `0x0E` | Rain event | in |
| `0x0F` | Rain hourly | in |
| `0x10` | Rain daily | in |
| `0x11` | Rain weekly | in |
| `0x12` | Rain monthly | in |
| `0x13` | Rain yearly | in |
| `0x14` | Rain total | in |
| `0x15` | Solar radiation | W/m2 |
| `0x16` | UV index | index |
| `0x17` | UV index | index |

#### GET /get_device_info

Returns device information:

```json
{
  "model": "GW2000A",
  "frequency": "915M",
  "date": "2024-01-15",
  "firmware": "V2.2.8"
}
```

#### GET /get_sensors_info

Returns information about all registered wireless sensors:

```json
{
  "wh65": {
    "id": 1,
    "model": "WH65B",
    "battery": 0,
    "signal": 4
  }
}
```

Battery values: `0` = OK, `1` = low. Signal: `0-4` scale.

#### GET /get_units_info

Returns the configured display units (temperature F/C, pressure inHg/hPa, wind mph/m/s, rain in/mm).

### Custom Server (Push Mode)

The GW2000 can be configured to POST sensor data to a custom HTTP endpoint at configurable intervals (minimum 16 seconds). This is configured in the web interface under "Weather Services" > "Customized". The POST body uses either Ecowitt protocol format or Wunderground-compatible format.

Ecowitt protocol POST body example:

```
PASSKEY=xxxx&stationtype=GW2000A_V2.2.8&dateutc=2024-01-15+12:00:00&tempinf=68.9&humidityin=52&baromrelin=29.83&baromabsin=29.79&tempf=72.5&humidity=45&winddir=180&windspeedmph=5.4&windgustmph=8.9&maxdailygust=12.3&rainratein=0.000&eventrainin=0.000&dailyrainin=0.000&weeklyrainin=0.453&monthlyrainin=2.567&yearlyrainin=12.567&solarradiation=156.78&uv=3&dateutc=2024-01-15+12:00:00
```

This push mode is ideal for Haus integration -- configure the GW2000 to push data to Haus's HTTP endpoint, eliminating the need for polling entirely.

## AI Capabilities

The Ecowitt GW2000 is an excellent candidate for AI chat integration. The AI concierge could:

- Provide natural-language weather reports ("It's 72 degrees and partly cloudy with 45% humidity")
- Report wind conditions with context ("Wind is 5 mph from the south, gusting to 9 mph -- good conditions for grilling")
- Rain summaries ("No rain today, but 0.45 inches this week and 12.5 inches this year")
- UV index warnings ("UV index is 8 -- high. Consider sunscreen if going outside")
- Solar radiation data for solar panel performance correlation
- Barometric pressure trends for weather prediction ("Pressure has been dropping steadily -- storm likely approaching")
- Historical comparisons ("Today is 5 degrees warmer than yesterday at this time")
- Indoor vs. outdoor comfort comparisons

## Quirks & Notes

- **ESP32-based** -- The GW2000 runs on an Espressif ESP32, which means its MAC address matches generic Espressif OUIs. Many other IoT devices also use ESP32, so OUI alone is not sufficient for identification.
- **No authentication** -- The local API has zero security. Anyone on the local network can read all sensor data. This is fine for home use but worth noting for security-conscious deployments.
- **Sensor ID scheme** -- The sensor data uses hex IDs (0x01, 0x02, etc.) that must be mapped to human-readable names. The mapping is documented in the Ecowitt API PDF.
- **Unit configuration** -- Sensor values are returned in whatever units the GW2000 is configured to display (F vs. C, inHg vs. hPa, etc.). Haus should query `/get_units_info` first to determine the active unit system, or normalize to metric internally.
- **Custom server push** -- The most reliable integration method. Configure the GW2000 to push data to a Haus endpoint every 16-60 seconds. This avoids polling and provides near-real-time updates.
- **Fine Offset compatibility** -- The GW2000 is compatible with the broader Fine Offset ecosystem. Sensors sold under Ambient Weather, Froggit, Sainlogic, and other brands often use the same 915 MHz protocol and work with the GW2000.
- **Firmware updates** -- OTA firmware updates available through the WS View Plus mobile app. Ecowitt releases updates periodically that add sensor support and fix bugs.
- **Port 49123 UDP** -- Used by the WS View Plus app for device discovery and configuration. Sends/receives binary protocol messages. Not useful for sensor data retrieval.
- **Multiple channel support** -- The GW2000 supports up to 8 temperature/humidity channels (WH31 sensors), 8 soil moisture channels, 4 PM2.5 channels, and more. Each channel appears with a numbered suffix in the API response.
- **868 MHz vs. 915 MHz** -- The GW2000 comes in regional variants. The US version uses 915 MHz, the EU version uses 868 MHz. Sensors must match the gateway's frequency.

## Similar Devices

- [airthings-wave-plus](airthings-wave-plus.md) -- Indoor air quality sensor (complementary: outdoor weather vs. indoor air)
- [nest-learning-thermostat](nest-learning-thermostat.md) -- Thermostat with temperature/humidity (indoor only)
- [ecobee-smart-thermostat-premium](ecobee-smart-thermostat-premium.md) -- Thermostat with remote temperature sensors
