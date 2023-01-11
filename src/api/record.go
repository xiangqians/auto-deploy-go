// record
// @author xiangqian
// @date 23:32 2023/01/10
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

const RecordColumns = "IFNULL(r.id, 0) AS 'id', " +
	"IFNULL(r.pull_stime, 0) AS 'pull_stime', IFNULL(r.pull_etime, 0) AS 'pull_etime', IFNULL(r.pull_status, 0) AS 'pull_status', IFNULL(r.pull_rem, '') AS 'pull_rem', " +
	"IFNULL(r.commit_id, '') AS 'commit_id', IFNULL(r.rev_msg, '') AS 'rev_msg', " +
	"IFNULL(r.build_stime, 0) AS 'build_stime', IFNULL(r.build_etime, 0) AS 'build_etime', IFNULL(r.build_status, 0) AS 'build_status', IFNULL(r.build_rem, '') AS 'build_rem', " +
	"IFNULL(r.pack_stime, 0) AS 'pack_stime', IFNULL(r.pack_etime, 0) AS 'pack_etime', IFNULL(r.pack_status, 0) AS 'pack_status', IFNULL(r.pack_rem, '') AS 'pack_rem', " +
	"IFNULL(r.ul_stime, 0) AS 'ul_stime', IFNULL(r.ul_etime, 0) AS 'ul_etime', IFNULL(r.ul_status, 0) AS 'ul_status', IFNULL(r.ul_rem, '') AS 'ul_rem', " +
	"IFNULL(r.unpack_stime, 0) AS 'unpack_stime', IFNULL(r.unpack_etime, 0) AS 'unpack_etime', IFNULL(r.unpack_status, 0) AS 'unpack_status', IFNULL(r.unpack_rem, '') AS 'unpack_rem', " +
	"IFNULL(r.deploy_stime, 0) AS 'deploy_stime', IFNULL(r.deploy_etime, 0) AS 'deploy_etime', IFNULL(r.deploy_status, 0) AS 'deploy_status', IFNULL(r.deploy_rem, '') AS 'deploy_rem', " +
	"IFNULL(r.status, 0) AS 'status', " +
	"IFNULL(r.rem, '') AS 'rem', " +
	"IFNULL(r.add_time, 0) AS 'add_time' "

func RecordIndex(pContext *gin.Context) {
	HtmlPage(pContext, "record/index.html", func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error) {
		itemIdStr := pContext.Param("itemId")
		itemId, err := strconv.ParseInt(itemIdStr, 10, 8)
		if err != nil || itemId <= 0 {
			return typ.Page[typ.Record]{}, nil, err
		}

		item, page, err := PageRecord(pContext, pageReq, itemId)
		return page, gin.H{"item": item}, err
	})
}

func RecordDel(pContext *gin.Context) {
	redirect := func(itemId int64, message any) {
		Redirect(pContext, fmt.Sprintf("/record/index/%v", itemId), message, nil)
	}

	// item id
	itemIdStr := pContext.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 8)
	if err != nil || itemId <= 0 {
		redirect(itemId, err)
		return
	}

	// record id
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 8)
	if err != nil || id <= 0 {
		redirect(itemId, err)
		return
	}

	// 判断 itemId 是否是当前用户拥有者
	item, err := Item(pContext, itemId, false)
	if err != nil || item.Id == 0 {
		redirect(itemId, err)
		return
	}

	// 逻辑删除record信息
	_, err = db.Del("UPDATE record SET del_flag = 1, upd_time = ? WHERE id = ?", time.Now().Unix(), id)
	redirect(itemId, err)
	return
}

// PageRecord record分页查询
func PageRecord(pContext *gin.Context, pageReq typ.PageReq, itemId int64) (typ.Item, typ.Page[typ.Record], error) {
	// 判断当前用户是否拥有访问 itemId 的权限
	item, err := Item(pContext, itemId, true)
	if err != nil || item.Id == 0 {
		return item, typ.Page[typ.Record]{}, err
	}

	// 分页查询record
	sql := fmt.Sprintf("SELECT %s FROM record r WHERE r.del_flag = 0 AND r.item_id = ? ORDER BY r.add_time DESC", RecordColumns)
	page, err := db.Page[typ.Record](pageReq, sql, itemId)
	return item, page, err
}
