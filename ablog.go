package main

import (
	"github.com/fuxiaohei/ablog/cmd"
	"github.com/fuxiaohei/ablog/sys"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	sys.Init()

	cmd.Init()

	sys.Run()
}
