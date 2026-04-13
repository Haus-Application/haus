---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "amazon-smart-plug"
name: "Amazon Smart Plug"
manufacturer: "Amazon.com Services LLC"
brand: "Amazon"
model: "B01MZEEFNX"
model_aliases: ["B089DR29T6", "Amazon Smart Plug (Works with Alexa)", "SP-A1"]
device_type: "cloud_plug"
category: "smart_home"
product_line: "Amazon Smart Home"
release_year: 2018
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "cloud"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: true
  local_only_capable: false
  protocols_spoken: ["wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "FC:65:DE"        # Amazon Technologies Inc.
    - "40:B4:CD"        # Amazon Technologies Inc.
    - "A0:02:DC"        # Amazon Technologies Inc.
    - "F0:F0:A4"        # Amazon Technologies Inc.
    - "F0:27:2D"        # Amazon Technologies Inc.
    - "74:C2:46"        # Amazon Technologies Inc.
    - "68:54:FD"        # Amazon Technologies Inc.
    - "44:00:49"        # Amazon Technologies Inc.
    - "AC:63:BE"        # Amazon Technologies Inc.
  mdns_services: []     # No mDNS service advertisement for local control
  mdns_txt_keys: []
  default_ports: []     # No open ports for local API
  signature_ports: [8443]   # TLS connection to Amazon cloud (outbound only)
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^amazon-[0-9a-f]{12}$"
    - "^[0-9a-f]{12}\\.amazon$"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP interface

# --- HAUS INTEGRATION ---
integration:
  status: "not_feasible"
  integration_key: ""
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "on_off"

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "HTTPS"
  encoding: "Protobuf"
  auth_method: "oauth2"
  auth_detail: "Communication exclusively through Amazon's Alexa cloud infrastructure. Uses proprietary TLS protocol to Amazon endpoints. No local API, no local control path. Alexa Smart Home Skill API requires Amazon Developer account and OAuth2 linked to user's Alexa account."
  base_url_template: ""
  tls: true
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
  product_page: "https://www.amazon.com/Amazon-Smart-Plug-works-Alexa/dp/B01MZEEFNX"
  api_docs: "https://developer.amazon.com/en-US/docs/alexa/smarthome/understand-the-smart-home-skill-api.html"
  developer_portal: "https://developer.amazon.com/alexa"
  support: "https://www.amazon.com/gp/help/customer/display.html"
  community_forum: ""
  image_url: ""
  fcc_id: "2AQHT-4252"

# --- TAGS ---
tags: ["wifi", "cloud_only", "alexa", "no_local_api", "plug", "amazon", "not_feasible"]
---

# Amazon Smart Plug

## What It Is

The Amazon Smart Plug is an inexpensive WiFi smart plug sold by Amazon, designed exclusively for the Alexa ecosystem. It plugs into a standard wall outlet and provides on/off control for any connected appliance up to 15A. The plug is set up through the Alexa app and communicates solely through Amazon's cloud infrastructure. Unlike virtually every other smart plug on the market, the Amazon Smart Plug has absolutely no local API, no local control mechanism, and no way to interact with it without going through Amazon's servers. It is a pure cloud device -- if your internet goes down or Amazon's services experience an outage, the plug becomes uncontrollable (though it retains its last state). Amazon sells it at or below cost as an Alexa ecosystem lock-in product, often bundled with Echo devices. For Haus, this device is classified as "not_feasible" because there is no viable local control path.

## How Haus Discovers It

Haus can detect the Amazon Smart Plug on the network but cannot control it:

1. **OUI Match** -- MAC addresses beginning with Amazon OUIs (`FC:65:DE`, `40:B4:CD`, `A0:02:DC`, `F0:F0:A4`, `F0:27:2D`, `74:C2:46`, `68:54:FD`, `44:00:49`, `AC:63:BE`) indicate an Amazon device. However, these OUIs are shared across all Amazon devices (Echo, Fire TV, Kindle, Ring, etc.), so further identification is needed.
2. **DHCP Hostname** -- The plug may register with a hostname pattern like `amazon-xxxxxxxxxxxx`.
3. **Network Traffic Analysis** -- The plug maintains persistent TLS connections to Amazon endpoints (typically `*.amazon.com` or `*.amazonaws.com` on port 8443 or 443). Observing these connections can confirm it is an Amazon IoT device.
4. **No Open Ports** -- Unlike smart plugs with local APIs, the Amazon Smart Plug has no open TCP or UDP ports accepting inbound connections. A port scan returns nothing useful. This absence itself is a signal.
5. **Haus Classification** -- When Haus detects an Amazon-OUI device with no open ports and no mDNS/SSDP advertisements, it can flag it as a probable Amazon cloud-only device and inform the user that local control is not possible.

## Pairing / Authentication

### Alexa App Setup (Only Method)

1. Open the Alexa app on a mobile device.
2. Go to Devices > Add Device > Plug > Amazon Smart Plug.
3. The app puts the plug into setup mode (solid orange light).
4. The phone connects to the plug's temporary WiFi AP to provision WiFi credentials.
5. The plug connects to the home network and registers with Amazon's cloud.
6. The plug appears in the Alexa app and can be controlled via voice or the app.

### Why Local Control is Not Possible

- The plug does not expose any ports on the local network.
- It does not respond to UPnP/SSDP discovery.
- It does not advertise via mDNS.
- The device communicates with Amazon's infrastructure over TLS using a certificate-pinned connection with proprietary protocol framing.
- Amazon has not published and does not provide a local API.
- The firmware is signed and encrypted, preventing community firmware replacement (unlike ESP-based devices that can be flashed with Tasmota or ESPHome).
- The plug uses a custom Amazon chipset, not a standard ESP8266/ESP32, making hardware-level exploits impractical.

## API Reference

### Alexa Smart Home Skill API (Cloud Only)

The only programmatic way to control the Amazon Smart Plug is through the Alexa Smart Home Skill API, which requires:

1. An Amazon Developer account.
2. A registered Alexa Smart Home Skill.
3. OAuth2 account linking between the Alexa user's account and the skill.
4. An AWS Lambda function to handle skill directives.

This is a cloud-to-cloud integration path designed for third-party skill developers, not for local home automation systems.

#### Power Control Interface

```json
{
  "directive": {
    "header": {
      "namespace": "Alexa.PowerController",
      "name": "TurnOn",
      "payloadVersion": "3",
      "messageId": "..."
    },
    "endpoint": {
      "endpointId": "amazon-smart-plug-xxxx"
    },
    "payload": {}
  }
}
```

Directives: `TurnOn`, `TurnOff`

This requires the full Alexa skill chain and is not practical for Haus integration.

## AI Capabilities

The AI concierge cannot chat with or control the Amazon Smart Plug through Haus. If a user asks about it, the AI should:

- **Explain the limitation** -- "This Amazon Smart Plug can only be controlled through Alexa. It has no local API and requires Amazon's cloud."
- **Suggest alternatives** -- Recommend locally-controllable plugs like Kasa, Shelly, or Meross.
- **Detect and inform** -- If Haus discovers an Amazon device on the network, inform the user it is not locally controllable.

## Quirks & Notes

- **No Local API Whatsoever:** This cannot be overstated. The Amazon Smart Plug is a completely closed, cloud-dependent device. There is no local HTTP API, no UPnP, no mDNS, no MQTT, no CoAP, no nothing. It is a brick without internet.
- **Custom Silicon:** Unlike many budget smart plugs that use ESP8266/ESP32 chips (which can be reflashed with open-source firmware), the Amazon Smart Plug uses Amazon's custom chipset. This prevents flashing with Tasmota, ESPHome, or other community firmware.
- **Loss Leader Pricing:** Amazon frequently sells the plug for $5-15 USD, well below manufacturing cost. The purpose is Alexa ecosystem lock-in, not profit from the hardware.
- **Shared OUIs:** Amazon's MAC address OUIs are shared across ALL Amazon devices. A MAC starting with `FC:65:DE` could be an Echo Dot, Fire TV Stick, Ring camera, or this plug. OUI alone is insufficient for identification.
- **Matter Conspicuously Absent:** As of 2026, Amazon has not added Matter support to the Amazon Smart Plug. This is likely intentional to maintain ecosystem lock-in.
- **Routine Integration:** The plug works with Alexa Routines, which can be triggered by time, other Alexa devices, or voice. This is the only "automation" capability, and it runs entirely in Amazon's cloud.
- **Power Outage Behavior:** After a power outage, the plug restores its last known state once internet connectivity is re-established. During the period between power restoration and cloud reconnection, the plug's state is unpredictable.
- **Energy Monitoring:** The Amazon Smart Plug does NOT have energy monitoring. It is purely on/off.
- **15A Rating:** Standard US outlet rating. Suitable for lamps, fans, coffee makers, but not high-draw appliances like space heaters that approach the 15A limit.

## Similar Devices

- **kasa-smart-plug** -- TP-Link Kasa plug with excellent local XOR protocol, same price range
- **wemo-smart-plug** -- Belkin Wemo with local SOAP/UPnP API
- **meross-smart-plug-mss110** -- Meross plug with local HTTP API
- **shelly-plug-s** -- Shelly plug with best-in-class local RPC API
- **wyze-plug** -- Wyze plug, mostly cloud-dependent but some community local control efforts
