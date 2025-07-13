package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenBrowser opens the default browser with the given URL
func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		// Try xdg-open first, then fall back to common browsers
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", url)
		} else {
			// Try common browsers
			browsers := []string{"firefox", "google-chrome", "chromium", "opera"}
			for _, browser := range browsers {
				if _, err := exec.LookPath(browser); err == nil {
					cmd = exec.Command(browser, url)
					break
				}
			}
			if cmd == nil {
				return fmt.Errorf("no suitable browser found")
			}
		}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	
	return cmd.Start()
}