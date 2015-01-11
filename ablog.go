package main

import (
	"github.com/fuxiaohei/ablog/cmd"
	"github.com/fuxiaohei/ablog/core"
	"github.com/fuxiaohei/ablog/mod"
	"runtime"
)

func main() {
	// multiple processes
	runtime.GOMAXPROCS(runtime.NumCPU())

	// init commands
	cmd.Init()

	// init modules
	mod.Init()

	// start engine
	core.Start()

}
