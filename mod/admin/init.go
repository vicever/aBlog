package admin

import (
	"github.com/fuxiaohei/ablog/core"
	"github.com/lunny/tango"
)

func Init() {
	var t *tango.Tango = core.Tango

	t.Get("/admin/dashboard", new(DashboardAction))
}
