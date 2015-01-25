package cmd

import (
	"ablog/core"
	"fmt"
	"github.com/codegangsta/cli"
	"time"
)

type User struct {
	Id           int64  `model:"pk"`
	Name         string `model:"unique"`
	Password     string
	PasswordSalt string

	NickName string
	Email    string `model:"unique"`
	Url      string
	Bio      string

	CreateTime    time.Time
	LastLoginTime time.Time

	Role   string `model:"index"`
	Social map[string]string
}

var TestCommand cli.Command = cli.Command{
	Name:  "test",
	Usage: "ablog functionalities testing",
	Action: func(ctx *cli.Context) {
		fmt.Println(core.Model.Save(&User{
			Id:    66,
			Email: "test@qq.com",
			Name:  "test",
			Role:  "admin",
		}))
		fmt.Println(core.Model.Register(new(User)))

	},
}
