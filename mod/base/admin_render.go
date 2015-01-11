package base

import (
	"github.com/tango-contrib/renders"
	"path/filepath"
)

type AdminRenders struct {
	renders.Renderer
}

func (a *AdminRenders) Render(tpl string, v interface{}) {
	if err := a.Renderer.Render(filepath.Join("admin", tpl), v); err != nil {
		panic(err)
	}
}
