package controller

import (
	"ablog/mvc/model"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
	"net/http"
	"time"
)

/*
===== login controller
*/

type LoginController struct {
	tango.Ctx
	xsrf.Checker

	AdminRender
}

func (lc *LoginController) Get() {
	// check cookie existing
	if lc.isLogged() {
		lc.Redirect("/admin/articles", 302)
		return
	}

	lc.Render("login.html", renders.T{
		"XsrfFormHtml": lc.XsrfFormHtml,
	})
}

func (lc *LoginController) isLogged() bool {
	// check cookie existing
	cookie := lc.Cookies().Get("authorize")
	if cookie == nil || cookie.Value == "" {
		return false
	}
	// check token available
	token, err := model.GetToken(cookie.Value)
	if token == nil || err != nil {
		return false
	}
	return true
}

func (lc *LoginController) Post() {
	user := lc.Req().FormValue("username")
	password := lc.Req().FormValue("password")
	keep := lc.Req().FormValue("keep")

	// get user by name
	u, err := model.GetUserByName(user)
	if err != nil || u == nil || u.Name != user {
		lc.Redirect("/login?error=no-user", 302)
		return
	}

	// check password
	if !u.IsSamePassword(password) {
		lc.Redirect("/login?error=wrong-password", 302)
		return
	}

	// create token
	var expireTime int64 = 3 * 24 * 3600
	if keep != "" {
		expireTime = 30 * 24 * 3600
	}
	token := model.NewToken(u.Id, expireTime, "web")
	//token.Ip = ".."
	token.UserAgent = lc.Req().UserAgent()

	// save token, save cookie
	cookie := http.Cookie{
		Name:     "authorize",
		Value:    token.Value,
		Expires:  time.Now().Add(time.Duration(expireTime) * time.Second),
		MaxAge:   int(expireTime),
		HttpOnly: true,
	}
	lc.Cookies().Set(&cookie)
	token.Save()

	// update user's login time
	u.LastLoginTime = time.Now()
	u.Save()

	// redirect into article list-page
	lc.Redirect("/admin/articles", 302)
}

/*
==== logout controller
*/

type LogoutController struct {
	tango.Ctx
}

func (lgt *LogoutController) Get() {
	// get this token, remove it
	cookie := lgt.Cookies().Get("authorize")
	if cookie != nil {
		if token, _ := model.GetToken(cookie.Value); token != nil {
			token.Remove()
		}
	}

	// set a expired cookie,
	// to make the using cookie illegal
	if cookie != nil {
		cookie = &http.Cookie{
			Name:     "authorize",
			Value:    "",
			Expires:  time.Now().Add(-3600 * time.Second),
			MaxAge:   -3600,
			HttpOnly: true,
		}
		lgt.Cookies().Set(cookie)
	}

	lgt.Redirect("/login", 302)
}