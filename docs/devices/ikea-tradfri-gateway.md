---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "ikea-tradfri-gateway"
name: "IKEA TRADFRI Gateway"
manufacturer: "IKEA of Sweden"
brand: "IKEA TRADFRI"
model: "E1526"
model_aliases: ["TRADFRI Gateway", "IKEA Gateway"]
device_type: "tradfri_gateway"
category: "smart_home"
product_line: "TRADFRI"
release_year: 2017
discontinued: true
price_range: "$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: true
  cloud_api: true
  cloud_required_for_setup: true
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["zigbee", "ethernet", "wifi"]

# --- NETWORK FINGERPRINTS ---
network:
  mac_prefixes: ["94:B9:7E", "CC:50:E3", "DC:EF:CA", "60:01:94", "B4:E1:EB"]
  mdns_services: []
  mdns_txt_keys: []
  default_ports: [5684, 80]
  signature_ports: [5684]
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: ["^TRADFRI-Gateway-[a-f0-9]+$", "^GW-[a-f0-9]+$"]
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "planned"
  integration_key: "tradfri"
  polling_interval_sec: 10
  websocket_event: "tradfri:state"
  setup_type: "api_key"
  ai_chattable: false
  haus_milestone: "M11"

# --- CAPABILITIES ---
capabilities: ["on_off", "brightness", "color", "color_temp", "scenes", "groups"]

# --- PROTOCOL ---
protocol:
  type: "coap"
  port: 5684
  transport: "UDP"
  encoding: "CBOR"
  auth_method: "api_key"
  auth_detail: "Pre-shared key derived from security code on gateway label. Initial handshake exchanges the 16-character security code for a per-identity PSK via CoAP POST to /15011/9063."
  base_url_template: "coaps://{ip}:5684"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "gateway"
  power_source: "mains"
  mounting: "shelf"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee"]

# --- LINKS ---
links:
  product_page: "https://www.ikea.com/us/en/p/tradfri-gateway-white-00337813/"
  api_docs: "https://github.com/glenndehaan/ikea-tradfri-coap-docs"
  developer_portal: ""
  support: "https://www.ikea.com/us/en/customer-service/"
  community_forum: "https://github.com/home-assistant/core/tree/dev/homeassistant/components/tradfri"
  image_url: ""
  fcc_id: "2AHFL-E1526"

# --- TAGS ---
tags: ["zigbee_hub", "coap", "dtls", "ikea", "deprecated", "replaced_by_dirigera"]
---

# IKEA TRADFRI Gateway

## What It Is

The IKEA TRADFRI Gateway (model E1526) is a Zigbee coordinator that serves as the central hub for IKEA's TRADFRI smart home ecosystem. It connects to the home network via Ethernet and bridges communication between the IKEA Home Smart app (or local API clients) and TRADFRI Zigbee devices including bulbs, LED drivers, blinds, outlets, motion sensors, and remote controls. It was one of the most affordable smart home hubs on the market at roughly $30. IKEA has discontinued the TRADFRI Gateway in favor of the newer DIRIGERA hub, but millions of units remain deployed.

## How Haus Discovers It

1. **OUI match**: The gateway's Ethernet MAC address will match IKEA OUI prefixes (94:B9:7E, CC:50:E3, DC:EF:CA, 60:01:94, or B4:E1:EB — all registered to IKEA of Sweden).
2. **Hostname pattern**: The gateway typically advertises itself via DHCP with a hostname matching `TRADFRI-Gateway-*` or `GW-*`.
3. **Port probe**: A CoAP DTLS handshake attempt on UDP port 5684 will confirm the gateway is speaking CoAPs. This is the signature port — very few consumer devices listen on 5684.
4. **CoAP probe**: A CoAP GET to `coaps://{ip}:5684/15001` (with a valid PSK) returns the device list, confirming it is a TRADFRI gateway.

Note: The TRADFRI Gateway does NOT advertise mDNS or SSDP services, which makes it less discoverable than Hue bridges. OUI + port 5684 is the primary detection path.

## Pairing / Authentication

The TRADFRI Gateway uses DTLS with a Pre-Shared Key (PSK) for all CoAP communication. The pairing flow is:

1. **Locate the Security Code**: A 16-character alphanumeric code is printed on the label on the bottom of the gateway. This is the "master key."
2. **Generate an identity**: Send a CoAP POST to `/15011/9063` using the security code as the PSK identity `Client_identity` and the security code as the PSK:
   ```json
   {"9090": "my_haus_identity"}
   ```
3. **Receive a per-identity PSK**: The gateway responds with a new PSK tied to the chosen identity:
   ```json
   {"9091": "generated_psk_string", "9029": "1.17.0044"}
   ```
4. **Store the PSK**: All subsequent CoAP requests use the identity + generated PSK for DTLS authentication. The security code is never used again.
5. **Verify**: CoAP GET to `coaps://{ip}:5684/15001` should return the full device list.

The gateway supports a maximum of roughly 10 PSK identities simultaneously.

## API Reference

The TRADFRI Gateway speaks CoAP (RFC 7252) over DTLS 1.2 (RFC 6347) on UDP port 5684. All payloads use a combination of CBOR encoding and JSON-like integer-keyed objects derived from the LWM2M / IPSO Smart Object specification.

### Key Endpoints

| Path | Method | Description |
|------|--------|-------------|
| `/15001` | GET | List all devices (returns array of device IDs) |
| `/15001/{id}` | GET | Get device state |
| `/15001/{id}` | PUT | Set device state |
| `/15004` | GET | List all groups |
| `/15004/{id}` | GET | Get group state |
| `/15004/{id}` | PUT | Set group state |
| `/15005` | GET | List all scenes (moods) |
| `/15005/{group_id}/{scene_id}` | PUT | Activate a scene |
| `/15006` | GET | List notifications |
| `/15010` | GET | List smart tasks (timers/schedules) |
| `/15011/9063` | POST | Generate new PSK identity |
| `/15011/15012` | PUT | Gateway reboot |
| `/15011/9034` | GET | Gateway info (firmware, uptime) |

### Device State Object Keys

| Key | Type | Description |
|-----|------|-------------|
| 3311 | array | Light resource array |
| 5706 | string | Color hex (e.g., "f1e0b5") |
| 5707 | int | Hue (0-65279) |
| 5708 | int | Saturation (0-65279) |
| 5709 | int | Color X (CIE) |
| 5710 | int | Color Y (CIE) |
| 5711 | int | Color temperature (mireds, 250-454) |
| 5712 | int | Transition time (ms / 100) |
| 5850 | int | On/Off (0 or 1) |
| 5851 | int | Brightness (0-254) |
| 9001 | string | Device name |
| 9002 | int | Created timestamp |
| 9003 | int | Instance ID |

### Example: Turn On a Light at 50% Brightness

```
CoAP PUT coaps://{ip}:5684/15001/65537
Payload: {"3311": [{"5850": 1, "5851": 127}]}
```

### Example: Set Color Temperature to Warm White

```
CoAP PUT coaps://{ip}:5684/15001/65537
Payload: {"3311": [{"5711": 400}]}
```

### Observation (Push Updates)

CoAP supports the Observe option (RFC 7641). By sending a GET with the Observe flag on a device or group endpoint, the gateway will push state updates when changes occur. This reduces polling overhead significantly.

```
CoAP GET (Observe) coaps://{ip}:5684/15001/65537
```

## AI Capabilities

AI chat integration is not planned for the TRADFRI Gateway since the device is discontinued and being replaced by DIRIGERA. If added, the AI could control lights (on/off, brightness, color temperature), activate scenes, and manage groups through the CoAP API.

## Quirks & Notes

- **Discontinued**: IKEA stopped selling the TRADFRI Gateway in 2023-2024, replacing it with the DIRIGERA hub. Firmware updates may cease.
- **PSK limit**: The gateway supports roughly 10 simultaneous PSK identities. If the limit is reached, old identities must be removed via the IKEA app or a factory reset.
- **CoAP library required**: Go does not have CoAP/DTLS in the standard library. The `plgd-dev/go-coap` and `pion/dtls` libraries are the most mature Go options.
- **Integer-keyed JSON**: The API uses numeric keys from the OMA LWM2M spec, not human-readable field names. A mapping layer is essential.
- **Rate limiting**: The gateway can become unresponsive if flooded with CoAP requests. A maximum of roughly 1 request per 100ms is safe.
- **Firmware OTA**: The gateway updates its own firmware and child device firmware via IKEA's cloud. There is no local firmware update path.
- **No mDNS**: Unlike most modern hubs, the TRADFRI Gateway does not advertise any mDNS services, making OUI matching critical for discovery.
- **Zigbee mesh**: The gateway acts as the Zigbee coordinator. TRADFRI mains-powered devices (bulbs, outlets) act as Zigbee routers, extending the mesh.

## Similar Devices

- [ikea-dirigera-hub](ikea-dirigera-hub.md) — The replacement hub with Matter support
- [philips-hue-bridge](philips-hue-bridge.md) — Similar hub-based lighting approach, but with a REST/EventStream API
- [lutron-caseta-bridge](lutron-caseta-bridge.md) — Another bridge-based lighting system with a local API
