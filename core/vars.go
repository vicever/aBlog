package core

// global variables
type CoreVars struct {
	Name        string
	Version     string
	VersionDate string

	OfficalUrl string
	GithubUrl  string

	Author          string
	AuthorUrl       string
	AuthorGithubUrl string
}

// new default variables
func NewVars() CoreVars {
	return CoreVars{
		Name:        "aBlog",
		Version:     "0.2.0",
		VersionDate: "20150505",

		OfficalUrl: "https://github.com/fuxiaohei/aBlog",
		GithubUrl:  "https://github.com/fuxiaohei/aBlog",

		Author:          "FuXiaohei",
		AuthorUrl:       "http://fuxiaohei.me",
		AuthorGithubUrl: "https://github.com/fuxiaohei",
	}
}
