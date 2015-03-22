package model

import (
	"github.com/fuxiaohei/aBlog/core"
	"time"
)

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

// create new category
func CreateCategory(uid int64, name, slug, desc string) (*Category, error) {
	c := &Category{
		Uid:         uid,
		Name:        name,
		Slug:        slug,
		Description: desc,
		Count:       0,
		CreateTime:  time.Now().Unix(),
	}
	if _, err := core.Db.Insert(c); err != nil {
		return nil, err
	}
	return c, nil
}

// update category by cid and uid
func UpdateCategory(uid, cid int64, name, slug, desc string) (*Category, error) {
	c := &Category{
		Name:        name,
		Slug:        slug,
		Description: desc,
		Count:       0,
		CreateTime:  time.Now().Unix(),
	}
	if _, err := core.Db.Cols("name,slug,description").Where("id = ? AND uid = ?", cid, uid).Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

// get category by column
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

// get user's categories list
func GetCategoriesByUser(uid int64, countOrder bool) []*Category {
	s := core.Db.Where("uid = ?", uid)
	if countOrder {
		s.OrderBy("count DESC")
	} else {
		s.OrderBy("id ASC")
	}
	cs := make([]*Category, 0)
	if err := s.Find(&cs); err != nil {
		return nil
	}
	return cs
}

// delete category by id
func DeleteCategory(cid int64) error {
	if _, err := core.Db.Where("category_id = ?", cid).Delete(new(CategoryArticle)); err != nil {
		return err
	}
	if _, err := core.Db.Where("id = ?", cid).Delete(new(Category)); err != nil {
		return err
	}
	return nil
}
