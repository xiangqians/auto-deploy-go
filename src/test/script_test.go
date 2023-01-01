// Script
// @author xiangqian
// @date 12:16 2022/12/28
package test

import (
	"auto-deploy-go/src/deploy"
	"log"
	"testing"
)

func TestIni(t *testing.T) {
	iniText := `
[build]
mvn clean
#mvn package

[target]
./target/jenkins-test-1.0-SNAPSHOT.jar
target/jenkins-test-1.0-SNAPSHOT.jar.original
#./target/classes
#target/maven-archiver

[deploy]
#! /bin/bash
#3

4
mv ./test test`
	script := deploy.ParseScriptTxt(iniText)
	log.Printf("Build:\n%v\n", script.Build)
	log.Printf("Target:\n%v\n", script.Target)
	log.Printf("Deploy:\n%v\n", script.Deploy)
}
