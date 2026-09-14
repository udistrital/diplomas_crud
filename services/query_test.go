package services

import "testing"

func TestParseQueryFilters(t *testing.T) {
	t.Parallel()

	filters, err := ParseQueryFilters("activo:true,codigo_estudiante:10000000001,uuid_documento:11111111-1111-1111-1111-111111111111")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if filters["activo"] != true {
		t.Fatalf("expected activo=true, got %#v", filters["activo"])
	}

	if filters["codigo_estudiante"] != int64(10000000001) {
		t.Fatalf("expected codigo_estudiante=10000000001, got %#v", filters["codigo_estudiante"])
	}

	if filters["uuid_documento"] != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("expected uuid_documento filter, got %#v", filters["uuid_documento"])
	}
}

func TestParseQueryFiltersRejectsInvalidFilter(t *testing.T) {
	t.Parallel()

	if _, err := ParseQueryFilters("activo"); err == nil {
		t.Fatal("expected error for invalid filter")
	}
}
