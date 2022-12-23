// item
// @author xiangqian
// @date 15:02 2022/12/20
package api

import (
	"auto-deploy-go/src/db"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// Item 项目
type Item struct {
	Abs
	UserId     int64  `form:"userId"`                                    // 项目所属用户id
	Name       string `form:"name" binding:"required,trim,min=3,max=10"` // 名称
	GitId      int64  `form:"gitId"`                                     // 项目所属Git id
	GitName    string `form:"gitName"`                                   // 项目所属Git Name
	RepoUrl    string `form:"repoUrl"`                                   // Git仓库地址
	Branch     string `form:"branch"`                                    // 分支名
	ServerId   int64  `form:"serverId"`                                  // 项目所属Server id
	ServerName string `form:"serverName"`                                // 项目所属Server Name
	Ini        string `form:"ini"`                                       // 脚本
}

func init() {
	// 注册 Item 模型
	gob.Register(Item{})
}

// Stage 自动化部署阶段
type Stage int8

const (
	StagePull   Stage = iota + 1 // 拉取资源
	StageBuild                   // 构建
	StagePack                    // 打包
	StageUl                      // upload上传
	StageDeploy                  // 部署
)

func ItemIndex(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "item/index.html", gin.H{
		"items": Items(pContext),
	})
}

func ItemAddPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	item := session.Get("item")
	message := session.Get("message")
	session.Delete("item")
	session.Delete("message")
	session.Save()

	if item == nil {
		item = Item{}
	}

	pContext.HTML(http.StatusOK, "item/add.html", gin.H{
		"gits":    Gits(pContext),
		"servers": Servers(pContext),
		"item":    item,
		"message": message,
	})
}

func ItemAdd(pContext *gin.Context) {
	item := Item{}
	err := ShouldBind(pContext, &item)

	message := ""
	//if err == nil {
	//	err = com.VerifyText(item.Name, 60)
	//	if err != nil {
	//		message = fmt.Sprintf(err.Error(), i18n.MustGetMessage("i18n.name"))
	//	}
	//}

	if err != nil {
		session := sessions.Default(pContext)
		session.Set("item", item)
		if message != "" {
			session.Set("message", message)
		} else {
			session.Set("message", err.Error())
		}
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/item/addpage")
		return
	}

	user := GetUser(pContext)
	db.Add("INSERT INTO `item` (`user_id`, `name`, `git_id`, `repo_url`, `branch`, `server_id`, `ini`, `rem`, `create_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, item.Name, item.GitId, item.RepoUrl, item.Branch, item.ServerId, item.Ini, item.Rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func ItemUpd(pContext *gin.Context) {
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func ItemDel(pContext *gin.Context) {
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func Items(pContext *gin.Context) []Item {
	user := GetUser(pContext)
	items := make([]Item, 1)
	err := db.Qry(&items, "SELECT i.id, i.`name`, i.git_id, IFNULL(g.`name`, '') AS 'git_name', i.repo_url, i.branch, i.server_id, IFNULL(s.`name`, '') AS 'server_name', i.rem, i.create_time, i.update_time FROM item i LEFT JOIN git g ON g.del_flag = 0 AND g.id = i.git_id LEFT JOIN server s ON s.del_flag = 0 AND s.id = i.server_id WHERE i.del_flag = 0 AND i.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if items[0].Id == 0 {
		items = nil
	}

	return items
}
