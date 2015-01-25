package cmd

import (
	"ablog/core"
	"github.com/codegangsta/cli"
)

// rewrite cli app template

func init() {
	cli.AppHelpTemplate = appHelpTemplate
	core.Cmd.Commands = []cli.Command{
		InstallCommand,
		ServerCommand,
		TestCommand,
	}
}

var appHelpTemplate string = `{{.Name}} {{.Version}} - {{.Usage}}

usage:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
`
