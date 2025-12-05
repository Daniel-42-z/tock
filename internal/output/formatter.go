package output

import (
	"encoding/json"
	"fmt"
	"os"
	"tock/internal/scheduler"
)

// Print displays the task information.
func Print(previous *scheduler.TaskEvent, current *scheduler.TaskEvent, next *scheduler.TaskEvent, asJSON bool, showTime bool, noTaskText string) error {
	if asJSON {
		return printJSON(previous, current, next)
	}
	// For natural language, we only print one based on flags (logic handled in main)
	// But wait, main calls this with either current OR next if not JSON.
	// If JSON, we want both.
	// Let's adjust the signature to take both, and main decides what to pass.
	// Actually, for natural language, we might still want to support printing just one.
	// Let's assume if asJSON is false, we print 'current' (which might be the 'next' task if the flag was set? No, main logic needs change).

	// If we are in natural mode, we print 'current' if it's not nil, or handle no task.
	// But main logic was: if --next, get next. if not, get current.
	// Now user says: -j should output BOTH.
	// So main needs to fetch BOTH if -j is set.

	return printNatural(current, showTime, noTaskText)
}

type jsonOutput struct {
	Previous *scheduler.TaskEvent `json:"previous"`
	Current  *scheduler.TaskEvent `json:"current"`
	Next     *scheduler.TaskEvent `json:"next"`
}

func printJSON(previous *scheduler.TaskEvent, current *scheduler.TaskEvent, next *scheduler.TaskEvent) error {
	out := jsonOutput{
		Previous: previous,
		Current:  current,
		Next:     next,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func printNatural(task *scheduler.TaskEvent, showTime bool, noTaskText string) error {
	if task == nil {
		if noTaskText != "" {
			fmt.Println(noTaskText)
		} else {
			fmt.Println("No task currently.")
		}
		return nil
	}

	if showTime {
		fmt.Printf("%s (%s - %s)\n", task.Name, task.StartTime.Format("15:04"), task.EndTime.Format("15:04"))
	} else {
		fmt.Println(task.Name)
	}
	return nil
}
