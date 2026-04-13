---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "pentair-screenlogic"
name: "Pentair ScreenLogic Interface"
manufacturer: "Pentair Water Pool and Spa Inc."
brand: "Pentair"
model: "ScreenLogic"
model_aliases: ["ScreenLogic2", "IntelliTouch ScreenLogic", "EasyTouch ScreenLogic", "522104"]
device_type: "pool_controller"
category: "smart_home"
product_line: "Pentair IntelliTouch / EasyTouch"
release_year: 2012
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["wifi", "ethernet"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["00:C0:33"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [80, 500]
  signature_ports: [80, 500]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^Pentair.*", "^screenlogic.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: "ScreenLogic"
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "pentair"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "none"
  ai_chattable: false
  haus_milestone: "post-V1"

# --- CAPABILITIES ---
capabilities: ["on_off", "temperature"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 80
  transport: "TCP"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Custom binary protocol over TCP. The ScreenLogic interface uses UDP broadcast on port 1444 for discovery, then TCP connections on a dynamically assigned port (typically 80 or a high port) for commands. No authentication required on the local network. Messages are binary-encoded with a custom header format."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "gateway"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.pentair.com/en-us/products/residential/pool-spa-equipment/pool-automation/screenlogic2-interface-kit.html"
  api_docs: ""
  developer_portal: ""
  support: "https://www.pentair.com/en-us/support.html"
  community_forum: "https://www.troublefreepool.com"
  image_url: ""
  fcc_id: "BUQ-SL2"

# --- TAGS ---
tags: ["pool", "spa", "pool-controller", "binary-protocol", "udp-discovery", "no-auth", "temperature", "pump", "lights", "chemistry", "pentair"]
---

# Pentair ScreenLogic Interface

## What It Is

> The Pentair ScreenLogic Interface is a WiFi/Ethernet adapter that connects to Pentair IntelliTouch, EasyTouch, and SunTouch pool and spa automation systems, enabling remote monitoring and control. It bridges the pool controller's RS-485 bus to the home network, providing access to pool/spa temperature, pump circuits, lighting, heater, chemistry (if IntelliChem is installed), and scheduling. The ScreenLogic uses a custom binary protocol over TCP for local communication and supports Pentair's cloud service for remote access. The protocol has been extensively reverse-engineered by the community, with libraries available in Python, Node.js, and other languages.

## How Haus Discovers It

1. **UDP broadcast discovery** -- Send a UDP broadcast to port 1444. ScreenLogic adapters respond with their name, IP address, and gateway port number
2. **OUI match** -- Pentair MAC prefix: `00:C0:33`
3. **Port probe** -- HTTP on port 80 may serve a basic status page; the primary protocol port is returned during UDP discovery
4. **Hostname pattern** -- DHCP hostname may contain `Pentair` or `screenlogic`

### UDP Discovery Packet

Send to broadcast address on UDP port 1444:
```
Bytes: 01 00 00 00
```

**Response:**
```
Gateway Name: "Pentair: 00-C0-33-XX-XX-XX"
Gateway IP: {ip_address}
Gateway Port: {port}  (typically 80)
Gateway Type: {type}
Gateway Subtype: {subtype}
```

## Pairing / Authentication

> No authentication required. The ScreenLogic protocol is completely open on the local network. Any device that can reach the ScreenLogic adapter's TCP port can send commands and read state.

### Connection Sequence

1. Discover via UDP broadcast on port 1444
2. Open TCP connection to the returned IP and port
3. Send a Connect message (message ID 0)
4. Send a Challenge message (message ID 14) -- returns a challenge response but no actual auth is enforced
5. Send a Login message (message ID 27) -- responds with gateway version
6. Begin sending status queries and commands

## API Reference

### Binary Protocol Overview

All messages use a binary format with the following header:

```
Bytes 0-1:   Sender ID (uint16 LE)
Bytes 2-3:   Message ID (uint16 LE)
Bytes 4-7:   Data Length (uint32 LE)
Bytes 8+:    Message data (variable length)
```

### Key Message Types

| Message ID | Name | Direction | Description |
|-----------|------|-----------|-------------|
| 0 | Connect | Client -> GW | Initial connection |
| 14 | Challenge | Client -> GW | Auth challenge (informational) |
| 27 | Login | Client -> GW | Login with version |
| 12500 | GetPoolStatus | Client -> GW | Get pool/spa status |
| 12502 | GetPumpStatus | Client -> GW | Get pump status |
| 12504 | GetChemData | Client -> GW | Get chemistry data |
| 12506 | GetScheduleData | Client -> GW | Get schedules |
| 12530 | SetCircuitState | Client -> GW | Turn circuit on/off |
| 12528 | SetSetPoint | Client -> GW | Set temperature setpoint |
| 12532 | SetHeatMode | Client -> GW | Set heater mode |
| 12576 | SetLightMode | Client -> GW | Set light color/mode |

### Pool Status (Message 12500)

**Request:** Header only, no data payload.

**Response structure:**
```
OK (4 bytes)
Freeze Mode (4 bytes)
Pool/Spa Status:
  Current Temperature (4 bytes) - pool temp in degrees
  Heat Set Point (4 bytes) - target temp
  Heat Mode (4 bytes) - 0=off, 1=solar, 2=solar_preferred, 3=heat_pump, 4=gas
  Heat Status (4 bytes) - 0=off, non-zero=heating
Spa section: (same structure)
Air Temperature (4 bytes)
Solar Temperature (4 bytes)
Circuits[]: (array)
  Circuit ID (4 bytes)
  State (4 bytes) - 0=off, 1=on
  Color Mode (4 bytes)
  Color Set (4 bytes)
  Color Position (4 bytes)
  Color Speed (4 bytes)
  Delay (4 bytes)
```

### Set Circuit State (Message 12530)

Turn a circuit (pump, light, feature, etc.) on or off:

**Data payload:**
```
Controller Index (4 bytes LE): 0
Circuit ID (4 bytes LE): {circuit_id}
State (4 bytes LE): 0=off, 1=on
```

### Common Circuit IDs

| ID | Default Assignment |
|----|--------------------|
| 1 | Spa |
| 2 | Spa Jets |
| 3 | Spa Blower |
| 5 | Pool |
| 6 | Pool Light |
| 7 | Spa Light |
| 8 | Cleaner |
| 9 | Water Feature |
| 10 | Spillway |

Circuit assignments are configurable; query the circuit list for actual mappings.

### Set Temperature Setpoint (Message 12528)

**Data payload:**
```
Controller Index (4 bytes LE): 0
Body (4 bytes LE): 0=pool, 1=spa
Temperature (4 bytes LE): {degrees_f}
```

### Chemistry Data (Message 12504)

Returns IntelliChem data (if installed):

```
pH (4 bytes): pH * 100 (e.g., 740 = 7.40)
ORP (4 bytes): ORP in mV
Saturation Index (4 bytes): SI * 100
Salt Level (4 bytes): ppm
TDS (4 bytes): ppm
Calcium Hardness (4 bytes): ppm
CYA (4 bytes): ppm
Alkalinity (4 bytes): ppm
```

## AI Capabilities

> AI integration planned. When available:
> - Report pool and spa temperatures (current and setpoints)
> - Turn pool/spa pumps, lights, and features on/off
> - Set temperature setpoints
> - Report heater status and mode
> - Report chemistry data (pH, ORP, salt, alkalinity) if IntelliChem installed
> - Control pool light colors/modes
> - Report air and solar collector temperatures

## Quirks & Notes

- **Binary protocol** -- All communication uses a custom binary protocol, not REST/JSON; requires careful byte-level parsing
- **No authentication** -- The ScreenLogic protocol has zero security; any device on the local network has full control of the pool system
- **UDP discovery on port 1444** -- Send a 4-byte packet to discover ScreenLogic adapters; they respond with their gateway name and TCP port
- **Little-endian encoding** -- All multi-byte integers are little-endian
- **Circuit IDs are configurable** -- Default circuit assignments vary by pool configuration; always query the circuit list rather than hardcoding IDs
- **Temperature units** -- Temperatures are in Fahrenheit by default; check the controller configuration for unit setting
- **Community libraries** -- `node-screenlogic` (Node.js), `screenlogicpy` (Python), and others provide well-tested implementations of the binary protocol
- **Connection keepalive** -- TCP connections may time out after periods of inactivity; implement keepalive or reconnection logic
- **IntelliTouch vs EasyTouch** -- The ScreenLogic interface works with both IntelliTouch (high-end) and EasyTouch (mid-range) automation systems; capabilities vary based on the underlying controller
- **Firmware updates** -- ScreenLogic firmware can be updated via the Pentair app; protocol changes between versions are generally backward-compatible
- **RS-485 bus** -- The ScreenLogic adapter connects to the pool controller's RS-485 communication bus; it acts as another device on the bus alongside indoor panels and remotes
- **IntelliChem optional** -- Chemistry data is only available if a Pentair IntelliChem controller is installed and connected to the automation system

## Similar Devices

> - [Flo by Moen Smart Water Shutoff](flo-by-moen-smart-water-shutoff.md) -- Water monitoring and control (different category but similar home infrastructure)
