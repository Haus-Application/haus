---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "lutron-caseta-bridge"
name: "Lutron Caseta Smart Bridge"
manufacturer: "Lutron Electronics"
brand: "Caseta"
model: "L-BDG2-WH"
model_aliases: ["L-BDGPRO2-WH", "Caseta Smart Bridge", "Caseta Smart Bridge Pro", "L-BDG3-WH"]
device_type: "caseta_bridge"
category: "lighting"
product_line: "Caseta"
release_year: 2014
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["ethernet", "wifi", "clear_connect"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["08:F1:EA", "28:9A:4B", "58:47:CA"]
  mdns_services: ["_leap._tcp"]
  mdns_txt_keys: []
  default_ports: [8081, 8083]
  signature_ports: [8083]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^Lutron-[A-F0-9]+$", "^Caseta-[A-F0-9]+$"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 8081
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: ""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "caseta"
  polling_interval_sec: 0
  websocket_event: "caseta:state"
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "scenes"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 8083
  transport: "TLS"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "LEAP protocol over TLS on port 8083. Standard bridge requires certificate-based mutual TLS auth obtained via Lutron cloud OAuth. Smart Bridge Pro also offers telnet on port 23 with LIPII protocol (no auth)."
  base_url_template: "tls://{ip}:8083"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["clear_connect"]

# --- LINKS ---
links:
  product_page: "https://www.casetawireless.com/"
  api_docs: "https://www.lutron.com/TechnicalDocumentLibrary/040249.pdf"
  developer_portal: "https://developer.lutron.com/"
  support: "https://www.lutron.com/en-US/Support/Pages/default.aspx"
  community_forum: "https://github.com/home-assistant/core/tree/dev/homeassistant/components/lutron_caseta"
  image_url: ""
  fcc_id: "2ABFCL-BDG2-WH"

# --- TAGS ---
tags: ["bridge", "lutron", "leap", "clear_connect", "professional_grade", "caseta", "dimmer"]
---

# Lutron Caseta Smart Bridge

## What It Is

The Lutron Caseta Smart Bridge is the central hub for Lutron's Caseta wireless lighting control system. It bridges the home network (Ethernet) to Lutron's proprietary Clear Connect RF protocol, which communicates with Caseta in-wall dimmers, switches, fan controllers, Pico remote controls, and Serena motorized shades. Caseta is one of the most reliable smart lighting systems on the market, used extensively by electricians and professional installers. The system comes in two bridge variants: the Standard Bridge (L-BDG2-WH, roughly $80) which supports the LEAP protocol over TLS, and the Smart Bridge Pro (L-BDGPRO2-WH, roughly $100-140) which additionally supports the legacy LIPII (Lutron Integration Protocol II) over telnet on port 23. Lutron has also released a third-generation bridge (L-BDG3-WH) with Thread/Matter support.

## How Haus Discovers It

1. **OUI match**: Lutron bridges use Ethernet with MAC prefixes registered to Lutron Electronics (08:F1:EA, 28:9A:4B, 58:47:CA).
2. **mDNS**: The bridge advertises `_leap._tcp` via mDNS. This is the primary and most reliable discovery method. The mDNS service name includes the bridge serial number.
3. **Hostname pattern**: DHCP hostname matches `Lutron-*` followed by a hex string (the serial number).
4. **Port probe**: TLS connection on port 8083 (LEAP protocol). The Smart Bridge Pro also has port 23 open (telnet/LIPII).
5. **Certificate inspection**: The TLS certificate on port 8083 contains Lutron-specific subject fields including the bridge serial number.

## Pairing / Authentication

Caseta bridge pairing is different depending on the bridge model and protocol used.

### LEAP Protocol (Both Bridge Models, Port 8083)

The LEAP protocol uses mutual TLS (client certificate authentication). Obtaining the certificates requires a complex flow through the Lutron cloud:

1. **Generate a CSR (Certificate Signing Request)**: Create an RSA 2048-bit keypair and CSR locally.
2. **Authenticate with Lutron cloud**: Obtain an OAuth access token via the Lutron developer portal or by intercepting the Lutron app's OAuth flow.
3. **Submit CSR to Lutron**: POST the CSR to `https://device-login.lutron.com/api/v1/pair` with the OAuth token and the bridge serial number.
4. **Press button on bridge**: The user must press the small button on the back of the Caseta bridge within 30 seconds.
5. **Retrieve signed certificate**: Poll `https://device-login.lutron.com/api/v1/pair/{pairing_id}` until the signed client certificate and the bridge's CA certificate are returned.
6. **Connect via mutual TLS**: Use the signed client certificate and private key to establish a TLS connection to port 8083. Both client and server certificates must be validated.

This flow is complex but the resulting certificates are long-lived (multi-year).

### LIPII Protocol (Smart Bridge Pro Only, Port 23)

The Smart Bridge Pro exposes a telnet interface on port 23 with the LIPII (Lutron Integration Protocol II). This is dramatically simpler:

1. Connect via TCP to port 23.
2. **Login prompt**: Username `lutron`, password `integration`.
3. Send and receive LIPII text commands.

This is the preferred integration path for the Smart Bridge Pro due to its simplicity.

## API Reference

### LEAP Protocol (Port 8083)

LEAP (Lutron Extensible Application Protocol) communicates over TLS with JSON messages. It is a request/response and event-subscription protocol.

**Connection**: Establish mutual TLS to `{ip}:8083` using the paired certificates. The protocol is line-delimited JSON over the TLS stream (not HTTP).

**Subscribe to all events:**
```json
{"CommuniqueType": "SubscribeRequest", "Header": {"Url": "/device/status/event"}}
```

**Query all devices:**
```json
{"CommuniqueType": "ReadRequest", "Header": {"Url": "/device"}}
```

**Query a zone (device) status:**
```json
{"CommuniqueType": "ReadRequest", "Header": {"Url": "/zone/1/status"}}
```

**Set zone level (dimmer):**
```json
{
  "CommuniqueType": "CreateRequest",
  "Header": {"Url": "/zone/1/commandprocessor"},
  "Body": {
    "Command": {
      "CommandType": "GoToLevel",
      "Parameter": [{"Type": "Level", "Value": 75}]
    }
  }
}
```

**Response/Event format:**
```json
{
  "CommuniqueType": "ReadResponse",
  "Header": {"StatusCode": "200 OK", "Url": "/zone/1/status"},
  "Body": {
    "ZoneStatus": {
      "href": "/zone/1",
      "Level": 75,
      "Zone": {"href": "/zone/1"}
    }
  }
}
```

### LIPII Protocol (Smart Bridge Pro, Port 23)

LIPII is a simple text-based protocol. Commands use the format: `#ACTION,INTEGRATION_ID,ACTION_NUMBER,PARAMETER`

**Monitor all events (after login):**
```
#MONITORING,255,1
```

**Query device level:**
```
?OUTPUT,2,1
```

Response:
```
~OUTPUT,2,1,75.00
```

**Set device level:**
```
#OUTPUT,2,1,75
```

**Set device level with fade time (seconds):**
```
#OUTPUT,2,1,75,5
```

**Button press on Pico remote:**
```
~DEVICE,3,2,3
```
(Device 3, button 2, action 3=press)

### LIPII Action Numbers

| Action | Description |
|--------|-------------|
| 1 | Set/Get level |
| 2 | Start raising |
| 3 | Start lowering |
| 4 | Stop raising/lowering |
| 6 | Flash zone |

### LIPII Button Actions

| Action | Description |
|--------|-------------|
| 2 | Press |
| 3 | Release |
| 4 | Hold |

### Device Types in Caseta

| Device | Capabilities |
|--------|-------------|
| PD-6WCL | In-wall dimmer (0-100%) |
| PD-5WS | In-wall switch (on/off) |
| PD-FSQN | Fan speed controller (off/low/medium/medium-high/high) |
| PJ2 | Pico remote (button events only) |
| Serena | Motorized shade (0-100% open) |

## AI Capabilities

Not yet planned for initial implementation, but when added the AI could control dimmers ("dim the living room to 50%"), toggle switches, set fan speeds, and monitor Pico remote button presses for automation triggers.

## Quirks & Notes

- **Two bridge models, two protocols**: The Standard Bridge only supports LEAP (TLS, complex pairing). The Smart Bridge Pro supports both LEAP and LIPII (telnet, trivial pairing). For Haus, the Smart Bridge Pro with LIPII is strongly preferred.
- **Clear Connect RF**: Lutron uses its own 434 MHz RF protocol (Clear Connect), not WiFi, Zigbee, or Z-Wave. This means Caseta devices have near-zero WiFi impact and exceptional reliability. The tradeoff is that only Lutron bridges can communicate with Caseta devices.
- **75 device limit**: The Caseta bridge supports a maximum of 75 devices (dimmers, switches, Pico remotes all count).
- **Event-driven**: Both LEAP and LIPII support real-time event subscription, eliminating the need for polling. When a user physically toggles a dimmer or presses a Pico button, the bridge immediately sends an event.
- **LEAP is not HTTP**: Despite using JSON and URL-like paths, LEAP is NOT HTTP REST. It is a persistent TLS connection with line-delimited JSON messages. Do not attempt to use `net/http`.
- **Certificate storage**: The LEAP client certificate and CA certificate must be stored securely. They are difficult to re-obtain if lost (requires repeating the cloud OAuth flow).
- **Third-generation bridge**: Lutron's newest bridge (L-BDG3-WH) adds Thread Border Router and Matter support. LEAP remains the local API, but Matter-over-Thread opens an additional integration path.
- **Fan speed zones**: Fan controllers use levels 0/25/50/75/100 mapping to off/low/medium/medium-high/high. Intermediate values are rounded.
- **Pico remotes are input-only**: Pico remotes generate button events but cannot be "controlled." They are useful for triggering Haus automations.
- **Default telnet credentials**: The LIPII telnet login (lutron/integration) is the same on every Smart Bridge Pro. There is no way to change it. This is a known security concern but is by Lutron's design for integrator access.

## Similar Devices

- [philips-hue-bridge](philips-hue-bridge.md) — Similar bridge-based approach, different RF protocol (Zigbee vs Clear Connect)
- [ikea-dirigera-hub](ikea-dirigera-hub.md) — Another bridge/hub with local API
- [ikea-tradfri-gateway](ikea-tradfri-gateway.md) — Bridge-based, different protocol (CoAP)
