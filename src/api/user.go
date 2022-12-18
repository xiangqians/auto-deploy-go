// user
// @author xiangqian
// @date 18:10 2022/12/18
package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Name     string `json:"name"`     // 用户名
	Nickname string `json:"nickname"` // 昵称
	Passwd   string `json:"passwd"`   // 密码
}

// admin

// 用户注册
func UserReg(pContext *gin.Context) {
}

// 用户登录html
func UserLoginHtml(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "login.html", gin.H{
		"test": "test",
	})
}

// 用户登录
func UserLogin(pContext *gin.Context) {
	// 初始化session对象
	session := sessions.Default(pContext)

	// 设置session数据
	session.Set("username", "admin")

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
