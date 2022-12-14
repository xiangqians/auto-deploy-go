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
	"log"
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

	build := script.Build
	if build != nil && len(build) > 0 {
		for _, cmd := range build {
			cd, err := util.Cd(resPath)
			if err != nil {
				updETime(typ.StepBuild, recordId, err, nil)
				return err
			}

			cmd = fmt.Sprintf("%s && %s", cd, cmd)
			pCmd, err := util.Command(cmd)
			if err != nil {
				updETime(typ.StepBuild, recordId, err, nil)
				return err
			}

			buf, err := pCmd.CombinedOutput()
			if err != nil {
				updETime(typ.StepBuild, recordId, err, buf)
				return err
			}
		}
	}

	updETime(typ.StepBuild, recordId, nil, nil)
	return nil
}

// DockerBuild docker容器构建
// 将源码cp到docker容器内， 并在docker容器内执行 [build] 脚本，将 build 结果再cp到 auto-deploy-go 应用所在的服务器上
// 低风险
func dockerBuild(script typ.Script, recordId int64, resPath string, container string) error {
	updSTime(typ.StepBuild, recordId)

	build := script.Build
	if build != nil && len(build) > 0 {
		// 判断是否支持 sudo 命令
		sudo := true
		pCmd, err := util.Command("sudo")
		if err != nil {
			log.Println(err)
			sudo = false
		}
		_, err = pCmd.CombinedOutput()
		if err != nil {
			log.Println(err)
			sudo = false
		}
		log.Printf("sudo: %v\n", sudo)

		// 容器资源路径
		containerResPath := fmt.Sprintf("/tmp/auto-deploy/%s", resPath[strings.LastIndex(resPath, "item"):])
		log.Printf("containerResPath: %v\n", containerResPath)

		// 容器文件拷贝（容器启动与否，拷贝命令都会生效）
		// 1、从宿主机拷贝文件到容器
		// $ docker cp [OPTIONS] SRC_PATH|- CONTAINER:DEST_PATH
		// $ docker cp [宿主机（外部主机）文件/文件夹路径] [容器名]:[容器文件/文件夹路径]
		// 2、从容器拷贝文件到宿主机
		// $ docker cp [OPTIONS] CONTAINER:SRC_PATH DEST_PATH|-
		// $ docker cp [容器ID/名称]:[容器文件/文件夹路径] [宿主机（外部主机）文件/文件夹路径]

		// 命令类型
		type CmdTyp int8
		const (
			CmdTypDef CmdTyp = iota
			CmdTypExc
			CmdTypCdExc
		)

		type CmdInfo struct {
			Cmd    string
			CmdTyp CmdTyp
		}

		cmdInfos := make([]CmdInfo, len(build)+3)
		index := 0
		// 删除容器资源路径
		cmdInfos[index] = CmdInfo{Cmd: fmt.Sprintf("rm -rf %s", containerResPath), CmdTyp: CmdTypExc}
		index++
		// 创建容器资源路径
		cmdInfos[index] = CmdInfo{Cmd: fmt.Sprintf("mkdir -p %s", containerResPath), CmdTyp: CmdTypExc}
		index++
		// 将源码cp到docker容器内
		// 将 test 目录下所有文件cp到容器 /tmp/test 目录下
		// $ docker cp test/. auto-deploy-build-env:/tmp/test/
		cmdInfos[index] = CmdInfo{Cmd: fmt.Sprintf("docker cp %s/. %s:%s/", resPath, container, containerResPath), CmdTyp: CmdTypDef}
		index++
		for _, cmd := range build {
			cmdInfos[index] = CmdInfo{Cmd: cmd, CmdTyp: CmdTypCdExc}
			index++
		}

		// 执行 build 命令集
		for _, cmdInfo := range cmdInfos {
			cmd := cmdInfo.Cmd
			cmdTyp := cmdInfo.CmdTyp
			switch cmdTyp {
			case CmdTypDef:

			case CmdTypExc:
				cmd = fmt.Sprintf("docker exec -u root %s /bin/bash -c \"%s\"", container, cmd)
			case CmdTypCdExc:
				// linux在宿主机执行docker容器环境内命令
				// $ sudo docker exec -u root auto-deploy-build-env /bin/bash -c "cd /root && ./test.sh"
				cmd = fmt.Sprintf("docker exec -u root %s /bin/bash -c \"cd %s && %s\"", container, containerResPath, cmd)
			}

			if sudo {
				cmd = "sudo " + cmd
			}
			log.Printf("%s\n", cmd)
			pCmd, err = util.Command(cmd)
			if err != nil {
				updETime(typ.StepBuild, recordId, err, nil)
				return err
			}

			buf, err := pCmd.CombinedOutput()
			if err != nil {
				updETime(typ.StepBuild, recordId, err, buf)
				return err
			}
		}

		// 将docker容器内编译结果cp到 auto-deploy-go 应用所在的服务器上
		// 将 /tmp/test 目录下所有文件cp到宿主机 test 目录下
		// $ docker cp auto-deploy-build-env:/tmp/test/. ./test/
		cmd := fmt.Sprintf("docker cp %s:%s/. %s/", container, containerResPath, resPath)
		if sudo {
			cmd = "sudo " + cmd
		}
		log.Printf("%s\n", cmd)
		pCmd, err = util.Command(cmd)
		if err != nil {
			updETime(typ.StepBuild, recordId, err, nil)
			return err
		}
		buf, err := pCmd.CombinedOutput()
		if err != nil {
			updETime(typ.StepBuild, recordId, err, buf)
			return err
		}
	}

	updETime(typ.StepBuild, recordId, nil, nil)
	return nil
}
