package controller

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
)

var _ Auther = new(AdminProfileController)

type AdminProfileController struct {
	tango.Ctx
	xsrf.Checker

	AdminRender
	AdminAutherController
}

func (apc *AdminProfileController) Get() {
	apc.Render("profile.html", renders.T{
		"AuthUser":     apc.AuthUser,
		"XsrfFormHtml": apc.XsrfFormHtml,
	})
}

func (apc *AdminProfileController) Post() {
	user := apc.AuthUser

	// todo: binding refactor
	user.Name = apc.Req().FormValue("username")
	user.NickName = apc.Req().FormValue("nickname")
	user.Email = apc.Req().FormValue("email")
	user.Url = apc.Req().FormValue("url")
	user.Bio = apc.Req().FormValue("bio")

	if err := user.Save(); err != nil {
		apc.Redirect("/admin/profile?error=save-fail", 302)
		return
	}

	apc.Redirect("/admin/profile?success=update", 302)
}
