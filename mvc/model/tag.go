package model

import "github.com/fuxiaohei/aBlog/core"

type Tag struct {
	Id   int64 `xorm:"pk autoincr"`
	Uid  int64
	Name string
	Slug string `xorm:"unique(tag-slug)"`
}

type TagArticle struct {
	ArticleId int64
	TagId     int64
}

func GetTagBy(column string, v interface{}) *Tag {
	t := new(Tag)
	if _, err := core.Db.Where(column+" = ?", v).Get(t); err != nil {
		return nil
	}
	if t.Id == 0 {
		return nil
	}
	return t
}

func GetTagsByUser(uid int64) []*Tag {
	ts := make([]*Tag, 0)
	if err := core.Db.Where("uid = ?", uid).Find(ts); err != nil {
		return nil
	}
	return ts
}

func DeleteTag(tid int64) error {
	if _, err := core.Db.Where("tag_id = ?", tid).Delete(new(TagArticle)); err != nil {
		return err
	}
	if _, err := core.Db.Where("id = ?", tid).Delete(new(Tag)); err != nil {
		return err
	}
	return nil
}
