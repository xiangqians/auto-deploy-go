// rx
// @author xiangqian
// @date 22:16 2022/12/24
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"errors"
	"fmt"
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

func RxShareItemPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()

	var shareItems []typ.Item

	// rx id
	idStr := pContext.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil && id > 0 {
		shareItems = RxShareItems(pContext, id)
	}

	pContext.HTML(http.StatusOK, "rx/shareitem.html", gin.H{
		"shareItems": shareItems,
		"message":    message,
	})
}

func RxShareItemAdd(pContext *gin.Context) {
}

func RxShareItemDel(pContext *gin.Context) {
	redirect := func(id int64, err error) {
		session := sessions.Default(pContext)
		if err != nil {
			session.Set("message", err.Error())
		}
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/rx/shareitempage?id=%v", id))
	}

	// rx id
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		redirect(id, err)
		return
	}

	// item id
	itemIdStr := pContext.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		redirect(id, err)
		return
	}
	if itemId <= 0 {
		redirect(id, nil)
		return
	}

	// rx
	user := GetUser(pContext)
	rx := typ.Rx{}
	err = db.Qry(&rx, "SELECT r.id, r.item_ids FROM rx r WHERE r.del_flag = 0 AND r.owner_id = ? AND r.id = ?", user.Id, id)
	if err != nil {
		redirect(id, err)
		return
	}
	if rx.Id == 0 {
		redirect(id, errors.New(i18n.MustGetMessage("i18n.rxNotExist")))
		return
	}

	// update
	itemIds := strings.ReplaceAll(rx.ItemIds, fmt.Sprintf(",%v,", itemId), ",")
	if itemIds != rx.ItemIds {
		db.Del("UPDATE rx SET item_ids = ?, upd_time = ? WHERE owner_id = ? AND id = ?", itemIds, time.Now().Unix(), user.Id, id)
	}
	redirect(id, nil)
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

// RxShareItems
// id: rx id
func RxShareItems(pContext *gin.Context, id int64) []typ.Item {
	if id <= 0 {
		return nil
	}

	user := GetUser(pContext)
	shareItems := make([]typ.Item, 1)
	err := db.Qry(&shareItems, "SELECT i.id, i.`name`, r.id AS 'rx_id', i.rem, i.add_time, i.upd_time FROM item i JOIN rx r ON r.del_flag = 0 AND r.item_ids LIKE ('%,' || i.id || ',%') WHERE i.del_flag = 0 AND( r.owner_id = ? OR r.sharer_id = ?) AND r.id = ? GROUP BY i.id", user.Id, user.Id, id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if shareItems[0].Id == 0 {
		shareItems = nil
	}

	return shareItems
}

func Rxs(pContext *gin.Context) []typ.Rx {
	user := GetUser(pContext)
	rxs := make([]typ.Rx, 1)
	err := db.Qry(&rxs, "SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name', IFNULL(ou.nickname, '') AS 'owner_nickname', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', IFNULL(su.nickname, '') AS 'sharer_nickname', r.item_ids, COUNT(DISTINCT i.id) AS 'share_item_count', r.rem, r.add_time, r.upd_time FROM rx r LEFT JOIN user ou ON ou.del_flag = 0 AND ou.id = r.owner_id LEFT JOIN user su ON su.del_flag = 0 AND su.id = r.sharer_id LEFT JOIN item i ON i.del_flag = 0 AND r.item_ids LIKE ('%,' || i.id || ',%') WHERE r.del_flag = 0 AND( r.owner_id = ? OR r.sharer_id = ?) GROUP BY r.id", user.Id, user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if rxs[0].Id == 0 {
		rxs = nil
	}

	return rxs
}
