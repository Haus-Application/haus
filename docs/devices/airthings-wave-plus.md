---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "airthings-wave-plus"
name: "Airthings Wave Plus"
manufacturer: "Airthings ASA"
brand: "Airthings"
model: "2930"
model_aliases: ["Wave Plus", "Airthings 2930", "Airthings Wave+"]
device_type: "airthings_air_quality"
category: "climate"
product_line: "Airthings Wave"
release_year: 2018
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "58:93:D8"              # Airthings ASA
  mdns_services: []           # BLE device, no mDNS
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
  status: "planned"
  integration_key: "airthings"
  polling_interval_sec: 300   # BLE readings every 5 minutes; cloud API every 5 minutes
  websocket_event: "airthings:state"
  setup_type: "oauth2"
  ai_chattable: true
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["temperature", "humidity"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Airthings API uses OAuth2 Client Credentials flow. Register at dashboard.airthings.com/integrations/api-integration to obtain client_id and client_secret. Token endpoint: https://accounts-api.airthings.com/v1/token. Bearer token in Authorization header. Alternatively, BLE GATT direct reads possible with no authentication."
  base_url_template: "https://ext-api.airthings.com/v1"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "sensor"
  power_source: "battery"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.airthings.com/wave-plus"
  api_docs: "https://developer.airthings.com/docs"
  developer_portal: "https://developer.airthings.com/"
  support: "https://www.airthings.com/support"
  community_forum: "https://forum.airthings.com/"
  image_url: ""
  fcc_id: "2AMJG-2930"

# --- TAGS ---
tags: ["air_quality", "radon", "co2", "voc", "temperature", "humidity", "barometric_pressure", "ble", "cloud_api", "oauth2", "battery"]
---

# Airthings Wave Plus

## What It Is

The Airthings Wave Plus is a battery-powered indoor air quality monitor that measures six environmental parameters: radon (the leading cause of lung cancer among non-smokers), carbon dioxide (CO2), total volatile organic compounds (TVOCs), temperature, humidity, and barometric pressure. It is a circular wall-mounted unit roughly the size of a smoke detector, powered by two AA batteries with approximately 16 months of battery life. The sensor communicates via Bluetooth Low Energy (BLE) -- you can wave your hand in front of it to see a color-coded air quality indicator (green/yellow/red), view readings in the Airthings app on your phone, or access historical data through the Airthings cloud API. The Airthings Hub (sold separately) can bridge BLE readings to the cloud via WiFi, enabling continuous remote monitoring without a phone in BLE range. Airthings is a Norwegian company that is the world leader in consumer radon detection technology.

## How Haus Discovers It

The Airthings Wave Plus is a BLE device with no WiFi or Ethernet interface. Discovery options:

### Via BLE (Direct)

1. **BLE scan** -- The Wave Plus advertises as a BLE peripheral. During BLE scanning, it appears with the local name `Airthings Wave+` or a serial number-based name.
2. **Manufacturer data** -- The BLE advertisement includes Airthings' Bluetooth SIG company ID (`0x0334`, decimal 820) in the manufacturer-specific data field.
3. **MAC prefix** -- Airthings devices use the OUI prefix `58:93:D8`.
4. **Service UUIDs** -- The device advertises custom BLE service UUID `b42e1c08-ade7-11e4-89d3-123b93f75cba` (Airthings sensor service).

### Via Airthings Hub (Cloud)

1. **Hub discovery** -- If an Airthings Hub (SmartLink or View Plus hub) is on the network, Haus discovers it via OUI matching or mDNS.
2. **Cloud API query** -- After OAuth2 authentication, the Airthings cloud API returns all devices on the account, including Wave Plus units.

### Via Airthings Cloud API

If the user connects their Airthings account, Haus queries the cloud API to enumerate all devices without needing BLE proximity.

## Pairing / Authentication

### BLE Direct (No Auth Required)

1. **BLE scan** -- Discover the Wave Plus via BLE advertisements.
2. **GATT connect** -- Connect to the BLE GATT server. No pairing or PIN required for reading sensor data (the device allows unauthenticated reads).
3. **Read characteristics** -- Read the sensor data characteristic to get current readings.

### Airthings Cloud API (OAuth2)

1. **Register API client** -- Go to `https://dashboard.airthings.com/integrations/api-integration` and create an API client. This provides a `client_id` and `client_secret`.
2. **Obtain access token** -- POST to `https://accounts-api.airthings.com/v1/token`:
   ```
   grant_type=client_credentials&client_id={client_id}&client_secret={client_secret}&scope=read:device
   ```
3. **Receive bearer token** -- The response contains an `access_token` (expires in 3600 seconds) and token type.
4. **API calls** -- Include `Authorization: Bearer {access_token}` header on all API requests.
5. **Token refresh** -- Request a new token before expiry using the same client credentials flow.

## API Reference

### BLE GATT Interface

The Wave Plus exposes sensor data via a custom BLE GATT service:

**Service UUID:** `b42e1c08-ade7-11e4-89d3-123b93f75cba`

**Sensor Data Characteristic UUID:** `b42e2a68-ade7-11e4-89d3-123b93f75cba`

Reading this characteristic returns a 20-byte binary payload:

| Byte Offset | Length | Description | Conversion |
|-------------|--------|-------------|------------|
| 0 | 1 | Version | Raw value |
| 1 | 1 | Humidity | Value / 2.0 = % RH |
| 2-3 | 2 | Ambiguous light / unused | -- |
| 4-5 | 2 | Radon short-term average | Raw = Bq/m3 |
| 6-7 | 2 | Radon long-term average | Raw = Bq/m3 |
| 8-9 | 2 | Temperature | Value / 100.0 = Celsius |
| 10-11 | 2 | Atmospheric pressure | Value / 50.0 = hPa |
| 12-13 | 2 | CO2 | Raw = ppm |
| 14-15 | 2 | TVOC | Raw = ppb |

All multi-byte values are little-endian unsigned integers.

**Command Characteristic UUID:** `b42e2d06-ade7-11e4-89d3-123b93f75cba`

Used to request current sensor readings (write `0x01` to trigger a measurement update).

### Airthings Cloud REST API

Base URL: `https://ext-api.airthings.com/v1`

| Path | Method | Description |
|------|--------|-------------|
| `/devices` | GET | List all devices on account |
| `/devices/{serialNumber}` | GET | Get device info |
| `/devices/{serialNumber}/latest-samples` | GET | Get latest sensor readings |
| `/devices/{serialNumber}/samples` | GET | Get historical samples (with time range) |
| `/locations` | GET | List locations |
| `/locations/{locationId}` | GET | Get location details |

#### Latest Samples Response

```json
{
  "data": {
    "battery": 100,
    "co2": 658.0,
    "humidity": 52.0,
    "pm1": null,
    "pm25": null,
    "pressure": 1013.25,
    "radonShortTermAvg": 45.0,
    "temp": 22.3,
    "time": 1705305600,
    "voc": 125.0
  }
}
```

#### Historical Samples

```
GET /devices/{serialNumber}/samples?start=2024-01-01T00:00:00Z&end=2024-01-15T00:00:00Z&resolution=HOUR
```

Resolution options: `HOUR`, `FOUR_HOURS`, `DAY`, `WEEK`.

### Rate Limits

- **Cloud API:** 120 requests per hour per API client
- **BLE:** Sensor updates every 5 minutes internally. Reading more frequently returns the same values.

## AI Capabilities

The Airthings Wave Plus is an excellent candidate for AI chat. The AI concierge could:

- Provide natural-language air quality summaries ("Indoor air quality is good. CO2 is 658 ppm, VOCs are low at 125 ppb, and humidity is comfortable at 52%")
- Radon monitoring alerts ("Radon levels have been elevated at 150 Bq/m3 over the past week -- consider increasing ventilation")
- Context-aware advice ("CO2 is at 1200 ppm in the bedroom -- opening a window or running the HVAC fan would help")
- Correlate indoor environment with outdoor weather data (pair with Ecowitt)
- Track trends over time ("Humidity has been increasing all week -- now at 65%. Running the dehumidifier would prevent mold risk")
- Health-oriented recommendations based on WHO and EPA guidelines for radon, CO2, and VOC levels

## Quirks & Notes

- **BLE range limitation** -- The Wave Plus is BLE-only (no WiFi). To read data, you need a BLE-capable device within approximately 10 meters. The optional Airthings Hub bridges data to the cloud over WiFi, but adds cost.
- **5-minute sensor interval** -- The Wave Plus samples all sensors every 5 minutes. Reading the BLE characteristic more frequently returns stale data. The cloud API similarly provides 5-minute granularity.
- **Radon requires calibration time** -- Radon readings need approximately 24 hours for short-term averages and 30 days for long-term averages to stabilize. Initial readings after installation or battery replacement may be inaccurate.
- **BLE GATT reads are unauthenticated** -- Anyone within BLE range can read sensor data. There is no pairing PIN or encryption requirement for data reads. This simplifies integration but means the data is not private within BLE range.
- **Cloud API rate limit** -- 120 requests per hour. With multiple devices, this requires careful request budgeting. For a single Wave Plus polled every 5 minutes, that is 12 requests per hour -- well within limits.
- **AA batteries** -- Two standard AA batteries. The 16-month battery life is excellent for continuous sensing with BLE advertising.
- **Radon units** -- The API returns radon in Bq/m3 (becquerels per cubic meter). US users typically think in pCi/L (picocuries per liter). Conversion: 1 pCi/L = 37 Bq/m3. The EPA action level is 4 pCi/L (148 Bq/m3).
- **No local hub API** -- Even with the Airthings Hub, there is no local API. The hub simply bridges BLE data to the Airthings cloud. All API access goes through the cloud endpoint.
- **Multiple Wave Plus units** -- The cloud API supports multiple devices per account. Each device is identified by its serial number.
- **View Plus upgrade** -- The Airthings View Plus (model 2960) includes a built-in WiFi hub, PM2.5 sensor, and a display. It is the premium upgrade path but still uses the same cloud API.

## Similar Devices

- [ecowitt-weather-station-gw2000](ecowitt-weather-station-gw2000.md) -- Outdoor weather station (complementary: indoor air vs. outdoor weather)
- [nest-learning-thermostat](nest-learning-thermostat.md) -- Thermostat with temperature/humidity sensing
- [ecobee-smart-thermostat-premium](ecobee-smart-thermostat-premium.md) -- Thermostat with air quality monitoring features
