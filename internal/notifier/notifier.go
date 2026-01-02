// Package notifier provides cross-platform desktop notification support.
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
	case "darwin":
		return sendDarwin(title, message)
	case "windows":
		return sendWindows(title, message)
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

func sendDarwin(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, escapeQuotes(message), escapeQuotes(title))
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
}

func sendWindows(title, message string) error {
	// PowerShell script to show a balloon tip
	// We use a small delay to ensure the balloon has time to appear before the icon is disposed
	script := fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms
$notify = New-Object System.Windows.Forms.NotifyIcon
$notify.Icon = [System.Drawing.SystemIcons]::Information
$notify.Visible = $true
$notify.ShowBalloonTip(10000, "%s", "%s", [System.Windows.Forms.ToolTipIcon]::None)
Start-Sleep -s 5
$notify.Visible = $false
$notify.Dispose()
`, escapeQuotes(title), escapeQuotes(message))

	cmd := exec.Command("powershell", "-Command", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
}

func escapeQuotes(s string) string {
	// Simple escaping for double quotes
	// In a real app, we might want more robust escaping depending on the shell
	// For now, replacing " with ' or \" is a basic safeguard
	var result []rune
	for _, r := range s {
		if r == '"' {
			result = append(result, '\\', '"')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
