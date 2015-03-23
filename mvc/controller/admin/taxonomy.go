package admin

import (
	"github.com/fuxiaohei/aBlog/mvc/action"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"github.com/lunny/tango"
	"strconv"
)

type TaxonomyController struct {
	tango.Ctx
	BaseController
}

func (tc *TaxonomyController) Get() {
	tc.AssignAuth()
	tc.Assign("Title", "Taxonomy")
	tc.Assign("IsPage_Taxonomy", true)
	// get all categories
	result := action.Call(action.GetCategories, map[string]string{
		"uid":   strconv.FormatInt(tc.AuthUser.Id, 10),
		"order": "true",
	})
	if result.Meta.Status {
		tc.Assign("Categories", result.Data["categories"].([]*model.Category))
	}

	// get single category if updating
	if cid, _ := strconv.ParseInt(tc.Req().FormValue("category"), 10, 64); cid > 0 {
		params := make(map[string]string)
		params["uid"] = strconv.FormatInt(tc.AuthUser.Id, 10)
		params["cid"] = tc.Req().FormValue("category")
		result := action.Call(action.GetCategory, params)
		if result.Meta.Status {
			tc.Assign("Category", result.Data["category"].(*model.Category))
		}
	}
	tc.Render("taxonomy.html")
}

func (tc *TaxonomyController) Post() {
	if tc.Req().FormValue("category") == "true" {
		params := make(map[string]string)
		params["name"] = tc.Req().FormValue("name")
		params["slug"] = tc.Req().FormValue("slug")
		params["desc"] = tc.Req().FormValue("desc")
		params["uid"] = strconv.FormatInt(tc.AuthUser.Id, 10)
		params["cid"] = tc.Req().FormValue("cid")
		var result action.ActionResult
		if params["cid"] == "" {
			result = action.Call(action.CreateCategory, params)
		} else {
			result = action.Call(action.UpdateCategory, params)
		}
		if result.Meta.Status {
			tc.Redirect("/admin/taxonomy?category=" + strconv.FormatInt(result.Data["category"].(*model.Category).Id, 10))
			return
		}
		tc.Redirect("/admin/taxonomy?category=0")
		return
	}
}
