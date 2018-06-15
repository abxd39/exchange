package model

import (
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
	. "digicon/proto/common"
)

// 买卖(广告)表
type Ads struct {
	Id          uint64  `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid         uint64  `xorm:"INT(10)" json:"uid"`                // 用户ID
	TypeId      uint32  `xorm:"TINYINT(1)" json:"type_id"`         // 类型:1出售 2购买
	TokenId     uint32  `xorm:"INT(10)" json:"token_id"`           // 货币类型
	TokenName   string  `xorm:"VARBINARY(36)" json:"token_name"`   // 货币名称
	Price       float64 `xorm:"DECIMAL(20,6)" json:"price"`        // 单价
	Num         float64 `xorm:"DECIMAL(20,6)" json:"num"`          // 数量
	Premium     int32   `xorm:"INT(10)" json:"premium"`            // 溢价
	AcceptPrice float64 `xorm:"DECIMAL(20,6)" json:"accept_price"` // 可接受最低[高]单价
	MinLimit    uint32  `xorm:"INT(10)" json:"min_limit"`          // 最小限额
	MaxLimit    uint32  `xorm:"INT(10)" json:"max_limit"`          // 最大限额
	IsTwolevel  uint32  `xorm:"TINYINT(1)" json:"is_twolevel"`     // 是否要通过二级认证:0不通过 1通过
	Pays        string  `xorm:"VARBINARY(50)" json:"pays"`         // 支付方式:以 , 分隔: 1,2,3
	Remarks     string  `xorm:"VARBINARY(512)" json:"remarks"`     // 交易备注
	Reply       string  `xorm:"VARBINARY(512)" json:"reply"`       // 自动回复问候语
	IsUsd       uint32  `xorm:"TINYINT(1)" json:"is_usd"`          // 是否美元支付:0否 1是
	States      uint32  `xorm:"TINYINT(1)" json:"states"`          // 状态:0下架 1上架
	CreatedTime string  `xorm:"DATETIME" json:"created_time"`
	UpdatedTime string  `xorm:"DATETIME" json:"updated_time"`
}

func (this *Ads) Add() int {
	_, err := dao.DB.GetMysqlConn().Insert(this)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}

func (this *Ads) Update() int {
	_, err := dao.DB.GetMysqlConn().Id(this.Id).Update(this)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}
