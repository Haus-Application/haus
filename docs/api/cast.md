# Google Cast Local API

## Overview

Google Cast devices (Chromecast, Nest Hub, NVIDIA Shield, Google Home, Cast-enabled TVs) expose a local **HTTP API on port 8008** for device information. Control requires the Cast SDK. The info API requires no authentication.

## Supported Devices

- Google Chromecast (all generations)
- Google Nest Hub / Hub Max
- Google Home / Home Mini
- NVIDIA Shield TV
- LG webOS TVs with Cast
- Any device advertising `_googlecast._tcp` via mDNS

## Info Endpoints

Base URL: `http://{device_ip}:8008`

### Device Info

```
GET /setup/eureka_info?params=name,build_info,detail,device_info,opt_in
```

**Response:**
```json
{
  "name": "Living Room TV",
  "build_info": {
    "cast_build_revision": "1.56.313396",
    "cast_control_version": 1
  },
  "device_info": {
    "model_name": "Chromecast",
    "manufacturer": "Google Inc."
  },
  "detail": {
    "locale": {
      "display_string": "English (United States)"
    }
  }
}
```

### Configured Networks

```
GET /setup/configured_networks
```

Returns WiFi network configuration.

### Supported App IDs

```
GET /setup/supported_app_ids
```

Returns list of supported Cast app IDs.

## Control (Cast SDK Required)

The info API on port 8008 is read-only. To control Cast devices (play media, adjust volume, etc.), you need the Google Cast SDK which uses a proprietary protocol on port 8009 (TLS).

### Cast Protocol (Port 8009)

- Protocol: Protobuf over TLS
- Authentication: TLS with device certificate
- Capabilities: media playback, volume, app launch
- SDK: https://developers.google.com/cast

## Additional Ports

| Port | Protocol | Purpose |
|------|----------|---------|
| 8008 | HTTP | Device info API (read-only) |
| 8009 | TLS/Protobuf | Cast control protocol |
| 8443 | HTTPS | Secure device info |

## Thread Border Router (Nest Hub)

Google Nest Hub Max devices also function as Thread border routers, advertising `_meshcop._udp` via mDNS. This enables Matter-over-Thread device commissioning.

## Haus Integration

- **Discovery:** mDNS `_googlecast._tcp` + HTTP probe on port 8008
- **Info:** `GET /setup/eureka_info` for device name and model
- **Control:** Not yet implemented (requires Cast SDK)
- **Status:** Read-only — can identify but not control

## AI Chat Capabilities

When chatting with a Cast device, the AI can:
- **Report device identity** — name, model from eureka_info
- **Report network status** — on the network, responding on port 8008
- Cast control (play/pause/volume) requires the Cast SDK and is not yet implemented

For Nest Hub Max and Nest displays that are also Google Nest devices: if the user has connected their Google account via OAuth, the AI can access camera feeds and describe what the camera sees using vision AI (live snapshot analysis).
