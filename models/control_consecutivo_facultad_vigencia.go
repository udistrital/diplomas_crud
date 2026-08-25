package models

import "time"

type ControlConsecutivoFacultadVigencia struct {
	Id                        int64     `orm:"column(id);pk;auto" json:"id"`
	FacultadId                int64     `orm:"column(facultad_id)" json:"facultad_id"`
	Vigencia                  int       `orm:"column(vigencia)" json:"vigencia"`
	UltimoConsecutivoFacultad int       `orm:"column(ultimo_consecutivo_facultad)" json:"ultimo_consecutivo_facultad"`
	UltimoFolio               int       `orm:"column(ultimo_folio)" json:"ultimo_folio"`
	LibroActual               int       `orm:"column(libro_actual)" json:"libro_actual"`
	FoliosPorLibro            int       `orm:"column(folios_por_libro)" json:"folios_por_libro"`
	Activo                    bool      `orm:"column(activo)" json:"activo"`
	FechaCreacion             time.Time `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion         time.Time `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *ControlConsecutivoFacultadVigencia) TableName() string {
	return "control_consecutivo_facultad_vigencia"
}

func (t *ControlConsecutivoFacultadVigencia) TableIndex() [][]string {
	return [][]string{
		{"FacultadId", "Vigencia"},
	}
}

func (t *ControlConsecutivoFacultadVigencia) TableUnique() [][]string {
	return [][]string{
		{"FacultadId", "Vigencia"},
	}
}
