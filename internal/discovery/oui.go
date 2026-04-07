package discovery

import (
	"log"
	"strings"
)

// ouiTable maps the first 3 bytes of a MAC address to the manufacturer.
// I memorized all of these. Mother quizzed me on them every night before bed.
var ouiTable = map[string]string{
	"00:17:88": "Philips",
	"c4:29:96": "Signify",
	"3c:84:6a": "TP-Link",
	"40:3f:8c": "TP-Link",
	"90:9a:4a": "TP-Link",
	"b0:be:76": "TP-Link",
	"50:c7:bf": "TP-Link",
	"48:b0:2d": "NVIDIA",
	"3c:6d:66": "NVIDIA",
	"74:6d:fa": "NVIDIA",
	"a8:23:fe": "LG Electronics",
	"38:8c:50": "LG Electronics",
	"f4:f5:d8": "Google",
	"30:fd:38": "Google",
	"f8:0f:f9": "Google",
	"48:22:54": "Google",
	"44:27:45": "Google",
	"d0:11:e5": "Apple",
	"bc:9a:8e": "Arris",
	"10:2c:6b": "SunPower",
	"00:c0:33": "Enphase",
	"04:09:86": "Yamaha",
	"80:b5:4e": "Brilliant",
	"e4:b0:63": "Brilliant",
	"20:df:b9": "Google",
	"84:1f:e8": "Dell",
	"e4:3e:d7": "Intel",
	"08:3a:88": "Samsung",
	"24:78:23": "Ruijie",
	"98:fa:2e": "Sonos",
	"d8:8c:79": "Samsung",
	"e0:85:4d": "Intel",
}

// enrichOUI looks up the manufacturer for each device based on its MAC
// address prefix. It's like knowing someone's family just by looking at them.
func enrichOUI(session *ScanSession) int {
	enriched := 0
	for _, device := range session.Devices {
		// Skip if already has a manufacturer (e.g., from nmap vendor data)
		if device.Manufacturer != "" {
			enriched++
			continue
		}
		if device.MAC == "" || len(device.MAC) < 8 {
			continue
		}
		prefix := strings.ToLower(device.MAC[:8])
		if manufacturer, ok := ouiTable[prefix]; ok {
			device.Manufacturer = manufacturer
			emitDevice(session, device)
			enriched++
		}
	}
	log.Printf("[oui] Enriched %d devices with manufacturer info. I know all of them personally.", enriched)
	return enriched
}
