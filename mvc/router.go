package mvc

import (
	"github.com/fuxiaohei/aBlog/core"
	"github.com/fuxiaohei/aBlog/mvc/controller/api"
	"github.com/lunny/tango"
)

func Init() {
	apiGroup := tango.NewGroup()
	apiGroup.Any("/user/authorize", new(api.AuthorizeController))
	apiGroup.Post("/user/is_authorize", new(api.AuthorizeCheckController))
	apiGroup.Any("/user/info", new(api.UserInfoController))
	apiGroup.Post("/user/password", new(api.UserPasswordController))
	apiGroup.Any("/category", new(api.CategoryController))
	apiGroup.Get("/categories", new(api.CategoriesController))
	apiGroup.Any("/tag", new(api.TagController))
	apiGroup.Get("/tags", new(api.TagsController))
	core.Server.Group("/api", apiGroup)
}
