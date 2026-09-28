package models

import "time"

type RegistroGrado struct {
	Id                int64           `orm:"column(id);pk;auto" json:"id"`
	DiplomaDigital    *DiplomaDigital `orm:"column(diploma_digital_id);rel(one)" json:"diploma_digital,omitempty"`
	FolioRegistro     *FolioRegistro  `orm:"column(folio_registro_id);rel(fk)" json:"folio_registro"`
	NumeroRegistro    int             `orm:"column(numero_registro)" json:"numero_registro"`
	FechaRegistro     time.Time       `orm:"column(fecha_registro);type(timestamp)" json:"fecha_registro"`
	Activo            bool            `orm:"column(activo)" json:"activo"`
	FechaCreacion     time.Time       `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion time.Time       `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *RegistroGrado) TableName() string { return "registro_grado" }

func (t *RegistroGrado) TableUnique() [][]string {
	return [][]string{{"DiplomaDigital"}, {"FolioRegistro", "NumeroRegistro"}}
}
