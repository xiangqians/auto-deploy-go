// Template
// @author xiangqian
// @date 21:45 2022/12/23
package app

import (
	"auto-deploy-go/src/api"
	"auto-deploy-go/src/typ"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func Add(i1 any, i2 int64) int64 {
	if v, r := i1.(int); r {
		return int64(v) + i2
	}

	if v, r := i1.(int8); r {
		return int64(v) + i2
	}

	if v, r := i1.(int16); r {
		return int64(v) + i2
	}

	if v, r := i1.(int32); r {
		return int64(v) + i2
	}

	if v, r := i1.(int64); r {
		return v + i2
	}

	return 0
}

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

			//r := math.Abs(float64(unix1 - unix2)) // s
			//return fmt.Sprintf("%ss", strconv.FormatFloat(r, 'f', 2, 64))

			r := unix1 - unix2
			if r < 0 {
				r = -r
			}
			return fmt.Sprintf("%ss", strconv.FormatInt(r, 10))
		},

		// +1
		"Add1": func(i any) int64 {
			return Add(i, 1)
		},

		"Minus1": func(i any) int64 {
			return Add(i, -1)
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

		"ItemTime": func(record typ.Record) string {
			unixDiff := func(stime, etime int64, status byte) int64 {
				if stime == 0 {
					return 0
				}

				if etime == 0 {
					etime = time.Now().Unix()
				}
				return etime - stime
			}

			var r int64 = 0

			// pull
			r += unixDiff(record.PullStime, record.PullEtime, record.PullStatus)
			// build
			r += unixDiff(record.BuildStime, record.BuildEtime, record.BuildStatus)
			// pack
			r += unixDiff(record.PackStime, record.PackEtime, record.PackStatus)
			// ul
			r += unixDiff(record.UlStime, record.UlEtime, record.UlStatus)
			// unpack
			r += unixDiff(record.UnpackStime, record.UnpackEtime, record.UnpackStatus)
			// deploy
			r += unixDiff(record.DeployStime, record.DeployEtime, record.DeployStatus)
			return fmt.Sprintf("%ss", strconv.FormatInt(r, 10))
		},

		"Put": func(h gin.H, key string, value any) string {
			h[key] = value
			return ""
		},

		"IsAdminUser": func(user typ.User) bool {
			return api.IsAdminUser(nil, user)
		},

		"Length": func(i any) int {
			if v, r := i.(string); r {
				return len(v)
			}
			return 0
		},

		"SubString": func(str string, start, end int, fill string) string {
			if start >= 0 && start <= end && end <= len(str) {
				return str[start:end] + fill
			}
			return ""
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
					subFileInfos, sfierr := pFile.Readdir(-1)
					if sfierr == nil {
						for _, subFileInfo := range subFileInfos {
							subfname := subFileInfo.Name()
							var files []string
							if fname == "com" {
								files = []string{fmt.Sprintf("%s/%s", matche, subfname)}
							} else {
								files = getFiles(fmt.Sprintf("%s/%s", matche, subfname))
							}
							renderer.AddFromFilesFuncs(fmt.Sprintf("%s/%s", fname, subfname), pEngine.FuncMap, files...)
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
