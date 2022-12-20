// index
// @author xiangqian
// @date 21:03 2022/12/18
package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(pContext *gin.Context) {
	pUser := GetUser(pContext)
	if pUser == nil {
		return
	}
	pContext.HTML(http.StatusOK, "index.html", gin.H{
		"username": pUser.Name,
		"nickname": pUser.Nickname,
	})
}
