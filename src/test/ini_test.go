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
mvn clean
mvn package

[target]
./target/jenkins-test-1.0-SNAPSHOT.jar
target/jenkins-test-1.0-SNAPSHOT.jar.original
./target/classes
target/maven-archiver

[deploy]
3

4
mv ./test test`
	ini := api.ParseIniText(iniText)
	log.Printf("Build:\n%v\n", ini.Build)
	log.Printf("Target:\n%v\n", ini.Target)
	log.Printf("Deploy:\n%v\n", ini.Deploy)
}
