package main

import "github.com/beego/beego/v2/client/orm/migration"

type RegistroGradoModel20260928 struct {
	migration.Migration
}

func init() {
	m := &RegistroGradoModel20260928{}
	m.Created = "20260928_000001"
	migration.Register("RegistroGradoModel20260928", m)
}

func (m *RegistroGradoModel20260928) Up() {
	m.SQL(`

ALTER TABLE diplomas.diploma_digital
ALTER COLUMN consecutivo_diploma SET DEFAULT nextval('diplomas.seq_consecutivo_diploma'),
ALTER COLUMN consecutivo_facultad DROP NOT NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas' AND table_name = 'diploma_digital' AND column_name = 'folio'
    ) THEN
        ALTER TABLE diplomas.diploma_digital ALTER COLUMN folio DROP NOT NULL;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas' AND table_name = 'diploma_digital' AND column_name = 'acta'
    ) THEN
        ALTER TABLE diplomas.diploma_digital ALTER COLUMN acta DROP NOT NULL;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'diplomas' AND table_name = 'diploma_digital' AND column_name = 'libro'
    ) THEN
        ALTER TABLE diplomas.diploma_digital ALTER COLUMN libro DROP NOT NULL;
    END IF;
END;
$$;

-- 4. PARAMETRIZACION REGISTRAL POR FACULTAD
-- ============================================================
--
-- Ejemplo:
--
-- Facultad 14 / vigencia 2027
--
-- libro inicial                = 1
-- folio inicial                = 1
-- registro inicial             = 1
-- consecutivo facultad inicial = 1
-- folios por libro             = 500
-- registros por folio          = 6
--
-- ============================================================

CREATE TABLE IF NOT EXISTS diplomas.parametrizacion_registro_facultad (
    id SERIAL NOT NULL,

    facultad_id INTEGER NOT NULL,
    vigencia INTEGER NOT NULL,

    libro_inicial INTEGER NOT NULL DEFAULT 1,
    folio_inicial INTEGER NOT NULL DEFAULT 1,
    registro_inicial INTEGER NOT NULL DEFAULT 1,

    consecutivo_facultad_inicial INTEGER NOT NULL DEFAULT 1,

    folios_por_libro INTEGER NOT NULL DEFAULT 500,
    registros_por_folio INTEGER NOT NULL DEFAULT 6,

    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_parametrizacion_registro_facultad
        PRIMARY KEY (id),

    CONSTRAINT uq_parametrizacion_registro_facultad
        UNIQUE (
            facultad_id,
            vigencia
        ),

    CONSTRAINT ck_parametrizacion_vigencia
        CHECK (
            vigencia BETWEEN 2000 AND 2200
        ),

    CONSTRAINT ck_parametrizacion_libro_inicial
        CHECK (libro_inicial > 0),

    CONSTRAINT ck_parametrizacion_folio_inicial
        CHECK (folio_inicial > 0),

    CONSTRAINT ck_parametrizacion_registro_inicial
        CHECK (registro_inicial > 0),

    CONSTRAINT ck_parametrizacion_consecutivo_inicial
        CHECK (consecutivo_facultad_inicial > 0),

    CONSTRAINT ck_parametrizacion_folios_por_libro
        CHECK (folios_por_libro > 0),

    CONSTRAINT ck_parametrizacion_registros_por_folio
        CHECK (registros_por_folio > 0)
);


COMMENT ON TABLE diplomas.parametrizacion_registro_facultad IS
'Configuracion registral de una facultad para una vigencia. Define desde que libro, folio, registro y consecutivo inicia y la capacidad de libros y folios.';

COMMENT ON COLUMN diplomas.parametrizacion_registro_facultad.facultad_id IS
'Referencia externa a la facultad en SGA/Oikos.';

COMMENT ON COLUMN diplomas.parametrizacion_registro_facultad.consecutivo_facultad_inicial IS
'Primer consecutivo institucional de la facultad que debe asignarse para la vigencia.';

COMMENT ON COLUMN diplomas.parametrizacion_registro_facultad.folios_por_libro IS
'Cantidad maxima de folios que contiene un libro.';

COMMENT ON COLUMN diplomas.parametrizacion_registro_facultad.registros_por_folio IS
'Cantidad predeterminada de registros permitidos en cada folio.';


-- ============================================================
-- 5. LIBROS
-- ============================================================

CREATE TABLE IF NOT EXISTS diplomas.libro_registro (
    id SERIAL NOT NULL,

    parametrizacion_registro_facultad_id INTEGER NOT NULL,

    numero_libro INTEGER NOT NULL,

    cantidad_folios INTEGER NOT NULL,

    fecha_apertura DATE NOT NULL DEFAULT CURRENT_DATE,
    fecha_cierre DATE,

    cerrado BOOLEAN NOT NULL DEFAULT FALSE,
    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_libro_registro
        PRIMARY KEY (id),

    CONSTRAINT fk_libro_registro_parametrizacion
        FOREIGN KEY (
            parametrizacion_registro_facultad_id
        )
        REFERENCES diplomas.parametrizacion_registro_facultad(id),

    CONSTRAINT uq_libro_registro_numero
        UNIQUE (
            parametrizacion_registro_facultad_id,
            numero_libro
        ),

    CONSTRAINT ck_libro_registro_numero
        CHECK (numero_libro > 0),

    CONSTRAINT ck_libro_registro_cantidad_folios
        CHECK (cantidad_folios > 0),

    CONSTRAINT ck_libro_registro_fechas
        CHECK (
            fecha_cierre IS NULL
            OR fecha_cierre >= fecha_apertura
        )
);


COMMENT ON TABLE diplomas.libro_registro IS
'Libro de registro de grados perteneciente a una facultad y vigencia.';

COMMENT ON COLUMN diplomas.libro_registro.cantidad_folios IS
'Cantidad maxima de folios disponibles en este libro.';


-- ============================================================
-- 6. FOLIOS
-- ============================================================

CREATE TABLE IF NOT EXISTS diplomas.folio_registro (
    id SERIAL NOT NULL,

    libro_registro_id INTEGER NOT NULL,

    numero_folio INTEGER NOT NULL,

    cantidad_registros INTEGER NOT NULL,

    cerrado BOOLEAN NOT NULL DEFAULT FALSE,
    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_folio_registro
        PRIMARY KEY (id),

    CONSTRAINT fk_folio_registro_libro
        FOREIGN KEY (libro_registro_id)
        REFERENCES diplomas.libro_registro(id),

    CONSTRAINT uq_folio_registro_numero
        UNIQUE (
            libro_registro_id,
            numero_folio
        ),

    CONSTRAINT ck_folio_registro_numero
        CHECK (numero_folio > 0),

    CONSTRAINT ck_folio_registro_cantidad
        CHECK (cantidad_registros > 0)
);


COMMENT ON TABLE diplomas.folio_registro IS
'Folio u hoja perteneciente a un libro de registro.';

COMMENT ON COLUMN diplomas.folio_registro.cantidad_registros IS
'Cantidad maxima de registros admitidos en este folio. Permite modificar individualmente la capacidad de un folio.';


-- ============================================================
-- 7. CONTROL ATOMICO POR FACULTAD/VIGENCIA
-- ============================================================
--
-- Esta tabla es especialmente importante.
--
-- consecutivo facultad:
-- 1,2,3,4,5,6,7,8...
--
-- numero registro:
-- 1,2,3,4,5,6 | 1,2,3...
--
-- ============================================================

CREATE TABLE IF NOT EXISTS diplomas.control_registro_facultad (
    id SERIAL NOT NULL,

    parametrizacion_registro_facultad_id INTEGER NOT NULL,

    libro_actual_id INTEGER,
    folio_actual_id INTEGER,

    ultimo_consecutivo_facultad INTEGER NOT NULL DEFAULT 0,
    ultimo_numero_registro INTEGER NOT NULL DEFAULT 0,

    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_control_registro_facultad
        PRIMARY KEY (id),

    CONSTRAINT fk_control_registro_parametrizacion
        FOREIGN KEY (
            parametrizacion_registro_facultad_id
        )
        REFERENCES diplomas.parametrizacion_registro_facultad(id),

    CONSTRAINT fk_control_registro_libro
        FOREIGN KEY (libro_actual_id)
        REFERENCES diplomas.libro_registro(id),

    CONSTRAINT fk_control_registro_folio
        FOREIGN KEY (folio_actual_id)
        REFERENCES diplomas.folio_registro(id),

    CONSTRAINT uq_control_registro_parametrizacion
        UNIQUE (
            parametrizacion_registro_facultad_id
        ),

    CONSTRAINT ck_control_registro_ultimo_consecutivo
        CHECK (
            ultimo_consecutivo_facultad >= 0
        ),

    CONSTRAINT ck_control_registro_ultimo_numero
        CHECK (
            ultimo_numero_registro >= 0
        )
);


COMMENT ON TABLE diplomas.control_registro_facultad IS
'Control transaccional para asignar consecutivo de facultad y posicion Libro/Folio/Registro sin duplicados.';

COMMENT ON COLUMN diplomas.control_registro_facultad.ultimo_consecutivo_facultad IS
'Ultimo consecutivo continuo asignado a la facultad durante la vigencia.';

COMMENT ON COLUMN diplomas.control_registro_facultad.ultimo_numero_registro IS
'Ultimo numero de registro utilizado dentro del folio actual.';


-- ============================================================
-- 8. DIPLOMA DIGITAL
-- ============================================================

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

    consecutivo_diploma BIGINT NOT NULL
        DEFAULT nextval('diplomas.seq_consecutivo_diploma'),

    /*
     * Se asigna atomicamente mediante asignar_registro_grado().
     *
     * Puede estar NULL mientras el diploma fue creado pero aun
     * no ha sido asentado en el libro de la facultad.
     */
    consecutivo_facultad INTEGER,

    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_diploma_digital
        PRIMARY KEY (id),

    CONSTRAINT ck_vigencia_diploma_digital
        CHECK (
            vigencia BETWEEN 2000 AND 2200
        ),

    CONSTRAINT ck_nombre_estudiante_diploma_digital
        CHECK (
            btrim(nombre_estudiante) <> ''
        ),

    CONSTRAINT ck_numero_documento_estudiante_diploma_digital
        CHECK (
            btrim(numero_documento_estudiante) <> ''
        ),

    CONSTRAINT ck_titulo_otorgado_diploma_digital
        CHECK (
            btrim(titulo_otorgado) <> ''
        ),

    CONSTRAINT ck_consecutivo_facultad_diploma_digital
        CHECK (
            consecutivo_facultad IS NULL
            OR consecutivo_facultad > 0
        ),

    CONSTRAINT fk_diploma_digital_documento_digital
        FOREIGN KEY (documento_digital_id)
        REFERENCES diplomas.documento_digital(id),

    CONSTRAINT uq_documento_digital_diploma_digital
        UNIQUE (documento_digital_id),

    CONSTRAINT uq_consecutivo_diploma_diploma_digital
        UNIQUE (consecutivo_diploma),

    CONSTRAINT uq_facultad_vigencia_consecutivo_diploma_digital
        UNIQUE (
            facultad_id,
            vigencia,
            consecutivo_facultad
        )
);


COMMENT ON TABLE diplomas.diploma_digital IS
'Datos propios del diploma emitido. Conserva snapshot de la informacion impresa y los consecutivos global y por facultad.';

COMMENT ON COLUMN diplomas.diploma_digital.consecutivo_diploma IS
'Consecutivo global unico del diploma.';

COMMENT ON COLUMN diplomas.diploma_digital.consecutivo_facultad IS
'Consecutivo continuo por facultad/vigencia. Es independiente del numero de registro dentro del folio.';


-- ============================================================
-- 9. REGISTRO DE GRADO
-- ============================================================
--
-- Aquí queda el segundo número:
-- numero_registro.
--
-- Ejemplo:
--
-- consecutivo_facultad = 7
-- libro                 = 1
-- folio                 = 2
-- numero_registro       = 1
--
-- ============================================================

CREATE TABLE IF NOT EXISTS diplomas.registro_grado (
    id SERIAL NOT NULL,

    diploma_digital_id INTEGER NOT NULL,
    folio_registro_id INTEGER NOT NULL,

    numero_registro INTEGER NOT NULL,

    fecha_registro TIMESTAMP NOT NULL DEFAULT now(),

    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_registro_grado
        PRIMARY KEY (id),

    CONSTRAINT fk_registro_grado_diploma
        FOREIGN KEY (diploma_digital_id)
        REFERENCES diplomas.diploma_digital(id),

    CONSTRAINT fk_registro_grado_folio
        FOREIGN KEY (folio_registro_id)
        REFERENCES diplomas.folio_registro(id),

    CONSTRAINT uq_registro_grado_diploma
        UNIQUE (diploma_digital_id),

    CONSTRAINT uq_registro_grado_folio_numero
        UNIQUE (
            folio_registro_id,
            numero_registro
        ),

    CONSTRAINT ck_registro_grado_numero
        CHECK (
            numero_registro > 0
        )
);


COMMENT ON TABLE diplomas.registro_grado IS
'Asiento registral del diploma dentro de un libro y folio.';

COMMENT ON COLUMN diplomas.registro_grado.numero_registro IS
'Numero del registro dentro del folio. Se reinicia en 1 cuando se inicia un nuevo folio.';


-- ============================================================
-- 10. CEREMONIA DE GRADO
-- ============================================================

CREATE TABLE IF NOT EXISTS diplomas.ceremonia_grado (
    id SERIAL NOT NULL,

    fecha_ceremonia DATE NOT NULL,
    hora_ceremonia TIME,

    lugar CHARACTER VARYING(300) NOT NULL,

    cantidad_estudiantes_programada INTEGER NOT NULL,

    observacion CHARACTER VARYING(500),

    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_ceremonia_grado
        PRIMARY KEY (id),

    CONSTRAINT ck_ceremonia_grado_lugar
        CHECK (
            btrim(lugar) <> ''
        ),

    CONSTRAINT ck_ceremonia_grado_cantidad
        CHECK (
            cantidad_estudiantes_programada > 0
        )
);


COMMENT ON TABLE diplomas.ceremonia_grado IS
'Ceremonia de grado definida por fecha, hora y lugar. Una ceremonia puede agrupar varias facultades y programas curriculares.';


-- ============================================================
-- 11. CEREMONIA - FACULTAD - PROGRAMA
-- ============================================================
--
-- Permite:
--
-- misma facultad -> varias fechas
-- misma fecha    -> varias facultades
-- facultad       -> varios programas
--
-- ============================================================

CREATE TABLE IF NOT EXISTS diplomas.ceremonia_facultad_programa (
    id SERIAL NOT NULL,

    ceremonia_grado_id INTEGER NOT NULL,

    facultad_id INTEGER NOT NULL,
    programa_academico_id INTEGER NOT NULL,

    cantidad_estudiantes_programada INTEGER,

    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_ceremonia_facultad_programa
        PRIMARY KEY (id),

    CONSTRAINT fk_ceremonia_facultad_programa_ceremonia
        FOREIGN KEY (ceremonia_grado_id)
        REFERENCES diplomas.ceremonia_grado(id),

    CONSTRAINT uq_ceremonia_facultad_programa
        UNIQUE (
            ceremonia_grado_id,
            facultad_id,
            programa_academico_id
        ),

    CONSTRAINT ck_ceremonia_facultad_programa_cantidad
        CHECK (
            cantidad_estudiantes_programada IS NULL
            OR cantidad_estudiantes_programada >= 0
        )
);


COMMENT ON TABLE diplomas.ceremonia_facultad_programa IS
'Relaciona una ceremonia de grado con las facultades y programas curriculares participantes.';

COMMENT ON COLUMN diplomas.ceremonia_facultad_programa.facultad_id IS
'Referencia externa a la facultad en SGA/Oikos.';

COMMENT ON COLUMN diplomas.ceremonia_facultad_programa.programa_academico_id IS
'Referencia externa al programa curricular en SGA.';


-- ============================================================
-- 12. ESTUDIANTES ASIGNADOS A CEREMONIA
-- ============================================================

CREATE TABLE IF NOT EXISTS diplomas.ceremonia_estudiante (
    id SERIAL NOT NULL,

    ceremonia_facultad_programa_id INTEGER NOT NULL,

    documento_digital_id INTEGER NOT NULL,

    activo BOOLEAN NOT NULL DEFAULT TRUE,

    fecha_creacion TIMESTAMP NOT NULL DEFAULT now(),
    fecha_modificacion TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT pk_ceremonia_estudiante
        PRIMARY KEY (id),

    CONSTRAINT fk_ceremonia_estudiante_programacion
        FOREIGN KEY (
            ceremonia_facultad_programa_id
        )
        REFERENCES diplomas.ceremonia_facultad_programa(id),

    CONSTRAINT fk_ceremonia_estudiante_documento
        FOREIGN KEY (documento_digital_id)
        REFERENCES diplomas.documento_digital(id),

    CONSTRAINT uq_ceremonia_estudiante
        UNIQUE (
            ceremonia_facultad_programa_id,
            documento_digital_id
        )
);


COMMENT ON TABLE diplomas.ceremonia_estudiante IS
'Relaciona cada estudiante/documento digital con la ceremonia, facultad y programa donde sera graduado.';


-- ============================================================

-- 21. INICIALIZAR PARAMETRIZACION
-- ============================================================
--
-- Esta funcion crea:
--
-- libro inicial
-- folio inicial
-- control de consecutivos
--
-- ============================================================

CREATE OR REPLACE FUNCTION
diplomas.inicializar_parametrizacion_registro(
    p_parametrizacion_id INTEGER
)
RETURNS VOID AS $$
DECLARE

    v_param
        diplomas.parametrizacion_registro_facultad%ROWTYPE;

    v_libro_id INTEGER;

    v_folio_id INTEGER;

BEGIN

    SELECT *
    INTO v_param
    FROM diplomas.parametrizacion_registro_facultad
    WHERE id = p_parametrizacion_id
      AND activo IS TRUE;


    IF NOT FOUND THEN

        RAISE EXCEPTION
            'No existe parametrizacion activa con id %',
            p_parametrizacion_id;

    END IF;


    -- --------------------------------------------------------
    -- Crear libro inicial
    -- --------------------------------------------------------

    INSERT INTO diplomas.libro_registro (
        parametrizacion_registro_facultad_id,
        numero_libro,
        cantidad_folios
    )
    VALUES (
        v_param.id,
        v_param.libro_inicial,
        v_param.folios_por_libro
    )

    ON CONFLICT (
        parametrizacion_registro_facultad_id,
        numero_libro
    )

    DO UPDATE
       SET activo = TRUE

    RETURNING id
    INTO v_libro_id;


    -- --------------------------------------------------------
    -- Crear folio inicial
    -- --------------------------------------------------------

    INSERT INTO diplomas.folio_registro (
        libro_registro_id,
        numero_folio,
        cantidad_registros
    )
    VALUES (
        v_libro_id,
        v_param.folio_inicial,
        v_param.registros_por_folio
    )

    ON CONFLICT (
        libro_registro_id,
        numero_folio
    )

    DO UPDATE
       SET activo = TRUE

    RETURNING id
    INTO v_folio_id;


    -- --------------------------------------------------------
    -- Crear control
    -- --------------------------------------------------------

    INSERT INTO diplomas.control_registro_facultad (
        parametrizacion_registro_facultad_id,
        libro_actual_id,
        folio_actual_id,
        ultimo_consecutivo_facultad,
        ultimo_numero_registro
    )
    VALUES (
        v_param.id,
        v_libro_id,
        v_folio_id,
        v_param.consecutivo_facultad_inicial - 1,
        v_param.registro_inicial - 1
    )

    ON CONFLICT (
        parametrizacion_registro_facultad_id
    )
    DO NOTHING;

END;
$$ LANGUAGE plpgsql;


COMMENT ON FUNCTION
diplomas.inicializar_parametrizacion_registro(INTEGER) IS
'Inicializa libro, folio y control de consecutivos para una facultad/vigencia.';


-- ============================================================
-- 22. ASIGNAR REGISTRO DE GRADO
-- ============================================================
--
-- REGLAS:
--
-- consecutivo_facultad siempre aumenta.
--
-- registro:
-- 1 2 3 4 5 6
--
-- siguiente:
-- folio + 1
-- registro = 1
--
-- cuando se terminan los folios:
-- libro + 1
-- folio = 1
-- registro = 1
--
-- ============================================================

CREATE OR REPLACE FUNCTION
diplomas.asignar_registro_grado(
    p_diploma_digital_id INTEGER
)
RETURNS TABLE (

    registro_grado_id INTEGER,

    consecutivo_facultad_asignado INTEGER,

    numero_libro INTEGER,

    numero_folio INTEGER,

    numero_registro INTEGER

) AS $$

DECLARE

    v_diploma
        diplomas.diploma_digital%ROWTYPE;

    v_param
        diplomas.parametrizacion_registro_facultad%ROWTYPE;

    v_control
        diplomas.control_registro_facultad%ROWTYPE;

    v_libro
        diplomas.libro_registro%ROWTYPE;

    v_folio
        diplomas.folio_registro%ROWTYPE;

    v_nuevo_consecutivo INTEGER;

    v_nuevo_registro INTEGER;

    v_nuevo_libro INTEGER;

    v_nuevo_folio INTEGER;

    v_registro_grado_id INTEGER;

BEGIN

    -- --------------------------------------------------------
    -- Obtener diploma
    -- --------------------------------------------------------

    SELECT *
    INTO v_diploma
    FROM diplomas.diploma_digital
    WHERE id = p_diploma_digital_id
      AND activo IS TRUE;


    IF NOT FOUND THEN

        RAISE EXCEPTION
            'No existe diploma digital activo con id %',
            p_diploma_digital_id;

    END IF;


    -- --------------------------------------------------------
    -- Validar que no este registrado
    -- --------------------------------------------------------

    IF EXISTS (

        SELECT 1
        FROM diplomas.registro_grado
        WHERE diploma_digital_id = p_diploma_digital_id
          AND activo IS TRUE

    ) THEN

        RAISE EXCEPTION
            'El diploma % ya tiene registro de grado asignado',
            p_diploma_digital_id;

    END IF;


    -- --------------------------------------------------------
    -- Obtener parametrizacion
    -- --------------------------------------------------------

    SELECT *
    INTO v_param
    FROM diplomas.parametrizacion_registro_facultad
    WHERE facultad_id = v_diploma.facultad_id
      AND vigencia = v_diploma.vigencia
      AND activo IS TRUE;


    IF NOT FOUND THEN

        RAISE EXCEPTION
            'No existe parametrizacion activa para facultad % y vigencia %',
            v_diploma.facultad_id,
            v_diploma.vigencia;

    END IF;


    -- --------------------------------------------------------
    -- Inicializar si es la primera vez
    -- --------------------------------------------------------

    PERFORM
        diplomas.inicializar_parametrizacion_registro(
            v_param.id
        );


    -- --------------------------------------------------------
    -- BLOQUEO DEL CONTROL
    --
    -- Esto impide que dos solicitudes simultaneas reciban
    -- el mismo consecutivo.
    -- --------------------------------------------------------

    SELECT *
    INTO v_control
    FROM diplomas.control_registro_facultad
    WHERE parametrizacion_registro_facultad_id = v_param.id
      AND activo IS TRUE
    FOR UPDATE;


    -- --------------------------------------------------------
    -- Obtener libro actual
    -- --------------------------------------------------------

    SELECT *
    INTO v_libro
    FROM diplomas.libro_registro
    WHERE id = v_control.libro_actual_id
    FOR UPDATE;


    -- --------------------------------------------------------
    -- Obtener folio actual
    -- --------------------------------------------------------

    SELECT *
    INTO v_folio
    FROM diplomas.folio_registro
    WHERE id = v_control.folio_actual_id
    FOR UPDATE;


    -- --------------------------------------------------------
    -- CALCULAR CONSECUTIVO FACULTAD
    -- --------------------------------------------------------

    v_nuevo_consecutivo :=
        v_control.ultimo_consecutivo_facultad + 1;


    -- --------------------------------------------------------
    -- CALCULAR REGISTRO DEL FOLIO
    -- --------------------------------------------------------

    v_nuevo_registro :=
        v_control.ultimo_numero_registro + 1;


    -- ========================================================
    -- FOLIO LLENO
    -- ========================================================

    IF v_nuevo_registro > v_folio.cantidad_registros THEN


        UPDATE diplomas.folio_registro
        SET cerrado = TRUE
        WHERE id = v_folio.id;


        v_nuevo_folio :=
            v_folio.numero_folio + 1;


        v_nuevo_registro := 1;


        -- ====================================================
        -- LIBRO LLENO
        -- ====================================================

        IF v_nuevo_folio > v_libro.cantidad_folios THEN


            UPDATE diplomas.libro_registro
            SET
                cerrado = TRUE,
                fecha_cierre = CURRENT_DATE
            WHERE id = v_libro.id;


            v_nuevo_libro :=
                v_libro.numero_libro + 1;


            v_nuevo_folio := 1;


            -- ------------------------------------------------
            -- Crear nuevo libro
            -- ------------------------------------------------

            INSERT INTO diplomas.libro_registro (
                parametrizacion_registro_facultad_id,
                numero_libro,
                cantidad_folios
            )
            VALUES (
                v_param.id,
                v_nuevo_libro,
                v_param.folios_por_libro
            )

            ON CONFLICT (
                parametrizacion_registro_facultad_id,
                numero_libro
            )

            DO UPDATE
               SET activo = TRUE

            RETURNING *
            INTO v_libro;

        END IF;


        -- ----------------------------------------------------
        -- Crear siguiente folio
        -- ----------------------------------------------------

        INSERT INTO diplomas.folio_registro (
            libro_registro_id,
            numero_folio,
            cantidad_registros
        )
        VALUES (
            v_libro.id,
            v_nuevo_folio,
            v_param.registros_por_folio
        )

        ON CONFLICT (
            libro_registro_id,
            numero_folio
        )

        DO UPDATE
           SET activo = TRUE

        RETURNING *
        INTO v_folio;

    END IF;


    -- ========================================================
    -- ACTUALIZAR CONSECUTIVO DEL DIPLOMA
    -- ========================================================

    UPDATE diplomas.diploma_digital

    SET consecutivo_facultad =
        v_nuevo_consecutivo

    WHERE id = v_diploma.id;


    -- ========================================================
    -- CREAR REGISTRO DE GRADO
    -- ========================================================

    INSERT INTO diplomas.registro_grado (

        diploma_digital_id,

        folio_registro_id,

        numero_registro

    )
    VALUES (

        v_diploma.id,

        v_folio.id,

        v_nuevo_registro

    )

    RETURNING id
    INTO v_registro_grado_id;


    -- ========================================================
    -- ACTUALIZAR CONTROL
    -- ========================================================

    UPDATE diplomas.control_registro_facultad

    SET
        libro_actual_id =
            v_libro.id,

        folio_actual_id =
            v_folio.id,

        ultimo_consecutivo_facultad =
            v_nuevo_consecutivo,

        ultimo_numero_registro =
            v_nuevo_registro

    WHERE id = v_control.id;


    -- ========================================================
    -- RESPUESTA
    -- ========================================================

    registro_grado_id :=
        v_registro_grado_id;

    consecutivo_facultad_asignado :=
        v_nuevo_consecutivo;

    numero_libro :=
        v_libro.numero_libro;

    numero_folio :=
        v_folio.numero_folio;

    numero_registro :=
        v_nuevo_registro;


    RETURN NEXT;

END;
$$ LANGUAGE plpgsql;


COMMENT ON FUNCTION
diplomas.asignar_registro_grado(INTEGER) IS
'Asigna atomicamente consecutivo de facultad y ubicacion Libro/Folio/Registro a un diploma.';


-- ============================================================
-- 23. VISTA CONSOLIDADA
-- ============================================================

CREATE OR REPLACE VIEW
diplomas.vw_registro_grado_completo AS

SELECT

    rg.id AS registro_grado_id,

    dd.id AS diploma_digital_id,

    dd.documento_digital_id,

    dd.facultad_id,

    dd.vigencia,

    dd.consecutivo_diploma,

    dd.consecutivo_facultad,

    lr.numero_libro,

    fr.numero_folio,

    rg.numero_registro,

    dd.fecha_grado,

    dd.nombre_estudiante,

    dd.numero_documento_estudiante,

    dd.titulo_otorgado,

    rg.fecha_registro

FROM diplomas.registro_grado rg

INNER JOIN diplomas.diploma_digital dd
    ON dd.id = rg.diploma_digital_id

INNER JOIN diplomas.folio_registro fr
    ON fr.id = rg.folio_registro_id

INNER JOIN diplomas.libro_registro lr
    ON lr.id = fr.libro_registro_id

WHERE rg.activo IS TRUE;


COMMENT ON VIEW
diplomas.vw_registro_grado_completo IS
'Consulta consolidada de consecutivo de facultad, libro, folio y registro del diploma.';


-- ============================================================

`)
}

func (m *RegistroGradoModel20260928) Down() {
	m.SQL(`
DROP VIEW IF EXISTS diplomas.vw_registro_grado_completo;
DROP FUNCTION IF EXISTS diplomas.asignar_registro_grado(INTEGER);
DROP FUNCTION IF EXISTS diplomas.inicializar_parametrizacion_registro(INTEGER);
DROP TABLE IF EXISTS diplomas.ceremonia_estudiante;
DROP TABLE IF EXISTS diplomas.ceremonia_facultad_programa;
DROP TABLE IF EXISTS diplomas.ceremonia_grado;
DROP TABLE IF EXISTS diplomas.registro_grado;
DROP TABLE IF EXISTS diplomas.control_registro_facultad;
DROP TABLE IF EXISTS diplomas.folio_registro;
DROP TABLE IF EXISTS diplomas.libro_registro;
DROP TABLE IF EXISTS diplomas.parametrizacion_registro_facultad;
`)
}
