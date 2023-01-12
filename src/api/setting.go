// setting
// @author xiangqian
// @date 13:39 2023/01/07
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SettingIndexPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()

	setting, _ := Setting()
	pContext.HTML(http.StatusOK, "setting/index.html", gin.H{
		"user":    SessionUser(pContext),
		"setting": setting,
		"message": message,
	})
}

func SettingAllowRegFlagUpd(pContext *gin.Context) {
	settingByteColumnUpd(pContext, "allow_reg_flag")
	return
}

func SettingBuildLevelUpd(pContext *gin.Context) {
	settingByteColumnUpd(pContext, "build_level")
	return
}

func SettingSudoFlagUpd(pContext *gin.Context) {
	settingByteColumnUpd(pContext, "sudo_flag")
	return
}

func settingByteColumnUpd(pContext *gin.Context, name string) {
	redirect := func(err error) {
		session := sessions.Default(pContext)
		if err != nil {
			session.Set("message", err.Error())
		}
		session.Save()
		pContext.Redirect(http.StatusMovedPermanently, "/setting/index")
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

func Setting() (typ.Setting, error) {
	setting := typ.Setting{}
	_, err := db.Qry(&setting, "SELECT `sudo_flag`, `allow_reg_flag`, `build_level` FROM `setting` LIMIT 1")
	if err != nil {
		return setting, err
	}
	return setting, nil
}
