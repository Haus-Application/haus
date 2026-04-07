# SunPower PVS (Power Vision System) API

## Overview

SunPower PVS solar monitoring gateways expose a local **HTTPS API on port 443** with self-signed certificates. Authentication uses HTTP Basic Auth to obtain a session cookie. The default username is always `ssm_owner`.

## Authentication

### Login

```
GET https://{pvs_ip}/auth?login
Authorization: Basic {base64(ssm_owner:password)}
```

**Response:** Sets a `session` cookie in the response headers.

The password is typically:
- Found on a sticker inside the PVS unit
- Available in the SunPower monitoring app
- Set during installation

### Session Usage

Include the session cookie in all subsequent requests:
```
Cookie: session={token_from_auth}
```

Sessions expire — if you get a 403 or 400, re-authenticate.

## Endpoints

### Live Data

```
GET https://{pvs_ip}/vars?match=livedata&fmt=obj
Cookie: session={token}
```

**Response:** Flat JSON map of path → value strings:
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

Returns ALL system data including per-panel inverter readings. Keys include:
- `/sys/devices/inverter/{n}/sn` — panel serial number
- `/sys/devices/inverter/{n}/ltea3phsumKwh` — lifetime production per panel
- `/sys/devices/inverter/{n}/p3phsumKw` — current production per panel
- `/sys/devices/inverter/{n}/tHtsnkDegc` — panel temperature
- `/sys/info/serialnum` — PVS serial number

### Device-Specific Data

```
GET https://{pvs_ip}/vars?match=sys/devices&fmt=obj
Cookie: session={token}
```

Returns only the device/inverter data subset.

**Note:** The `/cgi-bin/dl_cgi` endpoints (DeviceList, SystemInfo) return 403 Forbidden on most PVS firmware versions. Use `/vars` endpoints instead.

## Error Handling

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | No session cookie present |
| 403 | Session expired or invalid |

Both 400 and 403 mean "re-authenticate" — the PVS firmware uses them inconsistently.

## TLS

The PVS uses a self-signed certificate issued by "SunStrong Management LLC" with CN=pvs.local. Clients must skip TLS verification.

## Haus Integration

- **Discovery:** TLS cert CN=pvs.local, mDNS `_pvs6._tcp`, manufacturer=SunPower
- **Auth:** `POST /api/devices/{ip}/auth` with password
- **Polling:** Fetch `/vars?match=livedata&fmt=obj` every 30 seconds
- **Data:** Production, consumption, grid, battery metrics
