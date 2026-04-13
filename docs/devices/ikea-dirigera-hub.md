---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ikea-dirigera-hub"
name: "IKEA DIRIGERA Hub"
manufacturer: "IKEA of Sweden"
brand: "IKEA Home Smart"
model: "E2112"
model_aliases: ["DIRIGERA Hub for smart products"]
device_type: "dirigera_hub"
category: "smart_home"
product_line: "DIRIGERA"
release_year: 2022
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["zigbee", "thread", "matter", "ethernet", "wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["94:B9:7E", "CC:50:E3", "DC:EF:CA", "60:01:94", "B4:E1:EB", "34:25:BE"]
  mdns_services: ["_ihsp._tcp"]
  mdns_txt_keys: ["version", "name"]
  default_ports: [443, 8443]
  signature_ports: [8443]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^DIRIGERA-[a-f0-9]+$"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8443
    path: "/v1"
    method: "GET"
    expect_status: 401
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "dirigera"
  polling_interval_sec: 10
  websocket_event: "dirigera:state"
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp", "scenes", "groups"]

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 8443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "OAuth2 code challenge flow. POST /v1/oauth/authorize with code_challenge, press button on hub, then POST /v1/oauth/token to exchange code for bearer token."
  base_url_template: "https://{ip}:8443/v1"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee", "thread", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.ikea.com/us/en/p/dirigera-hub-for-smart-products-white-smart-50503409/"
  api_docs: "https://github.com/Leggin/dirigera"
  developer_portal: ""
  support: "https://www.ikea.com/us/en/customer-service/"
  community_forum: "https://github.com/home-assistant/core/tree/dev/homeassistant/components/ikea_home_smart"
  image_url: ""
  fcc_id: "2AHFL-E2112"

# --- TAGS ---
tags: ["zigbee_hub", "thread_border_router", "matter", "ikea", "tradfri_successor", "oauth2"]
---

# IKEA DIRIGERA Hub

## What It Is

The IKEA DIRIGERA Hub (model E2112) is IKEA's second-generation smart home hub, replacing the TRADFRI Gateway. It is a significant upgrade: it supports Zigbee 3.0, Thread (acting as a Thread Border Router), and Matter, while also maintaining backward compatibility with existing TRADFRI Zigbee devices. The hub connects via Ethernet (Gigabit) or WiFi and is managed through the IKEA Home Smart app. At roughly $60, it remains one of the most affordable Matter-capable hubs on the market. The DIRIGERA also exposes a local HTTPS REST API that has been reverse-engineered by the community, making local control possible without cloud dependency.

## How Haus Discovers It

1. **OUI match**: The DIRIGERA's Ethernet or WiFi MAC address will match IKEA OUI prefixes (94:B9:7E, CC:50:E3, DC:EF:CA, 60:01:94, B4:E1:EB, or 34:25:BE).
2. **mDNS**: The DIRIGERA advertises `_ihsp._tcp` (IKEA Home Smart Protocol) via mDNS. The TXT record includes `version` and `name` fields.
3. **Hostname pattern**: DHCP hostname typically matches `DIRIGERA-*` followed by a hex identifier.
4. **Port probe**: HTTPS on port 8443 is the signature port. A GET to `/v1` returns HTTP 401 (unauthorized), confirming the API is present.
5. **TLS certificate inspection**: The self-signed TLS certificate on port 8443 contains IKEA-specific subject fields.

## Pairing / Authentication

The DIRIGERA uses an OAuth2-like code challenge flow for local API authentication:

1. **Generate a code challenge**: Create a random `code_verifier` (43-128 characters), then compute `code_challenge = base64url(sha256(code_verifier))`.
2. **Request authorization**: POST to `https://{ip}:8443/v1/oauth/authorize`:
   ```json
   {
     "challenge": "{code_challenge}",
     "name": "Haus Hub"
   }
   ```
   The response returns a `code` value.
3. **Press the action button**: The user must physically press the round button on top of the DIRIGERA hub within 60 seconds. This authorizes the pending pairing request.
4. **Exchange for token**: POST to `https://{ip}:8443/v1/oauth/token`:
   ```json
   {
     "code": "{code_from_step_2}",
     "name": "Haus Hub",
     "grant_type": "authorization_code",
     "code_verifier": "{code_verifier}"
   }
   ```
5. **Receive bearer token**: The response contains an `access_token` (a long-lived bearer token). Store it securely.
6. **All subsequent requests**: Include `Authorization: Bearer {access_token}` header.

## API Reference

The DIRIGERA exposes a RESTful HTTPS API on port 8443. All responses are JSON. The TLS certificate is self-signed, so certificate verification must be disabled or the cert pinned.

### Key Endpoints

| Path | Method | Description |
|------|--------|-------------|
| `/v1/devices` | GET | List all paired devices |
| `/v1/devices/{id}` | GET | Get specific device state |
| `/v1/devices/{id}` | PATCH | Update device attributes |
| `/v1/rooms` | GET | List all rooms |
| `/v1/rooms/{id}` | GET | Get room details with devices |
| `/v1/scenes` | GET | List all scenes |
| `/v1/scenes/{id}` | POST | Trigger a scene |
| `/v1/home` | GET | Get home overview |
| `/v1/hub/status` | GET | Hub status and firmware version |
| `/v1/music` | GET | Sonos integration status (if linked) |
| `/v1/oauth/authorize` | POST | Begin OAuth flow |
| `/v1/oauth/token` | POST | Exchange code for bearer token |

### Device Object Structure

```json
{
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "type": "light",
  "deviceType": "light",
  "createdAt": "2023-01-01T00:00:00.000Z",
  "isReachable": true,
  "lastSeen": "2024-01-01T12:00:00.000Z",
  "attributes": {
    "customName": "Living Room Lamp",
    "model": "TRADFRI bulb E27 CWS 806lm",
    "manufacturer": "IKEA of Sweden",
    "firmwareVersion": "1.0.012",
    "isOn": true,
    "lightLevel": 75,
    "colorTemperature": 2700,
    "colorTemperatureMin": 2202,
    "colorTemperatureMax": 4000,
    "colorHue": 0,
    "colorSaturation": 0
  },
  "room": {
    "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "name": "Living Room"
  }
}
```

### Example: Set Light Brightness

```
PATCH https://{ip}:8443/v1/devices/{id}
Authorization: Bearer {token}
Content-Type: application/json

[{"attributes": {"isOn": true, "lightLevel": 50}}]
```

### WebSocket Events

The DIRIGERA hub also supports a WebSocket connection for real-time state updates:

```
wss://{ip}:8443/v1/events
Authorization: Bearer {token}
```

Events are JSON objects with `type`, `data`, and `id` fields. This eliminates the need for polling.

## AI Capabilities

AI chat integration is not planned for the initial DIRIGERA milestone, but when implemented the AI could manage lights (on/off, brightness, color, color temperature), activate scenes, query device reachability, and manage rooms.

## Quirks & Notes

- **Self-signed TLS**: The API uses HTTPS with a self-signed certificate. Go's `http.Client` must be configured with `InsecureSkipVerify: true` or the certificate must be pinned after first connection.
- **OAuth2 flow requires physical button press**: There is no way to bypass the button press. The user must be physically near the hub during pairing.
- **Token is long-lived**: The bearer token does not expire and does not need refreshing. Store securely.
- **Thread Border Router**: The DIRIGERA acts as a Thread Border Router, enabling Thread/Matter devices to communicate over IP. This is key for future Matter integration.
- **Backward compatible**: All existing TRADFRI Zigbee devices work with DIRIGERA without re-pairing (if factory reset first).
- **Undocumented API**: The local REST API is not officially documented by IKEA. It was reverse-engineered from the IKEA Home Smart app's traffic. It could change with firmware updates.
- **PATCH array format**: Device updates use a PATCH with a JSON array payload, not a plain object. This is unusual and easy to get wrong.
- **Matter support**: While DIRIGERA supports Matter, it acts as a Matter controller, not a Matter bridge that exposes devices. This means Haus cannot discover DIRIGERA's child devices via Matter; the local REST API is the integration path.
- **Rate limiting**: The hub can handle moderate request rates but may become sluggish under heavy polling. Use the WebSocket for state updates instead.

## Similar Devices

- [ikea-tradfri-gateway](ikea-tradfri-gateway.md) — The predecessor hub (CoAP protocol, discontinued)
- [philips-hue-bridge](philips-hue-bridge.md) — Similar bridge-based approach with a more mature local API
- [lutron-caseta-bridge](lutron-caseta-bridge.md) — Another bridge-based system with local API
