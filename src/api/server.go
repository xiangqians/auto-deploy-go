// Server
// @author xiangqian
// @date 15:45 2022/12/22
package api

import (
	"auto-deploy-go/src/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Abs
	UserId int64  // Server所属用户id
	Name   string // 名称
	Host   string // host
	Port   int    // 端口
	User   string // 用户
	Passwd string // 密码
}

func ServerIndex(pContext *gin.Context) {
	user := GetUser(pContext)
	servers := make([]Server, 1)
	err := db.Qry(&servers, "SELECT s.id, s.`name`, s.`host`, s.`port`, s.`user`, s.rem, s.create_time, s.update_time FROM server s WHERE s.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return
	}

	if servers[0].Id == 0 {
		servers = nil
	}

	pContext.HTML(http.StatusOK, "server/index.html", gin.H{
		"servers": servers,
	})
}

func ServerAddPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	name := session.Get("name")
	host := session.Get("host")
	port := session.Get("port")
	user := session.Get("user")
	rem := session.Get("rem")
	message := session.Get("message")
	session.Delete("name")
	session.Delete("host")
	session.Delete("port")
	session.Delete("user")
	session.Delete("rem")
	session.Delete("message")
	session.Save()

	pContext.HTML(http.StatusOK, "server/add.html", gin.H{
		"name":    name,
		"host":    host,
		"port":    port,
		"user":    user,
		"rem":     rem,
		"message": message,
	})
}

func ServerAdd(pContext *gin.Context) {
	name := strings.TrimSpace(pContext.PostForm("name"))
	host := strings.TrimSpace(pContext.PostForm("host"))
	port := strings.TrimSpace(pContext.PostForm("port"))
	_user := strings.TrimSpace(pContext.PostForm("user"))
	passwd := strings.TrimSpace(pContext.PostForm("passwd"))
	rem := strings.TrimSpace(pContext.PostForm("rem"))
	user := GetUser(pContext)
	db.Add("INSERT INTO `server` (`user_id`, `name`, `host`, `port`, `user`, `passwd`, `rem`, `create_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, name, host, port, _user, passwd, rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}

func ServerUpd(pContext *gin.Context) {
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}

func ServerDel(pContext *gin.Context) {
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}
