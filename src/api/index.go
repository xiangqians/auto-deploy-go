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

func IndexPage(pContext *gin.Context) {
	user := GetUser(pContext)

	//SELECT
	//i.id, i.`name`, g.`name` AS 'gname', s.`name` AS 'sname', i.rem, i.create_time, i.update_time
	//FROM item i
	//LEFT JOIN git g ON g.id = i.git_id
	//LEFT JOIN server s ON s.id = i.server_id
	//WHERE i.del_flag = 0
	items := make([]Item, 1)
	err := db.Qry(&items, "SELECT i.id, i.`name`, i.rem, i.create_time, i.update_time FROM item i WHERE i.del_flag = 0 AND i.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
	}
	if items[0].Id == 0 {
		items = nil
	}

	pContext.HTML(http.StatusOK, "index.html", gin.H{
		"user":  user,
		"items": items,
	})
}
