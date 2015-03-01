package model

type Comment struct {
	Id       int64 `xorm:"pk autoincr"`
	ParentId int64 `xorm:"notnull default 0"`

	Email        string `xorm:"notnull"`
	Url          string
	CreateTime   int64 `xorm:"created"`
	ApprovedTime int64
	ContentType  string `xorm:"notnull"`
	Content      string `xorm:"notnull"`

	Status  int8  `xorm:"notnull default 1"`
	RelType int8  `xorm:"notnull default 1"`
	RelId   int64 `xorm:"notnull"`

	Ip        string
	UserAgent string

	Uid      int64 `xorm:"notnull"`
	UserRole int8  `xorm:"notnull"`
}
