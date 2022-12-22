// item
// @author xiangqian
// @date 15:02 2022/12/20
package api

import (
	"auto-deploy-go/src/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 项目
type Item struct {
	Abs
	UserId   int64  // 项目所属用户id
	Name     string // 名称
	GitId    int64  // 项目所属Git id
	RepoUrl  string // Git仓库地址
	Branch   string // 分支名
	ServerId int64  // 项目所属Server id
	Cmd      string // 构建命令
	Script   string // 脚本，目前支持 #!/dockerfile, #!/static 解析

	//
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
