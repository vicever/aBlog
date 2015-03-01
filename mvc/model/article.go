package model

type Article struct {
	Id         int64  `xorm:"pk autoincr"`
	Uid        int64  `xorm:"notnull"`
	Title      string `xorm:"notnull"`
	Slug       string `xorm:"notnull unique(article-slug)"`
	CreateTime int64  `xorm:"created"`
	UpdateTime int64  `xorm:"updated"`

	ContentType string `xorm:"notnull"`
	Brief       string `xorm:"notnull"`
	Content     string `xorn:"notnull"`

	Status    int8 `xorm:"default 3"`
	IsComment int8 `xorm:"default 3"`
	IsTop     int8 `xorm:"default 1"`

	Comments int `xorm:"default 0"`
	Hits     int `xorm:"default 1"`

	tags       []*Tag `xorm:"-"`
	categories []*Category `xorm:"-"`
	comments []*Comment `xorm:"-"`
}
