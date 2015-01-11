package admin

import "github.com/fuxiaohei/ablog/mod/base"

type DashboardAction struct {
	base.AdminRenders
}

func (d *DashboardAction) Get() {
	d.Render("dashboard.html", nil)
}
