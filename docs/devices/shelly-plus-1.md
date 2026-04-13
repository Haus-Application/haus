---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "shelly-plus-1"
name: "Shelly Plus 1"
manufacturer: "Allterco Robotics EOOD"
brand: "Shelly"
model: "SNSW-001X16EU"
model_aliases: ["SNSW-001P16EU", "Shelly Plus 1", "shellyplus1", "ShellyPlus1"]
device_type: "shelly_relay"
category: "smart_home"
product_line: "Plus"
release_year: 2022
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "34:B7:DA"        # Allterco Robotics (primary OUI for Shelly devices)
    - "E8:DB:84"        # Espressif Systems (ESP32 chipset used in Shelly Plus line)
    - "C8:F0:9E"        # Espressif Systems (alternate)
    - "3C:61:05"        # Espressif Systems (alternate)
    - "EC:62:60"        # Allterco / Shelly (newer production)
  mdns_services:
    - "_shelly._tcp"     # Shelly-specific mDNS service
    - "_http._tcp"       # Standard HTTP service advertisement
  mdns_txt_keys:
    - "gen"             # Generation: "2" for Plus/Pro line
    - "id"              # Device ID (e.g., "shellyplus1-xxxxxxxxxxxx")
    - "model"           # Model identifier
    - "app"             # Application name ("Plus1")
    - "ver"             # Firmware version
    - "arch"            # Architecture ("esp32")
  default_ports: [80]
  signature_ports: [80]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^ShellyPlus1-[0-9A-Fa-f]{12}$"
    - "^shellyplus1-[0-9a-f]{12}$"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/shelly"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"app\":\"Plus1\""
    headers: {}
  - port: 80
    path: "/rpc/Shelly.GetDeviceInfo"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"model\":\"SNSW-001"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "shelly"
  polling_interval_sec: 5
  websocket_event: "shelly:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities:
  - "on_off"

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 80
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "none"
  auth_detail: "No authentication by default. Optional digest authentication can be enabled via device settings. Gen2 RPC API also available over WebSocket (ws://{ip}/rpc) and MQTT."
  base_url_template: "http://{ip}/rpc"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "switch"
  power_source: "hardwired"
  mounting: "in_wall"
  indoor_outdoor: "both"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.shelly.com/en/products/shop/shelly-plus-1"
  api_docs: "https://shelly-api-docs.shelly.cloud/gen2/"
  developer_portal: "https://shelly-api-docs.shelly.cloud/"
  support: "https://support.shelly.cloud/"
  community_forum: "https://www.facebook.com/groups/ShellyIoTCommunitySupport"
  image_url: ""
  fcc_id: "2AZERPLUS1"

# --- TAGS ---
tags: ["wifi", "bluetooth", "relay", "in_wall", "local_api", "rpc", "mqtt", "websocket", "esp32", "no_auth_default", "gen2"]
---

# Shelly Plus 1

## What It Is

The Shelly Plus 1 is a compact WiFi-connected smart relay module from Allterco Robotics (Shelly), designed for in-wall or DIN-rail mounting. It provides a single dry-contact relay capable of switching loads up to 16A (resistive) at 110-240VAC or 24-48VDC. The device is tiny (approximately 42x36x17mm), making it small enough to fit behind a standard wall switch or inside a junction box. The Shelly Plus line is built on the ESP32 platform and represents Gen2 of Shelly's product line, featuring the excellent Gen2 RPC (Remote Procedure Call) API over HTTP, WebSocket, and MQTT. Shelly devices are widely regarded as having the best local API in the consumer smart home market -- well-documented, stable, and fully featured without any cloud dependency. The Plus 1 supports Bluetooth Low Energy for initial setup and can be configured entirely without an internet connection. It also supports Shelly scripts (mJS JavaScript runtime) for on-device automation logic.

## How Haus Discovers It

1. **mDNS Discovery** -- The Shelly Plus 1 advertises `_shelly._tcp.local.` via multicast DNS. The TXT record includes `gen=2`, `id=shellyplus1-xxxxxxxxxxxx`, `model=SNSW-001X16EU`, `app=Plus1`, `ver={firmware_version}`, and `arch=esp32`. This is the primary and most reliable discovery method.
2. **OUI Match** -- MAC addresses beginning with `34:B7:DA` or `EC:62:60` (Allterco) or Espressif OUIs (`E8:DB:84`, `C8:F0:9E`, `3C:61:05`) may indicate a Shelly device. ESP32 OUIs are shared with many other ESP32-based products, so further fingerprinting is needed.
3. **HTTP Fingerprint** -- `GET http://{ip}/shelly` returns a JSON object identifying the device:
   ```json
   {
     "name": null,
     "id": "shellyplus1-xxxxxxxxxxxx",
     "mac": "XXXXXXXXXXXX",
     "slot": 0,
     "model": "SNSW-001X16EU",
     "gen": 2,
     "fw_id": "20240101-000000/1.0.0-gxxxxxxxx",
     "ver": "1.0.0",
     "app": "Plus1",
     "auth_en": false,
     "auth_domain": null
   }
   ```
4. **Hostname Pattern** -- Shelly Plus 1 devices register with hostname `ShellyPlus1-XXXXXXXXXXXX` (12 hex characters from MAC) on DHCP.
5. **Port Probe** -- Port 80/TCP responds with the Gen2 RPC API. Unlike many IoT devices, Shelly's port 80 serves a proper API, not just a web interface.

## Pairing / Authentication

No pairing or authentication is required by default. The Shelly Plus 1 responds to unauthenticated HTTP, WebSocket, and MQTT requests out of the box.

### Optional Authentication

The device supports optional HTTP Digest Authentication, which can be enabled via:

```
POST http://{ip}/rpc/Shelly.SetAuth
Content-Type: application/json

{
  "user": "admin",
  "realm": "shellyplus1-xxxxxxxxxxxx",
  "ha1": "{MD5(user:realm:password)}"
}
```

The `ha1` parameter is the MD5 hash of `admin:shellyplus1-xxxxxxxxxxxx:password`. Once enabled, all RPC requests require HTTP Digest auth.

### Bluetooth Provisioning

New devices can be configured via BLE using the Shelly app. The BLE provisioning sets WiFi credentials and basic configuration without needing a separate AP setup step. After WiFi connection, BLE can be disabled to save power.

## API Reference

The Shelly Plus 1 uses the Shelly Gen2 RPC API. All methods are accessible via three transports.

### Transports

| Transport | URL | Notes |
|-----------|-----|-------|
| HTTP | `http://{ip}/rpc/{Method}` | GET with query params or POST with JSON body |
| WebSocket | `ws://{ip}/rpc` | JSON-RPC 2.0 over WebSocket, bidirectional |
| MQTT | topic `{device_id}/rpc` | JSON-RPC 2.0 over MQTT |

### Device Info

```
GET http://{ip}/rpc/Shelly.GetDeviceInfo
```

**Response:**
```json
{
  "id": "shellyplus1-xxxxxxxxxxxx",
  "name": "Kitchen Light",
  "mac": "XXXXXXXXXXXX",
  "model": "SNSW-001X16EU",
  "gen": 2,
  "fw_id": "20240101-000000/1.0.0-gxxxxxxxx",
  "ver": "1.0.0",
  "app": "Plus1",
  "profile": "switch",
  "auth_en": false,
  "auth_domain": "shellyplus1-xxxxxxxxxxxx"
}
```

### Get Switch Status

```
GET http://{ip}/rpc/Switch.GetStatus?id=0
```

**Response:**
```json
{
  "id": 0,
  "source": "http",
  "output": true,
  "temperature": {"tC": 42.3, "tF": 108.1},
  "aenergy": {"total": 1234.56, "by_minute": [12.3, 11.8, 12.1], "minute_ts": 1712000060}
}
```

- `output` -- boolean, current relay state (true = closed/on)
- `temperature` -- internal device temperature (for thermal protection monitoring)
- `aenergy` -- only on models with energy monitoring (Plus 1PM, not base Plus 1)

### Set Switch State

```
GET http://{ip}/rpc/Switch.Set?id=0&on=true
```

Or via POST:
```
POST http://{ip}/rpc/Switch.Set
Content-Type: application/json

{"id": 0, "on": true}
```

- `id` -- switch component ID (0 for Plus 1, since it has only one relay)
- `on` -- boolean, desired state

### Toggle Switch

```
GET http://{ip}/rpc/Switch.Toggle?id=0
```

Toggles the relay to the opposite state. Returns the new state.

### Get Configuration

```
GET http://{ip}/rpc/Switch.GetConfig?id=0
```

**Response:**
```json
{
  "id": 0,
  "name": "Kitchen Light",
  "in_mode": "follow",
  "initial_state": "restore_last",
  "auto_on": false,
  "auto_on_delay": 0,
  "auto_off": false,
  "auto_off_delay": 0
}
```

### Set Configuration

```
POST http://{ip}/rpc/Switch.SetConfig
Content-Type: application/json

{
  "id": 0,
  "config": {
    "name": "Kitchen Light",
    "in_mode": "follow",
    "initial_state": "restore_last",
    "auto_on": true,
    "auto_on_delay": 300,
    "auto_off": false,
    "auto_off_delay": 0
  }
}
```

**Input modes (`in_mode`):**
- `"follow"` -- switch input directly controls relay (standard wall switch behavior)
- `"flip"` -- each toggle of the physical switch flips the relay state
- `"detached"` -- physical switch is disconnected from relay, only sends events (for smart bulb setups)
- `"momentary"` -- momentary/push-button input

**Initial state (`initial_state`):**
- `"on"` -- relay on after power restore
- `"off"` -- relay off after power restore
- `"restore_last"` -- restore state before power loss
- `"match_input"` -- match the physical switch position

### WiFi Status

```
GET http://{ip}/rpc/WiFi.GetStatus
```

Returns current WiFi connection info including SSID, RSSI, IP address.

### System Status

```
GET http://{ip}/rpc/Sys.GetStatus
```

Returns uptime, RAM usage, filesystem usage, firmware info.

### WebSocket Real-Time Notifications

Connect to `ws://{ip}/rpc` and send:
```json
{"id": 1, "src": "haus", "method": "Switch.GetStatus", "params": {"id": 0}}
```

The device also pushes unsolicited notifications on state changes:
```json
{
  "src": "shellyplus1-xxxxxxxxxxxx",
  "dst": "haus",
  "method": "NotifyStatus",
  "params": {
    "ts": 1712000000.00,
    "switch:0": {"id": 0, "output": true, "source": "button"}
  }
}
```

### Action Webhooks

The device supports configurable webhooks (actions) that fire on state changes:
```
POST http://{ip}/rpc/Webhook.Create
Content-Type: application/json

{
  "cid": 0,
  "enable": true,
  "event": "switch.on",
  "urls": ["http://{haus_ip}:{port}/webhook/shelly"]
}
```

Events: `switch.on`, `switch.off`, `switch.toggle`

## AI Capabilities

When the AI concierge is chatting with a Shelly Plus 1, it can:

- **Toggle the relay** on or off
- **Report current state** -- on/off, internal temperature, uptime
- **Configure behavior** -- input mode (follow, flip, detached, momentary), auto-on/off timers, initial state after power loss
- **Report WiFi info** -- SSID, signal strength, IP address
- **Explain input modes** -- describe how different physical switch configurations work
- **Monitor device health** -- temperature, RAM usage, firmware version

## Quirks & Notes

- **Best-in-Class Local API:** The Shelly Gen2 RPC API is widely considered the gold standard for local smart home device APIs. It is comprehensive, well-documented, stable across firmware updates, and works over HTTP, WebSocket, and MQTT simultaneously.
- **No Authentication by Default:** Out of the box, anyone on the network can control the device. Users should be informed about enabling digest auth if security is a concern.
- **ESP32 Platform:** Built on the Espressif ESP32, which means MAC addresses may appear under Espressif OUIs rather than Allterco. The mDNS `_shelly._tcp` advertisement is the most reliable identifier.
- **Scripting Engine:** The Plus 1 includes an mJS (Mongoose JS) scripting runtime that can execute user-defined scripts directly on the device for local automations. This runs independently of any hub.
- **Tiny Form Factor:** At 42x36x17mm, it fits behind most wall switches. However, the wiring can be tight in shallow switch boxes (especially US 14 cubic inch boxes).
- **Dry Contact Relay:** The relay is a dry contact, meaning it can switch AC or DC loads. This makes it versatile for controlling pumps, valves, garage doors, or any switched load.
- **Bluetooth Setup Only:** BLE is used for initial configuration via the Shelly app. It is not used for ongoing control. BLE can be disabled after setup.
- **No Dimming:** The Plus 1 is a relay (on/off only). For dimming, use the Shelly Plus Dimmer 0-10V or Shelly Plus 2PM in roller shutter mode.
- **Firmware Updates:** OTA firmware updates are available through the Shelly cloud or via local RPC (`Shelly.Update` method). Shelly's firmware update track record is excellent, with regular updates that generally do not break the local API.
- **Temperature Protection:** The device monitors its internal temperature and will shut off the relay if it exceeds approximately 95 degrees C to prevent damage.
- **MQTT Support:** The device supports MQTT natively. MQTT can be configured to connect to any broker (not just Shelly's cloud). This makes it compatible with Home Assistant's MQTT integration, Mosquitto, or any standard MQTT infrastructure.
- **Power Monitoring:** The base Plus 1 does NOT have power monitoring. The Plus 1PM variant adds voltage, current, and power measurement.

## Similar Devices

- **shelly-plus-1pm** -- Same form factor with power monitoring (voltage, current, watts)
- **shelly-plus-2pm** -- Dual-relay version with power monitoring, also supports roller shutter mode
- **shelly-plus-plug-s** -- Plug-in version with power monitoring, same Gen2 API
- **sonoff-basic-r4** -- Similar WiFi relay concept but uses eWeLink cloud (limited local API)
- **kasa-smart-plug** -- TP-Link plug, XOR protocol local control, different form factor
- **wemo-smart-plug** -- Belkin plug, SOAP/UPnP API, similar use case
