---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "tailwind-iq3"
name: "Tailwind iQ3 Smart Garage Controller"
manufacturer: "Tailwind"
brand: "Tailwind"
model: "iQ3"
model_aliases: ["TW3000", "iQ3-1", "iQ3-2"]
device_type: "garage_controller"
category: "security"
product_line: "Tailwind iQ"
release_year: 2023
discontinued: false
price_range: "$$"

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
    - "24:62:AB"        # Espressif (ESP32 module used in iQ3)
    - "A4:CF:12"        # Espressif Systems
    - "DC:54:75"        # Espressif Systems
  mdns_services:
    - "_http._tcp"      # Local HTTP API advertised via mDNS
  mdns_txt_keys:
    - "tailwind"        # Device identifier in TXT records
  default_ports: [80]   # Local HTTP API
  signature_ports: [80]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^tailwind"
    - "^Tailwind"
    - "^TW[0-9]+"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: "Tailwind"
    server_header: ""
    body_contains: "tailwind"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "tailwind"
  polling_interval_sec: 10
  websocket_event: ""
  setup_type: "api_key"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "garage_open_close"

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 80
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "Local API uses a token generated in the Tailwind app. Send token as a query parameter or in a custom header. The local API is documented and officially supported by Tailwind. UDP discovery protocol also available for LAN device discovery."
  base_url_template: "http://{ip}"
  tls: false
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
  product_page: "https://gotailwind.com/products/tailwind-iq3"
  api_docs: "https://github.com/paulwieland/tailwind"
  developer_portal: ""
  support: "https://gotailwind.com/pages/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2A6MKTW3000"

# --- TAGS ---
tags: ["garage", "local_api", "http_rest", "developer_friendly", "esp32", "multi_door", "tailwind", "udp_discovery"]
---

# Tailwind iQ3 Smart Garage Controller

## What It Is

The Tailwind iQ3 is a Wi-Fi connected smart garage door controller manufactured by Tailwind. It distinguishes itself from competitors like MyQ by being explicitly developer-friendly with a documented local HTTP API, making it one of the most integration-friendly garage controllers available. The iQ3 supports up to 3 garage doors from a single unit (the "3" in iQ3), using individual door sensors for each door and a multi-relay output to control separate openers. It connects to Wi-Fi (2.4GHz) and communicates with the Tailwind cloud for remote access, but the local API works fully without internet. The device is powered by a USB-C adapter (5V) and mounts on the garage wall near the opener(s). It supports auto-open via vehicle Bluetooth detection, scheduling, and geofencing.

## How Haus Discovers It

1. **OUI Match** -- The iQ3 uses an Espressif ESP32 module, so MAC prefixes will be Espressif OUIs (`24:62:AB`, `A4:CF:12`, `DC:54:75`). These OUIs are shared with many IoT devices and are not conclusive alone.
2. **mDNS Discovery** -- The iQ3 advertises `_http._tcp` via mDNS with `tailwind` in the TXT records. This is the most reliable local discovery method.
3. **HTTP Fingerprint** -- `GET http://{ip}/` on port 80 returns a response identifying the device as a Tailwind controller.
4. **UDP Discovery** -- The iQ3 supports a UDP broadcast discovery protocol. Sending a specific discovery packet to port 9988 (UDP broadcast) returns device information including model, firmware version, and number of configured doors.
5. **Hostname Pattern** -- The device typically appears with hostname `tailwind*` or `TW*` in DHCP tables.

## Pairing / Authentication

### Tailwind App Setup

1. **Account Creation:** Create a Tailwind account in the mobile app.
2. **Device Pairing:** The app discovers the iQ3 via BLE during initial setup and provisions Wi-Fi credentials.
3. **Door Sensor Pairing:** Each door sensor (magnetic tilt sensor) is paired with the controller.
4. **Opener Wiring:** The controller's relay outputs are wired to each garage door opener's button terminals (dry contact closure).

### Local API Token

1. **Token Generation:** In the Tailwind app, navigate to device settings and generate a local API token. This creates a long-lived authentication token.
2. **Token Usage:** Include the token in API requests as a query parameter (`?token={token}`) or in a custom header.
3. **Token Management:** Tokens can be regenerated (which invalidates the old token) from the app.

### Security Notes

- The local API uses HTTP (not HTTPS) -- all communication is unencrypted on the local network.
- The API token provides full control of all doors, so it should be treated as a sensitive credential.
- Tailwind is notably pro-integration, explicitly supporting third-party local access.

## API Reference

**Base URL:** `http://{ip}`

### Get Device Status

```
GET http://{ip}/status?token={token}
```

**Response:**
```json
{
  "product": "iQ3",
  "firmware": "10.10",
  "data": {
    "door1": {
      "state": "closed",
      "last_changed": "2024-01-01T12:00:00Z"
    },
    "door2": {
      "state": "open",
      "last_changed": "2024-01-01T11:30:00Z"
    },
    "door3": {
      "state": "closed",
      "last_changed": "2024-01-01T10:00:00Z"
    }
  }
}
```

Door states: `open`, `closed`, `opening`, `closing`

### Control Door

```
POST http://{ip}/cmd?token={token}
Content-Type: application/json

{
  "door": 1,
  "cmd": "open"
}
```

Commands: `open`, `close`, `toggle`

Door numbers: `1`, `2`, `3` (corresponding to the connected doors)

**Response:**
```json
{
  "result": "OK",
  "door": 1,
  "cmd": "open"
}
```

### Get Device Info

```
GET http://{ip}/device?token={token}
```

Returns device information:
```json
{
  "product": "iQ3",
  "firmware": "10.10",
  "mac": "24:62:AB:XX:XX:XX",
  "ip": "192.168.1.100",
  "ssid": "HomeNetwork",
  "rssi": -45,
  "doors_configured": 2,
  "uptime": 86400
}
```

### UDP Discovery Protocol

**Discovery Request:** Send a UDP broadcast to port 9988:
```
TAILWIND_DISCOVER
```

**Discovery Response:** The iQ3 responds with a JSON payload:
```json
{
  "product": "iQ3",
  "mac": "24:62:AB:XX:XX:XX",
  "ip": "192.168.1.100",
  "firmware": "10.10"
}
```

### Firmware Update Check

```
GET http://{ip}/firmware?token={token}
```

Returns current firmware version and whether an update is available.

## AI Capabilities

AI integration is not currently planned but the local API makes it very feasible. If implemented, the AI concierge could:

- Report the state of each garage door (open/closed)
- Open/close specific doors by name/number with security confirmation
- Report how long a door has been in its current state
- Alert if a door has been open for an unusual duration
- Report device connectivity and signal strength

## Quirks & Notes

- **Developer-Friendly:** Tailwind is the anti-MyQ. They explicitly support and document local API access, have worked with the Home Assistant community, and encourage third-party integration. This is their primary competitive advantage over MyQ/Chamberlain.
- **Multi-Door Support:** A single iQ3 unit controls up to 3 garage doors. Each door gets its own tilt sensor and relay output. The controller has 3 wired relay terminals -- one per door.
- **ESP32-Based:** The iQ3 runs on an Espressif ESP32 microcontroller with Wi-Fi. This means the local API is implemented on a resource-constrained embedded device -- keep polling intervals reasonable (10+ seconds).
- **HTTP Only (No TLS):** The local API uses plain HTTP. All traffic including the API token is unencrypted on the local network. This is typical for embedded IoT devices but worth noting for security-conscious deployments.
- **USB-C Power:** Powered by a USB-C cable and wall adapter (5V/2A). The device does not have a battery backup -- power outage means no smart control (though the garage door opener itself may have backup battery).
- **Tilt Sensors:** The iQ3 uses wireless tilt sensors (battery-powered, CR2032) that mount on the garage door. The sensors detect door position by orientation change (tilt). They communicate with the controller via a sub-GHz radio link (not Wi-Fi or BLE). Battery life is approximately 1-2 years.
- **Vehicle Bluetooth Detection:** The iQ3 can detect your vehicle's Bluetooth signal as you approach and auto-open the garage door. This works with most cars that have Bluetooth (it scans for known MAC addresses). This is more reliable than phone geofencing for garage auto-open.
- **Cloud Features:** While the local API handles basic open/close/status, some features (scheduling, geofencing, shared access, activity log) are managed through the Tailwind cloud. The cloud is not required for basic local operation.
- **Home Assistant Integration:** An official Home Assistant integration exists for the Tailwind iQ3, using the local API. This confirms the local API is stable and supported.
- **Firmware Updates:** OTA firmware updates are delivered via the Tailwind cloud through the app. The ESP32 platform supports reliable OTA updates.

## Similar Devices

- **myq-smart-garage-hub** -- Market leader but cloud-locked with hostile API policy
- **meross-smart-garage-opener** -- Budget option with local MQTT protocol
