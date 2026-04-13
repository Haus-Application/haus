---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "meross-smart-garage-opener"
name: "Meross Smart Garage Door Opener (MSG100)"
manufacturer: "Meross Technology Co., Ltd."
brand: "Meross"
model: "MSG100"
model_aliases: ["MSG100HK", "MSG200", "MSG100-UN"]
device_type: "garage_controller"
category: "security"
product_line: "Meross Smart"
release_year: 2019
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
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "48:E1:E9"        # Meross Technology
    - "34:29:8F"        # Meross Technology
  mdns_services: []     # Meross devices do not advertise mDNS
  mdns_txt_keys: []
  default_ports: []     # No standard open ports for discovery
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Meross"
    - "^meross"
    - "^MSG"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # Uses MQTT, not HTTP

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "meross"
  polling_interval_sec: 15
  websocket_event: ""
  setup_type: "password"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "garage_open_close"

# --- PROTOCOL ---
protocol:
  type: "mqtt"
  port: 8883
  transport: "TCP"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "Meross uses MQTT for both cloud and local communication. Cloud MQTT broker at mqtt.meross.com (port 8883 TLS). Local MQTT communication uses the same protocol but on the device's IP directly (port 80 HTTP with MQTT-like JSON payloads). Authentication via signing key derived from user credentials and device UUID. Messages are signed with MD5 HMAC using a key derived from the user's cloud token."
  base_url_template: ""
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "controller"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.meross.com/product/smart-garage-door-opener"
  api_docs: ""
  developer_portal: ""
  support: "https://www.meross.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AQ3Z-MSG100"

# --- TAGS ---
tags: ["garage", "mqtt", "local_mqtt", "meross", "budget", "homekit_variant", "reverse_engineered", "md5_signing"]
---

# Meross Smart Garage Door Opener (MSG100)

## What It Is

The Meross Smart Garage Door Opener (MSG100) is a budget-friendly Wi-Fi garage door controller manufactured by Meross Technology, a Chinese smart home device company. It connects to your existing garage door opener via a dry-contact relay and includes a wired magnetic sensor to detect door open/closed state. The device connects to Wi-Fi (2.4GHz) and communicates via MQTT protocol -- both to Meross's cloud and locally on the LAN. The MSG100 is notable for its low price point (typically $30-40 USD), wide compatibility with garage door openers, and a well-understood local MQTT protocol that has been reverse-engineered by the community. A HomeKit-compatible variant (MSG100HK) is also available. The MSG200 model supports up to 3 doors. The Meross protocol is the same across their entire product line (smart plugs, switches, garage openers), so understanding one Meross device enables integration with all of them.

## How Haus Discovers It

1. **OUI Match** -- Devices with MAC prefixes `48:E1:E9` or `34:29:8F` (Meross Technology) are flagged during network scanning.
2. **Hostname Pattern** -- Meross devices may appear with hostnames starting with `Meross`, `meross`, or their model number in DHCP tables.
3. **Cloud API Discovery** -- After Meross account authentication, the cloud API returns a list of all devices with their LAN IP addresses, UUIDs, and hardware/firmware versions.
4. **Local Protocol Probe** -- Meross devices listen on port 80 for local HTTP-wrapped MQTT-like messages. A probe message to this port with the correct format will elicit a response confirming the device type.

## Pairing / Authentication

### Meross App Setup

1. **Account Creation:** Create a Meross account in the mobile app.
2. **Device Pairing:** Put the MSG100 in pairing mode (press button for 5 seconds). The app connects to the device's temporary AP (access point), provisions Wi-Fi credentials, and registers the device with Meross's cloud.
3. **Sensor Installation:** Mount the magnetic reed sensor on the garage door and connect its wire to the MSG100.
4. **Opener Wiring:** Connect the MSG100's relay output (two wires) to the garage door opener's wall button terminals.

### Authentication Architecture

Meross uses a signing-based authentication system for both cloud and local communication:

1. **Cloud Login:** `POST https://iotx-us.meross.com/v1/Auth/Login` (or region-specific endpoint) with:
   ```json
   {
     "email": "{email}",
     "password": "{md5_of_password}"
   }
   ```
   Note: The password is sent as an MD5 hash, not plaintext. Returns `token`, `key`, and `userid`.

2. **Message Signing:** All MQTT messages (cloud and local) are signed using:
   - `messageId` -- random UUID
   - `timestamp` -- Unix timestamp (seconds)
   - `sign` -- MD5 hash of `messageId + key + timestamp`
   
   The `key` is derived from the cloud login response. This signing mechanism prevents replay attacks and unauthorized commands.

3. **MQTT Cloud Broker:** Devices connect to `mqtt.meross.com:8883` (TLS) or region-specific brokers. Topics are structured as `/appliance/{device_uuid}/subscribe` and `/appliance/{device_uuid}/publish`.

4. **Local HTTP/MQTT:** Devices also accept commands on their LAN IP address at port 80 via HTTP POST, using the same JSON message format and signing as the MQTT cloud protocol.

### Security Notes

- The MD5-based signing is cryptographically weak by modern standards but functional for home automation.
- Local communication on port 80 is unencrypted HTTP.
- The `key` obtained from cloud login is required for both cloud and local signing, creating a one-time cloud dependency.
- The protocol has been fully reverse-engineered and documented by several community projects (meross-iot Python library, MerossIot .NET library).

## API Reference

### Meross Message Format

All Meross communication (cloud MQTT and local HTTP) uses the same JSON message envelope:

```json
{
  "header": {
    "from": "/appliance/{uuid}/publish",
    "messageId": "{random-uuid}",
    "method": "GET",
    "namespace": "Appliance.System.All",
    "payloadVersion": 1,
    "sign": "{md5_signature}",
    "timestamp": 1704067200,
    "triggerSrc": "iOSLocal",
    "uuid": "{device_uuid}"
  },
  "payload": {}
}
```

**Methods:** `GET`, `SET`, `PUSH` (device-initiated events)

### Get Device Status (All Info)

**Namespace:** `Appliance.System.All`
**Method:** `GET`

```json
{
  "header": {
    "method": "GET",
    "namespace": "Appliance.System.All",
    ...
  },
  "payload": {}
}
```

Returns comprehensive device information including hardware version, firmware version, Wi-Fi signal, and all capability states.

### Get Garage Door State

**Namespace:** `Appliance.GarageDoor.State`
**Method:** `GET`

**Response payload:**
```json
{
  "state": [{
    "channel": 0,
    "open": 0,
    "lmTime": 1704067200
  }]
}
```

- `channel` -- door index (0 for MSG100 single door, 0-2 for MSG200)
- `open` -- 0 = closed, 1 = open
- `lmTime` -- last modified timestamp (Unix seconds)

### Open/Close Garage Door

**Namespace:** `Appliance.GarageDoor.State`
**Method:** `SET`

```json
{
  "header": {
    "method": "SET",
    "namespace": "Appliance.GarageDoor.State",
    ...
  },
  "payload": {
    "state": {
      "channel": 0,
      "open": 1,
      "uuid": "{device_uuid}"
    }
  }
}
```

Set `open: 1` to open, `open: 0` to close.

### Local HTTP Endpoint

```
POST http://{device_ip}/config
Content-Type: application/json

{full message JSON including header with sign}
```

The local endpoint accepts the same message format as the MQTT protocol. The device responds with the result in the same JSON envelope format.

### MQTT Topics (Cloud)

- **Device Subscribe (commands to device):** `/appliance/{uuid}/subscribe`
- **Device Publish (events from device):** `/appliance/{uuid}/publish`
- **App Topic:** `/app/{userid}-{appid}/subscribe`

### Common Namespaces

- `Appliance.System.All` -- Full device state
- `Appliance.System.Online` -- Online/offline status
- `Appliance.System.Firmware` -- Firmware info
- `Appliance.System.Hardware` -- Hardware info
- `Appliance.GarageDoor.State` -- Door open/close state and control
- `Appliance.Control.ToggleX` -- Generic toggle (for relay devices)

## AI Capabilities

AI integration is not currently planned but the local protocol makes it feasible. If implemented, the AI concierge could:

- Report garage door state (open/closed)
- Open/close the door with security confirmation
- Report how long the door has been in current state
- Report device connectivity and signal strength
- Alert if the door has been open for too long

## Quirks & Notes

- **Meross Protocol Universality:** The MSG100 uses the exact same MQTT/HTTP protocol as all other Meross devices (smart plugs MSS110/MSS210, switches, power strips). The only difference is the namespace for device-specific capabilities (e.g., `GarageDoor.State` vs `Control.ToggleX`). Implementing Meross protocol support enables integration with their entire product line.
- **One-Time Cloud Dependency:** The signing `key` must be obtained via the cloud API login. Once obtained and stored, all subsequent local communication can happen without the cloud. However, if the key changes (account password change, device re-pairing), a new cloud login is needed.
- **MD5 Everywhere:** Meross uses MD5 extensively -- for password hashing, message signing, and device UUIDs. This is not ideal from a cryptographic perspective but is well-understood and consistent.
- **Port 80 Local API:** The local API on port 80 uses HTTP POST to `/config`. Despite being on an HTTP port, the messages are not really REST -- they use the same MQTT-style namespace/method/payload format wrapped in HTTP.
- **MSG200 Multi-Door:** The MSG200 variant supports up to 3 doors using `channel` indices 0, 1, 2. Each channel has an independent relay and sensor. The protocol is identical except for the channel parameter.
- **MSG100HK (HomeKit):** The HomeKit variant adds Apple HomeKit support via the HAP protocol. It can be controlled both via HomeKit and the standard Meross cloud/local protocol simultaneously.
- **Magnetic Sensor Wired:** Unlike the Tailwind iQ3's wireless tilt sensors, the MSG100 uses a wired magnetic reed sensor. The sensor connects to the controller via a long wire (included, approximately 6 feet). The sensor mounts on the garage door frame and a magnet mounts on the door itself. This is simple and reliable but requires running a wire.
- **Budget Price:** At $30-40 USD, the MSG100 is one of the cheapest smart garage controllers available. The build quality is appropriate for the price -- functional but not premium.
- **Firmware Updates:** OTA firmware updates are delivered via the Meross cloud. Some updates have changed protocol behavior, though the core message format has remained stable.
- **Region-Specific Servers:** Meross uses different cloud servers by region (US, EU, AP). The device is bound to a region during setup. Cloud API endpoints and MQTT brokers differ by region: `iotx-us.meross.com`, `iotx-eu.meross.com`, etc.

## Similar Devices

- **myq-smart-garage-hub** -- Market leader but cloud-locked
- **tailwind-iq3** -- Premium option with official local API and multi-door support
- **meross-smart-plug-mss110** -- Same protocol family, different device type
