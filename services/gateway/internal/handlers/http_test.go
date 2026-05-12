package handlers

import (
	"testing"

	"github.com/pulsecity/services/gateway/internal/domain"
)

func TestFindNarrativeChoice(t *testing.T) {
	choices := []domain.NarrativeChoice{
		{ID: "build_culture", Label: "Empezá por identidad y cultura"},
		{ID: "win_now", Label: "Acelerá para competir rapido"},
	}

	choice, ok := findNarrativeChoice(choices, "win_now")
	if !ok {
		t.Fatal("expected to find choice")
	}
	if choice.Label != "Acelerá para competir rapido" {
		t.Fatalf("unexpected choice label %q", choice.Label)
	}

	if _, ok := findNarrativeChoice(choices, "missing"); ok {
		t.Fatal("expected missing choice lookup to fail")
	}
}
