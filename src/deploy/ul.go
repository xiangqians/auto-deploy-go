// upload
// @author xiangqian
// @date 22:31 2022/12/31
package deploy

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

func UlAndDeploy(item typ.Item, recordId int64, packName, ulPath string) error {
	updSTime(typ.StepUl, recordId)

	server := typ.Server{}
	err := db.Qry(&server, "SELECT s.id, s.`host`, s.`port`, s.`user`, s.passwd FROM server s WHERE s.del_flag = 0 AND s.id = ?", item.ServerId)
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}

	if server.Id == 0 {
		err = errors.New("server does not exist")
		updETime(typ.StepUl, recordId, err)
		return err
	}

	// 建立 ssh client
	config := &ssh.ClientConfig{
		User:            server.User,
		Auth:            []ssh.AuthMethod{ssh.Password(server.Passwd)},
		Timeout:         5 * time.Minute,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", server.Host, server.Port)
	pSshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}
	defer pSshClient.Close()

	exec := func(cmd string) error {
		// 开启一个 session，用于执行一个命令
		session, eerr := pSshClient.NewSession()
		if eerr != nil {
			return eerr
		}
		defer session.Close()
		_, eerr = session.CombinedOutput(cmd)
		return eerr
	}

	// 删除上传路径
	err = exec(fmt.Sprintf("rm -rf %s", ulPath))
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}

	// 创建上传路径
	err = exec(fmt.Sprintf("mkdir -p %s", ulPath))
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}

	// 基于ssh client, 创建 sftp 客户端
	pSftpClient, err := sftp.NewClient(pSshClient)
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}
	defer pSftpClient.Close()

	// 上传文件
	pSrcFile, err := os.Open(packName)
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}
	defer pSrcFile.Close()
	ulName := fmt.Sprintf("%s/%s", ulPath, typ.PackName)
	pDstFile, err := pSftpClient.Create(ulName)
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}
	defer pDstFile.Close()
	buf := make([]byte, 100*1024*1024) // 100 MB
	_, err = io.CopyBuffer(pDstFile, pSrcFile, buf)
	//_, err = io.CopyN(pDstFile, pSrcFile, 100*1024*1024) // 100 MB -> EOF ?
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}

	// 解压
	err = exec(fmt.Sprintf("unzip -o %s -d %s", ulName, ulPath))
	if err != nil {
		updETime(typ.StepUl, recordId, err)
		return err
	}

	updETime(typ.StepUl, recordId, nil)

	// #################### Deploy  ####################

	updSTime(typ.StepDeploy, recordId)
	cmd := fmt.Sprintf("cd %s && chmod +x %s && ./%s", ulPath, typ.DeployName, typ.DeployName)
	err = exec(cmd)
	if err != nil {
		updETime(typ.StepDeploy, recordId, err)
		return err
	}
	updETime(typ.StepDeploy, recordId, nil)
	return nil
}
