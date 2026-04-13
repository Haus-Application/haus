---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "wyze-bulb-color"
name: "Wyze Bulb Color"
manufacturer: "Wyze Labs"
brand: "Wyze"
model: "WLPA19C"
model_aliases: ["WLPA19C", "Wyze Bulb Color"]
device_type: "cloud_only_bulb"
category: "lighting"
product_line: "Wyze Bulb"
release_year: 2021
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
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["2C:AA:8E", "7C:78:B2", "A4:DA:22", "D0:3F:27"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["Wyze.*", "wyze.*"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "not_feasible"
  integration_key: ""
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: ""
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp"]

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "TLS"
  encoding: "binary"
  auth_method: "oauth2"
  auth_detail: "Wyze cloud authentication via email/password + 2FA; API tokens obtained from auth.wyze.com; MQTT TLS to Wyze cloud servers for device control"
  base_url_template: "https://api.wyzecam.com/"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "bulb"
  power_source: "mains"
  mounting: "ceiling"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.wyze.com/products/wyze-bulb-color"
  api_docs: ""
  developer_portal: ""
  support: "https://support.wyze.com/"
  community_forum: "https://forums.wyze.com/"
  image_url: ""
  fcc_id: "2AUIUWLPA19C"

# --- TAGS ---
tags: ["wifi", "cloud_only", "no_local_api", "budget", "not_feasible"]
---

# Wyze Bulb Color

## What It Is

The Wyze Bulb Color is an affordable WiFi smart bulb from Wyze Labs that supports full RGB color (16 million colors), tunable white (2700K-6500K), and 1100 lumens brightness at the E26 form factor. It is positioned as a budget-friendly smart bulb (typically under $8). However, the bulb is entirely cloud-dependent — all control flows through Wyze's proprietary cloud servers, and there is no documented or usable local API. This makes it not feasible for Haus integration, which requires local-first device control.

## How Haus Discovers It

Haus can detect Wyze bulbs on the network but cannot control them:

1. **OUI Match**: MAC addresses beginning with `2C:AA:8E`, `7C:78:B2`, `A4:DA:22`, or `D0:3F:27` are associated with Wyze devices (some of these OUIs are shared with the Wyze chipset vendor Tuya/Espressif, so false positives are possible).
2. **Network Presence**: The bulb will appear as a WiFi client on the network, maintaining persistent outbound connections to Wyze cloud servers.
3. **No Local Services**: The bulb does not expose any local HTTP servers, mDNS services, SSDP responses, or open ports. There is nothing to probe locally.

Haus will classify this device as "detected_only" — acknowledged on the network but not controllable.

## Pairing / Authentication

Setup is done exclusively through the Wyze mobile app:

1. Create a Wyze account (email + password).
2. Open the Wyze app and start "Add Device" flow.
3. The bulb is initially provisioned over BLE — the app sends WiFi credentials to the bulb.
4. The bulb connects to WiFi and registers with Wyze cloud.
5. All subsequent control flows through the Wyze cloud (MQTT over TLS).

There is no way to pair or authenticate the bulb for local control.

## API Reference

There is no official public API for Wyze devices. The ecosystem is closed.

### Unofficial / Reverse-Engineered Information

Community projects (such as `wyze-sdk` on GitHub and the Home Assistant Wyze integration) have reverse-engineered parts of the Wyze cloud API:

- **Auth Endpoint**: `https://auth-prod.api.wyze.com/api/user/login` — email/password authentication, returns access/refresh tokens.
- **Device List**: `https://api.wyzecam.com/app/v2/home_page/get_object_list` — returns all devices on the account.
- **Bulb Control**: Device commands are sent via MQTT to Wyze cloud brokers, which relay to the device. The MQTT protocol uses a proprietary binary format.
- **2FA Required**: Wyze enforces 2FA on accounts, complicating automated API access.
- **Rate Limiting**: Aggressive rate limiting on cloud endpoints; frequent requests may trigger account lockout.

### Why Local Control Is Not Feasible

- The Wyze Bulb Color uses a Tuya-derived chipset but does NOT support the Tuya local protocol. Wyze has customized the firmware to remove local control capabilities.
- The bulb communicates exclusively with Wyze cloud servers over TLS-encrypted MQTT.
- There are no open local ports and no local API of any kind.
- Attempts to use Tuya local key extraction do not work because Wyze has replaced the Tuya cloud backend with their own.
- Wyze has shown no interest in providing local API access or Matter/HomeKit support for their bulb products.

## AI Capabilities

Not applicable — this device cannot be integrated with Haus.

## Quirks & Notes

- **Extremely Affordable**: At under $8 per bulb, the Wyze Bulb Color is one of the cheapest RGB smart bulbs available. The trade-off is complete cloud dependency.
- **Cloud Outage = No Control**: When Wyze servers go down (which has happened multiple times), all Wyze bulbs become uncontrollable. They retain their last state but cannot be changed.
- **Tuya Heritage**: The hardware uses a Tuya/Espressif chipset, but the firmware has been heavily customized. Standard Tuya local protocols do not work.
- **No Matter Roadmap**: As of early 2026, Wyze has not announced Matter support for any of their bulb products.
- **Security Concerns**: Wyze has had multiple data breaches and security incidents. User credentials sent to their cloud should be considered a risk.
- **Unofficial Integrations Break Frequently**: The reverse-engineered cloud API used by Home Assistant and other projects breaks regularly when Wyze changes their backend.
- **Recommendation for Users**: If a user has Wyze bulbs and wants local control, recommend replacing them with LIFX, Nanoleaf Shapes, or other locally-controllable alternatives.

## Similar Devices

- **lifx-a19-color** — Similar WiFi color bulb but with full local LAN protocol (recommended alternative)
- **nanoleaf-essentials-a19** — Thread/Matter color bulb, local control via Matter
- **govee-rgbic-led-strip** — Budget brand but offers a local LAN API unlike Wyze
