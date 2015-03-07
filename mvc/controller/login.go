package controller
import (
    "github.com/tango-contrib/renders"
    "github.com/fuxiaohei/aBlog/lib/theme"
    "github.com/lunny/tango"
)

type LoginController struct {
    tango.Ctx
    renders.Renderer
}

func (lc *LoginController) Get() {
    lc.Render(theme.File("admin/login.html"))
}

func (lc *LoginController) Post(){
    lc.ResponseWriter.Write([]byte("login post"))
}