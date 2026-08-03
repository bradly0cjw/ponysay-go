package assets

import (
	"strings"

	"ponysay-go/pkg/color"
)

// FormatColumnisedList formats a slice of item strings into a multi-column matrix string
// matching Python ponysay's _columnise_list and _print_columnised logic.
func FormatColumnisedList(items []string, availableWidth int) string {
	if len(items) == 0 {
		return ""
	}

	if availableWidth <= 0 {
		availableWidth = 80
	}

	separation := 2

	maxItemLength := 0
	for _, item := range items {
		w := color.DisplayWidth(item)
		if w > maxItemLength {
			maxItemLength = w
		}
	}

	numColumns := (availableWidth + separation) / (maxItemLength + separation)
	if numColumns < 1 {
		numColumns = 1
	}

	columnLength := (len(items)-1)/numColumns + 1

	type cellWithSpacing struct {
		display string
		spacing int
	}

	var columns [][]cellWithSpacing
	for i := 0; i < len(items); i += columnLength {
		end := i + columnLength
		if end > len(items) {
			end = len(items)
		}
		var col []cellWithSpacing
		for j := i; j < end; j++ {
			w := color.DisplayWidth(items[j])
			sp := maxItemLength - w + separation
			col = append(col, cellWithSpacing{display: items[j], spacing: sp})
		}
		columns = append(columns, col)
	}

	var sb strings.Builder
	for r := 0; r < columnLength; r++ {
		spacing := 0
		for c := 0; c < len(columns); c++ {
			if r < len(columns[c]) {
				cell := columns[c][r]
				if spacing > 0 {
					sb.WriteString(strings.Repeat(" ", spacing))
				}
				sb.WriteString(cell.display)
				spacing = cell.spacing
			}
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	return sb.String()
}
