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
	// JSON mode outputs all three tasks (previous, current, next).
	// Natural language mode outputs only the 'current' task (which main sets based on flags).

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
