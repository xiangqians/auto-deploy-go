// Git
// @author xiangqian
// @date 11:48 2022/12/22
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

type Git struct {
	Abs
	UserId int64  // Git所属用户id
	Name   string // 名称
	User   string // 用户
	Passwd string // 密码
}

func GitIndex(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "git/index.html", gin.H{
		"gits": Gits(pContext),
	})
}

func GitAddPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	name := session.Get("name")
	user := session.Get("user")
	rem := session.Get("rem")
	message := session.Get("message")
	session.Delete("name")
	session.Delete("rem")
	session.Delete("message")
	session.Save()

	pContext.HTML(http.StatusOK, "git/add.html", gin.H{
		"name":    name,
		"user":    user,
		"rem":     rem,
		"message": message,
	})
}

func GitAdd(pContext *gin.Context) {
	name := strings.TrimSpace(pContext.PostForm("name"))
	_user := strings.TrimSpace(pContext.PostForm("user"))
	passwd := strings.TrimSpace(pContext.PostForm("passwd"))
	rem := strings.TrimSpace(pContext.PostForm("rem"))
	user := GetUser(pContext)
	db.Add("INSERT INTO `git` (`user_id`, `name`, `user`, `passwd`, `rem`, `create_time`) VALUES (?, ?, ?, ?, ?, ?)",
		user.Id, name, _user, passwd, rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/git/index")
}

func GitUpd(pContext *gin.Context) {
	pContext.Redirect(http.StatusMovedPermanently, "/git/index")
}

func GitDel(pContext *gin.Context) {
	pContext.Redirect(http.StatusMovedPermanently, "/git/index")
}

func Gits(pContext *gin.Context) []Git {
	user := GetUser(pContext)
	gits := make([]Git, 1)
	err := db.Qry(&gits, "SELECT g.id, g.`name`, g.`user`, g.rem, g.create_time, g.update_time FROM git g WHERE g.del_flag = 0 AND g.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if gits[0].Id == 0 {
		gits = nil
	}

	return gits
}
