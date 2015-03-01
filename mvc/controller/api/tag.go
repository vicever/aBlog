package api

import (
	"github.com/fuxiaohei/aBlog/mvc/action"
	"strconv"
)

type TagController struct {
	AuthBase
}

func (tc *TagController) Get() {
	// need auth
	b, t := tc.IsAuthorized()
	if !b {
		tc.ResponseWriter.WriteHeader(401)
		return
	}

	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["tid"] = tc.Req().FormValue("tid")
	params["slug"] = tc.Req().FormValue("slug")
	// use cid before
	tc.ServeJson(action.Call(action.GetTag, params))
}

func (tc *TagController) Post() {
	// need auth
	b, t := tc.IsAuthorized()
	if !b {
		tc.ResponseWriter.WriteHeader(401)
		return
	}

	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["tid"] = tc.Req().FormValue("tid")
	params["name"] = tc.Req().FormValue("name")
	params["slug"] = tc.Req().FormValue("slug")
	// no cid, to create
	if params["tid"] == "" {
		tc.ServeJson(action.Call(action.CreateTag, params))
		return
	}
	// there is tid, to update
	tc.ServeJson(action.Call(action.UpdateTag, params))
}

func (tc *TagController) Delete() {
	// need auth
	b, t := tc.IsAuthorized()
	if !b {
		tc.ResponseWriter.WriteHeader(401)
		return
	}

	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["tid"] = tc.Req().FormValue("tid")
	tc.ServeJson(action.Call(action.RemoveTag, params))
}

type TagsController struct {
	AuthBase
}

func (tc *TagsController) Get() {
	// need auth
	b, t := tc.IsAuthorized()
	if !b {
		tc.ResponseWriter.WriteHeader(401)
		return
	}
	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	tc.ServeJson(action.Call(action.GetTags, params))
}
