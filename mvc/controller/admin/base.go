package admin

import (
	"github.com/fuxiaohei/aBlog/lib/theme"
	"github.com/fuxiaohei/aBlog/mvc/controller/idl"
)

type BaseController struct {
	idl.AuthRedirecter
	idl.ThemeRenderer
}

func (bc *BaseController) AssignAuth() {
	bc.Assign("AuthUser", bc.AuthUser)
	bc.Assign("AuthToken", bc.AuthToken)
}

func (bc *BaseController) Render(file string) {
	bc.ThemeRenderer.Render(theme.AdminFile(file))
}
