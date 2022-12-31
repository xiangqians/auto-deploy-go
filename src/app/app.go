// App
// https://github.com/gin-gonic/gin
// @author xiangqian
// @date 18:00 2022/12/18
package app

import (
	"auto-deploy-go/src/api"
	"auto-deploy-go/src/arg"
	"auto-deploy-go/src/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Run() {
	// init
	logger.Init()
	api.Init()

	// validation
	//regTrimSpaceValidation()

	// Gin ReleaseMode
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	// default Engine
	pEngine := gin.Default()

	// int
	intHtmlTemplate(pEngine)
	userSessionMiddleware(pEngine)
	userI18nMiddleware(pEngine)
	userStaticMiddleware(pEngine)
	userPermMiddleware(pEngine)
	initRoute(pEngine)

	// arg parse
	arg.Parse()

	// addr
	addr := fmt.Sprintf(":%v", strconv.FormatInt(int64(arg.Port), 10))

	// run
	pEngine.Run(addr)
}
