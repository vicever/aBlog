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
	user_email_key = "user:email:%s"
	user_name_key  = "user:name:%s"
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
	key = fmt.Sprintf(user_name_key, u.Name)
	if err = core.Db.Set(key, util.Int642Bytes(u.Id)); err != nil {
		return err
	}
	key = fmt.Sprintf(user_email_key, u.Email)
	if err = core.Db.Set(key, util.Int642Bytes(u.Id)); err != nil {
		return err
	}

	return nil
}

// save new user email
func (u *User) SaveNewEmail(email string) error {
	// delete old index
	key := fmt.Sprintf(user_email_key, u.Email)
	core.Db.Del(key)

	// save new index
	u.Email = email
	return u.Save()
}

// save new user name
func (u *User) SaveNewName(name string) error {
	// delete old index
	key := fmt.Sprintf(user_email_key, u.Email)
	core.Db.Del(key)

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
	key := fmt.Sprintf(user_name_key, u.Name)
	if err = core.Db.Del(key); err != nil {
		return err
	}
	key = fmt.Sprintf(user_email_key, u.Email)
	if err = core.Db.Del(key); err != nil {
		return err
	}

	// delete user data
	key = fmt.Sprintf(user_id_key, u.Id)
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
	key := fmt.Sprintf(user_email_key, email)
	bytes, err := core.Db.Get(key)
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
	key := fmt.Sprintf(user_name_key, name)
	bytes, err := core.Db.Get(key)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, err
	}
	id := util.Bytes2Int64(bytes)
	return GetUserById(id)
}
