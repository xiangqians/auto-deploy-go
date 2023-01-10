// index
// @author xiangqian
// @date 21:03 2022/12/18
package api

import (
	"auto-deploy-go/src/arg"
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/deploy"
	"auto-deploy-go/src/typ"
	"auto-deploy-go/src/util"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// itemIdMap chanel
var _itemIdMap map[int64]int8
var itemIdMapChan chan map[int64]int8

func init() {
	_itemIdMap = make(map[int64]int8, 16)
	itemIdMapChan = make(chan map[int64]int8, 1)
	itemIdMapChan <- _itemIdMap
}

func IndexPage(pContext *gin.Context) {
	// 如果是admin账号登录
	if IsAdminUser(pContext, typ.User{}) {
		IndexAdminPage(pContext)
		return
	}

	html := func(page typ.Page[typ.Record], message any) {
		if v, r := message.(error); r {
			message = v.Error()
		}
		pContext.HTML(http.StatusOK, "index.html", gin.H{
			"user":    GetUser(pContext),
			"page":    page,
			"message": message, // 没有消息就是最好的消息
		})
	}

	pageReq := typ.PageReq{Current: 1, Size: 10}
	err := ShouldBind(pContext, &pageReq)
	if err != nil {
		html(typ.Page[typ.Record]{}, err)
		return
	}

	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()

	page, err := LastRecordPage(pContext, pageReq, 0)
	if err != nil {
		if message != nil {
			message = fmt.Sprintf("%v, %s", message, err)
		} else {
			message = err
		}
	}
	html(page, message)
	return
}

func Deploy(pContext *gin.Context) {
	// redirect func
	redirect := func(status int8, message string) {
		session := sessions.Default(pContext)
		session.Set("status", status)
		session.Set("message", message)
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/")
	}

	// itemId
	itemIdStr := pContext.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		redirect(1, err.Error())
		return
	}
	if itemId <= 0 {
		redirect(1, i18n.MustGetMessage("i18n.itemNotExist"))
		return
	}

	// itemLastRecords
	page, err := LastRecordPage(pContext, typ.PageReq{Current: 1, Size: 10}, itemId)
	if err != nil {
		redirect(1, err.Error())
		return
	}
	if page.Data == nil {
		redirect(1, i18n.MustGetMessage("i18n.itemNotExist"))
		return
	}

	// itemLastRecord
	//itemLastRecord := itemLastRecords[0]
	//if itemLastRecord.Status == typ.StatusInDeploy {
	//	redirect("项目已在部署中")
	//	return
	//}

	// 阻塞获取 chanel 中的 map
	itemIdMap := <-itemIdMapChan
	// 再将 map 添加到 channel
	defer func() {
		itemIdMapChan <- itemIdMap
	}()
	// get <- itemIdMap
	_, r := itemIdMap[itemId]
	if r {
		redirect(1, i18n.MustGetMessage("i18n.itemInDeploy"))
		return
	}

	// item
	item := typ.Item{}
	_, err = db.Qry(&item, "SELECT i.id, i.`name`, i.git_id, i.repo_url, i.branch, i.server_id, i.script, i.rem FROM item i  WHERE i.del_flag = 0 AND i.id = ?", itemId)
	if err != nil {
		redirect(1, err.Error())
		return
	}
	if item.Id == 0 {
		redirect(1, fmt.Sprintf("Item does not exist, %v", itemId))
		return
	}

	// add record
	recordId, err := db.Add("INSERT INTO record(build_env_id, item_id, `status`, `add_time`) VALUES(?, ?, ?)", 0, itemId, typ.StatusInDeploy, time.Now().Unix())
	if err != nil {
		redirect(1, err.Error())
		return
	}

	// put -> itemIdMap
	itemIdMap[itemId] = 1

	// 异步部署
	go asynDeploy(item, recordId)

	redirect(0, i18n.MustGetMessage("i18n.itemDeployStarted"))
	return
}

func asynDeploy(item typ.Item, recordId int64) {
	// updRecord func
	updRecord := func(err error) {
		status := typ.StatusDeploySuccess
		rem := ""
		if err != nil {
			status = typ.StatusDeployExc
			rem = err.Error()
		}
		// update record
		db.Upd("UPDATE record SET `status` = ?, rem = ?, `upd_time` = ? where id = ?", status, rem, time.Now().Unix(), recordId)

		// 阻塞获取 chanel 中的 map
		itemIdMap := <-itemIdMapChan
		// 再将 map 添加到 channel
		defer func() {
			itemIdMapChan <- itemIdMap
		}()
		// 删除key
		delete(itemIdMap, item.Id)
	}

	// delete If Exist
	delIfExist := func(path string) {
		if util.IsExistOfPath(path) {
			util.DelDir(path)
		}
		util.Mkdir(path)
	}

	// base path
	basePath := fmt.Sprintf("%v/item%v", arg.TmpDir, item.Id)

	// localRepoPath
	resPath := fmt.Sprintf("%v/res", basePath)
	delIfExist(resPath)

	// 1. pull
	err := deploy.Pull(item, recordId, resPath)
	if err != nil {
		updRecord(err)
		return
	}

	// script
	script := deploy.ParseScriptTxt(item.Script)

	// 2. build
	err = deploy.Build(script, recordId, resPath)
	if err != nil {
		updRecord(err)
		return
	}

	// 3. pack
	packPath := fmt.Sprintf("%v/pack", basePath)
	delIfExist(packPath)
	packName := fmt.Sprintf("%s/%s", packPath, typ.PackName)
	err = deploy.Pack(script, recordId, resPath, packName)
	if err != nil {
		updRecord(err)
		return
	}

	// 上传到服务路径
	ulPath := fmt.Sprintf("auto-deploy/item%v", item.Id)

	// 4&5. ul and deploy
	err = deploy.UlAndDeploy(item, recordId, packName, ulPath)
	if err != nil {
		updRecord(err)
		return
	}

	// deploy success
	updRecord(nil)
}

func LastRecordPage(pContext *gin.Context, pageReq typ.PageReq, itemId int64) (typ.Page[typ.Record], error) {
	user := GetUser(pContext)
	sql := "SELECT IFNULL(r.id, 0) AS 'id', i.id AS 'item_id', i.`name` AS 'item_name', i.rem AS 'item_rem', " +
		"IFNULL(r.pull_stime, 0) AS 'pull_stime', IFNULL(r.pull_etime, 0) AS 'pull_etime', IFNULL(r.pull_status, 0) AS 'pull_status', IFNULL(r.pull_rem, '') AS 'pull_rem', " +
		"IFNULL(r.commit_id, '') AS 'commit_id', IFNULL(r.rev_msg, '') AS 'rev_msg', " +
		"IFNULL(r.build_stime, 0) AS 'build_stime', IFNULL(r.build_etime, 0) AS 'build_etime', IFNULL(r.build_status, 0) AS 'build_status', IFNULL(r.build_rem, '') AS 'build_rem', " +
		"IFNULL(r.pack_stime, 0) AS 'pack_stime', IFNULL(r.pack_etime, 0) AS 'pack_etime', IFNULL(r.pack_status, 0) AS 'pack_status', IFNULL(r.pack_rem, '') AS 'pack_rem', " +
		"IFNULL(r.ul_stime, 0) AS 'ul_stime', IFNULL(r.ul_etime, 0) AS 'ul_etime', IFNULL(r.ul_status, 0) AS 'ul_status', IFNULL(r.ul_rem, '') AS 'ul_rem', " +
		"IFNULL(r.unpack_stime, 0) AS 'unpack_stime', IFNULL(r.unpack_etime, 0) AS 'unpack_etime', IFNULL(r.unpack_status, 0) AS 'unpack_status', IFNULL(r.unpack_rem, '') AS 'unpack_rem', " +
		"IFNULL(r.deploy_stime, 0) AS 'deploy_stime', IFNULL(r.deploy_etime, 0) AS 'deploy_etime', IFNULL(r.deploy_status, 0) AS 'deploy_status', IFNULL(r.deploy_rem, '') AS 'deploy_rem', " +
		"IFNULL(r.status, 0) AS 'status', " +
		"IFNULL(r.rem, '') AS 'rem', " +
		"IFNULL(r.add_time, 0) AS 'add_time' " +
		"FROM item i " +
		"LEFT JOIN record r ON r.del_flag = 0 AND r.item_id = i.id " +
		"LEFT JOIN record rt ON rt.del_flag = 0 AND rt.item_id = r.item_id AND r.add_time < rt.add_time "

	sql += "WHERE i.del_flag = 0 "
	//sql += "AND i.user_id IN(SELECT DISTINCT(owner_id) FROM rx WHERE del_flag = 0 AND sharer_id = ? UNION ALL SELECT %v) "
	sql += fmt.Sprintf("AND (i.user_id = %v OR EXISTS(SELECT 1 FROM rx rx WHERE rx.del_flag = 0 AND rx.sharer_id = %v AND (',' || rx.item_ids || ',') LIKE ('%%,' || i.id || ',%%') )) ", user.Id, user.Id)
	if itemId > 0 {
		sql += fmt.Sprintf("AND i.id = %v ", strconv.FormatInt(itemId, 10))
	}
	sql += "GROUP BY i.id, r.id HAVING COUNT(rt.id) < 1"
	return db.Page[typ.Record](pageReq, sql)
}
