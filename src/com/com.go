// common
// @author xiangqian
// @date 22:31 2022/12/20
package com

import "log"

// const DataDir = "./data"
const DataDir = "C:\\Users\\xiangqian\\Desktop\\tmp\\auto-deploy\\data"

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
