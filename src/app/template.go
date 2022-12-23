// Template
// @author xiangqian
// @date 21:45 2022/12/23
package app

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

// 初始化HTML模板
func intHtmlTemplate(pEngine *gin.Engine) {
	// 自定义模板函数，为了获取 i18n 文件中 key 对应的 value
	pEngine.SetFuncMap(template.FuncMap{
		"Localize": i18n.GetMessage,
		"UnixToTime": func(unix int64) string {
			if unix == 0 {
				return ""
			}
			t := time.Unix(unix, 0)
			return t.Format("2006/01/02 15:04:05")
		},
		"Add1": func(i int) int {
			return i + 1
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
							renderer.AddFromFilesFuncs(fmt.Sprintf("%s/%s", fname, subfname), pEngine.FuncMap, fmt.Sprintf("%s/%s", matche, subfname))
						}
					}

				} else
				// /*
				{
					renderer.AddFromFilesFuncs(name, pEngine.FuncMap, matche)
				}
			}
			pFile.Close()
		}

		return renderer
	}("./templates/*")
}
