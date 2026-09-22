package models

import "github.com/beego/beego/v2/client/orm"

const defaultDriver = "postgres"

func MustInitORM() {
	orm.RegisterDriver(defaultDriver, orm.DRPostgres)
	orm.RegisterModel(
		new(DocumentoDigital),
		new(ControlConsecutivoFacultadVigencia),
		new(DiplomaDigital),
		new(FirmaFirmante),
		new(FirmaDocumento),
		new(HistoricoEstadoDocumento),
	)
}
