package pony

import (
	"strings"
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

	if !p.HasInfo {
		t.Errorf("Expected HasInfo to be true")
	}

	if len(p.BodyLines) < 2 {
		t.Errorf("Expected at least 2 body lines")
	}
}

func TestFormatInfo(t *testing.T) {
	rawInfo := "NAME: Derpy\nCOAT: grey\n\nThis is a comment"
	formatted := FormatInfo(rawInfo)
	if !strings.Contains(formatted, "\x1b[1mNAME\x1b[22m: Derpy") {
		t.Errorf("FormatInfo failed to format tag NAME, got %q", formatted)
	}
	if !strings.Contains(formatted, "This is a comment") {
		t.Errorf("FormatInfo missing comment, got %q", formatted)
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

func BenchmarkRenderPonyWithBalloon(b *testing.B) {
	raw := `$$$
NAME: Test
$$$
$balloon5$
     $\$
  ▄▄▄▄▄▄▄
  █ █ █ █`
	p, _ := ParsePony("test", raw)
	balloonLines := []string{"┌───┐", "│Hi │", "└───┘"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.RenderPonyWithBalloon(balloonLines, "╲", "")
	}
}
