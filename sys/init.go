package sys

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/core"
)

func Init() {
	core.Cmd.Commands = []cli.Command{
		CmdInit,
	}
}
