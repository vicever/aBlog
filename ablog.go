package main

import (
	"github.com/fuxiaohei/ablog/core"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	core.Run()
}
