package api

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/coalson/haus/internal/db"
)

// enrichMu ensures only one Nest enrichment pass runs at a time. Multiple
// triggers (startup goroutine + OAuth callback) could collide otherwise.
var enrichMu sync.Mutex

// EnrichNestDevices queries the Google SDM API for the user's Nest devices
// (cameras, displays, doorbells, thermostats) and upgrades the matching DB
// rows with the SDM-derived display name, correct device_type, and the
// capabilities implied by each type. This lets a Nest Hub Max that the
// network scan originally tagged as a Thread border router light up with
// camera chat tools once the user signs in to Google.
//
// Safe to call multiple times — no-op if Google isn't connected, idempotent
// when it is. Cheap enough to run on every OAuth completion.
func (s *Server) EnrichNestDevices() {
	enrichMu.Lock()
	defer enrichMu.Unlock()

	client, err := s.GetGoogleClient()
	if err != nil {
		log.Printf("[enrich] Google not connected, skipping Nest enrichment: %v", err)
		return
	}

	sdmDevices, err := client.ListDevices()
	if err != nil {
		log.Printf("[enrich] Failed to list Nest devices: %v", err)
		return
	}

	dbDevices, err := db.LoadAllDevices(s.DB)
	if err != nil {
		log.Printf("[enrich] Failed to load devices: %v", err)
		return
	}

	// Build SDM entries we can work with: name + type suffix (THERMOSTAT, CAMERA,
	// DISPLAY, DOORBELL). ParentRelations[0] carries the human-readable room name.
	type sdmEntry struct {
		DisplayName string
		TypeSuffix  string // upper-case: "THERMOSTAT" | "CAMERA" | "DISPLAY" | "DOORBELL"
	}
	var entries []sdmEntry
	for _, d := range sdmDevices {
		displayName := ""
		if len(d.ParentRelations) > 0 {
			displayName = d.ParentRelations[0].DisplayName
		}
		if displayName == "" {
			continue
		}
		parts := strings.Split(d.Type, ".")
		entries = append(entries, sdmEntry{
			DisplayName: displayName,
			TypeSuffix:  strings.ToUpper(parts[len(parts)-1]),
		})
	}
	if len(entries) == 0 {
		log.Printf("[enrich] SDM returned no named devices — nothing to enrich.")
		return
	}

	// Partition DB devices into the Google population we'll match against.
	var googleDB []db.DeviceRow
	for _, row := range dbDevices {
		if strings.EqualFold(row.Manufacturer, "Google") || strings.Contains(strings.ToLower(row.Manufacturer), "nest") {
			googleDB = append(googleDB, row)
		}
	}

	used := make(map[string]bool) // ip -> already matched this pass
	consumed := make(map[int]bool) // entry idx -> already matched this pass
	updated := 0

	// Pass 1: exact case-insensitive name match against any SDM entry.
	for idx, e := range entries {
		if consumed[idx] {
			continue
		}
		lname := strings.ToLower(e.DisplayName)
		for _, row := range googleDB {
			if used[row.IP] {
				continue
			}
			if strings.ToLower(row.Name) == lname {
				applyEnrichment(s.DB, row, e.DisplayName, e.TypeSuffix)
				used[row.IP] = true
				consumed[idx] = true
				updated++
				break
			}
		}
	}

	// Pass 2: substring match both directions (e.g. SDM "Kitchen" ↔ DB "Kitchen display").
	for idx, e := range entries {
		if consumed[idx] {
			continue
		}
		lname := strings.ToLower(e.DisplayName)
		for _, row := range googleDB {
			if used[row.IP] {
				continue
			}
			drow := strings.ToLower(row.Name)
			if drow != "" && (strings.Contains(drow, lname) || strings.Contains(lname, drow)) {
				applyEnrichment(s.DB, row, e.DisplayName, e.TypeSuffix)
				used[row.IP] = true
				consumed[idx] = true
				updated++
				break
			}
		}
	}

	// Pass 3: round-robin across still-generic-named Google devices.
	for idx, e := range entries {
		if consumed[idx] {
			continue
		}
		for _, row := range googleDB {
			if used[row.IP] {
				continue
			}
			if row.Name == "" || strings.HasPrefix(row.Name, "Google .") || strings.HasPrefix(row.Name, "Device .") {
				applyEnrichment(s.DB, row, e.DisplayName, e.TypeSuffix)
				used[row.IP] = true
				consumed[idx] = true
				updated++
				break
			}
		}
	}

	log.Printf("[enrich] Updated %d Google device(s) from Nest SDM API", updated)
}

// applyEnrichment sets name, device_type, and capabilities for one DB row
// based on an SDM entry type suffix.
func applyEnrichment(database *sql.DB, row db.DeviceRow, sdmDisplayName, sdmTypeSuffix string) {
	// Type + pretty suffix + capabilities per SDM type.
	var (
		deviceType  string
		prettyType  string
		capsList    []string
	)
	switch sdmTypeSuffix {
	case "CAMERA":
		deviceType = "nest_camera"
		prettyType = "Camera"
		capsList = []string{"camera_stream", "camera_snapshot", "motion"}
	case "DISPLAY":
		deviceType = "nest_camera" // Nest Hub/Hub Max have a camera; reuse the tool plumbing
		prettyType = "Display"
		capsList = []string{"camera_stream", "camera_snapshot", "motion", "media_playback"}
	case "DOORBELL":
		deviceType = "nest_camera"
		prettyType = "Doorbell"
		capsList = []string{"camera_stream", "camera_snapshot", "motion", "doorbell"}
	case "THERMOSTAT":
		deviceType = "nest_thermostat"
		prettyType = "Thermostat"
		capsList = []string{"thermostat", "temperature", "humidity"}
	default:
		deviceType = "nest_device"
		prettyType = sdmTypeSuffix
		capsList = []string{}
	}

	newName := sdmDisplayName + " " + prettyType

	// Preserve existing "good" name if it already looks curated (not a generic
	// "Google .X" / "Device .X" placeholder).
	finalName := newName
	if row.Name != "" && !strings.HasPrefix(row.Name, "Google .") && !strings.HasPrefix(row.Name, "Device .") {
		finalName = row.Name
	}

	if err := db.EnrichGoogleDeviceType(database, row.IP, deviceType, finalName, capsList); err != nil {
		log.Printf("[enrich] %s: update failed: %v", row.IP, err)
		return
	}
	log.Printf("[enrich] %s → %s (%s, caps=%v)", row.IP, finalName, sdmTypeSuffix, capsList)
}
