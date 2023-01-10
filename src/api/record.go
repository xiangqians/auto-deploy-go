// record
// @author xiangqian
// @date 23:32 2023/01/10
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"github.com/gin-gonic/gin"
)

func RecordIndex(pContext *gin.Context) {
	HtmlPage(pContext, "record/index.html", func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error) {
		page, err := PageRecord(pContext, pageReq)
		return page, nil, err
	})
}

func PageRecord(pContext *gin.Context, pageReq typ.PageReq) (typ.Page[typ.Record], error) {
	//user := GetUser(pContext)
	return db.Page[typ.Record](pageReq, "SELECT r.id, IFNULL(i.`name`, '') AS 'item_name', r.build_env_id, IFNULL(be.value, '') AS 'build_env_value', r.pull_stime, r.pull_etime, r.pull_status, IFNULL(r.pull_rem, '') AS 'pull_rem', r.commit_id, IFNULL(r.rev_msg, '') AS 'rev_msg', r.build_stime, r.build_etime, r.build_status, IFNULL(r.build_rem, '') AS 'build_rem', r.pack_stime, r.pack_etime, r.pack_status, IFNULL(r.pack_rem, '') AS 'pack_rem', r.ul_stime, r.ul_etime, r.ul_status, IFNULL(r.ul_rem, '') AS 'ul_rem', r.unpack_stime, r.unpack_etime, r.unpack_status, IFNULL(r.unpack_rem, '') AS 'unpack_rem', r.deploy_stime, r.deploy_etime, r.deploy_status, IFNULL(r.deploy_rem, '') AS 'deploy_rem', r.`status`, IFNULL(r.rem, '') AS 'rem', r.del_flag, r.add_time, r.upd_time FROM record r LEFT JOIN build_env be ON be.id = r.build_env_id LEFT JOIN item i ON r.id = r.item_id GROUP BY r.id")
}
