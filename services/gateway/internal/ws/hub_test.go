package ws

import "testing"

func TestGameActivationTracksFirstAndLastConnection(t *testing.T) {
	hub := NewHub()

	if !hub.ActivateGame("game-1", "session-1") {
		t.Fatal("expected first connection to activate game session")
	}
	if hub.ActivateGame("game-1", "session-2") {
		t.Fatal("expected second connection to reuse active game session")
	}

	sessionID, last := hub.DeactivateGame("game-1")
	if last {
		t.Fatal("expected first disconnect to keep game session active")
	}
	if sessionID != "session-1" {
		t.Fatalf("expected original session id, got %q", sessionID)
	}

	sessionID, last = hub.DeactivateGame("game-1")
	if !last {
		t.Fatal("expected second disconnect to close game session")
	}
	if sessionID != "session-1" {
		t.Fatalf("expected original session id, got %q", sessionID)
	}
}
