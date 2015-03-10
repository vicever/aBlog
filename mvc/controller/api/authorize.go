package api

import (
	"github.com/fuxiaohei/aBlog/mvc/action"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"github.com/lunny/tango"
)

type AuthorizeController struct {
	tango.Ctx
}

func (ac *AuthorizeController) Post() {
	params := make(map[string]string)
	params["user"] = ac.Req().FormValue("user")
	params["password"] = ac.Req().FormValue("password")
	params["ip"] = ac.Req().RemoteAddr
	params["mark"] = "web-page"
	result := action.Call(action.Authorize, params)
	if result.Meta.ErrorCode == action.ERR_INVALID_PARAMS {
		ac.ResponseWriter.WriteHeader(400)
	}
	ac.ServeJson(result)
}

type AuthorizeCheckController struct {
	tango.Ctx
}

func (ac *AuthorizeCheckController) Post() {
	params := make(map[string]string)
	params["uid"] = ac.Req().FormValue("uid")
	params["token"] = ac.Req().FormValue("token")
	result := action.Call(action.IsAuthorized, params)
	ac.ServeJson(result)
}

// auth base struct
type AuthBase struct {
	tango.Ctx
}

func (ab *AuthBase) IsAuthorized() (bool, *model.Token) {
	params := make(map[string]string)
	params["uid"] = ab.Req().FormValue("uid")
	params["token"] = ab.Req().FormValue("token")
	result := action.Call(action.IsAuthorized, params)
	if result.Meta.Status {
		m := result.Data.(map[string]interface{})
		return true, m["token"].(*model.Token)
	}
	return false, nil
}
