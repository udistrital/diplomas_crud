package models

import "time"

type FirmaDocumento struct {
	Id                int64                       `orm:"column(id);pk;auto" json:"id"`
	DocumentoDigital  *DocumentoDigital           `orm:"column(documento_digital_id);rel(fk)" json:"documento_digital"`
	FirmaFirmante     *FirmaFirmante              `orm:"column(firma_firmante_id);rel(fk);null" json:"firma_firmante,omitempty"`
	RolFirmanteId     int64                       `orm:"column(rol_firmante_id)" json:"rol_firmante_id"`
	TerceroIdFirmante int64                       `orm:"column(tercero_id_firmante)" json:"tercero_id_firmante"`
	Activo            bool                        `orm:"column(activo)" json:"activo"`
	FechaCreacion     time.Time                   `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion time.Time                   `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
	Historicos        []*HistoricoEstadoDocumento `orm:"reverse(many)" json:"-"`
}

func (t *FirmaDocumento) TableName() string {
	return "firma_documento"
}

func (t *FirmaDocumento) TableIndex() [][]string {
	return [][]string{
		{"DocumentoDigital"},
	}
}
