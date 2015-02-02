package controller

import (
	"ablog/mvc/model"
	"fmt"
	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
	"net/url"
	"strconv"
	"time"
)

/*
===== admin article write controller
*/

type ArticleWriteController struct {
	tango.Ctx
	xsrf.Checker

	AdminBaseController
}

func (awc *ArticleWriteController) Get() {
	awc.Assign("IsArticlePage", true)
	awc.Assign("XsrfFormHtml", awc.XsrfFormHtml)
	awc.RenderAdmin("article_write.html")
}

func (awc *ArticleWriteController) Post() {
	//awc.Assign("IsArticlePage", true)

	article := &model.Article{
		AuthorId:   awc.AuthUser.Id,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	article.Title = awc.Req().FormValue("title")
	article.GenId()
	article.GenGuid()

	article.Slug = url.QueryEscape(awc.Req().FormValue("slug"))
	article.Content = awc.Req().FormValue("content")
	article.ContentType = awc.Req().FormValue("content-type")
	article.GenBrief()

	article.Status, _ = strconv.Atoi(awc.Req().FormValue("is-private"))
	article.CommentStatus, _ = strconv.Atoi(awc.Req().FormValue("is-comment"))

	article.Count = &model.ArticleCount{
		Hits:     1,
		Comments: 0,
	}
	article.Count.Words, _ = strconv.Atoi(awc.Req().FormValue("words"))

	fmt.Println(article)
}
