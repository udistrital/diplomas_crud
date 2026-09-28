package main

import "github.com/beego/beego/v2/client/orm/migration"

type RefineDocumentoDigitalModel20260915 struct {
	migration.Migration
}

func init() {
	m := &RefineDocumentoDigitalModel20260915{}
	m.Created = "20260915_000001"
	migration.Register("RefineDocumentoDigitalModel20260915", m)
}

func (m *RefineDocumentoDigitalModel20260915) Up() {
	m.SQL(`
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'documento_digital'
          AND column_name = 'tipo_documento_id'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'documento_digital'
          AND column_name = 'tipo_documento_digital_id'
    ) THEN
        ALTER TABLE diplomas.documento_digital
        RENAME COLUMN tipo_documento_id TO tipo_documento_digital_id;
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'diploma_digital'
          AND column_name = 'nombre_graduando'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'diploma_digital'
          AND column_name = 'nombre_estudiante'
    ) THEN
        ALTER TABLE diplomas.diploma_digital
        RENAME COLUMN nombre_graduando TO nombre_estudiante;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'diploma_digital'
          AND column_name = 'documento_identidad'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'diploma_digital'
          AND column_name = 'numero_documento_estudiante'
    ) THEN
        ALTER TABLE diplomas.diploma_digital
        RENAME COLUMN documento_identidad TO numero_documento_estudiante;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'diploma_digital'
          AND column_name = 'titulo_conferido'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'diploma_digital'
          AND column_name = 'titulo_otorgado'
    ) THEN
        ALTER TABLE diplomas.diploma_digital
        RENAME COLUMN titulo_conferido TO titulo_otorgado;
    END IF;
END $$;

ALTER TABLE diplomas.diploma_digital
ADD COLUMN IF NOT EXISTS tipo_documento_estudiante CHARACTER VARYING(20);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas'
          AND table_name = 'diploma_digital'
          AND column_name = 'tipo_documento'
    ) THEN
        UPDATE diplomas.diploma_digital
        SET tipo_documento_estudiante = tipo_documento
        WHERE tipo_documento_estudiante IS NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM diplomas.diploma_digital
        WHERE tipo_documento_estudiante IS NULL OR btrim(tipo_documento_estudiante) = ''
    ) THEN
        ALTER TABLE diplomas.diploma_digital
        ALTER COLUMN tipo_documento_estudiante SET NOT NULL;
    END IF;
END $$;

ALTER TABLE diplomas.diploma_digital
DROP COLUMN IF EXISTS tipo_documento,
DROP COLUMN IF EXISTS tipo_documento_estudiante_id;

DROP INDEX IF EXISTS diplomas.idx_documento_digital_tipo_documento_id;

CREATE INDEX IF NOT EXISTS idx_documento_digital_tipo_documento_digital_id
ON diplomas.documento_digital (tipo_documento_digital_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_documento_digital_tramite_activo
ON diplomas.documento_digital (
    tipo_documento_digital_id,
    codigo_estudiante,
    COALESCE(programa_academico_id, 0),
    COALESCE(periodo_id, 0),
    COALESCE(vigencia, 0)
)
WHERE activo IS TRUE;

COMMENT ON COLUMN diplomas.documento_digital.tipo_documento_digital_id IS 'Referencia externa al parametro que identifica el tipo de documento digital emitido: acta de grado, diploma normal o diploma IDE. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.diploma_digital.nombre_estudiante IS 'Nombre completo del estudiante tal como queda impreso en el diploma emitido. Se copia desde la fuente academica/personas al momento de emision para trazabilidad historica.';
COMMENT ON COLUMN diplomas.diploma_digital.tipo_documento_estudiante IS 'Tipo de documento de identidad del estudiante tal como llega desde SGA 1 al momento de emision.';
COMMENT ON COLUMN diplomas.diploma_digital.numero_documento_estudiante IS 'Numero de documento del estudiante tal como llega desde SGA 1 y queda impreso en el diploma emitido.';
COMMENT ON COLUMN diplomas.diploma_digital.municipio_expedicion IS 'Municipio de expedicion del documento del estudiante tal como llega desde SGA 1.';
COMMENT ON COLUMN diplomas.diploma_digital.titulo_otorgado IS 'Titulo academico otorgado tal como llega desde SGA 1 y queda impreso en el diploma emitido.';
`)
}

func (m *RefineDocumentoDigitalModel20260915) Down() {
	m.SQL(`
DROP INDEX IF EXISTS diplomas.uq_documento_digital_tramite_activo;
DROP INDEX IF EXISTS diplomas.idx_documento_digital_tipo_documento_digital_id;
`)
}
