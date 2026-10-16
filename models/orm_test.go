package models

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/beego/beego/v2/client/orm"
)

func TestModelosRegistradosAutomaticamente(t *testing.T) {
	modelos := []struct {
		nombre string
		modelo interface{}
	}{
		{"documento_digital", new(DocumentoDigital)},
		{"parametrizacion_registro_facultad", new(ParametrizacionRegistroFacultad)},
		{"libro_registro", new(LibroRegistro)},
		{"folio_registro", new(FolioRegistro)},
		{"control_registro_facultad", new(ControlRegistroFacultad)},
		{"diploma_digital", new(DiplomaDigital)},
		{"registro_grado", new(RegistroGrado)},
		{"ceremonia_grado", new(CeremoniaGrado)},
		{"ceremonia_facultad_programa", new(CeremoniaFacultadPrograma)},
		{"ceremonia_estudiante", new(CeremoniaEstudiante)},
		{"firma_firmante", new(FirmaFirmante)},
		{"firma_documento", new(FirmaDocumento)},
		{"historico_estado_documento", new(HistoricoEstadoDocumento)},
	}

	for _, tc := range modelos {
		t.Run(tc.nombre, func(t *testing.T) {
			assertModeloRegistrado(t, tc.modelo)
		})
	}
}

func TestBootstrapResuelveRelacionesYNombresDeTabla(t *testing.T) {
	MustInitORM()
	db := sql.OpenDB(conectorORMPrueba{})
	defer db.Close()

	o, err := orm.NewOrmWithDB(defaultDriver, "modelos_registrados_test", db)
	if err != nil {
		t.Fatalf("inicializar Beego ORM: %v", err)
	}

	modelos := []struct {
		nombre string
		modelo interface{}
	}{
		{"documento_digital", new(DocumentoDigital)},
		{"parametrizacion_registro_facultad", new(ParametrizacionRegistroFacultad)},
		{"libro_registro", new(LibroRegistro)},
		{"folio_registro", new(FolioRegistro)},
		{"control_registro_facultad", new(ControlRegistroFacultad)},
		{"diploma_digital", new(DiplomaDigital)},
		{"registro_grado", new(RegistroGrado)},
		{"ceremonia_grado", new(CeremoniaGrado)},
		{"ceremonia_facultad_programa", new(CeremoniaFacultadPrograma)},
		{"ceremonia_estudiante", new(CeremoniaEstudiante)},
		{"firma_firmante", new(FirmaFirmante)},
		{"firma_documento", new(FirmaDocumento)},
		{"historico_estado_documento", new(HistoricoEstadoDocumento)},
	}

	for _, tc := range modelos {
		if consulta := o.QueryTable(tc.modelo); consulta == nil {
			t.Errorf("modelo %s sin QuerySet", tc.nombre)
		}
		if consulta := o.QueryTable(tc.nombre); consulta == nil {
			t.Errorf("tabla %s sin QuerySet", tc.nombre)
		}
	}
}

type conectorORMPrueba struct{}

func (conectorORMPrueba) Connect(context.Context) (driver.Conn, error) {
	return conexionORMPrueba{}, nil
}

func (conectorORMPrueba) Driver() driver.Driver { return driverORMPrueba{} }

type driverORMPrueba struct{}

func (driverORMPrueba) Open(string) (driver.Conn, error) { return conexionORMPrueba{}, nil }

type conexionORMPrueba struct{}

func (conexionORMPrueba) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("la prueba de registro ORM no ejecuta consultas")
}

func (conexionORMPrueba) Close() error { return nil }

func (conexionORMPrueba) Begin() (driver.Tx, error) {
	return nil, errors.New("la prueba de registro ORM no inicia transacciones")
}

func (conexionORMPrueba) Ping(context.Context) error { return nil }

func assertModeloRegistrado(t *testing.T, modelo interface{}) {
	t.Helper()
	defer func() {
		recuperado := recover()
		if recuperado == nil {
			t.Fatal("el modelo no estaba registrado al importar el paquete models")
		}
		if !strings.Contains(fmt.Sprint(recuperado), "repeat Register") {
			t.Fatalf("RegisterModel produjo un panic inesperado: %v", recuperado)
		}
	}()
	orm.RegisterModel(modelo)
}
