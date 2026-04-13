---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "august-wifi-smart-lock"
name: "August WiFi Smart Lock (4th Gen)"
manufacturer: "August Home, Inc. (Yale / ASSA ABLOY)"
brand: "August"
model: "AUG-SL05-M01-S01"
model_aliases: ["AUG-SL05-M01-G01", "AUG-SL05-M01-S02", "SL05"]
device_type: "smart_lock"
category: "security"
product_line: "August Smart Lock"
release_year: 2020
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "D0:03:DF"        # August Home Inc.
    - "78:9C:85"        # August Home Inc.
  mdns_services: []     # August locks do not advertise mDNS
  mdns_txt_keys: []
  default_ports: []     # No open ports -- communicates outbound only
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^August"
    - "^august-lock"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP services

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "august"
  polling_interval_sec: 30
  websocket_event: ""
  setup_type: "oauth2"
  ai_chattable: false
  haus_milestone: ""

# --- CAPABILITIES ---
capabilities:
  - "lock_unlock"
  - "battery_level"

# --- PROTOCOL ---
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "oauth2"
  auth_detail: "August/Yale Access uses a proprietary API. Authentication via POST to https://api-production.august.com/session with email/phone and password, returns access_token and user_id. 2FA verification code sent via email/phone. Also supports local BLE control using the yalexs-ble protocol with offline keys obtained via cloud API."
  base_url_template: "https://api-production.august.com"
  tls: true
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "lock"
  power_source: "battery"
  mounting: "door"
  indoor_outdoor: "indoor"
  wireless_radios: ["wifi", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://august.com/products/august-wifi-smart-lock"
  api_docs: ""
  developer_portal: ""
  support: "https://support.august.com/"
  community_forum: ""
  image_url: ""
  fcc_id: "2ABQA-SL05"

# --- TAGS ---
tags: ["smart_lock", "wifi", "ble", "august", "yale", "assa_abloy", "auto_lock", "auto_unlock", "retrofit", "no_official_api"]
---

# August WiFi Smart Lock (4th Gen)

## What It Is

The August WiFi Smart Lock (4th Gen) is a retrofit smart deadbolt lock manufactured by August Home, a subsidiary of Yale (ASSA ABLOY). It installs over your existing deadbolt on the interior side of the door, preserving your existing keys for manual use. The 4th Gen model is notable for being the first August lock with built-in Wi-Fi, eliminating the need for a separate August Connect Wi-Fi bridge that previous generations required. It features BLE (Bluetooth Low Energy) for proximity-based control and auto-unlock, Wi-Fi for remote access and cloud integration, auto-lock with a configurable timer, DoorSense (a magnetic sensor that detects if the door is open or closed), and support for virtual keys shared via the August app. The lock is powered by two CR123A batteries with a typical lifespan of 3-6 months depending on usage.

## How Haus Discovers It

1. **OUI Match** -- During network scan, devices with MAC prefixes `D0:03:DF` or `78:9C:85` are flagged as August Home devices.
2. **Hostname Pattern** -- August locks may appear with hostnames starting with `August` in DHCP tables.
3. **No Local Probing** -- August locks expose no local HTTP/TCP API and no mDNS/SSDP services. The Wi-Fi connection is used exclusively for outbound communication to August's cloud.
4. **BLE Discovery** -- August locks advertise via BLE with a manufacturer-specific advertising payload. BLE scanning can detect nearby August locks by their BLE advertisement data (manufacturer ID and service UUIDs).

## Pairing / Authentication

August does not provide an official public API. Third-party integrations use the reverse-engineered August API (used by projects like `py-august` and `yalexs`).

### Cloud API Authentication

1. **Session Creation:** `POST https://api-production.august.com/session` with:
   ```json
   {
     "identifier": "email:{email}",
     "installId": "{uuid}",
     "password": "{password}"
   }
   ```
   Headers must include `x-august-api-key`, `x-kease-api-key`, and `Content-Type: application/json`. The API key values are embedded in the August mobile app and have been extracted by reverse-engineering efforts.

2. **2FA Verification:** August sends a verification code via email or SMS:
   ```
   POST https://api-production.august.com/validation/email
   {"value": "{email}"}
   ```
   Then validate the code:
   ```
   POST https://api-production.august.com/validate/email
   {"code": "{code}", "email": "{email}"}
   ```

3. **Token Usage:** The `x-august-access-token` header from the session response is used for all subsequent API calls. Tokens are long-lived but may expire.

### BLE (Local) Control

The August lock also supports local BLE control using the `yalexs-ble` protocol:

1. **Offline Key Retrieval:** Obtain the lock's offline BLE key via the cloud API (`GET /locks/{lock_id}` includes the `offline_key` and `offline_slot`).
2. **BLE Connection:** Connect to the lock via BLE using its advertising UUID.
3. **Command Encryption:** Commands are encrypted using the offline key with AES-CBC. The protocol uses a challenge-response handshake before accepting commands.
4. **Local Operation:** Once the BLE key is obtained (one-time cloud dependency), the lock can be controlled locally via BLE without internet access.

### Security Notes

- The unofficial API keys are extracted from August's mobile app and may change with app updates.
- August/Yale has transitioned to the "Yale Access" brand and API, though the `api-production.august.com` endpoints remain functional.
- BLE offline keys provide a path to internet-independent control but require initial cloud authentication.

## API Reference

All cloud endpoints use the `x-august-access-token` header.

**Base URL:** `https://api-production.august.com`

### List Locks

```
GET /users/locks/mine
```

Returns all locks associated with the account. Each lock includes:
- `LockID` -- unique lock identifier
- `LockName` -- user-assigned name
- `HouseName` -- associated house name
- `Type` -- lock type identifier
- `SerialNumber` -- hardware serial
- `battery` -- battery level (0.0 - 1.0)

### Get Lock Status

```
GET /locks/{lock_id}/status
```

Returns current lock state:
- `status` -- `"locked"` or `"unlocked"`
- `doorState` -- `"open"`, `"closed"`, or `"unknown"` (from DoorSense)
- `retryCount` -- number of retries for the last operation
- `totalUpdates` -- total state changes

### Lock

```
PUT /remoteoperate/{lock_id}/lock
```

Locks the deadbolt remotely via Wi-Fi. Returns updated status.

### Unlock

```
PUT /remoteoperate/{lock_id}/unlock
```

Unlocks the deadbolt remotely. Returns updated status.

### Get Activity Log

```
GET /houses/{house_id}/activities
```

Returns access history including manual key entries, app unlocks, auto-locks, and auto-unlocks with timestamps and user attribution.

## AI Capabilities

AI integration is not currently planned. If implemented, the AI concierge could:

- Report lock/unlock status and door open/closed state (via DoorSense)
- Lock/unlock the door via voice command (with appropriate security confirmation)
- Report battery level and warn when replacement is needed
- Show recent access activity log
- Report auto-lock status and countdown

## Quirks & Notes

- **Retrofit Design:** The August lock installs on the interior side of an existing deadbolt. The exterior keyhole and hardware remain unchanged, so existing physical keys still work. Installation typically takes 10-15 minutes with a screwdriver.
- **DoorSense:** A small magnetic sensor mounts on the door frame to detect whether the door is physically open or closed. This enables "auto-lock when closed" functionality and prevents locking when the door is open (which would be pointless). DoorSense calibration can be finicky on some door/frame combinations.
- **Auto-Unlock:** Uses phone BLE proximity to automatically unlock when you approach the door. This feature uses geofencing (leave home zone, then return) combined with BLE proximity to prevent false unlocks while you are home. Reliability varies -- it works well for most users but some report inconsistent behavior.
- **Auto-Lock:** Configurable timer (30 seconds to 5 minutes) to automatically re-lock after unlocking. Can also be set to auto-lock only when DoorSense detects the door is closed.
- **Battery Life:** Two CR123A lithium batteries. August advertises 3-6 months. Heavy Wi-Fi usage (frequent remote operations) drains batteries faster than BLE-only use. The lock sends low-battery push notifications at approximately 20%.
- **Wi-Fi vs BLE:** The lock prioritizes BLE for nearby phone control (lower latency, lower power) and uses Wi-Fi for remote access and cloud sync. Wi-Fi can be disabled to extend battery life if remote access is not needed.
- **Yale Access Transition:** August has been merging with Yale under the ASSA ABLOY umbrella. The August app has been replaced by the Yale Access app in some markets. The cloud API endpoint may transition from `api-production.august.com` to a Yale-branded domain.
- **Thread/Matter:** The 4th Gen WiFi lock does NOT support Thread or Matter. Yale's newer Assure Lock 2 is the Thread/Matter-capable successor in the same product family.
- **Z-Wave/Zigbee:** Not supported. This lock is Wi-Fi + BLE only.

## Similar Devices

- **yale-assure-lock-2** -- Yale's successor with Thread/Matter support (same parent company)
- **schlage-encode-plus** -- Competing Wi-Fi lock with HomeKit support
- **level-lock-plus** -- Invisible smart lock with Thread/Matter
- **kwikset-halo-touch** -- Wi-Fi lock with fingerprint reader
