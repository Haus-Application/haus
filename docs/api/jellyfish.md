# JellyFish Lighting Controller API

## Overview

JellyFish outdoor LED lighting controllers expose a **WebSocket API on port 9000**. The protocol is fully asynchronous — commands are sent as JSON and responses arrive independently. No authentication is required.

## Connection

```
WebSocket: ws://{controller_ip}:9000/
```

- Handshake timeout: 5 seconds
- Messages are JSON objects
- Commands use `"cmd": "toCtlrGet"` for queries and `"cmd": "toCtlrSet"` for control
- Responses echo the requested data type as the top-level key

## Endpoints

### Get Controller Name

**Send:**
```json
{"cmd": "toCtlrGet", "get": [["ctlrName"]]}
```

**Response:**
```json
{"ctlrName": "Happy lites"}
```

### Get Zones

**Send:**
```json
{"cmd": "toCtlrGet", "get": [["zones"]]}
```

**Response:**
```json
{
  "zones": {
    "Zone": {"numPixels": 300},
    "Zone1": {"numPixels": 150}
  }
}
```

Each zone represents a physical LED strip segment with a pixel count.

### Get Pattern List

**Send:**
```json
{"cmd": "toCtlrGet", "get": [["patternFileList"]]}
```

**Response:**
```json
{
  "patternFileList": [
    {"folders": "Easter", "name": "Easter Colors", "readOnly": true},
    {"folders": "Accent", "name": "All Lights Warm White 3000K", "readOnly": false},
    {"folders": "Holiday", "name": "Christmas", "readOnly": true}
  ]
}
```

Pattern file paths are constructed as `{folders}/{name}` (e.g., `"Accent/All Lights Warm White 3000K"`).

### Get Zone State

**Send:**
```json
{"cmd": "toCtlrGet", "get": [["runPattern", "Zone1", "Zone"]]}
```

**Response (one per zone):**
```json
{
  "runPattern": {
    "state": 1,
    "zoneName": ["Zone1"],
    "file": "Accent/White",
    "id": ""
  }
}
```

- `state: 1` = playing, `state: 0` = off

### Play a Pattern (Turn On)

**Send:**
```json
{
  "cmd": "toCtlrSet",
  "runPattern": {
    "state": 1,
    "zoneName": ["Zone1", "Zone"],
    "file": "Accent/All Lights Warm White 3000K",
    "id": "",
    "data": ""
  }
}
```

- `state: 1` = start playing
- `zoneName` = array of zone names to activate
- `file` = pattern file path (`folder/name`)

### Turn Off

**Send:**
```json
{
  "cmd": "toCtlrSet",
  "runPattern": {
    "state": 0,
    "zoneName": ["Zone1", "Zone"],
    "file": "",
    "id": "",
    "data": ""
  }
}
```

## Firmware Update UI

The controller also serves a web UI on **port 8080** for firmware updates. This is an HTML page with jQuery for checking and applying updates.

## Haus Integration

- **Discovery:** mDNS service `_jellyfishV2._tcp` + TCP check on port 9000
- **Probe:** WebSocket connect to get zones, patterns, controller name
- **Control:** `POST /api/devices/{ip}/jellyfish` with `{action, zones, pattern}`
- **AI Chat:** The AI can query zone states and control patterns via WebSocket
- **No polling:** Stateless — query on demand via WebSocket
- **No auth:** No authentication required — direct local WebSocket connection

## AI Chat Capabilities

When chatting with a JellyFish device, the AI can:
- Query current zone states (on/off, active pattern)
- Turn zones on/off
- Play specific patterns on zones
- List available patterns and zones

The AI uses a direct WebSocket connection to port 9000 for real-time queries.
