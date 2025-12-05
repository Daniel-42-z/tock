package scheduler

import (
	"testing"
	"time"
	"tock/internal/config"
)

func TestGetCurrentTask(t *testing.T) {
	cfg := &config.Config{
		CycleDays: 7,
		Days: []config.Day{
			{
				ID: 1, // Monday
				Tasks: []config.Task{
					{Name: "Task A", Start: "09:00", End: "10:00"},
				},
			},
		},
	}
	sched := New(cfg)

	// Test case: Monday 09:30 (Should match)
	// 2024-01-01 was a Monday
	now := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
	task, err := sched.GetCurrentTask(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task == nil {
		t.Fatal("expected task, got nil")
	}
	if task.Name != "Task A" {
		t.Errorf("expected Task A, got %s", task.Name)
	}

	// Test case: Monday 10:30 (Should not match)
	now = time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC)
	task, err = sched.GetCurrentTask(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task != nil {
		t.Errorf("expected nil, got %v", task)
	}
}

func TestGetNextTask(t *testing.T) {
	cfg := &config.Config{
		CycleDays: 7,
		Days: []config.Day{
			{
				ID: 1, // Monday
				Tasks: []config.Task{
					{Name: "Task A", Start: "09:00", End: "10:00"},
					{Name: "Task B", Start: "11:00", End: "12:00"},
				},
			},
			{
				ID: 2, // Tuesday
				Tasks: []config.Task{
					{Name: "Task C", Start: "09:00", End: "10:00"},
				},
			},
		},
	}
	sched := New(cfg)

	// Case 1: Before Task A on Monday
	now := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	task, err := sched.GetNextTask(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task == nil || task.Name != "Task A" {
		t.Errorf("expected Task A, got %v", task)
	}

	// Case 2: Between Task A and Task B on Monday
	now = time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC)
	task, err = sched.GetNextTask(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task == nil || task.Name != "Task B" {
		t.Errorf("expected Task B, got %v", task)
	}

	// Case 3: After Task B on Monday (Should find Task C on Tuesday)
	now = time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)
	task, err = sched.GetNextTask(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task == nil || task.Name != "Task C" {
		t.Errorf("expected Task C, got %v", task)
	}
}

func TestCycleLogic(t *testing.T) {
	// 3-day cycle
	// Anchor: 2024-01-01 (Day 0)
	// 2024-01-02 (Day 1)
	// 2024-01-03 (Day 2)
	// 2024-01-04 (Day 0)
	cfg := &config.Config{
		CycleDays:  3,
		AnchorDate: "2024-01-01",
		Days: []config.Day{
			{
				ID: 0,
				Tasks: []config.Task{
					{Name: "Day 0 Task", Start: "10:00", End: "11:00"},
				},
			},
		},
	}
	sched := New(cfg)

	// Check 2024-01-04 (Should be Day 0)
	now := time.Date(2024, 1, 4, 10, 30, 0, 0, time.UTC)
	task, err := sched.GetCurrentTask(now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task == nil || task.Name != "Day 0 Task" {
		t.Errorf("expected Day 0 Task, got %v", task)
	}
}
