// item
// @author xiangqian
// @date 15:02 2022/12/20
package api

import (
	"auto-deploy-go/src/db"
	"bufio"
	"bytes"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Stage 自动化部署阶段
type Stage int8

const (
	TagBuild  string = "[build]"
	TagTarget        = "[target]"
	TagDeploy        = "[deploy]"
)

// Item 项目
type Item struct {
	Abs
	UserId     int64  `form:"userId"`                                              // 项目所属用户id
	Name       string `form:"name" binding:"required,excludes= ,min=1,max=60"`     // 名称
	GitId      int64  `form:"gitId" binding:"gte=0"`                               // 项目所属Git id
	GitName    string `form:"gitName"`                                             // 项目所属Git Name
	RepoUrl    string `form:"repoUrl" binding:"required,excludes= ,min=1,max=500"` // Git仓库地址
	Branch     string `form:"branch" binding:"required,excludes= ,min=1,max=60"`   // 分支名
	ServerId   int64  `form:"serverId" binding:"required,gt=0"`                    // 项目所属Server id
	ServerName string `form:"serverName"`                                          // 项目所属Server Name
	Ini        string `form:"ini" binding:"required,min=1,max=1000"`               // 脚本
}

type Ini struct {
	Build  string // [build]
	Target string // [target]
	Deploy string // [deploy]
}

func init() {
	// 注册 Item 模型
	gob.Register(Item{})
}

func ItemIndex(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()
	pContext.HTML(http.StatusOK, "item/index.html", gin.H{
		"items":   Items(pContext),
		"message": message,
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
		_item := Item{}
		idStr := pContext.Query("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil && id > 0 {
			user := GetUser(pContext)
			err = db.Qry(&_item, "SELECT i.id, i.`name`, i.git_id, i.repo_url, i.branch, i.server_id, i.ini, i.rem FROM item i  WHERE i.del_flag = 0 AND i.user_id = ? AND i.id = ?", user.Id, id)
			if err != nil {
				log.Println(err)
			}
		}
		item = _item
	}

	pContext.HTML(http.StatusOK, "item/add.html", gin.H{
		"gits":    Gits(pContext),
		"servers": Servers(pContext),
		"item":    item,
		"message": message,
	})
}

func ItemAdd(pContext *gin.Context) {
	item, err := itemPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Add("INSERT INTO `item` (`user_id`, `name`, `git_id`, `repo_url`, `branch`, `server_id`, `ini`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, item.Name, item.GitId, item.RepoUrl, item.Branch, item.ServerId, item.Ini, item.Rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func ItemUpd(pContext *gin.Context) {
	item, err := itemPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Upd("UPDATE `item` SET `name` = ?, `git_id` = ?, `repo_url` = ?, `branch` = ?, `server_id` = ?, `ini` = ?, `rem` = ?, upd_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
		item.Name, item.GitId, item.RepoUrl, item.Branch, item.ServerId, item.Ini, item.Rem, time.Now().Unix(), user.Id, item.Id)
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func ItemDel(pContext *gin.Context) {
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil {
		user := GetUser(pContext)
		db.Del("UPDATE item SET del_flag = 1, upd_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	}
	pContext.Redirect(http.StatusMovedPermanently, "/item/index")
}

func itemPreAddOrUpd(pContext *gin.Context) (Item, error) {
	item := Item{}
	err := ShouldBind(pContext, &item)
	if err != nil {
		session := sessions.Default(pContext)
		session.Set("item", item)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/item/addpage")
	}

	return item, err
}

func Items(pContext *gin.Context) []Item {
	user := GetUser(pContext)
	items := make([]Item, 1)
	err := db.Qry(&items, "SELECT i.id, i.`name`, i.git_id, IFNULL(g.`name`, '') AS 'git_name', i.repo_url, i.branch, i.server_id, IFNULL(s.`name`, '') AS 'server_name', i.rem, i.add_time, i.upd_time FROM item i LEFT JOIN git g ON g.del_flag = 0 AND g.id = i.git_id LEFT JOIN server s ON s.del_flag = 0 AND s.id = i.server_id WHERE i.del_flag = 0 AND i.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if items[0].Id == 0 {
		items = nil
	}

	return items
}

func ParseIniText(iniTxt string) Ini {
	ini := Ini{}
	pReader := bufio.NewReader(bytes.NewBufferString(iniTxt))

	set := func(txt, ty string) {
		switch ty {
		case TagBuild:
			ini.Build = txt

		case TagTarget:
			ini.Target = txt

		case TagDeploy:
			ini.Deploy = txt

		default:
		}
	}

	ty := ""
	txt := ""
	for {
		line, err := pReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			continue
		}

		if line == "" {
			continue
		}

		switch line {
		case TagBuild:
			set(txt, ty)
			txt = ""
			ty = line
			continue

		case TagTarget:
			set(txt, ty)
			txt = ""
			ty = line
			continue

		case TagDeploy:
			set(txt, ty)
			txt = ""
			ty = line
			continue

		default:
			txt += line + "\n"
		}
	}
	set(txt, ty)

	return ini
}
