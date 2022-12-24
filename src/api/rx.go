// rx
// @author xiangqian
// @date 22:16 2022/12/24
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

type Rx struct {
	Abs
	Name       string `form:"name" binding:"required,min=1,max=60"` // 名称
	OwnerId    int64  // 拥有者id
	OwnerName  string // 拥有者名称
	SharerId   int64  // 共享者id
	SharerName string // 共享者名称
}

func init() {
	// 注册 Rx 模型
	gob.Register(Rx{})
}

func RxIndex(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "rx/index.html", gin.H{
		"rxs": Rxs(pContext),
	})
}

func RxAddPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	rx := session.Get("rx")
	message := session.Get("message")
	session.Delete("rx")
	session.Delete("message")
	session.Save()

	if rx == nil {
		rx = Rx{}
	}
	pContext.HTML(http.StatusOK, "rx/add.html", gin.H{
		"rx":      rx,
		"message": message,
	})
}

func RxAdd(pContext *gin.Context) {
	rx := Rx{}
	err := ShouldBind(pContext, &rx)

	rx.Name = strings.TrimSpace(rx.Name)
	rx.Rem = strings.TrimSpace(rx.Rem)

	if err != nil {
		session := sessions.Default(pContext)
		session.Set("rx", rx)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/rx/addpage")
		return
	}

	user := GetUser(pContext)
	db.Add("INSERT INTO `rx` (`name`, `owner_id`, `rem`, `create_time`) VALUES (?, ?, ?, ?)",
		rx.Name, user.Id, rx.Rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/rx/index")
}

func RxJoin(pContext *gin.Context) {
	codeStr := pContext.PostForm("code")
	code, err := strconv.ParseInt(codeStr, 10, 64)
	if err == nil {
		user := GetUser(pContext)
		db.Del("UPDATE rx SET sharer_id = ?, update_time = ? WHERE del_flag = 0 AND sharer_id = 0 AND owner_id != ? AND id = ?",
			user.Id, time.Now().Unix(), user.Id, code)
	}
	pContext.Redirect(http.StatusMovedPermanently, "/rx/index")
}

func RxDel(pContext *gin.Context) {
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil {
		user := GetUser(pContext)
		db.Del("UPDATE rx SET del_flag = 1, update_time = ? WHERE (owner_id = ? OR sharer_id = ?) AND id = ?", time.Now().Unix(), user.Id, user.Id, id)
	}
	pContext.Redirect(http.StatusMovedPermanently, "/rx/index")
}

func Rxs(pContext *gin.Context) []Rx {
	user := GetUser(pContext)
	rxs := make([]Rx, 1)
	err := db.Qry(&rxs, " SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', r.rem, r.create_time, r.update_time FROM rx r LEFT JOIN user ou ON ou.del_flag = 0 AND ou.id = r.owner_id LEFT JOIN user su ON su.del_flag = 0 AND su.id = r.sharer_id WHERE r.del_flag = 0 AND( owner_id = ? OR sharer_id = ?) GROUP BY r.id", user.Id, user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if rxs[0].Id == 0 {
		rxs = nil
	}

	return rxs
}
