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
	pEngine.NoRoute(func(pContext *gin.Context) {
		// pContext.HTML() 返回HTML
		pContext.HTML(http.StatusOK, "404.html", gin.H{
			"user":       api.GetUser(pContext),
			"goVersion":  runtime.Version(),
			"ginVersion": gin.Version,
			"time":       time.Now(),
		})
	})

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
		gitRouterGroup.Any("/index", api.GitIndex)
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

	// admin
	adminRouterGroup := pEngine.Group("/admin")
	{
		adminRouterGroup.PUT("/allowregflag/:value", api.AdminAllowRegFlagUpd)
		adminRouterGroup.PUT("/buildlevel/:value", api.AdminBuildLevelUpd)
		adminRouterGroup.PUT("/sudoflag/:value", api.AdminSudoFlagUpd)
		adminRouterGroup.POST("/buildenv/:value", api.AdminBuildEnvAdd)
		adminRouterGroup.DELETE("/buildenv/:value", api.AdminBuildEnvDel)
	}

	// ws
	pEngine.GET("/ws", api.Ws)
}
