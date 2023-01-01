// item
// @author xiangqian
// @date 15:02 2022/12/20
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ItemIndex(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()
	pContext.HTML(http.StatusOK, "item/index.html", gin.H{
		"user":    GetUser(pContext),
		"items":   Items(pContext),
		"message": message,
	})
}

func ItemAddPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	item := session.Get("item")
	message := session.Get("message")
	session.Delete("item")
	session.Delete("message")
	session.Save()

	if item == nil {
		_item := typ.Item{}
		idStr := pContext.Query("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil && id > 0 {
			user := GetUser(pContext)
			err = db.Qry(&_item, "SELECT i.id, i.`name`, i.git_id, i.repo_url, i.branch, i.server_id, i.script, i.rem FROM item i  WHERE i.del_flag = 0 AND i.user_id = ? AND i.id = ?", user.Id, id)
			if err != nil {
				log.Println(err)
			}
		}
		item = _item
	}

	pContext.HTML(http.StatusOK, "item/add.html", gin.H{
		"gits":    Gits(pContext),
		"servers": Servers(pContext),
		"item":    item,
		"message": message,
	})
}

func ItemAdd(pContext *gin.Context) {
	item, err := itemPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Add("INSERT INTO `item` (`user_id`, `name`, `git_id`, `repo_url`, `branch`, `server_id`, `script`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, item.Name, item.GitId, item.RepoUrl, item.Branch, item.ServerId, item.Script, item.Rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func ItemUpd(pContext *gin.Context) {
	item, err := itemPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Upd("UPDATE `item` SET `name` = ?, `git_id` = ?, `repo_url` = ?, `branch` = ?, `server_id` = ?, `script` = ?, `rem` = ?, upd_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
		item.Name, item.GitId, item.RepoUrl, item.Branch, item.ServerId, item.Script, item.Rem, time.Now().Unix(), user.Id, item.Id)
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func ItemDel(pContext *gin.Context) {
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil {
		user := GetUser(pContext)
		db.Del("UPDATE item SET del_flag = 1, upd_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	}
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func itemPreAddOrUpd(pContext *gin.Context) (typ.Item, error) {
	item := typ.Item{}
	err := ShouldBind(pContext, &item)
	if err != nil {
		session := sessions.Default(pContext)
		session.Set("item", item)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/item/addpage")
	}

	return item, err
}

func Items(pContext *gin.Context) []typ.Item {
	user := GetUser(pContext)
	items := make([]typ.Item, 1)
	err := db.Qry(&items, "SELECT i.id, i.`name`, i.git_id, IFNULL(g.`name`, '') AS 'git_name', i.repo_url, i.branch, i.server_id, IFNULL(s.`name`, '') AS 'server_name', i.rem, i.add_time, i.upd_time FROM item i LEFT JOIN git g ON g.del_flag = 0 AND g.id = i.git_id LEFT JOIN server s ON s.del_flag = 0 AND s.id = i.server_id WHERE i.del_flag = 0 AND i.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if items[0].Id == 0 {
		items = nil
	}

	return items
}
