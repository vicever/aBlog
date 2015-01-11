package user

import (
	"github.com/fuxiaohei/ablog/mod/base"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
)

type LoginAction struct {
	tango.Ctx
	base.AdminRenders
	xsrf.Checker
}

func (l *LoginAction) Get() {
	l.Render("login.html", renders.T{
		"XsrfFormHtml": l.XsrfFormHtml,
	})
}

func (l *LoginAction) Post() {

}
