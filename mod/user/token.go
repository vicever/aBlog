package user

import (
	"fmt"
	"github.com/fuxiaohei/ablog/sys"
	"time"
)

const (
	tokenValueKey = "token:%s"

	EVENT_TOKEN_CREATE = "token-create"
)

type Token struct {
	Uid       int64
	Value     string
	ValueSalt string

	CreateTime time.Time
	ExpireTime time.Time

	Src       string
	UserAgent string
}

func (t *Token) IsExpired() bool {
	return time.Now().Unix()-t.ExpireTime.Unix() > 0
}

func CreateToken(uid int64, expire int, src, userAgent string) (*Token, error) {
	now := time.Now()
	token := &Token{
		Uid:        uid,
		CreateTime: now,
		ExpireTime: now.Add(time.Duration(expire) * time.Second),
		Src:        src,
		UserAgent:  userAgent,
	}

	// create token
	tokenValue := fmt.Sprint("%d:%d:%s:%s", uid, now.Unix(), src, userAgent)
	token.Value, token.ValueSalt = createPassword(tokenValue)

	// write to database
	key := fmt.Sprintf(tokenValueKey, token.Value)
	if err := sys.NoDb.SetJson(key, token); err != nil {
		return nil, err
	}

	sys.Event.Call(EVENT_TOKEN_CREATE, token)
	return token, nil
}

func GetToken(value string) (*Token, *User, error) {
	token := MustGetToken(value)
	if token == nil || token.IsExpired() {
		// invalid token value, delete it
		key := fmt.Sprintf(tokenValueKey, value)
		sys.NoDb.Del(key)
		return nil, nil, nil
	}
	user, err := Get(token.Uid)
	return token, user, err
}

func MustGetToken(value string) *Token {
	// get token first
	key := fmt.Sprintf(tokenValueKey, value)
	token := new(Token)
	err := sys.NoDb.GetJson(key, token)
	if err != nil {
		return nil
	}
	if token.Uid == 0 || token.ValueSalt == "" {
		return nil
	}

	return token
}
