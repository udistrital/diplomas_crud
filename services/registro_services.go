package services

import "github.com/udistrital/diplomas_crud/models"

type ParametrizacionRegistroFacultadService struct {
	SimpleCRUDService[models.ParametrizacionRegistroFacultad]
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
