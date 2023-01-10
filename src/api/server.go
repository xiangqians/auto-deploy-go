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
	HtmlPage(pContext, "server/index.html", func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error) {
		page, err := PageServer(pContext, pageReq)
		return page, nil, err
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
	ServerAddOrUpd(pContext)
}

func ServerUpd(pContext *gin.Context) {
	ServerAddOrUpd(pContext)
}

func ServerAddOrUpd(pContext *gin.Context) {
	redirect := func(server typ.Server, message any) {
		Redirect(pContext, "/server/addpage", message, map[string]any{"server": server})
	}

	server := typ.Server{}
	err := ShouldBind(pContext, &server)
	server.Name = strings.TrimSpace(server.Name)
	server.Host = strings.TrimSpace(server.Host)
	server.User = strings.TrimSpace(server.User)
	server.Passwd = strings.TrimSpace(server.Passwd)
	server.Rem = strings.TrimSpace(server.Rem)
	if err != nil {
		redirect(server, err)
		return
	}

	user := GetUser(pContext)
	if pContext.Request.Method == http.MethodPost {
		_, err = db.Add("INSERT INTO `server` (`user_id`, `name`, `host`, `port`, `user`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			user.Id, server.Name, server.Host, server.Port, server.User, server.Passwd, server.Rem, time.Now().Unix())
	} else if pContext.Request.Method == http.MethodPut {
		_, err = db.Upd("UPDATE `server` SET `name` = ?, `host` = ?, `port` = ?, `user` = ?, `passwd` = ?, `rem` = ?, upd_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
			server.Name, server.Host, server.Port, server.User, server.Passwd, server.Rem, time.Now().Unix(), user.Id, server.Id)
	}

	Redirect(pContext, "/server/index", err, nil)
	return
}

func ServerDel(pContext *gin.Context) {
	redirect := func(message any) {
		Redirect(pContext, "/server/index", message, nil)
	}

	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		redirect(err)
		return
	}

	user := GetUser(pContext)
	db.Del("UPDATE server SET del_flag = 1, upd_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	redirect(nil)
	return
}

func PageServer(pContext *gin.Context, pageReq typ.PageReq) (typ.Page[typ.Server], error) {
	user := GetUser(pContext)
	return db.Page[typ.Server](pageReq, "SELECT s.id, s.`name`, s.`host`, s.`port`, s.`user`, s.rem, s.add_time, s.upd_time FROM server s WHERE s.del_flag = 0 AND s.user_id = ?", user.Id)
}
