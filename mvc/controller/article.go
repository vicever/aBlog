package controller

import (
	"ablog/mvc/model"
	"github.com/lunny/tango"
)

/*
===== admin articles controller
*/

var _ Auther = new(AdminArticleController)

type AdminArticleController struct {
	tango.Ctx
	authUser *model.User
}

func (aAc *AdminArticleController) SetAuthUser(u *model.User) {
	aAc.authUser = u
}

func (aAc *AdminArticleController) AuthFailRedirect() string {
	return "/login?error=auth-fail"
}
