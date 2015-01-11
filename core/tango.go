package core

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
	"path/filepath"
	"time"
)

var Tango *tango.Tango

func init() {
	Tango = tango.Classic()
	Tango.Use(
		renders.New(renders.Options{
			Reload:    true,
			Directory: "theme",
		}),
		xsrf.New(time.Minute*5),
	)

	Tango.Get("/theme/(.*)", new(themeAction))
}

type themeAction struct {
	tango.Ctx
	//tango.Compress
}

func (t *themeAction) Get() {
	path := t.Req().URL.Path
	if ext := filepath.Ext(path); ext == ".html" {
		t.NotFound()
		return
	}
	t.ServeFile(path[1:])
}
