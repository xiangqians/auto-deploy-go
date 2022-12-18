// user
// @author xiangqian
// @date 18:10 2022/12/18
package api

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

type User struct {
	Name     string `json:"name"`     // 用户名
	Nickname string `json:"nickname"` // 昵称
	Passwd   string `json:"passwd"`   // 密码
}

var users []User

func init() {
	// 读取json文件
	pFile, err := os.Open("./data/user.json")
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

// 用户注册
func UserReg(pContext *gin.Context) {
}

// 用户登录html
func UserLoginHtml(pContext *gin.Context) {
	session := sessions.Default(pContext)
	pContext.HTML(http.StatusOK, "login.html", gin.H{
		"username": session.Get("username"),
		"message":  session.Get("message"),
	})
	session.Delete("username")
	session.Delete("message")
	session.Save()
}

// 用户登录
func UserLogin(pContext *gin.Context) {
	name := pContext.PostForm("name")
	passwd := pContext.PostForm("passwd")
	var user *User = nil
	for _, u := range users {
		if u.Name == name && u.Passwd == passwd {
			user = &u
		}
	}

	// 初始化session对象
	session := sessions.Default(pContext)

	if user == nil {
		session.Set("username", name)
		session.Set("message", "用户名或密码错误")
		session.Save()
		pContext.Redirect(http.StatusFound, "/user/login")
		return
	}

	// 设置session数据
	session.Set("_username", name)

	// 保存session数据
	session.Save()

	// 重定向
	pContext.Redirect(http.StatusFound, "/")
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
	pContext.Redirect(http.StatusFound, "/user/login")
}
