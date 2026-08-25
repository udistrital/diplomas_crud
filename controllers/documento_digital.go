package controllers

import (
	"net/http"
	"strconv"

	"github.com/udistrital/diplomas_crud/models"
	"github.com/udistrital/diplomas_crud/services"
)

type DocumentoDigitalController struct {
	APIController
	Service  services.DocumentoDigitalService
	Workflow services.DocumentoDigitalWorkflowService
}

func (c *DocumentoDigitalController) URLMapping() {}

func (c *DocumentoDigitalController) Post() {
	payload, err := decodeBody[models.DocumentoDigital](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := c.Service.Create(payload)
	if err != nil {
		c.writeError(http.StatusInternalServerError, err.Error())
		return
	}

	c.writeJSON(http.StatusCreated, item)
}

func (c *DocumentoDigitalController) GetAll() {
	items, err := c.Service.List(c.GetString("query"))
	if err != nil {
		c.writeError(http.StatusBadRequest, err.Error())
		return
	}

	c.writeJSON(http.StatusOK, items)
}

func (c *DocumentoDigitalController) GetOne() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	item, err := c.Service.GetByID(id)
	if err != nil {
		c.writeError(mapReadError(err), "documento_digital not found")
		return
	}

	c.writeJSON(http.StatusOK, item)
}

func (c *DocumentoDigitalController) Put() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	payload, err := decodeBody[models.DocumentoDigital](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := c.Service.Update(id, payload)
	if err != nil {
		c.writeError(mapReadError(err), err.Error())
		return
	}

	c.writeJSON(http.StatusOK, item)
}

func (c *DocumentoDigitalController) Delete() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	if err := c.Service.SoftDelete(id); err != nil {
		c.writeError(mapReadError(err), err.Error())
		return
	}

	c.Ctx.Output.SetStatus(http.StatusNoContent)
}

func (c *DocumentoDigitalController) CambiarEstado() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	payload, err := decodeBody[services.CambiarEstadoDocumentoInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := c.Workflow.CambiarEstado(id, payload)
	if err != nil {
		c.writeError(mapReadError(err), err.Error())
		return
	}

	c.writeJSON(http.StatusCreated, item)
}

func (c *DocumentoDigitalController) ActualizarUUID() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	payload, err := decodeBody[services.ActualizarUUIDDocumentoInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := c.Workflow.ActualizarUUID(id, payload)
	if err != nil {
		c.writeError(mapReadError(err), err.Error())
		return
	}

	c.writeJSON(http.StatusOK, item)
}

func (c *DocumentoDigitalController) RegistrarFirma() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	payload, err := decodeBody[services.RegistrarFirmaDocumentoInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := c.Workflow.RegistrarFirma(id, payload)
	if err != nil {
		c.writeError(mapReadError(err), err.Error())
		return
	}

	c.writeJSON(http.StatusCreated, item)
}

func (c *DocumentoDigitalController) CrearDiploma() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	payload, err := decodeBody[services.CrearDiplomaDigitalInput](&c.Controller)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := c.Workflow.CrearDiploma(id, payload)
	if err != nil {
		c.writeError(mapReadError(err), err.Error())
		return
	}

	c.writeJSON(http.StatusCreated, item)
}
