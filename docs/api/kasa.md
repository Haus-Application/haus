# TP-Link Kasa Smart Switch Protocol

## Overview

TP-Link Kasa smart switches, dimmers, and fans communicate over a local TCP protocol on **port 9999**. The protocol uses XOR encryption with key `0xAB` (171) — no authentication required. All communication is JSON wrapped in a 4-byte big-endian length prefix followed by the XOR-encrypted payload.

## Supported Devices

| Model | Type | Capabilities |
|-------|------|-------------|
| HS200 | Switch | on/off |
| HS220 | Dimmer | on/off, brightness (0-100%) |
| HS210 | 3-Way Switch | on/off |
| KS220 | Dimmer | on/off, brightness (0-100%) |

Fan devices are detected by alias containing "fan" and support speed control mapped to brightness tiers (1=25%, 2=50%, 3=75%, 4=100%).

## Protocol

### Encryption

```
XOR Encrypt: key = 171; for each byte: key = key XOR byte; output key
XOR Decrypt: key = 171; for each byte: output = byte XOR key; key = byte
```

All messages are prefixed with a 4-byte big-endian uint32 containing the payload length.

### Connection

- TCP connect to device IP on port 9999
- 1 second connect timeout
- Send encrypted JSON command
- Read 4-byte length prefix, then read that many bytes
- Decrypt response

## Endpoints

### Query Device State

```json
{"system":{"get_sysinfo":{}}}
```

**Response fields:**
- `alias` — device display name (e.g., "Kitchen Lights")
- `model` — hardware model (e.g., "HS220(US)")
- `relay_state` — 0=off, 1=on
- `brightness` — 0-100 (dimmers only)
- `dev_name` — device description
- `mac` — MAC address
- `deviceId` — unique device ID
- `on_time` — seconds since last power on

### Turn On/Off

```json
{"system":{"set_relay_state":{"state":1}}}
```

- `state: 1` = on
- `state: 0` = off

### Set Brightness (Dimmers Only)

```json
{"smartlife.iot.dimmer":{"set_brightness":{"brightness":75}}}
```

- `brightness` — integer 0-100
- Only works on HS220/KS220 dimmers
- Device must be on for brightness to take effect

### Set Fan Speed

Fan speed is mapped to brightness values:

| Speed | Brightness | Label |
|-------|-----------|-------|
| 1 | 25 | Low |
| 2 | 50 | Medium |
| 3 | 75 | High |
| 4 | 100 | Max |

Use the set_brightness command with the mapped value.

## Haus Integration

- **Discovery:** TCP connect scan on port 9999 across subnet
- **Polling:** Every 10 seconds via `get_sysinfo`
- **Control:** `POST /api/kasa/devices/{ip}/state`, `/brightness`, `/fan-speed`
- **WebSocket:** Live state broadcast via `kasa:state` event
