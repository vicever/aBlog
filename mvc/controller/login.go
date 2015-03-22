package controller

import (
	"github.com/fuxiaohei/aBlog/lib/theme"
	"github.com/fuxiaohei/aBlog/mvc/action"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"net/http"
	"strconv"
	"time"
)

type LoginController struct {
	tango.Ctx
	renders.Renderer
}

func (lc *LoginController) checkAuthorize() bool {
	c := lc.Cookies().Get("auth")
	if c == nil {
		// no need to check
		return false
	}
	c2 := lc.Cookies().Get("auth_uid")
	params := make(map[string]string)
	params["uid"] = c2.Value
	params["token"] = c.Value
	result := action.Call(action.IsAuthorized, params)
	return result.Meta.Status // true value means logged
}

func (lc *LoginController) Get() {
	// check cookie
	if lc.checkAuthorize() {
		lc.Redirect("/admin/dashboard", 302)
		return
	}
	lc.Render(theme.AdminFile("login.html"))
}

func (lc *LoginController) Post() {
	params := make(map[string]string)
	params["user"] = lc.Req().FormValue("user")
	params["password"] = lc.Req().FormValue("password")
	params["ip"] = lc.Req().RemoteAddr
	params["mark"] = "web-page"
	// call authorize action
	result := action.Call(action.Authorize, params)
	if !result.Meta.Status {
		// auth fail
		lc.Redirect("/login?err="+strconv.Itoa(result.Meta.ErrorCode), 302)
		return
	}

	// auth succeed, set cookie
	token := result.Data["token"].(*model.Token)
	cookie := &http.Cookie{
		Name:     "auth",
		Value:    token.Value,
		Expires:  time.Unix(token.ExpireTime, 0),
		HttpOnly: true,
	}
	lc.Cookies().Set(cookie)
	cookie = &http.Cookie{
		Name:     "auth_uid",
		Value:    strconv.FormatInt(token.Uid, 10),
		Expires:  time.Unix(token.ExpireTime, 0),
		HttpOnly: true,
	}
	lc.Cookies().Set(cookie)
	lc.Redirect("/admin/dashboard")
}

type LogoutController struct {
	tango.Ctx
}

func (lc *LogoutController) Get() {
	c := lc.Cookies().Get("auth")
	if c != nil {
		// call unauthorize to delete token in database
		c2 := lc.Cookies().Get("auth_uid")
		params := make(map[string]string)
		params["uid"] = c2.Value
		params["token"] = c.Value
		action.Call(action.Unauthorize, params)
		lc.Cookies().Del("auth")
		lc.Cookies().Del("auth_uid")
	}
	lc.Redirect("/login", 302)
}
