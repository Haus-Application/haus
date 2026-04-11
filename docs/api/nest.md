# Google Nest SDM API Integration

## Overview

Haus integrates with Google Nest devices via the Smart Device Management (SDM) API. Nest cameras, thermostats, doorbells, and displays are cloud-only — they have no local API. All communication goes through Google's servers using OAuth2 tokens.

**Important:** Nest devices don't respond to local network probes. They show up in ARP/ping scans by their Google MAC addresses (prefixes: `7c:10:15`, `b0:09:da`, `18:b4:30`, `18:7f:88`, `f8:0f:f9`) but have no open ports. Control requires the SDM cloud API.

## Setup

### Cost

- **$5 one-time fee** per Google account for Device Access registration
- This is a developer fee — end users don't pay it in Commercial tier
- Google is NOT currently accepting new Commercial Development applications

### Prerequisites

1. Google Cloud project with SDM API enabled
2. Device Access project from https://console.nest.google.com/device-access
3. OAuth2 Web Application client (client ID + secret)
4. User's Google account linked to their Nest devices
5. go2rtc for camera live streaming (handles WebRTC negotiation)

### Environment Variables

```
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-...
GOOGLE_PROJECT_ID=bf59ddff-ff13-4b97-8123-4aca84e315dd
```

## OAuth2 Flow

Uses Google's Nest-specific partner connection flow (NOT standard Google OAuth).

### Step 1: Redirect to Authorization

```
GET https://nestservices.google.com/partnerconnections/{project_id}/auth
  ?redirect_uri={redirect_uri}
  &access_type=offline
  &prompt=consent
  &client_id={client_id}
  &response_type=code
  &scope=https://www.googleapis.com/auth/sdm.service
```

**Critical:** The URL is `nestservices.google.com`, NOT `accounts.google.com`. `access_type=offline` is required for refresh tokens. `prompt=consent` ensures a refresh token is always returned.

### Step 2: Exchange Code for Tokens

```
POST https://www.googleapis.com/oauth2/v4/token
Content-Type: application/x-www-form-urlencoded

client_id={client_id}&client_secret={client_secret}&code={code}&grant_type=authorization_code&redirect_uri={redirect_uri}
```

### Step 3: Refresh Token

Access tokens expire after 1 hour. Refresh tokens don't expire unless revoked.

```
POST https://www.googleapis.com/oauth2/v4/token

client_id={client_id}&client_secret={client_secret}&refresh_token={token}&grant_type=refresh_token
```

**Note:** Google's refresh response does NOT always return a new refresh token. Preserve the original.

## Device Types

| Type | SDM String | Haus device_type |
|------|-----------|-----------------|
| Thermostat | `sdm.devices.types.THERMOSTAT` | `nest_thermostat` |
| Camera | `sdm.devices.types.CAMERA` | `nest_camera` |
| Doorbell | `sdm.devices.types.DOORBELL` | `nest_camera` |
| Display (Nest Hub) | `sdm.devices.types.DISPLAY` | `nest_camera` |

## API Endpoints

Base URL: `https://smartdevicemanagement.googleapis.com/v1`

All requests require: `Authorization: Bearer {access_token}`

### List Devices

```
GET /v1/enterprises/{project_id}/devices
```

### Execute Command

```
POST /v1/enterprises/{project_id}/devices/{device_id}:executeCommand

{
  "command": "sdm.devices.commands.ThermostatMode.SetMode",
  "params": { "mode": "HEAT" }
}
```

## Thermostat Commands

| Command | Params |
|---------|--------|
| `ThermostatMode.SetMode` | `{"mode": "HEAT\|COOL\|HEATCOOL\|OFF"}` |
| `ThermostatTemperatureSetpoint.SetHeat` | `{"heatCelsius": 22.0}` |
| `ThermostatTemperatureSetpoint.SetCool` | `{"coolCelsius": 24.0}` |
| `ThermostatTemperatureSetpoint.SetRange` | `{"heatCelsius": 20.0, "coolCelsius": 24.0}` |

## Camera Streaming

### Architecture

Nest cameras use **WebRTC** (Google Home app cameras) or **RTSP** (legacy Nest app cameras). Direct browser-to-Google WebRTC has SDP compatibility issues. Haus uses **go2rtc as a middleware**:

```
Browser → WebRTC → Haus proxy → go2rtc → Google SDM API → Nest Camera
```

go2rtc handles the Nest-specific WebRTC negotiation, token management, and stream extension. The browser just does standard WebRTC `recvonly`.

### go2rtc Configuration

go2rtc streams use the `nest:` source format:

```yaml
streams:
  living_room: "nest:?client_id={id}&client_secret={secret}&refresh_token={token}&project_id={project}&device_id={device_id}"
```

The `device_id` is the last segment of the SDM device name (after `devices/`).

### Haus Camera Proxy Endpoints

```
GET  /api/cameras                      — List go2rtc streams
POST /api/cameras/{id}/webrtc          — WebRTC SDP exchange (proxy to go2rtc)
GET  /api/cameras/{id}/stream          — MP4/MSE stream (proxy to go2rtc)
```

### WebRTC Flow (Browser)

```javascript
pc = new RTCPeerConnection({ iceServers: [] })
pc.addTransceiver('video', { direction: 'recvonly' })
pc.addTransceiver('audio', { direction: 'recvonly' })
pc.ontrack = (e) => { video.srcObject = e.streams[0] }

offer = await pc.createOffer()
await pc.setLocalDescription(offer)

answer = await fetch(`/api/cameras/${streamId}/webrtc`, {
  method: 'POST',
  body: JSON.stringify({ type: 'offer', sdp: offer.sdp })
}).then(r => r.json())

await pc.setRemoteDescription(new RTCSessionDescription(answer))
```

### Direct SDM Camera Commands (without go2rtc)

#### Generate RTSP Stream

```json
{
  "command": "sdm.devices.commands.CameraLiveStream.GenerateRtspStream",
  "params": {}
}
```

Response includes `streamUrls.rtspUrl`, `streamExtensionToken`, `expiresAt`. Stream expires after 5 minutes.

#### Generate WebRTC Stream

```json
{
  "command": "sdm.devices.commands.CameraLiveStream.GenerateWebRtcStream",
  "params": { "offerSdp": "v=0\r\n..." }
}
```

Response includes `answerSdp`, `mediaSessionId`, `expiresAt`. SDP must use Unified format, support Trickle ICE, use Opus for audio. Answer must be used within 30 seconds.

#### Extend / Stop Stream

```json
{"command": "CameraLiveStream.ExtendRtspStream", "params": {"streamExtensionToken": "..."}}
{"command": "CameraLiveStream.StopRtspStream", "params": {"streamExtensionToken": "..."}}
```

## Device Discovery & Enrichment

Nest devices are discovered on the local network by their MAC addresses but cannot be probed locally. Haus enriches them with names from the SDM API:

1. **Network scan** finds devices with Google MAC prefixes → creates entries like "Google .89"
2. **On startup**, Haus queries the SDM API for device list (thermostats, cameras, displays)
3. **Enrichment** matches unnamed Google devices to Nest SDM devices and updates names (e.g., "Google .89" → "Living Room Camera")
4. **Subsequent scans** preserve enriched names (DB `UpsertDevice` won't overwrite good names with generic ones)

### MAC Address Prefixes (OUI)

| Prefix | Vendor |
|--------|--------|
| `7c:10:15` | Google/Nest |
| `b0:09:da` | Google |
| `18:7f:88` | Google |
| `18:b4:30` | Google (Nest Labs) |
| `f8:0f:f9` | Google |
| `48:d6:d5` | Google |

## Key Traits Reference

| Trait | Fields | Device Types |
|-------|--------|-------------|
| `sdm.devices.traits.Info` | `customName` | All |
| `sdm.devices.traits.Temperature` | `ambientTemperatureCelsius` | Thermostat, Display |
| `sdm.devices.traits.Humidity` | `ambientHumidityPercent` | Thermostat, Display |
| `sdm.devices.traits.ThermostatMode` | `mode`, `availableModes` | Thermostat |
| `sdm.devices.traits.ThermostatTemperatureSetpoint` | `heatCelsius`, `coolCelsius` | Thermostat |
| `sdm.devices.traits.ThermostatHvac` | `status` (HEATING/COOLING/OFF) | Thermostat |
| `sdm.devices.traits.CameraLiveStream` | `maxVideoResolution`, `videoCodecs`, `audioCodecs` | Camera, Doorbell, Display |
| `sdm.devices.traits.CameraImage` | `maxImageResolution` | Camera, Doorbell |
| `sdm.devices.traits.DoorbellChime` | (event-only) | Doorbell |

## Haus Integration

- **Discovery:** Google MAC prefix detection during network scan
- **Auth:** `GET /api/google/auth` → OAuth redirect → `GET /api/google/callback`
- **Status:** `GET /api/google/status` → `{"connected": true}`
- **Devices:** `GET /api/google/devices` → list all Nest devices with traits
- **Camera streaming:** go2rtc WebRTC proxy via `/api/cameras/{id}/webrtc`
- **Thermostat control:** Via AI chat (tool use) or direct SDM command execution
- **Auto-enrichment:** Device names updated from SDM API on startup
- **Performance:** Local probe skipped for `nest_camera` devices (cloud-only, no local ports)

## AI Chat Capabilities

When chatting with a Nest camera, the AI can:
- **See through the camera** — captures a live snapshot via go2rtc and uses Claude's vision AI to describe the scene in detail (people, objects, lighting, activity)
- **Report connection status** — confirms streaming is active
- **Describe capabilities** — live streaming, motion detection, person detection

When chatting with a Nest thermostat, the AI can:
- **Query temperature and humidity** — real-time readings from the SDM API
- **Report thermostat mode** — HEAT, COOL, HEATCOOL, OFF
- **Set temperature and mode** — via SDM command execution

When chatting with a Nest display (Hub Max), the AI can:
- **See through the camera** — same vision capabilities as cameras
- **Stream live video** — auto-starts WebRTC stream on page load

The AI speaks as the device: "I can see your living room — there's a couch, coffee table, and the TV is on."

## Limitations

- Nest cameras are cloud-only — no local streaming without go2rtc + Google API
- Camera streams expire after 5 minutes (go2rtc handles auto-extension)
- No local API for any Nest device — all communication through Google's servers
- Google Commercial Development program (unlimited users) is currently paused
- Sandbox tier limited to 25 users across 5 structures
- Two-way audio not available through SDM API

## Documentation

- Official: https://developers.google.com/nest/device-access
- Device Access Console: https://console.nest.google.com/device-access
- go2rtc: https://github.com/AlexxIT/go2rtc
