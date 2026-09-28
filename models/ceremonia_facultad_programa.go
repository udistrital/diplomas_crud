package models

import "time"

type CeremoniaFacultadPrograma struct {
	Id                            int64           `orm:"column(id);pk;auto" json:"id"`
	CeremoniaGrado                *CeremoniaGrado `orm:"column(ceremonia_grado_id);rel(fk)" json:"ceremonia_grado"`
	FacultadId                    int64           `orm:"column(facultad_id)" json:"facultad_id"`
	ProgramaAcademicoId           int64           `orm:"column(programa_academico_id)" json:"programa_academico_id"`
	CantidadEstudiantesProgramada *int            `orm:"column(cantidad_estudiantes_programada);null" json:"cantidad_estudiantes_programada,omitempty"`
	Activo                        bool            `orm:"column(activo)" json:"activo"`
	FechaCreacion                 time.Time       `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion             time.Time       `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *CeremoniaFacultadPrograma) TableName() string { return "ceremonia_facultad_programa" }

func (t *CeremoniaFacultadPrograma) TableUnique() [][]string {
	return [][]string{{"CeremoniaGrado", "FacultadId", "ProgramaAcademicoId"}}
}
