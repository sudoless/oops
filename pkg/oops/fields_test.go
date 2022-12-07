package oops

import "testing"

func TestError_Fields(t *testing.T) {
	t.Parallel()

	err := errTest.Yeet("")

	if err.FieldsMap() != nil {
		t.Errorf("expected nil fields list, got %v", err.FieldsMap())
	}

	_ = err.Fields(
		F("k1", "v1"),
		F("k2", "v2"),
		F("k3", "v3"),
	)

	if len(err.FieldsMap()) != 3 {
		t.Errorf("expected 3 fields, got %d", len(err.FieldsMap()))
	}

	fieldsMap := err.FieldsMap()

	if len(fieldsMap) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fieldsMap))
	}

	if fieldsMap["k1"] != "v1" {
		t.Errorf("expected field k1=v1, got %s", fieldsMap["k1"])
	}

	if fieldsMap["k2"] != "v2" {
		t.Errorf("expected field k2=v2, got %s", fieldsMap["k2"])
	}

	if fieldsMap["k3"] != "v3" {
		t.Errorf("expected field k3=v3, got %s", fieldsMap["k3"])
	}
}

func TestError_FieldsMap_overwrite(t *testing.T) {
	t.Parallel()

	err := errTest.Yeet("")

	_ = err.Fields(F("k1", "v1"), F("k2", "v2"))
	_ = err.Fields(F("k3", "v3"), F("k1", "v3"))

	fieldsMap := err.FieldsMap()

	if len(fieldsMap) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fieldsMap))
	}

	if fieldsMap["k1"] != "v3" {
		t.Errorf("expected field k1=v3, got %s", fieldsMap["k1"])
	}
}
