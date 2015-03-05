package controller
import (
    "github.com/fuxiaohei/aBlog/lib/theme"
    "github.com/tango-contrib/renders"
    "fmt"
)

type HomeController struct {
    theme.Themed
    renders.Renderer
}

func (hc *HomeController) Get() error{
    fmt.Println(hc.ThemeFile("index.html"))
    return hc.Render(hc.ThemeFile("index.html"))
}