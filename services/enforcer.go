package services

import (
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"

	"blog/conf"
	"blog/databases"
)

var enforcer *casbin.Enforcer

func init() {
	config := conf.Config
	a, err := gormAdapter.NewAdapterByDB(databases.DefaultDb)
	if err != nil {
		panic(err)
	}
	enforcer, err = casbin.NewEnforcer(config.Server.CasModelPath, a)
	if err != nil {
		panic(err)
	}
	enforcer.EnableLog(true)
}
