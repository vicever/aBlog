package admin

import (
	"github.com/fuxiaohei/aBlog/mvc/action"
	"github.com/lunny/tango"
	"strconv"
)

type ProfileController struct {
	tango.Ctx
	BaseController
}

func (pc *ProfileController) Get() {
	pc.AssignAuth()
	pc.Assign("IsUpdate", pc.Req().FormValue("update") == "true")
	pc.Assign("IsPassword", pc.Req().FormValue("password") == "true")
	pc.Assign("ErrorMessage", pc.Req().FormValue("error"))
	pc.Assign("Title", "Profile")
	pc.Render("profile.html")
}

func (pc *ProfileController) Post() {
	params := make(map[string]string)
	params["name"] = pc.Req().FormValue("name")
	params["nick"] = pc.Req().FormValue("nick")
	params["email"] = pc.Req().FormValue("email")
	params["url"] = pc.Req().FormValue("url")
	params["bio"] = pc.Req().FormValue("bio")
	params["uid"] = strconv.FormatInt(pc.AuthUser.Id, 10)

	result := action.Call(action.UpdateUser, params)
	// pc.Assign("IsUpdate", result.Meta.Status)
	if !result.Meta.Status {
		pc.Redirect("/admin/profile?update=false&error=" + strconv.Itoa(result.Meta.ErrorCode))
		return
	}
	pc.Redirect("/admin/profile?update=true")
}

type PasswordController struct {
	tango.Ctx
	BaseController
}

func (pc *PasswordController) Post() {
	params := make(map[string]string)
	params["new"] = pc.Req().FormValue("new")
	params["old"] = pc.Req().FormValue("old")
	params["uid"] = strconv.FormatInt(pc.AuthUser.Id, 10)

	result := action.Call(action.UpdateUserPassword, params)
	if !result.Meta.Status {
		pc.Redirect("/admin/profile?password=false&error=" + strconv.Itoa(result.Meta.ErrorCode))
		return
	}
	pc.Redirect("/admin/profile?password=true")
}
