---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "eufy-video-doorbell-dual"
name: "Eufy Video Doorbell Dual"
manufacturer: "Anker Innovations (Eufy)"
brand: "Eufy Security"
model: "E8213"
model_aliases: ["E8213181", "T8213", "T82131W1"]
device_type: "doorbell_camera"
category: "security"
product_line: "Eufy Security"
release_year: 2022
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "98:8E:79"        # Eufy / Anker Innovations
    - "7C:C2:94"        # Eufy / Anker Innovations
    - "70:2C:09"        # Anker Innovations
  mdns_services: []     # Eufy devices do not consistently advertise mDNS
  mdns_txt_keys: []
  default_ports: []     # No standard open ports; RTSP available if enabled
  signature_ports: [8554]  # RTSP port when enabled in settings
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^eufy"
    - "^T8213"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No HTTP-based discovery

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "eufy"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "camera_stream"
  - "doorbell"
  - "motion"

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "Eufy uses a proprietary cloud API for device management. Authentication via POST to Eufy's API endpoint with email/password returns a token. P2P protocol used for direct local streaming when on same network. RTSP available as secondary option if enabled in eufy app settings."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "camera"
  power_source: "hardwired"
  mounting: "door"
  indoor_outdoor: "outdoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.eufy.com/products/video-doorbells"
  api_docs: ""
  developer_portal: ""
  support: "https://www.eufy.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AOKB-T8213"

# --- TAGS ---
tags: ["local_storage", "dual_camera", "homebase", "anker", "doorbell", "hardwired", "p2p_streaming", "no_subscription"]
---

# Eufy Video Doorbell Dual

## What It Is

The Eufy Video Doorbell Dual is a hardwired video doorbell with a unique dual-camera design manufactured by Anker Innovations under the Eufy Security brand. It features two separate camera lenses: a main upper camera for visitor face-level video (2K resolution, 160-degree FOV) and a secondary lower "package camera" (1600x1200, 120-degree down-facing FOV) that monitors packages on the doorstep. The dual-camera approach solves the common doorbell blind spot problem where packages left directly below the camera are invisible. It includes two-way audio, on-device AI for human and package detection, and a built-in 4.3-inch display for delivery instructions. Video storage is local -- either on the paired HomeBase 2 or via the eufy app's local P2P streaming -- with no mandatory cloud subscription. This makes it one of the most privacy-conscious video doorbells available.

## How Haus Discovers It

1. **OUI Match** -- During network scan, devices with MAC prefixes `98:8E:79`, `7C:C2:94`, or `70:2C:09` are flagged as Eufy/Anker devices.
2. **Hostname Pattern** -- Eufy devices may appear with hostnames starting with `eufy` or their model number in DHCP tables.
3. **HomeBase Discovery** -- The Eufy HomeBase 2 is the primary network device (the doorbell communicates to the HomeBase via a proprietary Wi-Fi link). The HomeBase is what appears on the LAN. Some community research has found HomeBase devices responding on port 8554 (RTSP) when enabled.
4. **RTSP Probe** -- If RTSP is enabled in the eufy app, the doorbell stream may be accessible at `rtsp://{homebase_ip}:8554/{channel}`. This requires manual user configuration in the eufy Security app.

## Pairing / Authentication

Eufy does not provide an official API. The device ecosystem uses a combination of cloud authentication and peer-to-peer (P2P) local streaming.

### Setup Flow

1. **Eufy App Required:** Initial setup is exclusively through the eufy Security mobile app. The doorbell is paired to a HomeBase 2 (or HomeBase 3) station via a proprietary wireless protocol.
2. **Account Creation:** A eufy/Anker account is required for initial setup and remote access.
3. **HomeBase Pairing:** The HomeBase connects to the LAN via Ethernet. The doorbell communicates to the HomeBase wirelessly (not directly to Wi-Fi router in most configurations).
4. **RTSP Enablement (Optional):** In the eufy Security app, navigate to the doorbell settings and enable "RTSP Stream" under the camera settings. This exposes the video feed on the HomeBase's IP address at port 8554.

### Third-Party Integration Methods

1. **RTSP (Best Option):** When manually enabled, provides a standard RTSP stream accessible locally. Stream URL format: `rtsp://{homebase_ip}:8554/{stream_id}`. This is the most reliable integration path but must be enabled by the user.
2. **P2P Protocol (Reverse-Engineered):** Community projects (notably `eufy-security-client` on GitHub) have reverse-engineered the P2P protocol that eufy uses for local streaming. This protocol uses UDP for peer discovery and establishes a direct encrypted video stream.
3. **Cloud API (Reverse-Engineered):** The eufy cloud API has been partially reverse-engineered. Authentication is via email/password to an Anker API endpoint, returning session tokens. This API provides device listing, push notification history, and some device control.

### Security Notes

- Eufy faced a major privacy controversy in late 2022 when security researchers discovered that camera thumbnails were being uploaded to AWS cloud servers despite marketing claims of "local only" storage. Eufy subsequently improved their practices and added end-to-end encryption.
- The P2P protocol uses encryption but the security of the reverse-engineered implementation varies by community project.
- RTSP streams on the local network are unencrypted.

## API Reference

No official API exists. Integration options:

### RTSP Stream (When Enabled)

```
rtsp://{homebase_ip}:8554/{stream_id}
```

Standard RTSP stream accessible by any RTSP client (VLC, go2rtc, ffmpeg). The `stream_id` is specific to each camera channel. The main camera and package camera may have separate stream IDs.

### Community P2P Protocol

The `eufy-security-client` (Node.js) and `go2rtc` projects implement Eufy's P2P protocol:

1. **Authentication:** Login to Eufy's cloud to get P2P connection details (device serial, P2P addresses, encryption keys).
2. **UDP Discovery:** Send discovery packets to the device's P2P address.
3. **Stream Negotiation:** Establish an encrypted video stream via the P2P connection.
4. **go2rtc Integration:** go2rtc natively supports Eufy cameras via the P2P protocol, providing RTSP/WebRTC re-streaming.

### Push Notifications (Cloud)

Event notifications (doorbell press, motion, person detected, package detected) are delivered via push notification infrastructure. The reverse-engineered cloud API can retrieve notification history.

## AI Capabilities

AI integration is not currently planned due to lack of official API. If implemented via RTSP, the AI concierge could:

- Show live video from both the main and package cameras
- Report doorbell press events (via P2P or cloud push notification interception)
- Display dual-camera views for comprehensive doorstep monitoring

## Quirks & Notes

- **Dual Camera System:** The two cameras operate independently. The main camera captures visitor faces at eye level (2K, 160-degree FOV). The package camera points downward at approximately 45 degrees (1600x1200, 120-degree FOV). In the eufy app, both views can be displayed simultaneously in a split-screen or picture-in-picture mode.
- **HomeBase Required:** Unlike some standalone Eufy cameras, the Doorbell Dual requires a HomeBase 2 or HomeBase 3 station. The HomeBase provides local storage (16GB eMMC built-in, expandable via USB HDD up to 16TB on HomeBase 3), acts as an alarm siren, and manages the P2P connection.
- **Hardwired Only:** This model requires existing doorbell wiring (16-24V AC, 30VA transformer recommended). There is no battery option. The existing mechanical or electronic chime continues to work.
- **Built-in Display:** The 4.3-inch color display on the front of the doorbell can show custom messages, delivery instructions, or animated responses to visitors. This is configured via the eufy app.
- **Privacy Controversy:** The 2022 discovery that Eufy was uploading thumbnails to cloud servers despite "local only" claims damaged trust. Eufy/Anker responded by adding end-to-end encryption options and improving transparency, but the incident is worth noting for privacy-conscious users.
- **No Subscription Required:** All core features (recording, smart detection, two-way talk) work without a subscription. Eufy does offer optional cloud storage plans but the local-first approach is the primary model.
- **RTSP Limitations:** When RTSP is enabled, the stream quality may be reduced compared to the P2P stream. Additionally, enabling RTSP may disable some features in the eufy app (such as two-way audio via RTSP).
- **On-Device AI:** Human detection and package detection run on the device/HomeBase locally, not in the cloud. Detection accuracy is generally good but not as refined as Google's or Ring's cloud-based ML models.

## Similar Devices

- **ring-video-doorbell-4** -- Cloud-dependent alternative with pre-roll video
- **nest-doorbell-battery** -- Google's cloud-based alternative with SDM API support
- **arlo-essential-video-doorbell** -- Arlo's cloud-based alternative with head-to-toe view
