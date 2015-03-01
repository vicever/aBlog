package core

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

type CoreDb struct {
	*xorm.Engine
}

func NewCoreDb(config CoreConfig) *CoreDb {
	engine, err := xorm.NewEngine(config.DbDriver, config.DbSQLiteFile)
	if err != nil {
		panic(err)
	}
	if err = engine.Ping(); err != nil {
		panic(err)
	}
	return &CoreDb{engine}
}
