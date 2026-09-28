package models

import "time"

type FolioRegistro struct {
	Id                int64          `orm:"column(id);pk;auto" json:"id"`
	LibroRegistro     *LibroRegistro `orm:"column(libro_registro_id);rel(fk)" json:"libro_registro"`
	NumeroFolio       int            `orm:"column(numero_folio)" json:"numero_folio"`
	CantidadRegistros int            `orm:"column(cantidad_registros)" json:"cantidad_registros"`
	Cerrado           bool           `orm:"column(cerrado)" json:"cerrado"`
	Activo            bool           `orm:"column(activo)" json:"activo"`
	FechaCreacion     time.Time      `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion time.Time      `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *FolioRegistro) TableName() string { return "folio_registro" }

func (t *FolioRegistro) TableUnique() [][]string {
	return [][]string{{"LibroRegistro", "NumeroFolio"}}
}
