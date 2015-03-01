package action

import (
	"github.com/fuxiaohei/aBlog/core"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"strconv"
	"time"
)

const ERR_CATEGORY_NO_FOUND = 3001
const ERR_CATEGORY_NOT_OWNER = 3002

func init() {
	errorMap[ERR_CATEGORY_NO_FOUND] = "no-category"
	errorMap[ERR_CATEGORY_NOT_OWNER] = "not-category-owner"
}

func GetCategory(params ActionParam) ActionResult {
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	cid, _ := strconv.ParseInt(params["cid"], 10, 64)
	if uid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	// get by cid
	if cid > 0 {
		c := model.GetCategoryBy("id", cid)
		if c == nil {
			return NewResultError(ERR_CATEGORY_NO_FOUND)
		}
		if c.Uid != uid {
			return NewResultError(ERR_CATEGORY_NOT_OWNER)
		}
		return NewResult(map[string]interface{}{
			"category": c,
		})
	}
	// get slug
	slug := params["slug"]
	if slug != "" {
		c := model.GetCategoryBy("slug", slug)
		if c == nil {
			return NewResultError(ERR_CATEGORY_NO_FOUND)
		}
		if c.Uid != uid {
			return NewResultError(ERR_CATEGORY_NOT_OWNER)
		}
		return NewResult(map[string]interface{}{
			"category": c,
		})
	}
	return NewResultError(ERR_INVALID_PARAMS)
}

func GetCategories(params ActionParam) ActionResult {
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	order := params["order"]
	if uid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	categories := model.GetCategiresByUser(uid, order != "")
	return NewResult(map[string]interface{}{
		"categories": categories,
	})
}

func CreateCategory(params ActionParam) ActionResult {
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	name := params["name"]
	slug := params["slug"]
	desc := params["desc"]
	if name == "" || slug == "" || uid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	if desc == "" {
		desc = name
	}
	c := &model.Category{
		Uid:         uid,
		Name:        name,
		Slug:        slug,
		Description: desc,
		Count:       0,
		CreateTime:  time.Now().Unix(),
	}
	if _, err := core.Db.Insert(c); err != nil {
		return NewSystemError(err)
	}
	return NewResult(map[string]interface{}{
		"category": c,
	})
}

func UpdateCategory(params ActionParam) ActionResult {
	cid, _ := strconv.ParseInt(params["cid"], 10, 64)
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	name := params["name"]
	slug := params["slug"]
	desc := params["desc"]
	if name == "" || slug == "" || uid == 0 || cid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	if desc == "" {
		desc = name
	}
	c := &model.Category{
		Name:        name,
		Slug:        slug,
		Description: desc,
		Count:       0,
		CreateTime:  time.Now().Unix(),
	}
	if _, err := core.Db.Cols("name,slug,description").Where("id = ? AND uid = ?", cid, uid).Update(c); err != nil {
		return NewSystemError(err)
	}
	return NewResult(map[string]interface{}{
		"category": c,
	})
}

func RemoveCategory(params ActionParam) ActionResult {
	cid, _ := strconv.ParseInt(params["cid"], 10, 64)
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	if uid == 0 || cid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	c := model.GetCategoryBy("id", cid)
	if c == nil {
		return NewResultError(ERR_CATEGORY_NO_FOUND)
	}
	if c.Uid != uid {
		return NewResultError(ERR_CATEGORY_NOT_OWNER)
	}
	if err := model.DeleteCategory(cid); err != nil {
		return NewSystemError(err)
	}
	// todo : update article category to 0
	return NewResult(map[string]interface{}{
		"category": c,
	})
}
