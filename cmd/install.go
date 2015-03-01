package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/aBlog/core"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"time"
)

var installCmd cli.Command = cli.Command{
	Name:   "install",
	Usage:  "install aBlog",
	Action: installCmdFunc,
}

func installCmdFunc(ctx *cli.Context) {
	// if installed, stop
	if core.Config.IsFiled() {
		println("installed")
		return
	}

	// write config file
	core.Config.InstallTime = time.Now().Unix()
	if err := core.Config.WriteFile(); err != nil {
		panic(err)
	}

	// connect and init database
	core.InitDb()
	core.Db.ShowSQL = true

	if err := core.Db.Sync2(
		new(model.User),
		new(model.Token),
		new(model.Category),
		new(model.CategoryArticle),
	); err != nil {
		panic(err)
	}

	// prepare default data
	prepareDefaultData()
}

func prepareDefaultData() {
	model.CreateUser("admin", "admin", "admin@example.com", int8(1))
}
