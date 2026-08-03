package color

import (
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// StripANSI removes all ANSI escape codes from a string.
func StripANSI(s string) string {
	if strings.IndexByte(s, '\x1b') == -1 {
		return s
	}
	return ansiRegex.ReplaceAllString(s, "")
}

// DisplayWidth returns the visible terminal cell width of a string.
func DisplayWidth(s string) int {
	if strings.IndexByte(s, '\x1b') == -1 {
		isASCII := true
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				isASCII = false
				break
			}
		}
		if isASCII {
			return len(s)
		}
	}

	clean := StripANSI(s)
	width := 0
	for _, r := range clean {
		w := runewidth.RuneWidth(r)
		if w > 0 {
			width += w
		}
	}
	return width
}

// ApplyColor wraps text in ANSI color sequence if specified.
func ApplyColor(text string, colorCode string) string {
	if colorCode == "" {
		return text
	}
	// If it's a simple number or standard format, ensure \x1b[ ... m
	if !strings.HasPrefix(colorCode, "\x1b[") {
		if strings.HasSuffix(colorCode, "m") {
			colorCode = "\x1b[" + colorCode
		} else {
			colorCode = "\x1b[" + colorCode + "m"
		}
	}
	return colorCode + text + "\x1b[0m"
}
