# Sked

Sked is a CLI tool that helps you keep track of your schedule. It reads your timetable from a TOML or CSV configuration and tells you what you should be doing right now or what's coming up next.

## Features

- **Flexible Configuration**: Supports both TOML (for complex cycles) and CSV (for simple weekly schedules).
- **Cycle Support**: Handle non-standard schedules (e.g., 6-day school cycles) using TOML.
- **Output Formats**: Natural language or JSON output for integration with scripts.
- **Continuous Mode**: Watch mode for status bars (Waybar, Polybar, etc.).
- **Desktop Notifications**: Native notifications on Linux, macOS, and Windows.

## Installation

```bash
# Install with Go (requires Go 1.25+)
go install github.com/Daniel-42-z/sked/cmd/sked@latest
```

```bash
# Build from source
mkdir -p build
go build -o build/sked ./cmd/sked

# Install to system (example)
sudo cp build/sked /usr/local/bin/
```

## Usage

```bash
sked                  # Show current task
sked --next           # Show next task
sked --time           # Include time range
sked --json           # Output as JSON
sked --watch          # Run in continuous mode
sked --watch --notify-ahead 5m # Notify 5m before tasks (requires --watch), will still show task info in output based on the other flags
sked --config my.toml # Use specific config file
```

## Configuration

### TOML (Recommended for complex cycles)

```toml
cycle_days = 7
anchor_date = "2024-01-01" # Optional anchor for cycle calculation

[[day]]
id = 1 # Monday (if 7-day cycle)
tasks = [
  { name = "Math", start = "09:00", end = "10:00" }
]
```

### Overrides

You can temporarily override a specific date's schedule.

```toml
# Treat a specific date as another day (e.g., make Jan 3rd a Monday)
[[override]]
date = "2024-01-03"
use_day_id = 1

# Mark a specific date as an OFF day (no tasks)
[[override]]
date = "2024-01-02"
is_off = true
```

### CSV (Simple weekly schedule)

```csv
Start,End,Mon,Tue,Wed,Thu,Fri,Sat,Sun
09:00,10:00,Math,History,Math,History,Math,,
```

Note: Tasks named `/` are ignored and treated as empty time slots.
