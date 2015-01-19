package controller

import (
	"ablog/mvc/model"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
	"time"
)

type LoginController struct {
	tango.Ctx
	xsrf.Checker

	AdminRender
}

func (lc *LoginController) Get() {
	lc.Render("login.html", renders.T{
		"XsrfFormHtml": lc.XsrfFormHtml,
	})
}

func (lc *LoginController) Post() {
	user := lc.Req().FormValue("user")
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
	/*cookie := http.Cookie{
		Name:     "authorize",
		Value:    token.Value,
		Expires:  time.Now().Add(time.Duration(expireTime) * time.Second),
		MaxAge:   expireTime,
		HttpOnly: true,
	}*/
	token.Save()

	// update user's login time
	u.LastLoginTime = time.Now()
	u.Save()

	// redirect into article list-page
	lc.Redirect("/admin/articles", 302)
}
