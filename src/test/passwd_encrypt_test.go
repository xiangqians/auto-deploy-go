// Passwd Test
// @author xiangqian
// @date 11:48 2022/12/28
package test

import (
	"auto-deploy-go/src/api"
	"log"
	"testing"
)

func TestPasswdEncrypt(t *testing.T) {
	passwd := "autodeployadmin"
	str := api.PasswdEncrypt(passwd)
	log.Printf("%v\n", str)
	// 64f7156950d2bb280b9459a114361f4b

	passwd = "autodeploydemo"
	str = api.PasswdEncrypt(passwd)
	log.Printf("%v\n", str)
	// 28bd95214db32d1cd10a8301bacbd588
}
