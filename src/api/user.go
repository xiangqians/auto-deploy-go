// user
// @author xiangqian
// @date 18:10 2022/12/18
package api

import (
	"auto-deploy-go/src/db"
	"encoding/gob"
	"errors"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

const SessionKeyUser = "_user_"

type User struct {
	Id         int64  // 主键id
	Name       string // 用户名
	Nickname   string // 昵称
	Passwd     string // 密码
	Rem        string // 备注
	DelFlag    byte   // 删除标识，0-正常，1-删除
	CreateTime int64  // 创建时间（时间戳，s）
	UpdateTime int64  // 修改时间（时间戳，s）
}

func init() {
	// 注册 User 模型
	gob.Register(User{})
}

// 用户注册html
func UserRegPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	username := session.Get("username")
	nickname := session.Get("nickname")
	rem := session.Get("rem")
	message := session.Get("message")
	session.Delete("username")
	session.Delete("nickname")
	session.Delete("rem")
	session.Delete("message")
	session.Save()

	pContext.HTML(http.StatusOK, "user/reg.html", gin.H{
		"username": username,
		"nickname": nickname,
		"rem":      rem,
		"message":  message,
	})
}

// 用户注册
func UserReg(pContext *gin.Context) {
	name := strings.TrimSpace(pContext.PostForm("name"))
	nickname := strings.TrimSpace(pContext.PostForm("nickname"))
	rem := strings.TrimSpace(pContext.PostForm("rem"))
	err := VerifyUserName(name)
	if err != nil {
		session := sessions.Default(pContext)
		session.Set("username", name)
		session.Set("nickname", nickname)
		session.Set("rem", rem)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/regpage")
		return
	}

	passwd := strings.TrimSpace(pContext.PostForm("passwd"))
	db.Add("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`, `create_time`) VALUES (?, ?, ?, ?, ?)",
		name, nickname, passwd, rem, time.Now().Unix())

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

	// 初始化session对象
	session := sessions.Default(pContext)

	var user User
	err := db.Qry(&user, "SELECT u.id, u.`name`, u.nickname, u.rem, u.create_time, u.update_time FROM `user` u WHERE u.del_flag = 0 AND u.`name` = ? AND u.passwd = ? LIMIT 1", name, passwd)
	if err != nil {
		log.Fatalln(err)
	}

	if user.Id == 0 {
		session.Set("username", name)
		session.Set("message", i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
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
	user := GetUser(pContext)
	pContext.HTML(http.StatusOK, "user/stg.html", gin.H{
		"username": user.Name,
		"nickname": user.Nickname,
		"rem":      user.Rem,
	})
}

func UserStgUpd(pContext *gin.Context) {
	name := strings.TrimSpace(pContext.PostForm("name"))
	nickname := strings.TrimSpace(pContext.PostForm("nickname"))
	rem := strings.TrimSpace(pContext.PostForm("rem"))
	user := GetUser(pContext)

	var err error = nil
	if user.Name != name {
		err = VerifyUserName(name)
	}
	if err != nil {
		session := sessions.Default(pContext)
		session.Set("username", name)
		session.Set("nickname", nickname)
		session.Set("rem", rem)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/regpage")
		return
	}

	passwd := strings.TrimSpace(pContext.PostForm("passwd"))
	db.Add("UPDATE `user` SET `name` = ?, nickname = ?, `passwd` = ?, rem = ?, update_time = ? WHERE id = ?",
		name, nickname, passwd, rem, time.Now().Unix(), user.Id)

	// 更新session中User信息
	session := sessions.Default(pContext)
	user.Name = name
	user.Nickname = nickname
	user.Rem = rem
	session.Set(SessionKeyUser, user)
	session.Save()

	pContext.Redirect(http.StatusMovedPermanently, "/user/stgpage")
}

func VerifyUserName(name string) error {
	if name == "" {
		return errors.New(i18n.MustGetMessage("i18n.userCannotEmpty"))
	}

	var id int64
	err := db.Qry(&id, "SELECT u.id FROM `user` u WHERE u.`name` = ? LIMIT 1", name)
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
