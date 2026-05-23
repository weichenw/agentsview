package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SchedulerSettings holds user-configurable scheduler options.
type SchedulerSettings struct {
	Timezone string `json:"timezone"`
}

// defaultTimezone is the fallback when no settings file exists.
const defaultTimezone = "Australia/Sydney"

func settingsPath(dataDir string) string {
	return filepath.Join(dataDir, "scheduler-settings.json")
}

// LoadSettings reads scheduler-settings.json from the data dir.
// Returns a default settings struct if the file does not exist.
func LoadSettings(dataDir string) (*SchedulerSettings, error) {
	path := settingsPath(dataDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &SchedulerSettings{Timezone: defaultTimezone}, nil
		}
		return nil, fmt.Errorf("scheduler: read settings: %w", err)
	}
	var settings SchedulerSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("scheduler: parse settings: %w", err)
	}
	if settings.Timezone == "" {
		settings.Timezone = defaultTimezone
	}
	return &settings, nil
}

// SaveSettings atomically writes scheduler-settings.json to the data dir.
func SaveSettings(dataDir string, settings *SchedulerSettings) error {
	path := settingsPath(dataDir)
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("scheduler: marshal settings: %w", err)
	}
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		return fmt.Errorf("scheduler: write settings temp: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("scheduler: rename settings: %w", err)
	}
	return nil
}
