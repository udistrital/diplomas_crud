package models

import "github.com/beego/beego/v2/client/orm"

const defaultDriver = "postgres"

func MustInitORM() {
	orm.RegisterDriver(defaultDriver, orm.DRPostgres)
	orm.RegisterModel(
		new(DocumentoDigital),
		new(ParametrizacionRegistroFacultad),
		new(LibroRegistro),
		new(FolioRegistro),
		new(ControlRegistroFacultad),
		new(DiplomaDigital),
		new(RegistroGrado),
		new(CeremoniaGrado),
		new(CeremoniaFacultadPrograma),
		new(CeremoniaEstudiante),
		new(FirmaFirmante),
		new(FirmaDocumento),
		new(HistoricoEstadoDocumento),
	)
}
