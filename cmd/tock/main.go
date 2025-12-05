package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tock/internal/config"
	"tock/internal/output"
	"tock/internal/scheduler"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	jsonFmt    bool
	showTime   bool
	nextTask   bool
	watchMode  bool
	noTaskText string
)

var rootCmd = &cobra.Command{
	Use:   "tock",
	Short: "A CLI timetable tool",
	Long:  `tock reads your timetable configuration and tells you what you should be doing.`,
	RunE:  run,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/tock/tock.toml)")
	rootCmd.Flags().BoolVarP(&jsonFmt, "json", "j", false, "output in JSON format")
	rootCmd.Flags().BoolVarP(&showTime, "time", "t", false, "show time ranges in output")
	rootCmd.Flags().BoolVarP(&nextTask, "next", "n", false, "show next task instead of current")
	rootCmd.Flags().BoolVarP(&watchMode, "watch", "w", false, "continuous mode (watch for changes)")
	rootCmd.Flags().StringVar(&noTaskText, "no-task-text", "No task currently.", "text to display when no task is found")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// 1. Resolve config file path
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		// Try TOML first, then CSV
		tomlPath := filepath.Join(home, ".config", "tock", "tock.toml")
		csvPath := filepath.Join(home, ".config", "tock", "tock.csv")

		if _, err := os.Stat(tomlPath); err == nil {
			cfgFile = tomlPath
		} else if _, err := os.Stat(csvPath); err == nil {
			cfgFile = csvPath
		} else {
			return fmt.Errorf("no config file found at %s or %s", tomlPath, csvPath)
		}
	}

	// 2. Load Config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// 3. Initialize Scheduler
	sched := scheduler.New(cfg)

	// 4. Handle Watch Mode
	if watchMode {
		return runWatch(sched)
	}

	// 5. Output
	now := time.Now()
	var currentTask, nextTaskEvent *scheduler.TaskEvent

	// If JSON, we want both
	if jsonFmt {
		currentTask, err = sched.GetCurrentTask(now)
		if err != nil {
			return err
		}
		nextTaskEvent, err = sched.GetNextTask(now)
		if err != nil {
			return err
		}
	} else {
		// Natural language mode: depends on flag
		if nextTask {
			// If user asked for next, we treat it as the "primary" task to print
			currentTask, err = sched.GetNextTask(now)
		} else {
			currentTask, err = sched.GetCurrentTask(now)
		}
		if err != nil {
			return err
		}
	}

	return output.Print(currentTask, nextTaskEvent, jsonFmt, showTime, noTaskText)
}

func runWatch(sched *scheduler.Scheduler) error {
	ticker := time.NewTicker(30 * time.Second)

	for {
		now := time.Now()
		var currentTask, nextTaskEvent *scheduler.TaskEvent
		var err error

		if jsonFmt {
			currentTask, err = sched.GetCurrentTask(now)
			if err == nil {
				nextTaskEvent, err = sched.GetNextTask(now)
			}
		} else {
			if nextTask {
				currentTask, err = sched.GetNextTask(now)
			} else {
				currentTask, err = sched.GetCurrentTask(now)
			}
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			output.Print(currentTask, nextTaskEvent, jsonFmt, showTime, noTaskText)
		}

		<-ticker.C
	}
}
