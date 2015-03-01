package api

import (
	"github.com/fuxiaohei/aBlog/mvc/action"
	"strconv"
)

type UserInfoController struct {
	AuthBase
}

func (uc *UserInfoController) Get() {
	// check auth
	b, t := uc.IsAuthorized()
	if !b {
		uc.ResponseWriter.WriteHeader(401)
		return
	}

	// get user by token's owner id
	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["oid"] = params["uid"]
	result := action.Call(action.GetUser, params)
	uc.ServeJson(result)
}

func (uc *UserInfoController) Post() {
	b, t := uc.IsAuthorized()
	if !b {
		uc.ResponseWriter.WriteHeader(401)
		return
	}
	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["name"] = uc.Req().FormValue("name")
	params["nick"] = uc.Req().FormValue("nick")
	params["email"] = uc.Req().FormValue("email")
	params["url"] = uc.Req().FormValue("url")
	params["bio"] = uc.Req().FormValue("bio")
	uc.ServeJson(action.Call(action.UpdateUser, params))
}
