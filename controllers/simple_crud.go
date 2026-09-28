package controllers

import (
	"net/http"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type simpleCRUDService[T any] interface {
	Create(*T) (*T, error)
	GetByID(int64) (*T, error)
	List(string) ([]*T, error)
	Update(int64, *T) (*T, error)
	SoftDelete(int64) error
}

type SimpleCRUDController[T any] struct {
	APIController
	Service simpleCRUDService[T]
}

func (c *SimpleCRUDController[T]) URLMapping() {}

func (c *SimpleCRUDController[T]) Prepare() {
	if c.Service == nil {
		c.Service = defaultCRUDService[T]()
	}
}

func (c *SimpleCRUDController[T]) Post() {
	payload, err := decodeBody[T](&c.Controller)
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

func (c *SimpleCRUDController[T]) GetAll() {
	items, err := c.Service.List(c.GetString("query"))
	if err != nil {
		c.writeError(http.StatusBadRequest, err.Error())
		return
	}
	c.writeJSON(http.StatusOK, items)
}

func (c *SimpleCRUDController[T]) GetOne() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}
	item, err := c.Service.GetByID(id)
	if err != nil {
		c.writeError(mapReadError(err), "resource not found")
		return
	}
	c.writeJSON(http.StatusOK, item)
}

func (c *SimpleCRUDController[T]) Put() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		c.writeError(http.StatusBadRequest, "invalid id")
		return
	}
	payload, err := decodeBody[T](&c.Controller)
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

func (c *SimpleCRUDController[T]) Delete() {
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

func defaultCRUDService[T any]() simpleCRUDService[T] {
	panic("default CRUD service not registered")
}

var _ = beego.Controller{}
