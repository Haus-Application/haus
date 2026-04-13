---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "kwikset-halo-touch"
name: "Kwikset Halo Touch WiFi Fingerprint Smart Lock"
manufacturer: "Kwikset (Spectrum Brands)"
brand: "Kwikset"
model: "99590-001"
model_aliases: ["99590-003", "99590-004", "959"]
device_type: "smart_lock"
category: "security"
product_line: "Kwikset Halo"
release_year: 2020
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
  protocols_spoken: ["wifi", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes:
    - "EC:FA:5C"        # Kwikset / Spectrum Brands
    - "C0:49:EF"        # Espressif (WiFi module used in Halo)
  mdns_services: []     # Halo does not advertise mDNS services
  mdns_txt_keys: []
  default_ports: []     # No open ports
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns:
    - "^Kwikset"
    - "^kwikset-halo"
    - "^ESP_"           # Some units show Espressif default hostname
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []   # No local HTTP services

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "kwikset"
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
  auth_detail: "Kwikset uses a proprietary cloud API accessed through the Kwikset App. OAuth 2.0 authentication via Spectrum Brands identity service. No official public API or developer program. Some community reverse engineering exists but the API is not well documented."
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
  product_page: "https://www.kwikset.com/halo-touch"
  api_docs: ""
  developer_portal: ""
  support: "https://www.kwikset.com/support"
  community_forum: ""
  image_url: ""
  fcc_id: "2AZAH-99590"

# --- TAGS ---
tags: ["smart_lock", "wifi", "fingerprint", "biometric", "kwikset", "spectrum_brands", "smartkey", "cloud_dependent", "no_official_api"]
---

# Kwikset Halo Touch WiFi Fingerprint Smart Lock

## What It Is

The Kwikset Halo Touch is a Wi-Fi smart deadbolt with an integrated capacitive fingerprint reader, manufactured by Kwikset, a brand of Spectrum Brands (formerly Hardware & Home Improvement). The fingerprint reader is the headline feature -- it stores up to 100 fingerprints (for up to 50 users) and unlocks in under a second, providing keyless entry without needing a phone, code, or key. The lock connects directly to Wi-Fi (2.4GHz) for remote access via the Kwikset app, with no bridge or hub required. It features Kwikset's patented SmartKey Security cylinder, which allows re-keying the lock in seconds without a locksmith. The lock has ANSI/BHMA Grade 2 residential security certification and is available in various trim styles (Contemporary, Traditional). Powered by four AA batteries.

## How Haus Discovers It

1. **OUI Match** -- Devices with MAC prefix `EC:FA:5C` (Kwikset/Spectrum Brands) or `C0:49:EF` (Espressif, the Wi-Fi module manufacturer) are flagged during network scanning. The Espressif OUI is shared with many IoT devices, so it is not conclusive alone.
2. **Hostname Pattern** -- Kwikset locks may appear with hostnames starting with `Kwikset` or `ESP_` (Espressif default) in DHCP tables.
3. **No Local Discovery** -- The Halo Touch does not advertise mDNS, SSDP, or any local services. It communicates outbound to Kwikset's cloud only.

## Pairing / Authentication

Kwikset does not provide an official public API or developer program.

### Kwikset App Setup

1. **Account Creation:** Create a Kwikset account in the mobile app.
2. **BLE Pairing:** The app discovers the lock via BLE during initial setup.
3. **Wi-Fi Provisioning:** The app sends Wi-Fi credentials to the lock over BLE.
4. **Cloud Registration:** The lock connects to Kwikset's cloud and registers with the user's account.
5. **Fingerprint Enrollment:** Fingerprints are enrolled directly on the lock by repeatedly placing a finger on the sensor (8-10 placements per fingerprint for a complete template). Enrollment is done via the lock's physical interface with audio/LED guidance, not through the app.

### Fingerprint System

- The lock uses a capacitive fingerprint sensor (not optical), which is more secure against spoofing.
- Up to 100 fingerprints across 50 user slots.
- Fingerprint templates are stored locally on the lock's secure element -- they are never transmitted to the cloud.
- The sensor is rated for outdoor use with moisture resistance.
- False rejection rate varies with finger condition (wet, dry, dirty fingers may fail).

### Security Notes

- No official API means no sanctioned third-party integration path.
- The Kwikset cloud API has had minimal reverse-engineering compared to August/Ring.
- Fingerprint data is stored exclusively on the lock hardware, not in the cloud or app.
- SmartKey re-keying: the lock can be rekeyed to match any Kwikset-compatible key in seconds using the included SmartKey tool.

## API Reference

### Kwikset Cloud API (Unofficial)

The Kwikset app communicates with Kwikset's cloud infrastructure. Limited reverse-engineering has revealed:

- **Authentication:** OAuth 2.0 flow via Spectrum Brands identity service.
- **Device Listing:** Returns associated locks with status (locked/unlocked), battery level, and recent activity.
- **Lock/Unlock:** Remote lock/unlock commands routed through the cloud.
- **Access Code Management:** Virtual access codes can be managed for temporary/scheduled access.
- **Activity Log:** Access history with timestamps and method (fingerprint, app, manual key).

Detailed endpoint documentation is not publicly available. The API is significantly less well-documented by the community compared to August or Ring APIs.

### No Local API

The Halo Touch has no local API. All remote operations go through Kwikset's cloud. The fingerprint sensor operates independently on the lock hardware and does not require network connectivity.

## AI Capabilities

AI integration is not currently planned due to limited API access. If implemented, the AI concierge could:

- Report lock/unlock status
- Lock/unlock the door remotely with security confirmation
- Report battery level
- Show recent access log with access method (fingerprint, app, key, auto-lock)
- Report number of enrolled fingerprints per user

## Quirks & Notes

- **Fingerprint Performance:** The capacitive fingerprint sensor is generally fast (under 1 second) and accurate, but performance degrades with wet, oily, or heavily calloused fingers. Cold weather can also affect sensor sensitivity. Users report needing to enroll the same finger multiple times (different angles) for reliable recognition.
- **No Matter/Thread:** The Halo Touch does not support Matter, Thread, HomeKit, Zigbee, or Z-Wave. It is Wi-Fi + BLE only with cloud-dependent remote access.
- **SmartKey Security:** Kwikset's patented SmartKey cylinder uses a sidebar mechanism (instead of traditional pin tumblers) that is resistant to bump keys and forced entry. The lock can be re-keyed to any Kwikset KW1 key using the SmartKey tool -- useful when moving into a new home or changing keys.
- **ANSI/BHMA Grade 2:** Residential security rating. Lower than Schlage's Grade 1 but standard for residential smart locks.
- **2.4GHz Only:** The lock only supports 2.4GHz Wi-Fi (802.11 b/g/n). This is typical for battery-powered IoT devices but can be inconvenient on networks that separate 2.4/5GHz SSIDs.
- **Battery Life:** Four AA batteries, advertised 6-9 months. The fingerprint sensor and Wi-Fi radio are the primary power consumers. The lock provides low-battery warnings via push notification and audible beeping.
- **Espressif Wi-Fi Module:** The lock uses an Espressif (ESP32 or ESP8266 family) Wi-Fi module internally. This is why some units show up with Espressif OUI prefixes rather than Kwikset OUIs.
- **No Auto-Unlock:** Unlike August and Level locks, the Halo Touch does not have geofencing-based auto-unlock. The fingerprint sensor is intended to be the primary hands-free entry method instead.
- **Cloud Dependency:** If Kwikset's cloud goes down or Spectrum Brands discontinues the service, remote lock/unlock and app features would stop working. The fingerprint sensor, physical key, and manual thumbturn continue to function offline.
- **Limited Smart Home Integration:** The Halo Touch works with Amazon Alexa and Google Assistant for voice control (lock only, not unlock for security reasons). No HomeKit support. No IFTTT. Limited integration ecosystem compared to competitors.

## Similar Devices

- **august-wifi-smart-lock** -- Wi-Fi + BLE retrofit lock with better third-party API support
- **yale-assure-lock-2** -- Modular lock with Thread/Matter and touchscreen keypad
- **schlage-encode-plus** -- Wi-Fi lock with Apple Home Key
- **level-lock-plus** -- Invisible Thread/Matter lock (no keypad or fingerprint)
