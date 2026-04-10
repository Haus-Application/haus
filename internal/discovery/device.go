package discovery

// Category represents the functional category of a discovered device.
type Category string

const (
	CategoryLighting  Category = "lighting"
	CategoryMedia     Category = "media"
	CategorySmartHome Category = "smart_home"
	CategoryEnergy    Category = "energy"
	CategoryNetwork   Category = "network"
	CategoryCompute   Category = "compute"
	CategoryUnknown   Category = "unknown"
)

// Device represents a discovered network device with all the details
// Mother could ever want to know about it.
type Device struct {
	IP           string            `json:"ip"`             // primary (IPv4 preferred)
	IPv6         []string          `json:"ipv6,omitempty"` // additional IPv6 addresses
	MAC          string            `json:"mac"`
	Hostname     string            `json:"hostname"`
	Name         string            `json:"name"`
	Manufacturer string            `json:"manufacturer"`
	Model        string            `json:"model"`
	DeviceType   string            `json:"device_type"`
	Category     Category          `json:"category"`
	Protocols    []string          `json:"protocols"`
	Services     []string          `json:"services"`
	OpenPorts    []int             `json:"open_ports"`
	Metadata     map[string]string `json:"metadata"`
}
