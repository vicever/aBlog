package core

import "github.com/codegangsta/cli"

type CoreCmd struct {
	*cli.App
}

// new command cli from vars
func NewCmd(vars CoreVars) *CoreCmd {
	app := cli.NewApp()
	app.Name = vars.Name
	app.Version = vars.Version
	app.Author = vars.Author
	app.Usage = vars.Name
	app.Commands = []cli.Command{}
	return &CoreCmd{app}
}

// add namned command
func (cmd *CoreCmd) AddCommand(c ...cli.Command) {
	cmd.Commands = append(cmd.Commands, c...)
}
