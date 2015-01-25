package model

import "ablog/core"

func Register() error {
	return core.Model.Register(new(User), new(Token))
}
