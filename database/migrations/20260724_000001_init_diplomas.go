package main

import "github.com/astaxie/beego/migration"

type InitDiplomas20260724 struct {
	migration.Migration
}

func init() {
	m := &InitDiplomas20260724{}
	m.Created = "20260724_000001"
	migration.Register("InitDiplomas20260724", m)
}

func (m *InitDiplomas20260724) Up() {
	m.SQL(`
CREATE SCHEMA IF NOT EXISTS diplomas;

COMMENT ON SCHEMA diplomas IS 'Esquema para la gestion de hashes, metadata, verificacion e historico de estados de diplomas digitales.';

CREATE TABLE IF NOT EXISTS diplomas.documento_digital (
    id SERIAL NOT NULL,
    tipo_documento_id INTEGER NOT NULL,
    estado_documento_id INTEGER NOT NULL,
    tercero_id INTEGER NOT NULL,
    programa_academico_id INTEGER,
    periodo_id INTEGER,
    vigencia INTEGER,
    hash_documento CHARACTER VARYING(128) NOT NULL,
    uuid_verificacion CHARACTER VARYING(100) NOT NULL,
    metadata JSONB,
    activo BOOLEAN NOT NULL DEFAULT true,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_documento_digital PRIMARY KEY (id),
    CONSTRAINT uq_hash_documento_documento_digital UNIQUE (hash_documento),
    CONSTRAINT uq_uuid_verificacion_documento_digital UNIQUE (uuid_verificacion)
);

COMMENT ON TABLE diplomas.documento_digital IS 'Tabla principal del documento digital. No almacena PDF ni base64. Guarda hash, metadata minima, estado actual e identificador de verificacion.';
COMMENT ON COLUMN diplomas.documento_digital.id IS 'Identificador unico del registro.';
COMMENT ON COLUMN diplomas.documento_digital.tipo_documento_id IS 'Referencia logica al parametro existente que identifica el tipo de documento digital.';
COMMENT ON COLUMN diplomas.documento_digital.estado_documento_id IS 'Referencia logica al parametro existente que identifica el estado actual del documento digital.';
COMMENT ON COLUMN diplomas.documento_digital.tercero_id IS 'Identificador del tercero titular del documento.';
COMMENT ON COLUMN diplomas.documento_digital.programa_academico_id IS 'Identificador del programa academico asociado.';
COMMENT ON COLUMN diplomas.documento_digital.periodo_id IS 'Identificador del periodo academico asociado.';
COMMENT ON COLUMN diplomas.documento_digital.vigencia IS 'Vigencia asociada al documento digital.';
COMMENT ON COLUMN diplomas.documento_digital.hash_documento IS 'Hash del documento firmado usado para verificacion.';
COMMENT ON COLUMN diplomas.documento_digital.uuid_verificacion IS 'Identificador unico utilizado por el QR para consultar la verificacion sin exponer el id interno.';
COMMENT ON COLUMN diplomas.documento_digital.metadata IS 'Metadata minima asociada al documento digital.';
COMMENT ON COLUMN diplomas.documento_digital.activo IS 'Indica si el registro se encuentra activo.';
COMMENT ON COLUMN diplomas.documento_digital.fecha_creacion IS 'Fecha de creacion del registro.';
COMMENT ON COLUMN diplomas.documento_digital.fecha_modificacion IS 'Fecha de ultima modificacion del registro.';

CREATE TABLE IF NOT EXISTS diplomas.historico_estado_documento (
    id SERIAL NOT NULL,
    documento_digital_id INTEGER NOT NULL,
    estado_anterior_id INTEGER,
    estado_nuevo_id INTEGER NOT NULL,
    tercero_id INTEGER NOT NULL,
    observacion CHARACTER VARYING(500),
    metadata JSONB,
    activo BOOLEAN NOT NULL DEFAULT true,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_historico_estado_documento PRIMARY KEY (id),
    CONSTRAINT fk_historico_estado_documento_documento_digital
        FOREIGN KEY (documento_digital_id)
        REFERENCES diplomas.documento_digital(id)
);

COMMENT ON TABLE diplomas.historico_estado_documento IS 'Historico de cambios de estado del documento digital. Permite auditar estado anterior, estado nuevo y tercero que realizo el cambio.';
COMMENT ON COLUMN diplomas.historico_estado_documento.id IS 'Identificador unico del registro.';
COMMENT ON COLUMN diplomas.historico_estado_documento.documento_digital_id IS 'Referencia al documento digital asociado.';
COMMENT ON COLUMN diplomas.historico_estado_documento.estado_anterior_id IS 'Referencia logica al parametro existente que identifica el estado anterior.';
COMMENT ON COLUMN diplomas.historico_estado_documento.estado_nuevo_id IS 'Referencia logica al parametro existente que identifica el nuevo estado.';
COMMENT ON COLUMN diplomas.historico_estado_documento.tercero_id IS 'Identificador del tercero que realizo el cambio de estado.';
COMMENT ON COLUMN diplomas.historico_estado_documento.observacion IS 'Observacion asociada al cambio de estado.';
COMMENT ON COLUMN diplomas.historico_estado_documento.metadata IS 'Metadata adicional del evento de cambio de estado.';
COMMENT ON COLUMN diplomas.historico_estado_documento.activo IS 'Indica si el registro se encuentra activo.';
COMMENT ON COLUMN diplomas.historico_estado_documento.fecha_creacion IS 'Fecha de creacion del registro.';
COMMENT ON COLUMN diplomas.historico_estado_documento.fecha_modificacion IS 'Fecha de ultima modificacion del registro.';

CREATE INDEX IF NOT EXISTS idx_documento_digital_tipo_documento_id
ON diplomas.documento_digital(tipo_documento_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_estado_documento_id
ON diplomas.documento_digital(estado_documento_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_tercero_id
ON diplomas.documento_digital(tercero_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_programa_academico_id
ON diplomas.documento_digital(programa_academico_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_periodo_id
ON diplomas.documento_digital(periodo_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_vigencia
ON diplomas.documento_digital(vigencia);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_documento_digital_id
ON diplomas.historico_estado_documento(documento_digital_id);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_estado_nuevo_id
ON diplomas.historico_estado_documento(estado_nuevo_id);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_tercero_id
ON diplomas.historico_estado_documento(tercero_id);
`)
}

func (m *InitDiplomas20260724) Down() {
	m.SQL(`
DROP TABLE IF EXISTS diplomas.historico_estado_documento;
DROP TABLE IF EXISTS diplomas.documento_digital;
DROP SCHEMA IF EXISTS diplomas;
`)
}
