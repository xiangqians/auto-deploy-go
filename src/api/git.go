// Git
// @author xiangqian
// @date 11:48 2022/12/22
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func GitList(pContext *gin.Context) {
	HtmlPage(pContext, "git/list.html", func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error) {
		page, err := GitPage(pContext, pageReq)
		return page, nil, err
	})
}

func GitAddPage(pContext *gin.Context) {
	git, err := Session[typ.Git](pContext, "git", true)
	if err != nil {
		id, _ := Query[int64](pContext, "id")
		if id > 0 {
			_, git, err = Git(pContext, id)
			if err != nil {
				log.Println(err)
			}
		}
	}

	Html(pContext, "git/add.html", func(pContext *gin.Context) (gin.H, error) {
		return gin.H{"git": git}, nil
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

	user := SessionUser(pContext)
	if pContext.Request.Method == http.MethodPost {
		_, err = db.Add("INSERT INTO `git` (`user_id`, `name`, `user`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?, ?)",
			user.Id, git.Name, git.User, git.Passwd, git.Rem, time.Now().Unix())
	} else if pContext.Request.Method == http.MethodPut {
		_, err = db.Upd("UPDATE git SET `name` = ?, `user` = ?, `passwd` = ?, `rem` = ?, upd_time = ? WHERE del_flag = 0 AND user_id = ? AND id = ?",
			git.Name, git.User, git.Passwd, git.Rem, time.Now().Unix(), user.Id, git.Id)
	}

	Redirect(pContext, "/git/list", err, nil)
	return
}

func GitDel(pContext *gin.Context) {
	redirect := func(message any) {
		Redirect(pContext, "/git/list", message, nil)
	}

	id, err := Param[int64](pContext, "id")
	if err != nil || id <= 0 {
		redirect(err)
		return
	}

	user := SessionUser(pContext)
	db.Del("UPDATE git SET del_flag = 1, upd_time = ? WHERE user_id = ? AND id = ?", time.Now().Unix(), user.Id, id)
	redirect(nil)
	return
}

func GitPage(pContext *gin.Context, pageReq typ.PageReq) (typ.Page[typ.Git], error) {
	return db.Page[typ.Git](pageReq, GitSql(pContext, 0))
}

func Git(pContext *gin.Context, id int64) (int64, typ.Git, error) {
	git := typ.Git{}
	count, err := db.Qry(&git, GitSql(pContext, id))
	return count, git, err
}

func GitSql(pContext *gin.Context, id int64) string {
	user := SessionUser(pContext)
	sql := fmt.Sprintf("SELECT g.id, g.`name`, g.`user`, g.rem, g.add_time, g.upd_time FROM git g WHERE g.del_flag = 0 AND g.user_id = %v ", user.Id)
	if id != 0 {
		sql += fmt.Sprintf("AND g.id = %v ", id)
	}
	return sql
}
