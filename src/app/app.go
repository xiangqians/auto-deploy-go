// App
// https://github.com/gin-gonic/gin
// @author xiangqian
// @date 18:00 2022/12/18
package app

import (
	"auto-deploy-go/src/api"
	"auto-deploy-go/src/logger"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Run() {
	// init
	logger.Init()
	api.InitValidateTrans()

	// Gin ReleaseMode
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	// default Engine
	pEngine := gin.Default()

	// int, middleware, route
	intHtmlTemplate(pEngine)
	userSessionMiddleware(pEngine)
	userI18nMiddleware(pEngine)
	userStaticMiddleware(pEngine)
	userPermMiddleware(pEngine)
	initRoute(pEngine)

	// port
	var port int
	flag.IntVar(&port, "port", 8080, "port")
	flag.Parse()

	// addr
	addr := fmt.Sprintf(":%v", strconv.FormatInt(int64(port), 10))

	// run
	pEngine.Run(addr)
}

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
	pEngine.GET("/", api.IndexPage)

	// user
	userRouterGroup := pEngine.Group("/user")
	{
		userRouterGroup.GET("/regpage", api.UserRegPage)
		userRouterGroup.GET("/loginpage", api.UserLoginPage)
		userRouterGroup.POST("/login", api.UserLogin)
		userRouterGroup.Any("/logout", api.UserLogout)
		userRouterGroup.GET("/stgpage", api.UserStgPage)
	}
	pEngine.POST("/user", api.UserAdd)
	pEngine.PUT("/user", api.UserUpd)

	// git
	gitRouterGroup := pEngine.Group("/git")
	{
		gitRouterGroup.GET("/index", api.GitIndex)
		gitRouterGroup.GET("/addpage", api.GitAddPage)
	}
	pEngine.POST("/git", api.GitAdd)
	pEngine.PUT("/git", api.GitUpd)
	pEngine.DELETE("/git", api.GitDel)

	// server
	serverRouterGroup := pEngine.Group("/server")
	{
		serverRouterGroup.GET("/index", api.ServerIndex)
		serverRouterGroup.GET("/addpage", api.ServerAddPage)
	}
	pEngine.POST("/server", api.ServerAdd)
	pEngine.PUT("/server", api.ServerUpd)
	pEngine.DELETE("/server", api.ServerDel)

	// item
	itemRouterGroup := pEngine.Group("/item")
	{
		itemRouterGroup.GET("/index", api.ItemIndex)
		itemRouterGroup.GET("/addpage", api.ItemAddPage)
	}
	pEngine.POST("/item", api.ItemAdd)
	pEngine.PUT("/item", api.ItemUpd)
	pEngine.DELETE("/item", api.ItemDel)

	// ws
	pEngine.GET("/ws", api.Ws)
}

func userPermMiddleware(pEngine *gin.Engine) {
	// 未授权拦截
	pEngine.Use(func(pContext *gin.Context) {
		reqPath := pContext.Request.URL.Path

		// 静态资源放行
		if strings.HasPrefix(reqPath, "/static") {
			pContext.Next()
			return
		}

		// 获取session对象
		session := sessions.Default(pContext)

		// isLogin
		isLogin := false
		if user, r := session.Get(api.SessionKeyUser).(api.User); r && user.Id != 0 {
			isLogin = true
		}

		if reqPath == "/user/regpage" || reqPath == "/user/loginpage" ||
			//(reqPath == "/user" && pContext.Request.Method == http.MethodPost) ||
			reqPath == "/user/login" {
			if isLogin {
				pContext.Redirect(http.StatusMovedPermanently, "/")
				pContext.Abort()
			} else {
				pContext.Next()
			}
			return
		}

		if !isLogin {
			// 重定向
			//pContext.Request.URL.Path = "/user/loginpage"
			//pEngine.HandleContext(pContext)
			pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")

			// 中止调用链
			pContext.Abort()
		}
	})
}

func userStaticMiddleware(pEngine *gin.Engine) {
	// 静态资源处理
	// https://github.com/gin-contrib/static
	pEngine.Use(static.Serve("/static", static.LocalFile("./static", false)))
}

func userI18nMiddleware(pEngine *gin.Engine) {
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
			// 从url中获取lang
			lang := strings.ToLower(strings.TrimSpace(pContext.Query("lang")))
			if lang != "" && !(lang == api.LocalesZh || lang == api.LocalesEn) {
				lang = ""
			}

			// 从session中获取lang
			//session := sessions.Default(pContext)
			session := sessions.Default(pContext)
			sessionLang := ""
			if v, r := session.Get("lang").(string); r {
				sessionLang = v
			}
			if lang == "" {
				lang = sessionLang
			}

			if lang == "" {
				// 从请求头获取 Accept-Language
				acceptLanguage := pContext.GetHeader("Accept-Language")
				// en,zh-CN;q=0.9,zh;q=0.8
				if strings.HasPrefix(acceptLanguage, api.LocalesZh) {
					lang = api.LocalesZh
				} else if strings.HasPrefix(acceptLanguage, api.LocalesEn) {
					lang = api.LocalesEn
				}
			}

			if lang == "" {
				lang = defaultLang
			}

			if sessionLang != lang {
				session.Set("lang", lang)
				session.Save()
			}
			return lang
		},
	)))
}

func userSessionMiddleware(pEngine *gin.Engine) {
	// 密钥
	keyPairs := []byte("123456")
	// 创建基于cookie的存储引擎
	//store := cookie.NewStore(keyPairs)
	// 创建基于mem（内存）的存储引擎，其实就是一个 map[interface]interface 对象
	store := memstore.NewStore(keyPairs)

	// 设置session中间件
	// session中间件基于内存（其他存储引擎支持：redis、mysql等）实现
	pEngine.Use(sessions.Sessions("autoDeploySessionId", // session & cookie 名称
		store))
}

// 初始化HTML模板
func intHtmlTemplate(pEngine *gin.Engine) {
	// 自定义模板函数，为了获取 i18n 文件中 key 对应的 value
	pEngine.SetFuncMap(template.FuncMap{
		"Localize": i18n.GetMessage,
	})

	// HTML模板
	//pEngine.LoadHTMLGlob("templates/*")
	//pEngine.LoadHTMLGlob("templates/**/*")
	// https://github.com/gin-contrib/multitemplate
	pEngine.HTMLRender = func(templatesDir string) multitemplate.Renderer {
		// if gin.DebugMode -> NewDynamic()
		pRenderer := multitemplate.NewRenderer()

		matches, err := filepath.Glob(templatesDir)
		if err != nil {
			panic(err)
		}

		// Generate our templates map from our layouts/ and includes/ directories
		for _, matche := range matches {
			pFile, ferr := os.Open(matche)
			if ferr != nil {
				continue
			}

			fileInfo, fierr := pFile.Stat()
			if fierr == nil {
				name := filepath.Base(matche)
				// /**/*
				if fileInfo.IsDir() {
					fname := fileInfo.Name()
					subFileInfos, sfierr := pFile.Readdir(-1)
					if sfierr == nil {
						for _, subFileInfo := range subFileInfos {
							subfname := subFileInfo.Name()
							pRenderer.AddFromFilesFuncs(fmt.Sprintf("%s/%s", fname, subfname), pEngine.FuncMap, fmt.Sprintf("%s/%s", matche, subfname))
						}
					}

				} else
				// /*
				{
					pRenderer.AddFromFilesFuncs(name, pEngine.FuncMap, matche)
				}
			}
			pFile.Close()
		}

		return pRenderer
	}("./templates/*")
}
