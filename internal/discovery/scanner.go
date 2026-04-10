package discovery

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/coalson/haus/internal/db"
	"github.com/google/uuid"
)

// ScanEvent is sent over the session's event channel so the API layer
// can stream progress to the frontend via SSE.
type ScanEvent struct {
	Type string      `json:"type"` // "stage", "device", "complete"
	Data interface{} `json:"data"`
}

// StageEvent describes a scan stage starting or finishing.
type StageEvent struct {
	Stage   string `json:"stage"`
	Status  string `json:"status"` // "running", "complete"
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}

// ScanSession tracks a single scan run. Devices are keyed by IP so we
// can enrich them across multiple stages without duplicates.
type ScanSession struct {
	ID      string
	Status  string // "running", "complete", "error"
	Devices map[string]*Device
	Events  chan ScanEvent
	Started time.Time
	localIP string // our own IP, so we can skip it in results
}

// Scanner orchestrates network discovery. It remembers the auto-detected
// subnet and keeps a map of active/completed scan sessions.
type Scanner struct {
	subnet string
	db     *sql.DB
	mu     sync.Mutex
	scans  map[string]*ScanSession
}

// NewScanner creates a Scanner with an auto-detected subnet from the
// host's network interfaces. I -- I get very anxious when I can't find
// a valid interface, but I try my best.
func NewScanner(database *sql.DB) *Scanner {
	subnet := detectSubnet()
	if subnet != "" {
		log.Printf("[scanner] Auto-detected subnet: %s.0/24. Mother's network looks healthy.", subnet)
	} else {
		log.Println("[scanner] WARNING: Could not auto-detect subnet. Mother's network is hiding from me.")
	}
	return &Scanner{
		subnet: subnet,
		db:     database,
		scans:  make(map[string]*ScanSession),
	}
}

// Subnet returns the auto-detected subnet prefix (e.g. "192.168.1").
// Returns an empty string if detection failed. Mother's network is shy sometimes.
func (s *Scanner) Subnet() string {
	return s.subnet
}

// detectSubnet finds the first non-loopback, up interface with an IPv4
// address and returns the first 3 octets (assuming /24).
func detectSubnet() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}
			return fmt.Sprintf("%d.%d.%d", ip[0], ip[1], ip[2])
		}
	}
	return ""
}

// StartScan creates a new scan session and launches the scan goroutine.
// If subnet is empty, uses the auto-detected one.
func (s *Scanner) StartScan(subnet string) *ScanSession {
	if subnet == "" {
		subnet = s.subnet
	}

	session := &ScanSession{
		ID:      uuid.New().String(),
		Status:  "running",
		Devices: make(map[string]*Device),
		Events:  make(chan ScanEvent, 256),
		Started: time.Now(),
		localIP: detectLocalIP(subnet),
	}

	s.mu.Lock()
	s.scans[session.ID] = session
	s.mu.Unlock()

	go s.run(session, subnet)
	return session
}

// GetSession retrieves a scan session by ID.
func (s *Scanner) GetSession(id string) *ScanSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.scans[id]
}

// detectLocalIP finds our own IP on the given subnet.
func detectLocalIP(subnet string) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}
			if strings.HasPrefix(ip.String(), subnet+".") {
				return ip.String()
			}
		}
	}
	return ""
}

// emitStage sends a stage event to the session channel.
func emitStage(session *ScanSession, stage, status, message string, count int) {
	session.Events <- ScanEvent{
		Type: "stage",
		Data: StageEvent{
			Stage:   stage,
			Status:  status,
			Message: message,
			Count:   count,
		},
	}
}

// emitDevice sends a device discovery/update event to the session channel.
func emitDevice(session *ScanSession, device *Device) {
	session.Events <- ScanEvent{
		Type: "device",
		Data: device,
	}
}

// run executes the full scan pipeline, stage by stage, just like Mother
// taught me: methodical, thorough, and slightly obsessive.
func (s *Scanner) run(session *ScanSession, subnet string) {
	defer close(session.Events)

	log.Printf("[scanner] Mother is scanning the network (%s.0/24)... I'm so excited!", subnet)

	// Stage 1: Host discovery (ping sweep + ARP)
	emitStage(session, "arp", "running", fmt.Sprintf("Discovering hosts on %s.0/24...", subnet), 0)
	scanHostsARP(session, subnet)
	emitStage(session, "arp", "complete", fmt.Sprintf("Found %d hosts", len(session.Devices)), len(session.Devices))

	// Stage 1b: IPv6 neighbor discovery
	emitStage(session, "ipv6", "running", "Discovering IPv6 neighbors...", 0)
	scanIPv6Neighbors(session)
	ipv6Count := 0
	for _, d := range session.Devices {
		ipv6Count += len(d.IPv6)
	}
	emitStage(session, "ipv6", "complete", fmt.Sprintf("Found %d IPv6 addresses", ipv6Count), ipv6Count)

	// Stage 2: OUI enrichment
	emitStage(session, "oui", "running", "Looking up manufacturers...", 0)
	identified := enrichOUI(session)
	emitStage(session, "oui", "complete", fmt.Sprintf("Identified %d manufacturers", identified), identified)

	// Stage 3: Port scanning via nmap (fast + accurate)
	emitStage(session, "ports", "running", fmt.Sprintf("Scanning ports on %d hosts (nmap)...", len(session.Devices)), 0)
	scanPortsNmap(session)
	portsFound := countTotalOpenPorts(session)
	emitStage(session, "ports", "complete", fmt.Sprintf("Found %d open ports", portsFound), portsFound)

	// Stage 4: Kasa probe
	kasaBefore := countByProtocol(session, "kasa")
	emitStage(session, "kasa", "running", "Probing Kasa devices...", 0)
	probeKasa(session)
	kasaFound := countByProtocol(session, "kasa") - kasaBefore
	emitStage(session, "kasa", "complete", fmt.Sprintf("Found %d Kasa devices", kasaFound), kasaFound)

	// Stage 5: Cast probe
	castBefore := countByProtocol(session, "cast")
	emitStage(session, "cast", "running", "Probing Cast devices...", 0)
	probeCast(session)
	castFound := countByProtocol(session, "cast") - castBefore
	emitStage(session, "cast", "complete", fmt.Sprintf("Found %d Cast devices", castFound), castFound)

	// Stage 6: mDNS scan
	mdnsBefore := len(session.Devices)
	emitStage(session, "mdns", "running", "Browsing mDNS services...", 0)
	scanMDNS(session)
	mdnsNew := len(session.Devices) - mdnsBefore
	svcCount := countTotalServices(session)
	emitStage(session, "mdns", "complete", fmt.Sprintf("Found %d services, %d new devices", svcCount, mdnsNew), svcCount)

	// Stage 7: Classification
	emitStage(session, "classify", "running", "Classifying devices...", 0)
	classify(session)
	emitStage(session, "classify", "complete", fmt.Sprintf("Classified %d devices", len(session.Devices)), len(session.Devices))

	// Persist all discovered devices to the database.
	if s.db != nil {
		persisted := 0
		for _, device := range session.Devices {
			protos, _ := json.Marshal(device.Protocols)
			svcs, _ := json.Marshal(device.Services)
			ports, _ := json.Marshal(device.OpenPorts)
			meta, _ := json.Marshal(device.Metadata)
			if err := db.UpsertDevice(s.db, device.IP, device.MAC, device.Hostname, device.Name,
				device.Manufacturer, device.Model, device.DeviceType, string(device.Category),
				string(protos), string(svcs), string(ports), string(meta)); err != nil {
				log.Printf("[scanner] Failed to persist device %s: %v", device.IP, err)
			} else {
				persisted++
			}
		}
		log.Printf("[scanner] Persisted %d devices to database.", persisted)
	}

	duration := time.Since(session.Started)
	session.Status = "complete"

	// Final complete event
	session.Events <- ScanEvent{
		Type: "complete",
		Data: map[string]interface{}{
			"device_count": len(session.Devices),
			"duration_ms":  duration.Milliseconds(),
		},
	}

	log.Printf("[scanner] Scan complete! Found %d devices in %v. I found them all, Mother!", len(session.Devices), duration)
}

// countTotalOpenPorts counts the total open ports across all devices.
func countTotalOpenPorts(session *ScanSession) int {
	count := 0
	for _, d := range session.Devices {
		count += len(d.OpenPorts)
	}
	return count
}

// countByProtocol counts devices that have a given protocol in their list.
func countByProtocol(session *ScanSession, proto string) int {
	count := 0
	for _, d := range session.Devices {
		for _, p := range d.Protocols {
			if p == proto {
				count++
				break
			}
		}
	}
	return count
}

// countTotalServices counts the total number of mDNS services across all devices.
func countTotalServices(session *ScanSession) int {
	count := 0
	for _, d := range session.Devices {
		count += len(d.Services)
	}
	return count
}

// classify applies final categorization rules to all discovered devices.
func classify(session *ScanSession) {
	networkManufacturers := map[string]bool{
		"Arris":   true,
		"Netgear": true,
		"Cisco":   true,
		"Ubiquiti": true,
	}

	for _, device := range session.Devices {
		// If category already set by a probe stage, keep it
		if device.Category != CategoryUnknown && device.Category != "" {
			goto finalize
		}

		// Manufacturer-based classification
		switch device.Manufacturer {
		case "TP-Link":
			device.Category = CategoryLighting
		case "Philips", "Signify":
			device.Category = CategoryLighting
			if device.DeviceType == "" {
				device.DeviceType = "hue_bridge"
			}
		case "Brilliant":
			device.Category = CategorySmartHome
			if device.DeviceType == "" {
				device.DeviceType = "brilliant_switch"
			}
		case "SunPower":
			device.Category = CategoryEnergy
			if device.DeviceType == "" {
				device.DeviceType = "solar_gateway"
			}
		case "Enphase":
			device.Category = CategoryEnergy
			if device.DeviceType == "" {
				device.DeviceType = "solar_gateway"
			}
		case "Yamaha":
			device.Category = CategoryMedia
			if device.DeviceType == "" {
				device.DeviceType = "av_receiver"
			}
		case "Sonos":
			device.Category = CategoryMedia
			if device.DeviceType == "" {
				device.DeviceType = "speaker"
			}
		case "Google":
			if !hasProtocol(device, "cast") && device.Category == CategoryUnknown {
				device.Category = CategorySmartHome
				if device.DeviceType == "" {
					device.DeviceType = "nest_device"
				}
			}
		case "Apple":
			device.Category = CategoryCompute
		case "NVIDIA":
			device.Category = CategoryMedia
			if device.DeviceType == "" {
				device.DeviceType = "shield_tv"
			}
		case "LG Electronics":
			device.Category = CategoryMedia
			if device.DeviceType == "" {
				device.DeviceType = "tv"
			}
		case "Samsung":
			// Samsung could be TV or phone — check ports
			if hasPort(device, 8008) || hasPort(device, 8009) {
				device.Category = CategoryMedia
			}
		case "Dell", "Intel":
			device.Category = CategoryCompute
		default:
			// Common compute hostname patterns
			lower := strings.ToLower(device.Hostname)
			if strings.Contains(lower, "macbook") || strings.Contains(lower, "imac") ||
				strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad") {
				device.Category = CategoryCompute
			}
			// Network equipment manufacturers
			if networkManufacturers[device.Manufacturer] {
				device.Category = CategoryNetwork
			}
		}

		if device.Category != CategoryUnknown && device.Category != "" {
			goto finalize
		}

	finalize:
		// Ensure device has a name
		if device.Name == "" {
			if device.Manufacturer != "" && device.Manufacturer != "Unknown" {
				parts := strings.Split(device.IP, ".")
				lastOctet := parts[len(parts)-1]
				device.Name = fmt.Sprintf("%s .%s", device.Manufacturer, lastOctet)
			} else if device.Hostname != "" {
				device.Name = device.Hostname
			} else {
				parts := strings.Split(device.IP, ".")
				lastOctet := parts[len(parts)-1]
				device.Name = fmt.Sprintf("Device .%s", lastOctet)
			}
		}

		// Ensure Metadata is initialized
		if device.Metadata == nil {
			device.Metadata = make(map[string]string)
		}

		// Ensure slice fields are not nil for clean JSON
		if device.Protocols == nil {
			device.Protocols = []string{}
		}
		if device.Services == nil {
			device.Services = []string{}
		}
		if device.OpenPorts == nil {
			device.OpenPorts = []int{}
		}

		// Default category
		if device.Category == "" {
			device.Category = CategoryUnknown
		}

		emitDevice(session, device)
	}
}
