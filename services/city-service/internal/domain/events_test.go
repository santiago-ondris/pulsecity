package domain

import "testing"

func TestEventSubjects(t *testing.T) {
	if SubjectMatchFinished != "partido.terminado" {
		t.Fatalf("SubjectMatchFinished = %q", SubjectMatchFinished)
	}
	if SubjectCityEconomyChange != "ciudad.economia_cambio" {
		t.Fatalf("SubjectCityEconomyChange = %q", SubjectCityEconomyChange)
	}
	if SubjectCityLandUpdated != "ciudad.suelo_actualizado" {
		t.Fatalf("SubjectCityLandUpdated = %q", SubjectCityLandUpdated)
	}
	if SubjectCityPatchDelta != "city.patch" {
		t.Fatalf("SubjectCityPatchDelta = %q", SubjectCityPatchDelta)
	}
}
