# Philips Hue Bridge API (v2)

## Overview

Philips Hue bridges expose a local HTTPS REST API (CLIP v2) on **port 443** with self-signed certificates. Authentication requires a one-time link-button pairing to obtain an API key. All subsequent requests use this key in the `hue-application-key` header.

## Pairing Flow

1. User presses the physical link button on the Hue bridge
2. Within 30 seconds, send: `POST https://{bridge_ip}/api`
3. Body: `{"devicetype":"haus#app","generateclientkey":true}`
4. Response contains `username` — this is your API key
5. Store and use for all future requests via header: `hue-application-key: {username}`

## Supported Resources

| Resource | Description |
|----------|-------------|
| Light | Individual bulbs — on/off, brightness, color XY, color temperature |
| Room | Group of lights — control all at once via grouped_light |
| Scene | Saved light states — recall to set all lights in a room |
| Grouped Light | Virtual light representing all lights in a room/zone |

## Endpoints

All endpoints use HTTPS with the `hue-application-key` header.

### List Lights

```
GET /clip/v2/resource/light
```

**Response:** Array of light objects with:
- `id` — UUID
- `metadata.name` — light name
- `on.on` — boolean
- `dimming.brightness` — 0-100
- `color.xy` — CIE xy color coordinates `{x: 0.3127, y: 0.3290}`
- `color_temperature.mirek` — color temperature in mirek (153-500)
- `owner.rid` — room/zone this light belongs to

### Control a Light

```
PUT /clip/v2/resource/light/{id}
```

**Body (all fields optional):**
```json
{
  "on": {"on": true},
  "dimming": {"brightness": 75.0},
  "color": {"xy": {"x": 0.4578, "y": 0.4101}}
}
```

### List Rooms

```
GET /clip/v2/resource/room
```

**Response:** Array of room objects with:
- `id` — UUID
- `metadata.name` — room name
- `children` — array of device references
- `services` — includes `grouped_light` reference for room-level control

### Control a Room

```
PUT /clip/v2/resource/grouped_light/{grouped_light_id}
```

Same body format as individual light control — affects all lights in the room.

### List Scenes

```
GET /clip/v2/resource/scene
```

**Response:** Array of scene objects with:
- `id` — UUID
- `metadata.name` — scene name
- `group.rid` — room this scene belongs to

### Activate a Scene

```
PUT /clip/v2/resource/scene/{id}
```

```json
{"recall": {"action": "active"}}
```

## Color Reference

| Color | CIE x | CIE y |
|-------|-------|-------|
| Warm White | 0.4578 | 0.4101 |
| Cool White | 0.3127 | 0.3290 |
| Red | 0.6750 | 0.3220 |
| Blue | 0.1532 | 0.0475 |
| Green | 0.1700 | 0.7000 |
| Purple | 0.2703 | 0.1398 |
| Orange | 0.5614 | 0.3944 |
| Pink | 0.3944 | 0.1990 |

## Haus Integration

- **Discovery:** `GET https://discovery.meethue.com/` returns local bridges
- **Pairing:** `POST /api/hue/pair` with bridge IP
- **Polling:** Every 5 seconds for lights, rooms, scenes
- **Control:** `PUT /api/hue/lights/{id}`, `/rooms/{id}`, `/scenes/{id}/activate`
- **WebSocket:** Live state broadcast via `hue:state` event

## Documentation

- Official: https://developers.meethue.com/develop/hue-api-v2/
