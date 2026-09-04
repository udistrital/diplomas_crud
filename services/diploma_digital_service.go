package services

import (
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type DiplomaDigitalService struct{}

var diplomaDigitalFilterMap = map[string]string{
	"id":                   "id",
	"documento_digital_id": "DocumentoDigital__Id",
	"facultad_id":          "facultad_id",
	"vigencia":             "vigencia",
	"consecutivo_diploma":  "consecutivo_diploma",
	"consecutivo_facultad": "consecutivo_facultad",
	"folio":                "folio",
	"acta":                 "acta",
	"libro":                "libro",
	"activo":               "activo",
}

func (s DiplomaDigitalService) Create(input *models.DiplomaDigital) (*models.DiplomaDigital, error) {
	input.Activo = true
	o := orm.NewOrm()
	id, err := o.Insert(input)
	if err != nil {
		return nil, fmt.Errorf("insert diploma_digital: %w", err)
	}
	input.Id = id
	return s.GetByID(id)
}

func (s DiplomaDigitalService) GetByID(id int64) (*models.DiplomaDigital, error) {
	o := orm.NewOrm()
	item := &models.DiplomaDigital{}
	if err := o.QueryTable(new(models.DiplomaDigital)).Filter("id", id).RelatedSel().One(item); err != nil {
		return nil, fmt.Errorf("read diploma_digital %d: %w", id, err)
	}
	return item, nil
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
	return items, nil
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
	current.ConsecutivoDiploma = input.ConsecutivoDiploma
	current.ConsecutivoFacultad = input.ConsecutivoFacultad
	current.Folio = input.Folio
	current.Acta = input.Acta
	current.Libro = input.Libro
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
