// rx
// @author xiangqian
// @date 22:16 2022/12/24
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RxIndex(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()
	pContext.HTML(http.StatusOK, "rx/index.html", gin.H{
		"user":    GetUser(pContext),
		"rxs":     Rxs(pContext),
		"message": message,
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
		_rx := typ.Rx{}
		idStr := pContext.Query("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil && id > 0 {
			user := GetUser(pContext)
			err = db.Qry(&_rx, "SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', r.rem, r.add_time, r.upd_time FROM rx r LEFT JOIN user ou ON ou.del_flag = 0 AND ou.id = r.owner_id LEFT JOIN user su ON su.del_flag = 0 AND su.id = r.sharer_id WHERE r.del_flag = 0 AND r.owner_id = ? AND r.id = ? GROUP BY r.id", user.Id, id)
			if err != nil {
				log.Println(err)
			}
		}
		rx = _rx
	}
	pContext.HTML(http.StatusOK, "rx/add.html", gin.H{
		"rx":      rx,
		"message": message,
	})
}

func RxAdd(pContext *gin.Context) {
	rx, err := rxPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Add("INSERT INTO `rx` (`name`, `owner_id`, `rem`, `add_time`) VALUES (?, ?, ?, ?)",
		rx.Name, user.Id, rx.Rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/rx/index")
}

func RxUpd(pContext *gin.Context) {
	rx, err := rxPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Upd("UPDATE `rx` SET `name` = ?, `rem` = ?, upd_time = ? WHERE `owner_id` = ? AND id = ?",
		rx.Name, rx.Rem, time.Now().Unix(), user.Id, rx.Id)
	pContext.Redirect(http.StatusMovedPermanently, "/rx/index")

}

func RxJoin(pContext *gin.Context) {
	codeStr := pContext.PostForm("code")
	code, err := strconv.ParseInt(codeStr, 10, 64)
	message := ""
	if err != nil {
		message = err.Error()
	} else {
		rx := typ.Rx{}
		err = db.Qry(&rx, "SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', r.rem, r.add_time, r.upd_time FROM rx r LEFT JOIN user ou ON ou.del_flag = 0 AND ou.id = r.owner_id LEFT JOIN user su ON su.del_flag = 0 AND su.id = r.sharer_id WHERE r.del_flag = 0 AND r.id = ? GROUP BY r.id", code)
		if rx.Id == 0 {
			message = i18n.MustGetMessage("i18n.invalidCode")
		} else {
			if rx.SharerId != 0 {
				message = i18n.MustGetMessage("i18n.codeHasBeenUsed")
			} else {
				user := GetUser(pContext)
				if rx.OwnerId == user.Id {
					message = i18n.MustGetMessage("i18n.yourCodeCannotBeSharedByYourself")
				} else {
					db.Upd("UPDATE rx SET sharer_id = ?, upd_time = ? WHERE del_flag = 0 AND sharer_id = 0 AND owner_id != ? AND id = ?",
						user.Id, time.Now().Unix(), user.Id, code)
				}
			}
		}
	}

	session := sessions.Default(pContext)
	session.Set("message", message)
	session.Save()
	pContext.Redirect(http.StatusMovedPermanently, "/rx/index")
}

func RxDel(pContext *gin.Context) {
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil {
		user := GetUser(pContext)
		db.Del("UPDATE rx SET del_flag = 1, upd_time = ? WHERE (owner_id = ? OR sharer_id = ?) AND id = ?", time.Now().Unix(), user.Id, user.Id, id)
	}
	pContext.Redirect(http.StatusMovedPermanently, "/rx/index")
}

func rxPreAddOrUpd(pContext *gin.Context) (typ.Rx, error) {
	rx := typ.Rx{}
	err := ShouldBind(pContext, &rx)

	rx.Name = strings.TrimSpace(rx.Name)
	rx.Rem = strings.TrimSpace(rx.Rem)

	if err != nil {
		session := sessions.Default(pContext)
		session.Set("rx", rx)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/rx/addpage")
	}

	return rx, err
}

func Rxs(pContext *gin.Context) []typ.Rx {
	user := GetUser(pContext)
	rxs := make([]typ.Rx, 1)
	err := db.Qry(&rxs, "SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', r.rem, r.add_time, r.upd_time FROM rx r LEFT JOIN user ou ON ou.del_flag = 0 AND ou.id = r.owner_id LEFT JOIN user su ON su.del_flag = 0 AND su.id = r.sharer_id WHERE r.del_flag = 0 AND( r.owner_id = ? OR r.sharer_id = ?) GROUP BY r.id", user.Id, user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if rxs[0].Id == 0 {
		rxs = nil
	}

	return rxs
}
