package main

import (
	"os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/lib/pq"
	apistatus "github.com/udistrital/utils_oas/v2/apiStatusLib"
	auditoria "github.com/udistrital/utils_oas/v2/auditoria"
	"github.com/udistrital/utils_oas/v2/customerror"
	"github.com/udistrital/utils_oas/v2/database"
	"github.com/udistrital/utils_oas/v2/security"
	"github.com/udistrital/utils_oas/v2/xray"

	"github.com/udistrital/diplomas_crud/models"
	_ "github.com/udistrital/diplomas_crud/routers"
)

func main() {
	models.MustInitORM()

	conn, err := database.BuildPostgresConnectionString()
	if err != nil {
		logs.Error("error consultando la cadena de conexión: %v", err)
		return
	}

	err = orm.RegisterDataBase("default", "postgres", conn)
	if err != nil {
		logs.Error("error al conectarse a la base de datos: %v", err)
		return
	}

	allowedOrigins := []string{"*.udistrital.edu.co"}
	if beego.BConfig.RunMode == beego.DEV {
		allowedOrigins = []string{"*"}
		orm.Debug = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		if _, err := os.Stat("swagger"); err == nil {
			beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		}
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"DELETE", "GET", "OPTIONS", "PATCH", "POST", "PUT"},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"User-Agent",
			"X-Amzn-Trace-Id",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	err = xray.InitXRay()
	if err != nil {
		logs.Error("error configurando AWS XRay: %v", err)
	}
	apistatus.Init()
	auditoria.InitMiddleware()
	beego.ErrorController(&customerror.CustomErrorController{})
	security.SetSecurityHeaders()

	beego.Run()
}
