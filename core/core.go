package core

import "os"

// global values
var (
	Vars   *coreVars
	Config *coreConfig
	Db     *coreDb
	Log    *coreLogger
	Cmd    *coreCmd
	Web    *coreWeb
)

func init() {
	// init global vars
	Vars = newCoreVars()

	// init config and load config file if exist
	Config = newCoreConfig()

	// prepare logger
	Log = newCoreLogger(Config.LogFile, "ABlog", true, false)

	// prepare cmd
	Cmd = newCoreCmd()

	// if installed, init nosql, and web server
	if Config.InstallTime > 0 {
		Db = newCoreDb(Config.Db)
		Web = newCoreWeb(Config.Server)
	}
}

func Run() {
	Cmd.Run(os.Args)
}

// prepare db manually
func PrepareDB() {
	Db = newCoreDb(Config.Db)
}

// prepare web server manually
func PrepareWeb() {
	Web = newCoreWeb(Config.Server)
}
