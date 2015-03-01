package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/aBlog/core"
	"github.com/fuxiaohei/aBlog/mvc"
)

var servCmd cli.Command = cli.Command{
	Name:   "serv",
	Usage:  "serv aBlog",
	Action: servCmdFunc,
}

func servCmdFunc(ctx *cli.Context) {
	// if not installed, stop
	if !core.Config.IsFiled() {
		println("installed")
		return
	}

	// connect db
	core.InitDb()

	// init server
	core.InitServer()

	// init mvc
	mvc.Init()

	// start server
	core.StartServer()
}
