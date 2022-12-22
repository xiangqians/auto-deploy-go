// index
// @author xiangqian
// @date 21:03 2022/12/18
package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexPage(pContext *gin.Context) {
	user := GetUser(pContext)
	pContext.HTML(http.StatusOK, "index.html", gin.H{
		"username": user.Name,
		"nickname": user.Nickname,
	})
}
