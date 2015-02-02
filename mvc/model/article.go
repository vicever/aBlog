package model

import (
	"ablog/core"
	"ablog/util"
	"fmt"
	"strings"
	"time"
)

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

func (a *Article) GenId() {
	diff := time.Now().Unix() - core.Config.InstallTime
	a.Id = diff/1800 + 1
}

func (a *Article) GenGuid() {
	a.Guid = util.MD5Short(a.Title, fmt.Sprint(a.CreateTime.UnixNano()))
}

func (a *Article) GenBrief() {
	str := strings.Split(a.Content, "<--more-->")
	a.Brief = str[0]
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
