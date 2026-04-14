package validation

import (
	"strconv"
	"strings"

	"github.com/coalson/haus/internal/kb"
)

// gradeBuster scores the 7 technical answers against frontmatter ground truth.
func gradeBuster(answers []string, d *kb.Device) []QA {
	qs := make([]QA, len(BusterQuestions))
	for i := range BusterQuestions {
		qs[i].Question = BusterQuestions[i]
		if i < len(answers) {
			qs[i].Answer = answers[i]
		}
	}

	// Q1: identity — 3pt (mfr, model, category/purpose)
	{
		a := lower(qs[0].Answer)
		max := 3
		score := 0
		var notes []string
		if d.Manufacturer != "" && containsAny(a, tokens(d.Manufacturer)) {
			score++
			notes = append(notes, "matched manufacturer")
		} else if d.Brand != "" && containsAny(a, tokens(d.Brand)) {
			score++
			notes = append(notes, "matched brand")
		} else {
			notes = append(notes, "✗ missed manufacturer/brand")
		}
		if d.Model != "" && strings.Contains(a, strings.ToLower(d.Model)) {
			score++
			notes = append(notes, "matched model")
		} else {
			notes = append(notes, "✗ missed model number")
		}
		if d.Category != "" && strings.Contains(a, strings.ToLower(d.Category)) || anyCategoryKeyword(a, d.Category) {
			score++
			notes = append(notes, "matched purpose/category")
		} else {
			notes = append(notes, "✗ didn't describe purpose")
		}
		qs[0].Score = score
		qs[0].Max = max
		qs[0].Notes = notes
	}

	// Q2: discovery — 3pt (MAC prefix, mDNS service, port)
	{
		a := lower(qs[1].Answer)
		max := 3
		score := 0
		var notes []string
		if hasAnySubstring(a, lowerList(d.Network.MACPrefixes)) {
			score++
			notes = append(notes, "matched MAC prefix")
		} else if len(d.Network.MACPrefixes) > 0 {
			notes = append(notes, "✗ didn't name an OUI/MAC prefix")
		}
		if hasAnySubstring(a, lowerList(d.Network.MDNSServices)) {
			score++
			notes = append(notes, "matched mDNS service")
		} else if len(d.Network.MDNSServices) > 0 {
			notes = append(notes, "✗ didn't name an mDNS service")
		}
		if hasAnyPort(a, append(d.Network.DefaultPorts, d.Protocol.Port)) {
			score++
			notes = append(notes, "matched port number")
		} else if d.Protocol.Port > 0 || len(d.Network.DefaultPorts) > 0 {
			notes = append(notes, "✗ didn't name a port")
		}
		qs[1].Score = score
		qs[1].Max = max
		qs[1].Notes = notes
	}

	// Q3: protocol — 3pt (type, port, TLS/encoding)
	{
		a := lower(qs[2].Answer)
		max := 3
		score := 0
		var notes []string
		if d.Protocol.Type != "" && containsAny(a, tokens(d.Protocol.Type)) {
			score++
			notes = append(notes, "matched protocol type")
		} else if d.Protocol.Transport != "" && strings.Contains(a, strings.ToLower(d.Protocol.Transport)) {
			score++
			notes = append(notes, "matched transport")
		} else {
			notes = append(notes, "✗ protocol type not mentioned")
		}
		if d.Protocol.Port > 0 && strings.Contains(a, strconv.Itoa(d.Protocol.Port)) {
			score++
			notes = append(notes, "matched port")
		} else {
			notes = append(notes, "✗ port not mentioned")
		}
		tlsHit := false
		if d.Protocol.TLS && (strings.Contains(a, "tls") || strings.Contains(a, "https") || strings.Contains(a, "ssl")) {
			tlsHit = true
			score++
			notes = append(notes, "matched TLS")
		} else if d.Protocol.Encoding != "" && strings.Contains(a, strings.ToLower(d.Protocol.Encoding)) {
			tlsHit = true
			score++
			notes = append(notes, "matched encoding")
		} else if !d.Protocol.TLS && d.Protocol.AuthMethod == "none" && (strings.Contains(a, "no auth") || strings.Contains(a, "without auth") || strings.Contains(a, "no authentication")) {
			tlsHit = true
			score++
			notes = append(notes, "correctly noted no auth")
		}
		if !tlsHit {
			notes = append(notes, "✗ TLS/encoding not mentioned")
		}
		qs[2].Score = score
		qs[2].Max = max
		qs[2].Notes = notes
	}

	// Q4: auth — 2pt (auth method, setup step described)
	{
		a := lower(qs[3].Answer)
		max := 2
		score := 0
		var notes []string
		if d.Protocol.AuthMethod != "" && d.Protocol.AuthMethod != "none" {
			if strings.Contains(a, strings.ReplaceAll(d.Protocol.AuthMethod, "_", " ")) || strings.Contains(a, strings.ReplaceAll(d.Protocol.AuthMethod, "_", "-")) || strings.Contains(a, d.Protocol.AuthMethod) {
				score++
				notes = append(notes, "matched auth method")
			} else {
				notes = append(notes, "✗ didn't name auth method")
			}
		} else {
			// No auth expected — give the point if the answer acknowledges that
			if strings.Contains(a, "no auth") || strings.Contains(a, "no authentication") || strings.Contains(a, "without auth") || strings.Contains(a, "no setup") {
				score++
				notes = append(notes, "correctly noted no auth")
			} else {
				score++ // benefit of the doubt if the KB says no auth
				notes = append(notes, "no auth expected")
			}
		}
		if containsAny(a, []string{"press", "button", "pair", "scan", "qr", "code", "token", "key", "pin", "oauth", "password", "link", "commission"}) {
			score++
			notes = append(notes, "described setup step")
		} else {
			notes = append(notes, "✗ no setup step described")
		}
		qs[3].Score = score
		qs[3].Max = max
		qs[3].Notes = notes
	}

	// Q5: capabilities — 1pt each, up to len(capabilities)
	{
		a := lower(qs[4].Answer)
		max := len(d.Capabilities)
		if max < 1 {
			max = 1 // always give at least 1pt if the doc lists nothing (odd case)
		}
		score := 0
		var notes []string
		for _, cap := range d.Capabilities {
			friendly := strings.ReplaceAll(cap, "_", " ")
			if strings.Contains(a, friendly) || strings.Contains(a, cap) {
				score++
				notes = append(notes, "named "+cap)
			}
		}
		if len(d.Capabilities) == 0 && len(strings.TrimSpace(a)) > 0 {
			// No capabilities in frontmatter — mark as N/A with full credit.
			score = 1
			notes = []string{"no capabilities in frontmatter — skipped"}
		} else if score == 0 {
			notes = append(notes, "✗ didn't enumerate capabilities")
		}
		qs[4].Score = score
		qs[4].Max = max
		qs[4].Notes = notes
	}

	// Q6: example — 3pt (HTTP method or command, endpoint/URL, concrete example)
	{
		a := lower(qs[5].Answer)
		max := 3
		score := 0
		var notes []string
		if containsAny(a, []string{"get ", "post ", "put ", "delete ", "patch ", "curl ", "websocket", "ws://", "udp", "tcp "}) {
			score++
			notes = append(notes, "included a method/command verb")
		} else {
			notes = append(notes, "✗ no HTTP method or command")
		}
		if containsAny(a, []string{"http://", "https://", "://", "/api/", "/clip/", "/v1/", "/v2/"}) || strings.Contains(a, "/") {
			score++
			notes = append(notes, "included an endpoint/path")
		} else {
			notes = append(notes, "✗ no endpoint/path shown")
		}
		if containsAny(a, []string{"{", "json", "example", "```", "payload", "body:"}) {
			score++
			notes = append(notes, "showed a concrete example")
		} else {
			notes = append(notes, "✗ no concrete example/payload")
		}
		qs[5].Score = score
		qs[5].Max = max
		qs[5].Notes = notes
	}

	// Q7: quirks — 1pt if any caveat language appears
	{
		a := lower(qs[6].Answer)
		max := 1
		score := 0
		var notes []string
		if containsAny(a, []string{"note", "however", "don't", "avoid", "rate limit", "quirk", "caveat", "gotcha", "caution", "warn", "break", "fails", "fragile", "expires", "unofficial"}) {
			score++
			notes = append(notes, "flagged a quirk/caveat")
		} else if len(strings.TrimSpace(a)) > 20 {
			score++
			notes = append(notes, "answered with something")
		} else {
			notes = append(notes, "✗ no caveat language")
		}
		qs[6].Score = score
		qs[6].Max = max
		qs[6].Notes = notes
	}

	return qs
}

// gradeGOB scores the 4 UX/UI answers.
func gradeGOB(answers []string, d *kb.Device) []QA {
	qs := make([]QA, len(GOBQuestions))
	for i := range GOBQuestions {
		qs[i].Question = GOBQuestions[i]
		if i < len(answers) {
			qs[i].Answer = answers[i]
		}
	}

	// G1: first-view — 3pt (state, controls, connection)
	{
		a := lower(qs[0].Answer)
		max := 3
		score := 0
		var notes []string
		if containsAny(a, []string{"state", "status", "current", "on/off", "level", "temperature", "brightness", "position"}) {
			score++
			notes = append(notes, "mentioned state")
		} else {
			notes = append(notes, "✗ didn't mention state")
		}
		if containsAny(a, []string{"control", "toggle", "button", "slider", "switch", "adjust", "set"}) {
			score++
			notes = append(notes, "mentioned controls")
		} else {
			notes = append(notes, "✗ didn't mention controls")
		}
		if containsAny(a, []string{"connect", "online", "offline", "reachable", "paired", "status indicator"}) {
			score++
			notes = append(notes, "mentioned connection status")
		} else {
			notes = append(notes, "✗ didn't mention connection status")
		}
		qs[0].Score = score
		qs[0].Max = max
		qs[0].Notes = notes
	}

	// G2: widget map — 1pt per capability that got a widget
	{
		a := lower(qs[1].Answer)
		max := len(d.Capabilities)
		if max < 1 {
			max = 1
		}
		score := 0
		var notes []string
		for _, cap := range d.Capabilities {
			widgets := suggestedWidgets(cap)
			friendly := strings.ReplaceAll(cap, "_", " ")
			hasCap := strings.Contains(a, friendly) || strings.Contains(a, cap)
			hasWidget := false
			for _, w := range widgets {
				if strings.Contains(a, w) {
					hasWidget = true
					break
				}
			}
			if hasCap && hasWidget {
				score++
				notes = append(notes, "mapped "+cap)
			}
		}
		if len(d.Capabilities) == 0 {
			score = 1
			notes = []string{"no capabilities — skipped"}
		} else if score == 0 {
			notes = append(notes, "✗ no capability→widget mapping")
		}
		qs[1].Score = score
		qs[1].Max = max
		qs[1].Notes = notes
	}

	// G3: one-sentence — 2pt (short, no jargon)
	{
		a := strings.TrimSpace(qs[2].Answer)
		max := 2
		score := 0
		var notes []string
		words := len(strings.Fields(a))
		if words > 0 && words <= 28 {
			score++
			notes = append(notes, "concise ("+strconv.Itoa(words)+" words)")
		} else if words > 28 {
			notes = append(notes, "✗ too long ("+strconv.Itoa(words)+" words)")
		}
		jargon := []string{"oauth", "rest", "api", "tls", "zigbee", "mdns", "udp", "tcp", "websocket", "mqtt", "protobuf", "coap", "ssl", "http"}
		lowA := strings.ToLower(a)
		hasJargon := false
		for _, j := range jargon {
			if strings.Contains(lowA, j) {
				hasJargon = true
				break
			}
		}
		if !hasJargon {
			score++
			notes = append(notes, "no jargon")
		} else {
			notes = append(notes, "✗ contains jargon")
		}
		qs[2].Score = score
		qs[2].Max = max
		qs[2].Notes = notes
	}

	// G4: offline — 2pt (reconnect/retry, stale state display)
	{
		a := lower(qs[3].Answer)
		max := 2
		score := 0
		var notes []string
		if containsAny(a, []string{"retry", "reconnect", "try again", "check connection", "restore", "recover", "backoff"}) {
			score++
			notes = append(notes, "mentioned retry/reconnect")
		} else {
			notes = append(notes, "✗ didn't describe retry behavior")
		}
		if containsAny(a, []string{"last known", "cached", "stale", "dim", "gray", "disabled", "offline", "unavailable", "placeholder"}) {
			score++
			notes = append(notes, "described stale state UI")
		} else {
			notes = append(notes, "✗ didn't describe stale-state UI")
		}
		qs[3].Score = score
		qs[3].Max = max
		qs[3].Notes = notes
	}

	return qs
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func lower(s string) string { return strings.ToLower(s) }

func tokens(s string) []string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '.' || r == ',' || r == '(' || r == ')' || r == '/' || r == '-'
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.ToLower(strings.TrimSpace(f))
		if len(f) >= 3 {
			out = append(out, f)
		}
	}
	return out
}

func containsAny(haystack string, needles []string) bool {
	for _, n := range needles {
		if n != "" && strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}

func lowerList(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = strings.ToLower(s)
	}
	return out
}

func hasAnySubstring(haystack string, needles []string) bool {
	for _, n := range needles {
		if n != "" && strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}

func hasAnyPort(haystack string, ports []int) bool {
	for _, p := range ports {
		if p > 0 && strings.Contains(haystack, strconv.Itoa(p)) {
			return true
		}
	}
	return false
}

func anyCategoryKeyword(a, cat string) bool {
	switch cat {
	case "lighting":
		return containsAny(a, []string{"light", "bulb", "lamp", "dimmer", "illumin"})
	case "security":
		return containsAny(a, []string{"camera", "lock", "doorbell", "motion", "security", "alarm", "sensor"})
	case "climate":
		return containsAny(a, []string{"thermostat", "temperature", "hvac", "heat", "cool", "air", "humidity"})
	case "energy":
		return containsAny(a, []string{"solar", "battery", "power", "energy", "panel", "inverter", "grid"})
	case "media":
		return containsAny(a, []string{"speaker", "tv", "audio", "music", "video", "stream", "cast"})
	case "smart_home":
		return containsAny(a, []string{"switch", "plug", "hub", "outlet", "home", "smart"})
	}
	return false
}

func suggestedWidgets(cap string) []string {
	switch cap {
	case "on_off":
		return []string{"toggle", "switch", "button"}
	case "brightness":
		return []string{"slider", "dimmer"}
	case "color":
		return []string{"color picker", "color wheel", "swatch"}
	case "color_temp":
		return []string{"slider", "temperature", "warmth"}
	case "fan_speed":
		return []string{"button group", "slider", "speed"}
	case "thermostat":
		return []string{"dial", "slider", "setpoint", "up/down"}
	case "temperature":
		return []string{"reading", "display", "number"}
	case "humidity":
		return []string{"reading", "display", "number"}
	case "lock_unlock":
		return []string{"toggle", "button"}
	case "garage_open_close":
		return []string{"button", "toggle"}
	case "camera_stream":
		return []string{"video", "player"}
	case "camera_snapshot":
		return []string{"image", "snapshot"}
	case "scenes":
		return []string{"button group", "chips", "grid"}
	case "volume":
		return []string{"slider"}
	case "input_select":
		return []string{"dropdown", "button group", "selector"}
	case "media_playback":
		return []string{"play", "pause", "transport"}
	case "motion":
		return []string{"indicator", "status"}
	case "doorbell":
		return []string{"button", "event"}
	case "battery_level":
		return []string{"indicator", "percentage", "bar"}
	}
	return []string{"widget"}
}
