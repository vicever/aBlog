package action

import (
	"github.com/fuxiaohei/aBlog/core"
	"github.com/fuxiaohei/aBlog/mvc/model"
	"strconv"
)

const ERR_AUTH_NO_USER = 2001       // get no user
const ERR_AUTH_WRONG_PWD = 2002     // wrong password for user
const ERR_AUTH_NO_TOKEN = 2003      // get no token
const ERR_AUTH_EXPIRED_TOKEN = 2004 // get expired token
const ERR_AUTH_NO_ACCESS = 2005     // no access to do

func init() {
	errorMap[ERR_AUTH_NO_USER] = "no-user"
	errorMap[ERR_AUTH_WRONG_PWD] = "wrong-password"
	errorMap[ERR_AUTH_NO_TOKEN] = "no-token"
	errorMap[ERR_AUTH_EXPIRED_TOKEN] = "expired-token"
	errorMap[ERR_AUTH_NO_ACCESS] = "no-access"
}

// user authorize,
// input "user, password"
// output "user:*model.User, token:*model.Token"
func Authorize(param ActionParam) ActionResult {
	// need user and password param
	user := param["user"]
	password := param["password"]
	if user == "" || password == "" {
		return NewResultError(ERR_INVALID_PARAMS)
	}

	// get user
	u := model.GetUserBy("name", user)
	if u == nil {
		return NewResultError(ERR_AUTH_NO_USER)
	}

	// check password
	if !u.CheckPassword(password) {
		return NewResultError(ERR_AUTH_WRONG_PWD)
	}

	// create token
	token := model.CreateToken(u.Id, 3*24*3600, param["ip"], param["mark"])
	return NewResult(map[string]interface{}{
		"user":  u,
		"token": token,
	})
}

// user's token is valid
// input "uid,token"
// output "uid:uid, token:*model.Token, is_authorized:true"
func IsAuthorized(param ActionParam) ActionResult {
	// need user_id and token value
	uid, _ := strconv.ParseInt(param["uid"], 10, 64)
	token := param["token"]
	if uid == 0 || token == "" {
		return NewResultError(ERR_INVALID_PARAMS)
	}

	// get token by uid and token
	t := model.GetToken(uid, token)
	if t == nil {
		return NewResultError(ERR_AUTH_NO_TOKEN)
	}
	if t.IsExpired() {
		model.RemoveToken(uid, token)
		return NewResultError(ERR_AUTH_EXPIRED_TOKEN)
	}
	return NewResult(map[string]interface{}{
		"uid":           uid,
		"token":         t,
		"is_authorized": true,
	})
}

// get user by other-user. it checkes other-user's access.
// input "uid,oid"
// output "user:*model.User"
func GetUser(params ActionParam) ActionResult {
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	oid, _ := strconv.ParseInt(params["oid"], 10, 64)

	// if the user get own user data, just return
	if uid == oid {
		u := model.GetUserBy("id", uid)
		return NewResult(map[string]interface{}{
			"user": u,
		})
	}

	// check oid's access
	o := model.GetUserBy("id", oid)
	if o == nil || o.Role != 1 {
		return NewResultError(ERR_AUTH_NO_ACCESS)
	}

	// if have access, get user info
	u := model.GetUserBy("id", uid)
	return NewResult(map[string]interface{}{
		"user": u,
	})
}

// update user by uid. it only updates general infomations, not all items.
// input "uid,name,nick,email,url,bio"
// output "user:*model.User" - updated user data
func UpdateUser(params ActionParam) ActionResult {
	// build data
	user := model.User{
		Name:  params["name"],
		Nick:  params["nick"],
		Email: params["email"],
		Url:   params["url"],
		Bio:   params["bio"],
	}
	if user.Nick == "" {
		user.Nick = user.Name
	}
	if user.Url == "" {
		user.Url = "#"
	}
	// update user
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	if _, err := core.Db.Cols("name,nick,email,url,bio").Where("id = ?", uid).Update(user); err != nil {
		return NewSystemError(err)
	}
	// get new data to return
	u := model.GetUserBy("id", uid)
	return NewResult(map[string]interface{}{
		"user": u,
	})
}

// update user's password. it needs old password to check user.
// input "uid,old,new"
// output null - only return meta, no data
func UpdateUserPassword(params ActionParam) ActionResult {
	uid, _ := strconv.ParseInt(params["uid"], 10, 64)
	oldPassword := params["old"]
	newPassword := params["new"]
	if len(newPassword) < 6 {
		return NewResultError(ERR_INVALID_PARAMS)
	}

	// get user
	u := model.GetUserBy("id", uid)
	if u == nil {
		return NewResultError(ERR_AUTH_NO_USER)
	}
	if !u.CheckPassword(oldPassword) {
		return NewResultError(ERR_AUTH_WRONG_PWD)
	}
	u.GeneratePassword(newPassword)
	if _, err := core.Db.Cols("password,password_salt").Where("id = ?", uid).Update(u); err != nil {
		return NewSystemError(err)
	}
	return NewResult(nil)
}
