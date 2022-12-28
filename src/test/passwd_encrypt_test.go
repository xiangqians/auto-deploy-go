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
	passwd := "12345678"
	str := api.PasswdEncrypt(passwd)
	log.Printf("%v\n", str)
	// 59f0f252c669f2908f5d211cf4eae714
}
