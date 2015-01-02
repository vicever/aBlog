package user

import (
	"fmt"
	"github.com/fuxiaohei/ablog/sys"
	"github.com/fuxiaohei/ablog/util"
	"time"
)

const (
	ROLE_ADMIN = iota + 1
	ROLE_WRITER
	ROLE_READER

	userIdKey    = "user:%d"
	userNameKey  = "user:name:%s"
	userEmailKey = "user:email:%s"

	EVENT_USER_CREATE = "user-create"
	EVENT_USER_UPDATE = "user-update"
	EVENT_USER_DELETE = "user-delete"
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

func (u *User) IsPassword(password string) bool {
	password, salt := createPassword(password)
	if password == u.Password && salt == u.PasswordSalt {
		return true
	}
	return false
}

func init() {
	sys.Event.On(EVENT_USER_CREATE, func(msg sys.EventMsg) {
		u, ok := msg.Value.(*User)
		if !ok {
			return
		}
		sys.NoDb.SetInt64(fmt.Sprintf(userNameKey, u.Name), u.Id)
		sys.NoDb.SetInt64(fmt.Sprintf(userEmailKey, u.Email), u.Id)
	})
	sys.Event.On(EVENT_USER_UPDATE, func(msg sys.EventMsg) {
		u, ok := msg.Value.(*User)
		if !ok {
			return
		}
		sys.NoDb.SetInt64(fmt.Sprintf(userNameKey, u.Name), u.Id)
		sys.NoDb.SetInt64(fmt.Sprintf(userEmailKey, u.Email), u.Id)
	})
}

func createUserId(t int64, diff int64) int64 {
	id := (t-sys.Config.InitTime.Unix())/1800 + diff
	key := fmt.Sprintf(userIdKey, id)
	if sys.NoDb.Exist(key) {
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
	if err := sys.NoDb.SetJson(fmt.Sprintf(userIdKey, u.Id), u); err != nil {
		return nil, err
	}

	sys.Event.Call(EVENT_USER_CREATE, u)

	return u, nil
}

// get user by id
func Get(uid int64) (*User, error) {
	u := new(User)
	err := sys.NoDb.GetJson(fmt.Sprintf(userIdKey, uid), u)
	if err != nil {
		return nil, err
	}
	if u.Id == 0 || len(u.Name) == 0 {
		return nil, nil
	}
	return u, nil
}

// get user by name
func GetByName(name string) (*User, error) {
	key := fmt.Sprintf(userNameKey, name)
	uid, err := sys.NoDb.GetInt64(key)
	if err != nil {
		return nil, err
	}
	u, err := Get(uid)
	if u == nil {
		sys.NoDb.Del(key)
	}
	return u, err
}

// get user by email
func GetByEmail(email string) (*User, error) {
	key := fmt.Sprintf(userEmailKey, email)
	uid, err := sys.NoDb.GetInt64(key)
	if err != nil {
		return nil, err
	}
	u, err := Get(uid)
	if u == nil {
		sys.NoDb.Del(key)
	}
	return u, err
}

// update user
func Update(u *User) error {
	if u.Id == 0 || len(u.Name) == 0 || len(u.Email) == 0 {
		return nil
	}
	err := sys.NoDb.SetJson(fmt.Sprintf(userIdKey, u.Id), u)
	if err == nil {
		sys.Event.Call(EVENT_USER_UPDATE, u)
	}
	return err
}

// update user's password
func UpdatePassword(uid int64, password string) error {
	user, err := Get(uid)
	if err != nil {
		return err
	}
	user.Password, user.PasswordSalt = createPassword(password)
	return Update(user)
}

// delete user
func Delete(uid int64) error {
	key := fmt.Sprintf(userIdKey, uid)
	err := sys.NoDb.Del(key)
	if err == nil {
		sys.Event.Call(EVENT_USER_DELETE, uid)
	}
	return err
}
