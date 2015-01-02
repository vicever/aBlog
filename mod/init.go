package mod

import (
	"github.com/fuxiaohei/ablog/mod/user"
	"github.com/fuxiaohei/ablog/sys"
	"github.com/lunny/tango"
)

func Init() {

	var t *tango.Tango = sys.Tango

	t.Any("/login", new(user.LoginRoute))
	t.Get("/logout", new(user.LogoutRoute))

}
