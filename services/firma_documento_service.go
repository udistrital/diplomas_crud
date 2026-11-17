package services

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type FirmaDocumentoService struct{}

var firmaDocumentoFilterMap = map[string]string{
	"id":                   "id",
	"documento_digital_id": "DocumentoDigital__Id",
	"firma_firmante_id":    "FirmaFirmante__Id",
	"rol_firmante_id":      "rol_firmante_id",
	"documento_identidad":  "documento_identidad",
	"activo":               "activo",
}

func (s FirmaDocumentoService) Create(input *models.FirmaDocumento) (*models.FirmaDocumento, error) {
	o := orm.NewOrm()
	if input == nil || input.DocumentoDigital == nil {
		return nil, ErrFirmaConflict
	}
	tx, err := o.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if err := bloquearFirmasDocumentoTx(tx, input.DocumentoDigital.Id, nil, false); err != nil {
		return nil, err
	}
	input.Activo = true
	id, err := tx.Insert(input)
	if err != nil {
		return nil, fmt.Errorf("insert firma_documento: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
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
	return nil, fmt.Errorf("%w: utilice el flujo transaccional de firma del documento", ErrFirmaConflict)
}

func (s FirmaDocumentoService) SoftDelete(id int64) error {
	return fmt.Errorf("%w: utilice la reversión validada del MID", ErrFirmaConflict)
}
