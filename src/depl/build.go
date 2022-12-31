// build
// @author xiangqian
// @date 21:42 2022/12/31
package depl

import (
	"auto-deploy-go/src/typ"
	"auto-deploy-go/src/util"
	"fmt"
)

func DefaultBuildEnv() {

}

func DockerBuildEnv() {

}

func Build(script typ.Script, recordId int64, resPath string) error {
	updSTime(typ.StageBuild, recordId)

	_build := script.Build
	if _build != nil && len(_build) > 0 {
		for _, cmd := range _build {
			cd, err := util.Cd(resPath)
			if err != nil {
				updETime(typ.StageBuild, recordId, err)
				return err
			}

			cmd = fmt.Sprintf("%s && %s", cd, cmd)
			pCmd, err := util.Command(cmd)
			if err != nil {
				updETime(typ.StageBuild, recordId, err)
				return err
			}

			_, err = pCmd.CombinedOutput()
			if err != nil {
				updETime(typ.StageBuild, recordId, err)
				return err
			}
		}
	}

	updETime(typ.StageBuild, recordId, nil)
	return nil
}
