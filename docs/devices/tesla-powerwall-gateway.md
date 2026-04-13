---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "tesla-powerwall-gateway"
name: "Tesla Powerwall + Gateway"
manufacturer: "Tesla Inc."
brand: "Tesla"
model: "Powerwall Gateway"
model_aliases: ["Tesla Energy Gateway", "TEG", "Gateway 2", "Powerwall 2", "Powerwall+", "Powerwall 3"]
device_type: "energy_gateway"
category: "energy"
product_line: "Tesla Energy"
release_year: 2016
discontinued: false
price_range: "$$$$"

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
  mac_prefixes: ["DC:44:27", "4C:FC:AA", "98:ED:5C"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [443]
  signature_ports: [443]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^teg-.*", "^Tesla.*", "^1118431.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 443
    path: "/api/status"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "gateway_din"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "tesla"
  polling_interval_sec: 15
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
  auth_method: "password"
  auth_detail: "Local HTTPS API on port 443 with self-signed certificate. Authentication via POST /api/login/Basic with email and password (installer or customer credentials). Returns an AuthCookie token. Newer firmware versions may require Tesla account OAuth token from auth.tesla.com."
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
  product_page: "https://www.tesla.com/powerwall"
  api_docs: ""
  developer_portal: ""
  support: "https://www.tesla.com/support/energy"
  community_forum: "https://teslamotorsclub.com/tmc/forums/powerwall.156/"
  image_url: ""
  fcc_id: "2AEIM-1099200"

# --- TAGS ---
tags: ["energy-monitoring", "battery", "solar", "self-signed-tls", "local-api", "tesla", "powerwall", "grid-backup", "time-of-use", "storm-watch"]
---

# Tesla Powerwall + Gateway

## What It Is

> The Tesla Powerwall is a home battery system paired with the Tesla Energy Gateway (TEG), which manages energy flow between solar panels, the battery, the home, and the grid. The Gateway serves as the central monitoring and control point, providing a local HTTPS REST API on port 443 with a self-signed certificate. It reports real-time data for solar production, battery charge/discharge, home consumption, and grid import/export. The Powerwall system supports backup power during outages, time-of-use energy optimization, and Storm Watch (automatic pre-charging before severe weather). The Gateway connects to Tesla's cloud for remote monitoring via the Tesla app, but the local API functions independently of internet connectivity.

## How Haus Discovers It

1. **OUI match** -- Tesla MAC prefixes: `DC:44:27`, `4C:FC:AA`, `98:ED:5C`
2. **Hostname pattern** -- DHCP hostname starts with `teg-` (Tesla Energy Gateway) or contains `Tesla`
3. **Port probe** -- HTTPS on port 443 with self-signed certificate
4. **HTTP fingerprint** -- `GET /api/status` returns JSON with `gateway_din` field confirming a Tesla Gateway
5. **TLS certificate inspection** -- Self-signed cert with organization "Tesla Motors" or "Tesla Inc"

## Pairing / Authentication

### Local API Authentication (Gateway 2 / Older Firmware)

```
POST https://{ip}/api/login/Basic
Content-Type: application/json

{
  "username": "customer",
  "password": "S{last5_of_password}",
  "email": "user@example.com",
  "force_sm_off": false
}
```

**Response:**
```json
{
  "email": "user@example.com",
  "firstname": "Tesla",
  "lastname": "Energy",
  "roles": ["Home_Owner"],
  "token": "eyJ...",
  "provider": "Basic",
  "loginTime": "2024-04-01T12:00:00.000Z"
}
```

The password for the `customer` login is typically the last 5 characters of the Gateway serial number, prefixed with `S`. For example, if the serial ends in `12345`, the password is `S12345`.

Include the token in subsequent requests:
```
Authorization: Bearer {token}
```

### Newer Firmware / Tesla Account Auth

Newer Gateway firmware versions require a Tesla account OAuth token:

1. Obtain an OAuth access token from `https://auth.tesla.com/oauth2/v3/token`
2. Use the Tesla Fleet API or the token directly with the local Gateway
3. Some endpoints may work without auth on newer firmware

### Haus Auth Flow

`POST /api/devices/{ip}/auth` with Gateway password (last 5 of serial) or Tesla account credentials. Haus authenticates and stores the token.

## API Reference

### System Status

**Gateway status (no auth required):**
```
GET https://{ip}/api/status
```

**Response:**
```json
{
  "din": "1232100-00-E--TG123456789012",
  "start_time": "2024-01-01 00:00:00 +0000",
  "up_time_seconds": "7890123s",
  "is_new": false,
  "version": "24.4.0 abc123",
  "git_hash": "abc123def456",
  "commission_count": 1,
  "device_type": "teg",
  "teg_type": "teg",
  "sync_type": "v2.1",
  "cellular_disabled": false,
  "can_reboot": true
}
```

### Aggregate Meters

**Energy flow data (requires auth):**
```
GET https://{ip}/api/meters/aggregates
Authorization: Bearer {token}
```

**Response:**
```json
{
  "site": {
    "last_communication_time": "2024-04-01T12:00:00Z",
    "instant_power": -2100.5,
    "instant_reactive_power": -123.4,
    "instant_apparent_power": 2104.1,
    "frequency": 60.01,
    "energy_exported": 12450234.56,
    "energy_imported": 3200123.45,
    "instant_average_voltage": 240.2,
    "instant_average_current": 8.75,
    "i_a_current": 4.3,
    "i_b_current": 4.45,
    "i_c_current": 0,
    "last_phase_voltage_communication_time": "2024-04-01T12:00:00Z",
    "timeout": 1500000000
  },
  "battery": {
    "instant_power": 500.0,
    "instant_reactive_power": 0,
    "instant_apparent_power": 500.0,
    "frequency": 60.01,
    "energy_exported": 5600123.45,
    "energy_imported": 5200234.56
  },
  "load": {
    "instant_power": 2130.5,
    "instant_reactive_power": 100.2,
    "instant_apparent_power": 2133.0,
    "frequency": 60.01,
    "energy_exported": 0,
    "energy_imported": 32780345.67
  },
  "solar": {
    "instant_power": 4230.0,
    "instant_reactive_power": -50.0,
    "instant_apparent_power": 4230.3,
    "frequency": 60.01,
    "energy_exported": 45230567.89,
    "energy_imported": 0
  }
}
```

**Key fields:**
- `site.instant_power` -- Grid power (negative = exporting to grid)
- `battery.instant_power` -- Battery power (positive = discharging, negative = charging)
- `load.instant_power` -- Home consumption
- `solar.instant_power` -- Solar production

### Battery State of Charge

```
GET https://{ip}/api/system_status/soe
Authorization: Bearer {token}
```

**Response:**
```json
{
  "percentage": 85.2345
}
```

### Grid Status

```
GET https://{ip}/api/system_status/grid_status
Authorization: Bearer {token}
```

**Response:**
```json
{
  "grid_status": "SystemGridConnected",
  "grid_services_active": true
}
```

Grid status values: `SystemGridConnected`, `SystemIslandedActive` (backup mode), `SystemTransitionToGrid`

### Battery Operation Mode

```
GET https://{ip}/api/operation
Authorization: Bearer {token}
```

**Response:**
```json
{
  "real_mode": "self_consumption",
  "backup_reserve_percent": 20.0,
  "freq_shift_load_shed_soe": 0,
  "freq_shift_load_shed_delta_f": 0
}
```

Operation modes: `self_consumption` (maximize solar use), `backup` (maximize battery reserve), `autonomous` (time-of-use optimization)

### Powerwalls

```
GET https://{ip}/api/powerwalls
Authorization: Bearer {token}
```

**Response:**
```json
{
  "powerwalls": [
    {
      "Type": "",
      "PackagePartNumber": "2012170-25-E",
      "PackageSerialNumber": "TG123456789012",
      "type": "acpw",
      "grid_state": "Grid_Compliant",
      "grid_reconnection_time_seconds": 0,
      "under_phase_detection": false,
      "updating": false,
      "commissioning_diagnostic": {},
      "energy_left": 10850.0,
      "total_energy": 13500.0,
      "nominal_energy_remaining": 11200.0
    }
  ],
  "gateway_din": "1232100-00-E--TG123456789012",
  "sync": {
    "updating": false,
    "commissioning_diagnostic": {}
  },
  "phase_detection": {}
}
```

### Complete Endpoint List

| Endpoint | Auth | Description |
|----------|------|-------------|
| `/api/status` | No | Gateway version, uptime, DIN |
| `/api/meters/aggregates` | Yes | Real-time power flow (solar, battery, grid, load) |
| `/api/system_status/soe` | Yes | Battery state of charge percentage |
| `/api/system_status/grid_status` | Yes | Grid connected/islanded status |
| `/api/operation` | Yes | Current operation mode and backup reserve |
| `/api/powerwalls` | Yes | Individual Powerwall details and energy |
| `/api/site_info` | Yes | Site name, timezone, tariff info |
| `/api/site_info/site_name` | Yes | Just the site name |
| `/api/networks` | Yes | Network configuration (WiFi, Ethernet) |
| `/api/system/update/status` | Yes | Firmware update status |
| `/api/solar_powerwall` | No | Solar + Powerwall summary |

## AI Capabilities

> AI integration planned. When available:
> - Report real-time solar production, battery charge, grid usage, and home consumption
> - Report battery state of charge percentage
> - Report grid status (connected vs backup/islanded)
> - Display current operation mode (self-consumption, backup, time-of-use)
> - Per-Powerwall energy data (for multi-Powerwall installations)
> - Alert on grid outages (grid status changes to islanded)
> - The AI speaks as the device: "I'm at 85% charge, producing 4.2 kW solar, and exporting 2.1 kW to the grid. The house is using 2.1 kW."

## Quirks & Notes

- **Self-signed TLS** -- The Gateway uses a self-signed certificate; all HTTPS clients must skip TLS verification
- **Auth evolution** -- Tesla has changed authentication requirements across firmware versions; older Gateways use simple password auth, newer ones may require Tesla account OAuth tokens
- **Password format** -- The local `customer` password is typically `S` + last 5 characters of the Gateway serial number; installer password uses different credentials
- **Token expiry** -- Auth tokens expire; re-authenticate on 401 or 403 responses
- **Negative = exporting** -- In `site.instant_power`, negative values mean exporting to the grid; in `battery.instant_power`, negative means charging
- **Multiple Powerwalls** -- A single Gateway can manage multiple Powerwall units (stacked for more capacity); each appears in the `/api/powerwalls` response
- **Grid services** -- If enrolled in a Tesla Virtual Power Plant (VPP) program, `grid_services_active` will be true and Tesla may remotely discharge the battery during grid events
- **Storm Watch** -- The Gateway automatically enters full backup mode when severe weather is forecasted in the area (via Tesla cloud)
- **Powerwall 3** -- The Powerwall 3 integrates the Gateway directly into the battery unit; the same local API applies
- **Firmware updates** -- Tesla pushes firmware updates automatically; the API may change between versions
- **Rate limiting** -- The Gateway handles approximately 1-2 requests per second comfortably; higher rates may cause timeouts
- **Energy values** -- `energy_exported` and `energy_imported` are cumulative Wh values since installation; calculate delta for period totals
- **Phase data** -- Split-phase data (i_a_current, i_b_current) available for 240V installations; useful for per-leg load balancing analysis

## Similar Devices

> - [Enphase Envoy-S](enphase-envoy-s.md) -- Enphase solar monitoring gateway (solar monitoring, no battery)
> - [SunPower PVS6](sunpower-pvs6.md) -- SunPower solar monitoring gateway
