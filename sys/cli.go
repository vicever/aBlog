package sys

import (
	"github.com/codegangsta/cli"
)

var Cli *cli.App

func newCli() *cli.App {
	app := cli.NewApp()
	app.Name = Vars.Name
	app.Usage = Vars.Description
	app.Author = Vars.Author
	app.Version = Vars.Version
	return app
}
