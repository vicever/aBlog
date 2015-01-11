package mod

import (
	"github.com/fuxiaohei/ablog/mod/admin"
	"github.com/fuxiaohei/ablog/mod/user"
)

func Init() {
	// init admin module
	admin.Init()

	// init user module
	user.Init()

}
