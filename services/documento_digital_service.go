package services

import (
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type DocumentoDigitalService struct{}

var documentoDigitalFilterMap = map[string]string{
	"id":                    "id",
	"tipo_documento_id":     "tipo_documento_id",
	"estado_documento_id":   "estado_documento_id",
	"tercero_id_estudiante": "tercero_id_estudiante",
	"programa_academico_id": "programa_academico_id",
	"periodo_id":            "periodo_id",
	"vigencia":              "vigencia",
	"uuid_documento":        "uuid_documento",
	"activo":                "activo",
}

func (s DocumentoDigitalService) Create(input *models.DocumentoDigital) (*models.DocumentoDigital, error) {
	input.Activo = true
	o := orm.NewOrm()
	id, err := o.Insert(input)
	if err != nil {
		return nil, fmt.Errorf("insert documento_digital: %w", err)
	}
	input.Id = id
	return input, nil
}

func (s DocumentoDigitalService) GetByID(id int64) (*models.DocumentoDigital, error) {
	o := orm.NewOrm()
	item := &models.DocumentoDigital{Id: id}
	if err := o.Read(item); err != nil {
		return nil, fmt.Errorf("read documento_digital %d: %w", id, err)
	}
	return item, nil
}

func (s DocumentoDigitalService) List(rawQuery string) ([]*models.DocumentoDigital, error) {
	filters, err := ParseQueryFilters(rawQuery)
	if err != nil {
		return nil, err
	}

	qs := orm.NewOrm().QueryTable(new(models.DocumentoDigital))
	for key, value := range filters {
		ormField, ok := documentoDigitalFilterMap[key]
		if !ok {
			return nil, fmt.Errorf("unsupported filter %q", key)
		}
		qs = qs.Filter(ormField, value)
	}

	var items []*models.DocumentoDigital
	_, err = qs.All(&items)
	if err != nil {
		return nil, fmt.Errorf("list documento_digital: %w", err)
	}
	return items, nil
}

func (s DocumentoDigitalService) Update(id int64, input *models.DocumentoDigital) (*models.DocumentoDigital, error) {
	current, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	current.TipoDocumentoId = input.TipoDocumentoId
	current.EstadoDocumentoId = input.EstadoDocumentoId
	current.TerceroIdEstudiante = input.TerceroIdEstudiante
	current.ProgramaAcademicoId = input.ProgramaAcademicoId
	current.PeriodoId = input.PeriodoId
	current.Vigencia = input.Vigencia
	current.UUIDDocumento = input.UUIDDocumento
	current.Activo = input.Activo

	o := orm.NewOrm()
	if _, err := o.Update(current); err != nil {
		return nil, fmt.Errorf("update documento_digital %d: %w", id, err)
	}

	return current, nil
}

func (s DocumentoDigitalService) SoftDelete(id int64) error {
	current, err := s.GetByID(id)
	if err != nil {
		return err
	}

	current.Activo = false
	_, err = orm.NewOrm().Update(current, "activo", "fecha_modificacion")
	if err != nil {
		return fmt.Errorf("soft delete documento_digital %d: %w", id, err)
	}

	return nil
}
