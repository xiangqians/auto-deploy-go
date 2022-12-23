// item
// @author xiangqian
// @date 15:02 2022/12/20
package api

import (
	"auto-deploy-go/src/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

// Item 项目
type Item struct {
	Abs
	UserId     int64  // 项目所属用户id
	Name       string // 名称
	GitId      int64  // 项目所属Git id
	GitName    string // 项目所属Git Name
	RepoUrl    string // Git仓库地址
	Branch     string // 分支名
	ServerId   int64  // 项目所属Server id
	ServerName string // 项目所属Server Name
	Ini        string // 脚本
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
	user := GetUser(pContext)
	items := make([]Item, 1)
	err := db.Qry(&items, "SELECT i.id, i.`name`, IFNULL(g.`name`, '') AS 'git_name', i.repo_url, i.branch, s.`name` AS 'server_name', i.rem, i.create_time, i.update_time FROM item i LEFT JOIN git g ON g.del_flag = 0 AND g.id = i.git_id LEFT JOIN server s ON s.del_flag = 0 AND s.id = i.server_id WHERE i.del_flag = 0 AND i.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
	}
	if items[0].Id == 0 {
		items = nil
	}

	pContext.HTML(http.StatusOK, "item/index.html", gin.H{
		"items": items,
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
		"item":    item,
		"message": message,
	})
}

func ItemAdd(pContext *gin.Context) {
	name := strings.TrimSpace(pContext.PostForm("name"))
	gitId := strings.TrimSpace(pContext.PostForm("gitId"))
	repoUrl := strings.TrimSpace(pContext.PostForm("repoUrl"))
	branch := strings.TrimSpace(pContext.PostForm("branch"))
	serverId := strings.TrimSpace(pContext.PostForm("serverId"))
	cmd := strings.TrimSpace(pContext.PostForm("cmd"))
	script := strings.TrimSpace(pContext.PostForm("script"))
	rem := strings.TrimSpace(pContext.PostForm("rem"))
	user := GetUser(pContext)
	db.Add("INSERT INTO `item` (`user_id`, `name`, `git_id`, `repo_url`, `branch`, `server_id`, `cmd`, `script`, `rem`, `create_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, name, gitId, repoUrl, branch, serverId, cmd, script, rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/")
}

func ItemUpd(pContext *gin.Context) {
	pContext.Redirect(http.StatusMovedPermanently, "/")
}

func ItemDel(pContext *gin.Context) {
	pContext.Redirect(http.StatusMovedPermanently, "/")
}
