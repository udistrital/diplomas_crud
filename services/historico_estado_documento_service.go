package services

import (
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type HistoricoEstadoDocumentoService struct{}

var historicoFilterMap = map[string]string{
	"id":                   "id",
	"documento_digital_id": "DocumentoDigital__Id",
	"estado_nuevo_id":      "estado_nuevo_id",
	"tercero_id":           "tercero_id",
	"activo":               "activo",
}

func (s HistoricoEstadoDocumentoService) Create(input *models.HistoricoEstadoDocumento) (*models.HistoricoEstadoDocumento, error) {
	input.Activo = true
	o := orm.NewOrm()
	id, err := o.Insert(input)
	if err != nil {
		return nil, fmt.Errorf("insert historico_estado_documento: %w", err)
	}
	input.Id = id
	return input, nil
}

func (s HistoricoEstadoDocumentoService) GetByID(id int64) (*models.HistoricoEstadoDocumento, error) {
	o := orm.NewOrm()
	item := &models.HistoricoEstadoDocumento{}
	if err := o.QueryTable(new(models.HistoricoEstadoDocumento)).Filter("id", id).RelatedSel().One(item); err != nil {
		return nil, fmt.Errorf("read historico_estado_documento %d: %w", id, err)
	}
	return item, nil
}

func (s HistoricoEstadoDocumentoService) List(rawQuery string) ([]*models.HistoricoEstadoDocumento, error) {
	filters, err := ParseQueryFilters(rawQuery)
	if err != nil {
		return nil, err
	}

	qs := orm.NewOrm().QueryTable(new(models.HistoricoEstadoDocumento))
	for key, value := range filters {
		ormField, ok := historicoFilterMap[key]
		if !ok {
			return nil, fmt.Errorf("unsupported filter %q", key)
		}

		qs = qs.Filter(ormField, value)
	}

	var items []*models.HistoricoEstadoDocumento
	_, err = qs.RelatedSel().All(&items)
	if err != nil {
		return nil, fmt.Errorf("list historico_estado_documento: %w", err)
	}

	return items, nil
}

func (s HistoricoEstadoDocumentoService) Update(id int64, input *models.HistoricoEstadoDocumento) (*models.HistoricoEstadoDocumento, error) {
	current, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	current.DocumentoDigital = input.DocumentoDigital
	current.EstadoAnteriorId = input.EstadoAnteriorId
	current.EstadoNuevoId = input.EstadoNuevoId
	current.TerceroId = input.TerceroId
	current.Observacion = input.Observacion
	current.Metadata = input.Metadata
	current.Activo = input.Activo

	o := orm.NewOrm()
	if _, err := o.Update(current); err != nil {
		return nil, fmt.Errorf("update historico_estado_documento %d: %w", id, err)
	}

	return current, nil
}

func (s HistoricoEstadoDocumentoService) SoftDelete(id int64) error {
	current, err := s.GetByID(id)
	if err != nil {
		return err
	}

	current.Activo = false
	_, err = orm.NewOrm().Update(current, "activo", "fecha_modificacion")
	if err != nil {
		return fmt.Errorf("soft delete historico_estado_documento %d: %w", id, err)
	}

	return nil
}
