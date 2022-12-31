// Type
// @author xiangqian
// @date 22:16 2022/12/31
package typ

// Tag
const (
	TagBuild  string = "[build]"
	TagTarget        = "[target]"
	TagDeploy        = "[deploy]"
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

const PackName string = "target.zip"
const DeployName string = "deploy.sh"

type Script struct {
	Build  []string // [build]
	Target []string // [target]
	Deploy string   // [deploy]
}
