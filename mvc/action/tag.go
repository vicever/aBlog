package action

import (
	"github.com/fuxiaohei/aBlog/core"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"strconv"
)

const ERR_TAG_NO_FOUND = 3101
const ERR_TAG_NOT_OWNER = 3102

func init() {
	errorMap[ERR_TAG_NO_FOUND] = "no-tag"
	errorMap[ERR_TAG_NOT_OWNER] = "not-tag-owner"
}

// get tags
// input "uid,cid,slug"
// ouput "tag:*model.Tag"
func GetTag(params ActionParam) ActionResult {
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	tid, _ := strconv.ParseInt(params["tid"], 10, 64)
	if uid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	// get by cid
	if tid > 0 {
		t := model.GetTagBy("id", tid)
		if t == nil {
			return NewResultError(ERR_CATEGORY_NO_FOUND)
		}
		if t.Uid != uid {
			return NewResultError(ERR_CATEGORY_NOT_OWNER)
		}
		return NewResult(map[string]interface{}{
			"tag": t,
		})
	}
	// get slug
	slug := params["slug"]
	if slug != "" {
		t := model.GetTagBy("slug", slug)
		if t == nil {
			return NewResultError(ERR_CATEGORY_NO_FOUND)
		}
		if t.Uid != uid {
			return NewResultError(ERR_CATEGORY_NOT_OWNER)
		}
		return NewResult(map[string]interface{}{
			"tag": t,
		})
	}
	return NewResultError(ERR_INVALID_PARAMS)
}

// get tags
// input "uid"
// output "tags:[]*model.Tag"
func GetTags(params ActionParam) ActionResult {
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	if uid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	tags := model.GetTagsByUser(uid)
	return NewResult(map[string]interface{}{
		"tags": tags,
	})
}

// create new category
// input "uid,name,slug"
// output "tag:*model.Tag"
func CreateTag(params ActionParam) ActionResult {
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	name := params["name"]
	slug := params["slug"]
	if name == "" || slug == "" || uid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	t := &model.Tag{
		Uid:  uid,
		Name: name,
		Slug: slug,
	}
	if _, err := core.Db.Insert(t); err != nil {
		return NewSystemError(err)
	}
	return NewResult(map[string]interface{}{
		"tag": t,
	})
}

// update tag
// input "uid,tid,name,slug"
// output "tag:*model.Tag"
func UpdateTag(params ActionParam) ActionResult {
	tid, _ := strconv.ParseInt(params["tid"], 10, 64)
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	name := params["name"]
	slug := params["slug"]
	if name == "" || slug == "" || uid == 0 || tid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	t := &model.Tag{
		Name: name,
		Slug: slug,
	}
	if _, err := core.Db.Cols("name,slug").Where("id = ? AND uid = ?", tid, uid).Update(t); err != nil {
		return NewSystemError(err)
	}
	return NewResult(map[string]interface{}{
		"tag": t,
	})
}

// remove tag
// input "uid,tid"
// output "tag:*model.Tag"
func RemoveTag(params ActionParam) ActionResult {
	tid, _ := strconv.ParseInt(params["tid"], 10, 64)
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	if uid == 0 || tid == 0 {
		return NewResultError(ERR_INVALID_PARAMS)
	}
	t := model.GetTagBy("id", tid)
	if t == nil {
		return NewResultError(ERR_CATEGORY_NO_FOUND)
	}
	if t.Uid != uid {
		return NewResultError(ERR_CATEGORY_NOT_OWNER)
	}
	if err := model.DeleteTag(tid); err != nil {
		return NewSystemError(err)
	}
	return NewResult(map[string]interface{}{
		"tag": t,
	})
}
