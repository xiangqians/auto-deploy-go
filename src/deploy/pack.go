// pack
// @author xiangqian
// @date 22:23 2022/12/31
package deploy

import (
	"auto-deploy-go/src/typ"
	"auto-deploy-go/src/util"
	"bufio"
	"errors"
	"fmt"
	"os"
)

func Pack(script typ.Script, recordId int64, resPath, packName string) error {
	updSTime(typ.StepPack, recordId)
	target := script.Target
	var names []string
	if target != nil && len(target) > 0 {
		for _, v := range target {
			path := fmt.Sprintf("%s/%s", resPath, v)
			if !util.IsExistOfPath(path) {
				err := errors.New(fmt.Sprintf("%s file does not exist", v))
				updETime(typ.StepPack, recordId, err, nil)
				return err
			}
			names = append(names, path)
		}
	}

	deployName := fmt.Sprintf("%s/%s", resPath, typ.DeployName)
	pDeployFile, err := os.Create(deployName)
	if err != nil {
		updETime(typ.StepPack, recordId, err, nil)
		return err
	}
	defer pDeployFile.Close()
	pWriter := bufio.NewWriter(pDeployFile)
	pWriter.WriteString(script.Deploy)
	pWriter.Flush()
	names = append(names, deployName)

	// zip
	err = util.Zip("", packName, names...)
	if err != nil {
		updETime(typ.StepPack, recordId, err, nil)
		return err
	}

	updETime(typ.StepPack, recordId, nil, nil)
	return nil
}
