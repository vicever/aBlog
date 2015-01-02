package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/mod/user"
	"github.com/fuxiaohei/ablog/sys"
	"time"
)

var CmdInit cli.Command = cli.Command{
	Name:   "init",
	Usage:  "init blog",
	Action: initFunc,
}

func initFunc(ctx *cli.Context) {
	if sys.Config.InitTime.Unix() > 0 {
		sys.Error("[Init] the blog has been inited")
		return
	}
	sys.Config.InitTime = time.Now()

	// close event channel
	sys.Event.EnableAsync = false

	// init nodb
	sys.InitNodb()
	sys.Debug("[Init] init nodb success !!")

	// init admin
	initAdmin()
	sys.Debug("[Init] init administrator success !!")

	// init config file
	sys.Config.ToFile()
	sys.Debug("[Init] init config file success !!")

}

func initAdmin() {
	_, err := user.Create("admin", "admin@example.com", "12345678", user.ROLE_ADMIN)
	if err != nil {
		sys.Fatal("[Init] init administrator error : %v", err)
		return
	}
}

func checkInit() {
	if sys.Config.InitTime.Unix() <= 0 {
		sys.Fatal("[Sys] please run `ablog init` to install")
		return
	}
}
