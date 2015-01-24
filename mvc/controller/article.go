package controller

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
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
	awc.RenderAdmin("article_write.html")
}
