package base

import (
	"github.com/tango-contrib/renders"
	"path/filepath"
)

type Renders struct {
	renders.Renderer
}

func (r *Renders) RenderAdmin(tpl string, v interface{}) {
	if err := r.Render(filepath.Join("admin", tpl), v); err != nil {
		panic(err)
	}
}
