package main

import "github.com/beego/beego/v2/client/orm/migration"

type AddDocumentoDigitalSGASnapshot20260915 struct {
	migration.Migration
}

func init() {
	m := &AddDocumentoDigitalSGASnapshot20260915{}
	m.Created = "20260915_000002"
	migration.Register("AddDocumentoDigitalSGASnapshot20260915", m)
}

func (m *AddDocumentoDigitalSGASnapshot20260915) Up() {
	m.SQL(`
ALTER TABLE diplomas.documento_digital
ADD COLUMN IF NOT EXISTS facultad_id INTEGER,
ADD COLUMN IF NOT EXISTS nombre_estudiante CHARACTER VARYING(250),
ADD COLUMN IF NOT EXISTS tipo_documento_estudiante CHARACTER VARYING(20),
ADD COLUMN IF NOT EXISTS numero_documento_estudiante CHARACTER VARYING(50),
ADD COLUMN IF NOT EXISTS municipio_expedicion CHARACTER VARYING(150),
ADD COLUMN IF NOT EXISTS titulo_otorgado CHARACTER VARYING(250);

CREATE INDEX IF NOT EXISTS idx_documento_digital_numero_documento_estudiante
ON diplomas.documento_digital (numero_documento_estudiante);

CREATE INDEX IF NOT EXISTS idx_documento_digital_facultad_id
ON diplomas.documento_digital (facultad_id);

COMMENT ON COLUMN diplomas.documento_digital.facultad_id IS 'Referencia externa a la facultad del estudiante en SGA/Oikos para asignar consecutivos del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.nombre_estudiante IS 'Nombre completo del estudiante consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.tipo_documento_estudiante IS 'Tipo de documento de identidad del estudiante consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.numero_documento_estudiante IS 'Numero de documento del estudiante consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.municipio_expedicion IS 'Municipio de expedicion del documento del estudiante consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.titulo_otorgado IS 'Titulo academico otorgado consultado desde SGA 1 para preparar la generacion del diploma.';
`)
}

func (m *AddDocumentoDigitalSGASnapshot20260915) Down() {
	m.SQL(`
DROP INDEX IF EXISTS diplomas.idx_documento_digital_numero_documento_estudiante;
DROP INDEX IF EXISTS diplomas.idx_documento_digital_facultad_id;

ALTER TABLE diplomas.documento_digital
DROP COLUMN IF EXISTS facultad_id,
DROP COLUMN IF EXISTS nombre_estudiante,
DROP COLUMN IF EXISTS tipo_documento_estudiante,
DROP COLUMN IF EXISTS numero_documento_estudiante,
DROP COLUMN IF EXISTS municipio_expedicion,
DROP COLUMN IF EXISTS titulo_otorgado;
`)
}
