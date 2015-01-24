package model

import (
	"ablog/core"
	"ablog/util"
	"fmt"
	"time"
)

var (
	token_value_key = "token:%s"
)

type Token struct {
	UserId     int64
	Value      string
	ExpireTime time.Time
	duration   int64

	From      string
	Ip        string
	UserAgent string
}

func NewToken(uid int64, duration int64, from string) *Token {
	tk := &Token{
		UserId:     uid,
		ExpireTime: time.Now().Add(time.Duration(duration) * time.Second),
		From:       from,
		duration:   duration,
	}
	tk.generateValue()
	return tk
}

func (tk *Token) generateValue() {
	tk.Value = util.Sha256(fmt.Sprintf("%d.%d", tk.UserId, tk.ExpireTime.Unix()), tk.From)
}

// check the token expiration
func (tk *Token) IsExpired() bool {
	return time.Now().Unix() > tk.ExpireTime.Unix()
}

// get user in token
func (tk *Token) GetUser() (*User, error) {
	return GetUserById(tk.UserId)
}

// save token
func (tk *Token) Save() error {
	key := fmt.Sprintf(token_value_key, tk.Value)
	if err := core.Db.SetJson(key, tk); err != nil {
		return err
	}
	if err := core.Db.SetExpire(key, tk.duration); err != nil {
		return err
	}
	return nil
}

// save expire-time-extended token
func (tk *Token) SaveExtend(duration int64) error {
	tk.duration = duration
	tk.ExpireTime.Add(time.Duration(duration) * time.Second)
	return tk.Save()
}

// remove token
func (tk *Token) Remove() error {
	key := fmt.Sprintf(token_value_key, tk.Value)
	if err := core.Db.Del(key); err != nil {
		return err
	}
	return nil
}

// get token with expiration checking
func GetToken(value string) (*Token, error) {
	tk, err := GetTokenIgnoreExpired(value)
	if err != nil {
		return nil, err
	}
	if tk.IsExpired() {
		return nil, nil
	}
	return tk, nil
}

// must get token
func GetTokenIgnoreExpired(value string) (*Token, error) {
	key := fmt.Sprintf(token_value_key, value)
	tk := &Token{}
	if err := core.Db.GetJson(key, tk); err != nil {
		return nil, err
	}
	if tk.Value != value {
		return nil, nil
	}
	return tk, nil
}
