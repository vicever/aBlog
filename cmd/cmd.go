package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/sys"
)

func Init() {
	sys.Cli.Commands = []cli.Command{
		CmdInit,
		CmdServ,
	}
}
