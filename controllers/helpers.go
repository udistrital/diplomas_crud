package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type APIController struct {
	beego.Controller
}

type apiError struct {
	Message string `json:"message"`
}

func (c *APIController) writeJSON(status int, payload interface{}) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = payload
	c.ServeJSON()
}

func (c *APIController) writeError(status int, message string) {
	c.writeJSON(status, apiError{Message: message})
}

func decodeBody[T any](controller *beego.Controller) (*T, error) {
	var payload T
	if err := json.NewDecoder(controller.Ctx.Request.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func mapReadError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if errors.Is(err, orm.ErrNoRows) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
