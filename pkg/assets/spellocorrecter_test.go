package assets

import (
	"testing"
)

func TestSpelloCorrecter(t *testing.T) {
	sc := NewSpelloCorrecter()

	candidates := []string{
		"fluttershy",
		"twilight",
		"applejack",
		"rarity",
		"rainbow-dash",
		"pinkie-pie",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"fluter-shy", "fluttershy"},
		{"twiligt", "twilight"},
		{"aplejak", "applejack"},
		{"rariti", "rarity"},
	}

	for _, tt := range tests {
		best, dist := sc.Correct(tt.input, candidates)
		if len(best) == 0 {
			t.Errorf("Expected match for '%s', got none", tt.input)
			continue
		}
		if best[0] != tt.expected {
			t.Errorf("For '%s', expected '%s', got '%s' (dist: %f)", tt.input, tt.expected, best[0], dist)
		}
	}

	// Test case outside of distance limit (5)
	_, dist := sc.Correct("completelydifferentunmatchedname", candidates)
	if dist <= 5.0 {
		t.Errorf("Expected distance > 5 for unmatched name, got %f", dist)
	}
}
