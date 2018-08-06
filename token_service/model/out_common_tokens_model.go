package model

// 货币类型表
type OutCommonTokens struct {
	Id   uint32 `xorm:"id pk autoincr" json:"id"`
	Name string `xorm:"name" json:"name"`
	Mark string `xorm:"mark" json:"mark"`
}

func (*OutCommonTokens) TableName() string {
	return "g_common.tokens" // 跨库，g_common
}
