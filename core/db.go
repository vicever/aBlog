package core

import (
	"github.com/lunny/nodb"
	"github.com/lunny/nodb/config"
)

var (
	Db *noDbClient
)

type noDbClient struct {
	Db    *nodb.Nodb
	DbSet *nodb.DB
}

func init() {
	ConnectDb()

}

func ConnectDb() {
	cfg := config.NewConfigDefault()
	cfg.DataDir = Vars.DataDirectory

	NoDb, err := nodb.Open(cfg)
	if err != nil {
		panic(err)
	}

	DbSet, err := NoDb.Select(Vars.DbIndex)
	if err != nil {
		panic(err)
	}

	client := &noDbClient{
		Db:    NoDb,
		DbSet: DbSet,
	}

	Db = client
}
