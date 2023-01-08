// index admin
// @author xiangqian
// @date 13:00 2023/01/08
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"github.com/gin-gonic/gin"
	"net/http"
)

var tables []typ.Table

func init() {
	tables = []typ.Table{
		{Name: "user", Desc: "i18n.userTable"},
		{Name: "rx", Desc: "i18n.rxTable"},
		{Name: "git", Desc: "i18n.gitTable"},
		{Name: "server", Desc: "i18n.serverTable"},
		{Name: "item", Desc: "i18n.itemTable"},
		{Name: "record", Desc: "i18n.recordTable"},
	}
}

func IndexAdminPage(pContext *gin.Context) {

	html := func(page any, err error) {
		message := ""
		if err != nil {
			message = err.Error()
		}
		pContext.HTML(http.StatusOK, "index_admin.html", gin.H{
			"user":    GetUser(pContext),
			"page":    page,
			"table":   tables[2],
			"tables":  tables,
			"message": message,
		})
	}

	pageReq := typ.PageReq{Current: 1, Size: 10}
	err := ShouldBind(pContext, &pageReq)
	//currentStr := pContext.Query("current")
	//sizeStr := pContext.Query("size")
	html(UserPage(pageReq.Current, pageReq.Size), err)
}

func UserPage(current int64, size uint8) typ.Page[typ.User] {
	page := typ.Page[typ.User]{
		Current: current,
		Size:    size,
	}

	page, err := db.Page[typ.User](current, size, "SELECT u.id, u.`name`, u.nickname, u.rem, u.del_flag, u.add_time, u.upd_time FROM `user` u")
	if err != nil {
		return page
	}

	return page
}
