package services

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/client/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type DiplomaDigitalService struct{}

var diplomaDigitalFilterMap = map[string]string{
	"id":                          "id",
	"documento_digital_id":        "DocumentoDigital__Id",
	"facultad_id":                 "facultad_id",
	"vigencia":                    "vigencia",
	"titulo_otorgado":             "titulo_otorgado",
	"nombre_estudiante":           "nombre_estudiante",
	"tipo_documento_estudiante":   "tipo_documento_estudiante",
	"numero_documento_estudiante": "numero_documento_estudiante",
	"municipio_expedicion":        "municipio_expedicion",
	"consecutivo_diploma":         "consecutivo_diploma",
	"consecutivo_facultad":        "consecutivo_facultad",
	"activo":                      "activo",
}

func (s DiplomaDigitalService) Create(input *models.DiplomaDigital) (*models.DiplomaDigital, error) {
	if input.DocumentoDigital == nil || input.DocumentoDigital.Id == 0 {
		return nil, errors.New("documento_digital.id is required")
	}

	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin create diploma_digital: %w", err)
	}

	var id int64
	err = tx.Raw(
		`INSERT INTO diploma_digital (
			documento_digital_id,
			facultad_id,
			vigencia,
			fecha_grado,
			nombre_estudiante,
			tipo_documento_estudiante,
			numero_documento_estudiante,
			municipio_expedicion,
			titulo_otorgado,
			activo
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, true) RETURNING id`,
		input.DocumentoDigital.Id,
		input.FacultadId,
		input.Vigencia,
		input.FechaGrado,
		input.NombreEstudiante,
		input.TipoDocumentoEstudiante,
		input.NumeroDocumentoEstudiante,
		input.MunicipioExpedicion,
		input.TituloOtorgado,
	).QueryRow(&id)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("insert diploma_digital: %w", err)
	}

	var registroID int64
	var consecutivoFacultad int
	var numeroLibro int
	var numeroFolio int
	var numeroRegistro int
	err = tx.Raw(
		`SELECT registro_grado_id, consecutivo_facultad_asignado, numero_libro, numero_folio, numero_registro
		 FROM diplomas.asignar_registro_grado(?)`,
		id,
	).QueryRow(&registroID, &consecutivoFacultad, &numeroLibro, &numeroFolio, &numeroRegistro)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("asignar registro_grado diploma_digital %d: %w", id, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create diploma_digital %d: %w", id, err)
	}

	_ = registroID
	_ = consecutivoFacultad
	_ = numeroLibro
	_ = numeroFolio
	_ = numeroRegistro
	return s.GetByID(id)
}

func (s DiplomaDigitalService) GetByID(id int64) (*models.DiplomaDigital, error) {
	o := orm.NewOrm()
	item := &models.DiplomaDigital{}
	if err := o.QueryTable(new(models.DiplomaDigital)).Filter("id", id).RelatedSel().One(item); err != nil {
		return nil, fmt.Errorf("read diploma_digital %d: %w", id, err)
	}
	return s.enrichRegistroGrado(item), nil
}

func (s DiplomaDigitalService) List(rawQuery string) ([]*models.DiplomaDigital, error) {
	filters, err := ParseQueryFilters(rawQuery)
	if err != nil {
		return nil, err
	}

	qs := orm.NewOrm().QueryTable(new(models.DiplomaDigital))
	for key, value := range filters {
		ormField, ok := diplomaDigitalFilterMap[key]
		if !ok {
			return nil, fmt.Errorf("unsupported filter %q", key)
		}
		qs = qs.Filter(ormField, value)
	}

	var items []*models.DiplomaDigital
	_, err = qs.RelatedSel().All(&items)
	if err != nil {
		return nil, fmt.Errorf("list diploma_digital: %w", err)
	}
	for _, item := range items {
		s.enrichRegistroGrado(item)
	}
	return items, nil
}

func (s DiplomaDigitalService) enrichRegistroGrado(item *models.DiplomaDigital) *models.DiplomaDigital {
	registro := &models.RegistroGrado{}
	err := orm.NewOrm().QueryTable(new(models.RegistroGrado)).
		Filter("DiplomaDigital__Id", item.Id).
		Filter("activo", true).
		RelatedSel().
		One(registro)
	if err == nil {
		item.RegistroGrado = registro
	}
	return item
}

func (s DiplomaDigitalService) Update(id int64, input *models.DiplomaDigital) (*models.DiplomaDigital, error) {
	current, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	current.DocumentoDigital = input.DocumentoDigital
	current.FacultadId = input.FacultadId
	current.Vigencia = input.Vigencia
	current.FechaGrado = input.FechaGrado
	current.NombreEstudiante = input.NombreEstudiante
	current.TipoDocumentoEstudiante = input.TipoDocumentoEstudiante
	current.NumeroDocumentoEstudiante = input.NumeroDocumentoEstudiante
	current.MunicipioExpedicion = input.MunicipioExpedicion
	current.TituloOtorgado = input.TituloOtorgado
	current.ConsecutivoDiploma = input.ConsecutivoDiploma
	current.ConsecutivoFacultad = input.ConsecutivoFacultad
	current.Activo = input.Activo

	o := orm.NewOrm()
	if _, err := o.Update(current); err != nil {
		return nil, fmt.Errorf("update diploma_digital %d: %w", id, err)
	}
	return s.GetByID(id)
}

func (s DiplomaDigitalService) SoftDelete(id int64) error {
	current, err := s.GetByID(id)
	if err != nil {
		return err
	}

	current.Activo = false
	_, err = orm.NewOrm().Update(current, "activo", "fecha_modificacion")
	if err != nil {
		return fmt.Errorf("soft delete diploma_digital %d: %w", id, err)
	}
	return nil
}
