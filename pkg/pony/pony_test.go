package pony

import (
	"testing"
)

func TestParsePony(t *testing.T) {
	raw := `$$$
NAME: Derpy
COAT: grey
BALLOON TOP: 4
$$$
$balloon9$
     $\$
      ▄▄▄`

	p, err := ParsePony("derpy", raw)
	if err != nil {
		t.Fatalf("Failed to parse pony: %v", err)
	}

	if p.Name != "derpy" {
		t.Errorf("Expected pony name derpy, got %s", p.Name)
	}

	if p.Metadata["NAME"] != "Derpy" {
		t.Errorf("Expected metadata NAME Derpy, got %s", p.Metadata["NAME"])
	}

	if p.BalloonTop != 4 {
		t.Errorf("Expected BalloonTop 4, got %d", p.BalloonTop)
	}

	if len(p.BodyLines) < 2 {
		t.Errorf("Expected at least 2 body lines")
	}
}

func TestRenderPonyWithBalloon(t *testing.T) {
	raw := `$$$
NAME: Test
$$$
$balloon5$
     $\$`

	p, _ := ParsePony("test", raw)
	balloonLines := []string{"┌───┐", "│Hi │", "└───┘"}
	rendered := p.RenderPonyWithBalloon(balloonLines, "╲", "")

	if len(rendered) == 0 {
		t.Errorf("Rendered pony should not be empty")
	}
}
