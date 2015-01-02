package setting

import (
	"github.com/fuxiaohei/ablog/mod/user"
	"github.com/lunny/tango"
)

var _ user.AuthUserRoute = new(GeneralRoute)

type GeneralRoute struct {
	AuthUser *user.User
	tango.Ctx
}

func (s *GeneralRoute) SetAuthUser(u *user.User) {
	s.AuthUser = u
}

func (s *GeneralRoute) IsFailRedirect() bool {
	return true
}

func (s GeneralRoute) Get() {
	s.Result = "user-setting-general"
}

func (s GeneralRoute) Post() {
	s.Result = "user-setting-post"
}
