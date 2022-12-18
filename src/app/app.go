// App
// @author xiangqian
// @date 18:00 2022/12/18
package app

import (
	"auto-deploy-go/src/api"
	"auto-deploy-go/src/logger"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/text/language"
	"html/template"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Run() {
	// logger
	logger.Init()

	// Gin ReleaseMode
	gin.SetMode(gin.ReleaseMode)

	// default Engine
	pEngine := gin.Default()

	// route
	route(pEngine)

	// port
	var port int
	flag.IntVar(&port, "port", 8080, "port")
	flag.Parse()

	// addr
	addr := fmt.Sprintf(":%v", strconv.FormatInt(int64(port), 10))

	// run
	pEngine.Run(addr)
}

func route(pEngine *gin.Engine) {
	// 自定义模板函数，为了获取 i18n 文件中 key 对应的 value
	pEngine.SetFuncMap(template.FuncMap{
		"Localize": i18n.GetMessage,
	})

	// HTML模板
	pEngine.LoadHTMLGlob("templates/*")

	// 创建基于cookie的存储引擎
	keyPairs := []byte("123456") // 密钥
	store := cookie.NewStore(keyPairs)
	// 设置session中间件
	// session中间件基于内存（其他存储引擎支持：redis、mysql等）实现时，其实就是一个 map[interface]interface 对象
	pEngine.Use(sessions.Sessions("auto-deploy-session", // session & cookie名字
		store))
	gob.Register(api.User{})

	// 未授权拦截
	pEngine.Use(func(pContext *gin.Context) {
		reqPath := pContext.Request.URL.Path
		if strings.HasPrefix(reqPath, "/static") {
			pContext.Next()
			return
		}

		// 获取session对象
		session := sessions.Default(pContext)

		// isLogin
		isLogin := false
		username := session.Get("_username")
		if v, r := username.(string); r && v != "" {
			isLogin = true
		}

		if isLogin && (reqPath == "/user/reg" || reqPath == "/user/login") {
			pContext.Redirect(http.StatusFound, "/")
			return
		}

		if isLogin {
			if reqPath == "/user/reg" || reqPath == "/user/login" {
				pContext.Redirect(http.StatusFound, "/")
			} else {
				pContext.Next()
			}
			return
		}

		if reqPath == "/user/reg" || reqPath == "/user/login" {
			pContext.Next()
			return
		}

		// 重定向
		pContext.Redirect(http.StatusFound, "/user/login")
	})

	// apply i18n middleware
	// https://github.com/gin-contrib/i18n
	pEngine.Use(i18n.Localize(i18n.WithBundle(&i18n.BundleCfg{
		RootPath:         "./i18n",
		AcceptLanguage:   []language.Tag{language.Chinese, language.English},
		DefaultLanguage:  language.Chinese,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	}), i18n.WithGetLngHandle(
		func(pContext *gin.Context, defaultLang string) string {
			lang := pContext.Query("lang")
			if lang == "" {
				return defaultLang
			}
			return lang
		},
	)))

	// 静态资源处理
	// https://github.com/gin-contrib/static
	pEngine.Use(static.Serve("/static", static.LocalFile("./static", false)))

	// 设置默认路由
	pEngine.NoRoute(func(pContext *gin.Context) {
		// pContext.HTML() 返回HTML
		pContext.HTML(http.StatusOK, "404.html", gin.H{
			"goVersion":  runtime.Version(),
			"ginVersion": gin.Version,
			"time":       time.Now(),
		})
	})

	// user
	userRouterGroup := pEngine.Group("/user")
	{
		userRouterGroup.POST("/reg", api.UserReg)
		userRouterGroup.GET("/login", api.UserLoginHtml)
		userRouterGroup.POST("/login", api.UserLogin)
		userRouterGroup.GET("/logout", api.UserLogout)
		userRouterGroup.POST("/logout", api.UserLogout)
	}

	// index
	pEngine.GET("/", api.Index)

}

func Uuid() string {
	return uuid.New().String()
}
