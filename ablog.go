package main

import (
	"github.com/fuxiaohei/aBlog/cmd"
	"github.com/fuxiaohei/aBlog/core"
)

func main() {

	cmd.Init()

	core.Run()
}
