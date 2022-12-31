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

// Step 自动化部署步骤
type Step int8

const (
	StepPull   Step = iota + 1 // 拉取资源
	StepBuild                  // 构建
	StepPack                   // 打包
	StepUl                     // upload上传
	StepDeploy                 // 部署
)

// Env type
const (
	EnvDefault string = "default"
	EnvDocker         = "docker:"
)

const PackName string = "target.zip"
const DeployName string = "deploy.sh"

type Script struct {
	Build  []string // [build]
	Target []string // [target]
	Deploy string   // [deploy]
}
