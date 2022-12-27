// Git
// @author xiangqian
// @date 11:48 2022/12/22
package api

import (
	"auto-deploy-go/src/db"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Git struct {
	Abs
	UserId int64  // Git所属用户id
	Name   string `form:"name" binding:"required,min=1,max=60"`    // 名称
	User   string `form:"user" binding:"required,min=1,max=60"`    // 用户
	Passwd string `form:"passwd" binding:"required,min=1,max=100"` // 密码
}

func init() {
	// 注册 Git 模型
	gob.Register(Git{})
}

func GitIndex(pContext *gin.Context) {
	pContext.HTML(http.StatusOK, "git/index.html", gin.H{
		"gits": Gits(pContext),
	})
}

func GitAddPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	git := session.Get("git")
	message := session.Get("message")
	session.Delete("git")
	session.Delete("message")
	session.Save()

	if git == nil {
		_git := Git{}
		idStr := pContext.Query("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil && id > 0 {
			user := GetUser(pContext)
			err = db.Qry(&_git, "SELECT g.id, g.`name`, g.`user`, g.rem, g.add_time, g.upd_time FROM git g WHERE g.del_flag = 0 AND g.user_id = ? AND g.id = ?", user.Id, id)
			if err != nil {
				log.Println(err)
			}
		}
		git = _git
	}

	pContext.HTML(http.StatusOK, "git/add.html", gin.H{
		"git":     git,
		"message": message,
	})
}

func GitAdd(pContext *gin.Context) {
	git, err := gitPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Add("INSERT INTO `git` (`user_id`, `name`, `user`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?, ?)",
		user.Id, git.Name, git.User, git.Passwd, git.Rem, time.Now().Unix())
	pContext.Redirect(http.StatusMovedPermanently, "/git/index")
}

func GitUpd(pContext *gin.Context) {
	git, err := gitPreAddOrUpd(pContext)
	if err != nil {
		return
	}

	user := GetUser(pContext)
	db.Upd("UPDATE git SET `name` = ?, `user` = ?, `passwd` = ?, `rem` = ?, upd_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
		git.Name, git.User, git.Passwd, git.Rem, time.Now().Unix(), user.Id, git.Id)
	pContext.Redirect(http.StatusMovedPermanently, "/git/index")
}

func GitDel(pContext *gin.Context) {
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil {
		user := GetUser(pContext)
		db.Del("UPDATE git SET del_flag = 1, upd_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	}
	pContext.Redirect(http.StatusMovedPermanently, "/git/index")
}

func gitPreAddOrUpd(pContext *gin.Context) (Git, error) {
	git := Git{}
	err := ShouldBind(pContext, &git)

	git.Name = strings.TrimSpace(git.Name)
	git.User = strings.TrimSpace(git.User)
	git.Passwd = strings.TrimSpace(git.Passwd)
	git.Rem = strings.TrimSpace(git.Rem)

	if err != nil {
		session := sessions.Default(pContext)
		session.Set("git", git)
		session.Set("message", err.Error())
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/git/addpage")
	}

	return git, err
}

func Gits(pContext *gin.Context) []Git {
	user := GetUser(pContext)
	gits := make([]Git, 1)
	err := db.Qry(&gits, "SELECT g.id, g.`name`, g.`user`, g.rem, g.add_time, g.upd_time FROM git g WHERE g.del_flag = 0 AND g.user_id = ?", user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if gits[0].Id == 0 {
		gits = nil
	}

	return gits
}
