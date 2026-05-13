package persistence

import (
	"testing"

	"github.com/pulsecity/services/gateway/internal/domain"
)

func TestOwnerKindForGame(t *testing.T) {
	if ownerKindForGame("guest_123", "") != domain.OwnerKindGuest {
		t.Fatal("expected guest owner kind")
	}
	if ownerKindForGame("", "user_123") != domain.OwnerKindUser {
		t.Fatal("expected user owner kind")
	}
	if ownerKindForGame("", "") != "" {
		t.Fatal("expected empty owner kind for invalid ownership")
	}
}

func TestHasExclusiveOwner(t *testing.T) {
	if !hasExclusiveOwner("guest_123", "") {
		t.Fatal("expected guest-only ownership to be valid")
	}
	if !hasExclusiveOwner("", "user_123") {
		t.Fatal("expected user-only ownership to be valid")
	}
	if hasExclusiveOwner("guest_123", "user_123") {
		t.Fatal("expected mixed ownership to be invalid")
	}
	if hasExclusiveOwner("", "") {
		t.Fatal("expected missing ownership to be invalid")
	}
}
