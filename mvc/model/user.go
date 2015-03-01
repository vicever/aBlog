package model

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/fuxiaohei/aBlog/core"
	"time"
)

type User struct {
	Id           int64  `xorm:"pk autoincr" json:"id"`
	Name         string `xorm:"notnull index(user-name)" json:"name"`
	Password     string `xorm:"notnull" json:"-"`
	PasswordSalt string `xorm:"notnull" json:"-"`

	Nick  string `json:"nick"`
	Email string `xorm:"notnull unique(user-email)" json:"email"`
	Url   string `json:"url"`
	Bio   string `json:"bio"`

	CreateTime int64 `xorm:"created" json:"created"`
	LoginTime  int64 `json:"logged"`

	Role   int8 `xorm:"notnull default 1" json:"role"`
	Status int8 `xorm:"notnull default 1" json:"status"`
}

func (u *User) GeneratePassword(password string) (string, string) {
	m := md5.New()
	m.Write([]byte(u.Name + u.Email))
	u.PasswordSalt = hex.EncodeToString(m.Sum(nil))

	s := sha256.New()
	s.Write([]byte(password))
	s.Write([]byte(u.PasswordSalt))
	u.Password = hex.EncodeToString(s.Sum(nil))
	return u.Password, u.PasswordSalt
}

func (u *User) CheckPassword(password string) bool {
	s := sha256.New()
	s.Write([]byte(password))
	s.Write([]byte(u.PasswordSalt))
	return hex.EncodeToString(s.Sum(nil)) == u.Password
}

func GetUserBy(column string, value interface{}) *User {
	u := new(User)
	_, err := core.Db.Where(column+" = ?", value).Get(u)
	if err != nil {
		return nil
	}
	return u
}

func CreateUser(name string, password string, email string, role int8) *User {
	u := &User{
		Name:       name,
		Nick:       name,
		Email:      email,
		Url:        "#",
		CreateTime: time.Now().Unix(),
		LoginTime:  time.Now().Unix(),
		Role:       role,
		Status:     1,
	}
	u.GeneratePassword(password)
	if _, err := core.Db.Insert(u); err != nil {
		return nil
	}
	return u
}
