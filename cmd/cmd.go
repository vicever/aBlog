package cmd

import "github.com/fuxiaohei/aBlog/core"

func Init() {
	core.Cmd.AddCommand(
		installCmd,
		servCmd,
	)
}
