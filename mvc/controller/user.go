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
	vars := renders.T{
		"AuthUser":         apc.AuthUser,
		"XsrfFormHtml":     apc.XsrfFormHtml,
		"IsProfileUpdate":  0,
		"IsPasswordUpdate": 0,
	}

	// parse profile updating result
	profileResult := apc.Req().FormValue("profile")
	if profileResult == "update" {
		vars["IsProfileUpdate"] = 2
	} else if profileResult != "" {
		vars["IsProfileUpdate"] = 1
	}

	// parse password updating result
	passwordResult := apc.Req().FormValue("password")
	switch passwordResult {
	case "old-error":
		vars["IsOldPasswordError"] = true
	case "confirm-error":
		vars["IsConfirmPasswordError"] = true
	case "update":
		vars["IsPasswordUpdate"] = 2
	default:
		if passwordResult != "" {
			vars["IsPasswordUpdate"] = 1
		}
	}

	apc.Render("profile.html", vars)
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

	AdminAutherController
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
