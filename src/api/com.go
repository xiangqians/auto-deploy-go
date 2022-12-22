// common
// @author xiangqian
// @date 13:46 2022/12/22
package api

// 抽象实体定义
type Abs struct {
	Id         int64  // 主键id
	Rem        string // 备注
	DelFlag    byte   // 删除标识，0-正常，1-删除
	CreateTime int64  // 创建时间（时间戳，s）
	UpdateTime int64  // 修改时间（时间戳，s）
}
