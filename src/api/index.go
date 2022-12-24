// index
// @author xiangqian
// @date 21:03 2022/12/18
package api

import (
	"auto-deploy-go/src/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ItemLastRecord struct {
	Id          int64
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
	CreateTime  int64  // CreateTime
}

func IndexPage(pContext *gin.Context) {
	user := GetUser(pContext)
	itemLastRecords := make([]ItemLastRecord, 1)
	sql := "SELECT IFNULL(r.id, 0) AS 'id', i.`name` AS 'item_name', i.rem AS 'item_rem', " +
		"IFNULL(r.pull_stime, 0) AS 'pull_stime', IFNULL(r.pull_etime, 0) AS 'pull_etime', IFNULL(r.pull_rem, '') AS 'pull_rem', " +
		"IFNULL(r.commit_id, '') AS 'commit_id', IFNULL(r.rev_msg, '') AS 'rev_msg', " +
		"IFNULL(r.build_stime, 0) AS 'build_stime', IFNULL(r.build_etime, 0) AS 'build_etime', IFNULL(r.build_rem, '') AS 'build_rem', " +
		"IFNULL(r.pack_stime, 0) AS 'pack_stime', IFNULL(r.pack_etime, 0) AS 'pack_etime', IFNULL(r.pack_rem, '') AS 'pack_rem', " +
		"IFNULL(r.ul_stime, 0) AS 'ul_stime', IFNULL(r.ul_etime, 0) AS 'ul_etime', IFNULL(r.ul_rem, '') AS 'ul_rem', " +
		"IFNULL(r.deploy_stime, 0) AS 'deploy_stime', IFNULL(r.deploy_etime, 0) AS 'deploy_etime', IFNULL(r.deploy_rem, '') AS 'deploy_rem', " +
		"IFNULL(r.status, 0) AS 'status', " +
		"IFNULL(r.rem, '') AS 'rem', " +
		"IFNULL(r.create_time, 0) AS 'create_time' " +
		"FROM item i " +
		"LEFT JOIN record r ON r.del_flag = 0 AND r.item_id = i.id " +
		"LEFT JOIN record rt ON rt.del_flag = 0 AND rt.item_id = r.item_id AND r.create_time < rt.create_time " +
		"WHERE i.del_flag = 0 AND i.user_id = ? " +
		"GROUP BY i.id, r.id " +
		"HAVING COUNT(rt.id) < 1"
	err := db.Qry(&itemLastRecords, sql, user.Id)
	if err != nil {
		log.Println(err)
	}
	if itemLastRecords[0].Id == 0 {
		itemLastRecords = nil
	}

	pContext.HTML(http.StatusOK, "index.html", gin.H{
		"user":            user,
		"itemLastRecords": itemLastRecords,
	})
}
