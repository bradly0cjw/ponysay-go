package assets

import (
	"strings"
	"unicode/utf8"

	"ponysay-go/pkg/term"
)

// stripANSI removes ANSI escape sequences from string.
func stripANSI(s string) string {
	var builder strings.Builder
	inSeq := false
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == 0x1b {
			inSeq = true
			continue
		}
		if inSeq {
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') {
				inSeq = false
			}
			continue
		}
		builder.WriteByte(b)
	}
	return builder.String()
}

// GetPonyWidth computes maximum line width of pony artwork (ignoring ANSI color sequences).
func GetPonyWidth(content string) int {
	maxWidth := 0
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		clean := stripANSI(l)
		w := utf8.RuneCountInString(clean)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

// FilterFittingPonies filters ponies whose width fits within terminal width constraint.
func (am *AssetManager) FilterFittingPonies(ponies []string, includeStandard bool, includeExtra bool) []string {
	termWidth := term.GetTerminalWidth()
	if termWidth <= 0 {
		return ponies
	}

	var fitting []string
	for _, p := range ponies {
		w, ok := am.ponyWidths[p]
		if !ok {
			_, content, err := am.getPonyFileLocked(p, includeStandard, includeExtra)
			if err == nil {
				w = GetPonyWidth(content)
				am.ponyWidths[p] = w
			}
		}
		if w > 0 && w <= termWidth {
			fitting = append(fitting, p)
		}
	}

	if len(fitting) > 0 {
		return fitting
	}
	return ponies
}
