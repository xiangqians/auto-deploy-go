// item
// @author xiangqian
// @date 15:02 2022/12/20
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ItemIndex(pContext *gin.Context) {
	HtmlPage(pContext, "item/index.html", func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error) {
		page, err := PageItem(pContext, pageReq, "")
		return page, nil, err
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
			_, err = db.Qry(&_item, "SELECT i.id, i.`name`, i.git_id, i.repo_url, i.branch, i.server_id, i.script, i.rem FROM item i  WHERE i.del_flag = 0 AND i.user_id = ? AND i.id = ?", user.Id, id)
			if err != nil {
				log.Println(err)
			}
		}
		item = _item
	}

	gitPage, _ := PageGit(pContext, typ.PageReq{Current: 1, Size: 100})
	serverPage, _ := PageServer(pContext, typ.PageReq{Current: 1, Size: 100})
	pContext.HTML(http.StatusOK, "item/add.html", gin.H{
		"gits":    gitPage.Data,
		"servers": serverPage.Data,
		"item":    item,
		"message": message,
	})
}

func ItemAdd(pContext *gin.Context) {
	ItemAddOrUpd(pContext)
}

func ItemUpd(pContext *gin.Context) {
	ItemAddOrUpd(pContext)
}

func ItemAddOrUpd(pContext *gin.Context) {
	redirect := func(item typ.Item, message any) {
		Redirect(pContext, "/item/addpage", message, map[string]any{"item": item})
	}

	item := typ.Item{}
	err := ShouldBind(pContext, &item)
	if err != nil {
		redirect(item, err)
		return
	}

	user := GetUser(pContext)
	if pContext.Request.Method == http.MethodPost {
		_, err = db.Add("INSERT INTO `item` (`user_id`, `name`, `git_id`, `repo_url`, `branch`, `server_id`, `script`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			user.Id, item.Name, item.GitId, item.RepoUrl, item.Branch, item.ServerId, item.Script, item.Rem, time.Now().Unix())
	} else if pContext.Request.Method == http.MethodPut {
		_, err = db.Upd("UPDATE `item` SET `name` = ?, `git_id` = ?, `repo_url` = ?, `branch` = ?, `server_id` = ?, `script` = ?, `rem` = ?, upd_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
			item.Name, item.GitId, item.RepoUrl, item.Branch, item.ServerId, item.Script, item.Rem, time.Now().Unix(), user.Id, item.Id)
	}

	Redirect(pContext, "/item/index", err, nil)
	return
}

func ItemDel(pContext *gin.Context) {
	redirect := func(message any) {
		Redirect(pContext, "/item/index", message, nil)
	}

	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		redirect(err)
		return
	}

	user := GetUser(pContext)
	db.Del("UPDATE item SET del_flag = 1, upd_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	redirect(nil)
	return
}

// Item 获取item信息
// id: item id
// sharer: 是否包含共享者
func Item(pContext *gin.Context, id int64, sharer bool) (typ.Item, error) {
	user := GetUser(pContext)
	sql := "SELECT i.id, i.user_id, i.`name`, i.git_id, i.repo_url, i.branch, i.server_id, i.script, i.rem "
	sql += "FROM item i "
	sql += "WHERE i.del_flag = 0 "
	sql += fmt.Sprintf("AND i.id = %v ", id)
	if sharer {
		sql += fmt.Sprintf("AND (i.user_id = %v OR EXISTS(SELECT 1 FROM rx rx WHERE rx.del_flag = 0 AND rx.sharer_id = %v AND (',' || rx.item_ids || ',') LIKE ('%%,' || i.id || ',%%') ))", user.Id, user.Id)
	} else {
		sql += fmt.Sprintf(" AND i.user_id = %v ", user.Id)
	}
	item := typ.Item{}
	_, err := db.Qry(&item, sql)
	return item, err
}

func PageItem(pContext *gin.Context, pageReq typ.PageReq, notLikeIds string) (typ.Page[typ.Item], error) {
	user := GetUser(pContext)
	sql := "SELECT i.id, i.`name`, i.git_id, IFNULL(g.`name`, '') AS 'git_name', i.repo_url, i.branch, i.server_id, IFNULL(s.`name`, '') AS 'server_name', i.rem, i.add_time, i.upd_time FROM item i LEFT JOIN git g ON g.del_flag = 0 AND g.id = i.git_id LEFT JOIN server s ON s.del_flag = 0 AND s.id = i.server_id WHERE i.del_flag = 0 AND i.user_id = ? "
	if notLikeIds != "" {
		sql += fmt.Sprintf("AND ',%s,' NOT LIKE ('%%,' || i.id || ',%%') ", notLikeIds)
	}
	sql += "GROUP BY i.id"
	return db.Page[typ.Item](pageReq, sql, user.Id)
}
