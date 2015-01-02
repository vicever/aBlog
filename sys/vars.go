package sys

var Vars sysVars = sysVars{
	Name:        "ablog",
	Description: "ablog engine",
	Homepage:    "http://github.com/fuxiaohei/ablog",
	Author:      "fuxiaohei",
	GitPage:     "http://fuxiaohei/ablog",

	Version:     "0.1",
	VersionDate: "2015-01-01",
}

type sysVars struct {
	Name        string
	Description string // blog description
	Homepage    string // homepage url
	Author      string // author
	GitPage     string // github page url

	Version     string // version number
	VersionDate string // version published date
}
