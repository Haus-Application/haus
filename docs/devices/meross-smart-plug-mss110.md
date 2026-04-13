---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "meross-smart-plug-mss110"
name: "Meross Smart WiFi Plug MSS110"
manufacturer: "Chengdu Meross Technology Co., Ltd."
brand: "Meross"
model: "MSS110"
model_aliases: ["MSS110HK", "MSS110US", "MSS110-1", "MSS110-2"]
device_type: "meross_plug"
category: "smart_home"
product_line: "Meross Smart"
release_year: 2018
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
    - "48:E1:E9"        # Meross / Chengdu Meross Technology
    - "34:29:8F"        # Meross devices (newer production)
  mdns_services: []     # Meross does not use mDNS for discovery
  mdns_txt_keys: []
  default_ports: [80]
  signature_ports: [80]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Meross"
    - "^meross"
    - "^mss110"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/config"
    method: "POST"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "\"method\""
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "meross"
  polling_interval_sec: 10
  websocket_event: "meross:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M6"

# --- CAPABILITIES ---
capabilities:
  - "on_off"

# --- PROTOCOL ---
protocol:
  type: "http_rest"
  port: 80
  transport: "HTTP"
  encoding: "JSON"
  auth_method: "api_key"
  auth_detail: "Local HTTP API uses a messageId, timestamp, and sign (MD5 hash of messageId + key + timestamp). The key is derived from cloud account pairing. Community libraries have reverse-engineered the signing process."
  base_url_template: "http://{ip}/config"
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
  product_page: "https://www.meross.com/en-gc/smart-plug/smart-wifi-plug-mini/29"
  api_docs: ""
  developer_portal: ""
  support: "https://www.meross.com/en-gc/support"
  community_forum: "https://community.home-assistant.io/t/meross-lan-local-integration/367514"
  image_url: ""
  fcc_id: "2AK5Z-MSS110"

# --- TAGS ---
tags: ["wifi", "no_hub", "plug", "meross", "local_http", "mqtt", "json_signing"]
---

# Meross Smart WiFi Plug MSS110

## What It Is

The Meross Smart WiFi Plug MSS110 is an affordable WiFi smart plug from Chengdu Meross Technology (also operating under the Refoss brand). It connects directly to a home WiFi network (2.4GHz) and provides on/off control for any connected appliance up to 15A. Meross plugs are popular in the Home Assistant community because, unlike many budget smart plugs, they have a documented local HTTP API that allows control without cloud dependency. The MSS110 is the basic single-outlet model. Meross also makes the MSS110HK variant with Apple HomeKit support, the MSS210 with energy monitoring, and the MSS425 power strip. All share the same local protocol. The device communicates with Meross cloud servers via MQTT (port 443) for the mobile app, but the local HTTP API on port 80 enables fully local control once the signing key is obtained.

## How Haus Discovers It

1. **OUI Match** -- MAC addresses beginning with `48:E1:E9` or `34:29:8F` are associated with Meross devices.
2. **Port Probe** -- The MSS110 listens on port 80/TCP for local HTTP API requests. A connection to port 80 that accepts POST requests is a candidate.
3. **HTTP Fingerprint** -- Sending a properly formatted JSON POST to `http://{ip}/config` with a valid Meross protocol message header and the `Appliance.System.All` method returns a JSON response containing device info including `hardware.type` (e.g., "mss110"), `hardware.version`, `firmware.version`, and `hardware.macAddress`.
4. **DHCP Hostname** -- Meross devices sometimes register with hostnames starting with "Meross" or the model number.
5. **Cloud MQTT Broker** -- Meross devices maintain persistent MQTT connections to Meross cloud servers (e.g., `iot.meross.com:443` or regional variants). Haus does not use this path but observing outbound MQTT traffic to known Meross endpoints can confirm device identity during network analysis.

## Pairing / Authentication

Meross uses a signed message authentication scheme for its local HTTP API. Unlike Wemo or LIFX, you cannot simply send unauthenticated commands.

### Authentication Mechanism

Every local HTTP request must include a JSON header with:
- `messageId` -- a UUID v4 string identifying the request
- `timestamp` -- Unix epoch timestamp in seconds
- `sign` -- an MD5 hash of the concatenation: `messageId + key + timestamp`

The `key` is a shared secret derived from the Meross cloud account. It can be obtained by:

1. **Cloud Login** -- Authenticate with the Meross cloud API at `https://iot.meross.com/v1/Auth/Login` using email/password to retrieve an access token. The token response contains the `key` field.
2. **Community Tools** -- Tools like `meross-iot` (Python library) and the Home Assistant `meross_lan` integration can extract the key.
3. **One-Time Extraction** -- The key is stable and does not rotate unless the user changes their Meross account password.

### Signing Example

```
messageId = "abc123-def456-..."
key = "0123456789abcdef0123456789abcdef"
timestamp = 1712000000

sign = MD5(messageId + key + timestamp)
     = MD5("abc123-def456-..." + "0123456789abcdef0123456789abcdef" + "1712000000")
```

## API Reference

All requests are HTTP POST to `http://{ip}/config` with `Content-Type: application/json`.

### Request Format

```json
{
  "header": {
    "messageId": "{uuid}",
    "method": "GET",
    "namespace": "Appliance.System.All",
    "timestamp": 1712000000,
    "sign": "{md5_hash}",
    "payloadVersion": 1,
    "from": "http://{ip}/config"
  },
  "payload": {}
}
```

### Get System Info (Appliance.System.All)

**Namespace:** `Appliance.System.All`

Returns complete device state including hardware info, firmware version, WiFi signal strength, and toggle state.

**Response payload:**
```json
{
  "all": {
    "system": {
      "hardware": {
        "type": "mss110",
        "subType": "us",
        "version": "2.0.0",
        "chipType": "mt7682",
        "uuid": "...",
        "macAddress": "48:e1:e9:xx:xx:xx"
      },
      "firmware": {
        "version": "2.1.16",
        "compileTime": "..."
      },
      "online": {
        "status": 1
      }
    },
    "digest": {
      "togglex": [
        {"channel": 0, "onoff": 1, "lmTime": 1712000000}
      ]
    }
  }
}
```

### Toggle On/Off (Appliance.Control.ToggleX)

**Namespace:** `Appliance.Control.ToggleX`

```json
{
  "header": {
    "messageId": "{uuid}",
    "method": "SET",
    "namespace": "Appliance.Control.ToggleX",
    "timestamp": 1712000000,
    "sign": "{md5_hash}",
    "payloadVersion": 1,
    "from": ""
  },
  "payload": {
    "togglex": {
      "channel": 0,
      "onoff": 1
    }
  }
}
```

- `onoff: 1` = ON
- `onoff: 0` = OFF
- `channel: 0` = primary outlet (multi-outlet models like MSS425 use channels 0-3)

### Get Toggle State (Appliance.Control.ToggleX)

Same namespace with `"method": "GET"` and empty payload. Returns current toggle state for all channels.

### WiFi Info (Appliance.System.DNDMode)

**Namespace:** `Appliance.System.DNDMode` with `"method": "GET"` returns the Do Not Disturb mode (LED on/off).

### Firmware Info (Appliance.System.Firmware)

**Namespace:** `Appliance.System.Firmware` returns current firmware details and available update info.

## AI Capabilities

When the AI concierge is chatting with a Meross plug, it can:

- **Toggle the plug** on or off
- **Report current state** -- on, off, which channel
- **Report device info** -- hardware version, firmware version, WiFi signal
- **Identify the device** by its configured name and model

## Quirks & Notes

- **Cloud Setup Required:** Initial WiFi provisioning must be done through the Meross app, which creates a cloud account and pairs the device. The local API key is derived from this cloud account. There is no way to provision a brand-new Meross device locally.
- **Key Extraction:** The `key` needed for local API signing must be obtained from the Meross cloud login. Once extracted, it can be stored and reused indefinitely without further cloud interaction. The `meross_lan` Home Assistant integration automates this.
- **Firmware Updates May Break Local API:** Meross has occasionally pushed firmware updates that changed or temporarily broke the local HTTP API. The community has generally found workarounds, but this is a risk.
- **MT7682 Chipset:** The MSS110 v2 uses a MediaTek MT7682 chipset. Earlier v1 models used different chips and a slightly different protocol variant.
- **Multi-Channel Devices:** The MSS425 power strip and MSS620 outdoor plug have multiple controllable channels (0-3 or 0-1). The MSS110 is single-channel (channel 0 only). Always specify channel in ToggleX commands.
- **MQTT Cloud Protocol:** The cloud communication uses MQTT over TLS on port 443 (not standard MQTT port 1883). The device connects to `iot.meross.com` or regional servers. This is separate from the local HTTP API.
- **Response Time:** Local HTTP API responses are typically under 100ms, making it very responsive for automation. Cloud round-trip via MQTT is 200-500ms.
- **No UPnP/SSDP:** Unlike Wemo, Meross devices do not announce themselves via SSDP. Discovery relies on OUI matching, port probing, and HTTP fingerprinting.
- **HomeKit Variants:** The MSS110HK variant adds Apple HomeKit (HAP over WiFi) support. The HAP protocol runs alongside the Meross local HTTP API on the same device.
- **Energy Monitoring:** The MSS110 does NOT have energy monitoring. The MSS210/MSS310 add wattage, voltage, and current reporting via the `Appliance.Control.Electricity` namespace.

## Similar Devices

- **meross-smart-plug-mss210** -- Same form factor with energy monitoring
- **meross-power-strip-mss425** -- 4-outlet smart power strip, multi-channel variant
- **kasa-smart-plug** -- TP-Link Kasa plug, different protocol but similar use case
- **wemo-smart-plug** -- Belkin Wemo plug, SOAP/UPnP instead of HTTP/JSON
- **shelly-plug-s** -- Shelly plug with excellent local REST API, similar local-first philosophy
