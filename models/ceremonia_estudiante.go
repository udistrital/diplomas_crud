package models

import "time"

type CeremoniaEstudiante struct {
	Id                        int64                      `orm:"column(id);pk;auto" json:"id"`
	CeremoniaFacultadPrograma *CeremoniaFacultadPrograma `orm:"column(ceremonia_facultad_programa_id);rel(fk)" json:"ceremonia_facultad_programa"`
	DocumentoDigital          *DocumentoDigital          `orm:"column(documento_digital_id);rel(fk)" json:"documento_digital"`
	Activo                    bool                       `orm:"column(activo)" json:"activo"`
	FechaCreacion             time.Time                  `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion         time.Time                  `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *CeremoniaEstudiante) TableName() string { return "ceremonia_estudiante" }

func (t *CeremoniaEstudiante) TableUnique() [][]string {
	return [][]string{{"CeremoniaFacultadPrograma", "DocumentoDigital"}}
}
