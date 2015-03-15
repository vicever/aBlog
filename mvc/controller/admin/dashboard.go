package admin

import (
	"fmt"
	"github.com/fuxiaohei/aBlog/lib/theme"
	"github.com/fuxiaohei/aBlog/mvc/controller/idl"
	"github.com/tango-contrib/renders"
)

type DashboardController struct {
	idl.AuthRedirecter
	renders.Renderer
}

func (dc *DashboardController) Get() {
	//println(dc.AuthUser.Name)
	fmt.Println(dc.Render(theme.AdminFile("dashboard.html")))
}
