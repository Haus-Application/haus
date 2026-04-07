package discovery

import (
	"encoding/xml"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// arpPattern matches lines like: ? (192.168.1.1) at aa:bb:cc:dd:ee:ff on en0 ...
var arpPattern = regexp.MustCompile(`\((\d+\.\d+\.\d+\.\d+)\) at ([0-9a-f:]+)`)

// nmap XML types for port scan results.
type nmapRun struct {
	Hosts []nmapHost `xml:"host"`
}

type nmapHost struct {
	Status    nmapStatus     `xml:"status"`
	Addresses []nmapAddress  `xml:"address"`
	Hostnames []nmapHostname `xml:"hostnames>hostname"`
	Ports     []nmapPort     `xml:"ports>port"`
}

type nmapStatus struct {
	State string `xml:"state,attr"`
}

type nmapAddress struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
	Vendor   string `xml:"vendor,attr"`
}

type nmapHostname struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type nmapPort struct {
	Protocol string      `xml:"protocol,attr"`
	PortID   int         `xml:"portid,attr"`
	State    nmapState   `xml:"state"`
	Service  nmapService `xml:"service"`
}

type nmapState struct {
	State string `xml:"state,attr"`
}

type nmapService struct {
	Name    string `xml:"name,attr"`
	Product string `xml:"product,attr"`
}

// scanHostsARP uses ping sweep + ARP to find all live hosts on the subnet.
// This finds way more devices than nmap -sn without root (~50 vs ~15).
func scanHostsARP(session *ScanSession, subnet string) {
	pingSweep(subnet)

	out, err := exec.Command("arp", "-a").Output()
	if err != nil {
		log.Printf("[arp] Failed to run arp -a: %v", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		matches := arpPattern.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		ip := matches[1]
		mac := matches[2]

		if mac == "ff:ff:ff:ff:ff:ff" || mac == "(incomplete)" {
			continue
		}
		if strings.HasPrefix(ip, "224.") || strings.HasPrefix(ip, "239.") || strings.HasSuffix(ip, ".255") {
			continue
		}
		if !strings.HasPrefix(ip, subnet+".") {
			continue
		}
		if ip == session.localIP {
			continue
		}

		hostname := ""
		if idx := strings.Index(line, " ("); idx > 0 {
			hostname = strings.TrimSpace(line[:idx])
			if hostname == "?" {
				hostname = ""
			}
		}

		device := &Device{
			IP:       ip,
			MAC:      mac,
			Hostname: hostname,
			Category: CategoryUnknown,
			Metadata: make(map[string]string),
		}

		session.Devices[ip] = device
		emitDevice(session, device)
	}

	log.Printf("[arp] Found %d hosts via ping sweep + ARP.", len(session.Devices))
}

// pingSweep sends a quick ping to every IP in the /24 subnet to populate
// the system ARP cache. Fires all concurrently with a short timeout.
func pingSweep(subnet string) {
	log.Printf("[arp] Ping sweeping %s.0/24...", subnet)
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50)

	for i := 1; i < 255; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(ip string) {
			defer wg.Done()
			defer func() { <-sem }()
			exec.Command("ping", "-c", "1", "-W", "300", ip).Run()
		}(fmt.Sprintf("%s.%d", subnet, i))
	}

	wg.Wait()
}

// scanPortsNmap uses nmap for fast, accurate port scanning on discovered hosts.
// Falls back to the Go port scanner if nmap fails.
func scanPortsNmap(session *ScanSession) {
	var ips []string
	for ip := range session.Devices {
		ips = append(ips, ip)
	}
	if len(ips) == 0 {
		return
	}

	portsStr := joinPorts(targetPorts)
	// -Pn: skip host discovery (we already know they're up from ARP).
	// Without this, nmap without root can't do ARP pings and declares hosts "down".
	args := []string{"-Pn", "-T4", "-p", portsStr, "--open", "-oX", "-"}
	args = append(args, ips...)

	log.Printf("[nmap] Port scanning %d hosts on ports %s...", len(ips), portsStr)
	out, err := exec.Command("nmap", args...).Output()
	if err != nil {
		log.Printf("[nmap] Port scan failed: %v. Falling back to Go scanner.", err)
		scanPorts(session)
		return
	}

	var result nmapRun
	if err := xml.Unmarshal(out, &result); err != nil {
		log.Printf("[nmap] Failed to parse XML: %v. Falling back to Go scanner.", err)
		scanPorts(session)
		return
	}

	portsFound := 0
	for _, host := range result.Hosts {
		var ip string
		for _, addr := range host.Addresses {
			if addr.AddrType == "ipv4" {
				ip = addr.Addr
				break
			}
		}
		if ip == "" {
			continue
		}

		device, ok := session.Devices[ip]
		if !ok {
			continue
		}

		for _, port := range host.Ports {
			if port.State.State == "open" {
				device.OpenPorts = append(device.OpenPorts, port.PortID)
				portsFound++
			}
		}
		emitDevice(session, device)
	}

	log.Printf("[nmap] Found %d open ports. Mother is impressed!", portsFound)
}

// joinPorts converts a port list to nmap's comma-separated format.
func joinPorts(ports []int) string {
	strs := make([]string, len(ports))
	for i, p := range ports {
		strs[i] = strconv.Itoa(p)
	}
	return strings.Join(strs, ",")
}
