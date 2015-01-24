package controller

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
)

var _ Auther = new(AdminProfileController)

type AdminProfileController struct {
	tango.Ctx
	xsrf.Checker

	AdminBaseController
}

func (apc *AdminProfileController) Get() {
	apc.Assign("XsrfFormHtml", apc.XsrfFormHtml)
	apc.Assign("IsProfileUpdate", 0)
	apc.Assign("IsPasswordUpdate", 0)

	// parse profile updating result
	profileResult := apc.Req().FormValue("profile")
	if profileResult == "update" {
		apc.Assign("IsProfileUpdate", 2)
	} else if profileResult != "" {
		apc.Assign("IsProfileUpdate", 1)
	}

	// parse password updating result
	passwordResult := apc.Req().FormValue("password")
	switch passwordResult {
	case "old-error":
		apc.Assign("IsOldPasswordError", true)
	case "confirm-error":
		apc.Assign("IsConfirmPasswordError", true)
	case "update":
		apc.Assign("IsPasswordUpdate", 2)
	default:
		if passwordResult != "" {
			apc.Assign("IsPasswordUpdate", 1)
		}
	}

	apc.RenderAdmin("profile.html")
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
		apc.Redirect("/admin/profile?profile=save-fail")
		return
	}

	apc.Redirect("/admin/profile?profile=update")
}

type AdminPasswordController struct {
	tango.Ctx
	xsrf.Checker

	AdminBaseController
}

func (apc *AdminPasswordController) Post() {
	user := apc.AuthUser

	oldPwd := apc.Req().FormValue("old-pwd")

	if !user.IsSamePassword(oldPwd) {
		apc.Redirect("/admin/profile?password=old-error")
		return
	}

	newPwd := apc.Req().FormValue("new-pwd")
	cfmPwd := apc.Req().FormValue("confirm-pwd")
	if newPwd != cfmPwd {
		apc.Redirect("/admin/profile?password=confirm-error")
		return
	}

	user.SetNewPassword(newPwd)
	if err := user.Save(); err != nil {
		apc.Redirect("/admin/profile?password=save-fail")
		return
	}

	apc.Redirect("/admin/profile?password=update", 302)

}
