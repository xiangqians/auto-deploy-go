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
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	netHttp "net/http"
	"strconv"
	"time"
)

var _itemIdMap map[int64]int8
var itemIdMapChan chan map[int64]int8

func init() {
	_itemIdMap = make(map[int64]int8, 16)
	itemIdMapChan = make(chan map[int64]int8, 1)
	itemIdMapChan <- _itemIdMap
}

func IndexPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()
	pContext.HTML(netHttp.StatusOK, "index.html", gin.H{
		"user":            GetUser(pContext),
		"itemLastRecords": getItemLastRecords(pContext, 0),
		"message":         message,
	})
}

func Deploy(pContext *gin.Context) {
	// redirect func
	redirect := func(message string) {
		session := sessions.Default(pContext)
		session.Set("message", message)
		session.Save()
		pContext.Redirect(netHttp.StatusMovedPermanently, "/")
	}

	// itemId
	itemIdStr := pContext.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		redirect(err.Error())
		return
	}

	// itemLastRecords
	itemLastRecords := getItemLastRecords(pContext, itemId)
	if itemLastRecords == nil {
		redirect("项目不存在")
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
	// get
	_, r := itemIdMap[itemId]
	if r {
		redirect("项目已在部署中")
		return
	}
	// put
	itemIdMap[itemId] = 1

	// item
	item := typ.Item{}
	err = db.Qry(&item, "SELECT i.id, i.`name`, i.git_id, i.repo_url, i.branch, i.server_id, i.script, i.rem FROM item i  WHERE i.del_flag = 0 AND i.id = ?", itemId)
	if err != nil {
		redirect(err.Error())
		return
	}

	// add record
	recordId, err := db.Add("INSERT INTO record(item_id, `status`, `add_time`) VALUES(?, ?, ?)", itemId, typ.StatusInDeploy, time.Now().Unix())
	if err != nil {
		redirect(err.Error())
		return
	}

	// 异步部署
	go asynDeploy(item, recordId)

	redirect("")
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

func getItemLastRecords(pContext *gin.Context, itemId int64) []typ.ItemLastRecord {
	itemLastRecords := make([]typ.ItemLastRecord, 1)
	user := GetUser(pContext)
	sql := "SELECT IFNULL(r.id, 0) AS 'id', i.id AS 'item_id', i.`name` AS 'item_name', i.rem AS 'item_rem', " +
		"IFNULL(r.pull_stime, 0) AS 'pull_stime', IFNULL(r.pull_etime, 0) AS 'pull_etime', IFNULL(r.pull_rem, '') AS 'pull_rem', " +
		"IFNULL(r.commit_id, '') AS 'commit_id', IFNULL(r.rev_msg, '') AS 'rev_msg', " +
		"IFNULL(r.build_stime, 0) AS 'build_stime', IFNULL(r.build_etime, 0) AS 'build_etime', IFNULL(r.build_rem, '') AS 'build_rem', " +
		"IFNULL(r.pack_stime, 0) AS 'pack_stime', IFNULL(r.pack_etime, 0) AS 'pack_etime', IFNULL(r.pack_rem, '') AS 'pack_rem', " +
		"IFNULL(r.ul_stime, 0) AS 'ul_stime', IFNULL(r.ul_etime, 0) AS 'ul_etime', IFNULL(r.ul_rem, '') AS 'ul_rem', " +
		"IFNULL(r.deploy_stime, 0) AS 'deploy_stime', IFNULL(r.deploy_etime, 0) AS 'deploy_etime', IFNULL(r.deploy_rem, '') AS 'deploy_rem', " +
		"IFNULL(r.status, 0) AS 'status', " +
		"IFNULL(r.rem, '') AS 'rem', " +
		"IFNULL(r.add_time, 0) AS 'add_time' " +
		"FROM item i " +
		"LEFT JOIN record r ON r.del_flag = 0 AND r.item_id = i.id " +
		"LEFT JOIN record rt ON rt.del_flag = 0 AND rt.item_id = r.item_id AND r.add_time < rt.add_time "

	sql += "WHERE i.del_flag = 0 AND i.user_id IN(SELECT DISTINCT(owner_id) FROM rx WHERE del_flag = 0 AND sharer_id = ? UNION ALL SELECT %v) "
	if itemId > 0 {
		sql += fmt.Sprintf("AND i.id = %v ", strconv.FormatInt(itemId, 10))
	}
	sql += "GROUP BY i.id, r.id HAVING COUNT(rt.id) < 1"
	sql = fmt.Sprintf(sql, user.Id)
	err := db.Qry(&itemLastRecords, sql, user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if itemLastRecords[0].ItemId == 0 {
		itemLastRecords = nil
	}

	return itemLastRecords
}
