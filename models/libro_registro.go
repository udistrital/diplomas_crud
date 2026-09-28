package models

import "time"

type LibroRegistro struct {
	Id                              int64                            `orm:"column(id);pk;auto" json:"id"`
	ParametrizacionRegistroFacultad *ParametrizacionRegistroFacultad `orm:"column(parametrizacion_registro_facultad_id);rel(fk)" json:"parametrizacion_registro_facultad"`
	NumeroLibro                     int                              `orm:"column(numero_libro)" json:"numero_libro"`
	CantidadFolios                  int                              `orm:"column(cantidad_folios)" json:"cantidad_folios"`
	FechaApertura                   time.Time                        `orm:"column(fecha_apertura);type(date)" json:"fecha_apertura"`
	FechaCierre                     *time.Time                       `orm:"column(fecha_cierre);type(date);null" json:"fecha_cierre,omitempty"`
	Cerrado                         bool                             `orm:"column(cerrado)" json:"cerrado"`
	Activo                          bool                             `orm:"column(activo)" json:"activo"`
	FechaCreacion                   time.Time                        `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion               time.Time                        `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *LibroRegistro) TableName() string { return "libro_registro" }

func (t *LibroRegistro) TableUnique() [][]string {
	return [][]string{{"ParametrizacionRegistroFacultad", "NumeroLibro"}}
}
