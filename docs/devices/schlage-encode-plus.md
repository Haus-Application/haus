---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "schlage-encode-plus"
name: "Schlage Encode Plus Smart WiFi Deadbolt"
manufacturer: "Schlage (Allegion)"
brand: "Schlage"
model: "BE499WB"
model_aliases: ["BE499WB1V", "BE499WB CEN", "BE499WB CAM"]
device_type: "smart_lock"
category: "security"
product_line: "Schlage Encode"
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
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "B4:A2:EB"        # Schlage / Allegion
    - "00:1A:22"        # Allegion (older)
  mdns_services: []     # Schlage Encode does not advertise mDNS services
  mdns_txt_keys: []
  default_ports: []     # No open ports -- outbound cloud only
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Schlage"
    - "^schlage-encode"
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP services

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "schlage"
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
  auth_detail: "Schlage uses the Schlage Home app with a proprietary cloud API (api.allegion.com). Authentication via OAuth 2.0 with Allegion's identity service. No official public API or developer program. Apple HomeKit provides local control via HAP (HomeKit Accessory Protocol) over Wi-Fi -- no cloud required once paired."
  base_url_template: ""
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
  product_page: "https://www.schlage.com/en/home/smart-locks/encode-plus.html"
  api_docs: ""
  developer_portal: ""
  support: "https://www.schlage.com/en/home/support.html"
  community_forum: ""
  image_url: ""
  fcc_id: "CEL-BE499WB"

# --- TAGS ---
tags: ["smart_lock", "wifi", "homekit", "home_key", "uwb", "allegion", "schlage", "keypad", "ansi_grade_1", "no_official_api"]
---

# Schlage Encode Plus Smart WiFi Deadbolt

## What It Is

The Schlage Encode Plus is a premium Wi-Fi smart deadbolt manufactured by Schlage, a brand of Allegion. It is best known as the first smart lock to support Apple Home Key, which allows users to unlock their door by tapping their iPhone or Apple Watch against the lock -- using NFC for close range and UWB (Ultra-Wideband) for hands-free "Express Mode" unlock when approaching with a qualifying device. The lock features a capacitive touchscreen keypad, built-in Wi-Fi (no bridge or hub required), Bluetooth LE, Apple HomeKit native support, ANSI/BHMA Grade 1 residential security rating (the highest available), a built-in alarm sensor, and Schlage's well-regarded physical security construction. It replaces the entire deadbolt assembly and is available in Century and Camelot trim styles. Powered by four AA batteries.

## How Haus Discovers It

1. **OUI Match** -- Devices with MAC prefixes `B4:A2:EB` or `00:1A:22` (Allegion/Schlage) are flagged during network scanning.
2. **Hostname Pattern** -- Schlage locks may appear with hostnames starting with `Schlage` in DHCP tables.
3. **No Local Discovery Services** -- The Encode Plus does not advertise mDNS or SSDP services. It connects to Schlage's cloud via Wi-Fi and speaks HAP (HomeKit Accessory Protocol) for local Apple HomeKit control.
4. **HomeKit Discovery** -- If implementing HomeKit/HAP integration, the lock advertises via Bonjour as a HomeKit accessory (`_hap._tcp`). However, HomeKit accessories use encrypted pairing and cannot be controlled without going through the HomeKit pairing flow.

## Pairing / Authentication

Schlage does not provide an official public API. Integration paths are limited.

### Schlage Home App (Cloud)

1. **Account Creation:** Create an account in the Schlage Home app.
2. **Lock Pairing:** The app discovers the lock via BLE, provisions Wi-Fi credentials, and registers the lock with Schlage's cloud (Allegion API infrastructure).
3. **Cloud Control:** The Schlage Home app communicates via Allegion's cloud API. No official developer documentation exists for this API.

### Apple HomeKit (Local Control - Best Path)

1. **HomeKit Setup Code:** The lock includes an 8-digit HomeKit setup code (XXX-XX-XXX format) printed in the documentation and on the lock interior.
2. **HAP Pairing:** The lock advertises as a HomeKit accessory via Bonjour (`_hap._tcp`). Pairing is done via the Apple Home app using the setup code.
3. **Local Control:** Once paired, HomeKit controls the lock locally via HAP over Wi-Fi (encrypted SRP + ChaCha20-Poly1305). No internet required.
4. **Home Key Setup:** After HomeKit pairing, Home Key can be added to Apple Wallet. The lock stores the Home Key credential (ECC key pair) and uses NFC/UWB for tap-to-unlock.

### Home Key (UWB Express Mode)

1. **UWB Ranging:** iPhone 11+ or Apple Watch Series 6+ with UWB support can unlock the door automatically as you approach (Express Mode).
2. **NFC Fallback:** All Home Key-capable devices also support NFC tap-to-unlock, which works even when the phone battery is critically low (power reserve mode for up to 5 hours after shutdown).
3. **Key Sharing:** Home Keys can be shared with family members via Apple Wallet invitations.

### Security Notes

- HomeKit uses end-to-end encryption (Ed25519 key pairs exchanged during pairing). Apple does not have access to control your lock.
- Home Key uses ECC key pairs stored in the Secure Element of iPhone/Apple Watch. The lock stores up to 10 Home Key credentials.
- There is no known local API independent of HomeKit. The Schlage cloud API is proprietary and undocumented.
- For non-Apple ecosystems, the lock is effectively cloud-dependent via the Schlage Home app.

## API Reference

### HomeKit HAP (Local)

The lock exposes standard HomeKit services via HAP:

**Lock Mechanism Service:**
- `LockCurrentState` -- 0 (Unsecured), 1 (Secured), 2 (Jammed), 3 (Unknown)
- `LockTargetState` -- 0 (Unsecured), 1 (Secured)

**Battery Service:**
- `BatteryLevel` -- 0-100 percentage
- `StatusLowBattery` -- 0 (normal), 1 (low)
- `ChargingState` -- 0 (not charging)

**Lock Management Service:**
- `LockControlPoint` -- write-only for advanced lock operations
- `Version` -- firmware version string

Controlling the lock via HAP requires implementing the HomeKit Accessory Protocol, which involves:
- SRP (Secure Remote Password) pairing
- Ed25519 key pair exchange
- Encrypted sessions using ChaCha20-Poly1305
- HKDF-SHA-512 for key derivation

### Schlage Cloud API (Unofficial)

The Schlage Home app communicates with `api.allegion.com`. The API uses OAuth 2.0 for authentication. Specific endpoint documentation is not publicly available. Community reverse-engineering efforts have identified endpoints for:
- Device listing
- Lock/unlock commands
- Access code management
- Lock status polling
- Activity log retrieval

## AI Capabilities

AI integration is not currently planned. If implemented (likely via HomeKit bridge or cloud API), the AI concierge could:

- Report lock/unlock status
- Lock/unlock the door with security confirmation
- Report battery level
- Manage keypad access codes
- Report recent activity

## Quirks & Notes

- **Apple Home Key Exclusive:** The Encode Plus was the first and for some time the only lock supporting Apple Home Key with UWB Express Mode. This feature is exclusive to the Apple ecosystem -- Android users can use the keypad and Schlage app but not Home Key.
- **ANSI/BHMA Grade 1:** This is the highest residential security grade, exceeding Grade 2 (which most smart locks achieve). Grade 1 requires 800,000 locking cycles, 10 door strikes, and 250,000 key cycles. This makes the Encode Plus one of the most physically secure smart locks.
- **No Matter Support:** Despite being released in 2022, the Encode Plus does not support Matter or Thread. Allegion has not announced Matter plans for this specific model. A firmware update could theoretically add Matter over Wi-Fi, but this has not been indicated.
- **No Z-Wave/Zigbee:** Wi-Fi + BLE only. Cannot be integrated with traditional smart home hubs that rely on Z-Wave or Zigbee.
- **Built-in Alarm:** The lock includes a built-in tamper detection sensor. When the lock detects forced entry attempts or tampering, it triggers a loud alarm tone from the lock itself.
- **Battery Life:** Four AA batteries, advertised 6-12 months. Wi-Fi usage pattern affects life significantly -- the lock goes into a low-power sleep between operations.
- **Fingerprint-Resistant Touchscreen:** The capacitive touchscreen has a matte coating designed to obscure fingerprint smudges, making it harder for observers to determine frequently used digits.
- **Auto-Lock:** Configurable auto-lock timer (15 seconds to 4 minutes). The lock includes a door position sensor (similar to DoorSense) on some models to prevent locking when the door is open.
- **Offline Operation:** The touchscreen keypad works regardless of Wi-Fi or phone availability. Access codes are stored locally on the lock (up to 100 codes).
- **Supply Issues:** The Encode Plus experienced significant supply shortages in 2022-2023, with prices inflated on secondary markets due to demand driven by Apple Home Key exclusivity.

## Similar Devices

- **yale-assure-lock-2** -- Yale's modular lock with Thread/Matter option
- **august-wifi-smart-lock** -- August's retrofit Wi-Fi lock (same company family via ASSA ABLOY)
- **level-lock-plus** -- Invisible Thread/Matter lock
- **kwikset-halo-touch** -- Wi-Fi lock with fingerprint reader
