package assets

import (
	"testing"
)

func TestAssetManagerPonies(t *testing.T) {
	am := NewAssetManager()

	// Test listing ponies
	ponies := am.ListPonies(false, false)
	if len(ponies) == 0 {
		t.Fatalf("Expected non-empty list of MLP ponies")
	}

	// Test getting derpy pony
	name, content, err := am.GetPonyFile("derpy", false)
	if err != nil {
		t.Fatalf("Failed to get derpy pony: %v", err)
	}
	if name != "derpy" {
		t.Errorf("Expected pony name derpy, got %s", name)
	}
	if len(content) == 0 {
		t.Errorf("Pony content should not be empty")
	}

	// Test getting random pony
	rName, rContent, err := am.GetRandomPonyFile(false)
	if err != nil {
		t.Fatalf("Failed to get random pony: %v", err)
	}
	if rName == "" || len(rContent) == 0 {
		t.Errorf("Random pony should have valid name and content")
	}
}

func TestAssetManagerBalloons(t *testing.T) {
	am := NewAssetManager()

	balloons := am.ListBalloons(false)
	if len(balloons) == 0 {
		t.Fatalf("Expected non-empty list of balloon styles")
	}

	content, err := am.GetBalloonContent("unicode", false)
	if err != nil {
		t.Fatalf("Failed to get unicode balloon content: %v", err)
	}
	if len(content) == 0 {
		t.Errorf("Balloon content should not be empty")
	}
}

func TestAssetManagerQuotes(t *testing.T) {
	am := NewAssetManager()

	quoters := am.ListQuoters()
	if len(quoters) == 0 {
		t.Fatalf("Expected non-empty list of quoters")
	}

	// Test random quote
	pName, qText, err := am.GetPonyQuote(nil)
	if err != nil {
		t.Fatalf("Failed to get random pony quote: %v", err)
	}
	if pName == "" || qText == "" {
		t.Errorf("Quote result should have valid pony name and text")
	}

	// Test specific quote for Pinkie Pie
	pNamePinkie, qPinkie, err := am.GetPonyQuote([]string{"pinkie"})
	if err != nil {
		t.Fatalf("Failed to get Pinkie quote: %v", err)
	}
	if pNamePinkie != "pinkie" {
		t.Errorf("Expected quote pony to be pinkie, got %s", pNamePinkie)
	}
	if qPinkie == "" {
		t.Errorf("Pinkie quote should not be empty")
	}
}
