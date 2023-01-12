// Route
// @author xiangqian
// @date 21:47 2022/12/23
package app

import (
	"auto-deploy-go/src/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

func initRoute(pEngine *gin.Engine) {
	// 设置默认路由
	handler := func(pContext *gin.Context) {
		// pContext.HTML() 返回HTML
		pContext.HTML(http.StatusOK, "404.html", gin.H{
			"user":       api.GetUser(pContext),
			"goVersion":  runtime.Version(),
			"ginVersion": gin.Version,
			"time":       time.Now(),
		})
	}
	pEngine.Any("/404", handler)
	pEngine.NoRoute(handler)

	// index
	pEngine.Any("/", api.IndexPage)
	pEngine.POST("/deploy/:itemId", api.Deploy)

	// user
	userRouterGroup := pEngine.Group("/user")
	{
		userRouterGroup.Any("/regpage", api.UserRegPage)
		userRouterGroup.Any("/loginpage", api.UserLoginPage)
		userRouterGroup.POST("/login", api.UserLogin)
		userRouterGroup.POST("/logout", api.UserLogout)
		userRouterGroup.Any("/stgpage", api.UserStgPage)
	}
	pEngine.POST("/user", api.UserAdd)
	pEngine.PUT("/user", api.UserUpd)

	// git
	gitRouterGroup := pEngine.Group("/git")
	{
		gitRouterGroup.Any("/list", api.GitList)
		gitRouterGroup.Any("/addpage", api.GitAddPage)
	}
	pEngine.POST("/git", api.GitAdd)
	pEngine.PUT("/git", api.GitUpd)
	pEngine.DELETE("/git/:id", api.GitDel)

	// server
	serverRouterGroup := pEngine.Group("/server")
	{
		serverRouterGroup.Any("/index", api.ServerIndex)
		serverRouterGroup.Any("/addpage", api.ServerAddPage)
	}
	pEngine.POST("/server", api.ServerAdd)
	pEngine.PUT("/server", api.ServerUpd)
	pEngine.DELETE("/server/:id", api.ServerDel)

	// item
	itemRouterGroup := pEngine.Group("/item")
	{
		itemRouterGroup.Any("/index", api.ItemIndex)
		itemRouterGroup.Any("/addpage", api.ItemAddPage)
	}
	pEngine.POST("/item", api.ItemAdd)
	pEngine.PUT("/item", api.ItemUpd)
	pEngine.DELETE("/item/:id", api.ItemDel)

	// rx
	rxRouterGroup := pEngine.Group("/rx")
	{
		rxRouterGroup.Any("/index", api.RxIndex)
		rxRouterGroup.Any("/addpage", api.RxAddPage)
		rxRouterGroup.POST("/join", api.RxJoin)
		rxRouterGroup.Any("/shareitempage", api.RxShareItemPage)
		rxRouterGroup.GET("/notshareitems/:id", api.RxNotShareItems)
		rxRouterGroup.POST("/shareitem/:id/:itemId", api.RxShareItemAdd)
		rxRouterGroup.DELETE("/shareitem/:id/:itemId", api.RxShareItemDel)
	}
	pEngine.POST("/rx", api.RxAdd)
	pEngine.PUT("/rx", api.RxUpd)
	pEngine.DELETE("/rx/:id", api.RxDel)

	// record
	recordRouterGroup := pEngine.Group("/record")
	{
		recordRouterGroup.Any("/index/:itemId", api.RecordIndex)
	}
	pEngine.DELETE("/record/:itemId/:id", api.RecordDel)

	// setting
	settingRouterGroup := pEngine.Group("/setting")
	{
		settingRouterGroup.Any("/index", api.SettingIndexPage)
		settingRouterGroup.PUT("/allowregflag/:value", api.SettingAllowRegFlagUpd)
		settingRouterGroup.PUT("/buildlevel/:value", api.SettingBuildLevelUpd)
		settingRouterGroup.PUT("/sudoflag/:value", api.SettingSudoFlagUpd)
	}

	// buildenv
	buildEnvRouterGroup := pEngine.Group("/buildenv")
	{
		buildEnvRouterGroup.Any("/index", api.BuildEnvIndex)
		buildEnvRouterGroup.Any("/addpage", api.BuildEnvAddPage)
		buildEnvRouterGroup.PUT("/:id/enable", api.BuildEnvEnableOrDisable)
		buildEnvRouterGroup.PUT("/:id/disable", api.BuildEnvEnableOrDisable)
	}
	pEngine.POST("/buildenv", api.BuildEnvAdd)
	pEngine.PUT("/buildenv", api.BuildEnvUpd)
	pEngine.DELETE("/buildenv/:id", api.BuildEnvDel)

	// ws
	pEngine.GET("/ws", api.Ws)
}
