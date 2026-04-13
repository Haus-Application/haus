---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "leviton-decora-smart-dimmer"
name: "Leviton Decora Smart WiFi Dimmer"
manufacturer: "Leviton Manufacturing Co., Inc."
brand: "Leviton"
model: "DW6HD"
model_aliases: ["DW6HD-1BZ", "D26HD", "D26HD-1BW", "D26HD-2RW", "DW1KD", "DERA"]
device_type: "leviton_dimmer"
category: "lighting"
product_line: "Decora Smart"
release_year: 2017
discontinued: false
price_range: "$$"

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
    - "E8:F2:E2"        # Leviton Manufacturing (observed on Decora Smart WiFi devices)
    - "64:52:99"        # Espressif Systems (some Decora Smart variants use ESP32)
  mdns_services:
    - "_hap._tcp"       # HomeKit Accessory Protocol (2nd gen D26HD models)
  mdns_txt_keys:
    - "md"              # model name in HAP advertisement
    - "id"              # device ID
    - "sf"              # status flags (1 = not paired)
    - "ci"              # category identifier
  default_ports: [80, 443]
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Leviton"
    - "^DW6HD"
    - "^D26HD"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints:
  - port: 80
    path: "/"
    method: "GET"
    expect_status: 200
    title_contains: ""
    server_header: ""
    body_contains: "Leviton"
    headers: {}

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "leviton"
  polling_interval_sec: 30
  websocket_event: "leviton:state"
  setup_type: "oauth2"
  ai_chattable: true
  haus_milestone: "M8"

# --- CAPABILITIES ---
capabilities:
  - "on_off"
  - "brightness"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "My Leviton cloud API requires OAuth2 authentication via https://my.leviton.com/api/. Login with email/password to receive access token and refresh token. No documented local API for 1st gen (DW6HD). 2nd gen (D26HD) supports HomeKit (HAP) which provides limited local control via Apple Home ecosystem."
  base_url_template: "https://my.leviton.com/api"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "switch"
  power_source: "hardwired"
  mounting: "in_wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.leviton.com/products/dw6hd-1bz"
  api_docs: ""
  developer_portal: ""
  support: "https://www.leviton.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "LS8-DW6HD"

# --- TAGS ---
tags: ["wifi", "in_wall", "dimmer", "cloud_api", "leviton", "decora", "no_local_api", "homekit_v2"]
---

# Leviton Decora Smart WiFi Dimmer

## What It Is

The Leviton Decora Smart WiFi Dimmer (models DW6HD and D26HD) is an in-wall smart dimmer switch that connects directly to WiFi and is primarily controlled through the My Leviton app and cloud service. Leviton is one of the oldest and most established electrical device manufacturers in North America, and the Decora Smart line brings their traditional Decora-style aesthetic into the smart home. The DW6HD (1st generation) supports 600W incandescent/300W LED loads and communicates exclusively through the My Leviton cloud. The D26HD (2nd generation) adds Apple HomeKit support via HAP (HomeKit Accessory Protocol) over WiFi, providing some local control through the Apple Home ecosystem. Leviton also makes Z-Wave variants (DZ6HD) for hub-based systems. The WiFi models require a neutral wire and support single-pole, 3-way (with DD0SR companion), and multi-way configurations.

## How Haus Discovers It

1. **OUI Match** -- MAC addresses beginning with `E8:F2:E2` (Leviton) or `64:52:99` (Espressif, used in some variants) may indicate a Decora Smart device.
2. **mDNS Discovery (2nd Gen Only)** -- The D26HD advertises `_hap._tcp.local.` for HomeKit. TXT records include the `md` (model description), `id` (device ID), `sf` (status flags: 1 means unpaired), and `ci` (category identifier: 2 = bridge, 5 = lightbulb). This only helps identify 2nd-gen models.
3. **DHCP Hostname** -- Leviton devices may register with hostnames containing "Leviton", "DW6HD", or "D26HD".
4. **Cloud API Enumeration** -- With My Leviton account credentials, Haus can query the cloud API to enumerate all registered devices on the account, including model numbers and current state. This is the most reliable discovery method but requires user account linkage.
5. **Port Probe** -- Port 80 may respond on some models with basic device info, but this is not reliable for identification.

## Pairing / Authentication

### My Leviton Cloud (Primary)

1. User creates a My Leviton account at `https://my.leviton.com/` or via the My Leviton mobile app.
2. Device is added to the account via the app, which provisions WiFi credentials to the dimmer.
3. Haus integrates by having the user provide My Leviton credentials (email/password).
4. Haus authenticates via the My Leviton API:
   ```
   POST https://my.leviton.com/api/Person/login
   Content-Type: application/json

   {"email": "user@example.com", "password": "..."}
   ```
5. The response includes an `id` (access token) and `userId` used for all subsequent API calls.

### HomeKit (2nd Gen D26HD)

1. The D26HD shows a HomeKit setup code (8-digit) on the device or in packaging.
2. Pair via Apple Home app or a HAP-compatible controller.
3. Haus does not currently implement HAP client, but this is a potential future local control path.

## API Reference

### My Leviton Cloud API

**Base URL:** `https://my.leviton.com/api`

All requests require the `Authorization` header with the access token from login.

#### Login

```
POST /Person/login
Content-Type: application/json

{"email": "user@example.com", "password": "..."}
```

**Response:**
```json
{
  "id": "access_token_here",
  "userId": "user_id_here",
  "ttl": 1209600
}
```

#### List Residences

```
GET /Person/{userId}/Residences
Authorization: {access_token}
```

Returns array of residences, each containing an `id`.

#### List Devices (IotSwitches)

```
GET /Residences/{residenceId}/IotSwitches
Authorization: {access_token}
```

**Response:** Array of switch objects:
```json
{
  "id": "switch_id",
  "name": "Living Room Dimmer",
  "model": "DW6HD",
  "brightness": 75,
  "power": "ON",
  "connected": true,
  "serial": "...",
  "version": "1.5.24"
}
```

#### Control a Switch

```
PUT /IotSwitches/{switchId}
Authorization: {access_token}
Content-Type: application/json

{"power": "ON", "brightness": 75}
```

- `power` -- "ON" or "OFF"
- `brightness` -- 0-100 (integer percentage)

#### Get Switch State

```
GET /IotSwitches/{switchId}
Authorization: {access_token}
```

Returns the full switch object with current state.

### Limitations

- **No Local API (1st Gen):** The DW6HD has no documented or community-discovered local API. All control goes through My Leviton cloud.
- **Cloud Latency:** Commands take 500ms-2s round-trip through the cloud.
- **Token Expiry:** Access tokens expire after the TTL (default 1209600 seconds = 14 days). Haus must re-authenticate periodically.
- **Rate Limiting:** The My Leviton API has undocumented rate limits. Excessive polling may result in temporary blocks.

## AI Capabilities

When the AI concierge is chatting with a Leviton dimmer, it can:

- **Turn the dimmer on or off** via cloud command
- **Set brightness** as a percentage (0-100%)
- **Report current state** -- on/off, brightness level, connectivity
- **Report device info** -- model, firmware version, serial number
- **Warn about cloud dependency** -- inform user that cloud outages will prevent control

## Quirks & Notes

- **No Usable Local API:** This is the biggest limitation for Haus integration. The DW6HD communicates exclusively through the My Leviton cloud. If Leviton's servers are down or the internet is out, the switch still works physically (manual paddle control) but cannot be controlled remotely or via Haus.
- **2nd Gen HomeKit is Better:** The D26HD's HomeKit support via `_hap._tcp` means it can theoretically be controlled locally through HAP. If Haus implements a HAP controller in the future, this would enable local control of 2nd-gen models.
- **Z-Wave Variant:** The DZ6HD is the Z-Wave version of the same physical switch. It requires a Z-Wave hub but offers true local control. If the user has a Z-Wave setup, this is preferable for Haus.
- **Neutral Wire Required:** Unlike some smart switches, all Decora Smart WiFi dimmers require a neutral wire. This limits installation in older homes without neutral wires at the switch box.
- **Decora Aesthetic:** The switch uses Leviton's standard Decora paddle design and fits standard Decora wallplates. This is a strong selling point for users who want a professional-looking installation.
- **3-Way Wiring:** For 3-way setups, the DW6HD/D26HD pairs with the DD0SR-1Z companion dimmer (wired, not smart). The companion connects to the smart dimmer via traveler wire and does not need WiFi.
- **LED Indicator:** A small green LED on the switch indicates WiFi connection status. It blinks during provisioning and is solid when connected.
- **Firmware Updates:** Firmware updates are pushed via the My Leviton cloud. Users cannot manually flash firmware.
- **Max Load Ratings:** 600W incandescent/halogen, 300W LED/CFL. Exceeding these ratings can cause overheating or failure.

## Similar Devices

- **leviton-decora-zwave-dimmer** -- DZ6HD, Z-Wave variant with hub-based local control
- **inovelli-blue-series-switch** -- Zigbee dimmer with LED bar, local protocol
- **lutron-caseta-dimmer** -- Clear Connect protocol, requires Lutron bridge but extremely reliable
- **kasa-smart-dimmer** -- TP-Link Kasa dimmer, local XOR protocol control
- **zooz-zwave-switch-zen76** -- Z-Wave on/off switch (non-dimming)
