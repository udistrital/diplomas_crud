package services

import (
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/udistrital/diplomas_crud/models"
)

type ControlConsecutivoFacultadVigenciaService struct{}

var controlConsecutivoFilterMap = map[string]string{
	"id":                          "id",
	"facultad_id":                 "facultad_id",
	"vigencia":                    "vigencia",
	"ultimo_consecutivo_facultad": "ultimo_consecutivo_facultad",
	"ultimo_folio":                "ultimo_folio",
	"libro_actual":                "libro_actual",
	"folios_por_libro":            "folios_por_libro",
	"activo":                      "activo",
}

func (s ControlConsecutivoFacultadVigenciaService) Create(input *models.ControlConsecutivoFacultadVigencia) (*models.ControlConsecutivoFacultadVigencia, error) {
	input.Activo = true
	if input.LibroActual == 0 {
		input.LibroActual = 1
	}
	if input.FoliosPorLibro == 0 {
		input.FoliosPorLibro = 500
	}

	o := orm.NewOrm()
	id, err := o.Insert(input)
	if err != nil {
		return nil, fmt.Errorf("insert control_consecutivo_facultad_vigencia: %w", err)
	}
	input.Id = id
	return input, nil
}

func (s ControlConsecutivoFacultadVigenciaService) GetByID(id int64) (*models.ControlConsecutivoFacultadVigencia, error) {
	o := orm.NewOrm()
	item := &models.ControlConsecutivoFacultadVigencia{Id: id}
	if err := o.Read(item); err != nil {
		return nil, fmt.Errorf("read control_consecutivo_facultad_vigencia %d: %w", id, err)
	}
	return item, nil
}

func (s ControlConsecutivoFacultadVigenciaService) List(rawQuery string) ([]*models.ControlConsecutivoFacultadVigencia, error) {
	filters, err := ParseQueryFilters(rawQuery)
	if err != nil {
		return nil, err
	}

	qs := orm.NewOrm().QueryTable(new(models.ControlConsecutivoFacultadVigencia))
	for key, value := range filters {
		ormField, ok := controlConsecutivoFilterMap[key]
		if !ok {
			return nil, fmt.Errorf("unsupported filter %q", key)
		}
		qs = qs.Filter(ormField, value)
	}

	var items []*models.ControlConsecutivoFacultadVigencia
	_, err = qs.All(&items)
	if err != nil {
		return nil, fmt.Errorf("list control_consecutivo_facultad_vigencia: %w", err)
	}
	return items, nil
}

func (s ControlConsecutivoFacultadVigenciaService) Update(id int64, input *models.ControlConsecutivoFacultadVigencia) (*models.ControlConsecutivoFacultadVigencia, error) {
	current, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	current.FacultadId = input.FacultadId
	current.Vigencia = input.Vigencia
	current.UltimoConsecutivoFacultad = input.UltimoConsecutivoFacultad
	current.UltimoFolio = input.UltimoFolio
	current.LibroActual = input.LibroActual
	current.FoliosPorLibro = input.FoliosPorLibro
	current.Activo = input.Activo

	o := orm.NewOrm()
	if _, err := o.Update(current); err != nil {
		return nil, fmt.Errorf("update control_consecutivo_facultad_vigencia %d: %w", id, err)
	}
	return current, nil
}

func (s ControlConsecutivoFacultadVigenciaService) SoftDelete(id int64) error {
	current, err := s.GetByID(id)
	if err != nil {
		return err
	}

	current.Activo = false
	_, err = orm.NewOrm().Update(current, "activo", "fecha_modificacion")
	if err != nil {
		return fmt.Errorf("soft delete control_consecutivo_facultad_vigencia %d: %w", id, err)
	}
	return nil
}
