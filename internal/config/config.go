package config

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the top-level configuration structure.
type Config struct {
	CycleDays  int    `toml:"cycle_days"`
	AnchorDate string `toml:"anchor_date"`
	CSVPath    string `toml:"csv_path"`
	Days       []Day  `toml:"day"`
}

// Day represents a single day's schedule in the cycle.
type Day struct {
	ID    int    `toml:"id"`
	Tasks []Task `toml:"tasks"`
}

// Task represents a specific activity.
type Task struct {
	Name  string `toml:"name"`
	Start string `toml:"start"`
	End   string `toml:"end"`
}

// Load reads the configuration from the specified path.
// It detects the format based on the file extension (.toml or .csv).
func Load(path string) (*Config, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".toml":
		return LoadTOML(path)
	case ".csv":
		return LoadCSV(path)
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}

// LoadTOML reads a TOML configuration file.
func LoadTOML(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	// Set defaults
	cfg.CycleDays = 7

	dec := toml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}

	// Check for CSV redirection
	if cfg.CSVPath != "" {
		csvPath := cfg.CSVPath
		// If path is relative, resolve it relative to the TOML file
		if !filepath.IsAbs(csvPath) {
			csvPath = filepath.Join(filepath.Dir(path), csvPath)
		}
		return LoadCSV(csvPath)
	}

	return &cfg, nil
}

// LoadCSV reads a CSV configuration file.
// CSV format assumes a standard 7-day cycle.
// Header: Start,End,Mon,Tue,Wed,Thu,Fri,Sat,Sun (flexible day column order)
func LoadCSV(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("csv file is empty or missing header")
	}

	header := records[0]
	if len(header) < 3 {
		return nil, fmt.Errorf("header must have at least Start, End and one Day column")
	}

	// Map column index to day ID
	colToDay := make(map[int]int)
	startCol := -1
	endCol := -1

	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))
		if col == "start" || col == "time-start" {
			startCol = i
		} else if col == "end" || col == "time-end" {
			endCol = i
		} else {
			// Try to parse as day
			dayID, err := parseDayName(col)
			if err == nil {
				colToDay[i] = dayID
			}
		}
	}

	if startCol == -1 || endCol == -1 {
		return nil, fmt.Errorf("header must contain 'Start' and 'End' columns")
	}

	cfg := &Config{
		CycleDays: 7,
		Days:      make([]Day, 0),
	}
	dayMap := make(map[int][]Task)

	for _, record := range records[1:] {
		if len(record) <= startCol || len(record) <= endCol {
			continue // Skip invalid rows
		}

		start := strings.TrimSpace(record[startCol])
		end := strings.TrimSpace(record[endCol])

		if start == "" {
			continue // Skip rows without start time
		}

		for colIdx, dayID := range colToDay {
			if colIdx >= len(record) {
				continue
			}
			name := strings.TrimSpace(record[colIdx])
			if name != "" {
				task := Task{
					Name:  name,
					Start: start,
					End:   end,
				}
				dayMap[dayID] = append(dayMap[dayID], task)
			}
		}
	}

	// Convert map to slice
	for id, tasks := range dayMap {
		cfg.Days = append(cfg.Days, Day{
			ID:    id,
			Tasks: tasks,
		})
	}

	return cfg, nil
}

// parseDayName converts a day name (e.g., "Monday") to a cycle ID (0-6).
// Assumes 0=Sunday, 1=Monday, ..., 6=Saturday to match time.Weekday().
func parseDayName(name string) (int, error) {
	name = strings.ToLower(name)
	switch name {
	case "sunday", "sun":
		return 0, nil
	case "monday", "mon":
		return 1, nil
	case "tuesday", "tue":
		return 2, nil
	case "wednesday", "wed":
		return 3, nil
	case "thursday", "thu":
		return 4, nil
	case "friday", "fri":
		return 5, nil
	case "saturday", "sat":
		return 6, nil
	default:
		return -1, fmt.Errorf("invalid day name: %s", name)
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.CycleDays <= 0 {
		return fmt.Errorf("cycle_days must be positive")
	}
	if c.CycleDays != 7 && c.AnchorDate == "" {
		return fmt.Errorf("anchor_date is required for non-7-day cycles")
	}
	if c.AnchorDate != "" {
		_, err := time.Parse("2006-01-02", c.AnchorDate)
		if err != nil {
			return fmt.Errorf("invalid anchor_date format (expected YYYY-MM-DD): %w", err)
		}
	}
	// TODO: Validate time formats (HH:MM)
	return nil
}
