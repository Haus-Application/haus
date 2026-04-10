# Google Nest SDM API Integration

Haus integrates with Google Nest devices via the Smart Device Management (SDM)
API. This provides local-quality control of Nest thermostats, cameras, doorbells,
and displays through Google's cloud API.

## Prerequisites

1. A Google Cloud project with the SDM API enabled
2. A Device Access project ($5 one-time registration fee)
3. OAuth2 client credentials (client ID + client secret)
4. The user must have linked their Nest account to Google

## OAuth2 Flow

The SDM API uses Google OAuth2 with a Nest-specific partner connection flow.

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

Note: The authorization URL uses `nestservices.google.com`, NOT
`accounts.google.com`. This is the Device Access partner connection flow.

The `access_type=offline` parameter is required to receive a refresh token.
The `prompt=consent` parameter forces the consent screen to ensure a refresh
token is returned even if the user has previously authorized.

### Step 2: Exchange Authorization Code

```
POST https://www.googleapis.com/oauth2/v4/token
Content-Type: application/x-www-form-urlencoded

client_id={client_id}
&client_secret={client_secret}
&code={authorization_code}
&grant_type=authorization_code
&redirect_uri={redirect_uri}
```

Response:
```json
{
  "access_token": "ya29.a0...",
  "refresh_token": "1//0d...",
  "expires_in": 3600,
  "token_type": "Bearer",
  "scope": "https://www.googleapis.com/auth/sdm.service"
}
```

### Step 3: Refresh Access Token

Access tokens expire after 3600 seconds (1 hour). Use the refresh token to
obtain new access tokens without user interaction.

```
POST https://www.googleapis.com/oauth2/v4/token
Content-Type: application/x-www-form-urlencoded

client_id={client_id}
&client_secret={client_secret}
&refresh_token={refresh_token}
&grant_type=refresh_token
```

The refresh token does not expire unless the user revokes access.

## Device Types

| Type | SDM Type String |
|------|----------------|
| Thermostat | `sdm.devices.types.THERMOSTAT` |
| Camera | `sdm.devices.types.CAMERA` |
| Doorbell | `sdm.devices.types.DOORBELL` |
| Display | `sdm.devices.types.DISPLAY` |

## API Endpoints

Base URL: `https://smartdevicemanagement.googleapis.com/v1`

All requests require: `Authorization: Bearer {access_token}`

### List Devices

```
GET /v1/enterprises/{project_id}/devices
```

Response:
```json
{
  "devices": [
    {
      "name": "enterprises/{project_id}/devices/{device_id}",
      "type": "sdm.devices.types.THERMOSTAT",
      "traits": {
        "sdm.devices.traits.Info": { "customName": "Living Room" },
        "sdm.devices.traits.Temperature": { "ambientTemperatureCelsius": 22.5 },
        "sdm.devices.traits.Humidity": { "ambientHumidityPercent": 45.0 },
        "sdm.devices.traits.ThermostatMode": { "mode": "HEAT", "availableModes": ["HEAT", "COOL", "HEATCOOL", "OFF"] },
        "sdm.devices.traits.ThermostatTemperatureSetpoint": { "heatCelsius": 21.0 }
      },
      "parentRelations": [
        { "parent": "enterprises/{project_id}/structures/{structure_id}/rooms/{room_id}", "displayName": "Living Room" }
      ]
    }
  ]
}
```

### Get Single Device

```
GET /v1/enterprises/{project_id}/devices/{device_id}
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

### Set Mode

```json
{
  "command": "sdm.devices.commands.ThermostatMode.SetMode",
  "params": { "mode": "HEAT" }
}
```

Valid modes: `HEAT`, `COOL`, `HEATCOOL`, `OFF`

### Set Heat Temperature

```json
{
  "command": "sdm.devices.commands.ThermostatTemperatureSetpoint.SetHeat",
  "params": { "heatCelsius": 22.0 }
}
```

### Set Cool Temperature

```json
{
  "command": "sdm.devices.commands.ThermostatTemperatureSetpoint.SetCool",
  "params": { "coolCelsius": 24.0 }
}
```

### Set Temperature Range (HEATCOOL mode)

```json
{
  "command": "sdm.devices.commands.ThermostatTemperatureSetpoint.SetRange",
  "params": { "heatCelsius": 20.0, "coolCelsius": 24.0 }
}
```

## Camera Streaming

### Generate RTSP Stream

```json
{
  "command": "sdm.devices.commands.CameraLiveStream.GenerateRtspStream",
  "params": {}
}
```

Response:
```json
{
  "results": {
    "streamUrls": { "rtspUrl": "rtsps://..." },
    "streamToken": "...",
    "streamExtensionToken": "...",
    "expiresAt": "2024-01-01T00:05:00Z"
  }
}
```

The stream URL expires after 5 minutes. Use the extension token to renew.

### Extend Stream

```json
{
  "command": "sdm.devices.commands.CameraLiveStream.ExtendRtspStream",
  "params": { "streamExtensionToken": "..." }
}
```

### Stop Stream

```json
{
  "command": "sdm.devices.commands.CameraLiveStream.StopRtspStream",
  "params": { "streamExtensionToken": "..." }
}
```

## Key Traits Reference

| Trait | Fields | Device Types |
|-------|--------|-------------|
| `sdm.devices.traits.Info` | `customName` | All |
| `sdm.devices.traits.Temperature` | `ambientTemperatureCelsius` | Thermostat, Display |
| `sdm.devices.traits.Humidity` | `ambientHumidityPercent` | Thermostat, Display |
| `sdm.devices.traits.ThermostatMode` | `mode`, `availableModes` | Thermostat |
| `sdm.devices.traits.ThermostatTemperatureSetpoint` | `heatCelsius`, `coolCelsius` | Thermostat |
| `sdm.devices.traits.ThermostatHvac` | `status` (HEATING, COOLING, OFF) | Thermostat |
| `sdm.devices.traits.CameraLiveStream` | `maxVideoResolution`, `videoCodecs`, `audioCodecs` | Camera, Doorbell, Display |
| `sdm.devices.traits.CameraImage` | `maxImageResolution` | Camera, Doorbell |
| `sdm.devices.traits.CameraEventImage` | `maxImageResolution` | Camera, Doorbell |
| `sdm.devices.traits.DoorbellChime` | (event-only trait) | Doorbell |
