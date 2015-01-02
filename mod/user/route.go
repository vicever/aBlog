package user

import (
	"fmt"
	"github.com/lunny/tango"
	"net/http"
	"time"
)

const (
	TOKEN_COOKIE_NAME = "token-cookie"
)

type LoginRoute struct {
	tango.Ctx
}

func (login LoginRoute) Get() string {
	cookie, _ := login.Req().Cookie(TOKEN_COOKIE_NAME)
	fmt.Println(cookie)
	return "user-login"
}

func (login LoginRoute) Post() {
	var (
		username string = login.Req().FormValue("username")
		password string = login.Req().FormValue("password")
		src      string = login.Req().FormValue("src")
	)
	if src == "" {
		src = "web"
	}

	// get user
	u, _ := GetByName(username)
	if !u.IsPassword(password) {
		return
	}

	// create token
	expire := 3600 * 24 * 7
	token, _ := CreateToken(u.Id, expire, src, login.Req().UserAgent())
	if token != nil {
		login.setLoginCookie(token)
	}

	// todo : login returning
	login.ServeJson(map[string]interface{}{
		"ok":   true,
		"user": u.Id,
	})
}

func (login LoginRoute) setLoginCookie(t *Token) {
	cookie := http.Cookie{
		Name:     TOKEN_COOKIE_NAME,
		Value:    t.Value,
		Expires:  t.ExpireTime,
		HttpOnly: true,
	}
	http.SetCookie(login.ResponseWriter, &cookie)
}

type LogoutRoute struct {
	tango.Ctx
}

func (logout LogoutRoute) Get() {
	// clean cookie
	cookie := http.Cookie{
		Name:    TOKEN_COOKIE_NAME,
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(logout.ResponseWriter, &cookie)

	// redirect
	logout.Redirect("/login")
}

//

type AuthUserRoute interface {
	SetAuthUser(*User)
	IsFailRedirect() bool
}

func AuthorizeHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		if ctx.Action() == nil {
			ctx.Next()
			return
		}
		// only support AuthUserRoute implementation
		action, ok := ctx.Action().(AuthUserRoute)
		if !ok {
			ctx.Next()
			return
		}

		// get cookie
		cookie, _ := ctx.Req().Cookie(TOKEN_COOKIE_NAME)
		if cookie == nil {
			if action.IsFailRedirect() {
				ctx.Redirect("/login")
				return
			}
			ctx.Next()
			return
		}

		// get token and user
		token, user, _ := GetToken(cookie.Value)
		if token == nil || token.IsExpired() || user == nil {
			if action.IsFailRedirect() {
				ctx.Redirect("/login")
				return
			}
			ctx.Next()
			return
		}

		// assign action
		action.SetAuthUser(user)
		ctx.Next()
	}
}
