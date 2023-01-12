// rx
// @author xiangqian
// @date 22:16 2022/12/24
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
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
	HtmlPage(pContext, "rx/index.html", func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error) {
		page, err := PageRx(pContext, pageReq)
		return page, nil, err
	})
}

func RxAddPage(pContext *gin.Context) {
	html := func(rx any, message any) {
		Html(pContext, "rx/add.html", func(pContext *gin.Context) (gin.H, error) {
			return gin.H{"rx": rx}, nil
		})
	}

	session := sessions.Default(pContext)
	rx := session.Get("rx")
	session.Delete("rx")
	session.Save()
	if rx == nil {
		//_rx := typ.Rx{}
		//idStr := pContext.Query("id")
		//id, err := strconv.ParseInt(idStr, 10, 64)
		//if err == nil && id > 0 {
		//	user := SessionUser(pContext)
		//	_, err = db.Qry(&_rx, "SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', r.rem, r.add_time, r.upd_time FROM rx r LEFT JOIN user ou ON ou.del_flag = 0 AND ou.id = r.owner_id LEFT JOIN user su ON su.del_flag = 0 AND su.id = r.sharer_id WHERE r.del_flag = 0 AND r.owner_id = ? AND r.id = ? GROUP BY r.id", user.Id, id)
		//	if err != nil {
		//		html(rx, err)
		//		return
		//	}
		//}
		//rx = _rx
	}

	html(rx, nil)
	return
}

func RxAdd(pContext *gin.Context) {
	RxAddOrUpd(pContext)
}

func RxUpd(pContext *gin.Context) {
	RxAddOrUpd(pContext)
}

func RxAddOrUpd(pContext *gin.Context) {
	redirect := func(rx typ.Rx, message any) {
		Redirect(pContext, "/rx/addpage", message, map[string]any{"rx": rx})
	}

	rx := typ.Rx{}
	err := ShouldBind(pContext, &rx)
	rx.Name = strings.TrimSpace(rx.Name)
	rx.Rem = strings.TrimSpace(rx.Rem)
	if err != nil {
		redirect(rx, err)
		return
	}

	user := SessionUser(pContext)
	if pContext.Request.Method == http.MethodPost {
		_, err = db.Add("INSERT INTO `rx` (`name`, `owner_id`, `rem`, `add_time`) VALUES (?, ?, ?, ?)", rx.Name, user.Id, rx.Rem, time.Now().Unix())
	} else if pContext.Request.Method == http.MethodPut {
		_, err = db.Upd("UPDATE `rx` SET `name` = ?, `rem` = ?, upd_time = ? WHERE `owner_id` = ? AND id = ?", rx.Name, rx.Rem, time.Now().Unix(), user.Id, rx.Id)
	}

	Redirect(pContext, "/rx/index", err, nil)
	return
}

func RxJoin(pContext *gin.Context) {
	codeStr := pContext.PostForm("code")
	code, err := strconv.ParseInt(codeStr, 10, 64)
	message := ""
	if err != nil {
		message = err.Error()
	} else {
		rx, _, _ := db.Qry[typ.Rx]("SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', r.rem, r.add_time, r.upd_time FROM rx r LEFT JOIN user ou ON ou.del_flag = 0 AND ou.id = r.owner_id LEFT JOIN user su ON su.del_flag = 0 AND su.id = r.sharer_id WHERE r.del_flag = 0 AND r.id = ? GROUP BY r.id", code)
		if rx.Id == 0 {
			message = i18n.MustGetMessage("i18n.invalidCode")
		} else {
			if rx.SharerId != 0 {
				message = i18n.MustGetMessage("i18n.codeHasBeenUsed")
			} else {
				user := SessionUser(pContext)
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
	redirect := func(message any) {
		Redirect(pContext, "/rx/index", message, nil)
	}

	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		redirect(err)
		return
	}

	user := SessionUser(pContext)
	db.Del("UPDATE rx SET del_flag = 1, upd_time = ? WHERE (owner_id = ? OR sharer_id = ?) AND id = ?", time.Now().Unix(), user.Id, user.Id, id)

	redirect(nil)
	return
}

func RxShareItemPage(pContext *gin.Context) {
	HtmlPage(pContext, "rx/shareitem.html", func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error) {
		// rx id
		idStr := pContext.Query("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return typ.Page[typ.Item]{}, nil, err
		}

		// rx
		rx, err := Rx(pContext, id, true)
		if err != nil || rx.Id == 0 {
			return typ.Page[typ.Item]{}, nil, err
		}

		page, err := PageRxShareItem(pContext, pageReq, id)
		return page, gin.H{"rx": rx}, err
	})

}

// RxNotShareItems 获取尚未共享的项目集
func RxNotShareItems(pContext *gin.Context) {
	redirect := func(notShareItems []typ.Item) {
		pContext.JSON(http.StatusOK, gin.H{
			"notShareItems": notShareItems,
		})
	}

	// rx id
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		redirect(nil)
		return
	}

	// rx
	rx, err := Rx(pContext, id, false)
	if err != nil {
		log.Println(err)
		redirect(nil)
		return
	}

	// 只返回需要的字段
	page, err := PageItem(pContext, typ.PageReq{Current: 1, Size: 100}, rx.ItemIds)
	notShareItems := page.Data
	if notShareItems != nil {
		_notShareItems := make([]typ.Item, len(notShareItems))
		for i, notShareItem := range notShareItems {
			_notShareItems[i] = typ.Item{Abs: typ.Abs{Id: notShareItem.Id}, Name: notShareItem.Name}
		}
		notShareItems = _notShareItems
	}

	redirect(notShareItems)
	return
}

func RxShareItemAdd(pContext *gin.Context) {
	redirect := func(id int64, message any) {
		Redirect(pContext, fmt.Sprintf("/rx/shareitempage?id=%v", id), message, nil)
	}

	// rx id
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		redirect(id, err)
		return
	}

	// item id
	itemIdStr := pContext.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil || itemId <= 0 {
		redirect(id, err)
		return
	}

	// item
	item, err := Item(pContext, itemId, false)
	if err != nil || item.Id == 0 {
		redirect(id, i18n.MustGetMessage("i18n.itemNotExist"))
		return
	}

	// rx
	rx, err := Rx(pContext, id, false)
	if err != nil || rx.Id == 0 {
		log.Println(err)
		redirect(id, i18n.MustGetMessage("i18n.rxNotExist"))
		return
	}

	// 判断是否已经包含itemId
	itemIds := rx.ItemIds
	contains := strings.Contains(fmt.Sprintf(",%s,", itemIds), fmt.Sprintf(",%v,", itemId))
	if contains {
		redirect(id, nil)
		return
	}

	if itemIds != "" {
		itemIds += ","
	}
	itemIds += fmt.Sprintf("%v", itemId)
	user := SessionUser(pContext)
	db.Upd("UPDATE rx SET item_ids = ?, upd_time = ? WHERE owner_id = ? AND id = ?", itemIds, time.Now().Unix(), user.Id, id)

	redirect(id, nil)
	return
}

func RxShareItemDel(pContext *gin.Context) {
	redirect := func(id int64, message any) {
		Redirect(pContext, fmt.Sprintf("/rx/shareitempage?id=%v", id), message, nil)
	}

	// rx id
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		redirect(id, err)
		return
	}

	// item id
	itemIdStr := pContext.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil || itemId <= 0 {
		redirect(id, err)
		return
	}

	// rx
	rx, err := Rx(pContext, id, false)
	if err != nil || rx.Id == 0 {
		log.Println(err)
		redirect(id, i18n.MustGetMessage("i18n.rxNotExist"))
		return
	}

	// update
	itemIds := fmt.Sprintf(",%s,", rx.ItemIds)
	contains := strings.Contains(itemIds, fmt.Sprintf(",%v,", itemId))
	if contains {
		itemIds = strings.ReplaceAll(itemIds, fmt.Sprintf(",%v,", itemId), ",")
		if itemIds == "," {
			itemIds = ""
		}
		if itemIds != "" {
			itemIds = itemIds[1 : len(itemIds)-1]
		}
		user := SessionUser(pContext)
		db.Del("UPDATE rx SET item_ids = ?, upd_time = ? WHERE owner_id = ? AND id = ?", itemIds, time.Now().Unix(), user.Id, id)
	}
	redirect(id, nil)
}

// PageRxShareItem
// id: rx id
func PageRxShareItem(pContext *gin.Context, pageReq typ.PageReq, id int64) (typ.Page[typ.Item], error) {
	if id <= 0 {
		return typ.Page[typ.Item]{}, nil
	}

	user := SessionUser(pContext)
	return db.Page[typ.Item](pageReq, "SELECT i.id, i.`name`, r.id AS 'rx_id', r.owner_id, i.rem, i.add_time, i.upd_time FROM item i JOIN rx r ON r.del_flag = 0 AND (',' || r.item_ids || ',') LIKE ('%,' || i.id || ',%') WHERE i.del_flag = 0 AND( r.owner_id = ? OR r.sharer_id = ?) AND r.id = ? GROUP BY i.id", user.Id, user.Id, id)
}

func Rx(pContext *gin.Context, id int64, sharer bool) (typ.Rx, error) {
	user := SessionUser(pContext)
	sql := "SELECT r.id, r.`name`, r.owner_id, r.sharer_id, r.item_ids, r.rem FROM rx r "
	sql += "WHERE r.del_flag = 0 "
	if sharer {
		sql += fmt.Sprintf("AND (r.owner_id = %v OR r.sharer_id = %v) ", user.Id, user.Id)
	} else {
		sql += fmt.Sprintf("AND r.owner_id = %v ", user.Id)
	}
	sql += "AND r.id = ? "
	rx, _, err := db.Qry[typ.Rx](sql, id)
	return rx, err
}

func PageRx(pContext *gin.Context, pageReq typ.PageReq) (typ.Page[typ.Rx], error) {
	user := SessionUser(pContext)
	return db.Page[typ.Rx](pageReq, "SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name', IFNULL(ou.nickname, '') AS 'owner_nickname', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', IFNULL(su.nickname, '') AS 'sharer_nickname', r.item_ids, COUNT(DISTINCT i.id) AS 'share_item_count', r.rem, r.add_time, r.upd_time FROM rx r LEFT JOIN user ou ON ou.del_flag = 0 AND ou.id = r.owner_id LEFT JOIN user su ON su.del_flag = 0 AND su.id = r.sharer_id LEFT JOIN item i ON i.del_flag = 0 AND (',' || r.item_ids || ',') LIKE ('%,' || i.id || ',%') WHERE r.del_flag = 0 AND( r.owner_id = ? OR r.sharer_id = ?) GROUP BY r.id", user.Id, user.Id)
}
