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
	tc.Render("taxonomy.html")
}
