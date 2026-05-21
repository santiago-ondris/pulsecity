package domain

import "testing"

func TestNewHealth(t *testing.T) {
	health := NewHealth()

	if health.Service != ServiceName {
		t.Fatalf("Service = %q, want %q", health.Service, ServiceName)
	}
	if health.Status != "ok" {
		t.Fatalf("Status = %q, want ok", health.Status)
	}
}
