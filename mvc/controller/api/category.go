package api

import (
	"github.com/fuxiaohei/aBlog/mvc/action"
	"strconv"
)

type CategoryController struct {
	AuthBase
}

func (cc *CategoryController) Get() {
	// need auth
	b, t := cc.IsAuthorized()
	if !b {
		cc.ResponseWriter.WriteHeader(401)
		return
	}

	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["cid"] = cc.Req().FormValue("cid")
	params["slug"] = cc.Req().FormValue("slug")
	// use cid before
	cc.ServeJson(action.Call(action.GetCategory, params))
}

func (cc *CategoryController) Post() {
	// need auth
	b, t := cc.IsAuthorized()
	if !b {
		cc.ResponseWriter.WriteHeader(401)
		return
	}

	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["cid"] = cc.Req().FormValue("cid")
	params["name"] = cc.Req().FormValue("name")
	params["slug"] = cc.Req().FormValue("slug")
	params["desc"] = cc.Req().FormValue("desc")
	// no cid, to create
	if params["cid"] == "" {
		cc.ServeJson(action.Call(action.CreateCategory, params))
		return
	}
	// there is cid, to update
	cc.ServeJson(action.Call(action.UpdateCategory, params))
}

func (cc *CategoryController) Delete() {
	// need auth
	b, t := cc.IsAuthorized()
	if !b {
		cc.ResponseWriter.WriteHeader(401)
		return
	}

	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["cid"] = cc.Req().FormValue("cid")
	cc.ServeJson(action.Call(action.RemoveCategory, params))
}

type CategoriesController struct {
	AuthBase
}

func (cc *CategoriesController) Get() {
	// need auth
	b, t := cc.IsAuthorized()
	if !b {
		cc.ResponseWriter.WriteHeader(401)
		return
	}
	params := make(map[string]string)
	params["uid"] = strconv.FormatInt(t.Uid, 64)
	params["order"] = cc.Req().FormValue("order")
	cc.ServeJson(action.Call(action.GetCategories, params))
}
