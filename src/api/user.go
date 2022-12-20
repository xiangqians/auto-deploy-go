// user
// @author xiangqian
// @date 18:10 2022/12/18
package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

const SessionKeyUsername = "_username_"

const userJsonFname = "./data/user/user.json"

type Server struct {
	Name   string `json:"name"`   // 名称，唯一
	Desc   string `json:"desc"`   // 服务描述
	Host   string `json:"host"`   // host
	Port   int    `json:"port"`   // 端口，默认22
	User   string `json:"user"`   // 用户名
	Passwd string `json:"passwd"` // 密码
}

type Git struct {
	Name   string `json:"name"`   // 名称，唯一
	Desc   string `json:"desc"`   // 服Git描述
	User   string `json:"user"`   // 用户名
	Passwd string `json:"passwd"` // 密码
}

type User struct {
	Name     string   `json:"name"`     // 用户名
	Nickname string   `json:"nickname"` // 昵称
	Passwd   string   `json:"passwd"`   // 密码
	Servers  []Server `json:"servers"`  // 用户所拥有的服务列表
	Gits     []Git    `json:"gits"`     // 用户所拥有的Git列表
}

// configuration

var users []User

func init() {
	// 读取json文件
	pFile, err := os.Open(userJsonFname)
	if err != nil {
		panic(err)
	}
	defer pFile.Close()

	pDecoder := json.NewDecoder(pFile)
	err = pDecoder.Decode(&users)
	if err != nil {
		panic(err)
	}

	log.Printf("users = %v\n", users)
}

// 用户注册html
func UserRegPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	username := session.Get("username")
	nickname := session.Get("nickname")
	message := session.Get("message")
	session.Delete("username")
	session.Delete("nickname")
	session.Delete("message")
	session.Save()

	pContext.HTML(http.StatusOK, "user/reg.html", gin.H{
		"username": username,
		"nickname": nickname,
		"message":  message,
	})
}

// 用户注册
func UserReg(pContext *gin.Context) {
	name := strings.TrimSpace(pContext.PostForm("name"))
	nickname := strings.TrimSpace(pContext.PostForm("nickname"))
	err := VerifyUserName(name)
	if err != nil {
		session := sessions.Default(pContext)
		session.Set("username", name)
		session.Set("nickname", nickname)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/reg")
		return
	}

	passwd := strings.TrimSpace(pContext.PostForm("passwd"))
	users = append(users, User{
		Name:     name,
		Nickname: nickname,
		Passwd:   passwd,
	})

	// 将用户信息序列化到本地
	FlushUser()

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
	verify := false
	for _, user := range users {
		if user.Name == name && user.Passwd == passwd {
			verify = true
			break
		}
	}

	// 初始化session对象
	session := sessions.Default(pContext)

	if !verify {
		session.Set("username", name)
		session.Set("message", i18n.MustGetMessage("i18n.usernameOrPasswordIncorrect"))
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")
		return
	}

	// 设置session数据
	session.Set(SessionKeyUsername, name)

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

func UserAccountPage(pContext *gin.Context) {
	pUser := GetUser(pContext)
	username := pUser.Name
	nickname := pUser.Nickname
	pContext.HTML(http.StatusOK, "user/settings.html", gin.H{
		"username": username,
		"nickname": nickname,
		"type":     "account",
	})
}

func UserAccountUpd(pContext *gin.Context) {
	nickname := strings.TrimSpace(pContext.PostForm("nickname"))
	passwd := strings.TrimSpace(pContext.PostForm("passwd"))
	pUser := GetUser(pContext)
	pUser.Nickname = nickname
	pUser.Passwd = passwd
	FlushUser()

	// 个人信息修改成功后重定向到当前页面
	pContext.Redirect(http.StatusMovedPermanently, "/user/accountpage")
}

func UserGitPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	name := session.Get("name")
	desc := session.Get("desc")
	user := session.Get("user")
	message := session.Get("message")
	session.Delete("name")
	session.Delete("desc")
	session.Delete("user")
	session.Delete("message")
	session.Save()

	pUser := GetUser(pContext)
	pContext.HTML(http.StatusOK, "user/settings.html", gin.H{
		"gits":    pUser.Gits,
		"name":    name,
		"desc":    desc,
		"user":    user,
		"message": message,
		"type":    "git",
	})
}

func UserGitAdd(pContext *gin.Context) {
	//type Git struct {
	//	Name   string `json:"name"`   // 名称，唯一
	//	Desc   string `json:"desc"`   // 服Git描述
	//	User   string `json:"user"`   // 用户名
	//	Passwd string `json:"passwd"` // 密码
	//}
	//
	name := strings.TrimSpace(pContext.PostForm("name"))
	desc := strings.TrimSpace(pContext.PostForm("desc"))
	user := strings.TrimSpace(pContext.PostForm("user"))
	passwd := strings.TrimSpace(pContext.PostForm("passwd"))

	errI18nKey := ""
	if name == "" {
		errI18nKey = "i18n.nameCannotEmpty"
	} else if user == "" {
		errI18nKey = "i18n.usernameCannotEmpty"
	} else if passwd == "" {
		errI18nKey = "i18n.passwordCannotEmpty"
	}

	pUser := GetUser(pContext)
	if errI18nKey == "" && pUser.Gits != nil {
		for _, v := range pUser.Gits {
			if v.Name == name {
				errI18nKey = "i18n.nameAlreadyExists"
				break
			}
		}
	}

	if errI18nKey != "" {
		session := sessions.Default(pContext)
		session.Set("name", name)
		session.Set("desc", desc)
		session.Set("user", user)
		session.Set("message", i18n.MustGetMessage(errI18nKey))
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/gitpage")
		return
	}

	// save
	pUser.Gits = append(pUser.Gits, Git{
		Name:   name,
		Desc:   desc,
		User:   user,
		Passwd: passwd,
	})
	FlushUser()

	// 个人信息修改成功后重定向到当前页面
	pContext.Redirect(http.StatusMovedPermanently, "/user/gitpage")
}

func UserGitUpd(pContext *gin.Context) {

}

func UserGitDel(pContext *gin.Context) {

}

func VerifyUserName(name string) error {
	if name == "" {
		return errors.New(i18n.MustGetMessage("i18n.usernameCannotEmpty"))
	}

	for _, user := range users {
		if user.Name == name {
			return errors.New(i18n.MustGetMessage("i18n.usernameAlreadyExists"))
		}
	}

	return nil
}

func FlushUser() {
	// 将用户信息序列化到本地
	pFile, err := os.Create(userJsonFname)
	if err != nil {
		log.Fatal(err)
	}

	pEncoder := json.NewEncoder(pFile)
	err = pEncoder.Encode(users)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUser(pContext *gin.Context) *User {
	session := sessions.Default(pContext)
	username := ""
	if v, r := session.Get(SessionKeyUsername).(string); r && v != "" {
		username = v
	}

	for i, l := 0, len(users); i < l; i++ {
		user := users[i]
		if user.Name == username {
			// 注意两者区别，否则无法修改 user 信息
			//return &user
			return &users[i]
		}
	}

	return nil
}
