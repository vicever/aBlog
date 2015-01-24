package model

import (
	"ablog/core"
	"ablog/util"
	"fmt"
	"time"
)

const (
	USER_ROLE_ADMIN = "admin"
)

var (
	user_id_key    = "user:id:%d"
	user_email_key = "user:email"
	user_name_key  = "user:name"
)

type User struct {
	Id           int64
	Name         string
	Password     string
	PasswordSalt string

	NickName string
	Email    string
	Url      string
	Bio      string

	CreateTime    time.Time
	LastLoginTime time.Time

	Role   string
	Social map[string]string
}

func NewUser(name, password, email, role string) *User {
	salt := util.MD5Short(password, "")
	password = util.Sha256(password, salt)
	u := &User{
		Name:         name,
		Email:        email,
		Role:         role,
		Password:     password,
		PasswordSalt: salt,
		NickName:     name,
		CreateTime:   time.Now(),
		Social:       make(map[string]string),
	}
	u.Id = generateUserID()
	return u
}

// check password
func (u *User) IsSamePassword(password string) bool {
	salt := util.MD5Short(password, "")
	password = util.Sha256(password, salt)
	return password == u.Password
}

// set new password
func (u *User) SetNewPassword(password string) {
	salt := util.MD5Short(password, "")
	password = util.Sha256(password, salt)
	u.PasswordSalt = salt
	u.Password = password
}

func generateUserID() int64 {
	diff := time.Now().Unix() - core.Config.InstallTime
	return diff/300 + 1
}

// save user data
func (u *User) Save() error {
	var err error

	// save user data
	key := fmt.Sprintf(user_id_key, u.Id)
	if err = core.Db.SetJson(key, u); err != nil {
		return err
	}

	// save name and email indexes
	if err = core.Db.HSet(user_name_key, u.Name, util.Int642Bytes(u.Id)); err != nil {
		return err
	}
	if err = core.Db.HSet(user_email_key, u.Email, util.Int642Bytes(u.Id)); err != nil {
		return err
	}

	return nil
}

// save new user email
func (u *User) SaveNewEmail(email string) error {
	// delete old index
	core.Db.HDel(user_email_key, u.Email)

	// save new index
	u.Email = email
	return u.Save()
}

// save new user name
func (u *User) SaveNewName(name string) error {
	// delete old index
	core.Db.HDel(user_name_key, u.Name)

	// save new index
	u.Name = name
	return u.Save()
}

func (u *User) saveUserData() error {
	// save user data
	key := fmt.Sprintf(user_id_key, u.Id)
	if err := core.Db.SetJson(key, u); err != nil {
		return err
	}
	return nil
}

// save new password
func (u *User) SaveNewPassword(password string) error {
	salt := util.MD5Short(password, "")
	password = util.Sha256(password, salt)
	u.Password = password
	u.PasswordSalt = salt
	return u.saveUserData()
}

// remove user
func (u *User) Remove() error {
	var err error

	// delete indexes
	core.Db.HDel(user_email_key, u.Email)

	core.Db.HDel(user_name_key, u.Name)

	// delete user data
	key := fmt.Sprintf(user_id_key, u.Id)
	if err = core.Db.Del(key); err != nil {
		return err
	}

	return nil
}

// get user by id
func GetUserById(id int64) (*User, error) {
	key := fmt.Sprintf(user_id_key, id)
	u := &User{}
	if err := core.Db.GetJson(key, u); err != nil {
		return nil, err
	}
	if u.Id != id {
		return nil, nil
	}
	return u, nil
}

// get user by email
func GetUserByEmail(email string) (*User, error) {
	bytes, err := core.Db.HGet(user_email_key, email)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, err
	}
	id := util.Bytes2Int64(bytes)
	return GetUserById(id)
}

// get user by name
func GetUserByName(name string) (*User, error) {
	bytes, err := core.Db.HGet(user_name_key, name)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, err
	}
	id := util.Bytes2Int64(bytes)
	return GetUserById(id)
}
