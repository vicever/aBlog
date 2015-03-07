package model

type Setting struct {
	Id    int64  `xorm:"pk autoincr"`
	Key   string `xorm:"notnull unique(user-setting)"`
	Value string `xorm:"text"`
	Type  string `xorm:"notnull"`
	Uid   int64  `xorm:"notnull unique(user-setting)"`
}
