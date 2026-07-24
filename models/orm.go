package models

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	defaultAlias  = "default"
	defaultDriver = "postgres"
)

func MustInitORM() {
	if err := LoadDotEnv(".env"); err != nil {
		panic(fmt.Errorf("load .env: %w", err))
	}
	SyncEnvToAppConfig()

	orm.RegisterDriver(defaultDriver, orm.DRPostgres)
	orm.RegisterModel(new(DocumentoDigital), new(HistoricoEstadoDocumento))
}

func buildDSN() string {
	host := getConfig("DIPLOMAS_CRUD_PGHOST", "PGhost", "127.0.0.1")
	port := getConfig("DIPLOMAS_CRUD_PGPORT", "PGport", "5432")
	user := getConfig("DIPLOMAS_CRUD_PGUSER", "PGuser", "postgres")
	pass := getConfig("DIPLOMAS_CRUD_PGPASS", "PGpass", "postgres")
	name := getConfig("DIPLOMAS_CRUD_PGDB", "PGdb", "udistrital_academica")
	schema := getConfig("DIPLOMAS_CRUD_PGSCHEMA", "PGschema", "diplomas")
	sslMode := getConfig("DIPLOMAS_CRUD_PGSSLMODE", "PGsslmode", "disable")
	timeZone := getConfig("DIPLOMAS_CRUD_PGTIMEZONE", "PGtimezone", "America/Bogota")

	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s search_path=%s,public",
		user,
		pass,
		name,
		host,
		port,
		sslMode,
		timeZone,
		schema,
	)
}

func getConfig(envKey, confKey, fallback string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}

	if value := beego.AppConfig.String(confKey); value != "" {
		return value
	}

	return fallback
}
