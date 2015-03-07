package controller

import (
	"github.com/fuxiaohei/aBlog/lib/theme"
	"github.com/tango-contrib/renders"
)

type HomeController struct {
	renders.Renderer
}

func (hc *HomeController) Get() error {
	return hc.Render(theme.File("index.html"))
}
