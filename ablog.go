package main

import (
	_ "ablog/cmd"
	"ablog/core"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	core.Run()
}
