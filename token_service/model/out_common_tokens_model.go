package model

import (
	"digicon/common/errors"
	"digicon/token_service/dao"
)

// 货币类型表
type OutCommonTokens struct {
	Id   uint32 `xorm:"id pk autoincr" json:"id"`
	Name string `xorm:"name" json:"name"`
	Mark string `xorm:"mark" json:"mark"`
}

func (*OutCommonTokens) TableName() string {
	return "g_common.tokens" // 跨库，g_common
}

// 获取币种信息
func (t *OutCommonTokens) Get(id uint32) (*OutCommonTokens, error) {
	commonTokens := new(OutCommonTokens)
	has, err := dao.DB.GetCommonMysqlConn().ID(id).Get(commonTokens)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	if !has {
		return nil, errors.NewNormal("币种不存在")
	}

	return commonTokens, nil
}
