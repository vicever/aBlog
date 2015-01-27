package model

import (
	"ablog/core"
	"ablog/util"
	"time"
)

const (
	USER_ROLE_ADMIN = "admin"
)

type User struct {
	Id           int64  `model:"pk"`
	Name         string `model:"unique"`
	Password     string
	PasswordSalt string

	NickName string
	Email    string `model:"unique"`
	Url      string
	Bio      string

	CreateTime    time.Time
	LastLoginTime time.Time

	Role   string `model:"index"`
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
	return core.Model.Save(u)
}

func (u *User) Update() error {
	// get old user whole data
	oldUser := &User{Id: u.Id}
	core.Model.Get(oldUser)
	// remove it
	core.Model.Remove(oldUser)
	// save current
	return core.Model.Save(u)
}

// remove user
func (u *User) Remove() error {
	return core.Model.Remove(u)
}

// get user by id
func GetUserById(id int64) (*User, error) {
	user := &User{Id: id}
	err := core.Model.Get(user)
	return user, err
}

// get user by email
func GetUserByEmail(email string) (*User, error) {
	user := &User{Email: email}
	err := core.Model.GetBy(user, "Email")
	return user, err
}

// get user by name
func GetUserByName(name string) (*User, error) {
	user := &User{Name: name}
	err := core.Model.GetBy(user, "Name")
	return user, err
}
