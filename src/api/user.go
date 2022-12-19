// user
// @author xiangqian
// @date 18:10 2022/12/18
package api

import (
	"encoding/json"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

const SessionKeyUsername = "_username_"

const userJsonFileName = "./data/user.json"

type User struct {
	Name     string `json:"name"`     // 用户名
	Nickname string `json:"nickname"` // 昵称
	Passwd   string `json:"passwd"`   // 密码
}

var users []User

func init() {
	// 读取json文件
	pFile, err := os.Open(userJsonFileName)
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
func UserRegHtml(pContext *gin.Context) {
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
	i18nKey := ""
	verify := true
	if name == "" {
		verify = false
		i18nKey = "i18n.usernameCannotEmpty"
	} else {
		for _, user := range users {
			if user.Name == name {
				verify = false
				i18nKey = "i18n.usernameAlreadyExists"
				break
			}
		}
	}

	if !verify {
		session := sessions.Default(pContext)
		session.Set("username", name)
		session.Set("nickname", nickname)
		if i18nKey != "" {
			session.Set("message", i18n.MustGetMessage(i18nKey))
		}
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/user/reg")
		return
	}

	passwd := strings.TrimSpace(pContext.PostForm("passwd"))

	if nickname == "" {
		nickname = name
	}

	users = append(users, User{
		Name:     name,
		Nickname: nickname,
		Passwd:   passwd,
	})

	// 将用户信息序列化到本地
	pFile, err := os.Create(userJsonFileName)
	if err != nil {
		log.Fatal(err)
	}
	pEncoder := json.NewEncoder(pFile)
	err = pEncoder.Encode(users)
	if err != nil {
		log.Fatal(err)
	}

	// 用户注册成功后，重定向到登录页
	pContext.Redirect(http.StatusMovedPermanently, "/user/login")
}

// 用户登录html
func UserLoginHtml(pContext *gin.Context) {
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
		pContext.Redirect(http.StatusMovedPermanently, "/user/login")
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
	pContext.Redirect(http.StatusMovedPermanently, "/user/login")
}

func UserStgHtml(pContext *gin.Context) {

}
func UserStg(pContext *gin.Context) {

}

func GetUser(pContext *gin.Context) *User {
	session := sessions.Default(pContext)
	username := ""
	if v, r := session.Get(SessionKeyUsername).(string); r && v != "" {
		username = v
	}

	var user *User = nil
	for _, u := range users {
		if u.Name == username {
			user = &u
			break
		}
	}

	return user
}
