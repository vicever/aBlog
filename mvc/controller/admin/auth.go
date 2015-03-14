package admin

import (
	"fmt"
	"github.com/fuxiaohei/aBlog/mvc/action"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"github.com/lunny/tango"
)

type IAuthController interface {
	SetUser(*model.User, *model.Token) // set user and token to controller fields
	FailRedirect() string              // fail redirect
}

type AuthController struct {
	AuthUser  *model.User
	AuthToken *model.Token
}

func (ac *AuthController) SetUser(u *model.User, t *model.Token) {
	ac.AuthUser = u
	ac.AuthToken = t
}

func (ac *AuthController) FailRedirect() string {
	return "/login"
}

func AuthHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		controller, ok := ctx.Action().(IAuthController)
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
			fmt.Println(result)

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
