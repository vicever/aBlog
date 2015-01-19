package core

type coreVars struct {
	Name          string
	Version       string
	VersionDate   string
	VersionStatus string
	Description   string

	Official string

	Author       string
	AuthorEmail  string
	AuthorGithub string
}

func newCoreVars() *coreVars {
	vars := &coreVars{
		Name:          "ablog",
		Version:       "0.2",
		VersionDate:   "20150122",
		VersionStatus: "alpha",
		Description:   "an golang blog engine",

		Official: "http://github.com/fuxiaohei/ablog",

		Author:       "fuxiaohei",
		AuthorEmail:  "fuxiaohei@vip.qq.com",
		AuthorGithub: "https://github.com/fuxiaohei",
	}
	return vars
}
