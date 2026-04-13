---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "level-lock-plus"
name: "Level Lock+"
manufacturer: "Level Home, Inc."
brand: "Level"
model: "L-D11U"
model_aliases: ["C-L12U", "C-L11U", "L-DB1"]
device_type: "smart_lock"
category: "security"
product_line: "Level Lock"
release_year: 2023
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
  protocols_spoken: ["bluetooth", "thread", "matter"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: []      # Level locks use Thread (802.15.4), not WiFi -- no WiFi MAC on network
  mdns_services:
    - "_matter._tcp"    # Matter service advertisement (via Thread border router)
    - "_matterc._udp"   # Matter commissioning
  mdns_txt_keys:
    - "VP"              # Vendor/Product ID
    - "D"               # Discriminator
  default_ports: []     # Thread devices do not appear on IP network directly
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []  # Not on WiFi network
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No HTTP services

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
  auth_detail: "Matter protocol over Thread. Commissioning via BLE using PASE with setup code, then CASE for ongoing encrypted communication over Thread mesh. The lock communicates on the Thread 802.15.4 mesh and reaches IP networks via a Thread Border Router."
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
  product_page: "https://level.co/products/level-lock-plus"
  api_docs: ""
  developer_portal: ""
  support: "https://support.level.co/"
  community_forum: ""
  image_url: ""
  fcc_id: "2AQ5Q-LD11U"

# --- TAGS ---
tags: ["smart_lock", "matter", "thread", "invisible", "ble", "level", "deadbolt", "no_wifi", "cr2_battery", "industrial_design"]
---

# Level Lock+

## What It Is

The Level Lock+ is an invisible smart deadbolt manufactured by Level Home, Inc. Its defining characteristic is that it looks identical to a traditional deadbolt from both inside and outside the door -- all smart home electronics, motor, and batteries are miniaturized and hidden entirely within the deadbolt housing itself. There is no external keypad, no visible electronics, and no bulky interior escutcheon. The lock is controlled via Bluetooth from the Level app, Apple Home Key (NFC/UWB tap-to-unlock on iPhone/Apple Watch), physical key, and auto-unlock based on phone proximity. The Level Lock+ (2023 model) added Thread radio support, making it one of the first invisible smart locks to be Matter-certified. It uses a single CR2 battery with approximately one year of battery life. The lock is ANSI/BHMA Grade A certified.

## How Haus Discovers It

1. **No WiFi Presence** -- The Level Lock+ does not have Wi-Fi. It will NOT appear in network scans, ARP tables, or DHCP logs. It communicates exclusively via BLE and Thread (802.15.4).
2. **Thread Mesh Discovery** -- With a Thread Border Router present on the network, the Level Lock+ joins the Thread mesh and becomes addressable via IPv6. It advertises as a Matter device through the border router's mDNS proxy.
3. **Matter Discovery** -- Once on the Thread network, the lock is discoverable via mDNS (`_matter._tcp`, `_matterc._udp`) with Matter Vendor/Product identifiers for Level Home.
4. **BLE Discovery** -- The lock advertises via BLE for initial pairing and Matter commissioning. BLE scanning can detect it by Level's manufacturer-specific advertising data.

## Pairing / Authentication

### Matter Commissioning (Thread)

1. **Setup Code:** The lock includes a Matter QR code and 11-digit manual pairing code on the included card and inside the battery compartment.
2. **BLE Commissioning:** The Matter controller discovers the lock via BLE advertisement, connects using PASE with the setup code, and provisions Thread network credentials.
3. **Thread Joining:** The lock joins the Thread mesh network and becomes reachable via IPv6 through the Thread Border Router.
4. **CASE Session:** Ongoing communication uses CASE with device attestation certificates for mutual authentication.
5. **Multi-Admin:** The lock supports multi-admin, allowing simultaneous control from Apple Home, Google Home, Haus, and other Matter controllers.

### Apple Home Key

1. **HomeKit Pairing:** The lock pairs with Apple HomeKit via the Home app using the HomeKit setup code.
2. **Home Key Provisioning:** After HomeKit pairing, a Home Key can be added to Apple Wallet.
3. **Tap-to-Unlock:** Tap iPhone or Apple Watch to the lock for NFC unlock. iPhone 11+ / Apple Watch Series 6+ with UWB support also enable Express Mode for hands-free proximity unlock.

### Level App (BLE Direct)

1. **BLE Pairing:** The Level app discovers and pairs with the lock via BLE.
2. **Account Required:** A Level account is required for sharing access and remote features.
3. **Local BLE Control:** Once paired, the Level app can lock/unlock via BLE without internet.

### Security Notes

- Thread communication is encrypted at both the mesh layer and application layer.
- Matter CASE uses device attestation certificates rooted in the Distributed Compliance Ledger (DCL).
- The physical key remains functional at all times -- electronic failure does not lock you out.
- Auto-unlock uses phone BLE proximity with geofencing to prevent false triggers.

## API Reference

### Matter Door Lock Cluster

The lock exposes the Matter Door Lock cluster (Cluster ID 0x0101):

**Attributes:**
- `LockState` (0x0000) -- `Locked`, `Unlocked`, `NotFullyLocked`
- `LockType` (0x0001) -- `DeadBolt`
- `ActuatorEnabled` (0x0002) -- Motor enabled state
- `OperatingMode` (0x0025) -- `Normal`, `Vacation`, `Privacy`, `NoRemoteLockUnlock`

**Power Source Cluster:**
- `BatPercentRemaining` -- Battery percentage (0-200, divide by 2 for actual %)
- `BatChargeLevel` -- `OK`, `Warning`, `Critical`

**Commands:**
- `LockDoor` (0x00) -- Lock the deadbolt
- `UnlockDoor` (0x01) -- Unlock the deadbolt

**Events:**
- `LockOperation` -- Emitted on lock/unlock with operation source
- `DoorLockAlarm` -- Triggered on jam detection

**Note:** The Level Lock+ Matter implementation is relatively minimal compared to locks with keypads. It does not expose User, Credential, or Schedule clusters since it has no keypad for PIN entry.

## AI Capabilities

AI integration is not currently planned but will be supported when Matter integration is implemented in M11. With Matter, the AI concierge could:

- Report lock/unlock status
- Lock/unlock the door with security confirmation
- Report battery level and charge state
- Report recent lock operations and their sources

## Quirks & Notes

- **Truly Invisible:** The Level Lock+ is the only smart lock that is completely invisible. Both the exterior and interior look like a standard deadbolt. The entire electronics package (BLE/Thread radios, motor, accelerometer, battery) fits inside the deadbolt barrel itself. This is a remarkable engineering achievement.
- **No Keypad:** Because the lock has no visible electronics, there is no keypad for PIN codes. Access methods are limited to: physical key, phone app (BLE), Apple Home Key (NFC/UWB), Matter controller, auto-unlock, and shared access via the Level app.
- **CR2 Battery:** The lock uses a single CR2 lithium battery (not the more common CR123A or AA). CR2 batteries are available at most stores but less common. Level advertises approximately 1 year of battery life. The lock sends push notifications when battery is low.
- **Touch-to-Lock:** The interior of the lock has a capacitive touch sensor. A quick touch on the lock body activates locking. This is the only physical interface besides the key and thumbturn.
- **Thread Border Router Required:** For Matter/Thread operation, a Thread Border Router must be on the network. Apple HomePod Mini, Apple TV 4K (2nd gen+), Google Nest Hub (2nd gen), and some Amazon Echo devices serve as Thread Border Routers.
- **No Wi-Fi:** The lock has no Wi-Fi radio at all. Remote access (when away from home) requires either a Thread Border Router with a cloud-connected controller (Apple Home, Google Home) or the Level app with a Level Connect bridge (sold separately, $79).
- **ANSI/BHMA Grade A:** Level uses the newer ANSI/BHMA A156.36 standard for smart locks, where Grade A is the highest. This is equivalent to the traditional Grade 1 certification.
- **Standard Deadbolt Compatibility:** The Level Lock+ replaces a standard single-cylinder deadbolt. It fits standard US door preps (2-1/8" bore hole, 1" faceplate). Installation takes approximately 15 minutes.
- **Auto-Unlock Reliability:** Level's auto-unlock uses phone BLE proximity combined with geofencing. You must leave and return to the geofence zone. Some users report inconsistency, particularly on Android devices.
- **Firmware Updates:** OTA updates delivered via BLE from the Level app. Thread/Matter firmware updates may also propagate via the Thread network.

## Similar Devices

- **yale-assure-lock-2** -- Visible smart lock with Thread/Matter and touchscreen keypad
- **august-wifi-smart-lock** -- Retrofit lock with Wi-Fi (same invisible-from-outside concept but visible interior)
- **schlage-encode-plus** -- Wi-Fi lock with HomeKit Home Key but no Thread/Matter
- **kwikset-halo-touch** -- Wi-Fi lock with fingerprint reader
