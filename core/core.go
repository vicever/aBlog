package core

import (
	"os"
)

var (
	Vars   CoreVars    // global vars
	Cmd    *CoreCmd    // command struct
	Config *CoreConfig // configuration
	Db     *CoreDb     // database
	Server *CoreServer
)

func init() {
	Vars = NewVars()
	Cmd = NewCmd(Vars)

	// load config
	Config = NewConfig()
	if Config.IsFiled() {
		// if config file exist, read it
		if err := Config.ReadFile(); err != nil {
			panic(err)
		}
	}

}

func InitDb() {
	Db = NewCoreDb(*Config)
}

func InitServer() {
	Server = NewCoreServer()
}

func StartServer() {
	Server.Run(Config.TcpAddress)
}

// run core
func Run() {
	Cmd.Run(os.Args)
}
