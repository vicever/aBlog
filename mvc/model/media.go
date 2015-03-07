package model

type Media struct {
	Id  int64 `xorm:"pk autoincr"`
	Uid int64 `xorm:"notnull"`

	Title       string `xorm:"notnull"`
	Path        string `xorm:"notnull"`
	Type        string `xorm:"notnull"`
	ContentType string `xorm:"notnull"`
	Size        int32
	Description string
	CreateTime  int64
	DeleteTime  int64

	Downloads int32 `xorm:"default 0"`
}
