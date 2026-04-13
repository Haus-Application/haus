---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "govee-rgbic-led-strip"
name: "Govee RGBIC LED Strip H6167"
manufacturer: "Govee"
brand: "Govee"
model: "H6167"
model_aliases: ["H6167", "H61A0", "H6163", "H6154", "Govee RGBIC WiFi LED Strip"]
device_type: "govee_led_strip"
category: "lighting"
product_line: "Govee RGBIC"
release_year: 2021
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
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["A4:C1:38", "D4:AD:FC", "FA:21:4B"]
  mdns_services: ["_govee._tcp"]
  mdns_txt_keys: []
  default_ports: [4003]
  signature_ports: [4003]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["Govee.*", "ihoment.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "govee"
  polling_interval_sec: 5
  websocket_event: "govee:state"
  setup_type: "none"
  ai_chattable: true
  haus_milestone: "M6"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp", "scenes"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 4003
  transport: "UDP"
  encoding: "JSON"
  auth_method: "none"
  auth_detail: "LAN API requires no authentication; device must have LAN control enabled in Govee app settings"
  base_url_template: "udp://{ip}:4003"
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "strip"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.govee.com/products/rgbic-led-strip-lights"
  api_docs: "https://developer.govee.com/reference/lan-api"
  developer_portal: "https://developer.govee.com/"
  support: "https://www.govee.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AQA6-H6167"

# --- TAGS ---
tags: ["wifi", "bluetooth_le", "led_strip", "rgbic", "addressable", "lan_api", "udp", "budget"]
---

# Govee RGBIC LED Strip H6167

## What It Is

The Govee RGBIC LED Strip (model H6167) is an addressable LED light strip that supports individual segment color control (RGBIC = RGB with Independent Control). Unlike standard RGB strips that display one color at a time, RGBIC strips have an IC chip that allows different sections to show different colors simultaneously, enabling gradient effects, rainbow patterns, and dynamic scenes. The strip connects via WiFi and BLE, with setup done through the Govee mobile app. Critically for Haus, Govee offers a documented LAN API that allows local control over UDP without requiring cloud connectivity, making it a strong budget-friendly candidate for integration.

## How Haus Discovers It

1. **UDP Discovery Broadcast**: Govee LAN API devices respond to a discovery message. Send a JSON message to the multicast address `239.255.255.250` on UDP port 4001:
   ```json
   {"msg": {"cmd": "scan", "data": {"account_topic": "reserve"}}}
   ```
   Devices respond on UDP port 4002 with their device ID, model, IP address, and supported commands.

2. **mDNS Discovery**: Some Govee devices advertise `_govee._tcp.local.` via mDNS (firmware-dependent).

3. **OUI Match**: MAC addresses beginning with `A4:C1:38`, `D4:AD:FC` are associated with Govee devices (via Espressif/Tuya chipsets — some overlap with other brands).

4. **Port Probe**: UDP port 4003 is the Govee LAN API control port. Sending a `devStatus` command and receiving a valid JSON response confirms the device.

## Pairing / Authentication

No pairing or authentication is required for the Govee LAN API. However, there is an important prerequisite:

1. **Initial Setup**: The strip must first be set up via the Govee mobile app (WiFi provisioning via BLE).
2. **Enable LAN Control**: In the Govee app, navigate to the device settings and enable "LAN Control." This setting is off by default and must be explicitly turned on by the user.
3. Once LAN control is enabled, any device on the local network can send UDP commands without any authentication.

For the Govee cloud API (not used by Haus), an API key is obtained from the Govee Developer Portal.

## API Reference

### Govee LAN API Protocol

The LAN API operates over UDP with JSON-formatted messages.

#### Discovery

**Multicast Scan** (send to `239.255.255.250:4001`):

```json
{
  "msg": {
    "cmd": "scan",
    "data": {
      "account_topic": "reserve"
    }
  }
}
```

**Response** (received on UDP port 4002):

```json
{
  "msg": {
    "cmd": "scan",
    "data": {
      "ip": "192.168.1.100",
      "device": "AB:CD:EF:12:34:56",
      "sku": "H6167",
      "bleVersionHard": "3.01.01",
      "bleVersionSoft": "1.03.01",
      "wifiVersionHard": "1.00.01",
      "wifiVersionSoft": "1.02.03"
    }
  }
}
```

#### Control Commands (send to device IP on UDP port 4003)

**Turn On/Off**:

```json
{
  "msg": {
    "cmd": "turn",
    "data": {
      "value": 1
    }
  }
}
```

Value: `1` = on, `0` = off.

**Set Brightness**:

```json
{
  "msg": {
    "cmd": "brightness",
    "data": {
      "value": 80
    }
  }
}
```

Value: `0` to `100`.

**Set Color (RGB)**:

```json
{
  "msg": {
    "cmd": "colorwc",
    "data": {
      "color": {
        "r": 255,
        "g": 100,
        "b": 0
      },
      "colorTemInKelvin": 0
    }
  }
}
```

RGB values 0-255. Set `colorTemInKelvin` to 0 when using RGB color.

**Set Color Temperature**:

```json
{
  "msg": {
    "cmd": "colorwc",
    "data": {
      "color": {
        "r": 0,
        "g": 0,
        "b": 0
      },
      "colorTemInKelvin": 4000
    }
  }
}
```

Color temperature range: `2000` to `9000` Kelvin. Set RGB to 0 when using color temperature.

**Get Device Status**:

```json
{
  "msg": {
    "cmd": "devStatus",
    "data": {}
  }
}
```

**Status Response**:

```json
{
  "msg": {
    "cmd": "devStatus",
    "data": {
      "onOff": 1,
      "brightness": 80,
      "color": {
        "r": 255,
        "g": 100,
        "b": 0
      },
      "colorTemInKelvin": 0
    }
  }
}
```

### Govee Cloud HTTP API (Not used by Haus — for reference)

Base URL: `https://developer-api.govee.com/v1/`

Authentication: `Govee-API-Key: {api_key}`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/devices` | List all devices |
| PUT | `/devices/control` | Control a device |
| GET | `/devices/state` | Get device state |

## AI Capabilities

When the AI concierge is chatting with a Govee strip, it can:
- Turn the strip on/off
- Set brightness as a percentage
- Set a solid RGB color or color temperature
- Report current state (power, brightness, color, color temperature)
- Suggest ambient lighting configurations based on room activity

Note: The LAN API does not currently expose per-segment control or scene/effect activation. Those features are only available via BLE or the Govee cloud API.

## Quirks & Notes

- **LAN Control Must Be Enabled**: The LAN API is disabled by default. Users must enable it in the Govee app settings. Haus should detect the device on the network and prompt the user to enable LAN control if the device is not responding to UDP commands.
- **No Per-Segment Control via LAN**: The RGBIC independent color control (different colors per segment) is only available via BLE commands, not the LAN API. The LAN API sets a single color for the entire strip.
- **No Scene/Effect Control via LAN**: Dynamic scenes and music visualization modes cannot be activated via the LAN API. Only basic on/off, brightness, and color commands are supported locally.
- **UDP Reliability**: Like LIFX, the protocol is UDP-based, so messages can be lost. Implement retry logic for critical commands.
- **Discovery Multicast**: The scan command uses multicast `239.255.255.250:4001` — the same multicast group as SSDP but a different port. Some networks may block or throttle multicast.
- **Model Variations**: Govee has many RGBIC strip models (H6167, H61A0, H6163, H6154, etc.). The LAN API protocol is the same across supported models, but not all models support LAN control. Check the Govee developer portal for the list of LAN-supported devices.
- **Firmware Updates Required**: Older firmware versions may not support the LAN API at all. Ensure the device firmware is up to date via the Govee app.
- **Espressif Chipset**: Govee devices typically use Espressif ESP32 chips. The MAC OUI may show up as Espressif rather than Govee.

## Similar Devices

- **lifx-a19-color** — WiFi bulb with local UDP protocol, similar no-auth local control paradigm
- **nanoleaf-shapes** — WiFi panels with local REST API, more capable local API
- **wyze-bulb-color** — Budget brand but cloud-only, no local control (anti-pattern)
- **twinkly-smart-lights** — WiFi LED lights with local REST API
