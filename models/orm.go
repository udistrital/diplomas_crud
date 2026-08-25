package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

const defaultDriver = "postgres"

func MustInitORM() {
	if err := LoadDotEnv(".env"); err != nil {
		panic(fmt.Errorf("load .env: %w", err))
	}
	SyncEnvToAppConfig()

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
