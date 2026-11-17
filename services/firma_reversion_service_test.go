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
	_ "github.com/lib/pq"
)

// Opt-in integration test. The supplied DSN must refer to a disposable local database.
// A separate schema is created and dropped; no application tables are touched.
func TestReversionPostgres(t *testing.T) {
	dsn := os.Getenv("DIPLOMAS_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("DIPLOMAS_TEST_POSTGRES_DSN not configured")
	}
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	if u.Hostname() != "127.0.0.1" && u.Hostname() != "localhost" {
		t.Fatal("integration database must be local")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	schema := fmt.Sprintf("test_reversion_%d", time.Now().UnixNano())
	if _, err := db.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	defer db.Exec("DROP SCHEMA " + schema + " CASCADE")
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	if err := orm.RegisterDataBase("default", "postgres", u.String()); err != nil {
		t.Fatal(err)
	}
	o := orm.NewOrm()
	for _, ddl := range []string{
		`CREATE TABLE documento_digital (id bigint PRIMARY KEY, estado_documento_id bigint NOT NULL, uuid_documento uuid, fecha_modificacion timestamp DEFAULT now())`,
		`CREATE TABLE firma_documento (id bigserial PRIMARY KEY, documento_digital_id bigint NOT NULL REFERENCES documento_digital(id), firma_firmante_id bigint, rol_firmante_id bigint NOT NULL, documento_identidad bigint NOT NULL, activo boolean NOT NULL DEFAULT true, fecha_modificacion timestamp DEFAULT now())`,
		`CREATE TABLE historico_estado_documento (id bigserial PRIMARY KEY, documento_digital_id bigint, estado_anterior_id bigint, estado_nuevo_id bigint, documento_identidad bigint, rol_actor_id bigint, firma_documento_id bigint, observacion text)`,
	} {
		if _, err := o.Raw(ddl).Exec(); err != nil {
			t.Fatal(err)
		}
	}
	// Columns used by the existing RelatedSel response after successful registration.
	for _, ddl := range []string{
		`ALTER TABLE documento_digital ADD COLUMN tipo_documento_digital_id bigint DEFAULT 1, ADD COLUMN codigo_estudiante bigint DEFAULT 1, ADD COLUMN facultad_id bigint, ADD COLUMN programa_academico_id bigint, ADD COLUMN periodo_id bigint, ADD COLUMN vigencia integer, ADD COLUMN nombre_estudiante text, ADD COLUMN tipo_documento_estudiante text, ADD COLUMN numero_documento_estudiante text, ADD COLUMN municipio_expedicion text, ADD COLUMN titulo_otorgado text, ADD COLUMN activo boolean DEFAULT true, ADD COLUMN fecha_creacion timestamp DEFAULT now()`,
		`ALTER TABLE firma_documento ADD COLUMN fecha_creacion timestamp DEFAULT now()`,
		`CREATE TABLE firma_firmante (id bigint PRIMARY KEY, documento_identidad bigint, enlace_firma text, activo boolean, fecha_creacion timestamp, fecha_modificacion timestamp)`,
	} {
		if _, err := o.Raw(ddl).Exec(); err != nil {
			t.Fatal(err)
		}
	}
	reset := func() {
		t.Helper()
		for _, query := range []string{"TRUNCATE historico_estado_documento, firma_documento, documento_digital RESTART IDENTITY CASCADE", "INSERT INTO documento_digital(id,estado_documento_id) VALUES (40,8)", "INSERT INTO firma_documento(documento_digital_id,rol_firmante_id,documento_identidad) VALUES (40,158,123)"} {
			if _, err := o.Raw(query).Exec(); err != nil {
				t.Fatal(err)
			}
		}
	}
	input := func() *RevertirFirmaDocumentoInput {
		return &RevertirFirmaDocumentoInput{DocumentoIdentidad: 123, RolFirmanteId: 158, FirmasEsperadas: []int64{1}, EstadoNuevoId: 9}
	}
	assertState := func(active, hist, state int) {
		t.Helper()
		var a, h, s int
		if err := o.Raw("SELECT count(*) FROM firma_documento WHERE activo").QueryRow(&a); err != nil {
			t.Fatal(err)
		}
		if err := o.Raw("SELECT count(*) FROM historico_estado_documento").QueryRow(&h); err != nil {
			t.Fatal(err)
		}
		if err := o.Raw("SELECT estado_documento_id FROM documento_digital WHERE id=40").QueryRow(&s); err != nil {
			t.Fatal(err)
		}
		if a != active || h != hist || s != state {
			t.Fatalf("active=%d hist=%d state=%d", a, h, s)
		}
	}
	t.Run("soft delete and audit", func(t *testing.T) {
		reset()
		result, err := (DocumentoDigitalWorkflowService{}).RevertirFirma(40, input())
		if err != nil || result.Anuladas != 1 {
			t.Fatalf("result=%v err=%v", result, err)
		}
		assertState(0, 1, 9)
		var actor, signature int64
		var obs string
		if err := o.Raw("SELECT documento_identidad,firma_documento_id,observacion FROM historico_estado_documento").QueryRow(&actor, &signature, &obs); err != nil || actor != 123 || signature != 1 || obs == "" {
			t.Fatalf("audit err=%v", err)
		}
	})
	t.Run("wrong owner no writes", func(t *testing.T) {
		reset()
		in := input()
		in.DocumentoIdentidad = 456
		_, err := (DocumentoDigitalWorkflowService{}).RevertirFirma(40, in)
		if !errors.Is(err, ErrFirmaConflict) {
			t.Fatal(err)
		}
		assertState(1, 0, 8)
	})
	t.Run("finalized no writes", func(t *testing.T) {
		reset()
		o.Raw("UPDATE documento_digital SET uuid_documento='00000000-0000-0000-0000-000000000001'").Exec()
		_, err := (DocumentoDigitalWorkflowService{}).RevertirFirma(40, input())
		if !errors.Is(err, ErrFirmaConflict) {
			t.Fatal(err)
		}
		assertState(1, 0, 8)
	})
	t.Run("concurrent signature blocks stale reversal", func(t *testing.T) {
		reset()
		tx, err := o.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		if err := bloquearFirmasDocumentoTx(tx, 40, []int64{1}, false); err != nil {
			t.Fatal(err)
		}
		if _, err := upsertFirmaDocumentoTx(tx, 40, &RegistrarFirmaDocumentoInput{RolFirmanteId: 270, DocumentoIdentidad: 456}); err != nil {
			t.Fatal(err)
		}
		done := make(chan error, 1)
		go func() { _, err := (DocumentoDigitalWorkflowService{}).RevertirFirma(40, input()); done <- err }()
		select {
		case err := <-done:
			t.Fatalf("reversal did not wait on document lock: %v", err)
		case <-time.After(100 * time.Millisecond):
		}
		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}
		select {
		case err := <-done:
			if !errors.Is(err, ErrFirmaConflict) {
				t.Fatal(err)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("reversal did not finish")
		}
		assertState(2, 0, 8)
	})
	t.Run("stale signing after reversal rejected", func(t *testing.T) {
		reset()
		if _, err := (DocumentoDigitalWorkflowService{}).RevertirFirma(40, input()); err != nil {
			t.Fatal(err)
		}
		uuid := "00000000-0000-0000-0000-000000000002"
		_, err := (DocumentoDigitalWorkflowService{}).RegistrarFirma(40, &RegistrarFirmaDocumentoInput{RolFirmanteId: 100, DocumentoIdentidad: 456, FirmasEsperadas: []int64{1}, UUIDDocumento: &uuid})
		if !errors.Is(err, ErrFirmaConflict) {
			t.Fatal(err)
		}
		var stored string
		if err := o.Raw("SELECT COALESCE(uuid_documento::text,'') FROM documento_digital WHERE id=40").QueryRow(&stored); err != nil || stored != "" {
			t.Fatalf("unexpected final UUID: %s error=%v", stored, err)
		}
		assertState(0, 1, 9)
	})
	t.Run("audit failure rolls back reversal", func(t *testing.T) {
		reset()
		if _, err := o.Raw("ALTER TABLE historico_estado_documento ADD CONSTRAINT reject_test_audit CHECK (observacion = 'not permitted')").Exec(); err != nil {
			t.Fatal(err)
		}
		defer o.Raw("ALTER TABLE historico_estado_documento DROP CONSTRAINT reject_test_audit").Exec()
		if _, err := (DocumentoDigitalWorkflowService{}).RevertirFirma(40, input()); err == nil {
			t.Fatal("expected audit failure")
		}
		assertState(1, 0, 8)
	})
	t.Run("final signature UUID and state commit together", func(t *testing.T) {
		reset()
		uuid := "00000000-0000-0000-0000-000000000003"
		estado := int64(10)
		firma, err := (DocumentoDigitalWorkflowService{}).RegistrarFirma(40, &RegistrarFirmaDocumentoInput{RolFirmanteId: 100, DocumentoIdentidad: 456, FirmasEsperadas: []int64{1}, UUIDDocumento: &uuid, EstadoFirmadoId: &estado})
		if err != nil {
			t.Fatal(err)
		}
		if firma.DocumentoDigital == nil || firma.DocumentoDigital.UUIDDocumento == nil || *firma.DocumentoDigital.UUIDDocumento != uuid || firma.DocumentoDigital.EstadoDocumentoId != estado {
			t.Fatalf("final document missing in response: %+v", firma)
		}
		assertState(2, 1, 10)
	})
	t.Run("legacy writes cannot bypass reversal", func(t *testing.T) {
		reset()
		if err := (FirmaDocumentoService{}).SoftDelete(1); !errors.Is(err, ErrFirmaConflict) {
			t.Fatal(err)
		}
		if _, err := (FirmaDocumentoService{}).Update(1, nil); !errors.Is(err, ErrFirmaConflict) {
			t.Fatal(err)
		}
		assertState(1, 0, 8)
	})
}
