// Script
// @author xiangqian
// @date 21:53 2022/12/31
package deploy

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

func ParseScriptTxt(scriptTxt string) typ.Script {
	script := typ.Script{}
	pReader := bufio.NewReader(bytes.NewBufferString(scriptTxt))

	set := func(slice []string, ty string) {
		switch ty {
		case typ.TagBuild:
			script.Build = slice

		case typ.TagTarget:
			script.Target = slice

		case typ.TagDeploy:
			str := ""
			for _, v := range slice {
				str += v + "\n"
			}
			script.Deploy = str

		default:
		}
	}

	ty := ""
	var slice []string
	handleLine := func(line string) {
		if ty != typ.TagDeploy {
			if line == "" || strings.HasPrefix(line, "#") {
				return
			}
		}

		switch line {
		case typ.TagBuild:
			set(slice, ty)
			slice = nil
			ty = line

		case typ.TagTarget:
			set(slice, ty)
			slice = nil
			ty = line

		case typ.TagDeploy:
			set(slice, ty)
			slice = nil
			ty = line

		default:
			slice = append(slice, line)
		}
	}
	for {
		line, err := pReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				handleLine(line)
				break
			}
			log.Println(err)
			continue
		}

		handleLine(line)
	}
	set(slice, ty)

	return script
}

func updETime(Step typ.Step, recordId int64, err error) {
	etimeName := ""
	remName := ""
	switch Step {
	case typ.StepPull:
		etimeName = "pull_etime"
		remName = "pull_rem"

	case typ.StepBuild:
		etimeName = "build_etime"
		remName = "build_rem"

	case typ.StepPack:
		etimeName = "pack_etime"
		remName = "pack_rem"

	case typ.StepUl:
		etimeName = "ul_etime"
		remName = "ul_rem"

	case typ.StepDeploy:
		etimeName = "deploy_etime"
		remName = "deploy_rem"
	}

	var etime int64 = time.Now().Unix()
	rem := ""
	if err != nil {
		etime = -1
		rem = err.Error()
	}
	db.Upd(fmt.Sprintf("UPDATE record SET %s = ?, %s = ? WHERE id = ?", etimeName, remName), etime, rem, recordId)
}

func updSTime(Step typ.Step, recordId int64) {
	stimeName := ""
	switch Step {
	case typ.StepPull:
		stimeName = "pull_stime"

	case typ.StepBuild:
		stimeName = "build_stime"

	case typ.StepPack:
		stimeName = "pack_stime"

	case typ.StepUl:
		stimeName = "ul_stime"

	case typ.StepDeploy:
		stimeName = "deploy_stime"
	}

	db.Upd(fmt.Sprintf("UPDATE record SET %s = ? WHERE id = ?", stimeName), time.Now().Unix(), recordId)
}
