// Git
// @author xiangqian
// @date 11:48 2022/12/22
package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func UserGitPage(pContext *gin.Context) {
	//session := sessions.Default(pContext)
	//name := session.Get("name")
	//desc := session.Get("desc")
	//user := session.Get("user")
	//message := session.Get("message")
	//session.Delete("name")
	//session.Delete("desc")
	//session.Delete("user")
	//session.Delete("message")
	//session.Save()
	//
	//pUser := GetUser(pContext)
	//pContext.HTML(http.StatusOK, "user/stg.html", gin.H{
	//	"gits":    pUser.Gits,
	//	"name":    name,
	//	"desc":    desc,
	//	"user":    user,
	//	"message": message,
	//	"type":    "git",
	//})
}

func UserGitAdd(pContext *gin.Context) {
	//name := strings.TrimSpace(pContext.PostForm("name"))
	//desc := strings.TrimSpace(pContext.PostForm("desc"))
	//user := strings.TrimSpace(pContext.PostForm("user"))
	//passwd := strings.TrimSpace(pContext.PostForm("passwd"))
	//
	//errI18nKey := ""
	//if name == "" {
	//	errI18nKey = "i18n.nameCannotEmpty"
	//} else if user == "" {
	//	errI18nKey = "i18n.userCannotEmpty"
	//} else if passwd == "" {
	//	errI18nKey = "i18n.passwdCannotEmpty"
	//}
	//
	//pUser := GetUser(pContext)
	//if errI18nKey == "" && pUser.Gits != nil {
	//	for _, v := range pUser.Gits {
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
	//	session.Set("desc", desc)
	//	session.Set("user", user)
	//	session.Set("message", i18n.MustGetMessage(errI18nKey))
	//	session.Save()
	//	pContext.Redirect(http.StatusMovedPermanently, "/user/gitpage")
	//	return
	//}
	//
	//// save
	//pUser.Gits = append(pUser.Gits, Git{
	//	Name:   name,
	//	Desc:   desc,
	//	User:   user,
	//	Passwd: passwd,
	//})
	//FlushUser()

	pContext.Redirect(http.StatusMovedPermanently, "/user/gitpage")
}

func UserGitUpd(pContext *gin.Context) {
}

func UserGitDel(pContext *gin.Context) {

}

func UserServerPage(pContext *gin.Context) {
	//session := sessions.Default(pContext)
	//name := session.Get("name")
	//desc := session.Get("desc")
	//host := session.Get("host")
	//port := session.Get("port")
	//user := session.Get("user")
	//message := session.Get("message")
	//session.Delete("name")
	//session.Delete("desc")
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
	//	"desc":    desc,
	//	"host":    host,
	//	"port":    port,
	//	"user":    user,
	//	"message": message,
	//	"type":    "server",
	//})
}

func UserServerAdd(pContext *gin.Context) {
	//name := strings.TrimSpace(pContext.PostForm("name"))
	//desc := strings.TrimSpace(pContext.PostForm("desc"))
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
	//	session.Set("desc", desc)
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
	//	Desc:   desc,
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
