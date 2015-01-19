package core

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
	"path/filepath"
	"time"
)

type coreWeb struct {
	*tango.Tango
	opt serverConfig
}

func newCoreWeb(opt serverConfig) *coreWeb {
	t := tango.Classic()
	t.Use(
		renders.New(renders.Options{
			Reload:    true,
			Directory: "theme",
		}),
		xsrf.New(time.Minute*5),
	)

	// theme handler ,for theme static files
	t.Get("/theme/(.*)", new(themeAction))
	return &coreWeb{t, opt}
}

func (web *coreWeb) Run() {
	web.Tango.Run(web.opt.Addr + ":" + web.opt.Port)
}

type themeAction struct {
	tango.Ctx
}

func (t *themeAction) Get() {
	path := t.Req().URL.Path
	if ext := filepath.Ext(path); ext == ".html" {
		t.NotFound()
		return
	}
	t.ServeFile(path[1:])
}
