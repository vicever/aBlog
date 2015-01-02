package sys

import "os"

func Init() {
	// init config
	Config = newSysConfig()

	// if config is load, means config file is created
	if Config.IsLoad {
		// init nodb
		NoDb = newNodbClient()
	}

	// init command
	Cli = newCli()

}

func InitNodb() {
	NoDb = newNodbClient()
}

func Run() {
	Cli.Run(os.Args)
}
