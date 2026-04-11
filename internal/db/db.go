package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens a SQLite database at the given path with WAL mode and foreign keys enabled.
// Mother always said a well-organized database is the foundation of a good home.
func Open(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=on", path)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("[db] Database opened and migrated successfully. Mother would be proud.")
	return db, nil
}

// Migrate creates the initial tables if they don't exist.
func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS devices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT NOT NULL,
			mac TEXT NOT NULL DEFAULT '',
			hostname TEXT NOT NULL DEFAULT '',
			name TEXT NOT NULL DEFAULT '',
			manufacturer TEXT NOT NULL DEFAULT '',
			model TEXT NOT NULL DEFAULT '',
			device_type TEXT NOT NULL DEFAULT '',
			category TEXT NOT NULL DEFAULT 'unknown',
			protocols TEXT NOT NULL DEFAULT '[]',
			services TEXT NOT NULL DEFAULT '[]',
			open_ports TEXT NOT NULL DEFAULT '[]',
			metadata TEXT NOT NULL DEFAULT '{}',
			first_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(ip)
		);

		CREATE TABLE IF NOT EXISTS scans (
			id TEXT PRIMARY KEY,
			status TEXT NOT NULL DEFAULT 'running',
			device_count INTEGER NOT NULL DEFAULT 0,
			duration_ms INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS hue_config (
			id INTEGER PRIMARY KEY,
			bridge_ip TEXT NOT NULL,
			username TEXT NOT NULL,
			bridge_id TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS device_credentials (
			ip TEXT PRIMARY KEY,
			integration TEXT NOT NULL DEFAULT '',
			username TEXT NOT NULL DEFAULT '',
			password TEXT NOT NULL DEFAULT '',
			session TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS kasa_devices (
			ip TEXT PRIMARY KEY,
			alias TEXT NOT NULL DEFAULT '',
			model TEXT NOT NULL DEFAULT '',
			device_type TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS google_tokens (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			token_type TEXT NOT NULL DEFAULT 'Bearer',
			expiry DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

// HueConfig holds a stored Hue bridge configuration.
type HueConfig struct {
	BridgeIP string
	Username string
	BridgeID string
}

// SaveHueConfig inserts or replaces the Hue bridge configuration.
// There's only ever one row (id=1). Mother says one bridge is enough.
func SaveHueConfig(db *sql.DB, bridgeIP, username, bridgeID string) error {
	_, err := db.Exec(`
		INSERT OR REPLACE INTO hue_config (id, bridge_ip, username, bridge_id)
		VALUES (1, ?, ?, ?)
	`, bridgeIP, username, bridgeID)
	return err
}

// LoadHueConfig loads the stored Hue bridge configuration, if any.
func LoadHueConfig(db *sql.DB) (*HueConfig, error) {
	var cfg HueConfig
	err := db.QueryRow(`SELECT bridge_ip, username, bridge_id FROM hue_config WHERE id = 1`).
		Scan(&cfg.BridgeIP, &cfg.Username, &cfg.BridgeID)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveDeviceCredential stores auth credentials for a device.
func SaveDeviceCredential(db *sql.DB, ip, integration, username, password, session string) error {
	_, err := db.Exec(`
		INSERT OR REPLACE INTO device_credentials (ip, integration, username, password, session)
		VALUES (?, ?, ?, ?, ?)
	`, ip, integration, username, password, session)
	return err
}

// DeviceCredential holds stored auth info for a device.
type DeviceCredential struct {
	IP          string
	Integration string
	Username    string
	Password    string
	Session     string
}

// LoadDeviceCredential loads stored credentials for a device.
func LoadDeviceCredential(db *sql.DB, ip string) (*DeviceCredential, error) {
	var cred DeviceCredential
	err := db.QueryRow(`SELECT ip, integration, username, password, session FROM device_credentials WHERE ip = ?`, ip).
		Scan(&cred.IP, &cred.Integration, &cred.Username, &cred.Password, &cred.Session)
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

// DeleteHueConfig removes the stored Hue bridge configuration.
func DeleteHueConfig(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM hue_config WHERE id = 1`)
	return err
}

// GoogleTokens holds OAuth2 tokens for the Google Nest SDM API.
// Only one set of tokens per hub -- you don't need two identities
// unless you're running from the SEC.
type GoogleTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// SaveGoogleTokens stores (or replaces) the Google OAuth tokens.
// Single row, id=1. One set of credentials. No aliases.
func SaveGoogleTokens(db *sql.DB, accessToken, refreshToken string, expiresAt time.Time) error {
	_, err := db.Exec(`
		INSERT OR REPLACE INTO google_tokens (id, access_token, refresh_token, token_type, expiry)
		VALUES (1, ?, ?, 'Bearer', ?)
	`, accessToken, refreshToken, expiresAt.UTC().Format(time.RFC3339))
	return err
}

// LoadGoogleTokens loads the stored Google OAuth tokens, if any.
func LoadGoogleTokens(db *sql.DB) (*GoogleTokens, error) {
	var t GoogleTokens
	var expiryStr string
	err := db.QueryRow(`SELECT access_token, refresh_token, expiry FROM google_tokens WHERE id = 1`).
		Scan(&t.AccessToken, &t.RefreshToken, &expiryStr)
	if err != nil {
		return nil, err
	}
	t.ExpiresAt, _ = time.Parse(time.RFC3339, expiryStr)
	return &t, nil
}

// DeleteGoogleTokens removes stored Google OAuth tokens.
// When the feds come knocking, you shred the evidence.
func DeleteGoogleTokens(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM google_tokens WHERE id = 1`)
	return err
}

// UpsertDevice inserts a device or updates it if the IP already exists.
// Updates last_seen on every upsert; preserves first_seen.
func UpsertDevice(db *sql.DB, ip, mac, hostname, name, manufacturer, model, deviceType, category, protocols, services, openPorts, metadata string) error {
	_, err := db.Exec(`
		INSERT INTO devices (ip, mac, hostname, name, manufacturer, model, device_type, category, protocols, services, open_ports, metadata, first_seen, last_seen)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(ip) DO UPDATE SET
			mac = CASE WHEN excluded.mac != '' THEN excluded.mac ELSE devices.mac END,
			hostname = CASE WHEN excluded.hostname != '' THEN excluded.hostname ELSE devices.hostname END,
			name = CASE
				WHEN devices.name != '' AND devices.name NOT LIKE 'Device .%' AND devices.name NOT LIKE '% .%' THEN devices.name
				WHEN excluded.name != '' THEN excluded.name
				ELSE devices.name END,
			manufacturer = CASE WHEN excluded.manufacturer != '' THEN excluded.manufacturer ELSE devices.manufacturer END,
			model = CASE WHEN excluded.model != '' THEN excluded.model ELSE devices.model END,
			device_type = CASE
				WHEN devices.device_type != '' AND devices.device_type NOT IN ('', 'nest_device') THEN devices.device_type
				WHEN excluded.device_type != '' THEN excluded.device_type
				ELSE devices.device_type END,
			category = CASE WHEN excluded.category != '' AND excluded.category != 'unknown' THEN excluded.category ELSE devices.category END,
			protocols = CASE WHEN excluded.protocols != '[]' THEN excluded.protocols ELSE devices.protocols END,
			services = CASE WHEN excluded.services != '[]' THEN excluded.services ELSE devices.services END,
			open_ports = CASE WHEN excluded.open_ports != '[]' THEN excluded.open_ports ELSE devices.open_ports END,
			metadata = CASE WHEN excluded.metadata != '{}' THEN excluded.metadata ELSE devices.metadata END,
			last_seen = CURRENT_TIMESTAMP
	`, ip, mac, hostname, name, manufacturer, model, deviceType, category, protocols, services, openPorts, metadata)
	return err
}

// DeviceRow represents a device loaded from the database.
type DeviceRow struct {
	IP           string
	MAC          string
	Hostname     string
	Name         string
	Manufacturer string
	Model        string
	DeviceType   string
	Category     string
	Protocols    string // JSON array
	Services     string // JSON array
	OpenPorts    string // JSON array
	Metadata     string // JSON object
}

// LoadAllDevices returns all persisted devices ordered by last_seen desc.
func LoadAllDevices(db *sql.DB) ([]DeviceRow, error) {
	rows, err := db.Query(`
		SELECT ip, mac, hostname, name, manufacturer, model, device_type, category, protocols, services, open_ports, metadata
		FROM devices ORDER BY last_seen DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []DeviceRow
	for rows.Next() {
		var d DeviceRow
		if err := rows.Scan(&d.IP, &d.MAC, &d.Hostname, &d.Name, &d.Manufacturer, &d.Model, &d.DeviceType, &d.Category, &d.Protocols, &d.Services, &d.OpenPorts, &d.Metadata); err != nil {
			continue
		}
		devices = append(devices, d)
	}
	return devices, nil
}
