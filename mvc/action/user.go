package action

import (
	"github.com/fuxiaohei/aBlog/mvc/model"
	"strconv"
    "github.com/fuxiaohei/aBlog/core"
)

const ERR_AUTH_NO_USER = 2001
const ERR_AUTH_WRONG_PWD = 2002
const ERR_AUTH_NO_TOKEN = 2003
const ERR_AUTH_EXPIRED_TOKEN = 2004
const ERR_AUTH_NO_ACCESS = 2005

func init() {
	errorMap[ERR_AUTH_NO_USER] = "no-user"
	errorMap[ERR_AUTH_WRONG_PWD] = "wrong-password"
	errorMap[ERR_AUTH_NO_TOKEN] = "no-token"
	errorMap[ERR_AUTH_EXPIRED_TOKEN] = "expired-token"
    errorMap[ERR_AUTH_NO_ACCESS] = "no-access"
}

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

func GetUser(params ActionParam) ActionResult {
    uid,_ := strconv.ParseInt(params["uid"],10,64)
    oid,_ := strconv.ParseInt(params["oid"],10,64)

    // if the user get own user data, just return
    if uid == oid{
        u := model.GetUserBy("id",uid)
        return NewResult(map[string]interface{}{
            "user":u,
        })
    }

    // check oid's access
    o := model.GetUserBy("id",oid)
    if o == nil || o.Role != 1{
        return NewResultError(ERR_AUTH_NO_ACCESS)
    }

    // if have access, get user info
    u := model.GetUserBy("id",uid)
    return NewResult(map[string]interface{}{
        "user":u,
    })
}

func UpdateUser(params ActionParam) ActionResult {
    // build data
    user := model.User{
        Name:params["name"],
        Nick:params["nick"],
        Email:params["email"],
        Url:params["url"],
        Bio:params["bio"],
    }
    if user.Nick == ""{
        user.Nick = user.Name
    }
    if user.Url == ""{
        user.Url = "#"
    }
    // update user
    uid,_ := strconv.ParseInt(params["uid"],10,64)
    if _,err := core.Db.Cols("name,nick,email,url,bio").Where("id = ?",uid).Update(user);err != nil{
        return NewSystemError(err)
    }
    // get new data to return
    u := model.GetUserBy("id",uid)
    return NewResult(map[string]interface{}{
        "user":u,
    })
}
