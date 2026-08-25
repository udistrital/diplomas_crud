package routers

import (
	"github.com/astaxie/beego"

	"github.com/udistrital/diplomas_crud/controllers"
)

func init() {
	api := beego.NewNamespace("/v1",
		beego.NSNamespace("/documento_digital",
			beego.NSRouter("/", &controllers.DocumentoDigitalController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/:id", &controllers.DocumentoDigitalController{}, "get:GetOne;put:Put;delete:Delete"),
			beego.NSRouter("/:id/estado", &controllers.DocumentoDigitalController{}, "post:CambiarEstado"),
			beego.NSRouter("/:id/uuid", &controllers.DocumentoDigitalController{}, "put:ActualizarUUID"),
			beego.NSRouter("/:id/firma", &controllers.DocumentoDigitalController{}, "post:RegistrarFirma"),
			beego.NSRouter("/:id/diploma", &controllers.DocumentoDigitalController{}, "post:CrearDiploma"),
		),
		beego.NSNamespace("/diploma_digital",
			beego.NSRouter("/", &controllers.DiplomaDigitalController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/:id", &controllers.DiplomaDigitalController{}, "get:GetOne;put:Put;delete:Delete"),
		),
		beego.NSNamespace("/firma_firmante",
			beego.NSRouter("/", &controllers.FirmaFirmanteController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/:id", &controllers.FirmaFirmanteController{}, "get:GetOne;put:Put;delete:Delete"),
		),
		beego.NSNamespace("/firma_documento",
			beego.NSRouter("/", &controllers.FirmaDocumentoController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/:id", &controllers.FirmaDocumentoController{}, "get:GetOne;put:Put;delete:Delete"),
		),
		beego.NSNamespace("/control_consecutivo_facultad_vigencia",
			beego.NSRouter("/", &controllers.ControlConsecutivoFacultadVigenciaController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/:id", &controllers.ControlConsecutivoFacultadVigenciaController{}, "get:GetOne;put:Put;delete:Delete"),
		),
		beego.NSNamespace("/historico_estado_documento",
			beego.NSRouter("/", &controllers.HistoricoEstadoDocumentoController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/:id", &controllers.HistoricoEstadoDocumentoController{}, "get:GetOne;put:Put;delete:Delete"),
		),
	)

	beego.AddNamespace(api)
}
