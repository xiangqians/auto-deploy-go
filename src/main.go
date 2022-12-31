// main
// @author xiangqian
// @date 19:28 2022/12/03
package main

import "auto-deploy-go/src/app"

func main() {

	app.Run()

}

// Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
// https://github.com/mattn/go-sqlite3/issues/855
// https://github.com/mattn/go-sqlite3/issues/975
// require (
//		github.com/mattn/go-sqlite3 v2.0.3+incompatible
// )
//
// 解决方法1：拉取其他版本
// https://github.com/mattn/go-sqlite3
// Latest stable version is v1.14 or later, not v2.
// go get github.com/mattn/go-sqlite3@v1.14.16
//
// 解决方法2：在不同系统构建不同可执行包
