// Ini
// @author xiangqian
// @date 12:16 2022/12/28
package test

import (
	"auto-deploy-go/src/api"
	"log"
	"testing"
)

func TestIni(t *testing.T) {
	iniText := `
[build]
1


[target]
2


[deploy]
3

4

`
	ini := api.ParseIniText(iniText)
	log.Printf("Build:\n%v\n", ini.Build)
	log.Printf("Target:\n%v\n", ini.Target)
	log.Printf("Deploy:\n%v\n", ini.Deploy)
}
