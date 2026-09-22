package models

import "time"

type DiplomaDigital struct {
	Id                        int64             `orm:"column(id);pk;auto" json:"id"`
	DocumentoDigital          *DocumentoDigital `orm:"column(documento_digital_id);rel(one)" json:"documento_digital"`
	FacultadId                int64             `orm:"column(facultad_id)" json:"facultad_id"`
	Vigencia                  int               `orm:"column(vigencia)" json:"vigencia"`
	FechaGrado                time.Time         `orm:"column(fecha_grado);type(date)" json:"fecha_grado"`
	NombreEstudiante          string            `orm:"column(nombre_estudiante);size(250)" json:"nombre_estudiante"`
	TipoDocumentoEstudiante   string            `orm:"column(tipo_documento_estudiante);size(20)" json:"tipo_documento_estudiante"`
	NumeroDocumentoEstudiante string            `orm:"column(numero_documento_estudiante);size(50)" json:"numero_documento_estudiante"`
	MunicipioExpedicion       string            `orm:"column(municipio_expedicion);size(150);null" json:"municipio_expedicion,omitempty"`
	TituloOtorgado            string            `orm:"column(titulo_otorgado);size(250)" json:"titulo_otorgado"`
	ConsecutivoDiploma        int64             `orm:"column(consecutivo_diploma);unique" json:"consecutivo_diploma"`
	ConsecutivoFacultad       int               `orm:"column(consecutivo_facultad)" json:"consecutivo_facultad"`
	Folio                     int               `orm:"column(folio)" json:"folio"`
	Acta                      int               `orm:"column(acta)" json:"acta"`
	Libro                     int               `orm:"column(libro)" json:"libro"`
	Activo                    bool              `orm:"column(activo)" json:"activo"`
	FechaCreacion             time.Time         `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion         time.Time         `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *DiplomaDigital) TableName() string {
	return "diploma_digital"
}

func (t *DiplomaDigital) TableIndex() [][]string {
	return [][]string{
		{"FacultadId", "Vigencia"},
	}
}

func (t *DiplomaDigital) TableUnique() [][]string {
	return [][]string{
		{"DocumentoDigital"},
		{"ConsecutivoDiploma"},
		{"FacultadId", "Vigencia", "ConsecutivoFacultad"},
		{"FacultadId", "Vigencia", "Libro", "Folio", "Acta"},
	}
}
