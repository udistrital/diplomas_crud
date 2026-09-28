package models

import "time"

type ParametrizacionRegistroFacultad struct {
	Id                         int64     `orm:"column(id);pk;auto" json:"id"`
	FacultadId                 int64     `orm:"column(facultad_id)" json:"facultad_id"`
	Vigencia                   int       `orm:"column(vigencia)" json:"vigencia"`
	LibroInicial               int       `orm:"column(libro_inicial)" json:"libro_inicial"`
	FolioInicial               int       `orm:"column(folio_inicial)" json:"folio_inicial"`
	RegistroInicial            int       `orm:"column(registro_inicial)" json:"registro_inicial"`
	ConsecutivoFacultadInicial int       `orm:"column(consecutivo_facultad_inicial)" json:"consecutivo_facultad_inicial"`
	FoliosPorLibro             int       `orm:"column(folios_por_libro)" json:"folios_por_libro"`
	RegistrosPorFolio          int       `orm:"column(registros_por_folio)" json:"registros_por_folio"`
	Activo                     bool      `orm:"column(activo)" json:"activo"`
	FechaCreacion              time.Time `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion          time.Time `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *ParametrizacionRegistroFacultad) TableName() string {
	return "parametrizacion_registro_facultad"
}

func (t *ParametrizacionRegistroFacultad) TableUnique() [][]string {
	return [][]string{{"FacultadId", "Vigencia"}}
}
