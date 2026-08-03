package assets

import (
	"testing"
)

func TestAssetManagerPonies(t *testing.T) {
	am := NewAssetManager()

	// Test listing ponies
	ponies := am.ListPonies(true, false)
	if len(ponies) == 0 {
		t.Fatalf("Expected non-empty list of MLP ponies")
	}

	// Test getting derpy pony
	name, content, err := am.GetPonyFile("derpy", true, false)
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
	rName, rContent, err := am.GetRandomPonyFile(true, false)
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

	quoters := am.ListQuoters(true, false)
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

	// Test specific quote for Rara (Countess Coloratura)
	pNameRara, qRara, err := am.GetPonyQuote([]string{"rara"})
	if err != nil {
		t.Fatalf("Failed to get Rara quote: %v", err)
	}
	if pNameRara != "rara" || qRara == "" {
		t.Errorf("Expected quote pony to be rara, got %s: %s", pNameRara, qRara)
	}

	// Test fuzzy search in quote lookup (e.g., fluter-shy -> fluttershy)
	pNameFuzzy, qFuzzy, err := am.GetPonyQuote([]string{"fluter-shy"})
	if err != nil {
		t.Fatalf("Failed to get fuzzy quote for fluter-shy: %v", err)
	}
	if pNameFuzzy != "fluttershy" {
		t.Errorf("Expected fuzzy quote pony to resolve to fluttershy, got %s", pNameFuzzy)
	}
	if qFuzzy == "" || qFuzzy == "Zecora! Help me, I am mute!" {
		t.Errorf("Expected valid Fluttershy quote, got: %s", qFuzzy)
	}

	// Test variant quote lookup (e.g. lunafly should resolve pony name to lunafly with luna quote)
	pNameVariant, qVariant, err := am.GetPonyQuote([]string{"lunafly"})
	if err != nil {
		t.Fatalf("Failed to get variant quote for lunafly: %v", err)
	}
	if pNameVariant != "lunafly" {
		t.Errorf("Expected variant quote pony to resolve to lunafly, got %s", pNameVariant)
	}
	if qVariant == "" || qVariant == "Zecora! Help me, I am mute!" {
		t.Errorf("Expected valid Luna quote for lunafly variant, got: %s", qVariant)
	}

	// Test MASTER metadata tag quote lookup (e.g. woona has MASTER: luna)
	pNameWoona, qWoona, err := am.GetPonyQuote([]string{"woona"})
	if err != nil {
		t.Fatalf("Failed to get quote for woona: %v", err)
	}
	if pNameWoona != "woona" {
		t.Errorf("Expected MASTER quote pony to resolve to woona, got %s", pNameWoona)
	}
	if qWoona == "" || qWoona == "Zecora! Help me, I am mute!" {
		t.Errorf("Expected valid Luna quote for woona (via MASTER tag), got: %s", qWoona)
	}
}

func TestAssetManagerConcurrency(t *testing.T) {
	am := NewAssetManager()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				_, _, _ = am.GetPonyFile("derpy", true, false)
				_, _, _ = am.GetRandomPonyFile(true, false)
				_ = am.ListPonies(true, true)
				_, _ = am.GetBalloonContent("cowsay", false)
				_ = am.ListBalloons(false)
				_, _, _ = am.GetPonyQuote([]string{"pinkie"})
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
