package discovery

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

// serviceTypeMap maps known mDNS service types to device properties.
// I've cataloged every service type on Mother's network. Some of them
// I discovered at 3 AM when I couldn't sleep.
var serviceTypeMap = map[string]struct {
	DeviceType string
	Category   Category
	Protocol   string
}{
	"_hue._tcp":                {DeviceType: "hue_bridge", Category: CategoryLighting, Protocol: ""},
	"_huesync._tcp":            {DeviceType: "hue_sync", Category: CategoryMedia, Protocol: ""},
	"_brilliant._tcp":          {DeviceType: "brilliant_switch", Category: CategorySmartHome, Protocol: ""},
	"_googlecast._tcp":         {DeviceType: "", Category: CategoryMedia, Protocol: ""},
	"_airplay._tcp":            {DeviceType: "", Category: "", Protocol: "airplay"},
	"_spotify-connect._tcp":    {DeviceType: "", Category: "", Protocol: "spotify"},
	"_meshcop._udp":            {DeviceType: "thread_border_router", Category: "", Protocol: "thread"},
	"_nv_shield_remote._tcp":   {DeviceType: "shield_tv", Category: CategoryMedia, Protocol: ""},
	"_jellyfishV2._tcp":        {DeviceType: "jellyfish", Category: CategoryLighting, Protocol: ""},
	"_pvs6._tcp":               {DeviceType: "solar_gateway", Category: CategoryEnergy, Protocol: ""},
}

// knownServiceTypes are the service types we explicitly browse for.
// These are the ones that matter for device identification. Rather than
// relying solely on _services._dns-sd._udp (which can be flaky with some
// mDNS libraries), we browse these directly.
var knownServiceTypes = []string{
	"_hue._tcp",
	"_huesync._tcp",
	"_brilliant._tcp",
	"_googlecast._tcp",
	"_airplay._tcp",
	"_spotify-connect._tcp",
	"_meshcop._udp",
	"_nv_shield_remote._tcp",
	"_jellyfishV2._tcp",
	"_pvs6._tcp",
	"_raop._tcp",
	"_http._tcp",
}

// scanMDNS browses all known service types concurrently for speed.
// 12 types in parallel at 2s each = ~2s total instead of ~18s sequential.
func scanMDNS(session *ScanSession) {
	var wg sync.WaitGroup

	// Browse all known service types concurrently
	for _, svcType := range knownServiceTypes {
		wg.Add(1)
		go func(st string) {
			defer wg.Done()
			browseService(session, st)
		}(svcType)
	}

	wg.Wait()
	log.Printf("[mdns] Finished browsing %d service types.", len(knownServiceTypes))
}

// browseService browses a specific mDNS service type and enriches matching devices.
func browseService(session *ScanSession, svcType string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return
	}

	entries := make(chan *zeroconf.ServiceEntry, 64)

	go func() {
		for entry := range entries {
			processEntry(session, svcType, entry)
		}
	}()

	// Clean up service type for browsing -- remove domain suffix if present
	cleanType := strings.TrimSuffix(svcType, ".local.")
	cleanType = strings.TrimSuffix(cleanType, ".")

	err = resolver.Browse(ctx, cleanType, "local.", entries)
	if err != nil {
		return
	}

	<-ctx.Done()
}

// processEntry matches an mDNS entry to a device in our session and
// enriches it with service information. If the device wasn't found via
// ARP, we create a new entry — mDNS knows things ARP doesn't.
func processEntry(session *ScanSession, svcType string, entry *zeroconf.ServiceEntry) {
	// Find matching device by IP, or create a new one
	var device *Device
	var deviceIP string
	for _, ip := range entry.AddrIPv4 {
		deviceIP = ip.String()
		if d, ok := session.Devices[deviceIP]; ok {
			device = d
			break
		}
	}

	if device == nil && deviceIP != "" {
		// New device discovered via mDNS that wasn't in ARP table
		device = &Device{
			IP:       deviceIP,
			Category: CategoryUnknown,
			Metadata: make(map[string]string),
		}
		session.Devices[deviceIP] = device
	}

	if device == nil {
		return
	}

	// Add service to device's service list
	serviceName := strings.TrimSuffix(svcType, ".local.")
	serviceName = strings.TrimSuffix(serviceName, ".")
	if !hasService(device, serviceName) {
		device.Services = append(device.Services, serviceName)
	}

	// Apply known service type mappings
	for knownType, mapping := range serviceTypeMap {
		if matchesServiceType(serviceName, knownType) {
			if mapping.DeviceType != "" && device.DeviceType == "" {
				device.DeviceType = mapping.DeviceType
			}
			if mapping.Category != "" && (device.Category == CategoryUnknown || device.Category == "") {
				device.Category = mapping.Category
			}
			if mapping.Protocol != "" && !hasProtocol(device, mapping.Protocol) {
				device.Protocols = append(device.Protocols, mapping.Protocol)
			}
		}
	}

	// Use mDNS instance name if device has no name yet
	if device.Name == "" && entry.Instance != "" {
		device.Name = cleanMDNSName(entry.Instance)
	}

	emitDevice(session, device)
}

// cleanMDNSName strips backslash escaping from mDNS instance names.
// zeroconf and mDNS protocols escape spaces and special chars with backslashes.
func cleanMDNSName(name string) string {
	// Remove backslash escapes: "Hue\ Bridge" -> "Hue Bridge"
	result := strings.ReplaceAll(name, "\\", "")
	return strings.TrimSpace(result)
}

// hasService checks if a device already has a service in its list.
func hasService(device *Device, service string) bool {
	for _, s := range device.Services {
		if s == service {
			return true
		}
	}
	return false
}

// matchesServiceType checks if a discovered service type matches a known type.
func matchesServiceType(discovered, known string) bool {
	// Normalize: strip leading underscores, trailing dots, .local. suffix
	normalize := func(s string) string {
		s = strings.TrimSuffix(s, ".local.")
		s = strings.TrimSuffix(s, ".")
		s = strings.ToLower(s)
		return s
	}
	d := normalize(discovered)
	k := normalize(known)
	return d == k || strings.Contains(d, k) || strings.Contains(k, d)
}
