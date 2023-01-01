// Template
// @author xiangqian
// @date 21:45 2022/12/23
package app

import (
	"auto-deploy-go/src/typ"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 初始化HTML模板
func intHtmlTemplate(pEngine *gin.Engine) {
	// 自定义模板函数
	pEngine.SetFuncMap(template.FuncMap{
		// 为了获取 i18n 文件中 key 对应的 value
		"Localize": i18n.GetMessage,

		// unix时间戳转日期
		"UnixToTime": func(unix int64) string {
			if unix == 0 {
				return "-"
			}
			t := time.Unix(unix, 0)
			return t.Format("2006/01/02 15:04:05")
		},
		"UnixDiff": func(unix1, unix2 int64) string {
			if unix1 == 0 || unix2 == 0 {
				return "-"
			}

			r := math.Abs(float64(unix1 - unix2)) // s
			return fmt.Sprintf("%ss", strconv.FormatFloat(r, 'f', 2, 64))
		},

		// +1
		"Add1": func(i int) int {
			return i + 1
		},

		// 部署状态文本信息
		"DeployStatusText": func(status byte) string {
			switch status {
			// 1-部署中
			case typ.StatusInDeploy:
				return i18n.MustGetMessage("i18n.inDeploy")

			// 2-部署异常，
			case typ.StatusDeployExc:
				return i18n.MustGetMessage("i18n.deployExc")

			// 3-部署成功
			case typ.StatusDeploySuccess:
				return i18n.MustGetMessage("i18n.deploySuccess")

			default:
				return "-"
			}
		},

		//"Template": func(name string) string {
		//	var data any = nil
		//	re := pEngine.HTMLRender.Instance(name, data)
		//	if html, r := re.(render.HTML); r {
		//		strBuilder := &strings.Builder{}
		//		html.Template.Execute(strBuilder, data)
		//		return strBuilder.String()
		//	}
		//	return ""
		//},
	})

	// HTML模板
	//pEngine.LoadHTMLGlob("templates/*")
	//pEngine.LoadHTMLGlob("templates/**/*")
	// https://github.com/gin-contrib/multitemplate
	pEngine.HTMLRender = func(templatesDir string) render.HTMLRender {
		// if gin.DebugMode -> NewDynamic()
		renderer := multitemplate.NewRenderer()

		matches, err := filepath.Glob(templatesDir + "/*")
		if err != nil {
			panic(err)
		}

		coms, err := filepath.Glob(templatesDir + "/com/*")
		if err != nil {
			panic(err)
		}

		getFiles := func(s string) []string {
			files := make([]string, len(coms)+1)
			i := 0
			files[i] = s
			i++
			for _, e := range coms {
				files[i] = e
				i++
			}
			return files
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
					if fname != "com" {
						subFileInfos, sfierr := pFile.Readdir(-1)
						if sfierr == nil {
							for _, subFileInfo := range subFileInfos {
								subfname := subFileInfo.Name()
								files := getFiles(fmt.Sprintf("%s/%s", matche, subfname))
								renderer.AddFromFilesFuncs(fmt.Sprintf("%s/%s", fname, subfname), pEngine.FuncMap, files...)
							}
						}
					}
				} else
				// /*
				{
					files := getFiles(matche)
					renderer.AddFromFilesFuncs(name, pEngine.FuncMap, files...)
				}
			}
			pFile.Close()
		}

		return renderer
	}("./templates")
}
