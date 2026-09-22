package main

import "github.com/beego/beego/v2/client/orm/migration"

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
    tipo_documento_digital_id INTEGER NOT NULL,
    estado_documento_id INTEGER NOT NULL,
    codigo_estudiante BIGINT NOT NULL,
    facultad_id INTEGER,
    programa_academico_id INTEGER,
    periodo_id INTEGER,
    vigencia INTEGER,
    nombre_estudiante CHARACTER VARYING(250),
    tipo_documento_estudiante CHARACTER VARYING(20),
    numero_documento_estudiante CHARACTER VARYING(50),
    municipio_expedicion CHARACTER VARYING(150),
    titulo_otorgado CHARACTER VARYING(250),
    uuid_documento UUID,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_documento_digital PRIMARY KEY (id),
    CONSTRAINT ck_vigencia_documento_digital CHECK (vigencia IS NULL OR vigencia BETWEEN 2000 AND 2200),
    CONSTRAINT uq_uuid_documento_documento_digital UNIQUE (uuid_documento)
);

COMMENT ON TABLE diplomas.documento_digital IS 'Documento digital asociado al estudiante. Guarda ids externos del SGA, Terceros, parametros y uuid_documento entregado por firma_digital para almacenar y consultar el documento en S3.';
COMMENT ON COLUMN diplomas.documento_digital.tipo_documento_digital_id IS 'Referencia externa al parametro que identifica el tipo de documento digital emitido: acta de grado, diploma normal o diploma IDE. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.estado_documento_id IS 'Referencia externa al parametro que identifica el estado actual del documento. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.codigo_estudiante IS 'Codigo academico del estudiante aprobado a grado.';
COMMENT ON COLUMN diplomas.documento_digital.facultad_id IS 'Referencia externa a la facultad del estudiante en SGA/Oikos para asignar consecutivos del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.programa_academico_id IS 'Referencia externa al programa academico en SGA. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.periodo_id IS 'Referencia externa al periodo academico en SGA. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.documento_digital.nombre_estudiante IS 'Nombre completo del estudiante consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.tipo_documento_estudiante IS 'Tipo de documento de identidad del estudiante consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.numero_documento_estudiante IS 'Numero de documento del estudiante consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.municipio_expedicion IS 'Municipio de expedicion del documento del estudiante consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.titulo_otorgado IS 'Titulo academico otorgado consultado desde SGA 1 para preparar la generacion del diploma.';
COMMENT ON COLUMN diplomas.documento_digital.uuid_documento IS 'UUID del documento firmado usado por el sistema consumidor para almacenar y ubicar el archivo.';

CREATE TABLE IF NOT EXISTS diplomas.control_consecutivo_facultad_vigencia (
    id SERIAL NOT NULL,
    facultad_id INTEGER NOT NULL,
    vigencia INTEGER NOT NULL,
    ultimo_consecutivo_facultad INTEGER NOT NULL DEFAULT 0,
    ultimo_folio INTEGER NOT NULL DEFAULT 0,
    ultimo_acta INTEGER NOT NULL DEFAULT 0,
    libro_actual INTEGER NOT NULL DEFAULT 1,
    folios_por_libro INTEGER NOT NULL DEFAULT 500,
    actas_por_folio INTEGER NOT NULL DEFAULT 6,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_control_consecutivo_facultad_vigencia PRIMARY KEY (id),
    CONSTRAINT ck_vigencia_control_consecutivo_facultad_vigencia CHECK (vigencia BETWEEN 2000 AND 2200),
    CONSTRAINT ck_ultimo_folio_control_consecutivo_facultad_vigencia CHECK (ultimo_folio >= 0),
    CONSTRAINT ck_ultimo_acta_control_consecutivo_facultad_vigencia CHECK (ultimo_acta >= 0),
    CONSTRAINT ck_libro_actual_control_consecutivo_facultad_vigencia CHECK (libro_actual > 0),
    CONSTRAINT ck_folios_por_libro_control_consecutivo_facultad_vigencia CHECK (folios_por_libro > 0),
    CONSTRAINT ck_actas_por_folio_control_consecutivo_facultad_vigencia CHECK (actas_por_folio > 0),
    CONSTRAINT uq_facultad_id_vigencia_control_consecutivo_facultad_vigencia UNIQUE (facultad_id, vigencia)
);

COMMENT ON TABLE diplomas.control_consecutivo_facultad_vigencia IS 'Control atomico de consecutivo, libro, folio y acta por facultad/vigencia. Cada facultad maneja libros independientes y la vigencia permite reiniciar la configuracion cuando sea requerido.';
COMMENT ON COLUMN diplomas.control_consecutivo_facultad_vigencia.facultad_id IS 'Referencia externa a la facultad en SGA. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.control_consecutivo_facultad_vigencia.ultimo_folio IS 'Ultimo folio usado dentro del libro actual.';
COMMENT ON COLUMN diplomas.control_consecutivo_facultad_vigencia.ultimo_acta IS 'Ultima acta usada dentro del folio actual.';
COMMENT ON COLUMN diplomas.control_consecutivo_facultad_vigencia.libro_actual IS 'Libro vigente para la facultad y vigencia.';
COMMENT ON COLUMN diplomas.control_consecutivo_facultad_vigencia.folios_por_libro IS 'Cantidad parametrizable de folios permitidos por libro.';
COMMENT ON COLUMN diplomas.control_consecutivo_facultad_vigencia.actas_por_folio IS 'Cantidad parametrizable de actas permitidas por folio. Por defecto son 6.';

CREATE TABLE IF NOT EXISTS diplomas.diploma_digital (
    id SERIAL NOT NULL,
    documento_digital_id INTEGER NOT NULL,
    facultad_id INTEGER NOT NULL,
    vigencia INTEGER NOT NULL,
    fecha_grado DATE NOT NULL,
    nombre_estudiante CHARACTER VARYING(250) NOT NULL,
    tipo_documento_estudiante CHARACTER VARYING(20) NOT NULL,
    numero_documento_estudiante CHARACTER VARYING(50) NOT NULL,
    municipio_expedicion CHARACTER VARYING(150),
    titulo_otorgado CHARACTER VARYING(250) NOT NULL,
    consecutivo_diploma BIGINT NOT NULL,
    consecutivo_facultad INTEGER NOT NULL,
    folio INTEGER NOT NULL,
    acta INTEGER NOT NULL,
    libro INTEGER NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_diploma_digital PRIMARY KEY (id),
    CONSTRAINT ck_vigencia_diploma_digital CHECK (vigencia BETWEEN 2000 AND 2200),
    CONSTRAINT ck_nombre_estudiante_diploma_digital CHECK (btrim(nombre_estudiante) <> ''),
    CONSTRAINT ck_numero_documento_estudiante_diploma_digital CHECK (btrim(numero_documento_estudiante) <> ''),
    CONSTRAINT ck_titulo_otorgado_diploma_digital CHECK (btrim(titulo_otorgado) <> ''),
    CONSTRAINT fk_diploma_digital_documento_digital FOREIGN KEY (documento_digital_id) REFERENCES diplomas.documento_digital(id),
    CONSTRAINT uq_documento_digital_id_diploma_digital UNIQUE (documento_digital_id),
    CONSTRAINT uq_consecutivo_diploma_diploma_digital UNIQUE (consecutivo_diploma),
    CONSTRAINT uq_facultad_id_vigencia_consecutivo_facultad_diploma_digital UNIQUE (facultad_id, vigencia, consecutivo_facultad),
    CONSTRAINT uq_facultad_id_vigencia_libro_folio_acta_diploma_digital UNIQUE (facultad_id, vigencia, libro, folio, acta)
);

COMMENT ON TABLE diplomas.diploma_digital IS 'Datos propios del diploma creado al final del flujo por Rectoria. Incluye snapshot minimo de los datos impresos en el diploma para preservar el documento emitido aunque cambien las fuentes externas.';
COMMENT ON COLUMN diplomas.diploma_digital.facultad_id IS 'Referencia externa a la facultad en SGA. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.diploma_digital.nombre_estudiante IS 'Nombre completo del estudiante tal como queda impreso en el diploma emitido. Se copia desde la fuente academica/personas al momento de emision para trazabilidad historica.';
COMMENT ON COLUMN diplomas.diploma_digital.tipo_documento_estudiante IS 'Tipo de documento de identidad del estudiante tal como llega desde SGA 1 al momento de emision.';
COMMENT ON COLUMN diplomas.diploma_digital.numero_documento_estudiante IS 'Numero de documento del estudiante tal como llega desde SGA 1 y queda impreso en el diploma emitido.';
COMMENT ON COLUMN diplomas.diploma_digital.municipio_expedicion IS 'Municipio de expedicion del documento del estudiante tal como llega desde SGA 1.';
COMMENT ON COLUMN diplomas.diploma_digital.titulo_otorgado IS 'Titulo academico otorgado tal como llega desde SGA 1 y queda impreso en el diploma emitido.';
COMMENT ON COLUMN diplomas.diploma_digital.consecutivo_diploma IS 'Consecutivo global unico. Se asigna cuando Rectoria crea el diploma.';
COMMENT ON COLUMN diplomas.diploma_digital.consecutivo_facultad IS 'Consecutivo por facultad y vigencia.';
COMMENT ON COLUMN diplomas.diploma_digital.folio IS 'Folio por facultad y vigencia.';
COMMENT ON COLUMN diplomas.diploma_digital.acta IS 'Acta dentro del folio. Al completar actas_por_folio, el siguiente diploma inicia acta 1 del siguiente folio.';
COMMENT ON COLUMN diplomas.diploma_digital.libro IS 'Libro por facultad y vigencia. Al completar folios_por_libro, el siguiente diploma inicia un nuevo libro.';

CREATE TABLE IF NOT EXISTS diplomas.firma_firmante (
    id SERIAL NOT NULL,
    documento_identidad INTEGER NOT NULL,
    enlace_firma UUID NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_firma_firmante PRIMARY KEY (id),
    CONSTRAINT uq_enlace_firma_firma_firmante UNIQUE (enlace_firma)
);

COMMENT ON TABLE diplomas.firma_firmante IS 'Referencia a imagen de firma administrada por Nuxeo o servicio externo de firmas.';
COMMENT ON COLUMN diplomas.firma_firmante.documento_identidad IS 'Documento de identidad del firmante validado contra administrativa_amazon_api.';
COMMENT ON COLUMN diplomas.firma_firmante.enlace_firma IS 'UUID entregado por servicio externo para consumir la imagen de firma.';

CREATE UNIQUE INDEX IF NOT EXISTS idx_firma_firmante_documento_identidad_activo
ON diplomas.firma_firmante (documento_identidad)
WHERE activo IS TRUE;

CREATE TABLE IF NOT EXISTS diplomas.firma_documento (
    id SERIAL NOT NULL,
    documento_digital_id INTEGER NOT NULL,
    firma_firmante_id INTEGER,
    rol_firmante_id INTEGER NOT NULL,
    documento_identidad INTEGER NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT pk_firma_documento PRIMARY KEY (id),
    CONSTRAINT fk_firma_documento_documento_digital FOREIGN KEY (documento_digital_id) REFERENCES diplomas.documento_digital(id),
    CONSTRAINT fk_firma_documento_firma_firmante FOREIGN KEY (firma_firmante_id) REFERENCES diplomas.firma_firmante(id)
);

COMMENT ON TABLE diplomas.firma_documento IS 'Registro de firma aplicada a un documento digital. La firma digital, QR, hash y archivo final los administra firma_digital.';
COMMENT ON COLUMN diplomas.firma_documento.rol_firmante_id IS 'Referencia externa al parametro/rol del actor firmante. No se define FK por estar en otro servicio/esquema.';
COMMENT ON COLUMN diplomas.firma_documento.documento_identidad IS 'Documento de identidad del firmante validado contra administrativa_amazon_api.';

CREATE UNIQUE INDEX IF NOT EXISTS idx_firma_documento_documento_rol_activo
ON diplomas.firma_documento (documento_digital_id, rol_firmante_id)
WHERE activo IS TRUE;

CREATE TABLE IF NOT EXISTS diplomas.historico_estado_documento (
    id SERIAL NOT NULL,
    documento_digital_id INTEGER NOT NULL,
    estado_anterior_id INTEGER,
    estado_nuevo_id INTEGER NOT NULL,
    documento_identidad INTEGER,
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
COMMENT ON COLUMN diplomas.historico_estado_documento.documento_identidad IS 'Documento de identidad del actor que firma o ejecuta el cambio de estado.';
COMMENT ON COLUMN diplomas.historico_estado_documento.rol_actor_id IS 'Referencia externa al parametro/rol del actor que ejecuta el cambio. No se define FK por estar en otro servicio/esquema.';

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

CREATE INDEX IF NOT EXISTS idx_documento_digital_estado_documento_id
ON diplomas.documento_digital (estado_documento_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_codigo_estudiante
ON diplomas.documento_digital (codigo_estudiante);

CREATE INDEX IF NOT EXISTS idx_documento_digital_facultad_id
ON diplomas.documento_digital (facultad_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_programa_academico_id
ON diplomas.documento_digital (programa_academico_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_periodo_id
ON diplomas.documento_digital (periodo_id);

CREATE INDEX IF NOT EXISTS idx_documento_digital_vigencia
ON diplomas.documento_digital (vigencia);

CREATE INDEX IF NOT EXISTS idx_documento_digital_numero_documento_estudiante
ON diplomas.documento_digital (numero_documento_estudiante);

CREATE INDEX IF NOT EXISTS idx_control_consecutivo_facultad_vigencia_facultad_id_vigencia
ON diplomas.control_consecutivo_facultad_vigencia (facultad_id, vigencia);

CREATE INDEX IF NOT EXISTS idx_diploma_digital_facultad_id_vigencia
ON diplomas.diploma_digital (facultad_id, vigencia);

CREATE INDEX IF NOT EXISTS idx_firma_firmante_documento_identidad
ON diplomas.firma_firmante (documento_identidad);

CREATE INDEX IF NOT EXISTS idx_firma_firmante_enlace_firma
ON diplomas.firma_firmante (enlace_firma);

CREATE INDEX IF NOT EXISTS idx_firma_documento_documento_digital_id
ON diplomas.firma_documento (documento_digital_id);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_documento_fecha
ON diplomas.historico_estado_documento (documento_digital_id, fecha_creacion DESC);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_estado_nuevo_id
ON diplomas.historico_estado_documento (estado_nuevo_id);

CREATE INDEX IF NOT EXISTS idx_historico_estado_documento_documento_identidad
ON diplomas.historico_estado_documento (documento_identidad);
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
