package scheduler

import (
	"sked/internal/config"
	"testing"
	"time"
)

func TestOverrides(t *testing.T) {
	// Setup:
	// Cycle: 7 days
	// Mon (1): Task A
	// Tue (2): Task B
	// Wed (3): Task C

	// Override 1: Tue is OFF.
	// Override 2: Wed uses Mon schedule (Task A).

	monTasks := []config.Task{{Name: "Task A", Start: "09:00", End: "10:00"}}
	tueTasks := []config.Task{{Name: "Task B", Start: "09:00", End: "10:00"}}
	wedTasks := []config.Task{{Name: "Task C", Start: "09:00", End: "10:00"}}

	// Note: We manually populate the internal fields (Date, UseDayID) 
	// because we are bypassing config.Load() logic here.
	cfg := &config.Config{
		CycleDays: 7,
		Days: []config.Day{
			{ID: 1, Tasks: monTasks},
			{ID: 2, Tasks: tueTasks},
			{ID: 3, Tasks: wedTasks},
		},
		Overrides: []config.Override{
			{
				// Tuesday Jan 2, 2024 -> OFF
				DateStr: "2024-01-02",
				IsOff:   true,
				Date:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			{
				// Wednesday Jan 3, 2024 -> Use Mon (ID 1)
				DateStr:  "2024-01-03",
				UseDayID: 1,
				Date:     time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	sched := New(cfg)

	// 1. Test Normal Monday
	// Jan 1, 2024 is a Monday
	monDate := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
	task, err := sched.GetCurrentTask(monDate)
	if err != nil {
		t.Fatalf("Monday error: %v", err)
	}
	if task == nil || task.Name != "Task A" {
		t.Errorf("Expected Task A on Monday, got %v", task)
	}

	// 2. Test OFF Tuesday
	tueDate := time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC)
	task, err = sched.GetCurrentTask(tueDate)
	if err != nil {
		t.Fatalf("Tuesday error: %v", err)
	}
	if task != nil {
		t.Errorf("Expected no task on OFF Tuesday, got %v", task)
	}

	// 3. Test Override Wednesday (should act like Monday)
	wedDate := time.Date(2024, 1, 3, 9, 30, 0, 0, time.UTC)
	task, err = sched.GetCurrentTask(wedDate)
	if err != nil {
		t.Fatalf("Wednesday error: %v", err)
	}
	if task == nil || task.Name != "Task A" {
		t.Errorf("Expected Task A on Override Wednesday, got %v", task)
	}

	// 4. Test GetNextTask skipping OFF Tuesday
	// Start search from Monday 11:00 (after Task A).
	// Should skip Tuesday (OFF) and find Wednesday (Task A again due to override).
	searchDate := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	nextTask, err := sched.GetNextTask(searchDate)
	if err != nil {
		t.Fatalf("GetNextTask error: %v", err)
	}
	if nextTask == nil {
		t.Fatal("Expected next task, got nil")
	}

	// Validating it found the Wednesday instance
	if nextTask.StartTime.Day() != 3 {
		t.Errorf("Expected next task to be on Wednesday (Day 3), got Day %d", nextTask.StartTime.Day())
	}
	if nextTask.Name != "Task A" {
		t.Errorf("Expected next task to be Task A, got %s", nextTask.Name)
	}
}
