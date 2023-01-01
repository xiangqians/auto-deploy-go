// api
// @author xiangqian
// @date 22:33 2022/12/31
package typ

import "encoding/gob"

// 抽象实体定义
type Abs struct {
	Id      int64  `form:"id" binding:"gte=0"`    // 主键id
	Rem     string `form:"rem" binding:"max=200"` // 备注
	DelFlag byte   `form:"delFlag"`               // 删除标识，0-正常，1-删除
	AddTime int64  `form:"addTime"`               // 创建时间（时间戳，s）
	UpdTime int64  `form:"updTime"`               // 修改时间（时间戳，s）
}

// Git
type Git struct {
	Abs
	UserId int64  // Git所属用户id
	Name   string `form:"name" binding:"required,min=1,max=60"`    // 名称
	User   string `form:"user" binding:"required,min=1,max=60"`    // 用户
	Passwd string `form:"passwd" binding:"required,min=1,max=100"` // 密码
}

// Item 项目
type Item struct {
	Abs
	UserId     int64  `form:"userId"`                                              // 项目所属用户id
	Name       string `form:"name" binding:"required,excludes= ,min=1,max=60"`     // 名称
	GitId      int64  `form:"gitId" binding:"gte=0"`                               // 项目所属Git id
	GitName    string `form:"gitName"`                                             // 项目所属Git Name
	RepoUrl    string `form:"repoUrl" binding:"required,excludes= ,min=1,max=500"` // Git仓库地址
	Branch     string `form:"branch" binding:"required,excludes= ,min=1,max=60"`   // 分支名
	ServerId   int64  `form:"serverId" binding:"required,gt=0"`                    // 项目所属Server id
	ServerName string `form:"serverName"`                                          // 项目所属Server Name
	Script     string `form:"script" binding:"required,min=1,max=100000"`          // 脚本
}

type Rx struct {
	Abs
	Name       string `form:"name" binding:"required,min=1,max=60"` // 名称
	OwnerId    int64  // 拥有者id
	OwnerName  string // 拥有者名称
	SharerId   int64  // 共享者id
	SharerName string // 共享者名称
}

type Server struct {
	Abs
	UserId int64  // Server所属用户id
	Name   string `form:"name" binding:"required,min=1,max=60"`    // 名称
	Host   string `form:"host" binding:"required,min=1,max=60"`    // host
	Port   int    `form:"port" binding:"required,gt=0"`            // 端口
	User   string `form:"user" binding:"required,min=1,max=60"`    // 用户
	Passwd string `form:"passwd" binding:"required,min=1,max=100"` // 密码
}

type User struct {
	Abs
	Name     string `form:"name" binding:"required,excludes= ,min=1,max=60"`               // 用户名
	Nickname string `form:"nickname"binding:"max=60"`                                      // 昵称
	Passwd   string `form:"passwd" binding:"required,excludes= ,max=100"`                  // 密码
	RePasswd string `form:"rePasswd" binding:"required,excludes= ,max=100,eqfield=Passwd"` // retype Passwd
}

type ItemLastRecord struct {
	Id           int64
	ItemId       int64
	ItemName     string // item
	ItemRem      string
	PullStime    int64 // pull
	PullEtime    int64
	PullStatus   byte
	PullRem      string
	CommitId     string // commitId
	RevMsg       string // revMsg
	BuildStime   int64  // build
	BuildEtime   int64
	BuildStatus  byte
	BuildRem     string
	PackStime    int64 // pack
	PackEtime    int64
	PackStatus   byte
	PackRem      string
	UlStime      int64 // ul
	UlEtime      int64
	UlStatus     byte
	UlRem        string
	DeployStime  int64 // deploy
	DeployEtime  int64
	DeployStatus byte
	DeployRem    string
	Status       byte   // status
	Rem          string // Rem
	AddTime      int64  // AddTime
}

const (
	StatusInDeploy      byte = iota + 1 // 部署中
	StatusDeployExc                     // 部署异常
	StatusDeploySuccess                 // 部署成功
)

const (
	LocaleZh = "zh"
	LocaleEn = "en"
)

// 注册模型
func init() {
	gob.Register(Git{})
	gob.Register(Item{})
	gob.Register(Rx{})
	gob.Register(Server{})
	gob.Register(User{})
	gob.Register(ItemLastRecord{})
}
