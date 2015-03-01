package mvc

import (
	"github.com/fuxiaohei/aBlog/core"
    "github.com/lunny/tango"
    "github.com/fuxiaohei/aBlog/mvc/controller/api"
)

func Init() {
    apiGroup := tango.NewGroup()
    apiGroup.Any("/user/authorize",new(api.AuthorizeController))
    apiGroup.Post("/user/is_authorize",new(api.AuthorizeCheckController))
    apiGroup.Any("/user/info",new(api.UserInfoController))
	core.Server.Group("/api",apiGroup)
}
