package idl

import "github.com/tango-contrib/renders"

type ITheme interface {
	Assign(string, interface{}) // assign data
	Unassign(string)            // unassign data
	IsAssign(string) bool
	Render(string)                    // render template
	Assigned() map[string]interface{} // get assigned data
}

type ThemeRenderer struct {
	renders.Renderer
	viewData renders.T
}

func (tr *ThemeRenderer) Assign(key string, val interface{}) {
	if len(tr.viewData) == 0 {
		tr.viewData = make(renders.T)
		tr.viewData[key] = val
		return
	}
	tr.viewData[key] = val
}

func (tr *ThemeRenderer) Unassign(key string) {
	delete(tr.viewData, key)
}

func (tr *ThemeRenderer) Render(file string) {
	// render with data
	if err := tr.Renderer.Render(file, tr.viewData); err != nil {
		panic(err)
	}
}

func (tr *ThemeRenderer) Assigned() map[string]interface{} {
	return tr.viewData
}

func (tr *ThemeRenderer) IsAssign(key string) bool {
	if len(tr.viewData) == 0 {
		return false
	}
	return tr.viewData[key] != nil
}
