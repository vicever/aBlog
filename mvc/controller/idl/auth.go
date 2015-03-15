package idl

import (
	"github.com/fuxiaohei/aBlog/mvc/action"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"github.com/lunny/tango"
)

type IAuth interface {
	SetUser(*model.User, *model.Token) // set user and token to controller fields
	FailRedirect() string              // fail redirect
}

type AuthRedirecter struct {
	AuthUser  *model.User
	AuthToken *model.Token
}

func (ac *AuthRedirecter) SetUser(u *model.User, t *model.Token) {
	ac.AuthUser = u
	ac.AuthToken = t
}

func (ac *AuthRedirecter) FailRedirect() string {
	return "/login"
}

type AuthNoRedirecter struct {
	AuthUser  *model.User
	AuthToken *model.Token
}

func (ac *AuthNoRedirecter) SetUser(u *model.User, t *model.Token) {
	ac.AuthUser = u
	ac.AuthToken = t
}

func (ac *AuthNoRedirecter) FailRedirect() string {
	return ""
}

func AuthHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		controller, ok := ctx.Action().(IAuth)
		// apply to IAuthController
		if !ok {
			ctx.Next()
			return
		}
		// no auth cookie, ignore
		c := ctx.Cookies().Get("auth")
		if c != nil {
			c2 := ctx.Cookies().Get("auth_uid")
			params := make(map[string]string)
			params["uid"] = c2.Value
			params["token"] = c.Value
			result := action.Call(action.IsAuthorized, params)

			if result.Meta.Status {
				// auth success
				resultMap := result.Data.(map[string]interface{})
				token := resultMap["token"].(*model.Token)
				user := model.GetUserBy("id", token.Uid)
				controller.SetUser(user, token)
				ctx.Next()
				return
			}

		}

		// auth fail,go redirect
		if url := controller.FailRedirect(); url != "" {
			ctx.Redirect(url, 302)
			return
		}
		ctx.Next()
		return
	}

}
