package assets

import (
	"os"
	"testing"
)

func TestUCSRemapping(t *testing.T) {
	am := NewAssetManager()

	// Test without PONYSAY_UCS_ME
	os.Unsetenv("PONYSAY_UCS_ME")
	if remapped := am.RemapUCS("mjölna"); remapped != "mjölna" {
		t.Errorf("Expected 'mjölna' unchanged when PONYSAY_UCS_ME is unset, got '%s'", remapped)
	}

	// Test with PONYSAY_UCS_ME=1
	os.Setenv("PONYSAY_UCS_ME", "1")
	defer os.Unsetenv("PONYSAY_UCS_ME")

	tests := []struct {
		input    string
		expected string
	}{
		{"mjölna", "mjolna"},
		{"bifröst", "bifrost"},
		{"piñacolada", "pinacolada"},
		{"bárbara", "barbara"},
		{"fluttershy", "fluttershy"}, // Unmapped ASCII stays unchanged
	}

	for _, tt := range tests {
		if remapped := am.RemapUCS(tt.input); remapped != tt.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", tt.input, tt.expected, remapped)
		}
	}
}

func TestGetPonyWidthAndFilter(t *testing.T) {
	artwork := "\x1b[38;5;196mHello\x1b[0m\n\x1b[38;5;200mWorld Wide Web\x1b[0m"
	w := GetPonyWidth(artwork)
	if w != 14 { // "World Wide Web" is 14 chars
		t.Errorf("Expected pony width 14, got %d", w)
	}
}
