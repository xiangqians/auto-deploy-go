// admin
// @author xiangqian
// @date 20:58 2023/01/07
package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminIndexPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()

	setting, _ := Setting()
	pContext.HTML(http.StatusOK, "admin/index.html", gin.H{
		"user":    GetUser(pContext),
		"setting": setting,
		"message": message,
	})
}
