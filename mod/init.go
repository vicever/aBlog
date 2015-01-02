package mod

import (
	"github.com/fuxiaohei/ablog/mod/setting"
	"github.com/fuxiaohei/ablog/mod/user"
	"github.com/fuxiaohei/ablog/sys"
	"github.com/lunny/tango"
)

func Init() {

	var t *tango.Tango = sys.Tango

	// user authorize handler for AuthUserHandler interface implement
	t.Use(user.AuthorizeHandler())

	t.Any("/login", new(user.LoginRoute))
	t.Get("/logout", new(user.LogoutRoute))

	t.Get("/setting/general", new(setting.GeneralRoute))

}
