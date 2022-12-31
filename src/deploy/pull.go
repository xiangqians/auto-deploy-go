// pull
// @author xiangqian
// @date 22:12 2022/12/31
package deploy

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

func Pull(item typ.Item, recordId int64, resPath string) error {
	updSTime(typ.StepPull, recordId)

	_git := typ.Git{}
	err := db.Qry(&_git, "SELECT g.id, g.`user`, g.passwd FROM git g WHERE g.del_flag = 0 AND g.id = ?", item.GitId)
	if err != nil {
		updETime(typ.StepPull, recordId, err)
		return err
	}

	var auth transport.AuthMethod = nil
	if _git.Id != 0 {
		auth = &http.BasicAuth{
			Username: _git.User,
			Password: _git.Passwd,
		}
	}

	// Clones the repository into the given dir, just as a normal git clone does
	isBare := false
	pRepository, err := git.PlainClone(resPath, isBare, &git.CloneOptions{
		URL:           item.RepoUrl,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", item.Branch)), // branchName, tagName, commitId
		Progress:      os.Stdout,
		Auth:          auth,
	})
	if err != nil {
		updETime(typ.StepPull, recordId, err)
		return err
	}

	// 获取 HEAD 指向的分支
	// ... retrieves the branch pointed by HEAD
	pReference, err := pRepository.Head()
	if err != nil {
		updETime(typ.StepPull, recordId, err)
		return err
	}

	// ... retrieves the commit history
	pCommitIter, err := pRepository.Log(&git.LogOptions{From: pReference.Hash()})
	if err != nil {
		updETime(typ.StepPull, recordId, err)
		return err
	}

	// 最近一次提交信息
	pCommit, err := pCommitIter.Next()
	if err != nil {
		updETime(typ.StepPull, recordId, err)
		return err
	}

	_, err = db.Upd("UPDATE record SET commit_id = ?, rev_msg = ? WHERE id = ?", pCommit.ID().String(), pCommit.String(), recordId)
	if err != nil {
		updETime(typ.StepPull, recordId, err)
		return err
	}

	updETime(typ.StepPull, recordId, nil)
	return nil
}
