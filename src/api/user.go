// user
// @author xiangqian
// @date 18:10 2022/12/18
package api

import (
	"auto-deploy-go/src/com"
	"auto-deploy-go/src/db"
	"encoding/gob"
	"errors"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const SessionKeyUser = "_user_"

type User struct {
	Abs
	Name     string `form:"name" binding:"required,excludes= ,min=1,max=60"`               // 用户名
	Nickname string `form:"nickname"binding:"max=60"`                                      // 昵称
	Passwd   string `form:"passwd" binding:"required,excludes= ,max=100"`                  // 密码
	RePasswd string `form:"rePasswd" binding:"required,excludes= ,max=100,eqfield=Passwd"` // retype Passwd
}

func init() {
	// 注册 User 模型
	gob.Register(User{})
}

// 用户注册html
func UserRegPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	user := session.Get("user")
	message := session.Get("message")
	session.Delete("user")
	session.Delete("message")
	session.Save()

	if user == nil {
		user = User{}
	}

	pContext.HTML(http.StatusOK, "user/reg.html", gin.H{
		"user":    user,
		"message": message,
	})
}

// 用户注册
func UserAdd(pContext *gin.Context) {
	user := User{}
	err := ShouldBind(pContext, &user)

	if err == nil {
		err = VerifyUserNameAndPasswd(user.Name, user.Passwd)
	}

	if err == nil {
		err = VerifyDbUserName(user.Name)
	}

	session := sessions.Default(pContext)
	if err != nil {
		session.Set("user", user)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/regpage")
		return
	}

	db.Add("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`, `create_time`) VALUES (?, ?, ?, ?, ?)",
		user.Name, strings.TrimSpace(user.Nickname), user.Passwd, strings.TrimSpace(user.Rem), time.Now().Unix())

	session.Set("username", user.Name)
	session.Set("message", i18n.MustGetMessage("i18n.accountRegSuccess"))
	session.Save()

	// 用户注册成功后，重定向到登录页
	pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")
}

// 用户登录html
func UserLoginPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	username := session.Get("username")
	message := session.Get("message")
	session.Delete("username")
	session.Delete("message")
	session.Save()

	pContext.HTML(http.StatusOK, "user/login.html", gin.H{
		"username": username,
		"message":  message,
	})
}

// 用户登录
func UserLogin(pContext *gin.Context) {
	name := strings.TrimSpace(pContext.PostForm("name"))
	passwd := strings.TrimSpace(pContext.PostForm("passwd"))

	err := VerifyUserNameAndPasswd(name, passwd)

	var user User
	if err == nil {
		err = db.Qry(&user, "SELECT u.id, u.`name`, u.nickname, u.rem, u.create_time, u.update_time FROM `user` u WHERE u.del_flag = 0 AND u.`name` = ? AND u.passwd = ? LIMIT 1", name, passwd)
	}

	if err == nil && user.Id == 0 {
		err = errors.New(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
	}

	// 初始化session对象
	session := sessions.Default(pContext)

	if err != nil {
		session.Set("username", name)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")
		return
	}

	// session
	// 设置session数据
	session.Set(SessionKeyUser, user)
	// 保存session数据
	session.Save()

	// 重定向
	pContext.Redirect(http.StatusMovedPermanently, "/")
}

// 用户登出
func UserLogout(pContext *gin.Context) {
	// 解析session
	session := sessions.Default(pContext)

	// 清除session
	session.Clear()

	// 保存session数据
	session.Save()

	// 重定向
	pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")
}

func UserStgPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	user := session.Get("user")
	message := session.Get("message")
	session.Delete("user")
	session.Delete("message")
	session.Save()

	if user == nil {
		user = GetUser(pContext)
	}
	pContext.HTML(http.StatusOK, "user/stg.html", gin.H{
		"user":    user,
		"message": message,
	})
}

func UserUpd(pContext *gin.Context) {
	user := User{}
	err := ShouldBind(pContext, &user)

	if err == nil {
		err = VerifyUserNameAndPasswd(user.Name, user.Passwd)
	}

	sessionUser := GetUser(pContext)
	if err == nil && user.Name != sessionUser.Name {
		err = VerifyDbUserName(user.Name)
	}

	user.Nickname = strings.TrimSpace(user.Nickname)
	user.Rem = strings.TrimSpace(user.Rem)

	session := sessions.Default(pContext)
	if err != nil {
		session.Set("user", user)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/stgpage")
		return
	}

	db.Add("UPDATE `user` SET `name` = ?, nickname = ?, `passwd` = ?, rem = ?, update_time = ? WHERE id = ?",
		user.Name, user.Nickname, user.Passwd, user.Rem, time.Now().Unix(), sessionUser.Id)

	// 更新session中User信息
	sessionUser.Name = user.Name
	sessionUser.Nickname = user.Nickname
	sessionUser.Rem = user.Rem
	session.Set(SessionKeyUser, sessionUser)
	session.Save()

	pContext.Redirect(http.StatusMovedPermanently, "/user/stgpage")
}

func VerifyUserNameAndPasswd(name, passwd string) error {
	err := com.VerifyUserName(name)
	if err == nil {
		err = com.VerifyPasswd(passwd)
	}
	return err
}

func VerifyDbUserName(name string) error {
	var id int64
	err := db.Qry(&id, "SELECT u.id FROM `user` u WHERE u.del_flag = 0 AND u.`name` = ? LIMIT 1", name)
	if err != nil {
		return err
	}

	if id != 0 {
		return errors.New(i18n.MustGetMessage("i18n.userAlreadyExists"))
	}

	return nil
}

func GetUser(pContext *gin.Context) User {
	session := sessions.Default(pContext)
	var user User
	if v, r := session.Get(SessionKeyUser).(User); r {
		user = v
	}

	// 如果返回指针值，有可能会发生逃逸
	//return &user

	return user
}
