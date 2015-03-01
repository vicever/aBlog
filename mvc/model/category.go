package model

import "github.com/fuxiaohei/aBlog/core"

type Category struct {
	Id          int64  `xorm:"pk autoincr" json:"id"`
	Uid         int64  `json:"user_id"`
	Name        string `json:"name"`
	Slug        string `xorm:"unique(category-slug)" json:"slug"`
	Description string `json:"desc"`
	Count       int    `xorm:"default 0" json:"count"`
	CreateTime  int64  `xorm:"created" json:"created"`
}

type CategoryArticle struct {
	ArticleId  int64
	CategoryId int64
}

func GetCategoryBy(column string, v interface{}) *Category {
	c := new(Category)
	if _, err := core.Db.Where(column+" = ?", v).Get(c); err != nil {
		return nil
	}
	if c.Id == 0 {
		return nil
	}
	return c
}

func DeleteCategory(cid int64) error {
	if _, err := core.Db.Where("category_id = ?", cid).Delete(new(CategoryArticle)); err != nil {
		return err
	}
	if _, err := core.Db.Where("id = ?", cid).Delete(new(Category)); err != nil {
		return err
	}
	return nil
}
