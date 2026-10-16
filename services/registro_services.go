package services

import "github.com/udistrital/diplomas_crud/models"

type ParametrizacionRegistroFacultadService struct {
	SimpleCRUDService[models.ParametrizacionRegistroFacultad]
}

func (s ParametrizacionRegistroFacultadService) List(rawQuery string) ([]*models.ParametrizacionRegistroFacultad, error) {
	s.FilterMap = map[string]string{
		"id":          "id",
		"facultad_id": "FacultadId",
		"FacultadId":  "FacultadId",
		"vigencia":    "Vigencia",
		"Vigencia":    "Vigencia",
		"activo":      "Activo",
		"Activo":      "Activo",
	}
	return s.SimpleCRUDService.List(rawQuery)
}

type LibroRegistroService struct {
	SimpleCRUDService[models.LibroRegistro]
}
type FolioRegistroService struct {
	SimpleCRUDService[models.FolioRegistro]
}
type ControlRegistroFacultadService struct {
	SimpleCRUDService[models.ControlRegistroFacultad]
}
type RegistroGradoService struct {
	SimpleCRUDService[models.RegistroGrado]
}
type CeremoniaGradoService struct {
	SimpleCRUDService[models.CeremoniaGrado]
}
type CeremoniaFacultadProgramaService struct {
	SimpleCRUDService[models.CeremoniaFacultadPrograma]
}
type CeremoniaEstudianteService struct {
	SimpleCRUDService[models.CeremoniaEstudiante]
}

func (s CeremoniaEstudianteService) List(rawQuery string) ([]*models.CeremoniaEstudiante, error) {
	s.FilterMap = map[string]string{
		"id":                             "id",
		"ceremonia_facultad_programa_id": "CeremoniaFacultadPrograma__Id",
		"documento_digital_id":           "DocumentoDigital__Id",
		"CeremoniaFacultadPrograma__Id":  "CeremoniaFacultadPrograma__Id",
		"DocumentoDigital__Id":           "DocumentoDigital__Id",
		"activo":                         "activo",
	}
	return s.SimpleCRUDService.List(rawQuery)
}
