package models

import "time"

type CeremoniaGrado struct {
	Id                            int64     `orm:"column(id);pk;auto" json:"id"`
	FechaCeremonia                time.Time `orm:"column(fecha_ceremonia);type(date)" json:"fecha_ceremonia"`
	HoraCeremonia                 string    `orm:"column(hora_ceremonia);type(time);null" json:"hora_ceremonia,omitempty"`
	Lugar                         string    `orm:"column(lugar);size(300)" json:"lugar"`
	CantidadEstudiantesProgramada int       `orm:"column(cantidad_estudiantes_programada)" json:"cantidad_estudiantes_programada"`
	Observacion                   string    `orm:"column(observacion);size(500);null" json:"observacion,omitempty"`
	Activo                        bool      `orm:"column(activo)" json:"activo"`
	FechaCreacion                 time.Time `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion             time.Time `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *CeremoniaGrado) TableName() string { return "ceremonia_grado" }
