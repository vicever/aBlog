package cmd

import "github.com/codegangsta/cli"

var backupCommand = cli.Command{
	Name:   "backup",
	Usage:  "Backup static resources and users' contents to zip archive",
	Action: Install,
}
