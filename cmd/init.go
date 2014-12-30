package cmd

import "github.com/codegangsta/cli"

var InitCommand cli.Command = cli.Command{
	Name:   "init",
	Usage:  "init blog for first run",
	Action: initFunction,
}

func initFunction(ctx *cli.Context) {

}
