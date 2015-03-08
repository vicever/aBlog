package mvc

import (
	"github.com/fuxiaohei/aBlog/core"
	"github.com/fuxiaohei/aBlog/mvc/controller"
	"github.com/fuxiaohei/aBlog/mvc/controller/api"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

func Init() {

	// static files middleware
	core.Server.Use(tango.Static(tango.StaticOptions{
		Prefix:     "theme",
		RootPath:   "./theme",
		ListDir:    false,
		FilterExts: []string{".css", ".js", ".jpg", ".png", ".gif", ".eot", ".otf", ".svg", ".ttf", ".woff", ".woff2"},
	}))

	// render middleware
	core.Server.Use(renders.New(renders.Options{
		Reload:    true, // just in dev
		Directory: "theme",
	}))

	// api group
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

	// page router
	core.Server.Any("/login", new(controller.LoginController))
	core.Server.Get("/", new(controller.HomeController))
}
