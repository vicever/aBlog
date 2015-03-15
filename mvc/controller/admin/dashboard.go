package admin

import (
	"github.com/fuxiaohei/aBlog/lib/theme"
	"github.com/tango-contrib/renders"
    "fmt"
)

type DashboardController struct {
	AuthController
	renders.Renderer
}

func (dc *DashboardController) Get() {
	//println(dc.AuthUser.Name)
	fmt.Println(dc.Render(theme.AdminFile("dashboard.html")))
}
