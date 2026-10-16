package models

import "github.com/beego/beego/v2/client/orm"

const defaultDriver = "postgres"

func MustInitORM() {
	orm.RegisterDriver(defaultDriver, orm.DRPostgres)
}
