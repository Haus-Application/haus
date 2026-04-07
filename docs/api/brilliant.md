# Brilliant Smart Home Switch API

## Overview

Brilliant smart home control panels communicate via a **proprietary protocol on port 5455**. The protocol is not publicly documented. Brilliant devices are discovered via mDNS as `_brilliant._tcp`.

## Discovery

### mDNS

Service type: `_brilliant._tcp`

Example mDNS record:
```
Instance: 01663ad46e010003271243acea345b26
Port: 5455
TXT: home_id=01919f35f09c00024384dccbd839c147
```

The `home_id` links multiple Brilliant panels in the same home.

## Known Capabilities

Based on the Brilliant product line:
- **Light dimming** — built-in dimmer for wired lights
- **Scene control** — trigger scenes across connected devices
- **Motion detection** — built-in motion sensor
- **Intercom** — audio/video between panels
- **Music playback** — built-in speaker with Sonos support
- **Smart home control** — integrates with Ring, Sonos, Honeywell, SmartThings

## Protocol Details

The Brilliant protocol on port 5455 appears to be:
- Binary/protobuf-based (not JSON)
- Encrypted communication
- Requires cloud authentication via Brilliant app
- Local control may require pairing through the Brilliant mobile app first

## Integration Status

**Not yet controllable** — the Brilliant API is proprietary and not publicly documented. Potential approaches:

1. **Brilliant Cloud API** — may offer OAuth-based control (requires Brilliant account)
2. **Local protocol reverse engineering** — binary protocol on port 5455
3. **HomeKit bridge** — Brilliant supports HomeKit, which could be used as a control path
4. **Smart home hub integration** — Brilliant integrates with SmartThings, which has a documented API

## Haus Integration

- **Discovery:** mDNS `_brilliant._tcp` + TCP check on port 5455
- **Control:** Not yet implemented
- **Status:** Detected but not controllable
