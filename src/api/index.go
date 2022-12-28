// index
// @author xiangqian
// @date 21:03 2022/12/18
package api

import (
	"auto-deploy-go/src/db"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type ItemLastRecord struct {
	Id          int64
	ItemId      int64
	ItemName    string // item
	ItemRem     string
	PullStime   int64 // pull
	PullEtime   int64
	PullRem     string
	CommitId    string // commitId
	RevMsg      string // revMsg
	BuildStime  int64  // build
	BuildEtime  int64
	BuildRem    string
	PackStime   int64 // pack
	PackEtime   int64
	PackRem     string
	UlStime     int64 // ul
	UlEtime     int64
	UlRem       string
	DeployStime int64 // deploy
	DeployEtime int64
	DeployRem   string
	Status      byte   // status
	Rem         string // Rem
	AddTime     int64  // AddTime
}

func IndexPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()
	pContext.HTML(http.StatusOK, "index.html", gin.H{
		"user":            GetUser(pContext),
		"itemLastRecords": ItemLastRecords(pContext, 0),
		"message":         message,
	})
}

// 1. 集成git拉取代码
func Deploy(pContext *gin.Context) {
	message := ""
	itemIdStr := pContext.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		message = err.Error()
	} else {
		itemLastRecords := ItemLastRecords(pContext, itemId)
		if itemLastRecords == nil {
			message = "项目不存在"
		} else {

		}
	}

	session := sessions.Default(pContext)
	session.Set("message", message)
	session.Save()
	pContext.Redirect(http.StatusMovedPermanently, "/")
}

func ItemLastRecords(pContext *gin.Context, itemId int64) []ItemLastRecord {
	itemLastRecords := make([]ItemLastRecord, 1)
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
