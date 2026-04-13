---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "zooz-zwave-switch-zen76"
name: "Zooz Z-Wave Plus S2 On/Off Switch ZEN76"
manufacturer: "Zooz (The Smartest House)"
brand: "Zooz"
model: "ZEN76"
model_aliases: ["ZEN76 S2", "ZEN76 800", "ZEN76 V1.0", "700 Series ZEN76"]
device_type: "zwave_switch"
category: "lighting"
product_line: "800 Series"
release_year: 2021
discontinued: false
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: false
  cloud_api: false
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["zwave"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []       # Z-Wave devices do not appear on IP network
  mdns_services: []      # Not an IP device
  mdns_txt_keys: []
  default_ports: []      # No IP ports
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "zwave"
  polling_interval_sec: 0
  websocket_event: "zwave:state"
  setup_type: "app_pairing"
  ai_chattable: true
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities:
  - "on_off"

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "Z-Wave"
  encoding: "binary"
  auth_method: "none"
  auth_detail: "Z-Wave S2 (Security 2) authenticated inclusion using DSK (Device Specific Key) printed on device. Supports S2 Unauthenticated, S2 Authenticated, and non-secure inclusion modes."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "switch"
  power_source: "hardwired"
  mounting: "in_wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["zwave"]

# --- LINKS ---
links:
  product_page: "https://www.thesmartesthouse.com/products/zooz-z-wave-plus-s2-on-off-switch-zen76"
  api_docs: ""
  developer_portal: ""
  support: "https://www.support.getzooz.com/"
  community_forum: "https://community.home-assistant.io/"
  image_url: ""
  fcc_id: "2AZJN-ZEN76"

# --- TAGS ---
tags: ["zwave", "in_wall", "switch", "s2_security", "scene_control", "no_hub_no_ip", "zooz", "800_series"]
---

# Zooz Z-Wave Plus S2 On/Off Switch ZEN76

## What It Is

The Zooz ZEN76 is a Z-Wave Plus on/off switch sold by The Smartest House, a US-based Z-Wave specialty retailer that designs the Zooz product line. The ZEN76 is a simple toggle switch (no dimming) in a standard Decora paddle form factor. It uses the Z-Wave 700 series chip (or 800 series in the latest revision), supports Z-Wave S2 security for encrypted communication, and includes scene control via multi-tap paddle actions. It requires a Z-Wave hub/controller (such as a USB Z-Wave dongle with Z-Wave JS, Hubitat, SmartThings, or a Z-Wave-to-Matter bridge) and does not connect to WiFi or the IP network directly. Zooz switches are beloved in the Home Assistant community for their excellent Z-Wave compliance, extensive configuration parameters, competitive pricing, and responsive customer support. The ZEN76 requires a neutral wire and supports single-pole and 3-way configurations.

## How Haus Discovers It

The ZEN76 is a Z-Wave device and is not directly discoverable on the IP network. Haus discovers it through a Z-Wave controller:

1. **Z-Wave Inclusion** -- The switch is put into inclusion mode by tapping the upper paddle 3x quickly. The Z-Wave controller must be in inclusion/add mode simultaneously.
2. **Z-Wave Interview** -- Once included, the controller queries the device's Z-Wave command classes. The ZEN76 reports:
   - Manufacturer ID: `0x027A` (Zooz)
   - Product Type ID: `0x7000`
   - Product ID: `0xA006` (ZEN76)
   - Supported Command Classes:
     - `COMMAND_CLASS_BASIC` (0x20)
     - `COMMAND_CLASS_SWITCH_BINARY` (0x25)
     - `COMMAND_CLASS_ASSOCIATION` (0x85)
     - `COMMAND_CLASS_MULTI_CHANNEL_ASSOCIATION` (0x8E)
     - `COMMAND_CLASS_CONFIGURATION` (0x70)
     - `COMMAND_CLASS_CENTRAL_SCENE` (0x5B)
     - `COMMAND_CLASS_FIRMWARE_UPDATE` (0x7A)
     - `COMMAND_CLASS_MANUFACTURER_SPECIFIC` (0x72)
     - `COMMAND_CLASS_VERSION` (0x86)
     - `COMMAND_CLASS_SECURITY_2` (0x9F)
3. **Z-Wave-to-Matter Bridge (Future)** -- If Haus implements or supports a Z-Wave-to-Matter bridge, the ZEN76 could be exposed as a Matter device. Silicon Labs and other vendors are developing bridge solutions that translate Z-Wave devices into Matter fabric endpoints.

## Pairing / Authentication

### Z-Wave S2 Inclusion (Preferred)

1. Put the Z-Wave controller into inclusion mode.
2. Tap the upper paddle 3 times quickly on the ZEN76.
3. The controller prompts for the DSK (Device Specific Key) -- a 5-digit PIN printed on the switch and its packaging.
4. Enter the DSK. The controller and switch negotiate S2 Authenticated inclusion.
5. The switch is now securely included with AES-128 encrypted communication.

### Z-Wave Non-Secure Inclusion

Same process but skip the DSK step. Communication will not be encrypted. Not recommended for any security-sensitive application.

### Exclusion / Factory Reset

1. Put the Z-Wave controller into exclusion mode.
2. Tap the lower paddle 3 times quickly.
3. Alternative factory reset: tap upper paddle 10x rapidly, then immediately hold the upper paddle for 10+ seconds.

## API Reference

### Z-Wave Command Classes

The ZEN76 is controlled via standard Z-Wave command classes through a Z-Wave controller.

#### Switch Binary (0x25)

The primary control command class for on/off:

- **Set** -- `SWITCH_BINARY_SET` with value `0xFF` (ON) or `0x00` (OFF)
- **Get** -- `SWITCH_BINARY_GET` returns current state
- **Report** -- `SWITCH_BINARY_REPORT` sent on state change (value `0xFF` or `0x00`)

#### Central Scene (0x5B)

Reports multi-tap paddle actions for scene control:

| Action | Scene Number | Key Attribute |
|--------|-------------|---------------|
| 1x tap up | 1 | 0 (Key Pressed 1 time) |
| 2x tap up | 1 | 3 (Key Pressed 2 times) |
| 3x tap up | 1 | 4 (Key Pressed 3 times) |
| 4x tap up | 1 | 5 (Key Pressed 4 times) |
| 5x tap up | 1 | 6 (Key Pressed 5 times) |
| 1x tap down | 2 | 0 |
| 2x tap down | 2 | 3 |
| 3x tap down | 2 | 4 |
| 4x tap down | 2 | 5 |
| 5x tap down | 2 | 6 |
| Hold up | 1 | 2 (Key Held Down) |
| Hold down | 2 | 2 |
| Release up | 1 | 1 (Key Released) |
| Release down | 2 | 1 |

#### Configuration (0x70)

Key configuration parameters:

| Parameter | Size | Default | Description |
|-----------|------|---------|-------------|
| 1 | 1 | 0 | LED indicator mode: 0=on when load off, 1=on when load on, 2=always off, 3=always on |
| 2 | 4 | 0 | Auto-off timer (seconds, 0=disabled) |
| 4 | 4 | 0 | Auto-on timer (seconds, 0=disabled) |
| 6 | 1 | 0 | Scene control enable: 0=disabled, 1=enabled (must enable for multi-tap) |
| 7 | 1 | 1 | Smart bulb mode: 0=normal, 1=smart bulb (switch always sends commands but load stays on) |
| 10 | 1 | 0 | Local control disable: 0=enabled, 1=disabled (paddle still sends Z-Wave commands) |
| 12 | 1 | 0 | 3-way switch type: 0=toggle, 1=momentary |
| 13 | 1 | 0 | Report type for local paddle: 0=Binary Switch, 1=Basic |

### Z-Wave Association Groups

| Group | Max Nodes | Description |
|-------|-----------|-------------|
| 1 | 1 | Lifeline (controller) -- sends Binary Switch Reports and Central Scene notifications |
| 2 | 5 | Basic Set -- sends Basic Set ON/OFF to associated devices when paddle is used |

## AI Capabilities

When the AI concierge is chatting with a Zooz ZEN76, it can:

- **Turn the switch on or off**
- **Report current state** -- on or off
- **Configure parameters** -- LED indicator mode, auto-off timer, scene control enable/disable
- **Explain multi-tap actions** -- describe available paddle tap combinations for automations
- **Report device info** -- manufacturer, product ID, firmware version, S2 security status

## Quirks & Notes

- **No IP Network Presence:** Z-Wave operates on sub-GHz radio frequencies (908.42 MHz in North America). The ZEN76 has zero network visibility over WiFi/Ethernet. A Z-Wave USB controller or hub is mandatory.
- **Scene Control Must Be Enabled:** Multi-tap scene control is disabled by default (Parameter 6 = 0). It must be explicitly enabled via configuration. When disabled, only single-tap on/off works.
- **Smart Bulb Mode:** Parameter 7 enables "smart bulb mode" where the physical load relay stays on but the paddle still sends Z-Wave commands. This is designed for controlling smart bulbs that need constant power.
- **Z-Wave 700 vs 800 Series:** The original ZEN76 used the Z-Wave 700 series chip (Silicon Labs ZGM130S). Newer revisions use the 800 series (ZG23) with improved range and lower power. Both are fully backward-compatible.
- **S2 Security Framework:** The ZEN76 supports all S2 security levels. S2 Authenticated is recommended, which requires entering the 5-digit DSK during inclusion. The DSK is also encoded in the QR code on the device.
- **No OTA Without Controller:** Firmware updates are delivered via the Z-Wave OTA command class through the Z-Wave controller. There is no way to update firmware without a Z-Wave hub.
- **Affordable:** The ZEN76 is typically priced between $20-30 USD, making it one of the most affordable quality Z-Wave switches available.
- **Neutral Wire Required:** Like most Z-Wave switches, the ZEN76 requires a neutral wire at the switch box.
- **Z-Wave-to-Matter Bridge Potential:** As the smart home industry moves toward Matter, Z-Wave-to-Matter bridge devices (like the Nabu Casa Z-Wave bridge) could expose Z-Wave devices as Matter endpoints. This is the most likely path for Haus integration of Z-Wave devices.

## Similar Devices

- **zooz-zwave-dimmer-zen77** -- Zooz ZEN77, dimmer variant of ZEN76 with Level Control command class
- **inovelli-blue-series-switch** -- Zigbee/Matter switch with LED bar, more features but higher price
- **leviton-decora-smart-dimmer** -- WiFi dimmer, cloud-dependent but no hub required
- **ge-enbrighten-zwave-switch** -- GE/Jasco Z-Wave switch, similar Z-Wave approach with different configuration parameters
- **shelly-plus-1** -- WiFi relay for in-wall use, has local API (no Z-Wave hub needed)
