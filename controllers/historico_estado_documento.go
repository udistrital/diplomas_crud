package controllers

import (
	"net/http"
	"strconv"

	"github.com/udistrital/diplomas_crud/models"
	"github.com/udistrital/diplomas_crud/services"
)

type HistoricoEstadoDocumentoController struct {
	APIController
	Service services.HistoricoEstadoDocumentoService
}

func (c *HistoricoEstadoDocumentoController) URLMapping() {}

func (c *HistoricoEstadoDocumentoController) Post() {
	payload, err := decodeBody[models.HistoricoEstadoDocumento](&c.Controller)
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

func (c *HistoricoEstadoDocumentoController) GetAll() {
	items, err := c.Service.List(c.GetString("query"))
	if err != nil {
		c.writeError(http.StatusBadRequest, err.Error())
		return
	}

	c.writeJSON(http.StatusOK, items)
}

func (c *HistoricoEstadoDocumentoController) GetOne() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	item, err := c.Service.GetByID(id)
	if err != nil {
		c.writeError(mapReadError(err), "historico_estado_documento not found")
		return
	}

	c.writeJSON(http.StatusOK, item)
}

func (c *HistoricoEstadoDocumentoController) Put() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}

	payload, err := decodeBody[models.HistoricoEstadoDocumento](&c.Controller)
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

func (c *HistoricoEstadoDocumentoController) Delete() {
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
