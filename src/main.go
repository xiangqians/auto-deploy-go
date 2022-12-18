// main
// @author xiangqian
// @date 19:28 2022/12/03
package main

import (
	"auto-deploy-go/src/logger"
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func main() {
	// logger
	logger.Init()

	// default
	pEngine := gin.Default()

	// port
	var port int
	flag.IntVar(&port, "port", 8080, "port")
	flag.Parse()
	addr := ":" + strconv.FormatInt(int64(port), 10)

	pEngine.GET("/", func(pContext *gin.Context) {
		pContext.JSON(http.StatusOK, gin.H{
			"method":     "GET",
			"goVersion":  runtime.Version(),
			"ginVersion": gin.Version,
			"time":       time.Now(),
		})
	})

	// run
	pEngine.Run(addr)
}
