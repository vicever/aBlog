package core

import (
	"github.com/codegangsta/cli"
)

var Cmd *cli.App = newCli()

func newCli() *cli.App {
	app := cli.NewApp()
	app.Name = BLOG_NAME
	app.Usage = BLOG_DESC
	app.Version = BLOG_VERSION
	app.Author = BLOG_AUTHOR
	return app
}
