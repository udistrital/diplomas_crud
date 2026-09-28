package services

import (
	"reflect"
	"testing"

	"github.com/udistrital/diplomas_crud/models"
)

func TestDiplomaDigitalUsaRegistroGradoSeparado(t *testing.T) {
	t.Parallel()

	tipo := reflect.TypeOf(models.DiplomaDigital{})
	for _, field := range []string{"Libro", "Folio", "Acta"} {
		if _, ok := tipo.FieldByName(field); ok {
			t.Fatalf("DiplomaDigital no debe exponer campo legacy %s", field)
		}
	}
	if _, ok := tipo.FieldByName("RegistroGrado"); !ok {
		t.Fatalf("DiplomaDigital debe exponer RegistroGrado para libro/folio/registro")
	}
}
