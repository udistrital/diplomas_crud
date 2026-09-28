package models

import "time"

type ControlRegistroFacultad struct {
	Id                              int64                            `orm:"column(id);pk;auto" json:"id"`
	ParametrizacionRegistroFacultad *ParametrizacionRegistroFacultad `orm:"column(parametrizacion_registro_facultad_id);rel(fk)" json:"parametrizacion_registro_facultad"`
	LibroActual                     *LibroRegistro                   `orm:"column(libro_actual_id);rel(fk);null" json:"libro_actual,omitempty"`
	FolioActual                     *FolioRegistro                   `orm:"column(folio_actual_id);rel(fk);null" json:"folio_actual,omitempty"`
	UltimoConsecutivoFacultad       int                              `orm:"column(ultimo_consecutivo_facultad)" json:"ultimo_consecutivo_facultad"`
	UltimoNumeroRegistro            int                              `orm:"column(ultimo_numero_registro)" json:"ultimo_numero_registro"`
	Activo                          bool                             `orm:"column(activo)" json:"activo"`
	FechaCreacion                   time.Time                        `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion               time.Time                        `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *ControlRegistroFacultad) TableName() string { return "control_registro_facultad" }

func (t *ControlRegistroFacultad) TableUnique() [][]string {
	return [][]string{{"ParametrizacionRegistroFacultad"}}
}
