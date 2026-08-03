package assets

import (
	"strings"
	"testing"
)

func TestFormatColumnisedList(t *testing.T) {
	items := []string{"ascii", "cowsay", "linux-vt", "round", "unicode"}

	// Terminal width 80 -> all items fit on 1 row
	out80 := FormatColumnisedList(items, 80)
	expected80 := "ascii     cowsay    linux-vt  round     unicode\n\n"
	if out80 != expected80 {
		t.Errorf("FormatColumnisedList 80 width mismatch.\nGot: %q\nExpected: %q", out80, expected80)
	}

	// Terminal width 20 -> 2 columns, 3 rows
	out20 := FormatColumnisedList(items, 20)
	lines20 := strings.Split(strings.TrimRight(out20, "\n"), "\n")
	if len(lines20) != 3 {
		t.Fatalf("Expected 3 lines for width 20, got %d. Output:\n%s", len(lines20), out20)
	}
	if !strings.HasPrefix(lines20[0], "ascii") || !strings.Contains(lines20[0], "round") {
		t.Errorf("Line 0 expected ascii and round, got %q", lines20[0])
	}
	if !strings.HasPrefix(lines20[1], "cowsay") || !strings.Contains(lines20[1], "unicode") {
		t.Errorf("Line 1 expected cowsay and unicode, got %q", lines20[1])
	}
	if lines20[2] != "linux-vt" {
		t.Errorf("Line 2 expected linux-vt without trailing spacing, got %q", lines20[2])
	}

	// Empty list
	if outEmpty := FormatColumnisedList([]string{}, 80); outEmpty != "" {
		t.Errorf("Expected empty string for empty input, got %q", outEmpty)
	}
}
