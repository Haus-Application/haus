package discovery

import (
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// ndpPattern matches macOS `ndp -an` output lines:
//   fe80::1                         aa:bb:cc:dd:ee:ff  en0   permanent R
// Captures: IPv6 address, MAC address, interface name
var ndpPattern = regexp.MustCompile(`^([\w:]+)\s+([\w:]+)\s+(\w+)\s+`)

// linuxNeighPattern matches Linux `ip -6 neighbor show` output:
//   fe80::1 dev eth0 lladdr aa:bb:cc:dd:ee:ff REACHABLE
var linuxNeighPattern = regexp.MustCompile(`^([\w:]+)\s+dev\s+(\w+)\s+lladdr\s+([\w:]+)\s+`)

// scanIPv6Neighbors discovers IPv6 devices on the local network.
// On macOS: ping6 ff02::1 (multicast all-nodes) then parse ndp -an.
// On Linux: ip -6 neighbor show.
// I -- I get nervous about multicast, but Mother says it's fine as long
// as we use link-local scope. It's like shouting in the apartment, not
// shouting in the street.
func scanIPv6Neighbors(session *ScanSession) {
	iface := detectActiveInterface()
	if iface == "" {
		log.Println("[ipv6] Could not detect active interface. Mother's network is hiding again.")
		return
	}

	switch runtime.GOOS {
	case "darwin":
		scanIPv6Darwin(session, iface)
	case "linux":
		scanIPv6Linux(session, iface)
	default:
		log.Printf("[ipv6] Unsupported OS %q for IPv6 neighbor discovery.", runtime.GOOS)
	}
}

// scanIPv6Darwin uses ping6 + ndp on macOS.
func scanIPv6Darwin(session *ScanSession, iface string) {
	// Send multicast ping to populate neighbor table -- ff02::1 is all-nodes
	// on the link-local scope. It's like ringing every doorbell on the floor.
	log.Printf("[ipv6] Pinging ff02::1 on %s to populate neighbor table...", iface)
	exec.Command("ping6", "-c", "2", "-I", iface, "ff02::1").Run()

	out, err := exec.Command("ndp", "-an").Output()
	if err != nil {
		log.Printf("[ipv6] Failed to run ndp -an: %v", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	matched := 0
	for _, line := range lines {
		matches := ndpPattern.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}

		ipv6Addr := matches[1]
		mac := matches[2]
		entryIface := matches[3]

		// Only consider entries on our active interface
		if entryIface != iface {
			continue
		}

		// Skip incomplete or broadcast entries
		if mac == "(incomplete)" || mac == "ff:ff:ff:ff:ff:ff" {
			continue
		}

		// Normalize MAC to lowercase for consistent matching
		mac = strings.ToLower(mac)

		mergeIPv6(session, ipv6Addr, mac)
		matched++
	}

	log.Printf("[ipv6] Parsed %d IPv6 neighbors from ndp table on %s.", matched, iface)
}

// scanIPv6Linux uses `ip -6 neighbor show` on Linux.
func scanIPv6Linux(session *ScanSession, iface string) {
	log.Printf("[ipv6] Pinging ff02::1 on %s to populate neighbor table...", iface)
	exec.Command("ping", "-6", "-c", "2", "-I", iface, "ff02::1").Run()

	out, err := exec.Command("ip", "-6", "neighbor", "show").Output()
	if err != nil {
		log.Printf("[ipv6] Failed to run ip -6 neighbor show: %v", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	matched := 0
	for _, line := range lines {
		matches := linuxNeighPattern.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}

		ipv6Addr := matches[1]
		entryIface := matches[2]
		mac := matches[3]

		if entryIface != iface {
			continue
		}

		mac = strings.ToLower(mac)
		mergeIPv6(session, ipv6Addr, mac)
		matched++
	}

	log.Printf("[ipv6] Parsed %d IPv6 neighbors from ip -6 neighbor on %s.", matched, iface)
}

// mergeIPv6 associates an IPv6 address with an existing device (by MAC match)
// or creates a new device if no MAC match is found.
func mergeIPv6(session *ScanSession, ipv6Addr, mac string) {
	// First, try to find an existing device with this MAC
	for _, device := range session.Devices {
		if strings.EqualFold(device.MAC, mac) {
			if !containsString(device.IPv6, ipv6Addr) {
				device.IPv6 = append(device.IPv6, ipv6Addr)
			}
			return
		}
	}

	// Check if this IPv6 address is already a key in the devices map
	if _, ok := session.Devices[ipv6Addr]; ok {
		return
	}

	// No existing device found with this MAC -- create a new one keyed by IPv6
	if mac != "" {
		device := &Device{
			IP:       ipv6Addr,
			MAC:      mac,
			Category: CategoryUnknown,
			Metadata: make(map[string]string),
		}
		session.Devices[ipv6Addr] = device
	}
}

// containsString checks if a string slice contains a given value.
func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// detectActiveInterface finds the name of the first non-loopback, up
// interface that has an IPv4 address. Usually en0 on macOS, eth0 on Linux.
// Mother's favorite network interface -- the one she actually talks through.
func detectActiveInterface() string {
	switch runtime.GOOS {
	case "darwin":
		// On macOS, try to find the interface via route output
		out, err := exec.Command("route", "-n", "get", "default").Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "interface:") {
					return strings.TrimSpace(strings.TrimPrefix(line, "interface:"))
				}
			}
		}
		// Fallback
		return "en0"
	case "linux":
		out, err := exec.Command("ip", "route", "show", "default").Output()
		if err == nil {
			// Format: default via 192.168.1.1 dev eth0 ...
			fields := strings.Fields(string(out))
			for i, f := range fields {
				if f == "dev" && i+1 < len(fields) {
					return fields[i+1]
				}
			}
		}
		return "eth0"
	default:
		return ""
	}
}
