---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "yale-assure-lock-2"
name: "Yale Assure Lock 2"
manufacturer: "Yale (ASSA ABLOY)"
brand: "Yale"
model: "YRD450"
model_aliases: ["YRD410", "YRD420", "YRD430", "YRD450-WF1", "YRD450-ZW2", "YRD450-NR"]
device_type: "smart_lock"
category: "security"
product_line: "Yale Assure"
release_year: 2022
discontinued: false
price_range: "$$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "hybrid"
  local_api: false
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: false
  protocols_spoken: ["wifi", "bluetooth", "thread", "matter", "zigbee", "zwave"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "D0:03:DF"        # August Home Inc. (shared with August brand)
    - "78:9C:85"        # August Home Inc.
  mdns_services:
    - "_matter._tcp"    # When Thread/Matter module is installed
    - "_matterc._udp"   # Matter commissioning service
  mdns_txt_keys:
    - "VP"              # Matter Vendor/Product ID
    - "D"               # Matter discriminator
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Yale"
    - "^yale-lock"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP services

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "matter"
  polling_interval_sec: 0
  websocket_event: ""
  setup_type: "app_pairing"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities:
  - "lock_unlock"
  - "battery_level"

# --- PROTOCOL ---
protocol:
  type: "proprietary"
  port: 0
  transport: "UDP"
  encoding: "CBOR"
  auth_method: "none"
  auth_detail: "Matter protocol uses PASE (Passcode Authenticated Session Establishment) for commissioning with setup code, then CASE (Certificate Authenticated Session Establishment) for ongoing communication. Thread transport uses IPv6 mesh networking over 802.15.4 radio. Also supports Yale Access cloud API (same as August) when WiFi module is installed."
  base_url_template: ""
  tls: false
  tls_self_signed: false

# --- PHYSICAL ---
physical:
  form_factor: "lock"
  power_source: "battery"
  mounting: "door"
  indoor_outdoor: "indoor"
  wireless_radios: ["bluetooth_le", "thread"]

# --- LINKS ---
links:
  product_page: "https://shopyalehome.com/products/yale-assure-lock-2"
  api_docs: ""
  developer_portal: "https://developer.yale.com/"
  support: "https://shopyalehome.com/pages/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2ABQA-YRD450"

# --- TAGS ---
tags: ["smart_lock", "matter", "thread", "modular_radio", "yale", "assa_abloy", "keypad", "touchscreen", "multi_protocol"]
---

# Yale Assure Lock 2

## What It Is

The Yale Assure Lock 2 is a full-replacement smart deadbolt lock manufactured by Yale, a division of ASSA ABLOY. Unlike the August retrofit design, the Assure Lock 2 replaces the entire deadbolt assembly (interior and exterior). Its defining feature is the swappable radio module system -- the lock ships with a Bluetooth/Thread module by default, and users can purchase and swap in different modules for Wi-Fi, Zigbee, or Z-Wave connectivity. With the Thread module, it supports Matter, making it one of the first smart locks to achieve Matter certification. The lock features a capacitive touchscreen keypad on the exterior, tamper-proof design, ANSI/BHMA Grade 2 security rating, auto-lock, DoorSense, and virtual key sharing via the Yale Access app. Powered by four AA batteries.

## How Haus Discovers It

1. **OUI Match** -- Devices with MAC prefixes `D0:03:DF` or `78:9C:85` (August Home / Yale OUIs) are flagged during network scanning.
2. **Matter/Thread Discovery (Primary)** -- With the Thread module installed, the lock participates in a Thread mesh network and advertises via mDNS as a Matter device (`_matter._tcp` and `_matterc._udp`). The Matter discovery includes Vendor ID (Yale) and Product ID in TXT records.
3. **BLE Discovery** -- The lock advertises via BLE for phone pairing and Matter commissioning. Matter BLE advertisements include the discriminator and pairing hint.
4. **Wi-Fi Module (Alternate)** -- If the Wi-Fi module is installed instead, the lock connects to the LAN directly and uses the Yale Access (August) cloud API, discoverable via MAC OUI and hostname.
5. **Zigbee/Z-Wave (Alternate)** -- With Zigbee or Z-Wave modules, the lock is discoverable by the respective hub/coordinator protocol, not directly on the IP network.

## Pairing / Authentication

### Matter Commissioning (Thread Module)

1. **Setup Code:** The lock includes a Matter setup QR code and 11-digit manual pairing code on the documentation and inside the battery compartment.
2. **BLE Commissioning:** The Matter controller (Haus hub) connects to the lock via BLE using PASE (Passcode Authenticated Session Establishment) with the setup code.
3. **Thread Network Credentials:** During commissioning, the controller provides Thread network credentials to the lock so it joins the Thread mesh.
4. **CASE Establishment:** After commissioning, ongoing communication uses CASE (Certificate Authenticated Session Establishment) with device attestation certificates.
5. **Matter Fabric:** The lock is added to the Matter fabric and can be controlled by any Matter controller on the same fabric.

### Yale Access Cloud (Wi-Fi Module)

When the Wi-Fi module is installed, the lock uses the same Yale Access (August) cloud API as the August WiFi Smart Lock. See the `august-wifi-smart-lock` entry for API details.

### Security Notes

- Matter commissioning is a one-time process. The setup code should be treated as a secret -- anyone with the code can commission the device.
- The lock supports multiple Matter fabrics (multi-admin), allowing it to be controlled by multiple ecosystems simultaneously (Apple Home, Google Home, Haus, etc.).
- Thread communication is encrypted at the network layer (MLE) and at the application layer (Matter CASE sessions).

## API Reference

### Matter Door Lock Cluster

When connected via Matter, the lock exposes the standard Matter Door Lock cluster (Cluster ID 0x0101):

**Attributes:**
- `LockState` (0x0000) -- Current lock state: `NotFullyLocked`, `Locked`, `Unlocked`, `Unlatched`
- `LockType` (0x0001) -- Lock type: `DeadBolt`
- `ActuatorEnabled` (0x0002) -- Whether the lock motor is enabled
- `DoorState` (0x0003) -- Door physical state: `Open`, `Closed`, `DoorJammed`, `DoorForcedOpen` (from DoorSense)
- `NumberOfTotalUsersSupported` (0x0011) -- Max user count (typically 250)
- `NumberOfPINUsersSupported` (0x0012) -- Max PIN users
- `OperatingMode` (0x0025) -- `Normal`, `Vacation`, `Privacy`, `NoRemoteLockUnlock`
- `BatPercent` (from Power Source cluster) -- Battery percentage

**Commands:**
- `LockDoor` (0x00) -- Lock the deadbolt. Optional PIN code parameter.
- `UnlockDoor` (0x01) -- Unlock the deadbolt. Optional PIN code parameter.
- `SetUser` -- Add/modify a user with PIN code, credential type, and schedule.
- `GetUser` -- Retrieve user information.
- `ClearUser` -- Remove a user and all associated credentials.
- `SetCredential` -- Set a PIN code for a user.

**Events:**
- `DoorLockAlarm` -- Triggered on lock jam, forced entry, tamper alert
- `LockOperation` -- Emitted on every lock/unlock with source (manual, keypad, remote, auto-lock)
- `LockOperationError` -- Failed lock operations

### Yale Access Cloud API

When using the Wi-Fi module, the cloud API is identical to the August API documented in `august-wifi-smart-lock`. Key endpoints:

- `GET /users/locks/mine` -- List locks
- `GET /locks/{lock_id}/status` -- Get lock status
- `PUT /remoteoperate/{lock_id}/lock` -- Lock
- `PUT /remoteoperate/{lock_id}/unlock` -- Unlock

## AI Capabilities

AI integration is not currently planned but will be supported when Matter integration is implemented in M11. With Matter, the AI concierge could:

- Report lock/unlock status and door open/closed state
- Lock/unlock the door with security confirmation prompts
- Report battery level
- Manage PIN codes for users
- Report recent lock activity (lock operations, alarms)
- Set operating modes (normal, vacation, privacy)

## Quirks & Notes

- **Swappable Radio Modules:** The lock's modular design means the connectivity protocol depends entirely on which module is installed. The lock ships with a Bluetooth + Thread module. Alternate modules (Wi-Fi, Zigbee, Z-Wave) are purchased separately ($30-40 each). Only one module can be installed at a time.
- **Matter Certification:** The Assure Lock 2 with Thread module was among the first wave of Matter-certified door locks. It supports the Matter Door Lock device type with User, Credential, and Schedule features.
- **Thread Border Router Required:** The Thread module requires a Thread Border Router on the network to bridge between Thread mesh and IP. Apple TV 4K, HomePod Mini, Google Nest Hub (2nd Gen), and certain Amazon Echo devices function as Thread Border Routers.
- **Battery Life:** Four AA batteries. Yale advertises up to 9 months with normal use (approximately 10 lock/unlock cycles per day). Wi-Fi module draws more power than Thread, reducing battery life to approximately 3-4 months.
- **Touchscreen Keypad:** The exterior capacitive touchscreen supports PIN codes (4-8 digits), one-touch locking, and a privacy mode button. The screen illuminates on touch and is rated for outdoor use (IP55 weather resistance on touchscreen models).
- **DoorSense:** Same magnetic door position sensor as August locks. Detects open/closed/ajar states. Exposed as the `DoorState` attribute in the Matter Door Lock cluster.
- **ANSI/BHMA Grade 2:** Residential security rating. Meets ANSI A156.36 and UL 294 standards. Not as robust as Grade 1 (commercial) but appropriate for residential use.
- **Multi-Admin:** Matter supports multiple fabrics, meaning the lock can be simultaneously controlled by Apple Home, Google Home, Haus, and other Matter controllers without conflict.
- **Auto-Lock and Auto-Unlock:** Available via the Yale Access app. Auto-unlock uses phone BLE proximity (same technology as August). Auto-lock configurable from 30 seconds to 5 minutes.
- **Firmware Updates:** OTA firmware updates delivered via BLE from the Yale Access app. Thread/Matter firmware updates may eventually be supported via the Matter OTA protocol.

## Similar Devices

- **august-wifi-smart-lock** -- Same parent company, retrofit design, Wi-Fi + BLE only
- **schlage-encode-plus** -- Competing Wi-Fi lock with HomeKit (no Matter yet)
- **level-lock-plus** -- Invisible Thread/Matter lock
- **kwikset-halo-touch** -- Wi-Fi lock with fingerprint (no Matter)
