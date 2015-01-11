package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/core"
	"github.com/fuxiaohei/ablog/log"
)

var installCommand = cli.Command{
	Name:   "install",
	Usage:  "Install new ablog with default settings and data",
	Action: Install,
}

func Install(ctx *cli.Context) {
	if core.Config.HasFile() {
		log.Fatal("ablog had been installed")
	}

	// write config file
	if err := core.Config.WriteFile(); err != nil {
		log.Fatal("ablog install : %v", err)
	}

}
