// env test
// @author xiangqian
// @date 00:56 2023/01/01
package test

import (
	"auto-deploy-go/src/typ"
	"fmt"
	"testing"
)

func TestEnv(t *testing.T) {
	buildEnv := "docker:test"
	container := buildEnv[len(typ.EnvDocker):]
	fmt.Printf("%v\n", container)
}
