// item
// @author xiangqian
// @date 15:02 2022/12/20
package api

// 项目
type Item struct {
	Abs
	UserId   int64  // 项目所属用户id
	Name     string // 名称
	GitId    int64  // 项目所属Git id
	RepoUrl  string // Git仓库地址
	Branch   string // 分支名
	ServerId int64  // 项目所属Server id
	Cmd      string // 构建命令
	Script   string // 脚本，目前支持 #!/dockerfile, #!/static 解析
}
