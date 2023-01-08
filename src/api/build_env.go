// build env
// @author xiangqian
// @date 21:07 2023/01/07
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func BuildEnvIndex(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()
	pContext.HTML(http.StatusOK, "buildenv/index.html", gin.H{
		"user":      GetUser(pContext),
		"buildEnvs": BuildEnvs(),
		"message":   message,
	})
}

func BuildEnvAddPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	buildEnv := session.Get("buildEnv")
	message := session.Get("message")
	session.Delete("message")
	session.Save()

	if buildEnv == nil {
		buildEnv = typ.BuildEnv{}
		idStr := strings.TrimSpace(pContext.Query("id"))
		if idStr != "" {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err == nil && id > 0 {
				buildEnv, err = BuildEnv(id)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	pContext.HTML(http.StatusOK, "buildenv/add.html", gin.H{
		"message":  message,
		"buildEnv": buildEnv,
	})
}

func BuildEnvAdd(pContext *gin.Context) {
	buildEnvAddOrUpd(pContext)
}

func BuildEnvUpd(pContext *gin.Context) {
	buildEnvAddOrUpd(pContext)
}

func BuildEnvDel(pContext *gin.Context) {
	idStr := pContext.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err == nil {
		_, err = db.Del("DELETE FROM build_env WHERE id = ?", id)
	}

	if err != nil {
		session := sessions.Default(pContext)
		session.Set("message", err.Error())
		session.Save()
	}
	pContext.Redirect(http.StatusMovedPermanently, "/buildenv/index")
}

func buildEnvAddOrUpd(pContext *gin.Context) {
	buildEnv := typ.BuildEnv{}
	err := ShouldBind(pContext, &buildEnv)

	buildEnv.Value = strings.TrimSpace(buildEnv.Value)
	buildEnv.Rem = strings.TrimSpace(buildEnv.Rem)

	redirectToAddPage := func(buildEnv typ.BuildEnv, err error) {
		session := sessions.Default(pContext)
		session.Set("buildEnv", buildEnv)
		if err != nil {
			session.Set("message", err.Error())
		}
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/buildenv/addpage")
	}

	if err != nil {
		redirectToAddPage(buildEnv, err)
		return
	}

	// add
	if pContext.Request.Method == http.MethodPost {
		_, err = db.Add("INSERT INTO `build_env` (`value`, `rem`, `add_time`) VALUES (?, ?, ?)", buildEnv.Value, buildEnv.Rem, time.Now().Unix())

	} else
	// upd
	if pContext.Request.Method == http.MethodPut {
		_, err = db.Upd("UPDATE build_env SET `value` = ?, `rem` = ?, upd_time = ? WHERE id = ?", buildEnv.Value, buildEnv.Rem, time.Now().Unix(), buildEnv.Id)
	}

	if err != nil {
		redirectToAddPage(buildEnv, err)
		return
	}

	pContext.Redirect(http.StatusMovedPermanently, "/buildenv/index")
}

func BuildEnv(id int64) (typ.BuildEnv, error) {
	buildEnv := typ.BuildEnv{}
	err := db.Qry(&buildEnv, "SELECT be.`id`, be.`value`, be.`rem`, be.`disable_flag`, be.`add_time`, be.`upd_time` FROM `build_env` be WHERE be.id = ?", id)
	if err != nil {
		return buildEnv, err
	}

	if buildEnv.Id == 0 {
		return buildEnv, errors.New("no record")
	}

	return buildEnv, nil
}

func BuildEnvs() []typ.BuildEnv {
	buildEnvs := make([]typ.BuildEnv, 1)
	err := db.Qry(&buildEnvs, "SELECT be.`id`, be.`value`, be.`rem`, be.`disable_flag`, be.`add_time`, be.`upd_time` FROM `build_env` be")
	if err != nil {
		log.Println(err)
		return nil
	}

	if buildEnvs[0].Id == 0 {
		return nil
	}

	return buildEnvs
}
