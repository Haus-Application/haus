---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "enphase-envoy-s"
name: "Enphase Envoy-S Metered Gateway"
manufacturer: "Enphase Energy Inc."
brand: "Enphase"
model: "Envoy-S Metered"
model_aliases: ["ENV-S-WM-230", "ENV-S-WB-230", "Envoy-S Standard", "Envoy", "IQ Gateway", "IQ Gateway Metered"]
device_type: "solar_gateway"
category: "energy"
product_line: "Enphase IQ"
release_year: 2016
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "ethernet"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["00:1D:C0", "00:11:76"]
  mdns_services: ["_enphase-envoy._tcp"]
  mdns_txt_keys: ["serialnum", "protovers"]
  default_ports: [443, 80]
  signature_ports: [443]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^envoy.*", "^Envoy.*", "^enphase.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 443
    path: "/info.xml"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "envoy"
    headers: {}
  - port: 80
    path: "/production.json"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "production"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "enphase"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "password"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["energy_monitoring"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "Authentication varies by firmware version. Older firmware (< D7.x): HTTP Digest auth on port 80 with username 'envoy' and password = last 6 digits of serial number. Newer firmware (D7.x+/IQ Gateway): JWT token obtained from Enphase cloud (entrez.enphaseenergy.com), then used locally on port 443. Token valid for ~12 months."
  base_url_template: "https://{ip}"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "gateway"
  power_source: "hardwired"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://enphase.com/homeowners/home-solar-systems"
  api_docs: ""
  developer_portal: "https://developer-v4.enphase.com"
  support: "https://support.enphase.com"
  community_forum: "https://support.enphase.com/s/community"
  image_url: ""
  fcc_id: "TUF-ENV-S"

# --- TAGS ---
tags: ["solar", "energy-monitoring", "microinverter", "per-panel", "self-signed-tls", "jwt-auth", "digest-auth", "enphase", "envoy", "iq-gateway"]
---

# Enphase Envoy-S Metered Gateway

## What It Is

> The Enphase Envoy-S (now rebranded as IQ Gateway) is the central monitoring gateway for Enphase residential and commercial solar installations. It communicates with each Enphase microinverter (one per solar panel) via powerline communication (PLC) and aggregates real-time production data, consumption data (if CTs are installed), and per-panel performance metrics. The Envoy connects to the home network via WiFi or Ethernet and uploads data to Enphase's Enlighten cloud platform. It also serves a local HTTPS API on port 443 with a self-signed certificate, providing real-time production and consumption data without internet dependency. The "Metered" variant includes consumption monitoring CTs (current transformers) for net energy tracking.

## How Haus Discovers It

1. **mDNS** -- Advertises as `_enphase-envoy._tcp` with TXT records containing serial number and protocol version
2. **OUI match** -- Enphase MAC prefixes: `00:1D:C0`, `00:11:76`
3. **Hostname pattern** -- DHCP hostname starts with `envoy` or `Envoy`
4. **HTTP fingerprint** -- `GET /info.xml` on port 443 returns XML with device info; `GET /production.json` on port 80 returns production data (older firmware)
5. **TLS certificate inspection** -- Self-signed cert with organization containing "Enphase"

## Pairing / Authentication

### Older Firmware (pre-D7.x / Envoy-S)

Simple HTTP Digest authentication on port 80:

- **Username:** `envoy`
- **Password:** Last 6 digits of the Envoy serial number (found on the device label or in the Enlighten app)

```
GET http://{ip}/production.json
Authorization: Digest username="envoy", ...
```

Some endpoints on older firmware are available without authentication:
```
GET http://{ip}/production.json      (no auth on some firmware versions)
GET http://{ip}/api/v1/production    (no auth)
```

### Newer Firmware (D7.x+ / IQ Gateway)

JWT token-based authentication:

1. **Get token from Enphase cloud:**
```
POST https://entrez.enphaseenergy.com/tokens
Content-Type: application/json

{
  "session_id": "{session_id_from_login}",
  "serial_num": "{envoy_serial}",
  "username": "{enlighten_email}"
}
```

First, log in to get a session:
```
POST https://enlighten.enphaseenergy.com/login/login.json?user[email]={email}&user[password]={password}
```

2. **Use token locally:**
```
GET https://{ip}/api/v1/production
Authorization: Bearer {jwt_token}
```

The JWT token is valid for approximately 12 months. Store it and refresh before expiry.

### Haus Auth Flow

For older firmware: `POST /api/devices/{ip}/auth` with the Envoy serial number (password derived automatically).
For newer firmware: `POST /api/devices/{ip}/auth` with Enlighten email and password. Haus obtains the JWT token from Enphase cloud and stores it for local API access.

## API Reference

### Device Info

```
GET https://{ip}/info.xml
```

Returns XML with serial number, firmware version, and device type. No auth required on most firmware versions.

### Production Summary (Legacy, Port 80)

```
GET http://{ip}/production.json
```

**Response:**
```json
{
  "production": [
    {
      "type": "inverters",
      "activeCount": 20,
      "readingTime": 1712000000,
      "wNow": 4230,
      "whLifetime": 45230500
    },
    {
      "type": "eim",
      "activeCount": 1,
      "measurementType": "production",
      "readingTime": 1712000000,
      "wNow": 4230.45,
      "whLifetime": 45230567.89,
      "vahLifetime": 50123456.78,
      "varhLeadLifetime": 1234567.89,
      "varhLagLifetime": 2345678.90,
      "rmsCurrent": 17.65,
      "rmsVoltage": 239.8,
      "reactPwr": 123.45,
      "apprntPwr": 4250.0,
      "pwrFactor": 0.99
    }
  ],
  "consumption": [
    {
      "type": "eim",
      "activeCount": 1,
      "measurementType": "total-consumption",
      "readingTime": 1712000000,
      "wNow": 2130.50,
      "whLifetime": 32780345.67
    },
    {
      "type": "eim",
      "activeCount": 1,
      "measurementType": "net-consumption",
      "readingTime": 1712000000,
      "wNow": -2099.95,
      "whLifetime": 12450234.56
    }
  ],
  "storage": []
}
```

### Production API v1

```
GET https://{ip}/api/v1/production
Authorization: Bearer {token}
```

**Response:**
```json
{
  "wattHoursToday": 25430,
  "wattHoursSevenDays": 178500,
  "wattHoursLifetime": 45230567,
  "wattsNow": 4230
}
```

### Per-Inverter Data

```
GET https://{ip}/api/v1/production/inverters
Authorization: Bearer {token}
```

**Response:**
```json
[
  {
    "serialNumber": "122230012345",
    "lastReportDate": 1712000000,
    "devType": 1,
    "lastReportWatts": 215,
    "maxReportWatts": 290
  },
  {
    "serialNumber": "122230012346",
    "lastReportDate": 1712000000,
    "devType": 1,
    "lastReportWatts": 220,
    "maxReportWatts": 295
  }
]
```

### Inventory

```
GET https://{ip}/inventory.json
Authorization: Bearer {token}
```

Returns complete device inventory including microinverter serial numbers, firmware versions, and communication status.

### Key Metrics

| Endpoint | Metric | Unit | Auth Required |
|----------|--------|------|---------------|
| `/production.json` (production.wNow) | Current solar production | Watts | Varies |
| `/production.json` (consumption.wNow, total) | Current consumption | Watts | Varies |
| `/production.json` (consumption.wNow, net) | Net grid (negative = exporting) | Watts | Varies |
| `/api/v1/production` (wattsNow) | Current production | Watts | Yes (D7+) |
| `/api/v1/production` (wattHoursToday) | Today's production | Wh | Yes (D7+) |
| `/api/v1/production/inverters` | Per-panel production | Watts | Yes (D7+) |

## AI Capabilities

> AI integration planned. When available:
> - Report real-time solar production, consumption, and net grid export/import
> - Per-panel production monitoring (identify underperforming panels)
> - Daily, weekly, and lifetime energy production totals
> - Power factor and voltage reporting
> - Alert on offline or underperforming microinverters
> - The AI speaks as the device: "I'm producing 4.2 kW right now across 20 panels. Panel 15 is running a bit low at 180W."

## Quirks & Notes

- **Self-signed TLS** -- The Envoy uses a self-signed certificate; all HTTPS clients must skip TLS verification
- **Auth migration** -- Enphase has been migrating from simple digest auth to JWT tokens; firmware updates may change auth requirements without notice
- **JWT token lifetime** -- Tokens from entrez.enphaseenergy.com are valid for approximately 12 months; plan for automatic renewal
- **production.json availability** -- On older firmware, `/production.json` is available without auth on port 80; on newer firmware, it may redirect to HTTPS and require JWT auth
- **Metered vs Standard** -- The "Metered" Envoy-S includes consumption CTs; the "Standard" variant only reports production (no consumption/net data)
- **EIM vs Inverters** -- Two data types in production.json: `inverters` (aggregated from microinverter reports, updated every 5-15 minutes) and `eim` (from the electricity meter CTs, updated every second)
- **Negative net consumption** -- Negative values in net-consumption mean exporting to the grid
- **Powerline communication** -- Microinverters communicate with the Envoy via PLC (powerline communication) over the home's electrical wiring, not WiFi
- **Storage data** -- The `storage` array in production.json reports Enphase IQ Battery data if batteries are installed (similar to SunPower battery fields)
- **Rate limiting** -- The Envoy has limited processing power; polling more frequently than every 5 seconds may cause timeouts or slow responses
- **IQ Gateway rebrand** -- Newer Enphase gateways are branded "IQ Gateway" or "IQ Gateway Metered" but use the same API structure as the Envoy-S

## Similar Devices

> - [SunPower PVS6](sunpower-pvs6.md) -- SunPower's solar monitoring gateway with similar local API
> - [Tesla Powerwall Gateway](tesla-powerwall-gateway.md) -- Tesla's energy gateway with local REST API
