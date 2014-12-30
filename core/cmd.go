package core

import (
	"github.com/codegangsta/cli"
	"os"
)

var Cmd *cli.App = defaultCommander()

func defaultCommander() *cli.App {
	// todo : system language
	app := cli.NewApp()
	app.Name = Vars.Name
	app.Usage = Vars.Description
	app.Author = Vars.Author
	app.Version = Vars.Version
	return app
}

func Run() {
	Cmd.Run(os.Args)
}
