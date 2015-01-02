package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/mod"
	"github.com/fuxiaohei/ablog/sys"
	"github.com/lunny/tango"
)

var CmdServ cli.Command = cli.Command{
	Name:   "serv",
	Usage:  "run blog http server",
	Action: servFunc,
}

func servFunc(ctx *cli.Context) {
	checkInit()

	// init tango server
	sys.InitTango()

	// init modulars
	mod.Init()

	var t *tango.Tango = sys.Tango

	t.Run(sys.Config.Tango.AddrString)

}
