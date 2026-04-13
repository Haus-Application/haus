---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "philips-hue-bridge"
name: "Philips Hue Bridge"
manufacturer: "Signify Netherlands B.V."
brand: "Philips Hue"
model: "BSB002"
model_aliases: ["3241312018", "BSB002-US", "BSB002-EU"]
device_type: "hue_bridge"
category: "smart_home"
product_line: "Hue"
release_year: 2015
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["ethernet", "zigbee", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "00:17:88"        # Philips Lighting BV (primary OUI for Hue bridges)
    - "EC:B5:FA"        # Signify Netherlands B.V. (newer production runs)
  mdns_services:
    - "_hue._tcp"
  mdns_txt_keys:
    - "bridgeid"        # uppercase hex bridge ID (MAC without colons, padded)
    - "modelid"         # "BSB002"
    - "swversion"       # firmware version
  default_ports: [80, 443]
  signature_ports: [443]
  ssdp_search_target: "urn:schemas-upnp-org:device:basic:1"
  ssdp_server_string: "Linux/3.14.0 UPnP/1.0 IpBridge/1.x.x"
  hostname_patterns:
    - "^Philips-hue"
    - "^ecb5fa[0-9a-f]{6}$"
    - "^001788[0-9a-f]{6}$"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/api/0/config"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "nginx"
    body_contains: "\"modelid\":\"BSB002\""
    headers: {}
  - port: 443
    path: "/api/0/config"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "nginx"
    body_contains: "\"apiversion\""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "hue"
  polling_interval_sec: 5
  websocket_event: "hue:state"
  setup_type: "link_button"
  ai_chattable: true
  haus_milestone: "M3"

# --- CAPABILITIES ---
capabilities:
  - "on_off"
  - "brightness"
  - "color"
  - "color_temp"
  - "scenes"
  - "groups"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "link_button"
  auth_detail: "Press physical link button, POST /api with devicetype, receive API key. Use hue-application-key header for all subsequent requests."
  base_url_template: "https://{ip}/clip/v2"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.philips-hue.com/en-us/p/hue-bridge/046677458478"
  api_docs: "https://developers.meethue.com/develop/hue-api-v2/"
  developer_portal: "https://developers.meethue.com/"
  support: "https://www.philips-hue.com/en-us/support"
  community_forum: "https://developers.meethue.com/forum/"
  image_url: "https://www.philips-hue.com/en-us/p/hue-bridge/046677458478#checks-checks"
  fcc_id: "2ABA6-BSB002"

# --- TAGS ---
tags: ["zigbee_coordinator", "local_api", "clip_v2", "entertainment_api", "bridge", "signify"]
---

# Philips Hue Bridge

## What It Is

The Philips Hue Bridge (model BSB002) is the central hub for the Philips Hue ecosystem, manufactured by Signify Netherlands B.V. It connects to your home network via Ethernet and coordinates up to 50 Zigbee Light Link (ZLL) devices and up to 63 accessories (sensors, switches). The bridge exposes a local HTTPS REST API (CLIP v2) that allows full control of all paired devices without cloud dependency. It also supports Bluetooth for initial mobile setup and a cloud connection for remote access via the Hue app, but operates fully locally once configured. The BSB002 is the second-generation bridge (square form factor) and remains the current production model.

## How Haus Discovers It

1. **OUI Match** -- During network scan, any device with MAC prefix `00:17:88` or `EC:B5:FA` is flagged as a Signify/Philips device.
2. **mDNS Discovery** -- Browse for `_hue._tcp.local.` services. The bridge advertises TXT records including `bridgeid`, `modelid`, and `swversion`. The `bridgeid` is derived from the bridge MAC address (with FFFE inserted in the middle).
3. **SSDP Discovery** -- The bridge responds to UPnP M-SEARCH with search target `urn:schemas-upnp-org:device:basic:1`. The response includes the bridge description XML URL.
4. **Cloud Discovery Fallback** -- `GET https://discovery.meethue.com/` returns JSON array of bridges visible on the local network (via N-UPnP / STUN). This is a convenience endpoint and not required for local discovery.
5. **HTTP Fingerprint** -- `GET http://{ip}/api/0/config` returns the unauthenticated config including `modelid: "BSB002"`, `apiversion`, `swversion`, and `bridgeid`. Port 443 with self-signed TLS also works.
6. **Hostname Pattern** -- Bridges typically appear as `Philips-hue` or their MAC-derived hostname in DHCP.

## Pairing / Authentication

The Hue Bridge uses a physical link-button pairing model. No passwords, no cloud accounts required.

### Pairing Flow

1. User presses the physical **link button** on top of the bridge (the large circular button).
2. Within **30 seconds**, Haus sends:
   ```
   POST https://{bridge_ip}/api
   Content-Type: application/json

   {"devicetype": "haus#hub", "generateclientkey": true}
   ```
3. The bridge responds with:
   ```json
   [{"success": {"username": "abcdef1234567890...", "clientkey": "AABBCCDD..."}}]
   ```
   - `username` is the API key (40-character hex string) used for all REST requests.
   - `clientkey` is used for the Entertainment API (streaming/sync features).
4. If the link button was not pressed, the response is:
   ```json
   [{"error": {"type": 101, "address": "", "description": "link button not pressed"}}]
   ```
5. Haus stores the `username` as the integration credential. All subsequent API calls include the header:
   ```
   hue-application-key: {username}
   ```

### Security Notes

- The bridge uses a **self-signed TLS certificate**. Haus must pin or accept the certificate on first connection. The certificate CN contains the bridge ID.
- API keys do not expire but can be revoked via the API or a factory reset.
- The bridge enforces rate limiting: approximately **10 commands per second** for individual lights, **1 command per second** for groups.

## API Reference

All endpoints are on **port 443** over HTTPS with the `hue-application-key` header.

**Base URL:** `https://{bridge_ip}/clip/v2`

### List Lights

```
GET /clip/v2/resource/light
```

**Response:** Array of light objects:
- `id` -- UUID (v4)
- `metadata.name` -- user-assigned light name
- `on.on` -- boolean power state
- `dimming.brightness` -- 0.0-100.0 (float percentage)
- `color.xy` -- CIE 1931 xy color coordinates `{"x": 0.3127, "y": 0.3290}`
- `color_temperature.mirek` -- color temperature in mirek (153-500, where 153 = 6500K, 500 = 2000K)
- `color_temperature.mirek_valid` -- boolean, false if light is in color mode
- `owner.rid` -- reference to the device resource
- `type` -- "light"

### Control a Light

```
PUT /clip/v2/resource/light/{id}
```

**Body (all fields optional):**
```json
{
  "on": {"on": true},
  "dimming": {"brightness": 75.0},
  "color": {"xy": {"x": 0.4578, "y": 0.4101}},
  "color_temperature": {"mirek": 250},
  "dynamics": {"duration": 400}
}
```

**Notes:**
- `dynamics.duration` is transition time in milliseconds.
- Setting `color.xy` overrides `color_temperature` and vice versa.
- Brightness is independent of on/off state -- you can set brightness without turning on.

### List Rooms

```
GET /clip/v2/resource/room
```

**Response:** Array of room objects:
- `id` -- UUID
- `metadata.name` -- room name
- `children` -- array of `{"rid": "...", "rtype": "device"}` references
- `services` -- array including a `grouped_light` reference for room-level control

### Control a Room (Grouped Light)

```
PUT /clip/v2/resource/grouped_light/{grouped_light_id}
```

Same body format as individual light control. Affects all lights in the room simultaneously.

### List Scenes

```
GET /clip/v2/resource/scene
```

**Response:** Array of scene objects:
- `id` -- UUID
- `metadata.name` -- scene name (e.g., "Energize", "Relax", "Concentrate")
- `group.rid` -- room this scene belongs to
- `group.rtype` -- "room" or "zone"
- `actions` -- array of light states the scene applies

### Activate a Scene

```
PUT /clip/v2/resource/scene/{id}
```

```json
{"recall": {"action": "active"}}
```

### List Devices

```
GET /clip/v2/resource/device
```

Returns all paired Zigbee devices with their services (light, motion, temperature, button, etc.).

### Event Stream (SSE)

```
GET /eventstream/clip/v2
```

**Headers:**
```
hue-application-key: {key}
Accept: text/event-stream
```

Server-Sent Events stream providing real-time updates for all resource changes. Each event contains the full updated resource object. This is more efficient than polling for detecting state changes.

## Color Reference

CIE 1931 xy coordinates for common colors:

| Color | CIE x | CIE y | Description |
|-------|--------|--------|-------------|
| Warm White | 0.4578 | 0.4101 | ~2700K incandescent |
| Cool White | 0.3127 | 0.3290 | ~4000K neutral (D65 illuminant) |
| Red | 0.6750 | 0.3220 | Deep red |
| Blue | 0.1532 | 0.0475 | Deep blue |
| Green | 0.1700 | 0.7000 | Pure green |
| Purple | 0.2703 | 0.1398 | Violet-purple |
| Orange | 0.5614 | 0.3944 | Warm orange |
| Pink | 0.3944 | 0.1990 | Soft pink |

**Mirek-to-Kelvin conversion:** `K = 1,000,000 / mirek`
- 153 mirek = 6536K (coolest)
- 250 mirek = 4000K (neutral)
- 370 mirek = 2703K (warm)
- 500 mirek = 2000K (warmest, candlelight)

## AI Capabilities

When the AI concierge "chats as" a Hue Bridge, it can:

- **List all lights** with current on/off state, brightness, color, and room assignment
- **Toggle individual lights** by name with fuzzy matching ("turn on the kitchen light")
- **Set brightness** on any light or room (0-100%)
- **Change colors** using named colors (warm, cool, red, blue, green, purple, orange, pink) mapped to CIE xy coordinates
- **Set color temperature** by mirek value or descriptive terms (warm, cool, daylight)
- **List and activate scenes** by name ("set the living room to Relax")
- **Control entire rooms** -- turn all lights on/off, set room brightness, apply scenes
- **Report status** -- "I have 12 lights across 4 rooms. Living room lights are on at 75%."

The AI speaks in first person as the bridge, providing a natural conversational interface to the lighting system.

## Quirks & Notes

- **Self-signed TLS:** The bridge generates its own CA and leaf certificate. The leaf certificate CN is the bridge ID. Haus must skip TLS verification or pin the certificate on first pairing.
- **Rate Limits:** ~10 commands/sec for individual lights, ~1 command/sec for groups. Exceeding this causes commands to be queued or dropped silently.
- **Mirek vs Color:** Setting `color.xy` puts the light in "color mode" and `color_temperature.mirek_valid` becomes false. Setting `color_temperature.mirek` puts it back in white mode.
- **Bridge ID Format:** The bridge ID is the MAC address with `FFFE` inserted in the middle. For MAC `00:17:88:AB:CD:EF`, the bridge ID is `001788FFFEABCDEF`.
- **Legacy CLIP v1:** The bridge still supports the old `/api/{key}/lights` style API on port 80. Haus uses CLIP v2 exclusively but v1 is available as fallback.
- **Max Devices:** 50 lights, 63 accessories (sensors, switches, remotes), 12 Entertainment areas.
- **Zigbee Channel:** The bridge typically operates on Zigbee channel 11, 15, 20, or 25. This can conflict with Wi-Fi channels.
- **Firmware Updates:** The bridge auto-updates firmware by default via Signify's cloud. This can occasionally change API behavior.
- **EventStream:** The SSE event stream on `/eventstream/clip/v2` is the preferred method for real-time updates. It requires the same `hue-application-key` header. Events are JSON arrays containing `creationtime`, `data`, `id`, and `type` fields.

## Similar Devices

- **philips-hue-bulb-a19** -- The most common bulb controlled by this bridge
- **philips-hue-lightstrip-plus** -- LED strip controlled by this bridge
- **philips-hue-sync-box** -- HDMI sync device with its own API (separate from bridge)
