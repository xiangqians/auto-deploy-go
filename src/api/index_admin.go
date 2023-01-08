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

func IndexAdminPage(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "index_admin.html", gin.H{
		"user": GetUser(pContext),
		"page": UserPage(1, 10),
	})
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
