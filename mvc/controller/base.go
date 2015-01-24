package controller

import (
	"ablog/core"
	"github.com/tango-contrib/renders"
	"path/filepath"
)

type BaseController struct {
	ViewVars renders.T // view data map

	renders.Renderer
}

func (bc *BaseController) Assign(key string, value interface{}) {
	if bc.ViewVars == nil {
		bc.ViewVars = make(renders.T)
	}
	bc.ViewVars[key] = value
}

func (bc *BaseController) RenderAdmin(tpl string) {
	if err := bc.Renderer.Render(filepath.Join("admin", tpl), bc.ViewVars); err != nil {
		panic(err)
	}
}

func Register() {
	core.Web.Use(AuthHandler())

	// admin profile and password
	core.Web.Any("/admin/profile", new(AdminProfileController))
	core.Web.Post("/admin/profile/password", new(AdminPasswordController))

	// admin article
	core.Web.Any("/admin/article/write", new(ArticleWriteController))

	core.Web.Any("/login", new(LoginController))
	core.Web.Get("/logout", new(LogoutController))
}
