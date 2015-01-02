package setting

import (
	"fmt"
	"github.com/fuxiaohei/ablog/sys"
)

const (
	settingSysKey  = "setting:sys:%s"
	settingUserKey = "setting:user:%d:%s"
)

func SetSys(key, value string) error {
	key = fmt.Sprintf(settingSysKey, key)
	return sys.NoDb.Set(key, []byte(value))
}

func SetSysInt(key string, value int) error {
	key = fmt.Sprintf(settingSysKey, key)
	return sys.NoDb.SetInt(key, value)
}

func SetSysJson(key string, value interface{}) error {
	key = fmt.Sprintf(settingSysKey, key)
	return sys.NoDb.SetJson(key, value)
}

func SetUser(uid int64, key, value string) error {
	key = fmt.Sprintf(settingUserKey, uid, key)
	return sys.NoDb.Set(key, []byte(value))
}

func SetUserInt(uid int64, key string, value int) error {
	key = fmt.Sprintf(settingUserKey, uid, key)
	return sys.NoDb.SetInt(key, value)
}

func SetUserJson(uid int64, key string, value interface{}) error {
	key = fmt.Sprintf(settingUserKey, uid, key)
	return sys.NoDb.SetJson(key, value)
}

func GetSys(key string) (string, error) {
	key = fmt.Sprintf(settingSysKey, key)
	value, err := sys.NoDb.Get(key)
	return string(value), err
}

func GetSysInt(key string) (int, error) {
	key = fmt.Sprintf(settingSysKey, key)
	return sys.NoDb.GetInt(key)
}

func GetSysJson(key string, value interface{}) error {
	key = fmt.Sprintf(settingSysKey, key)
	return sys.NoDb.GetJson(key, value)
}

func GetUser(uid int64, key string) (string, error) {
	key = fmt.Sprintf(settingUserKey, uid, key)
	value, err := sys.NoDb.Get(key)
	return string(value), err
}

func GetUserInt(uid int64, key string) (int, error) {
	key = fmt.Sprintf(settingUserKey, uid, key)
	return sys.NoDb.GetInt(key)
}

func GetUserJson(uid int64, key string, value interface{}) error {
	key = fmt.Sprintf(settingUserKey, uid, key)
	return sys.NoDb.GetJson(key, value)
}
