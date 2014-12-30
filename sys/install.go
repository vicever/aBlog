package sys

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/core"
	"github.com/fuxiaohei/ablog/mod/user"
	"time"
)

var CmdInit cli.Command = cli.Command{
	Name:   "init",
	Usage:  "init blog for first run",
	Action: initFunction,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name", Usage: "admin username", Value: "admin"},
	},
}

func initFunction(ctx *cli.Context) {
	if IsInit() {
		core.Fatal("[Init] init locked, please remove 'data' directory to re-install")
	}
	// assign init time, use for id generator
	core.Vars.Status.InitTime = time.Now().Unix()

	// stop event async
	core.Event.EnableAsync = false

	// init admin
	initAdminUser(ctx.String("name"))

	// init lock file
	core.Db.SetInt64(core.Vars.InitLock, core.Vars.Status.InitTime)

	core.Info("[Init] init success !!!")
}

func IsInit() bool {
	return core.Db.Exist(core.Vars.InitLock)
}

func initAdminUser(name string) {
	_, err := user.Create(name, "admin@example.com", "admin", user.ROLE_ADMIN)
	if err != nil {
		core.Fatal("[Init] %v", err)
		return
	}
	core.Debug("[Init] create admin user '%s'", name)
}
