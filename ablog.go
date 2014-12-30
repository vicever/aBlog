package main

import (
	"github.com/fuxiaohei/ablog/core"
	"github.com/fuxiaohei/ablog/sys"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	sys.Init()

	core.Run()
}
