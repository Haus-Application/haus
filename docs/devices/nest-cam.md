---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "nest-cam"
name: "Google Nest Cam (Battery/Wired)"
manufacturer: "Google LLC"
brand: "Google Nest"
model: "Nest Cam (Battery)"
model_aliases: ["Nest Cam (Wired)", "Nest Cam Indoor", "Nest Cam Outdoor", "GA01317-US", "GA01998-US", "GJQ9T", "G6L1Y"]
device_type: "nest_camera"
category: "security"
product_line: "Nest"
release_year: 2021
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["7C:10:15", "B0:09:DA", "18:B4:30", "18:7F:88", "F8:0F:F9", "48:D6:D5"]
  mdns_services: ["_googlerpc._tcp", "_googlecast._tcp"]
  mdns_txt_keys: ["md", "fn"]
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^Nest-Cam.*", "^Google-Nest.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "nest"
  polling_interval_sec: 0
  websocket_event: "nest:state"
  setup_type: "oauth2"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities: ["camera_stream", "camera_snapshot", "motion"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Google Smart Device Management (SDM) API. OAuth2 via nestservices.google.com partner connection flow. Requires Device Access project ($5 one-time fee), Google Cloud project with SDM API enabled, and user's Google account linked to Nest devices. Access tokens expire after 1 hour; refresh tokens are long-lived."
  base_url_template: "https://smartdevicemanagement.googleapis.com/v1/enterprises/{project_id}"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "battery"
  mounting: "wall"
  indoor_outdoor: "both"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://store.google.com/category/cameras_doorbells"
  api_docs: "https://developers.google.com/nest/device-access"
  developer_portal: "https://console.nest.google.com/device-access"
  support: "https://support.google.com/googlenest"
  community_forum: "https://support.google.com/googlenest/community"
  image_url: ""
  fcc_id: "A4RGE2"

# --- TAGS ---
tags: ["cloud-only", "google", "sdm-api", "oauth2", "webrtc", "go2rtc", "ai-vision", "person-detection", "battery", "wired-option", "google-home"]
---

# Google Nest Cam (Battery/Wired)

## What It Is

> The Google Nest Cam is a WiFi security camera available in battery-powered and wired indoor/outdoor variants. It provides 1080p HDR video with a 130-degree field of view, on-device person/animal/vehicle detection (no subscription required for basic alerts), two-way audio, and 3 hours of free cloud event recording. The camera is entirely cloud-dependent -- there is no local API or local streaming capability. All video access goes through Google's Smart Device Management (SDM) API. Haus integrates via the SDM API with go2rtc as a WebRTC middleware layer for live streaming.

## How Haus Discovers It

1. **OUI match** -- Google/Nest devices use MAC prefixes `7C:10:15`, `B0:09:DA`, `18:B4:30`, `18:7F:88`, `F8:0F:F9`, and `48:D6:D5`. These are detected during network scanning but no local ports are open
2. **mDNS** -- Some Nest devices advertise `_googlerpc._tcp` or `_googlecast._tcp` services
3. **No port probe** -- Nest cameras have no open local ports; they communicate exclusively with Google's cloud
4. **Cloud enrichment** -- After OAuth authentication, Haus queries the SDM API for device list and matches unnamed Google MAC-prefix devices to Nest camera names (e.g., "Google .89" becomes "Living Room Camera")

## Pairing / Authentication

> Google Nest cameras use the SDM OAuth2 flow:
>
> 1. **Prerequisites:** Google Cloud project with SDM API enabled, Device Access project ($5 one-time fee) from console.nest.google.com/device-access, OAuth2 Web Application client credentials
>
> 2. **Authorization redirect:**
> ```
> GET https://nestservices.google.com/partnerconnections/{project_id}/auth
>   ?redirect_uri={redirect_uri}
>   &access_type=offline
>   &prompt=consent
>   &client_id={client_id}
>   &response_type=code
>   &scope=https://www.googleapis.com/auth/sdm.service
> ```
> **Critical:** URL is `nestservices.google.com`, NOT `accounts.google.com`. `access_type=offline` required for refresh tokens. `prompt=consent` ensures refresh token is returned.
>
> 3. **Token exchange:**
> ```
> POST https://www.googleapis.com/oauth2/v4/token
> Content-Type: application/x-www-form-urlencoded
>
> client_id={id}&client_secret={secret}&code={code}&grant_type=authorization_code&redirect_uri={uri}
> ```
>
> 4. **Token refresh:** Access tokens expire after 1 hour. Refresh with `grant_type=refresh_token`. Google's refresh response does NOT always return a new refresh token -- preserve the original.
>
> 5. **Haus endpoints:**
>    - `GET /api/google/auth` -- OAuth redirect
>    - `GET /api/google/callback` -- OAuth callback
>    - `GET /api/google/status` -- Connection status

## API Reference

> ### SDM API Base
>
> ```
> Base URL: https://smartdevicemanagement.googleapis.com/v1
> Auth: Authorization: Bearer {access_token}
> ```
>
> ### List Devices
>
> ```
> GET /v1/enterprises/{project_id}/devices
> ```
>
> Camera devices have type `sdm.devices.types.CAMERA`. Doorbells are `sdm.devices.types.DOORBELL`. Nest Hub displays with cameras are `sdm.devices.types.DISPLAY`.
>
> ### Camera Traits
>
> | Trait | Description |
> |-------|-------------|
> | `sdm.devices.traits.Info` | `customName` -- user-assigned device name |
> | `sdm.devices.traits.CameraLiveStream` | `maxVideoResolution`, `videoCodecs`, `audioCodecs` -- streaming capabilities |
> | `sdm.devices.traits.CameraImage` | `maxImageResolution` -- snapshot capabilities |
> | `sdm.devices.traits.DoorbellChime` | Doorbell press events (doorbell models only) |
>
> ### Generate WebRTC Stream
>
> ```json
> POST /v1/enterprises/{project_id}/devices/{device_id}:executeCommand
>
> {
>   "command": "sdm.devices.commands.CameraLiveStream.GenerateWebRtcStream",
>   "params": { "offerSdp": "v=0\r\n..." }
> }
> ```
>
> Response includes `answerSdp`, `mediaSessionId`, `expiresAt`. SDP must use Unified format with Trickle ICE and Opus audio. Answer must be used within 30 seconds.
>
> ### Generate RTSP Stream
>
> ```json
> POST /v1/enterprises/{project_id}/devices/{device_id}:executeCommand
>
> {
>   "command": "sdm.devices.commands.CameraLiveStream.GenerateRtspStream",
>   "params": {}
> }
> ```
>
> Response includes `streamUrls.rtspUrl`, `streamExtensionToken`, `expiresAt`. Stream expires after 5 minutes.
>
> ### Extend Stream
>
> ```json
> {
>   "command": "sdm.devices.commands.CameraLiveStream.ExtendRtspStream",
>   "params": { "streamExtensionToken": "..." }
> }
> ```
>
> ### Stop Stream
>
> ```json
> {
>   "command": "sdm.devices.commands.CameraLiveStream.StopRtspStream",
>   "params": { "streamExtensionToken": "..." }
> }
> ```
>
> ### go2rtc Streaming Architecture
>
> Nest cameras use WebRTC (Google Home cameras) or RTSP (legacy Nest app cameras). Direct browser-to-Google WebRTC has SDP compatibility issues. Haus uses go2rtc as middleware:
>
> ```
> Browser -> WebRTC -> Haus proxy -> go2rtc -> Google SDM API -> Nest Camera
> ```
>
> **go2rtc stream configuration:**
> ```yaml
> streams:
>   living_room: "nest:?client_id={id}&client_secret={secret}&refresh_token={token}&project_id={project}&device_id={device_id}"
> ```
>
> **Haus camera proxy endpoints:**
> ```
> GET  /api/cameras                     -- List go2rtc streams
> POST /api/cameras/{id}/webrtc         -- WebRTC SDP exchange
> GET  /api/cameras/{id}/stream         -- MP4/MSE stream
> ```
>
> **Browser WebRTC flow:**
> ```javascript
> pc = new RTCPeerConnection({ iceServers: [] })
> pc.addTransceiver('video', { direction: 'recvonly' })
> pc.addTransceiver('audio', { direction: 'recvonly' })
> pc.ontrack = (e) => { video.srcObject = e.streams[0] }
>
> offer = await pc.createOffer()
> await pc.setLocalDescription(offer)
>
> answer = await fetch(`/api/cameras/${streamId}/webrtc`, {
>   method: 'POST',
>   body: JSON.stringify({ type: 'offer', sdp: offer.sdp })
> }).then(r => r.json())
>
> await pc.setRemoteDescription(new RTCSessionDescription(answer))
> ```

## AI Capabilities

> When chatting with a Nest camera, the AI can:
> - **See through the camera** -- captures a live snapshot via go2rtc and uses Claude's vision AI to describe the scene in detail (people, objects, lighting, activity)
> - **Report connection status** -- confirms streaming is active
> - **Describe capabilities** -- live streaming, motion detection, person detection
>
> When chatting with a Nest display (Hub Max), the AI can:
> - **See through the camera** -- same vision capabilities as cameras
> - **Stream live video** -- auto-starts WebRTC stream on page load
>
> The AI speaks as the device: "I can see your living room -- there's a couch, coffee table, and the TV is on."

## Quirks & Notes

- **Cloud-only** -- No local API, no local streaming. All communication through Google's SDM API servers
- **Stream expiration** -- Camera streams expire after 5 minutes; go2rtc handles automatic stream extension
- **$5 developer fee** -- Device Access registration costs $5 per Google account (one-time)
- **Commercial program paused** -- Google is NOT currently accepting new Commercial Development applications. Sandbox tier limited to 25 users across 5 structures
- **Refresh token preservation** -- Google's refresh response does NOT always return a new refresh token. The original must be preserved
- **go2rtc dependency** -- Live streaming requires go2rtc running as a sidecar service for WebRTC/RTSP negotiation
- **Two-way audio** -- Not available through the SDM API
- **Battery version limitations** -- Battery model may enter deep sleep to conserve power; streaming requires waking the camera first
- **Person/animal/vehicle detection** -- On-device AI detection with 3 hours of free event recording. Nest Aware subscription ($8/mo) for 30-day history
- **Local probe skip** -- Haus skips local port probing for `nest_camera` devices since they have no local ports (performance optimization)

## Similar Devices

> - [Ring Indoor Cam](ring-indoor-cam.md) -- Amazon's cloud-only camera ecosystem
> - [Arlo Pro 5](arlo-pro-5.md) -- Cloud-primary with optional RTSP via SmartHub
> - [Wyze Cam v3](wyze-cam-v3.md) -- Budget camera with optional local RTSP
> - [Reolink Argus 3 Pro](reolink-argus-3-pro.md) -- Local-first camera (opposite philosophy)
> - [UniFi Protect G4 Bullet](unifi-protect-g4-bullet.md) -- Professional local-first camera
