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
		new(model.Tag),
		new(model.TagArticle),
		new(model.Article),
		new(model.Comment),
		new(model.Setting),
		new(model.Media),
	); err != nil {
		panic(err)
	}

	// prepare default data
	prepareDefaultData()
}

func prepareDefaultData() {

	// prepare admin user
	user := model.CreateUser("admin", "admin", "admin@example.com", model.USER_ROLE_ADMIN)

	// prepare default contents
	model.CreateCategory(user.Id, "default", "default", "this is a default category")
	model.CreateTag(user.Id, "blog", "blog")
}
