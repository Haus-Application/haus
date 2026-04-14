// Package kb loads and indexes the device knowledge base under docs/devices/.
// Each file is a markdown document with YAML frontmatter describing one device.
// The catalog lets runtime code match discovered devices to their rich KB entry.
package kb

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Device is a single knowledge base entry.
type Device struct {
	Slug              string              // filename without .md (e.g. "philips-hue-bridge")
	ID                string              // frontmatter id
	Name              string              // display name
	Manufacturer      string              // legal entity
	Brand             string              // consumer brand
	Model             string              // primary model number
	DeviceType        string              // haus internal type
	Category          string              // lighting | media | security | climate | ...
	IntegrationKey    string              // matches DeviceProbeResult.Integration
	IntegrationStatus string              // supported | read_only | detected_only | planned | not_feasible
	Capabilities      []string            // on_off | brightness | color | ...
	Protocol          ProtocolInfo        // protocol details
	Network           NetworkFingerprints // discovery hints
	Body              string              // markdown body (everything after closing ---)
	Raw               string              // full file contents
}

// ProtocolInfo mirrors the `protocol:` block in frontmatter.
type ProtocolInfo struct {
	Type          string // https_rest | tcp_xor | websocket_json | ...
	Port          int
	Transport     string // TCP | HTTP | HTTPS | WebSocket | UDP | TLS
	Encoding      string // JSON | Protobuf | XML | XOR-JSON | ...
	AuthMethod    string // none | api_key | basic_auth | oauth2 | link_button | ...
	BaseURL       string // e.g. "https://{ip}/clip/v2"
	TLS           bool
	TLSSelfSigned bool
}

// NetworkFingerprints mirrors the `network:` block in frontmatter.
type NetworkFingerprints struct {
	MACPrefixes      []string
	MDNSServices     []string
	DefaultPorts     []int
	SignaturePorts   []int
	HostnamePatterns []string
}

// Catalog is the indexed knowledge base.
type Catalog struct {
	All              []*Device
	BySlug           map[string]*Device
	ByIntegrationKey map[string][]*Device
	ByDeviceType     map[string][]*Device
	LoadedAt         time.Time
}

// Load scans a directory for *.md files (skipping _template.md and README.md),
// parses frontmatter for each, and returns an indexed catalog.
func Load(dir string) (*Catalog, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read kb dir %q: %w", dir, err)
	}

	c := &Catalog{
		BySlug:           make(map[string]*Device),
		ByIntegrationKey: make(map[string][]*Device),
		ByDeviceType:     make(map[string][]*Device),
		LoadedAt:         time.Now(),
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		// Skip template and index.
		if name == "_template.md" || name == "README.md" || strings.HasPrefix(name, "_") {
			continue
		}

		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("[kb] skip %s: %v", name, err)
			continue
		}

		slug := strings.TrimSuffix(name, ".md")
		dev, err := parseDevice(slug, data)
		if err != nil {
			log.Printf("[kb] skip %s: %v", name, err)
			continue
		}
		c.All = append(c.All, dev)
		c.BySlug[slug] = dev
		if dev.IntegrationKey != "" {
			c.ByIntegrationKey[dev.IntegrationKey] = append(c.ByIntegrationKey[dev.IntegrationKey], dev)
		}
		if dev.DeviceType != "" {
			c.ByDeviceType[dev.DeviceType] = append(c.ByDeviceType[dev.DeviceType], dev)
		}
	}

	log.Printf("[kb] Loaded %d devices, %d integration keys indexed", len(c.All), len(c.ByIntegrationKey))
	return c, nil
}

// Match finds the best KB entry for a runtime-discovered device.
// Priority:
//  1. Exact integration_key match (prefer one whose device_type also matches)
//  2. Exact device_type match
//  3. Fuzzy manufacturer+model match (case-insensitive Contains on either side)
//  4. nil if nothing matches
func (c *Catalog) Match(integration, deviceType, manufacturer, model string) *Device {
	if c == nil {
		return nil
	}

	// 1. integration_key exact
	if integration != "" {
		if matches := c.ByIntegrationKey[integration]; len(matches) > 0 {
			// Prefer one whose device_type also matches (substring either way)
			if deviceType != "" {
				dtLow := strings.ToLower(deviceType)
				for _, d := range matches {
					dtKB := strings.ToLower(d.DeviceType)
					if dtKB == dtLow || strings.Contains(dtKB, dtLow) || strings.Contains(dtLow, dtKB) {
						return d
					}
				}
			}
			return matches[0]
		}
	}

	// 2. device_type exact
	if deviceType != "" {
		if matches := c.ByDeviceType[deviceType]; len(matches) > 0 {
			return matches[0]
		}
	}

	// 3. Fuzzy manufacturer+model
	if manufacturer != "" || model != "" {
		mfr := strings.ToLower(manufacturer)
		mdl := strings.ToLower(model)
		for _, d := range c.All {
			dmfr := strings.ToLower(d.Manufacturer + " " + d.Brand)
			dmdl := strings.ToLower(d.Model)
			mfrHit := mfr != "" && dmfr != "" && (strings.Contains(dmfr, mfr) || strings.Contains(mfr, d.Manufacturer) && d.Manufacturer != "")
			mdlHit := mdl != "" && dmdl != "" && (strings.Contains(dmdl, mdl) || strings.Contains(mdl, dmdl))
			if mfrHit && mdlHit {
				return d
			}
			if mfrHit && model == "" {
				return d
			}
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Frontmatter parser (hand-rolled, YAML-ish subset)
// -----------------------------------------------------------------------------
//
// We only need to extract the specific fields we use. A full YAML parser would
// be overkill and would add a dep. Our frontmatter is machine-generated from a
// single template, so the shape is predictable.

func parseDevice(slug string, data []byte) (*Device, error) {
	raw := string(data)
	fmBytes, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, err
	}

	d := &Device{Slug: slug, Raw: raw, Body: body}
	fm := parseFrontmatter(string(fmBytes))

	d.ID = fm.Scalar("id")
	d.Name = fm.Scalar("name")
	d.Manufacturer = fm.Scalar("manufacturer")
	d.Brand = fm.Scalar("brand")
	d.Model = fm.Scalar("model")
	d.DeviceType = fm.Scalar("device_type")
	d.Category = fm.Scalar("category")
	d.Capabilities = fm.List("capabilities")

	d.IntegrationKey = fm.NestedScalar("integration", "integration_key")
	d.IntegrationStatus = fm.NestedScalar("integration", "status")

	d.Protocol.Type = fm.NestedScalar("protocol", "type")
	d.Protocol.Transport = fm.NestedScalar("protocol", "transport")
	d.Protocol.Encoding = fm.NestedScalar("protocol", "encoding")
	d.Protocol.AuthMethod = fm.NestedScalar("protocol", "auth_method")
	d.Protocol.BaseURL = fm.NestedScalar("protocol", "base_url_template")
	d.Protocol.Port = atoiSafe(fm.NestedScalar("protocol", "port"))
	d.Protocol.TLS = boolSafe(fm.NestedScalar("protocol", "tls"))
	d.Protocol.TLSSelfSigned = boolSafe(fm.NestedScalar("protocol", "tls_self_signed"))

	d.Network.MACPrefixes = fm.NestedList("network", "mac_prefixes")
	d.Network.MDNSServices = fm.NestedList("network", "mdns_services")
	d.Network.DefaultPorts = intList(fm.NestedList("network", "default_ports"))
	d.Network.SignaturePorts = intList(fm.NestedList("network", "signature_ports"))
	d.Network.HostnamePatterns = fm.NestedList("network", "hostname_patterns")

	return d, nil
}

// splitFrontmatter separates the YAML frontmatter from the markdown body.
// The frontmatter is the block between a leading "---\n" and the next "---\n".
func splitFrontmatter(data []byte) (fm []byte, body string, err error) {
	// Require leading "---\n" or "---\r\n"
	if !bytes.HasPrefix(data, []byte("---\n")) && !bytes.HasPrefix(data, []byte("---\r\n")) {
		return nil, string(data), fmt.Errorf("no frontmatter")
	}

	// Skip past opening fence
	rest := data
	rest = rest[bytes.IndexByte(rest, '\n')+1:]

	// Find closing fence
	closeIdx := -1
	scanner := bufio.NewScanner(bytes.NewReader(rest))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	offset := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			closeIdx = offset
			break
		}
		// +1 for the newline we stripped in Text()
		offset += len(line) + 1
	}
	if closeIdx < 0 {
		return nil, "", fmt.Errorf("no closing frontmatter fence")
	}

	fm = rest[:closeIdx]
	// Body starts after the closing fence line
	bodyStart := closeIdx
	// advance past the "---" line (including its newline)
	if bodyStart < len(rest) {
		nl := bytes.IndexByte(rest[bodyStart:], '\n')
		if nl >= 0 {
			bodyStart += nl + 1
		}
	}
	if bodyStart <= len(rest) {
		body = string(rest[bodyStart:])
	}
	return fm, body, nil
}

// fmap is the parsed frontmatter. Values can be:
//   - string (scalar)
//   - []string (inline list like `[a, b, c]`)
//   - map[string]any (nested block)
type fmap map[string]any

func (f fmap) Scalar(key string) string {
	v, ok := f[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (f fmap) List(key string) []string {
	v, ok := f[key]
	if !ok {
		return nil
	}
	if list, ok := v.([]string); ok {
		return list
	}
	return nil
}

func (f fmap) NestedScalar(outer, inner string) string {
	v, ok := f[outer]
	if !ok {
		return ""
	}
	m, ok := v.(fmap)
	if !ok {
		return ""
	}
	return m.Scalar(inner)
}

func (f fmap) NestedList(outer, inner string) []string {
	v, ok := f[outer]
	if !ok {
		return nil
	}
	m, ok := v.(fmap)
	if !ok {
		return nil
	}
	return m.List(inner)
}

// parseFrontmatter handles a narrow YAML subset matching the device template:
//   - top-level `key: value` scalars
//   - top-level `key: [a, b, c]` inline lists
//   - top-level `key: []` empty inline lists
//   - top-level `key:` followed by indented `  - item` block lists
//   - top-level `key:` followed by indented `  subkey: value` nested maps
//   - comments (# ...) are ignored
//   - empty lines are ignored
func parseFrontmatter(src string) fmap {
	root := fmap{}
	lines := strings.Split(src, "\n")

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		indent := leadingSpaces(line)

		// Only process top-level (indent 0) lines at the outer loop.
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || indent != 0 {
			i++
			continue
		}

		k, v, ok := splitKV(trimmed)
		if !ok {
			i++
			continue
		}
		v = stripComment(v)

		// Three cases based on what follows the colon:
		if v == "" {
			// Could be a nested map OR a block list. Peek ahead.
			i++
			nested, listItems, consumed := readBlock(lines, i)
			i += consumed
			if len(listItems) > 0 {
				root[k] = append([]string(nil), listItems...)
			} else if len(nested) > 0 {
				root[k] = nested
			} else {
				root[k] = ""
			}
			continue
		}

		// Inline list: [a, b, c]
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			root[k] = parseInlineList(v)
			i++
			continue
		}

		// Scalar
		root[k] = unquote(v)
		i++
	}

	return root
}

// readBlock reads indented lines starting at `start`. It returns either a
// nested map (if it sees `key: value` lines) or a list (if it sees `- item`
// lines). The caller is positioned at `start`; we return how many lines were
// consumed.
func readBlock(lines []string, start int) (nested fmap, list []string, consumed int) {
	nested = fmap{}
	i := start
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		indent := leadingSpaces(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			i++
			continue
		}
		if indent == 0 {
			// Back to top level — end of block.
			break
		}

		if strings.HasPrefix(trimmed, "- ") || trimmed == "-" {
			item := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
			item = stripComment(item)
			item = unquote(item)
			if item != "" {
				list = append(list, item)
			}
			i++
			continue
		}

		// Nested key: value (supports one more level of nesting, ignored deeper)
		k, v, ok := splitKV(trimmed)
		if !ok {
			i++
			continue
		}
		v = stripComment(v)
		if v == "" {
			// Sub-block (list of items or deeper map). Collect list items only;
			// we don't need deeper nesting for our schema.
			i++
			for i < len(lines) {
				sub := lines[i]
				subTrim := strings.TrimSpace(sub)
				subIndent := leadingSpaces(sub)
				if subTrim == "" || strings.HasPrefix(subTrim, "#") {
					i++
					continue
				}
				if subIndent <= indent {
					break
				}
				if strings.HasPrefix(subTrim, "- ") || subTrim == "-" {
					item := strings.TrimSpace(strings.TrimPrefix(subTrim, "-"))
					item = stripComment(item)
					item = unquote(item)
					// Store the raw item for deeper structures (we don't
					// parse http_fingerprints map entries, but we preserve them).
					if item != "" {
						// attach under key as string list
						prev, _ := nested[k].([]string)
						nested[k] = append(prev, item)
					}
					i++
				} else {
					// ignore deeper nested maps
					i++
				}
			}
			continue
		}
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			nested[k] = parseInlineList(v)
		} else {
			nested[k] = unquote(v)
		}
		i++
	}
	consumed = i - start
	return nested, list, consumed
}

func leadingSpaces(s string) int {
	n := 0
	for _, r := range s {
		if r == ' ' {
			n++
		} else if r == '\t' {
			n += 2
		} else {
			break
		}
	}
	return n
}

func splitKV(s string) (k, v string, ok bool) {
	idx := strings.Index(s, ":")
	if idx < 0 {
		return "", "", false
	}
	return strings.TrimSpace(s[:idx]), strings.TrimSpace(s[idx+1:]), true
}

func stripComment(s string) string {
	// only strip `# ...` that's not inside quotes
	inQuote := byte(0)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inQuote != 0 {
			if c == inQuote && (i == 0 || s[i-1] != '\\') {
				inQuote = 0
			}
			continue
		}
		if c == '"' || c == '\'' {
			inQuote = c
			continue
		}
		if c == '#' {
			return strings.TrimSpace(s[:i])
		}
	}
	return strings.TrimSpace(s)
}

func unquote(s string) string {
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}

func parseInlineList(s string) []string {
	inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, "["), "]"))
	if inner == "" {
		return nil
	}
	// Split on commas that aren't inside quotes
	var out []string
	var cur strings.Builder
	inQuote := byte(0)
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		if inQuote != 0 {
			if c == inQuote {
				inQuote = 0
			}
			cur.WriteByte(c)
			continue
		}
		if c == '"' || c == '\'' {
			inQuote = c
			cur.WriteByte(c)
			continue
		}
		if c == ',' {
			out = append(out, unquote(strings.TrimSpace(cur.String())))
			cur.Reset()
			continue
		}
		cur.WriteByte(c)
	}
	if cur.Len() > 0 {
		out = append(out, unquote(strings.TrimSpace(cur.String())))
	}
	// Drop empty entries
	filtered := out[:0]
	for _, s := range out {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func atoiSafe(s string) int {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return n
}

func boolSafe(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "yes", "y", "1":
		return true
	}
	return false
}

func intList(ss []string) []int {
	out := make([]int, 0, len(ss))
	for _, s := range ss {
		if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			out = append(out, n)
		}
	}
	return out
}
