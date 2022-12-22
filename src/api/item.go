// item
// @author xiangqian
// @date 15:02 2022/12/20
package api

// 项目
type Item struct {
	Name           string `json:"name"`           // 项目名称
	Rem            string `json:"Rem"`            // 项目描述
	LastDeployTime int64  `json:"lastDeployTime"` // 最近一次部署时间戳
	LastRevMessage string `json:"lastRevMessage"` // 最近一次修改信息
}

// 项目阶段
type Stage struct {
	StartTime int64  `json:"startTime"` // 项目阶段开始时间戳
	EndTime   int64  `json:"endTime"`   // 项目阶段结束时间戳
	Message   string `json:"message"`   // 项目阶段消息
}
