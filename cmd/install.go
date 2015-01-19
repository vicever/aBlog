package cmd

import (
	"ablog/core"
	"github.com/codegangsta/cli"
)

var InstallCommand cli.Command = cli.Command{
	Name:  "install",
	Usage: "install ablog engine",
	Action: func(_ *cli.Context) {
		// config file is exist,
		// it means be installed
		if core.Config.HasFile() {
			core.Log.Fatal("ABlog has been installed")
		}

		// write config file
		if err := core.Config.WriteFile(); err != nil {
			core.Log.Fatal("instaill fail : %v", err)
		}
	},
}
