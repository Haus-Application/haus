---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "sunpower-pvs6"
name: "SunPower PVS6 Solar Monitoring Gateway"
manufacturer: "SunPower"
brand: "SunPower"
model: "PVS6"
model_aliases: ["PVS5", "PVS 6"]
device_type: "solar_gateway"
category: "energy"
product_line: "SunPower Equinox"
release_year: 2018
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
  mac_prefixes: ["00:1D:C0", "00:40:AD"]
  mdns_services: ["_pvs6._tcp"]
  mdns_txt_keys: []
  default_ports: [443]
  signature_ports: [443]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^pvs.*", "^PVS.*", "^sunpower.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 443
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "sunpower"
  polling_interval_sec: 30
  websocket_event: "sunpower:state"
  setup_type: "password"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: []

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "session_cookie"
  auth_detail: "HTTP Basic Auth to /auth?login (username: ssm_owner) returns a session cookie. Include Cookie: session={token} in all subsequent requests. Re-authenticate on 400 or 403."
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
  product_page: "https://us.sunpower.com"
  api_docs: ""
  developer_portal: ""
  support: "https://us.sunpower.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "YDR-1001349"

# --- TAGS ---
tags: ["solar", "energy-monitoring", "pvs", "self-signed-tls", "session-auth", "per-panel", "inverter"]
---

# SunPower PVS6 Solar Monitoring Gateway

## What It Is

> The SunPower PVS6 (Power Vision System 6) is the local monitoring gateway for SunPower residential solar installations. It connects to micro-inverters on each solar panel via a proprietary powerline protocol and aggregates real-time production, consumption, grid import/export, and battery data (if equipped). The PVS serves a local HTTPS API on port 443 with a self-signed TLS certificate (CN=pvs.local, issued by SunStrong Management LLC). It provides per-panel granularity including individual inverter production, lifetime energy, and temperature readings.

## How Haus Discovers It

1. **mDNS** -- Advertises as `_pvs6._tcp` on the local network
2. **TLS certificate inspection** -- Self-signed cert with CN=pvs.local identifies SunPower PVS devices
3. **Port probe** -- HTTPS on port 443 with TLS verification disabled
4. **Manufacturer identification** -- Certificate issuer "SunStrong Management LLC" confirms SunPower

## Pairing / Authentication

### Login

```
GET https://{pvs_ip}/auth?login
Authorization: Basic {base64(ssm_owner:password)}
```

**Response:** Sets a `session` cookie in the response headers.

The password is typically:
- Found on a sticker inside the PVS unit
- Available in the SunPower monitoring app
- Set during professional installation

### Session Usage

Include the session cookie in all subsequent requests:
```
Cookie: session={token_from_auth}
```

Sessions expire without a fixed TTL. If you receive a 400 or 403 response, re-authenticate. The PVS firmware uses these status codes inconsistently -- both mean "re-authenticate."

### Haus Auth Flow

`POST /api/devices/{ip}/auth` with the PVS password. Haus stores the password in the database (encrypted) and re-authenticates automatically when the session expires.

## API Reference

### TLS Configuration

The PVS uses a self-signed certificate issued by "SunStrong Management LLC" with CN=pvs.local. All HTTPS clients must skip TLS verification.

### Live Data

```
GET https://{pvs_ip}/vars?match=livedata&fmt=obj
Cookie: session={token}
```

**Response:** Flat JSON map of path to value strings:
```json
{
  "livedata.production.p_3phsum_kw": "4.23",
  "livedata.grid.p_3phsum_kw": "-2.10",
  "livedata.consumption.p_3phsum_kw": "2.13",
  "livedata.battery.soc_pct": "85.0",
  "livedata.battery.p_3phsum_kw": "0.50",
  "livedata.production.e_3phsum_kwh": "45230.5",
  "livedata.grid.e_3phsum_kwh": "12450.2",
  "livedata.consumption.e_3phsum_kwh": "32780.3"
}
```

### Key Metrics

| Path | Description | Unit |
|------|-------------|------|
| `livedata.production.p_3phsum_kw` | Current solar production | kW |
| `livedata.grid.p_3phsum_kw` | Grid power (negative = exporting) | kW |
| `livedata.consumption.p_3phsum_kw` | House consumption | kW |
| `livedata.battery.soc_pct` | Battery state of charge | % |
| `livedata.battery.p_3phsum_kw` | Battery power (negative = charging) | kW |
| `livedata.production.e_3phsum_kwh` | Lifetime production | kWh |
| `livedata.grid.e_3phsum_kwh` | Total grid energy | kWh |
| `livedata.consumption.e_3phsum_kwh` | Total consumption | kWh |

### All System Data (includes panels/inverters)

```
GET https://{pvs_ip}/vars?match=sys&fmt=obj
Cookie: session={token}
```

Returns ALL system data including per-panel inverter readings. Key paths include:
- `/sys/devices/inverter/{n}/sn` -- panel serial number
- `/sys/devices/inverter/{n}/ltea3phsumKwh` -- lifetime production per panel
- `/sys/devices/inverter/{n}/p3phsumKw` -- current production per panel
- `/sys/devices/inverter/{n}/tHtsnkDegc` -- panel/heatsink temperature
- `/sys/info/serialnum` -- PVS serial number

### Device-Specific Data

```
GET https://{pvs_ip}/vars?match=sys/devices&fmt=obj
Cookie: session={token}
```

Returns only the device/inverter data subset.

### Error Handling

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | No session cookie present (re-authenticate) |
| 403 | Session expired or invalid (re-authenticate) |

Both 400 and 403 mean "re-authenticate" -- the PVS firmware uses them inconsistently.

**Note:** The `/cgi-bin/dl_cgi` endpoints (DeviceList, SystemInfo) return 403 Forbidden on most PVS firmware versions. Use `/vars` endpoints instead.

## AI Capabilities

> When chatting with a SunPower PVS, the AI can:
> - **Query live solar data** -- current production, consumption, grid import/export, battery charge
> - **Count panels** -- list all micro-inverters with serial numbers
> - **Per-panel data** -- lifetime production, current output, temperature per inverter
> - **Lifetime stats** -- total kWh produced, grid energy, consumption
> - **Uses stored credentials** -- authenticates automatically with the saved password
>
> The AI uses the `/vars?match=sys&fmt=obj` endpoint (not `/cgi-bin/` which returns 403).
> The AI speaks as the device: "I'm producing 3.5 kW right now with 14 panels. Battery is at 85%."

## Quirks & Notes

- **Self-signed TLS** -- the PVS certificate is self-signed (CN=pvs.local, issuer SunStrong Management LLC); all HTTPS clients must skip TLS verification
- **Session expiry is unpredictable** -- there is no documented TTL for session cookies; always handle 400/403 by re-authenticating
- **400 vs 403** -- the PVS firmware uses both status codes to mean "no valid session"; treat them identically
- **Battery fields optional** -- `livedata.battery.*` fields only appear on systems with battery storage (e.g., SunVault); solar-only systems will not return these fields
- **Values are strings** -- all values in the `/vars` response are strings, not numbers; parse to float for calculations
- **Polling interval** -- 30 seconds is recommended; the PVS updates its internal data approximately every 5-15 seconds
- **Avoid /cgi-bin/** -- the `/cgi-bin/dl_cgi` endpoints documented in older references return 403 on most current firmware versions; use `/vars` endpoints exclusively
- **Professional installation** -- the PVS is installed by SunPower certified installers and is typically mounted near the electrical panel

## Similar Devices

> No directly comparable devices in the Haus knowledge base. Other solar monitoring gateways (Enphase Envoy, SolarEdge) use different protocols but serve a similar role.
