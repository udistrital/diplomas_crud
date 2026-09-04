package models

import (
	"time"
)

type HistoricoEstadoDocumento struct {
	Id                 int64             `orm:"column(id);pk;auto" json:"id"`
	DocumentoDigital   *DocumentoDigital `orm:"column(documento_digital_id);rel(fk)" json:"documento_digital"`
	EstadoAnteriorId   *int64            `orm:"column(estado_anterior_id);null" json:"estado_anterior_id,omitempty"`
	EstadoNuevoId      int64             `orm:"column(estado_nuevo_id)" json:"estado_nuevo_id"`
	DocumentoIdentidad *int64            `orm:"column(documento_identidad);null" json:"documento_identidad,omitempty"`
	RolActorId         *int64            `orm:"column(rol_actor_id);null" json:"rol_actor_id,omitempty"`
	FirmaDocumento     *FirmaDocumento   `orm:"column(firma_documento_id);rel(fk);null" json:"firma_documento,omitempty"`
	Observacion        string            `orm:"column(observacion);size(500);null" json:"observacion,omitempty"`
	Activo             bool              `orm:"column(activo)" json:"activo"`
	FechaCreacion      time.Time         `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion  time.Time         `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
}

func (t *HistoricoEstadoDocumento) TableName() string {
	return "historico_estado_documento"
}

func (t *HistoricoEstadoDocumento) TableIndex() [][]string {
	return [][]string{
		{"DocumentoDigital"},
		{"EstadoNuevoId"},
		{"DocumentoIdentidad"},
	}
}
