package services

import (
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type FirmaDocumentoService struct{}

var firmaDocumentoFilterMap = map[string]string{
	"id":                   "id",
	"documento_digital_id": "DocumentoDigital__Id",
	"firma_firmante_id":    "FirmaFirmante__Id",
	"rol_firmante_id":      "rol_firmante_id",
	"tercero_id_firmante":  "tercero_id_firmante",
	"activo":               "activo",
}

func (s FirmaDocumentoService) Create(input *models.FirmaDocumento) (*models.FirmaDocumento, error) {
	input.Activo = true
	o := orm.NewOrm()
	id, err := o.Insert(input)
	if err != nil {
		return nil, fmt.Errorf("insert firma_documento: %w", err)
	}
	input.Id = id
	return s.GetByID(id)
}

func (s FirmaDocumentoService) GetByID(id int64) (*models.FirmaDocumento, error) {
	o := orm.NewOrm()
	item := &models.FirmaDocumento{}
	if err := o.QueryTable(new(models.FirmaDocumento)).Filter("id", id).RelatedSel().One(item); err != nil {
		return nil, fmt.Errorf("read firma_documento %d: %w", id, err)
	}
	return item, nil
}

func (s FirmaDocumentoService) List(rawQuery string) ([]*models.FirmaDocumento, error) {
	filters, err := ParseQueryFilters(rawQuery)
	if err != nil {
		return nil, err
	}

	qs := orm.NewOrm().QueryTable(new(models.FirmaDocumento))
	for key, value := range filters {
		ormField, ok := firmaDocumentoFilterMap[key]
		if !ok {
			return nil, fmt.Errorf("unsupported filter %q", key)
		}
		qs = qs.Filter(ormField, value)
	}

	var items []*models.FirmaDocumento
	_, err = qs.RelatedSel().All(&items)
	if err != nil {
		return nil, fmt.Errorf("list firma_documento: %w", err)
	}
	return items, nil
}

func (s FirmaDocumentoService) Update(id int64, input *models.FirmaDocumento) (*models.FirmaDocumento, error) {
	current, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	current.DocumentoDigital = input.DocumentoDigital
	current.FirmaFirmante = input.FirmaFirmante
	current.RolFirmanteId = input.RolFirmanteId
	current.TerceroIdFirmante = input.TerceroIdFirmante
	current.Activo = input.Activo

	o := orm.NewOrm()
	if _, err := o.Update(current); err != nil {
		return nil, fmt.Errorf("update firma_documento %d: %w", id, err)
	}
	return s.GetByID(id)
}

func (s FirmaDocumentoService) SoftDelete(id int64) error {
	current, err := s.GetByID(id)
	if err != nil {
		return err
	}

	current.Activo = false
	_, err = orm.NewOrm().Update(current, "activo", "fecha_modificacion")
	if err != nil {
		return fmt.Errorf("soft delete firma_documento %d: %w", id, err)
	}
	return nil
}
