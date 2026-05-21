package domain

import "testing"

func TestEventSubjects(t *testing.T) {
	tests := map[string]string{
		"day advanced":    SubjectTimeDayAdvanced,
		"match scheduled": SubjectMatchScheduled,
		"match finished":  SubjectMatchFinished,
	}

	if tests["day advanced"] != "tiempo.dia_avanzado" {
		t.Fatalf("SubjectTimeDayAdvanced = %q", tests["day advanced"])
	}
	if tests["match scheduled"] != "partido.programado" {
		t.Fatalf("SubjectMatchScheduled = %q", tests["match scheduled"])
	}
	if tests["match finished"] != "partido.terminado" {
		t.Fatalf("SubjectMatchFinished = %q", tests["match finished"])
	}
}
