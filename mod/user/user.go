package user

import (
	"fmt"
	"github.com/fuxiaohei/ablog/core"
	"github.com/fuxiaohei/ablog/util"
	"time"
)

const (
	ROLE_ADMIN = iota + 1
	ROLE_WRITER
	ROLE_READER

	userKeyString = "user:%d"
)

type User struct {
	Id           int64
	Name         string
	Password     string
	PasswordSalt string
	Email        string

	Nick   string
	Url    string
	Bio    string
	Social map[string]string

	Role        int
	CreateTime  int64
	LastLogTime int64
}

func createUserId(t int64, diff int64) int64 {
	id := (t-core.Vars.Status.InitTime)/1800 + diff
	key := fmt.Sprintf(userKeyString, id)
	if core.Db.Exist(key) {
		return createUserId(t, diff+1)
	}
	return id
}

func createPassword(password string) (encoded, salt string) {
	salt = util.MD5Short(password, "")
	encoded = util.Sha1(password, salt)
	return encoded, salt
}

// create new user with simple data
func Create(name, email, password string, role int) (*User, error) {
	now := time.Now()

	u := &User{
		Id:          createUserId(now.Unix(), 1),
		Name:        name,
		Email:       email,
		Nick:        name,
		Url:         "#",
		Social:      make(map[string]string),
		Role:        role,
		CreateTime:  now.Unix(),
		LastLogTime: now.Unix(),
	}

	u.Password, u.PasswordSalt = createPassword(password)

	// write to db
	if err := core.Db.SetJson(fmt.Sprintf(userKeyString, u.Id), u); err != nil {
		return nil, err
	}

	return u, nil
}
