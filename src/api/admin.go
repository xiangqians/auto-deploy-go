// admin
// @author xiangqian
// @date 13:39 2023/01/07
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AdminIndexPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()

	setting, _ := Setting()
	pContext.HTML(http.StatusOK, "admin/index.html", gin.H{
		"user":      GetUser(pContext),
		"setting":   setting,
		"buildEnvs": BuildEnvs(),
		"message":   message,
	})
}

func AdminAllowRegFlagUpd(pContext *gin.Context) {
	adminByteColumnUpd(pContext, "allow_reg_flag")
	return
}

func AdminBuildLevelUpd(pContext *gin.Context) {
	adminByteColumnUpd(pContext, "build_level")
	return
}

func AdminSudoFlagUpd(pContext *gin.Context) {
	adminByteColumnUpd(pContext, "sudo_flag")
	return
}

func adminByteColumnUpd(pContext *gin.Context, name string) {
	redirect := func(err error) {
		session := sessions.Default(pContext)
		if err != nil {
			session.Set("message", err.Error())
		}
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/")
	}

	if !IsAdminUser(pContext, typ.User{}) {
		redirect(nil)
		return
	}

	// value
	valueStr := pContext.Param("value")
	value, err := strconv.ParseInt(valueStr, 10, 8)
	if err != nil {
		redirect(err)
		return
	}

	_, err = db.Upd(fmt.Sprintf("UPDATE `setting` SET `%s` = ?", name), value)
	if err != nil {
		redirect(err)
		return
	}

	redirect(nil)
	return
}

func AdminBuildEnvAdd(pContext *gin.Context) {
	adminBuildEnvAddOrDel(pContext)
}

func AdminBuildEnvDel(pContext *gin.Context) {
	adminBuildEnvAddOrDel(pContext)
}

func adminBuildEnvAddOrDel(pContext *gin.Context) {
	redirect := func(err error) {
		session := sessions.Default(pContext)
		if err != nil {
			session.Set("message", err.Error())
		}
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/")
	}

	if !IsAdminUser(pContext, typ.User{}) {
		redirect(nil)
		return
	}

	// value
	value := strings.TrimSpace(pContext.Param("value"))
	if value == "" {
		redirect(errors.New("BuildEnv cannot be empty"))
		return
	}

	setting, err := Setting()
	if err != nil {
		redirect(err)
		return
	}

	buildEnvs := setting.BuildEnvs
	contains := buildEnvs != "" && strings.Contains(fmt.Sprintf(",%s,", buildEnvs), fmt.Sprintf(",%s,", value))
	if pContext.Request.Method == http.MethodPost {
		if contains {
			redirect(nil)
			return
		}

		if buildEnvs == "" {
			buildEnvs = value
		} else {
			buildEnvs += "," + value
		}

	} else if pContext.Request.Method == http.MethodDelete {
		if !contains {
			redirect(nil)
			return
		}
		buildEnvs = strings.ReplaceAll(fmt.Sprintf(",%s,", buildEnvs), fmt.Sprintf(",%s,", value), ",")
		if buildEnvs == "," {
			buildEnvs = ""
		} else {
			buildEnvs = buildEnvs[1 : len(buildEnvs)-1]
		}

	} else {
		redirect(nil)
		return
	}

	_, err = db.Upd("UPDATE `setting` SET `build_envs` = ?", buildEnvs)
	if err != nil {
		redirect(err)
		return
	}

	redirect(nil)
	return
}

func BuildEnvs() []typ.BuildEnv {
	setting, err := Setting()
	if err != nil {
		log.Println(err)
		return nil
	}

	buildEnvsStr := setting.BuildEnvs
	if buildEnvsStr == "" {
		return nil
	}

	buildEnvsArr := strings.Split(setting.BuildEnvs, ",")
	buildEnvs := make([]typ.BuildEnv, len(buildEnvsArr))
	for i, buildEnvStr := range buildEnvsArr {
		buildEnvs[i] = typ.BuildEnv{
			Value: buildEnvStr,
		}
	}

	return buildEnvs
}

func Setting() (typ.Setting, error) {
	setting := typ.Setting{}
	err := db.Qry(&setting, "SELECT `sudo_flag`, `allow_reg_flag`, `build_level`, `build_envs` FROM `setting` LIMIT 1")
	if err != nil {
		return setting, err
	}
	return setting, nil
}
