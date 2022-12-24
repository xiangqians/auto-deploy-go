// Server
// @author xiangqian
// @date 15:45 2022/12/22
package api

import (
	"auto-deploy-go/src/db"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Abs
	UserId int64  // Server所属用户id
	Name   string `form:"name" binding:"required,min=1,max=60"`    // 名称
	Host   string `form:"host" binding:"required,min=1,max=60"`    // host
	Port   int    `form:"port" binding:"required,gt=0"`            // 端口
	User   string `form:"user" binding:"required,min=1,max=60"`    // 用户
	Passwd string `form:"passwd" binding:"required,min=1,max=100"` // 密码
}

func init() {
	// 注册 Server 模型
	gob.Register(Server{})
}

func ServerIndex(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "server/index.html", gin.H{
		"servers": Servers(pContext),
	})
}

func ServerAddPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	server := session.Get("server")
	message := session.Get("message")
	session.Delete("server")
	session.Delete("message")
	session.Save()

	if server == nil {
		_server := Server{}
		idStr := pContext.Query("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil && id > 0 {
			user := GetUser(pContext)
			err = db.Qry(&_server, "SELECT s.id, s.`name`, s.`host`, s.`port`, s.`user`, s.rem, s.create_time, s.update_time FROM server s WHERE s.del_flag = 0 AND s.user_id = ? AND s.id = ?", user.Id, id)
			if err != nil {
				log.Println(err)
			}
		}
		server = _server
	}

	pContext.HTML(http.StatusOK, "server/add.html", gin.H{
		"server":  server,
		"message": message,
	})
}

func ServerAdd(pContext *gin.Context) {
	server, err := serverPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Add("INSERT INTO `server` (`user_id`, `name`, `host`, `port`, `user`, `passwd`, `rem`, `create_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, server.Name, server.Host, server.Port, server.User, server.Passwd, server.Rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}

func ServerUpd(pContext *gin.Context) {
	server, err := serverPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Upd("UPDATE `server` SET `name` = ?, `host` = ?, `port` = ?, `user` = ?, `passwd` = ?, `rem` = ?, update_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
		server.Name, server.Host, server.Port, server.User, server.Passwd, server.Rem, time.Now().Unix(), user.Id, server.Id)
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}

func ServerDel(pContext *gin.Context) {
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil {
		user := GetUser(pContext)
		db.Del("UPDATE server SET del_flag = 1, update_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	}
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}

func serverPreAddOrUpd(pContext *gin.Context) (Server, error) {
	server := Server{}
	err := ShouldBind(pContext, &server)

	server.Name = strings.TrimSpace(server.Name)
	server.Host = strings.TrimSpace(server.Host)
	server.User = strings.TrimSpace(server.User)
	server.Passwd = strings.TrimSpace(server.Passwd)
	server.Rem = strings.TrimSpace(server.Rem)

	if err != nil {
		session := sessions.Default(pContext)
		session.Set("server", server)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/server/addpage")
	}

	return server, err
}

func Servers(pContext *gin.Context) []Server {
	user := GetUser(pContext)
	servers := make([]Server, 1)
	err := db.Qry(&servers, "SELECT s.id, s.`name`, s.`host`, s.`port`, s.`user`, s.rem, s.create_time, s.update_time FROM server s WHERE s.del_flag = 0 AND s.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if servers[0].Id == 0 {
		servers = nil
	}
	return servers
}
