package controllers

import (
	"github.com/udistrital/diplomas_crud/models"
	"github.com/udistrital/diplomas_crud/services"
)

type ParametrizacionRegistroFacultadController struct {
	SimpleCRUDController[models.ParametrizacionRegistroFacultad]
}

func (c *ParametrizacionRegistroFacultadController) Prepare() {
	c.Service = services.ParametrizacionRegistroFacultadService{}
}

type LibroRegistroController struct {
	SimpleCRUDController[models.LibroRegistro]
}

func (c *LibroRegistroController) Prepare() { c.Service = services.LibroRegistroService{} }

type FolioRegistroController struct {
	SimpleCRUDController[models.FolioRegistro]
}

func (c *FolioRegistroController) Prepare() { c.Service = services.FolioRegistroService{} }

type ControlRegistroFacultadController struct {
	SimpleCRUDController[models.ControlRegistroFacultad]
}

func (c *ControlRegistroFacultadController) Prepare() {
	c.Service = services.ControlRegistroFacultadService{}
}

type RegistroGradoController struct {
	SimpleCRUDController[models.RegistroGrado]
}

func (c *RegistroGradoController) Prepare() { c.Service = services.RegistroGradoService{} }

type CeremoniaGradoController struct {
	SimpleCRUDController[models.CeremoniaGrado]
}

func (c *CeremoniaGradoController) Prepare() { c.Service = services.CeremoniaGradoService{} }

type CeremoniaFacultadProgramaController struct {
	SimpleCRUDController[models.CeremoniaFacultadPrograma]
}

func (c *CeremoniaFacultadProgramaController) Prepare() {
	c.Service = services.CeremoniaFacultadProgramaService{}
}

type CeremoniaEstudianteController struct {
	SimpleCRUDController[models.CeremoniaEstudiante]
}

func (c *CeremoniaEstudianteController) Prepare() { c.Service = services.CeremoniaEstudianteService{} }
