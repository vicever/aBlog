package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/core"
	"github.com/fuxiaohei/ablog/log"
	"github.com/lunny/tango"
)

var serverCommand = cli.Command{
	Name:   "server",
	Usage:  "Run own webserver to render and display public",
	Action: Server,
}

func Server(_ *cli.Context) {
	if !core.Config.HasFile() {
		log.Fatal("please install ablog first!")
	}
	// read config file
	if err := core.Config.ReadFile(); err != nil {
		log.Fatal("start server fail : %v", err)
	}

	var (
		server *tango.Tango = core.Tango
		addr   string       = fmt.Sprintf("%s:%d", core.Config.Server.Addr, core.Config.Server.Port)
	)

	server.Run(addr)

}
