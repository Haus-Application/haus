# Yamaha MusicCast / Extended Control API

## Overview

Yamaha AV receivers expose a local **HTTP REST API on port 80** called "Extended Control" (also known as MusicCast API). No authentication required. All commands are GET requests with parameters in the query string.

## Supported Devices

- RX-V6A, RX-V4A series
- RX-A series (Aventage)
- Any Yamaha receiver with MusicCast / Network Module

## Endpoints

Base URL: `http://{receiver_ip}/YamahaExtendedControl/v1`

### System Info

```
GET /system/getDeviceInfo
```

**Response:**
```json
{
  "response_code": 0,
  "model_name": "RX-V6A",
  "destination": "U",
  "device_id": "04098643797E",
  "system_version": 1.70,
  "api_version": 2.15,
  "serial_number": "Y726554RT"
}
```

### Get Main Zone Status

```
GET /main/getStatus
```

**Response:**
```json
{
  "response_code": 0,
  "power": "standby",
  "volume": 45,
  "mute": false,
  "input": "hdmi1",
  "sound_program": "straight"
}
```

- `power`: "on", "standby"
- `volume`: 0-161 (0 = -80dB, 161 = +16.5dB)
- `input`: "hdmi1"-"hdmi7", "av1"-"av7", "audio1"-"audio4", "tuner", "bluetooth", "airplay", "spotify", "net_radio"

### Power Control

```
GET /main/setPower?power=on
GET /main/setPower?power=standby
GET /main/setPower?power=toggle
```

### Volume Control

```
GET /main/setVolume?volume=50
GET /main/setVolume?volume=up&step=1
GET /main/setVolume?volume=down&step=1
```

### Mute

```
GET /main/setMute?enable=true
GET /main/setMute?enable=false
```

### Input Selection

```
GET /main/setInput?input=hdmi1
GET /main/setInput?input=bluetooth
GET /main/setInput?input=airplay
```

### Sound Program

```
GET /main/setSoundProgram?program=straight
GET /main/setSoundProgram?program=surr_decoder
GET /main/getSoundProgramList
```

### Multi-Zone (Zone 2, Zone 3)

Replace `main` with `zone2` or `zone3`:
```
GET /zone2/getStatus
GET /zone2/setPower?power=on
GET /zone2/setVolume?volume=30
GET /zone2/setInput?input=audio1
```

### Network Standby

```
GET /system/getNetworkStandby
GET /system/setNetworkStandby?standby=on
```

## Response Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Initializing |
| 2 | Internal Error |
| 3 | Invalid Request |
| 4 | Invalid Parameter |
| 5 | Guarded |

## Haus Integration

- **Discovery:** HTTP Server header contains `Network_Module` or model name (e.g., `RX-V6A`)
- **Probe:** `GET /system/getDeviceInfo` + `GET /main/getStatus`
- **Control:** Direct HTTP GET to receiver API endpoints
- **No auth:** All endpoints are open on the local network
