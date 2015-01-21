package controller

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

var _ Auther = new(AdminProfileController)

type AdminProfileController struct {
	tango.Ctx
	AdminRender
	AdminAutherController
}

func (aPc *AdminProfileController) Get() {
	aPc.Render("profile.html", renders.T{
		"AuthUser": aPc.AuthUser,
	})
}
