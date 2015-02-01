package model

import "time"

type Article struct {
	Id   int64  `model:"pk"`
	Guid string `model:"unique"`

	Title    string
	Slug     string `model:"unique"`
	AuthorId int64  `model:"index"`

	CreateTime time.Time
	UpdateTime time.Time

	Brief       string
	Content     string
	ContentType string

	Status        int `model:"index"`
	CommentStatus int

	Count *ArticleCount `model:"1-1" json:"-"` // article count
	Tags  []*ArticleTag `model:"1-n" json:"-"` // article tags
}

type ArticleCount struct {
	Words    int
	Hits     int
	Comments int
}

type ArticleTag struct {
	Id   int64  `model:"pk"`
	Name string `model:"unique"`
}
