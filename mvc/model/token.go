package model

import (
	"ablog/core"
	"ablog/util"
	"fmt"
	"time"
)

type Token struct {
	Id         int64 `model:"pk"`
	UserId     int64
	Value      string `model:"index"`
	ExpireTime time.Time
	duration   int64

	From      string
	Ip        string
	UserAgent string
}

func NewToken(uid int64, duration int64, from string) *Token {
	tk := &Token{
		Id:         generateTokenID(),
		UserId:     uid,
		ExpireTime: time.Now().Add(time.Duration(duration) * time.Second),
		From:       from,
		duration:   duration,
	}
	tk.generateValue()
	return tk
}

func generateTokenID() int64 {
	diff := time.Now().Unix() - core.Config.InstallTime
	return diff
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
	return core.Model.Save(tk)
}

// save expire-time-extended token
func (tk *Token) SaveExtend(duration int64) error {
	tk.duration = duration
	tk.ExpireTime.Add(time.Duration(duration) * time.Second)
	return tk.Save()
}

// remove token
func (tk *Token) Remove() error {
	return core.Model.Remove(tk)
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
	tk := &Token{Value: value}
	err := core.Model.GetBy(tk, "Value")
	return tk, err
}
