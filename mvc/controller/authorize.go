package controller

import (
	"ablog/mvc/model"
	"github.com/lunny/tango"
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

	BaseController
}

func (lc *LoginController) Get() {
	// check cookie existing
	if lc.isLogged() {
		lc.Redirect("/admin/articles", 302)
		return
	}

	lc.Assign("XsrfFormHtml", lc.XsrfFormHtml)
	lc.RenderAdmin("login.html")
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

	// delete cookie
	if cookie != nil {
		lgt.Cookies().Del("authorize")
	}

	lgt.Redirect("/login", 302)
}

/*
===== authorize interface
*/

type Auther interface {
	SetAuthUser(*model.User)
	AuthFailRedirect() string
}

type AdminBaseController struct {
	AuthUser *model.User
	BaseController
}

func (abc *AdminBaseController) SetAuthUser(u *model.User) {
	abc.AuthUser = u
	abc.Assign("AuthUser", u)
}

func (abc *AdminBaseController) AuthFailRedirect() string {
	return "/login?error=fail-auth"
}

func AuthHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		// only apply to AuthControllerInterface
		action, ok := ctx.Action().(Auther)
		if !ok {
			ctx.Next()
			return
		}
		// get cookie, then get user
		if cookie := ctx.Cookies().Get("authorize"); cookie != nil {
			if token, _ := model.GetToken(cookie.Value); token != nil {
				if user, _ := token.GetUser(); user != nil {
					action.SetAuthUser(user)
					ctx.Next() // remember call next
					return
				}
			}
		}

		// check redirect for auth-failure
		if redirectUrl := action.AuthFailRedirect(); redirectUrl != "" {
			ctx.Redirect(redirectUrl, 302)
			return
		}

		ctx.Next()
	}
}
