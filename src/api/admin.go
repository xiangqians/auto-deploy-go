// admin
// @author xiangqian
// @date 13:39 2023/01/07
package api

import (
	"github.com/gin-gonic/gin"
	netHttp "net/http"
)

func AdminIndexPage(pContext *gin.Context) {
	pContext.HTML(netHttp.StatusOK, "admin/index.html", gin.H{
		"user": GetUser(pContext),
	})
}
