# Tock Implementation Plan

## Goal Description
`tock` is a CLI tool to read a user's timetable from a TOML configuration and output the current or next task. It supports flexible schedules (non-7-day cycles), JSON/Natural Language output, and continuous mode for integration with status bars (Waybar, Quickshell).

## User Review Required
> [!IMPORTANT]
> **TOML Format Design**: Please review the proposed TOML structure below. It uses an `anchor_date` and `cycle_days` to support non-standard cycles (e.g., 6-day school cycles).

### Proposed TOML Format
Based on user feedback, the configuration will group tasks by day, which is more intuitive for school timetables.

```toml
# tock.toml

# Cycle configuration remains the same
cycle_days = 7
anchor_date = "2024-01-01" 

# Define the schedule by day
[[day]]
# The 'id' corresponds to the day index in the cycle (0-indexed).
# For a standard 7-day week starting Sunday: 1=Monday, 2=Tuesday...
id = 1 
tasks = [
  { name = "Wake Up", start = "07:00", end = "07:30" },
  { name = "Math", start = "08:00", end = "09:00" },
  { name = "History", start = "09:00", end = "10:00" }
]

[[day]]
id = 2
tasks = [
  { name = "Wake Up", start = "07:00", end = "07:30" },
  { name = "Science", start = "08:00", end = "09:00" },
  { name = "Gym", start = "18:00", end = "19:30" }
]

# Example of a non-standard 3-day cycle
# cycle_days = 3
# anchor_date = "2024-01-01"
# [[day]]
# id = 0 # Day 1 of the 3-day cycle
# tasks = [...]
```

### CSV Format (Alternative for 7-day cycles)
For simple standard weeks, a CSV format will be supported to reduce verbosity.
Format: `Start,End,Mon,Tue,Wed,Thu,Fri,Sat,Sun` (Header required)
Rows represent time slots. Empty cells mean no task for that day/time.

```csv
Start,End,Mon,Tue,Wed,Thu,Fri,Sat,Sun
09:00,10:00,Math,History,Math,History,Math,,
10:00,11:00,History,Math,History,Math,History,,
```

## Proposed Changes

### Library Selection
- **CLI Framework**: `github.com/spf13/cobra` - Standard for Go CLIs.
- **TOML Parsing**: `github.com/pelletier/go-toml/v2` - High performance.
- **CSV Parsing**: Standard `encoding/csv` package.
- **Time Handling**: Standard `time` package.

### Project Structure
```
tock/
├── cmd/
│   └── tock/
│       └── main.go      # Entry point
├── internal/
│   ├── config/          # Configuration loading (TOML & CSV)
│   ├── scheduler/       # Logic to determine current/next task
│   └── output/          # Formatters
├── go.mod
└── README.md
```

### Features
1.  **One-shot Mode**: `tock` (outputs current task)
2.  **Next Task**: `tock --next`
3.  **Output Flags**:
    - `--json`: Machine readable
    - `--format "..."`: Custom template? (Maybe v2, stick to simple NL/JSON for now)
    - `--show-time`: Include time ranges in output.
4.  **Continuous Mode**: `tock --watch` (blocks and outputs JSON stream or updates on change).

## Verification Plan

### Automated Tests
- Unit tests for `scheduler` package:
    - Test cycle calculation with `anchor_date`.
    - Test "current task" logic at boundaries (start/end times).
    - Test "next task" logic (including wrapping to next day/cycle).

### Manual Verification
- Create a `tock.toml` with a known schedule.
- Run `tock` at specific times (mocking time might be needed for manual test or just change system time/config).
- Run `tock --watch` and verify it updates when a task transition occurs.
