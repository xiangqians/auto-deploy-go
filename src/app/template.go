// Template
// @author xiangqian
// @date 21:45 2022/12/23
package app

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"path/filepath"
)

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
