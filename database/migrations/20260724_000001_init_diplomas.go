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

COMMENT ON SCHEMA diplomas IS 'Esquema funcional para la gestion de diplomas digitales. Guarda referencias a servicios externos, trazabilidad, firmas aplicadas y consecutivos propios del diploma.';

CREATE SEQUENCE IF NOT EXISTS diplomas.seq_consecutivo_diploma START WITH 1 INCREMENT BY 1;

CREATE TABLE IF NOT EXISTS diplomas.documento_digital (
    id SERIAL NOT NULL,
    tipo_documento_id INTEGER NOT NULL,
    estado_documento_id INTEGER NOT NULL,
    tercero_id_estudiante INTEGER NOT NULL,
    programa_academico_id INTEGER,
    periodo_id INTEGER,
    vigencia INTEGER,
    uuid_verificacion UUID,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_documento_digital PRIMARY KEY (id),
    CONSTRAINT ck_vigencia_documento_digital CHECK (vigencia IS NULL OR vigencia BETWEEN 2000 AND 2200),
    CONSTRAINT uq_uuid_verificacion_documento_digital UNIQUE (uuid_verificacion)
);

COMMENT ON TABLE diplomas.documento_digital IS 'Documento digital asociado al estudiante. Guarda ids externos del SGA, Terceros, parametros y uuid de verificacion entregado por firma_digital para consultar el documento en S3.';
COMMENT ON COLUMN diplomas.documento_digital.tipo_documento_id IS 'Referencia externa al parametro que identifica el tipo de documento. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.estado_documento_id IS 'Referencia externa al parametro que identifica el estado actual del documento. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.tercero_id_estudiante IS 'Referencia externa al id del estudiante en Terceros. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.programa_academico_id IS 'Referencia externa al programa academico en SGA. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.periodo_id IS 'Referencia externa al periodo academico en SGA. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.uuid_verificacion IS 'UUID entregado por firma_digital para consultar el documento generado en S3.';

CREATE TABLE IF NOT EXISTS diplomas.control_consecutivo_facultad_vigencia (
    id SERIAL NOT NULL,
    facultad_id INTEGER NOT NULL,
    vigencia INTEGER NOT NULL,
    ultimo_consecutivo_facultad INTEGER NOT NULL DEFAULT 0,
    ultimo_folio INTEGER NOT NULL DEFAULT 0,
    libro_actual INTEGER NOT NULL DEFAULT 1,
    folios_por_libro INTEGER NOT NULL DEFAULT 500,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_control_consecutivo_facultad_vigencia PRIMARY KEY (id),
    CONSTRAINT ck_vigencia_control_consecutivo_facultad_vigencia CHECK (vigencia BETWEEN 2000 AND 2200),
    CONSTRAINT ck_folios_por_libro_control_consecutivo_facultad_vigencia CHECK (folios_por_libro > 0),
    CONSTRAINT uq_facultad_id_vigencia_control_consecutivo_facultad_vigencia UNIQUE (facultad_id, vigencia)
);

COMMENT ON TABLE diplomas.control_consecutivo_facultad_vigencia IS 'Control atomico de consecutivo de facultad, folio y libro por facultad/vigencia. facultad_id referencia la facultad del SGA; no se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.control_consecutivo_facultad_vigencia.facultad_id IS 'Referencia externa a la facultad en SGA. No se define FK por estar en otro servicio/esquema.';

CREATE TABLE IF NOT EXISTS diplomas.diploma_digital (
    id SERIAL NOT NULL,
    documento_digital_id INTEGER NOT NULL,
    facultad_id INTEGER NOT NULL,
    vigencia INTEGER NOT NULL,
    fecha_grado DATE NOT NULL,
    consecutivo_diploma BIGINT NOT NULL,
    consecutivo_facultad INTEGER NOT NULL,
    folio INTEGER NOT NULL,
    libro INTEGER NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_diploma_digital PRIMARY KEY (id),
    CONSTRAINT ck_vigencia_diploma_digital CHECK (vigencia BETWEEN 2000 AND 2200),
    CONSTRAINT fk_diploma_digital_documento_digital FOREIGN KEY (documento_digital_id) REFERENCES diplomas.documento_digital(id),
    CONSTRAINT uq_documento_digital_id_diploma_digital UNIQUE (documento_digital_id),
    CONSTRAINT uq_consecutivo_diploma_diploma_digital UNIQUE (consecutivo_diploma),
    CONSTRAINT uq_facultad_id_vigencia_consecutivo_facultad_diploma_digital UNIQUE (facultad_id, vigencia, consecutivo_facultad),
    CONSTRAINT uq_facultad_id_vigencia_libro_folio_diploma_digital UNIQUE (facultad_id, vigencia, libro, folio)
);

COMMENT ON TABLE diplomas.diploma_digital IS 'Datos propios del diploma creado al final del flujo por Rectoria.';
COMMENT ON COLUMN diplomas.diploma_digital.facultad_id IS 'Referencia externa a la facultad en SGA. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.diploma_digital.consecutivo_diploma IS 'Consecutivo global unico. Se asigna cuando Rectoria crea el diploma.';
COMMENT ON COLUMN diplomas.diploma_digital.consecutivo_facultad IS 'Consecutivo por facultad y vigencia.';
COMMENT ON COLUMN diplomas.diploma_digital.folio IS 'Folio por facultad y vigencia.';
COMMENT ON COLUMN diplomas.diploma_digital.libro IS 'Libro por facultad y vigencia.';

CREATE TABLE IF NOT EXISTS diplomas.firma_firmante (
    id SERIAL NOT NULL,
    tercero_id_firmante INTEGER NOT NULL,
    enlace_firma UUID NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_firma_firmante PRIMARY KEY (id),
    CONSTRAINT uq_enlace_firma_firma_firmante UNIQUE (enlace_firma)
);

COMMENT ON TABLE diplomas.firma_firmante IS 'Referencia a imagen de firma administrada por Nuxeo o servicio externo de firmas.';
COMMENT ON COLUMN diplomas.firma_firmante.tercero_id_firmante IS 'Referencia externa al id del firmante en Terceros. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.firma_firmante.enlace_firma IS 'UUID entregado por servicio externo para consumir la imagen de firma.';

CREATE UNIQUE INDEX IF NOT EXISTS idx_firma_firmante_tercero_id_firmante_activo
ON diplomas.firma_firmante (tercero_id_firmante)
WHERE activo IS TRUE;

CREATE TABLE IF NOT EXISTS diplomas.firma_documento (
    id SERIAL NOT NULL,
    documento_digital_id INTEGER NOT NULL,
    firma_firmante_id INTEGER,
    rol_firmante_id INTEGER NOT NULL,
    tercero_id_firmante INTEGER NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_firma_documento PRIMARY KEY (id),
    CONSTRAINT fk_firma_documento_documento_digital FOREIGN KEY (documento_digital_id) REFERENCES diplomas.documento_digital(id),
    CONSTRAINT fk_firma_documento_firma_firmante FOREIGN KEY (firma_firmante_id) REFERENCES diplomas.firma_firmante(id)
);

COMMENT ON TABLE diplomas.firma_documento IS 'Registro de firma aplicada a un documento digital. La firma digital, QR, hash y archivo final los administra firma_digital.';
COMMENT ON COLUMN diplomas.firma_documento.rol_firmante_id IS 'Referencia externa al parametro/rol del actor firmante. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.firma_documento.tercero_id_firmante IS 'Referencia externa al id del firmante en Terceros. No se define FK por estar en otro servicio/esquema.';

CREATE UNIQUE INDEX IF NOT EXISTS idx_firma_documento_documento_rol_activo
ON diplomas.firma_documento (documento_digital_id, rol_firmante_id)
WHERE activo IS TRUE;

CREATE TABLE IF NOT EXISTS diplomas.historico_estado_documento (
    id SERIAL NOT NULL,
    documento_digital_id INTEGER NOT NULL,
    estado_anterior_id INTEGER,
    estado_nuevo_id INTEGER NOT NULL,
    tercero_id_firmante INTEGER,
    rol_actor_id INTEGER,
    firma_documento_id INTEGER,
    observacion CHARACTER VARYING(500),
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_historico_estado_documento PRIMARY KEY (id),
    CONSTRAINT fk_historico_estado_documento_documento_digital FOREIGN KEY (documento_digital_id) REFERENCES diplomas.documento_digital(id),
    CONSTRAINT fk_historico_estado_documento_firma_documento FOREIGN KEY (firma_documento_id) REFERENCES diplomas.firma_documento(id)
);

COMMENT ON TABLE diplomas.historico_estado_documento IS 'Auditoria de cambios de estado del documento digital.';
COMMENT ON COLUMN diplomas.historico_estado_documento.estado_anterior_id IS 'Referencia externa al parametro de estado anterior. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.historico_estado_documento.estado_nuevo_id IS 'Referencia externa al parametro de estado nuevo. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.historico_estado_documento.tercero_id_firmante IS 'Referencia externa al id del actor que firma o ejecuta el cambio de estado en Terceros. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.historico_estado_documento.rol_actor_id IS 'Referencia externa al parametro/rol del actor que ejecuta el cambio. No se define FK por estar en otro servicio/esquema.';

CREATE INDEX IF NOT EXISTS idx_documento_digital_tipo_documento_id
ON diplomas.documento_digital (tipo_documento_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_estado_documento_id
ON diplomas.documento_digital (estado_documento_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_tercero_id_estudiante
ON diplomas.documento_digital (tercero_id_estudiante);

CREATE INDEX IF NOT EXISTS idx_documento_digital_programa_academico_id
ON diplomas.documento_digital (programa_academico_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_periodo_id
ON diplomas.documento_digital (periodo_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_vigencia
ON diplomas.documento_digital (vigencia);

CREATE INDEX IF NOT EXISTS idx_control_consecutivo_facultad_vigencia_facultad_id_vigencia
ON diplomas.control_consecutivo_facultad_vigencia (facultad_id, vigencia);

CREATE INDEX IF NOT EXISTS idx_diploma_digital_facultad_id_vigencia
ON diplomas.diploma_digital (facultad_id, vigencia);

CREATE INDEX IF NOT EXISTS idx_firma_firmante_tercero_id_firmante
ON diplomas.firma_firmante (tercero_id_firmante);

CREATE INDEX IF NOT EXISTS idx_firma_firmante_enlace_firma
ON diplomas.firma_firmante (enlace_firma);

CREATE INDEX IF NOT EXISTS idx_firma_documento_documento_digital_id
ON diplomas.firma_documento (documento_digital_id);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_documento_fecha
ON diplomas.historico_estado_documento (documento_digital_id, fecha_creacion DESC);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_estado_nuevo_id
ON diplomas.historico_estado_documento (estado_nuevo_id);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_tercero_id_firmante
ON diplomas.historico_estado_documento (tercero_id_firmante);
`)
}

func (m *InitDiplomas20260724) Down() {
	m.SQL(`
DROP TABLE IF EXISTS diplomas.historico_estado_documento;
DROP TABLE IF EXISTS diplomas.firma_documento;
DROP TABLE IF EXISTS diplomas.firma_firmante;
DROP TABLE IF EXISTS diplomas.diploma_digital;
DROP TABLE IF EXISTS diplomas.control_consecutivo_facultad_vigencia;
DROP TABLE IF EXISTS diplomas.documento_digital;
DROP SEQUENCE IF EXISTS diplomas.seq_consecutivo_diploma;
DROP SCHEMA IF EXISTS diplomas;
`)
}
