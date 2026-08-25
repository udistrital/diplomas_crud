package services

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type FirmaFirmanteService struct{}

var firmaFirmanteFilterMap = map[string]string{
	"id":                  "id",
	"tercero_id_firmante": "tercero_id_firmante",
	"enlace_firma":        "enlace_firma",
	"activo":              "activo",
}

func (s FirmaFirmanteService) Create(input *models.FirmaFirmante) (*models.FirmaFirmante, error) {
	input.Activo = true
	o := orm.NewOrm()

	var current models.FirmaFirmante
	err := o.QueryTable(new(models.FirmaFirmante)).
		Filter("tercero_id_firmante", input.TerceroIdFirmante).
		Filter("activo", true).
		One(&current)
	if err == nil {
		current.EnlaceFirma = input.EnlaceFirma
		current.Activo = true
		if _, err := o.Update(&current); err != nil {
			return nil, fmt.Errorf("update firma_firmante activa: %w", err)
		}
		return &current, nil
	}
	if !errors.Is(err, orm.ErrNoRows) {
		return nil, fmt.Errorf("read firma_firmante activa: %w", err)
	}

	id, err := o.Insert(input)
	if err != nil {
		return nil, fmt.Errorf("insert firma_firmante: %w", err)
	}
	input.Id = id
	return input, nil
}

func (s FirmaFirmanteService) GetByID(id int64) (*models.FirmaFirmante, error) {
	o := orm.NewOrm()
	item := &models.FirmaFirmante{Id: id}
	if err := o.Read(item); err != nil {
		return nil, fmt.Errorf("read firma_firmante %d: %w", id, err)
	}
	return item, nil
}

func (s FirmaFirmanteService) List(rawQuery string) ([]*models.FirmaFirmante, error) {
	filters, err := ParseQueryFilters(rawQuery)
	if err != nil {
		return nil, err
	}

	qs := orm.NewOrm().QueryTable(new(models.FirmaFirmante))
	for key, value := range filters {
		ormField, ok := firmaFirmanteFilterMap[key]
		if !ok {
			return nil, fmt.Errorf("unsupported filter %q", key)
		}
		qs = qs.Filter(ormField, value)
	}

	var items []*models.FirmaFirmante
	_, err = qs.All(&items)
	if err != nil {
		return nil, fmt.Errorf("list firma_firmante: %w", err)
	}
	return items, nil
}

func (s FirmaFirmanteService) Update(id int64, input *models.FirmaFirmante) (*models.FirmaFirmante, error) {
	current, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	current.TerceroIdFirmante = input.TerceroIdFirmante
	current.EnlaceFirma = input.EnlaceFirma
	current.Activo = input.Activo

	o := orm.NewOrm()
	if _, err := o.Update(current); err != nil {
		return nil, fmt.Errorf("update firma_firmante %d: %w", id, err)
	}
	return current, nil
}

func (s FirmaFirmanteService) SoftDelete(id int64) error {
	current, err := s.GetByID(id)
	if err != nil {
		return err
	}

	current.Activo = false
	_, err = orm.NewOrm().Update(current, "activo", "fecha_modificacion")
	if err != nil {
		return fmt.Errorf("soft delete firma_firmante %d: %w", id, err)
	}
	return nil
}
