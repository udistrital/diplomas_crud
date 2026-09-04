package models

import "time"

type FirmaFirmante struct {
	Id                 int64             `orm:"column(id);pk;auto" json:"id"`
	DocumentoIdentidad int64             `orm:"column(documento_identidad)" json:"documento_identidad"`
	EnlaceFirma        string            `orm:"column(enlace_firma);unique" json:"enlace_firma"`
	Activo             bool              `orm:"column(activo)" json:"activo"`
	FechaCreacion      time.Time         `orm:"column(fecha_creacion);auto_now_add;type(timestamp)" json:"fecha_creacion"`
	FechaModificacion  time.Time         `orm:"column(fecha_modificacion);auto_now;type(timestamp)" json:"fecha_modificacion"`
	FirmasDocumento    []*FirmaDocumento `orm:"reverse(many)" json:"-"`
}

func (t *FirmaFirmante) TableName() string {
	return "firma_firmante"
}

func (t *FirmaFirmante) TableIndex() [][]string {
	return [][]string{
		{"DocumentoIdentidad"},
		{"EnlaceFirma"},
	}
}

func (t *FirmaFirmante) TableUnique() [][]string {
	return [][]string{
		{"EnlaceFirma"},
	}
}
