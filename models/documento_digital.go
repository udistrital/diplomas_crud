package models

import (
	"time"
)

type DocumentoDigital struct {
	Id                  int64                       `orm:"column(id);pk;auto" json:"id"`
	TipoDocumentoId     int64                       `orm:"column(tipo_documento_id)" json:"tipo_documento_id"`
	EstadoDocumentoId   int64                       `orm:"column(estado_documento_id)" json:"estado_documento_id"`
	TerceroId           int64                       `orm:"column(tercero_id)" json:"tercero_id"`
	ProgramaAcademicoId *int64                      `orm:"column(programa_academico_id);null" json:"programa_academico_id,omitempty"`
	PeriodoId           *int64                      `orm:"column(periodo_id);null" json:"periodo_id,omitempty"`
	Vigencia            *int                        `orm:"column(vigencia);null" json:"vigencia,omitempty"`
	HashDocumento       string                      `orm:"column(hash_documento);size(128);unique" json:"hash_documento"`
	UUIDVerificacion    string                      `orm:"column(uuid_verificacion);size(100);unique" json:"uuid_verificacion"`
	Metadata            JSONB                       `orm:"column(metadata);type(jsonb);null" json:"metadata,omitempty"`
	Activo              bool                        `orm:"column(activo)" json:"activo"`
	FechaCreacion       time.Time                   `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion   time.Time                   `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
	Historicos          []*HistoricoEstadoDocumento `orm:"reverse(many)" json:"-"`
}

func (t *DocumentoDigital) TableName() string {
	return "documento_digital"
}

func (t *DocumentoDigital) TableIndex() [][]string {
	return [][]string{
		{"TipoDocumentoId"},
		{"EstadoDocumentoId"},
		{"TerceroId"},
		{"ProgramaAcademicoId"},
		{"PeriodoId"},
		{"Vigencia"},
	}
}

func (t *DocumentoDigital) TableUnique() [][]string {
	return [][]string{
		{"HashDocumento"},
		{"UUIDVerificacion"},
	}
}
