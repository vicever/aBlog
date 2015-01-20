package cmd

import (
	"ablog/core"
	"ablog/mvc/model"
	"github.com/codegangsta/cli"
	"os"
)

var InstallCommand cli.Command = cli.Command{
	Name:  "install",
	Usage: "install ablog engine",
	Action: func(_ *cli.Context) {
		// config file is exist,
		// it means be installed
		if core.Config.HasFile() {
			core.Log.Fatal("ABlog has been installed")
		}

		// write config file
		if err := core.Config.WriteFile(); err != nil {
			core.Log.Fatal("instaill fail : %v", err)
		}

		// prepare db
		core.PrepareDB()

		// create admin user
		createAdminUser()

		core.Log.Info("ABlog is installed successfully !!")
	},
}

func installFailover() {
	core.Config.RemoveFile()
	os.RemoveAll("data")
}

func createAdminUser() {
	user := model.NewUser("admin", "admin", "admin@example.com", model.USER_ROLE_ADMIN)
	if err := user.Save(); err != nil {
		installFailover()
		core.Log.Fatal("create administrator fail")
	}
}
