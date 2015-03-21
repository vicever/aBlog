package admin

type DashboardController struct {
	BaseController
}

func (dc *DashboardController) Get() {
	dc.AssignAuth()
	dc.Assign("Title", "Dashboard")
	dc.Assign("IsPage_Dashboard", true)
	//println(dc.AuthUser.Name)
	dc.Render("dashboard.html")
}
