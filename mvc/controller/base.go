package controller

import (
	"ablog/core"
	"github.com/tango-contrib/renders"
	"path/filepath"
)

type AdminRender struct {
	renders.Renderer
}

func (a *AdminRender) Render(tpl string, v interface{}) {
	if err := a.Renderer.Render(filepath.Join("admin", tpl), v); err != nil {
		panic(err)
	}
}

func Register() {
	core.Web.Use(AuthHandler())

	core.Web.Any("/login", new(LoginController))
	core.Web.Get("/logout", new(LogoutController))
}
