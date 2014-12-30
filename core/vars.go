package core

var Vars coreVars = defaultVars()

type coreVars struct {
	Name        string // blog name
	Description string // blog description
	Homepage    string // homepage url
	Author      string // author
	GitPage     string // github page url

	Version     string // version number
	VersionDate string // version published date

	DataDirectory string // data store directory
	DbIndex       int    // data set index

	InitLock string // init lock key

	Status coreVarsStatus // global status
}

type coreVarsStatus struct {
	IsInit bool
}

func defaultVars() coreVars {
	vars := coreVars{
		Name:        "ablog",
		Description: "a blog engine",
		Homepage:    "http://github.com/fuxiaohei/ablog",
		Author:      "fuxiaohei",
		GitPage:     "http://github.com/fuxiaohei/ablog",

		Version:     "0.1",
		VersionDate: "2014-12-30",

		DataDirectory: "data",
		InitLock:      "init-lock",
	}
	vars.Status = coreVarsStatus{
		IsInit: false,
	}
	return vars
}
