package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/lib/pq"
	"github.com/udistrital/diplomas_crud/models"
)

// Esquema exclusivo de prueba en PostgreSQL local: no modifica tablas de la aplicación.
func TestCeremoniaIntegridadPostgres(t *testing.T) {
	dsn := os.Getenv("DIPLOMAS_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("DIPLOMAS_TEST_POSTGRES_DSN not configured")
	}
	u, err := url.Parse(dsn)
	if err != nil || (u.Hostname() != "127.0.0.1" && u.Hostname() != "localhost") {
		t.Fatal("se requiere PostgreSQL local desechable")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	schema := fmt.Sprintf("test_ceremonia_%d", time.Now().UnixNano())
	ns := pq.QuoteIdentifier(schema)
	if _, err := db.Exec("CREATE SCHEMA " + ns); err != nil {
		t.Fatal(err)
	}
	defer db.Exec("DROP SCHEMA " + ns + " CASCADE")
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	if err := orm.RegisterDataBase(schema, "postgres", u.String()); err != nil {
		t.Fatal(err)
	}
	o := orm.NewOrmUsingDB(schema)
	exec := func(query string) {
		t.Helper()
		if _, err := o.Raw(query).Exec(); err != nil {
			t.Fatal(err)
		}
	}
	exec(`CREATE TABLE documento_digital (
  id bigint PRIMARY KEY, tipo_documento_digital_id bigint DEFAULT 1, estado_documento_id bigint DEFAULT 1,
  codigo_estudiante bigint DEFAULT 1, facultad_id bigint, programa_academico_id bigint, periodo_id bigint, vigencia integer,
  nombre_estudiante text, tipo_documento_estudiante text, numero_documento_estudiante text, municipio_expedicion text,
  titulo_otorgado text, uuid_documento uuid, activo boolean DEFAULT true, fecha_creacion timestamp DEFAULT now(),fecha_modificacion timestamp DEFAULT now());
 CREATE TABLE ceremonia_grado (
  id bigint PRIMARY KEY, periodo_id bigint,fecha_ceremonia date DEFAULT '2026-11-23',hora_ceremonia time,lugar text DEFAULT 'Prueba',
  cantidad_estudiantes_programada integer DEFAULT 10,observacion text,activo boolean DEFAULT true,fecha_creacion timestamp DEFAULT now(),fecha_modificacion timestamp DEFAULT now());
 CREATE TABLE ceremonia_facultad_programa (
  id bigint PRIMARY KEY,ceremonia_grado_id bigint NOT NULL REFERENCES ceremonia_grado(id),facultad_id bigint DEFAULT 14,
  programa_academico_id bigint DEFAULT 0,cantidad_estudiantes_programada integer,activo boolean DEFAULT true,
  fecha_creacion timestamp DEFAULT now(),fecha_modificacion timestamp DEFAULT now());
 CREATE TABLE ceremonia_estudiante (
  id bigserial PRIMARY KEY,ceremonia_facultad_programa_id bigint NOT NULL REFERENCES ceremonia_facultad_programa(id),
  documento_digital_id bigint NOT NULL REFERENCES documento_digital(id),activo boolean DEFAULT true,
  fecha_creacion timestamp DEFAULT now(),fecha_modificacion timestamp DEFAULT now(),UNIQUE(ceremonia_facultad_programa_id,documento_digital_id));
 INSERT INTO documento_digital(id) VALUES(1);`)
	ceremonias := CeremoniaGradoService{SimpleCRUDService: SimpleCRUDService[models.CeremoniaGrado]{db: o}}
	relaciones := CeremoniaFacultadProgramaService{SimpleCRUDService: SimpleCRUDService[models.CeremoniaFacultadPrograma]{db: o}}
	estudiantes := CeremoniaEstudianteService{SimpleCRUDService: SimpleCRUDService[models.CeremoniaEstudiante]{db: o}}
	nuevo := func(relacion int64) *models.CeremoniaEstudiante {
		return &models.CeremoniaEstudiante{CeremoniaFacultadPrograma: &models.CeremoniaFacultadPrograma{Id: relacion}, DocumentoDigital: &models.DocumentoDigital{Id: 1}, Activo: true}
	}
	reset := func() {
		t.Helper()
		exec(`TRUNCATE ceremonia_estudiante,ceremonia_facultad_programa,ceremonia_grado RESTART IDENTITY;
  INSERT INTO ceremonia_grado(id,activo) VALUES(1,true),(2,true),(3,false);
  INSERT INTO ceremonia_facultad_programa(id,ceremonia_grado_id,activo) VALUES(1,1,true),(2,1,true),(3,2,true),(4,3,true),(5,2,false);`)
	}
	assertActiva := func(want bool) {
		t.Helper()
		var active bool
		if err := o.Raw("SELECT activo FROM ceremonia_grado WHERE id=1").QueryRow(&active); err != nil || active != want {
			t.Fatalf("activa=%v want=%v err=%v", active, want, err)
		}
	}
	ocupada := func() {
		t.Helper()
		exec("INSERT INTO ceremonia_estudiante(ceremonia_facultad_programa_id,documento_digital_id) VALUES(2,1)")
	}
	for _, tc := range []struct {
		nombre string
		op     func() error
	}{
		{"DELETE ceremonia", func() error { return ceremonias.SoftDelete(1) }},
		{"PUT ceremonia inactiva", func() error {
			item, err := ceremonias.GetByID(1)
			if err != nil {
				return err
			}
			item.Activo = false
			_, err = ceremonias.Update(1, item)
			return err
		}},
		{"DELETE relacion", func() error { return relaciones.SoftDelete(2) }},
		{"PUT relacion inactiva", func() error {
			item, err := relaciones.GetByID(2)
			if err != nil {
				return err
			}
			item.Activo = false
			_, err = relaciones.Update(2, item)
			return err
		}},
		{"cambio de ceremonia de relacion", func() error {
			item, err := relaciones.GetByID(2)
			if err != nil {
				return err
			}
			item.CeremoniaGrado = &models.CeremoniaGrado{Id: 2}
			_, err = relaciones.Update(2, item)
			return err
		}},
	} {
		t.Run(tc.nombre, func(t *testing.T) {
			reset()
			ocupada()
			if err := tc.op(); !errors.Is(err, ErrCeremoniaConflict) {
				t.Fatalf("esperaba conflicto: %v", err)
			}
			assertActiva(true)
		})
	}
	t.Run("relacion historica inactiva con asignados", func(t *testing.T) {
		reset()
		ocupada()
		exec("UPDATE ceremonia_facultad_programa SET activo=false WHERE id=2")
		if err := ceremonias.SoftDelete(1); !errors.Is(err, ErrCeremoniaConflict) {
			t.Fatal(err)
		}
	})
	t.Run("inactivos no bloquean y no pueden reactivarse", func(t *testing.T) {
		reset()
		exec("INSERT INTO ceremonia_estudiante(ceremonia_facultad_programa_id,documento_digital_id,activo) VALUES(1,1,false)")
		if err := ceremonias.SoftDelete(1); err != nil {
			t.Fatal(err)
		}
		assertActiva(false)
		item := nuevo(1)
		if _, err := estudiantes.Update(1, item); !errors.Is(err, ErrCeremoniaConflict) {
			t.Fatalf("reactivacion: %v", err)
		}
	})
	for _, id := range []int64{4, 5} {
		t.Run(fmt.Sprintf("asignar padre inactivo %d", id), func(t *testing.T) {
			reset()
			if _, err := estudiantes.Create(nuevo(id)); !errors.Is(err, ErrCeremoniaConflict) {
				t.Fatal(err)
			}
		})
	}
	t.Run("asignacion valida y respuesta completa", func(t *testing.T) {
		reset()
		item, err := estudiantes.Create(nuevo(1))
		if err != nil || item.Id <= 0 || item.CeremoniaFacultadPrograma.CeremoniaGrado.Id != 1 {
			t.Fatalf("item=%+v err=%v", item, err)
		}
	})
	t.Run("ceremonia vacia por PUT", func(t *testing.T) {
		reset()
		item, err := ceremonias.GetByID(1)
		if err != nil {
			t.Fatal(err)
		}
		item.Activo = false
		result, err := ceremonias.Update(1, item)
		if err != nil || result.Activo {
			t.Fatalf("result=%+v err=%v", result, err)
		}
		assertActiva(false)
	})
	t.Run("otra ceremonia ocupada no bloquea vacia", func(t *testing.T) {
		reset()
		if _, err := estudiantes.Create(nuevo(3)); err != nil {
			t.Fatal(err)
		}
		if err := ceremonias.SoftDelete(1); err != nil {
			t.Fatal(err)
		}
		assertActiva(false)
	})
	t.Run("edicion activa conserva asignaciones", func(t *testing.T) {
		reset()
		ocupada()
		item, err := ceremonias.GetByID(1)
		if err != nil {
			t.Fatal(err)
		}
		item.Lugar = "Nuevo lugar"
		result, err := ceremonias.Update(1, item)
		if err != nil || result.Lugar != "Nuevo lugar" {
			t.Fatalf("result=%+v err=%v", result, err)
		}
	})
	t.Run("fallo de escritura libera transaccion", func(t *testing.T) {
		reset()
		item := nuevo(1)
		item.DocumentoDigital.Id = 999
		if _, err := estudiantes.Create(item); err == nil {
			t.Fatal("esperaba error FK")
		}
		if err := ceremonias.SoftDelete(1); err != nil {
			t.Fatal(err)
		}
	})
	for _, tc := range []struct {
		nombre  string
		primera func(orm.TxOrmer) error
		segunda func() error
		activa  bool
	}{
		{"asignacion primero", func(tx orm.TxOrmer) error {
			if err := bloquearAsignacionActivaTx(tx, 1); err != nil {
				return err
			}
			_, err := tx.Insert(nuevo(1))
			return err
		}, func() error { return ceremonias.SoftDelete(1) }, true},
		{"inactivacion primero", func(tx orm.TxOrmer) error {
			if _, err := bloquearCeremoniaTx(tx, 1); err != nil {
				return err
			}
			if err := validarCeremoniaSinEstudiantesTx(tx, 1); err != nil {
				return err
			}
			_, err := tx.Raw("UPDATE ceremonia_grado SET activo=false WHERE id=1").Exec()
			return err
		}, func() error { _, err := estudiantes.Create(nuevo(1)); return err }, false},
		{"relacion inactiva primero", func(tx orm.TxOrmer) error {
			var id int64
			if err := tx.Raw("SELECT id FROM ceremonia_facultad_programa WHERE id=1 FOR UPDATE").QueryRow(&id); err != nil {
				return err
			}
			if err := validarRelacionSinEstudiantesTx(tx, 1); err != nil {
				return err
			}
			_, err := tx.Raw("UPDATE ceremonia_facultad_programa SET activo=false WHERE id=1").Exec()
			return err
		}, func() error { _, err := estudiantes.Create(nuevo(1)); return err }, true},
	} {
		t.Run(tc.nombre, func(t *testing.T) {
			reset()
			tx, err := o.Begin()
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()
			if err := tc.primera(tx); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- tc.segunda() }()
			select {
			case err := <-done:
				t.Fatalf("no espero bloqueo: %v", err)
			case <-time.After(150 * time.Millisecond):
			}
			if err := tx.Commit(); err != nil {
				t.Fatal(err)
			}
			select {
			case err := <-done:
				if !errors.Is(err, ErrCeremoniaConflict) {
					t.Fatalf("esperaba conflicto: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("no termino")
			}
			assertActiva(tc.activa)
		})
	}
}
