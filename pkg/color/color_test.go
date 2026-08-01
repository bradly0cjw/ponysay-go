package color

import (
	"testing"
)

func TestStripANSI(t *testing.T) {
	colored := "\x1b[31mRed Text\x1b[0m"
	clean := StripANSI(colored)
	if clean != "Red Text" {
		t.Errorf("Expected 'Red Text', got '%s'", clean)
	}
}

func TestDisplayWidth(t *testing.T) {
	strASCII := "Hello World"
	if DisplayWidth(strASCII) != 11 {
		t.Errorf("Expected width 11, got %d", DisplayWidth(strASCII))
	}

	strBlock := "▄▄▄"
	if DisplayWidth(strBlock) != 3 {
		t.Errorf("Expected width 3 for block graphics, got %d", DisplayWidth(strBlock))
	}
}
