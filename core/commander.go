package core

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/cmd"
	"os"
)

var Commander *cli.App = defaultCommander()

func defaultCommander() *cli.App {
	// todo : system language
	app := cli.NewApp()
	app.Name = Vars.Name
	app.Usage = Vars.Description
	app.Author = Vars.Author
	app.Version = Vars.Version

	// init commands
	app.Commands = []cli.Command{
		cmd.InitCommand,
	}
	return app
}

func Run() {
	Commander.Run(os.Args)
}
