---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "myq-smart-garage-hub"
name: "MyQ Smart Garage Hub"
manufacturer: "Chamberlain Group"
brand: "MyQ"
model: "MYQ-G0401"
model_aliases: ["MYQ-G0301", "G0401", "G0301", "821LMC", "MYQ-G0401-ES"]
device_type: "garage_controller"
category: "security"
product_line: "MyQ"
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
    - "64:52:99"        # Chamberlain Group
    - "7C:A7:B0"        # Chamberlain Group
    - "00:17:D1"        # Chamberlain Group (older)
  mdns_services: []     # MyQ does not advertise mDNS
  mdns_txt_keys: []
  default_ports: []     # No open ports -- outbound cloud only
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^MyQ"
    - "^myq-"
    - "^Chamberlain"
    - "^GW[0-9]+"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP services

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "myq"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "garage_open_close"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "MyQ uses OAuth 2.0 via Chamberlain's identity service. Authentication via POST to https://partner-identity.myq-cloud.com/connect/token. Third-party API access has been repeatedly blocked by Chamberlain -- they actively detect and shut down unofficial integrations. Ratwatch (Home Assistant MyQ integration) and others have been engaged in a multi-year cat-and-mouse with Chamberlain's API enforcement."
  base_url_template: "https://api.myqdevice.com/api/v5.2"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "hub"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi"]

# --- LINKS ---
links:
  product_page: "https://www.myq.com/smart-garage-hub"
  api_docs: ""
  developer_portal: "https://www.myq.com/myq-partnership"
  support: "https://support.myq.com/"
  community_forum: "https://community.myq.com/"
  image_url: ""
  fcc_id: "HBW7430"

# --- TAGS ---
tags: ["garage", "cloud_only", "chamberlain", "liftmaster", "api_hostile", "no_local_api", "walled_garden", "subscription"]
---

# MyQ Smart Garage Hub

## What It Is

The MyQ Smart Garage Hub is a Wi-Fi connected garage door controller manufactured by Chamberlain Group, the parent company of LiftMaster, Chamberlain, and Craftsman garage door opener brands. The hub connects to your Wi-Fi network and communicates with your existing garage door opener via a wired sensor (mounted on the garage door) and the hub's built-in radio transmitter that simulates a button press on the opener. It works with most garage door openers manufactured after 1993 that use standard safety sensors. The MyQ app provides remote monitoring (open/closed status) and control (open/close). MyQ is the dominant smart garage platform in North America -- Chamberlain Group controls approximately 60% of the residential garage door opener market, and many newer LiftMaster/Chamberlain openers have MyQ built-in. However, MyQ is notorious in the smart home community for its aggressively cloud-locked ecosystem and its history of actively blocking third-party API access.

## How Haus Discovers It

1. **OUI Match** -- Devices with MAC prefixes `64:52:99`, `7C:A7:B0`, or `00:17:D1` (Chamberlain Group) are flagged during network scanning.
2. **Hostname Pattern** -- MyQ hubs typically appear with hostnames starting with `MyQ`, `myq-`, `Chamberlain`, or `GW` in DHCP tables.
3. **No Local Discovery** -- MyQ devices expose no local API, no mDNS, no SSDP, and no open ports. The device communicates exclusively with MyQ's cloud servers via outbound HTTPS/MQTT connections.

## Pairing / Authentication

### MyQ App Setup

1. **Account Creation:** Create a MyQ account in the mobile app.
2. **Hub Pairing:** The app discovers the hub via BLE during setup and provisions Wi-Fi credentials.
3. **Door Sensor Installation:** Mount the tilt/magnetic sensor on the garage door and pair it with the hub.
4. **Opener Learning:** The hub "learns" your garage door opener's radio frequency by pressing the learn button on the opener.

### API Authentication (Historical -- Frequently Broken)

The MyQ API has gone through multiple revisions as Chamberlain blocks third-party access:

1. **Legacy API (v5/v5.1/v5.2):** Used email/password authentication via `POST https://api.myqdevice.com/api/v5/Login`. Required specific `MyQApplicationId` header values. Chamberlain rotated application IDs and added CAPTCHAs to block automated access.

2. **OAuth 2.0 Migration (2023+):** MyQ migrated to OAuth 2.0 via `partner-identity.myq-cloud.com`. The flow requires:
   - Browser-based login (with CAPTCHA)
   - OAuth authorization code exchange
   - Token refresh

3. **Active Blocking:** In October 2023, Chamberlain formally shut down all unauthorized third-party API access and published a blog post declaring it intentional. They implemented:
   - CAPTCHA on all login flows
   - Device fingerprinting
   - Rate limiting and IP blocking
   - Token revocation for detected automated clients
   - Legal threats to integration maintainers

### Current State of Third-Party Access

As of 2024-2025, the primary methods for third-party MyQ integration are:

1. **ratgdo (Local Hardware Bypass):** A community-developed hardware device (ratgdo.com) that connects directly to the garage door opener's control board, bypassing MyQ entirely. It provides local MQTT/ESPHome control. This is the recommended approach for users who want reliable local control.
2. **Home Assistant Cloud (Nabu Casa):** Chamberlain allows Nabu Casa (Home Assistant's cloud service) to use the MyQ API as an authorized partner, though this has been intermittent.
3. **Homebridge MyQ Plugin:** Community-maintained, frequently breaks and gets fixed as Chamberlain changes their blocking methods.

### Security Notes

- Chamberlain views unauthorized API access as a security concern and has been increasingly aggressive about blocking it.
- The MyQ account uses email/password authentication. 2FA is available but optional.
- Garage doors are a significant physical security concern -- unauthorized open commands could expose your home.

## API Reference

### MyQ Cloud API (Unofficial, Frequently Blocked)

**Base URL:** `https://api.myqdevice.com/api/v5.2` (subject to change)

**Note:** This documentation is for reference only. Chamberlain actively blocks all unauthorized third-party API access. These endpoints may not work without an authorized partner agreement.

### List Devices

```
GET /api/v5.2/Accounts/{account_id}/Devices
Authorization: Bearer {access_token}
```

Returns all MyQ devices including garage door openers, lamps, and gateways. Each device includes:
- `serial_number` -- unique device identifier
- `device_type` -- type identifier (e.g., `wifigaragedooropener`, `garagedooropener`, `lamp`)
- `name` -- user-assigned name
- `state` -- object containing `door_state` (`open`, `closed`, `opening`, `closing`, `stopped`), `last_update`, `online`

### Get Door State

```
GET /api/v5.2/Accounts/{account_id}/Devices/{serial}/state
Authorization: Bearer {access_token}
```

Returns current door state:
- `door_state` -- `open`, `closed`, `opening`, `closing`, `stopped`
- `last_update` -- ISO 8601 timestamp of last state change
- `online` -- boolean connectivity status

### Open/Close Door

```
PUT /api/v5.2/Accounts/{account_id}/Devices/{serial}/actions
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "action_type": "open"
}
```

Or `"action_type": "close"` to close. The door takes approximately 10-15 seconds to fully open or close.

### ratgdo Local API (Recommended Alternative)

The ratgdo hardware device provides a local REST/MQTT API:

**MQTT Topics:**
- `ratgdo/{device}/status/door` -- `open`, `closed`, `opening`, `closing`
- `ratgdo/{device}/command/door` -- Publish `open`, `close`, `toggle`
- `ratgdo/{device}/status/light` -- Opener light state
- `ratgdo/{device}/command/light` -- Control opener light

**ESPHome API:** When flashed with ESPHome firmware, ratgdo exposes a native ESPHome API for Home Assistant integration with auto-discovery.

## AI Capabilities

AI integration is not currently planned due to API access challenges. If implemented (likely via ratgdo or authorized partner API), the AI concierge could:

- Report garage door state (open/closed/opening/closing)
- Open/close the garage door with security confirmation
- Report when the door was last opened/closed
- Alert if the door has been open for an extended period

## Quirks & Notes

- **API Hostility:** Chamberlain is the most hostile major smart home manufacturer toward third-party integrations. They have:
  - Published blog posts explicitly stating they will block all unauthorized API access
  - Added CAPTCHAs to programmatic login flows
  - Rotated API keys and application IDs to break existing integrations
  - Sent cease-and-desist letters to integration developers
  - Implemented device fingerprinting and behavioral analysis to detect automated clients
  - In November 2023, blocked the Home Assistant MyQ integration entirely

- **Subscription Push:** In 2023, Chamberlain introduced a paid subscription tier for certain MyQ features, further frustrating users who felt they had already paid for the hardware.

- **ratgdo -- The Community Response:** The ratgdo (Rage Against The Garage Door Opener) hardware project emerged as a direct response to Chamberlain's API lockdown. It is a small circuit board ($30-40) that connects to the garage door opener's wiring and provides local MQTT/ESPHome/Home Assistant control without any cloud dependency. It works with Security+ 2.0 (Chamberlain/LiftMaster) and dry-contact openers.

- **Market Dominance:** Despite the API hostility, MyQ is unavoidable because Chamberlain Group manufactures the majority of residential garage door openers sold in North America (LiftMaster, Chamberlain, Craftsman brands). Many new openers ship with MyQ built-in.

- **Google Home / Amazon Alexa:** MyQ maintains authorized integrations with Google Home and Amazon Alexa for voice control. The Alexa integration requires an additional Key by Amazon device for auto-opening (security measure).

- **Apple HomeKit:** MyQ dropped HomeKit support. Users who want HomeKit must use a third-party bridge (like Homebridge/ratgdo).

- **Door State Accuracy:** The tilt/magnetic sensor can occasionally report incorrect states, especially if the door is partially open or the sensor is misaligned. There is typically a 5-10 second delay between physical state change and cloud state update.

- **Security Timeout:** For safety, MyQ requires re-authentication (app login) before allowing a door open command if the app has been inactive for an extended period.

## Similar Devices

- **tailwind-iq3** -- Integration-friendly garage controller with documented local API
- **meross-smart-garage-opener** -- Budget garage controller with local MQTT protocol
