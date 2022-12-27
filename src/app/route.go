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
			"goVersion":  runtime.Version(),
			"ginVersion": gin.Version,
			"time":       time.Now(),
		})
	})

	// index
	pEngine.Any("/", api.IndexPage)

	// user
	userRouterGroup := pEngine.Group("/user")
	{
		userRouterGroup.Any("/regpage", api.UserRegPage)
		userRouterGroup.Any("/loginpage", api.UserLoginPage)
		userRouterGroup.POST("/login", api.UserLogin)
		userRouterGroup.Any("/logout", api.UserLogout)
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
	}
	pEngine.POST("/rx", api.RxAdd)
	pEngine.PUT("/rx", api.RxUpd)
	pEngine.DELETE("/rx/:id", api.RxDel)

	// ws
	pEngine.GET("/ws", api.Ws)
}
