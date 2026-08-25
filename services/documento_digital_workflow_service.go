package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type CambiarEstadoDocumentoInput struct {
	EstadoNuevoId     int64  `json:"estado_nuevo_id"`
	TerceroIdFirmante *int64 `json:"tercero_id_firmante,omitempty"`
	RolActorId        *int64 `json:"rol_actor_id,omitempty"`
	FirmaDocumentoId  *int64 `json:"firma_documento_id,omitempty"`
	Observacion       string `json:"observacion,omitempty"`
}

type ActualizarUUIDDocumentoInput struct {
	UUIDVerificacion string `json:"uuid_verificacion"`
}

type RegistrarFirmaDocumentoInput struct {
	FirmaFirmanteId   *int64 `json:"firma_firmante_id,omitempty"`
	RolFirmanteId     int64  `json:"rol_firmante_id"`
	TerceroIdFirmante int64  `json:"tercero_id_firmante"`
	EstadoFirmadoId   *int64 `json:"estado_firmado_id,omitempty"`
	Observacion       string `json:"observacion,omitempty"`
}

type CrearDiplomaDigitalInput struct {
	FacultadId              int64     `json:"facultad_id"`
	Vigencia                int       `json:"vigencia"`
	FechaGrado              time.Time `json:"fecha_grado"`
	EstadoDocumentoCreadoId int64     `json:"estado_documento_creado_id"`
	UUIDVerificacion        *string   `json:"uuid_verificacion,omitempty"`
}

type DocumentoDigitalWorkflowService struct{}

func (s DocumentoDigitalWorkflowService) CambiarEstado(documentoID int64, input *CambiarEstadoDocumentoInput) (*models.HistoricoEstadoDocumento, error) {
	if input.EstadoNuevoId == 0 {
		return nil, errors.New("estado_nuevo_id is required")
	}

	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return nil, fmt.Errorf("begin cambiar estado: %w", err)
	}

	historicoID, err := registrarEstadoDocumentoTx(o, documentoID, input)
	if err != nil {
		_ = o.Rollback()
		return nil, err
	}

	if err := o.Commit(); err != nil {
		return nil, fmt.Errorf("commit cambiar estado documento_digital %d: %w", documentoID, err)
	}

	return HistoricoEstadoDocumentoService{}.GetByID(historicoID)
}

func (s DocumentoDigitalWorkflowService) ActualizarUUID(documentoID int64, input *ActualizarUUIDDocumentoInput) (*models.DocumentoDigital, error) {
	if input.UUIDVerificacion == "" {
		return nil, errors.New("uuid_verificacion is required")
	}

	o := orm.NewOrm()
	result, err := o.Raw(
		"UPDATE documento_digital SET uuid_verificacion = ?, fecha_modificacion = now() WHERE id = ?",
		input.UUIDVerificacion,
		documentoID,
	).Exec()
	if err != nil {
		return nil, fmt.Errorf("update uuid documento_digital %d: %w", documentoID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("check update uuid documento_digital %d: %w", documentoID, err)
	}
	if rows == 0 {
		return nil, fmt.Errorf("documento_digital %d no existe: %w", documentoID, orm.ErrNoRows)
	}

	return DocumentoDigitalService{}.GetByID(documentoID)
}

func (s DocumentoDigitalWorkflowService) RegistrarFirma(documentoID int64, input *RegistrarFirmaDocumentoInput) (*models.FirmaDocumento, error) {
	if input.RolFirmanteId == 0 {
		return nil, errors.New("rol_firmante_id is required")
	}
	if input.TerceroIdFirmante == 0 {
		return nil, errors.New("tercero_id_firmante is required")
	}

	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return nil, fmt.Errorf("begin registrar firma: %w", err)
	}

	firmaID, err := upsertFirmaDocumentoTx(o, documentoID, input)
	if err != nil {
		_ = o.Rollback()
		return nil, err
	}

	if input.EstadoFirmadoId != nil {
		_, err = registrarEstadoDocumentoTx(o, documentoID, &CambiarEstadoDocumentoInput{
			EstadoNuevoId:     *input.EstadoFirmadoId,
			TerceroIdFirmante: &input.TerceroIdFirmante,
			RolActorId:        &input.RolFirmanteId,
			FirmaDocumentoId:  &firmaID,
			Observacion:       input.Observacion,
		})
		if err != nil {
			_ = o.Rollback()
			return nil, err
		}
	}

	if err := o.Commit(); err != nil {
		return nil, fmt.Errorf("commit registrar firma documento_digital %d: %w", documentoID, err)
	}

	return FirmaDocumentoService{}.GetByID(firmaID)
}

func (s DocumentoDigitalWorkflowService) CrearDiploma(documentoID int64, input *CrearDiplomaDigitalInput) (*models.DiplomaDigital, error) {
	if input.FacultadId == 0 {
		return nil, errors.New("facultad_id is required")
	}
	if input.Vigencia == 0 {
		return nil, errors.New("vigencia is required")
	}
	if input.FechaGrado.IsZero() {
		return nil, errors.New("fecha_grado is required")
	}
	if input.EstadoDocumentoCreadoId == 0 {
		return nil, errors.New("estado_documento_creado_id is required")
	}

	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return nil, fmt.Errorf("begin crear diploma: %w", err)
	}

	diplomaID, err := crearDiplomaDigitalTx(o, documentoID, input)
	if err != nil {
		_ = o.Rollback()
		return nil, err
	}

	if err := o.Commit(); err != nil {
		return nil, fmt.Errorf("commit crear diploma documento_digital %d: %w", documentoID, err)
	}

	return DiplomaDigitalService{}.GetByID(diplomaID)
}

func registrarEstadoDocumentoTx(o orm.Ormer, documentoID int64, input *CambiarEstadoDocumentoInput) (int64, error) {
	var estadoAnteriorID int64
	err := o.Raw(
		"SELECT estado_documento_id FROM documento_digital WHERE id = ? FOR UPDATE",
		documentoID,
	).QueryRow(&estadoAnteriorID)
	if err != nil {
		return 0, fmt.Errorf("read estado documento_digital %d: %w", documentoID, err)
	}

	var historicoID int64
	err = o.Raw(
		`INSERT INTO historico_estado_documento (
			documento_digital_id,
			estado_anterior_id,
			estado_nuevo_id,
			tercero_id_firmante,
			rol_actor_id,
			firma_documento_id,
			observacion
		) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		documentoID,
		estadoAnteriorID,
		input.EstadoNuevoId,
		nullableInt64(input.TerceroIdFirmante),
		nullableInt64(input.RolActorId),
		nullableInt64(input.FirmaDocumentoId),
		input.Observacion,
	).QueryRow(&historicoID)
	if err != nil {
		return 0, fmt.Errorf("insert historico estado documento_digital %d: %w", documentoID, err)
	}

	_, err = o.Raw(
		"UPDATE documento_digital SET estado_documento_id = ?, fecha_modificacion = now() WHERE id = ?",
		input.EstadoNuevoId,
		documentoID,
	).Exec()
	if err != nil {
		return 0, fmt.Errorf("update estado documento_digital %d: %w", documentoID, err)
	}

	return historicoID, nil
}

func upsertFirmaDocumentoTx(o orm.Ormer, documentoID int64, input *RegistrarFirmaDocumentoInput) (int64, error) {
	var firmaID int64
	err := o.Raw(
		`UPDATE firma_documento
		 SET firma_firmante_id = ?, tercero_id_firmante = ?, activo = true, fecha_modificacion = now()
		 WHERE documento_digital_id = ? AND rol_firmante_id = ? AND activo IS TRUE
		 RETURNING id`,
		nullableInt64(input.FirmaFirmanteId),
		input.TerceroIdFirmante,
		documentoID,
		input.RolFirmanteId,
	).QueryRow(&firmaID)
	if err == nil {
		return firmaID, nil
	}
	if !errors.Is(err, orm.ErrNoRows) {
		return 0, fmt.Errorf("update firma_documento documento_digital %d: %w", documentoID, err)
	}

	err = o.Raw(
		`INSERT INTO firma_documento (
			documento_digital_id,
			firma_firmante_id,
			rol_firmante_id,
			tercero_id_firmante
		) VALUES (?, ?, ?, ?) RETURNING id`,
		documentoID,
		nullableInt64(input.FirmaFirmanteId),
		input.RolFirmanteId,
		input.TerceroIdFirmante,
	).QueryRow(&firmaID)
	if err != nil {
		return 0, fmt.Errorf("insert firma_documento documento_digital %d: %w", documentoID, err)
	}

	return firmaID, nil
}

func crearDiplomaDigitalTx(o orm.Ormer, documentoID int64, input *CrearDiplomaDigitalInput) (int64, error) {
	_, err := o.Raw(
		`INSERT INTO control_consecutivo_facultad_vigencia (facultad_id, vigencia)
		 VALUES (?, ?)
		 ON CONFLICT (facultad_id, vigencia) DO NOTHING`,
		input.FacultadId,
		input.Vigencia,
	).Exec()
	if err != nil {
		return 0, fmt.Errorf("ensure control consecutivo facultad %d vigencia %d: %w", input.FacultadId, input.Vigencia, err)
	}

	var controlID int64
	var ultimoConsecutivoFacultad int
	var ultimoFolio int
	var libroActual int
	var foliosPorLibro int
	err = o.Raw(
		`SELECT id, ultimo_consecutivo_facultad, ultimo_folio, libro_actual, folios_por_libro
		 FROM control_consecutivo_facultad_vigencia
		 WHERE facultad_id = ? AND vigencia = ? AND activo IS TRUE
		 FOR UPDATE`,
		input.FacultadId,
		input.Vigencia,
	).QueryRow(&controlID, &ultimoConsecutivoFacultad, &ultimoFolio, &libroActual, &foliosPorLibro)
	if err != nil {
		return 0, fmt.Errorf("lock control consecutivo facultad %d vigencia %d: %w", input.FacultadId, input.Vigencia, err)
	}

	siguienteConsecutivoFacultad := ultimoConsecutivoFacultad + 1
	siguienteLibro := libroActual
	siguienteFolio := ultimoFolio + 1
	if ultimoFolio >= foliosPorLibro {
		siguienteLibro = libroActual + 1
		siguienteFolio = 1
	}

	var consecutivoDiploma int64
	err = o.Raw("SELECT nextval('diplomas.seq_consecutivo_diploma')").QueryRow(&consecutivoDiploma)
	if err != nil {
		return 0, fmt.Errorf("next consecutivo diploma: %w", err)
	}

	var diplomaID int64
	err = o.Raw(
		`INSERT INTO diploma_digital (
			documento_digital_id,
			facultad_id,
			vigencia,
			fecha_grado,
			consecutivo_diploma,
			consecutivo_facultad,
			folio,
			libro
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		documentoID,
		input.FacultadId,
		input.Vigencia,
		input.FechaGrado,
		consecutivoDiploma,
		siguienteConsecutivoFacultad,
		siguienteFolio,
		siguienteLibro,
	).QueryRow(&diplomaID)
	if err != nil {
		return 0, fmt.Errorf("insert diploma_digital documento_digital %d: %w", documentoID, err)
	}

	_, err = o.Raw(
		`UPDATE control_consecutivo_facultad_vigencia
		 SET ultimo_consecutivo_facultad = ?, ultimo_folio = ?, libro_actual = ?, fecha_modificacion = now()
		 WHERE id = ?`,
		siguienteConsecutivoFacultad,
		siguienteFolio,
		siguienteLibro,
		controlID,
	).Exec()
	if err != nil {
		return 0, fmt.Errorf("update control consecutivo %d: %w", controlID, err)
	}

	_, err = o.Raw(
		`UPDATE documento_digital
		 SET uuid_verificacion = COALESCE(CAST(? AS uuid), uuid_verificacion), vigencia = COALESCE(vigencia, ?), fecha_modificacion = now()
		 WHERE id = ?`,
		nullableString(input.UUIDVerificacion),
		input.Vigencia,
		documentoID,
	).Exec()
	if err != nil {
		return 0, fmt.Errorf("update documento_digital %d after diploma: %w", documentoID, err)
	}

	_, err = registrarEstadoDocumentoTx(o, documentoID, &CambiarEstadoDocumentoInput{
		EstadoNuevoId: input.EstadoDocumentoCreadoId,
		Observacion:   "Diploma creado por Rectoria con consecutivos asignados",
	})
	if err != nil {
		return 0, err
	}

	return diplomaID, nil
}

func nullableInt64(value *int64) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func nullableString(value *string) interface{} {
	if value == nil {
		return nil
	}
	return *value
}
