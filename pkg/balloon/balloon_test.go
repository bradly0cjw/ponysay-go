package balloon

import (
	"testing"
)

func TestParseBalloon(t *testing.T) {
	b := ParseBalloon("", false)
	if b.Link == "" {
		t.Errorf("Default balloon link should not be empty")
	}

	bThink := ParseBalloon("", true)
	if bThink.Link != "o" {
		t.Errorf("Default think balloon link should be 'o', got %s", bThink.Link)
	}
}

func TestWrapText(t *testing.T) {
	msg := "This is a long message that needs to be wrapped properly."
	wrapped := WrapText(msg, 20)
	if len(wrapped) <= 1 {
		t.Errorf("Message should be wrapped into multiple lines")
	}
}

func TestFormatBalloonTemplate(t *testing.T) {
	b := ParseBalloon("", false)
	msgLines := []string{"I am just the cutest pony!"}
	result := b.FormatBalloon(msgLines, 0, 0, "")

	if len(result) < 3 {
		t.Fatalf("Formatted balloon should have at least 3 lines, got %d", len(result))
	}
}
