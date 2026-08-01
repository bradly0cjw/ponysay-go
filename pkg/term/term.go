package term

import (
	"os"

	"golang.org/x/term"
)

// GetTerminalWidth returns the width of the terminal stdout, defaulting to 80 if unattached.
func GetTerminalWidth() int {
	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		width, _, err := term.GetSize(fd)
		if err == nil && width > 0 {
			return width
		}
	}

	// Try stdin if stdout was piped
	stdinFd := int(os.Stdin.Fd())
	if term.IsTerminal(stdinFd) {
		width, _, err := term.GetSize(stdinFd)
		if err == nil && width > 0 {
			return width
		}
	}

	return 80
}
