package controllers

import (
	"errors"
	"fmt"
	"github.com/udistrital/diplomas_crud/services"
	"net/http"
	"testing"
)

func TestMapCeremoniaConflict(t *testing.T) {
	for _, tc := range []struct {
		nombre string
		err    error
		want   int
	}{
		{"conflicto de ceremonia", services.ErrCeremoniaConflict, http.StatusConflict},
		{"conflicto envuelto", fmt.Errorf("update ceremonia: %w", services.ErrCeremoniaConflict), http.StatusConflict},
		{"otro error", errors.New("otro error"), http.StatusInternalServerError},
	} {
		t.Run(tc.nombre, func(t *testing.T) {
			if got := mapReadError(tc.err); got != tc.want {
				t.Fatalf("status=%d want=%d", got, tc.want)
			}
		})
	}
}
