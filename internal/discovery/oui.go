package discovery

import (
	_ "embed"
	"encoding/csv"
	"log"
	"strings"
	"sync"
)

// ouiDataCSV is the IEEE OUI (MA-L) registry, trimmed to OUI + organization name.
// Sourced from https://standards-oui.ieee.org/oui/oui.csv.
// ~39k entries, ~1.2MB. Buster insisted on memorizing every OUI in the world.
//
//go:embed oui_data.csv
var ouiDataCSV []byte

// ouiOverrides are hand-curated display names. They take priority over the IEEE
// data because they produce nicer names ("Signify" instead of
// "SIGNIFY NETHERLANDS B.V.") and cover a few vendors whose MA-L entries we want
// to rebrand for Haus. Keys MUST be lowercase, colon-separated (xx:xx:xx).
var ouiOverrides = map[string]string{
	"00:17:88": "Philips Hue",
	"ec:b5:fa": "Philips Hue",
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
	"20:df:b9": "Google",
	"7c:10:15": "Google",
	"b0:09:da": "Google",
	"18:7f:88": "Google",
	"18:b4:30": "Google Nest",
	"d0:11:e5": "Apple",
	"bc:9a:8e": "Arris",
	"10:2c:6b": "SunPower",
	"00:c0:33": "Enphase",
	"04:09:86": "Yamaha",
	"80:b5:4e": "Brilliant",
	"e4:b0:63": "Brilliant",
	"84:1f:e8": "Dell",
	"e4:3e:d7": "Intel",
	"08:3a:88": "Samsung",
	"d8:8c:79": "Samsung",
	"98:fa:2e": "Sonos",
	"e0:85:4d": "Intel",
}

// ouiTable is the full lookup table built from ouiOverrides + the embedded IEEE
// CSV. Populated lazily on first use so binary startup stays snappy.
var (
	ouiTable     map[string]string
	ouiTableOnce sync.Once
)

// loadOUITable parses the embedded IEEE CSV and merges the hand-curated
// overrides on top. Called once, guarded by sync.Once.
func loadOUITable() {
	table := make(map[string]string, 40000)

	r := csv.NewReader(strings.NewReader(string(ouiDataCSV)))
	r.FieldsPerRecord = -1 // tolerate trailing whitespace / odd rows
	// Skip header
	if _, err := r.Read(); err != nil {
		log.Printf("[oui] Failed to read embedded OUI header: %v. Mother is disappointed.", err)
	}
	loaded := 0
	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		if len(row) < 2 {
			continue
		}
		// row[0] is like "286FB9" — normalize to "28:6f:b9"
		oui := strings.ToLower(strings.TrimSpace(row[0]))
		if len(oui) != 6 {
			continue
		}
		key := oui[0:2] + ":" + oui[2:4] + ":" + oui[4:6]
		table[key] = strings.TrimSpace(row[1])
		loaded++
	}

	// Apply overrides — our nicer names win
	for k, v := range ouiOverrides {
		table[k] = v
	}

	log.Printf("[oui] Loaded %d OUI entries (IEEE) + %d overrides. Buster knows them all.", loaded, len(ouiOverrides))
	ouiTable = table
}

// LookupOUI returns the manufacturer for a MAC address prefix, if known.
// Accepts full MAC or just the 3-byte prefix. Case-insensitive.
// Returns "" if the prefix is not registered with IEEE.
func LookupOUI(mac string) string {
	if mac == "" || len(mac) < 8 {
		return ""
	}
	ouiTableOnce.Do(loadOUITable)
	prefix := strings.ToLower(mac[:8])
	return ouiTable[prefix]
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
		if manufacturer := LookupOUI(device.MAC); manufacturer != "" {
			device.Manufacturer = manufacturer
			emitDevice(session, device)
			enriched++
		}
	}
	log.Printf("[oui] Enriched %d devices with manufacturer info. I know all of them personally.", enriched)
	return enriched
}
