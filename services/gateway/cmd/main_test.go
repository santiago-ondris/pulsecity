package main

import (
	"testing"

	"github.com/pulsecity/services/gateway/internal/domain"
)

func TestBuildOwnerIntroEventIncludesInteractiveChoices(t *testing.T) {
	game := domain.GameSetup{
		GameID:             "game-123",
		CityName:           "Nueva Aurora",
		FranchiseName:      "Lighthouses",
		InitialScenario:    "expansion",
		CityManagementMode: "owner_influence",
	}

	event := buildOwnerIntroEvent(game)

	if event.Type != "narrative.event" {
		t.Fatalf("expected narrative.event type, got %q", event.Type)
	}
	if event.Kind != "owner_intro" {
		t.Fatalf("expected owner_intro kind, got %q", event.Kind)
	}
	if len(event.Choices) != 3 {
		t.Fatalf("expected 3 owner intro choices, got %d", len(event.Choices))
	}
	if event.Metadata["city_management_mode"] != "owner_influence" {
		t.Fatalf("expected metadata city_management_mode to be owner_influence, got %q", event.Metadata["city_management_mode"])
	}
}
