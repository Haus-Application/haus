---
# ═══════════════════════════════════════════════════
# HAUS DEVICE KNOWLEDGE BASE
# ═══════════════════════════════════════════════════

# --- IDENTITY ---
id: "philips-hue-lightstrip-plus"
name: "Philips Hue Lightstrip Plus V4"
manufacturer: "Signify Netherlands B.V."
brand: "Philips Hue"
model: "LCL001"
model_aliases: ["929002269805", "555334", "555326", "800276", "LCL002", "LCL003"]
device_type: "hue_light"
category: "lighting"
product_line: "Hue"
release_year: 2021
discontinued: false
price_range: "$$"

# --- CONNECTIVITY ---
connectivity:
  mode: "local"
  local_api: false
  cloud_api: false
  cloud_required_for_setup: false
  internet_required: false
  local_only_capable: true
  protocols_spoken: ["zigbee", "bluetooth"]

# --- NETWORK FINGERPRINTS ---
# Zigbee device -- no direct IP network presence.
# Controlled exclusively via the Hue Bridge API.
network:
  mac_prefixes: []
  mdns_services: []
  mdns_txt_keys: []
  default_ports: []
  signature_ports: []
  ssdp_search_target: ""
  ssdp_server_string: ""
  hostname_patterns: []
  ip_ranges: []

# --- HTTP FINGERPRINTS ---
http_fingerprints: []

# --- HAUS INTEGRATION ---
integration:
  status: "supported"
  integration_key: "hue"
  polling_interval_sec: 5
  websocket_event: "hue:state"
  setup_type: "link_button"
  ai_chattable: true
  haus_milestone: "M3"

# --- CAPABILITIES ---
capabilities:
  - "on_off"
  - "brightness"
  - "color"
  - "color_temp"

# --- PROTOCOL ---
# No direct protocol -- accessed via bridge CLIP v2 API.
protocol:
  type: "https_rest"
  port: 443
  transport: "HTTPS"
  encoding: "JSON"
  auth_method: "link_button"
  auth_detail: "Accessed indirectly via Hue Bridge API. Bridge handles Zigbee communication."
  base_url_template: "https://{bridge_ip}/clip/v2/resource/light/{id}"
  tls: true
  tls_self_signed: true

# --- PHYSICAL ---
physical:
  form_factor: "strip"
  power_source: "mains"
  mounting: "wall"
  indoor_outdoor: "indoor"
  wireless_radios: ["zigbee", "bluetooth_le"]

# --- LINKS ---
links:
  product_page: "https://www.philips-hue.com/en-us/p/hue-white-and-color-ambiance-lightstrip-plus-base-v4-80-inch/046677555337"
  api_docs: "https://developers.meethue.com/develop/hue-api-v2/"
  developer_portal: "https://developers.meethue.com/"
  support: "https://www.philips-hue.com/en-us/support"
  community_forum: "https://developers.meethue.com/forum/"
  image_url: ""
  fcc_id: "2ABA6-LCL001"

# --- TAGS ---
tags: ["zigbee", "color", "ambiance", "lightstrip", "led_strip", "bluetooth_provisioning", "via_bridge", "gradient_capable"]
---

# Philips Hue Lightstrip Plus V4

## What It Is

The Philips Hue Lightstrip Plus V4 is a flexible adhesive-backed LED strip that produces tunable white light (2000K-6500K) and 16 million colors. The base kit includes an 80-inch (2 meter) strip with an integrated power supply and controller. It can be extended with 40-inch (1 meter) extension strips up to a total of 33 feet (10 meters). The V4 revision (model LCL001) introduced Bluetooth LE for direct app control alongside Zigbee, matching the rest of the modern Hue lineup.

The strip communicates over Zigbee 3.0 and requires a Hue Bridge for full functionality (scenes, automations, rooms, Entertainment API). Each strip appears as a single light resource in the bridge API -- the entire strip is one color/brightness at a time (unlike the Gradient Lightstrip which supports multi-zone color).

The Lightstrip Plus is commonly used for under-cabinet lighting, TV bias lighting, accent lighting behind furniture, and cove lighting installations.

## How Haus Discovers It

This device has **no direct network presence** -- it communicates exclusively over Zigbee via the Hue Bridge.

1. **Bridge Discovery** -- Haus first discovers and pairs with the Hue Bridge (see `philips-hue-bridge`).
2. **Device Enumeration** -- `GET /clip/v2/resource/device` lists all paired devices. The lightstrip appears with `product_data.model_id` of "LCL001" (or "LCL002" for the outdoor variant, "LCL003" for extensions).
3. **Light Resource** -- The strip exposes a single `light` service. `GET /clip/v2/resource/light` returns its state with the same structure as a bulb.
4. **Archetype Detection** -- The `metadata.archetype` field is typically `"hue_lightstrip"`, distinguishing it from bulbs.

## Pairing / Authentication

No separate pairing with Haus is required. The lightstrip is paired to the Hue Bridge via the Hue app's Zigbee pairing flow.

### Adding to the Bridge

1. Plug in the lightstrip power supply.
2. In the Hue app, tap "Add light" -- the bridge scans for new Zigbee devices.
3. The strip joins the bridge's Zigbee network within ~30 seconds.
4. Assign the strip to a room. It then appears in Haus via the bridge API.

### Factory Reset

- Power-cycle the strip 5 times (off 5 seconds, on 8 seconds each cycle). On the 5th power-on, the strip will blink to confirm reset.
- Alternatively, use a Hue Dimmer Switch held close to the strip's controller (hold all 4 buttons).

## API Reference

All control is via the Hue Bridge CLIP v2 API. See `philips-hue-bridge` for full endpoint documentation.

### Get Light State

```
GET /clip/v2/resource/light/{light_id}
```

**Response fields for a Lightstrip Plus:**
```json
{
  "id": "d4e5f6a7-...",
  "type": "light",
  "metadata": {
    "name": "TV Backlight",
    "archetype": "hue_lightstrip"
  },
  "on": {"on": true},
  "dimming": {
    "brightness": 80.0,
    "min_dim_level": 0.2
  },
  "color": {
    "xy": {"x": 0.3127, "y": 0.3290},
    "gamut": {
      "red":   {"x": 0.6915, "y": 0.3083},
      "green": {"x": 0.1700, "y": 0.7000},
      "blue":  {"x": 0.1532, "y": 0.0475}
    },
    "gamut_type": "C"
  },
  "color_temperature": {
    "mirek": 250,
    "mirek_valid": true,
    "mirek_schema": {"mirek_minimum": 153, "mirek_maximum": 500}
  }
}
```

### Control Light

```
PUT /clip/v2/resource/light/{light_id}
```

Same body format as any Hue light. See `philips-hue-bridge` for full control documentation.

### Entertainment API

The Lightstrip Plus is a popular choice for the Hue Entertainment API, which enables low-latency streaming of color data for media sync. When used in an Entertainment area, the strip receives color updates at up to 25Hz via the bridge's UDP streaming protocol (port 2100, using the `clientkey` from pairing).

## AI Capabilities

When the AI interacts with this lightstrip (via the bridge), it can:

- **Turn on/off** by name ("turn off the TV backlight")
- **Set brightness** ("dim the kitchen strip to 20%")
- **Set color** by name ("make the accent light purple")
- **Set color temperature** ("set the cabinet lights to warm white")
- **Include in scenes** -- the strip participates in room scenes alongside bulbs
- **Report state** ("The TV backlight is on, blue, at 80% brightness")

## Quirks & Notes

- **Single Zone:** Unlike the Hue Gradient Lightstrip, the Lightstrip Plus V4 is a single zone -- the entire strip is one uniform color at a time. For multi-zone gradient effects, see the Gradient Lightstrip (model LST004).
- **Extension Compatibility:** V4 extensions (LCL003) work with V4 base kits. V3 extensions are NOT compatible with V4 bases and vice versa -- the connector changed between generations.
- **Cut Points:** The strip can be cut at marked intervals (approximately every 13 inches / 33 cm). Cut sections cannot be reconnected without third-party connectors.
- **Power Limitations:** Maximum total length is 10 meters (base + extensions). Exceeding this causes brightness drop at the far end due to voltage loss.
- **Adhesive:** The 3M adhesive backing works best on clean, smooth surfaces. On textured surfaces, mounting clips are recommended.
- **Color Gamut C:** Same wide gamut as the A19 Gen 3+ bulbs, supporting the full range of Hue colors.
- **Lumen Output:** Approximately 1600 lumens for the 2-meter base strip at full white brightness. Color output varies by hue.
- **Wattage:** 20W for the 2-meter base strip.
- **Zigbee Repeater:** As a mains-powered device, the lightstrip controller acts as a Zigbee router in the mesh network.
- **Heat:** The controller module generates moderate heat during extended full-brightness use. Ensure adequate ventilation around the controller, not the strip itself.

## Similar Devices

- **philips-hue-bridge** -- Required hub for full functionality
- **philips-hue-bulb-a19** -- Standard bulb with the same color capabilities
- **philips-hue-sync-box** -- HDMI sync device that pairs well with lightstrips for TV ambient lighting
