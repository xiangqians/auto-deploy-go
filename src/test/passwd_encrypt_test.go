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
	passwd := "admin"
	str := api.PasswdEncrypt(passwd)
	log.Printf("%v\n", str)
	// 75b17d369a5ce9b50e1a608bee111cac

	passwd = "test"
	str = api.PasswdEncrypt(passwd)
	log.Printf("%v\n", str)
	// 682a990b43354f908437490c14bf3019
}
