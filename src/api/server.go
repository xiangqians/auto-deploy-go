// Server
// @author xiangqian
// @date 15:45 2022/12/22
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ServerIndex(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "server/index.html", gin.H{
		"user":    GetUser(pContext),
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
		_server := typ.Server{}
		idStr := pContext.Query("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil && id > 0 {
			user := GetUser(pContext)
			_, err = db.Qry(&_server, "SELECT s.id, s.`name`, s.`host`, s.`port`, s.`user`, s.rem, s.add_time, s.upd_time FROM server s WHERE s.del_flag = 0 AND s.user_id = ? AND s.id = ?", user.Id, id)
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
	db.Add("INSERT INTO `server` (`user_id`, `name`, `host`, `port`, `user`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, server.Name, server.Host, server.Port, server.User, server.Passwd, server.Rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}

func ServerUpd(pContext *gin.Context) {
	server, err := serverPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Upd("UPDATE `server` SET `name` = ?, `host` = ?, `port` = ?, `user` = ?, `passwd` = ?, `rem` = ?, upd_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
		server.Name, server.Host, server.Port, server.User, server.Passwd, server.Rem, time.Now().Unix(), user.Id, server.Id)
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}

func ServerDel(pContext *gin.Context) {
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil {
		user := GetUser(pContext)
		db.Del("UPDATE server SET del_flag = 1, upd_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	}
	pContext.Redirect(http.StatusMovedPermanently, "/server/index")
}

func serverPreAddOrUpd(pContext *gin.Context) (typ.Server, error) {
	server := typ.Server{}
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

func Servers(pContext *gin.Context) []typ.Server {
	user := GetUser(pContext)
	servers := make([]typ.Server, 1)
	_, err := db.Qry(&servers, "SELECT s.id, s.`name`, s.`host`, s.`port`, s.`user`, s.rem, s.add_time, s.upd_time FROM server s WHERE s.del_flag = 0 AND s.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if servers[0].Id == 0 {
		servers = nil
	}
	return servers
}
