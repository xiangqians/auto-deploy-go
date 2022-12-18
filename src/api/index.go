// index
// @author xiangqian
// @date 21:03 2022/12/18
package api

import (
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "index.html", gin.H{
		"username": i18n.MustGetMessage("username"),
		"password": i18n.MustGetMessage("password"),
		"submit":   i18n.MustGetMessage("submit"),
	})
}
