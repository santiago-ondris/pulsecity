package handlers

import (
	"net/http/httptest"
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

func TestGuestTokenFromRequest(t *testing.T) {
	request := httptest.NewRequest("GET", "/api/v1/games", nil)
	request.Header.Set("X-Guest-Token", " guest_123 ")

	token := guestTokenFromRequest(request)
	if token != "guest_123" {
		t.Fatalf("unexpected guest token %q", token)
	}
}

func TestGuestOwnsGame(t *testing.T) {
	game := domain.GameSetup{GameID: "game-1", GuestToken: "guest_123"}

	if !guestOwnsGame("guest_123", game) {
		t.Fatal("expected matching guest to own game")
	}
	if guestOwnsGame("guest_other", game) {
		t.Fatal("expected guest mismatch to fail ownership check")
	}
	if guestOwnsGame("", game) {
		t.Fatal("expected empty guest token to fail ownership check")
	}
}

func TestGameOwnedByUser(t *testing.T) {
	game := domain.GameSetup{GameID: "game-1", UserID: "user_123"}
	currentActor := actor{
		kind: "user",
		user: domain.User{UserID: "user_123"},
	}

	if !gameOwnedBy(currentActor, game) {
		t.Fatal("expected matching user to own game")
	}

	currentActor.user.UserID = "user_999"
	if gameOwnedBy(currentActor, game) {
		t.Fatal("expected user mismatch to fail ownership check")
	}
}

func TestValidateRegistrationInput(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		display string
		pass    string
		wantErr bool
	}{
		{name: "valid", email: "gm@pulsecity.test", display: "Jordan Vale", pass: "superpass", wantErr: false},
		{name: "missing email", email: "", display: "Jordan Vale", pass: "superpass", wantErr: true},
		{name: "invalid email", email: "bad", display: "Jordan Vale", pass: "superpass", wantErr: true},
		{name: "short display", email: "gm@pulsecity.test", display: "Jo", pass: "superpass", wantErr: true},
		{name: "short password", email: "gm@pulsecity.test", display: "Jordan Vale", pass: "short", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateRegistrationInput(test.email, test.display, test.pass)
			if test.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateLoginInput(t *testing.T) {
	if err := validateLoginInput("gm@pulsecity.test", "secretpass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateLoginInput("", "secretpass"); err == nil {
		t.Fatal("expected missing email to fail")
	}
	if err := validateLoginInput("bad", "secretpass"); err == nil {
		t.Fatal("expected invalid email to fail")
	}
	if err := validateLoginInput("gm@pulsecity.test", ""); err == nil {
		t.Fatal("expected empty password to fail")
	}
}
