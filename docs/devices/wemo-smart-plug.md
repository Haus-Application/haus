---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "wemo-smart-plug"
name: "Belkin Wemo Smart Plug"
manufacturer: "Belkin International, Inc."
brand: "Wemo"
model: "F7C063"
model_aliases: ["F7C063fc", "WSP080", "Wemo Mini", "F7C063-RM", "WeMo Insight F7C029"]
device_type: "wemo_plug"
category: "smart_home"
product_line: "Wemo"
release_year: 2016
discontinued: true
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
    - "94:10:3E"        # Belkin International (primary OUI)
    - "C4:41:1E"        # Belkin International
    - "EC:1A:59"        # Belkin International
    - "08:86:3B"        # Belkin International
    - "58:EF:68"        # Belkin International (newer production)
    - "B4:75:0E"        # Belkin International
  mdns_services: []     # Wemo does NOT use mDNS; relies entirely on SSDP/UPnP
  mdns_txt_keys: []
  default_ports: [49152, 49153, 49154, 49155]
  signature_ports: [49153]
  ssdp_search_target: "urn:Belkin:device:controllee:1"
  ssdp_server_string: "Unspecified, UPnP/1.0, Unspecified"
  hostname_patterns:
    - "^WeMo\\."
    - "^wemo"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 49153
    path: "/setup.xml"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: "Unspecified, UPnP/1.0, Unspecified"
    body_contains: "<deviceType>urn:Belkin:device:controllee:1</deviceType>"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "wemo"
  polling_interval_sec: 10
  websocket_event: "wemo:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M5"

# --- CAPABILITIES ---
capabilities:
  - "on_off"

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 49153
  transport: "HTTP"
  encoding: "XML"
  auth_method: "none"
  auth_detail: "No authentication on LAN. SOAP actions sent to /upnp/control/basicevent1 with SOAPAction header."
  base_url_template: "http://{ip}:{port}/upnp/control/basicevent1"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "plug"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.belkin.com/support-article/?articleNum=226142"
  api_docs: ""
  developer_portal: ""
  support: "https://www.belkin.com/support/"
  community_forum: ""
  image_url: ""
  fcc_id: "K7SF7C063"

# --- TAGS ---
tags: ["wifi", "upnp", "ssdp", "soap", "no_hub", "plug", "belkin", "no_auth_lan"]
---

# Belkin Wemo Smart Plug

## What It Is

The Belkin Wemo Smart Plug (model F7C063) is a WiFi-connected smart outlet adapter that plugs into a standard wall receptacle and provides on/off control of any connected appliance. It was one of the earliest mainstream smart plugs, predating many of today's smart home ecosystems. The Wemo line uses UPnP/SSDP for device discovery and a SOAP-over-HTTP local API for control, meaning it can be fully operated on the local network without internet access once initially set up through the Wemo app. Belkin has largely shifted the Wemo brand toward Apple HomeKit and Thread-based products, and the original F7C063 is discontinued, but millions remain in homes and the local SOAP API continues to function on existing firmware. The newer WSP080 (Wemo WiFi Smart Plug) uses a similar but updated protocol.

## How Haus Discovers It

1. **OUI Match** -- During network scan, any device with MAC prefix `94:10:3E`, `C4:41:1E`, `EC:1A:59`, `08:86:3B`, `58:EF:68`, or `B4:75:0E` is flagged as a potential Belkin device.
2. **SSDP Discovery** -- The Wemo plug responds to UPnP M-SEARCH multicast messages. Haus sends an M-SEARCH to `239.255.255.250:1900` with the search target `urn:Belkin:device:controllee:1` (for basic plugs) or `urn:Belkin:device:**` for broader Wemo discovery. The plug responds with its `LOCATION` header pointing to the device description XML, typically at `http://{ip}:{port}/setup.xml`.
3. **Setup XML Fingerprint** -- Fetching the setup.xml URL returns UPnP device description XML containing `<deviceType>urn:Belkin:device:controllee:1</deviceType>`, `<manufacturer>Belkin International Inc.</manufacturer>`, the `<friendlyName>` (user-assigned), `<modelName>Socket</modelName>`, `<modelNumber>1.0</modelNumber>`, and the `<serialNumber>` and `<UDN>`.
4. **Port Probe** -- Wemo devices listen on a dynamically assigned port in the 49152-49155 range (ephemeral UPnP ports). Port 49153 is most common but not guaranteed. The setup.xml location from SSDP provides the authoritative port.
5. **SOAP Probe** -- Sending a `GetBinaryState` SOAP request to `/upnp/control/basicevent1` and receiving a valid XML response confirms the device is a controllable Wemo plug.

## Pairing / Authentication

No pairing or authentication is required for local LAN control. The Wemo SOAP API has no auth mechanism — any device on the same network can send SOAP commands. Initial WiFi setup must be done through the Wemo mobile app, which puts the plug into a SoftAP mode, connects the phone to it, and provisions WiFi credentials. Once on the home network, the plug is immediately controllable locally.

## API Reference

The Wemo local API uses SOAP (Simple Object Access Protocol) over HTTP. All commands are sent as POST requests to the basicevent1 control URL.

### Endpoint

```
POST http://{ip}:{port}/upnp/control/basicevent1
```

The port is obtained from the SSDP LOCATION header (commonly 49153 but not guaranteed).

### Headers

All SOAP requests require:
```
Content-Type: text/xml; charset="utf-8"
SOAPACTION: "urn:Belkin:service:basicevent:1#{ActionName}"
```

### Get Binary State

Returns the current on/off state of the plug.

**SOAPAction:** `"urn:Belkin:service:basicevent:1#GetBinaryState"`

```xml
<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
            s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:GetBinaryState xmlns:u="urn:Belkin:service:basicevent:1">
    </u:GetBinaryState>
  </s:Body>
</s:Envelope>
```

**Response:**
```xml
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
            s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:GetBinaryStateResponse xmlns:u="urn:Belkin:service:basicevent:1">
      <BinaryState>1</BinaryState>
    </u:GetBinaryStateResponse>
  </s:Body>
</s:Envelope>
```

- `<BinaryState>1</BinaryState>` = ON
- `<BinaryState>0</BinaryState>` = OFF
- `<BinaryState>8</BinaryState>` = ON (standby/idle, seen on Wemo Insight)

### Set Binary State

Turns the plug on or off.

**SOAPAction:** `"urn:Belkin:service:basicevent:1#SetBinaryState"`

```xml
<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
            s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:SetBinaryState xmlns:u="urn:Belkin:service:basicevent:1">
      <BinaryState>1</BinaryState>
    </u:SetBinaryState>
  </s:Body>
</s:Envelope>
```

- Send `<BinaryState>1</BinaryState>` to turn ON
- Send `<BinaryState>0</BinaryState>` to turn OFF

### Get Friendly Name

**SOAPAction:** `"urn:Belkin:service:basicevent:1#GetFriendlyName"`

```xml
<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
            s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:GetFriendlyName xmlns:u="urn:Belkin:service:basicevent:1">
    </u:GetFriendlyName>
  </s:Body>
</s:Envelope>
```

### SSDP M-SEARCH

```
M-SEARCH * HTTP/1.1
HOST: 239.255.255.250:1900
MAN: "ssdp:discover"
MX: 5
ST: urn:Belkin:device:controllee:1
```

The plug responds with a unicast HTTP response containing the `LOCATION` header pointing to setup.xml.

### UPnP Event Subscription

Wemo supports UPnP eventing for state change notifications. Haus can subscribe to the basicevent1 service to receive callbacks when the plug state changes:

```
SUBSCRIBE /upnp/event/basicevent1 HTTP/1.1
HOST: {ip}:{port}
CALLBACK: <http://{haus_ip}:{callback_port}/>
NT: upnp:event
TIMEOUT: Second-600
```

The plug will POST event notifications to the callback URL when state changes.

## AI Capabilities

When the AI concierge is chatting with a Wemo plug, it can:

- **Toggle the plug** on or off
- **Report current state** -- "I'm currently ON" or "I'm currently OFF"
- **Identify the device** by its friendly name
- **Explain what it controls** -- the user names the plug during setup (e.g., "Coffee Maker", "Desk Lamp")

## Quirks & Notes

- **Dynamic Ports:** The Wemo SOAP API port is not fixed. It is typically 49152-49155 but can change after a reboot or firmware update. Always discover the port via SSDP rather than hardcoding it.
- **No mDNS:** Unlike many modern smart home devices, Wemo does not advertise via mDNS. Discovery relies entirely on SSDP/UPnP.
- **Slow SSDP Response:** Wemo devices can be slow to respond to M-SEARCH, sometimes taking 3-5 seconds. Use an MX value of at least 3.
- **Firmware Fragmentation:** Older Wemo firmware has known security vulnerabilities. Belkin released patches but many devices in the field remain unpatched. The SOAP API has remained stable across firmware versions.
- **BinaryState=8:** On Wemo Insight models (F7C029), a BinaryState of 8 means the device is on but in standby/idle (low power draw). Treat as ON for control purposes.
- **Subscription Expiry:** UPnP event subscriptions expire after the TIMEOUT period (default 600 seconds). Haus must re-subscribe periodically.
- **No TLS:** The local API is plain HTTP with no encryption. All SOAP commands are sent in cleartext on the LAN.
- **Wemo vs. Wemo Mini vs. WSP080:** The F7C063 is the original plug. The F7C063fc "Mini" uses the same SOAP API. The newer WSP080 also supports SOAP but additionally supports Apple HomeKit via HAP over WiFi. Haus should target the SOAP API for broadest compatibility.
- **XML Parsing:** The SOAP responses use standard XML namespaces. Haus should use a proper XML parser and not rely on regex for response parsing, as whitespace and namespace prefixes can vary between firmware versions.
- **Discontinued but Prevalent:** Belkin has largely moved on from WiFi SOAP-based Wemo products toward Thread/Matter, but millions of original Wemo plugs remain in use. This makes integration worthwhile despite the product being discontinued.

## Similar Devices

- **wemo-insight-plug** -- Wemo Insight (F7C029) with energy monitoring, same SOAP API plus GetInsightParams
- **wemo-light-switch** -- Wemo Light Switch (F7C030), in-wall switch variant with same SOAP protocol
- **kasa-smart-plug** -- TP-Link Kasa plug, different protocol (XOR-JSON on TCP 9999) but similar functionality
- **amazon-smart-plug** -- Amazon Smart Plug, cloud-only with no local API
- **meross-smart-plug-mss110** -- Meross plug with local HTTP+MQTT API
