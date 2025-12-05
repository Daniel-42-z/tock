package notifier

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Notifier handles sending desktop notifications.
type Notifier struct{}

// New creates a new Notifier.
func New() *Notifier {
	return &Notifier{}
}

// Send sends a notification with the given title and message.
func (n *Notifier) Send(title, message string) error {
	switch runtime.GOOS {
	case "linux":
		return sendLinux(title, message)
	// Add other platforms here if needed
	default:
		return fmt.Errorf("notifications not supported on %s", runtime.GOOS)
	}
}

func sendLinux(title, message string) error {
	cmd := exec.Command("notify-send", title, message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
}
