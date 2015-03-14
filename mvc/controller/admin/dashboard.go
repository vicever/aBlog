package admin

import (
	"github.com/fuxiaohei/aBlog/lib/theme"
	"github.com/tango-contrib/renders"
)

type DashboardController struct {
	AuthController
	renders.Renderer
}

func (dc *DashboardController) Get() {
	//println(dc.AuthUser.Name)
	dc.Render(theme.AdminFile("dashboard.html"))
}
