package core

import "github.com/codegangsta/cli"

type coreCmd struct {
	*cli.App
}

func newCoreCmd() *coreCmd {
	app := cli.NewApp()
	app.Name = Vars.Name
	app.Usage = Vars.Description
	app.Author = Vars.Author
	app.Email = Vars.AuthorEmail
	app.Version = Vars.Version + " " + Vars.VersionStatus
	app.Commands = []cli.Command{}
	return &coreCmd{app}
}
