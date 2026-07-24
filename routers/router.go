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
		),
		beego.NSNamespace("/historico_estado_documento",
			beego.NSRouter("/", &controllers.HistoricoEstadoDocumentoController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/:id", &controllers.HistoricoEstadoDocumentoController{}, "get:GetOne;put:Put;delete:Delete"),
		),
	)

	beego.AddNamespace(api)
}
