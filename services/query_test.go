package services

import "testing"

func TestParseQueryFilters(t *testing.T) {
	t.Parallel()

	filters, err := ParseQueryFilters("activo:true,tercero_id:123,hash_documento:abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if filters["activo"] != true {
		t.Fatalf("expected activo=true, got %#v", filters["activo"])
	}

	if filters["tercero_id"] != int64(123) {
		t.Fatalf("expected tercero_id=123, got %#v", filters["tercero_id"])
	}

	if filters["hash_documento"] != "abc" {
		t.Fatalf("expected hash_documento=abc, got %#v", filters["hash_documento"])
	}
}

func TestParseQueryFiltersRejectsInvalidFilter(t *testing.T) {
	t.Parallel()

	if _, err := ParseQueryFilters("activo"); err == nil {
		t.Fatal("expected error for invalid filter")
	}
}
