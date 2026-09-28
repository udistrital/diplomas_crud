package models

import (
	"time"
)

type DocumentoDigital struct {
	Id                        int64                       `orm:"column(id);pk;auto" json:"id"`
	TipoDocumentoDigitalId    int64                       `orm:"column(tipo_documento_digital_id)" json:"tipo_documento_digital_id"`
	EstadoDocumentoId         int64                       `orm:"column(estado_documento_id)" json:"estado_documento_id"`
	CodigoEstudiante          int64                       `orm:"column(codigo_estudiante)" json:"codigo_estudiante"`
	FacultadId                *int64                      `orm:"column(facultad_id);null" json:"facultad_id,omitempty"`
	ProgramaAcademicoId       *int64                      `orm:"column(programa_academico_id);null" json:"programa_academico_id,omitempty"`
	PeriodoId                 *int64                      `orm:"column(periodo_id);null" json:"periodo_id,omitempty"`
	Vigencia                  *int                        `orm:"column(vigencia);null" json:"vigencia,omitempty"`
	NombreEstudiante          string                      `orm:"column(nombre_estudiante);size(250);null" json:"nombre_estudiante,omitempty"`
	TipoDocumentoEstudiante   string                      `orm:"column(tipo_documento_estudiante);size(20);null" json:"tipo_documento_estudiante,omitempty"`
	NumeroDocumentoEstudiante string                      `orm:"column(numero_documento_estudiante);size(50);null" json:"numero_documento_estudiante,omitempty"`
	MunicipioExpedicion       string                      `orm:"column(municipio_expedicion);size(150);null" json:"municipio_expedicion,omitempty"`
	TituloOtorgado            string                      `orm:"column(titulo_otorgado);size(250);null" json:"titulo_otorgado,omitempty"`
	UUIDDocumento             *string                     `orm:"column(uuid_documento);null;unique" json:"uuid_documento,omitempty"`
	Activo                    bool                        `orm:"column(activo)" json:"activo"`
	FechaCreacion             time.Time                   `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion         time.Time                   `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
	Historicos                []*HistoricoEstadoDocumento `orm:"reverse(many)" json:"-"`
	Diploma                   *DiplomaDigital             `orm:"reverse(one)" json:"-"`
	Firmas                    []*FirmaDocumento           `orm:"reverse(many)" json:"-"`
}

func (t *DocumentoDigital) TableName() string {
	return "documento_digital"
}

func (t *DocumentoDigital) TableIndex() [][]string {
	return [][]string{
		{"TipoDocumentoDigitalId"},
		{"EstadoDocumentoId"},
		{"CodigoEstudiante"},
		{"FacultadId"},
		{"ProgramaAcademicoId"},
		{"PeriodoId"},
		{"Vigencia"},
		{"NumeroDocumentoEstudiante"},
	}
}

func (t *DocumentoDigital) TableUnique() [][]string {
	return [][]string{
		{"UUIDDocumento"},
	}
}
