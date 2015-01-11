package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/ablog/core"
)

func init() {
	cli.AppHelpTemplate = `{{.Name}} {{.Version}} - {{.Usage}}

commands:{{range .Commands}}
  {{.Name}}{{ "\t" }}{{.Usage}}{{end}}

author:
  {{if .Author}}{{.Author}}{{end}}

`
}

func Init() {
	core.Cmd.Commands = []cli.Command{
		installCommand,
		serverCommand,
		backupCommand,
	}
}
