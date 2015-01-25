package cmd

import (
	"ablog/core"
	"fmt"
	"github.com/codegangsta/cli"
)

var TestCommand cli.Command = cli.Command{
	Name:  "test",
	Usage: "ablog functionalities testing",
	Action: func(ctx *cli.Context) {
		key := "test-zset"
		core.Db.ZSet(key, 1, []byte("1111"))
		core.Db.ZSet(key, 9, []byte("2222"))
		core.Db.ZSet(key, 3, []byte("3333"))
		core.Db.ZSet(key, 4, []byte("4444"))

		fmt.Println(core.Db.ZPageAsc(key, 2, 2))
	},
}
