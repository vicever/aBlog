package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fuxiaohei/aBlog/core"
	"time"
)

type Token struct {
	Id  int64 `xorm:"pk autoincr" json:"-"`
	Uid int64 `xorm:"notnull index(user-id)" json:"u_id"`

	Value      string `xorm:"notnull index(token)" json:"value"`
	CreateTime int64  `xorm:"created" json:"created"`
	ExpireTime int64  `xorm:"notnull" json:"expired"`

	Ip   string `xorm:"notnull" json:"-"`
	Mark string `xorm:"notnull default 'web' index(mark)" json:"mark"`

	Extra string `json:"extra"`
}

func (t *Token) GenerateValue() string {
	s := sha256.New()
	s.Write([]byte(fmt.Sprint(t.Uid)))
	s.Write([]byte(fmt.Sprint(t.CreateTime)))
	s.Write([]byte(fmt.Sprint(t.Mark)))
	t.Value = hex.EncodeToString(s.Sum(nil))
	return t.Value
}

func (t *Token) IsExpired() bool {
	return time.Now().Unix() > t.ExpireTime
}

func CreateToken(uid int64, expire int64, ip, mark string, extra ...string) *Token {
	t := &Token{
		Uid:        uid,
		CreateTime: time.Now().Unix(),
		Ip:         ip,
		Mark:       mark,
	}
	t.ExpireTime = t.CreateTime + expire
	if len(extra) > 0 {
		t.Extra = extra[0]
	}
	t.GenerateValue()
	if _, err := core.Db.Insert(t); err != nil {
		return nil
	}
	return t
}

func GetToken(uid int64, token string) *Token {
	t := new(Token)
	if _, err := core.Db.Where("uid = ? AND value = ?", uid, token).Get(t); err != nil {
		return nil
	}
	return t
}

func RemoveToken(uid int64, token string) *Token {
	t := GetToken(uid, token)
	if t == nil {
		return nil
	}
	if _, err := core.Db.Where("uid = ? AND value = ?", uid, token).Delete(new(Token)); err != nil {
		return nil
	}
	return t
}
