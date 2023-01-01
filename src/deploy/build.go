// build
// @author xiangqian
// @date 21:42 2022/12/31
package deploy

import (
	"auto-deploy-go/src/arg"
	"auto-deploy-go/src/typ"
	"auto-deploy-go/src/util"
	"errors"
	"fmt"
	"strings"
)

func Build(script typ.Script, recordId int64, resPath string) error {
	if strings.HasPrefix(arg.BuildEnv, typ.EnvDocker) {
		container := arg.BuildEnv[len(typ.EnvDocker):]
		return dockerBuild(script, recordId, resPath, container)

	} else if arg.BuildEnv == typ.EnvDefault {
		return defaultBuild(script, recordId, resPath)
	}

	return errors.New(fmt.Sprintf("unknown -buildenv %v", arg.BuildEnv))
}

// DefaultBuild 默认构建
// 将会在 auto-deploy-go 应用所在的服务器上构建资源，
// 存在很大的风险，如果 [build] 脚本有刻意的命令，将会导致 auto-deploy-go db数据泄露，甚至服务器崩溃
func defaultBuild(script typ.Script, recordId int64, resPath string) error {
	updSTime(typ.StepBuild, recordId)

	_build := script.Build
	if _build != nil && len(_build) > 0 {
		for _, cmd := range _build {
			cd, err := util.Cd(resPath)
			if err != nil {
				updETime(typ.StepBuild, recordId, err)
				return err
			}

			cmd = fmt.Sprintf("%s && %s", cd, cmd)
			pCmd, err := util.Command(cmd)
			if err != nil {
				updETime(typ.StepBuild, recordId, err)
				return err
			}

			_, err = pCmd.CombinedOutput()
			if err != nil {
				updETime(typ.StepBuild, recordId, err)
				return err
			}
		}
	}

	updETime(typ.StepBuild, recordId, nil)
	return nil
}

// DockerBuild docker容器构建
// 将源码cp到docker容器内， 并在docker容器内执行 [build] 脚本，将 build 结果再cp到 auto-deploy-go 应用所在的服务器上
// 低风险
func dockerBuild(script typ.Script, recordId int64, resPath string, container string) error {
	updSTime(typ.StepBuild, recordId)

	// linux在宿主机执行docker容器环境内命令
	// sudo docker exec -it -u root auto-deploy-build-env /bin/bash -c "./test.sh"

	updETime(typ.StepBuild, recordId, nil)
	return nil
}
