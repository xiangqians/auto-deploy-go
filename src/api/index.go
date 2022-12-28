// index
// @author xiangqian
// @date 21:03 2022/12/18
package api

import (
	"auto-deploy-go/src/com"
	"auto-deploy-go/src/db"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	netHttp "net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Stage 自动化部署阶段
type Stage int8

const (
	StagePull   Stage = iota + 1 // 拉取资源
	StageBuild                   // 构建
	StagePack                    // 打包
	StageUl                      // upload上传
	StageDeploy                  // 部署
)

const (
	StatusInDeploy      byte = iota + 1 // 部署中
	StatusDeployExc                     // 部署异常
	StatusDeploySuccess                 // 部署成功
)

// 互斥锁
var lock sync.Mutex

type ItemLastRecord struct {
	Id          int64
	ItemId      int64
	ItemName    string // item
	ItemRem     string
	PullStime   int64 // pull
	PullEtime   int64
	PullRem     string
	CommitId    string // commitId
	RevMsg      string // revMsg
	BuildStime  int64  // build
	BuildEtime  int64
	BuildRem    string
	PackStime   int64 // pack
	PackEtime   int64
	PackRem     string
	UlStime     int64 // ul
	UlEtime     int64
	UlRem       string
	DeployStime int64 // deploy
	DeployEtime int64
	DeployRem   string
	Status      byte   // status
	Rem         string // Rem
	AddTime     int64  // AddTime
}

func IndexPage(pContext *gin.Context) {
	session := sessions.Default(pContext)
	message := session.Get("message")
	session.Delete("message")
	session.Save()
	pContext.HTML(netHttp.StatusOK, "index.html", gin.H{
		"user":            GetUser(pContext),
		"itemLastRecords": ItemLastRecords(pContext, 0),
		"message":         message,
	})
}

func Deploy(pContext *gin.Context) {
	// redirect func
	redirect := func(message string) {
		session := sessions.Default(pContext)
		session.Set("message", message)
		session.Save()
		pContext.Redirect(netHttp.StatusMovedPermanently, "/")
	}

	// itemId
	itemIdStr := pContext.Param("itemId")
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		redirect(err.Error())
		return
	}

	// 加锁
	lock.Lock()
	// 解锁
	defer lock.Unlock()

	// itemLastRecords
	itemLastRecords := ItemLastRecords(pContext, itemId)
	if itemLastRecords == nil {
		redirect("项目不存在")
		return
	}

	// itemLastRecord
	itemLastRecord := itemLastRecords[0]
	if itemLastRecord.Status == StatusInDeploy {
		redirect("项目已在部署中")
		return
	}

	// item
	item := Item{}
	err = db.Qry(&item, "SELECT i.id, i.`name`, i.git_id, i.repo_url, i.branch, i.server_id, i.ini, i.rem FROM item i  WHERE i.del_flag = 0 AND i.id = ?", itemId)
	if err != nil {
		redirect(err.Error())
		return
	}

	// add record
	recordId, err := db.Add("INSERT INTO record(item_id, `status`, `add_time`) VALUES(?, ?, ?)", itemId, StatusInDeploy, time.Now().Unix())
	if err != nil {
		redirect(err.Error())
		return
	}

	go func() {
		// updRecord func
		updRecord := func(err error) {
			status := StatusDeploySuccess
			rem := ""
			if err != nil {
				status = StatusDeployExc
				rem = err.Error()
			}
			// update record
			db.Upd("UPDATE record SET `status` = ?, rem = ?, `upd_time` = ? where id = ?", status, rem, time.Now().Unix(), recordId)
		}

		// localRepoPath
		localRepoPath := fmt.Sprintf("%v/tmp/item%v", com.DataDir, item.Id)
		if com.IsExist(localRepoPath) {
			com.DelDir(localRepoPath)
		}

		// git clone
		err = gitClone(item, recordId, localRepoPath)
		if err != nil {
			updRecord(err)
			return
		}

		// ini
		ini := ParseIniText(item.Ini)

		// build
		err = build(ini, recordId, localRepoPath)
		if err != nil {
			updRecord(err)
			return
		}

		// deploy success
		updRecord(nil)
	}()

	redirect("")
}

func build(ini Ini, recordId int64, localRepoPath string) error {
	updSTime(StageBuild, recordId)

	_build := ini.Build
	if _build != nil && len(_build) > 0 {
		for _, cmd := range _build {
			cd, err := com.Cd(localRepoPath)
			if err != nil {
				updETime(StageBuild, recordId, err)
				return err
			}

			cmd = fmt.Sprintf("%s && %s", cd, cmd)
			pCmd, err := com.Command(cmd)
			if err != nil {
				updETime(StageBuild, recordId, err)
				return err
			}

			_, err = pCmd.CombinedOutput()
			if err != nil {
				updETime(StageBuild, recordId, err)
				return err
			}
		}
	}

	updETime(StageBuild, recordId, nil)
	return nil
}

func gitClone(item Item, recordId int64, localRepoPath string) error {
	updSTime(StagePull, recordId)

	_git := Git{}
	err := db.Qry(&_git, "SELECT g.id, g.`user`, g.passwd FROM git g WHERE g.del_flag = 0 AND g.id = ?", item.GitId)
	if err != nil {
		updETime(StagePull, recordId, err)
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
	pRepository, err := git.PlainClone(localRepoPath, isBare, &git.CloneOptions{
		URL:           item.RepoUrl,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", item.Branch)), // branchName, tagName, commitId
		Progress:      os.Stdout,
		Auth:          auth,
	})
	if err != nil {
		updETime(StagePull, recordId, err)
		return err
	}

	// 获取 HEAD 指向的分支
	// ... retrieves the branch pointed by HEAD
	pReference, err := pRepository.Head()
	if err != nil {
		updETime(StagePull, recordId, err)
		return err
	}

	// ... retrieves the commit history
	pCommitIter, err := pRepository.Log(&git.LogOptions{From: pReference.Hash()})
	if err != nil {
		updETime(StagePull, recordId, err)
		return err
	}

	// 最近一次提交信息
	pCommit, err := pCommitIter.Next()
	if err != nil {
		updETime(StagePull, recordId, err)
		return err
	}

	db.Upd("UPDATE record SET commit_id = ?, rev_msg = ? WHERE id = ?", pCommit.ID().String(), pCommit.String(), recordId)

	updETime(StagePull, recordId, nil)
	return nil
}

func updETime(stage Stage, recordId int64, err error) {
	etimeName := ""
	remName := ""
	switch stage {
	case StagePull:
		etimeName = "pull_etime"
		remName = "pull_rem"

	case StageBuild:
		etimeName = "build_etime"
		remName = "build_rem"

	case StagePack:
		etimeName = "pack_etime"
		remName = "pack_rem"

	case StageUl:
		etimeName = "ul_etime"
		remName = "ul_rem"

	case StageDeploy:
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

func updSTime(stage Stage, recordId int64) {
	stimeName := ""
	switch stage {
	case StagePull:
		stimeName = "pull_stime"

	case StageBuild:
		stimeName = "build_stime"

	case StagePack:
		stimeName = "pack_stime"

	case StageUl:
		stimeName = "ul_stime"

	case StageDeploy:
		stimeName = "deploy_stime"
	}

	db.Upd(fmt.Sprintf("UPDATE record SET %s = ? WHERE id = ?", stimeName), time.Now().Unix(), recordId)
}

func ItemLastRecords(pContext *gin.Context, itemId int64) []ItemLastRecord {
	itemLastRecords := make([]ItemLastRecord, 1)
	user := GetUser(pContext)
	sql := "SELECT IFNULL(r.id, 0) AS 'id', i.id AS 'item_id', i.`name` AS 'item_name', i.rem AS 'item_rem', " +
		"IFNULL(r.pull_stime, 0) AS 'pull_stime', IFNULL(r.pull_etime, 0) AS 'pull_etime', IFNULL(r.pull_rem, '') AS 'pull_rem', " +
		"IFNULL(r.commit_id, '') AS 'commit_id', IFNULL(r.rev_msg, '') AS 'rev_msg', " +
		"IFNULL(r.build_stime, 0) AS 'build_stime', IFNULL(r.build_etime, 0) AS 'build_etime', IFNULL(r.build_rem, '') AS 'build_rem', " +
		"IFNULL(r.pack_stime, 0) AS 'pack_stime', IFNULL(r.pack_etime, 0) AS 'pack_etime', IFNULL(r.pack_rem, '') AS 'pack_rem', " +
		"IFNULL(r.ul_stime, 0) AS 'ul_stime', IFNULL(r.ul_etime, 0) AS 'ul_etime', IFNULL(r.ul_rem, '') AS 'ul_rem', " +
		"IFNULL(r.deploy_stime, 0) AS 'deploy_stime', IFNULL(r.deploy_etime, 0) AS 'deploy_etime', IFNULL(r.deploy_rem, '') AS 'deploy_rem', " +
		"IFNULL(r.status, 0) AS 'status', " +
		"IFNULL(r.rem, '') AS 'rem', " +
		"IFNULL(r.add_time, 0) AS 'add_time' " +
		"FROM item i " +
		"LEFT JOIN record r ON r.del_flag = 0 AND r.item_id = i.id " +
		"LEFT JOIN record rt ON rt.del_flag = 0 AND rt.item_id = r.item_id AND r.add_time < rt.add_time "

	sql += "WHERE i.del_flag = 0 AND i.user_id IN(SELECT DISTINCT(owner_id) FROM rx WHERE del_flag = 0 AND sharer_id = ? UNION ALL SELECT %v) "
	if itemId > 0 {
		sql += fmt.Sprintf("AND i.id = %v ", strconv.FormatInt(itemId, 10))
	}
	sql += "GROUP BY i.id, r.id HAVING COUNT(rt.id) < 1"
	sql = fmt.Sprintf(sql, user.Id)
	err := db.Qry(&itemLastRecords, sql, user.Id)
	if err != nil {
		log.Println(err)
		return nil
	}

	if itemLastRecords[0].ItemId == 0 {
		itemLastRecords = nil
	}

	return itemLastRecords
}
