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

type Server struct {
	Name   string `json:"name"`   // 名称，唯一
	Host   string `json:"host"`   // host
	Port   int    `json:"port"`   // 端口，默认22
	User   string `json:"user"`   // 用户名
	Passwd string `json:"passwd"` // 密码
}

type Git struct {
	Abs
	UserId int64  // Git所属用户id
	Name   string // 名称
	User   string // 用户
	Passwd string // 密码
}

func GitIndex(pContext *gin.Context) {
	user := GetUser(pContext)
	gits := make([]Git, 1)
	err := db.Qry(&gits, "select g.id, g.`name`, g.`user`, g.rem, g.create_time, g.update_time from git g where g.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return
	}

	if gits[0].Id == 0 {
		gits = nil
	}

	pContext.HTML(http.StatusOK, "git/index.html", gin.H{
		"gits": gits,
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

func UserServerPage(pContext *gin.Context) {
	//session := sessions.Default(pContext)
	//name := session.Get("name")
	//rem := session.Get("rem")
	//host := session.Get("host")
	//port := session.Get("port")
	//user := session.Get("user")
	//message := session.Get("message")
	//session.Delete("name")
	//session.Delete("rem")
	//session.Delete("host")
	//session.Delete("port")
	//session.Delete("user")
	//session.Delete("message")
	//session.Save()
	//
	//pUser := GetUser(pContext)
	//pContext.HTML(http.StatusOK, "user/stg.html", gin.H{
	//	"servers": pUser.Servers,
	//	"name":    name,
	//	"rem":    rem,
	//	"host":    host,
	//	"port":    port,
	//	"user":    user,
	//	"message": message,
	//	"type":    "server",
	//})
}

func UserServerAdd(pContext *gin.Context) {
	//name := strings.TrimSpace(pContext.PostForm("name"))
	//rem := strings.TrimSpace(pContext.PostForm("rem"))
	//host := strings.TrimSpace(pContext.PostForm("host"))
	//port := strings.TrimSpace(pContext.PostForm("port"))
	//user := strings.TrimSpace(pContext.PostForm("user"))
	//passwd := strings.TrimSpace(pContext.PostForm("passwd"))

	//errI18nKey := ""
	//if name == "" {
	//	errI18nKey = "i18n.nameCannotEmpty"
	//} else if host == "" {
	//	errI18nKey = "i18n.hostCannotEmpty"
	//} else if port == "" {
	//	errI18nKey = "i18n.portCannotEmpty"
	//} else if user == "" {
	//	errI18nKey = "i18n.userCannotEmpty"
	//} else if passwd == "" {
	//	errI18nKey = "i18n.passwdCannotEmpty"
	//}

	//pUser := GetUser(pContext)
	//if errI18nKey == "" && pUser.Servers != nil {
	//	for _, v := range pUser.Servers {
	//		if v.Name == name {
	//			errI18nKey = "i18n.nameAlreadyExists"
	//			break
	//		}
	//	}
	//}
	//
	//if errI18nKey != "" {
	//	session := sessions.Default(pContext)
	//	session.Set("name", name)
	//	session.Set("rem", rem)
	//	session.Set("host", host)
	//	session.Set("port", port)
	//	session.Set("user", user)
	//	session.Set("message", i18n.MustGetMessage(errI18nKey))
	//	session.Save()
	//	pContext.Redirect(http.StatusMovedPermanently, "/user/serverpage")
	//	return
	//}
	//
	//iPort, err := strconv.Atoi(port)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// save
	//pUser.Servers = append(pUser.Servers, Server{
	//	Name:   name,
	//	rem:   rem,
	//	Host:   host,
	//	Port:   iPort,
	//	User:   user,
	//	Passwd: passwd,
	//})
	//FlushUser()

	pContext.Redirect(http.StatusMovedPermanently, "/user/serverpage")
}

func UserServerUpd(pContext *gin.Context) {

}

func UserServerDel(pContext *gin.Context) {

}
