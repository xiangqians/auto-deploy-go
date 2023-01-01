// env test
// @author xiangqian
// @date 00:56 2023/01/01
package test

import (
	"auto-deploy-go/src/typ"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestEnv(t *testing.T) {
	buildEnv := "docker:test"
	container := buildEnv[len(typ.EnvDocker):]
	fmt.Printf("%v\n", container)
}

func TestResPath(t *testing.T) {
	resPath := "C:\\Users\\xiangqian\\Desktop\tmp\\auto-deploy\\tmp/item1/res"
	index := strings.LastIndex(resPath, "item")
	log.Printf("index: %v\n", index)
	log.Printf("resPath: %v\n", resPath[index:])
}
