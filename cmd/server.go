package cmd

import (
	"ablog/core"
	"ablog/mvc/controller"
	"ablog/mvc/model"
	"github.com/codegangsta/cli"
)

var ServerCommand cli.Command = cli.Command{
	Name:  "server",
	Usage: "run ablog web server",
	Action: func(ctx *cli.Context) {
		// check db and server preparation
		if core.Db == nil || core.Web == nil {
			core.Log.Fatal("did you install ABlog?")
		}

		// init models
		if err := model.Register(); err != nil {
			panic(err)
		}

		// init controllers' rules
		controller.Register()

		// start server
		core.Web.Run()
	},
}
