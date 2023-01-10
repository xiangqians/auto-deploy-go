// Git
// @author xiangqian
// @date 11:48 2022/12/22
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GitIndex(pContext *gin.Context) {
	HtmlPage(pContext, "git/index.html", func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error) {
		page, err := PageGit(pContext, pageReq)
		return page, nil, err
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
		_git := typ.Git{}
		idStr := pContext.Query("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil && id > 0 {
			user := GetUser(pContext)
			_, err = db.Qry(&_git, "SELECT g.id, g.`name`, g.`user`, g.rem, g.add_time, g.upd_time FROM git g WHERE g.del_flag = 0 AND g.user_id = ? AND g.id = ?", user.Id, id)
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
	GitAddOrUpd(pContext)
}

func GitUpd(pContext *gin.Context) {
	GitAddOrUpd(pContext)
}

func GitAddOrUpd(pContext *gin.Context) {
	redirect := func(git typ.Git, message any) {
		Redirect(pContext, "/git/addpage", message, map[string]any{"git": git})
	}

	git := typ.Git{}
	err := ShouldBind(pContext, &git)
	git.Name = strings.TrimSpace(git.Name)
	git.User = strings.TrimSpace(git.User)
	git.Passwd = strings.TrimSpace(git.Passwd)
	git.Rem = strings.TrimSpace(git.Rem)
	if err != nil {
		redirect(git, err)
		return
	}

	user := GetUser(pContext)
	if pContext.Request.Method == http.MethodPost {
		_, err = db.Add("INSERT INTO `git` (`user_id`, `name`, `user`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?, ?)",
			user.Id, git.Name, git.User, git.Passwd, git.Rem, time.Now().Unix())
	} else if pContext.Request.Method == http.MethodPut {
		_, err = db.Upd("UPDATE git SET `name` = ?, `user` = ?, `passwd` = ?, `rem` = ?, upd_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
			git.Name, git.User, git.Passwd, git.Rem, time.Now().Unix(), user.Id, git.Id)
	}

	Redirect(pContext, "/git/index", err, nil)
	return
}

func GitDel(pContext *gin.Context) {
	redirect := func(message any) {
		Redirect(pContext, "/git/index", message, nil)
	}

	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		redirect(err)
		return
	}

	user := GetUser(pContext)
	db.Del("UPDATE git SET del_flag = 1, upd_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	redirect(nil)
	return
}

func PageGit(pContext *gin.Context, pageReq typ.PageReq) (typ.Page[typ.Git], error) {
	user := GetUser(pContext)
	return db.Page[typ.Git](pageReq, "SELECT g.id, g.`name`, g.`user`, g.rem, g.add_time, g.upd_time FROM git g WHERE g.del_flag = 0 AND g.user_id = ?", user.Id)
}
