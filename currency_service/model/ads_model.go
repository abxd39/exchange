package model

import (
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
	. "digicon/proto/common"
)

// 买卖(广告)表
type Ads struct {
	Id          uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid         uint64 `xorm:"INT(10)" json:"uid"`              // 用户ID
	TypeId      uint32 `xorm:"TINYINT(1)" json:"type_id"`       // 类型:1出售 2购买
	TokenId     uint32 `xorm:"INT(10)" json:"token_id"`         // 货币类型
	TokenName   string `xorm:"VARBINARY(36)" json:"token_name"` // 货币名称
	Price       uint64 `xorm:"BIGINT(20)" json:"price"`         // 单价
	Num         uint64 `xorm:"BIGINT(20)" json:"num"`           // 数量
	Premium     int32  `xorm:"INT(10)" json:"premium"`          // 溢价
	AcceptPrice uint64 `xorm:"BIGINT(20)" json:"accept_price"`  // 可接受最低[高]单价
	MinLimit    uint32 `xorm:"INT(10)" json:"min_limit"`        // 最小限额
	MaxLimit    uint32 `xorm:"INT(10)" json:"max_limit"`        // 最大限额
	IsTwolevel  uint32 `xorm:"TINYINT(1)" json:"is_twolevel"`   // 是否要通过二级认证:0不通过 1通过
	Pays        string `xorm:"VARBINARY(50)" json:"pays"`       // 支付方式:以 , 分隔: 1,2,3
	Remarks     string `xorm:"VARBINARY(512)" json:"remarks"`   // 交易备注
	Reply       string `xorm:"VARBINARY(512)" json:"reply"`     // 自动回复问候语
	IsUsd       uint32 `xorm:"TINYINT(1)" json:"is_usd"`        // 是否美元支付:0否 1是
	States      uint32 `xorm:"TINYINT(1)" json:"states"`        // 状态:0下架 1上架
	Inventory   uint64 `xorm:"BIGINT(20)" json:"inventory"`     // 库存
	Fee         uint64 `xorm:"BIGINT(20)" json:"fee"`           // 手续费用
	CreatedTime string `xorm:"DATETIME" json:"created_time"`    // 创建时间
	UpdatedTime string `xorm:"DATETIME" json:"updated_time"`    // 修改时间
	IsDel       uint32 `xorm:"TINYINT(1)" json:"is_del"`        // 是否删除:0不删除 1删除
}

func (this *Ads) Get(id uint64) *Ads {

	data := new(Ads)
	isdata, err := dao.DB.GetMysqlConn().Id(id).Get(data)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}

	if !isdata {
		return nil
	}

	return data
}

func (this *Ads) Add() int {

	// 用户虚拟货币资产表
	isUcy := new(UserCurrency).Get(this.Uid, this.TokenId)
	if isUcy == nil || isUcy.Uid == 0 {
		return ERR_TOKEN_LESS
	}

	if this.Num > isUcy.Balance {
		return ERR_TOKEN_LESS
	}
	isUcy.Balance = isUcy.Balance - this.Num
	isUcy.Freeze = isUcy.Freeze + this.Num

	// 启用事务
	session := dao.DB.GetMysqlConn().NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	_, err = session.Where("uid=? AND token_id=?", this.Uid, this.TokenId).Update(isUcy)
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return ERRCODE_UNKNOWN
	}

	_, err = session.Insert(this)
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return ERRCODE_UNKNOWN
	}

	err = session.Commit()
	if err != nil {
		Log.Errorln(err.Error())
		session.Rollback()
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}

func (this *Ads) Update() int {

	isGet := this.Get(this.Id)
	if isGet == nil {
		return ERRCODE_ADS_NOTEXIST
	}

	_, err := dao.DB.GetMysqlConn().Id(this.Id).Update(this)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}

// 修改广告(买卖)状态
// id        uint64  广告ID
// status_id uint32  状态: 1下架 2上架 3正常(不删除) 4删除
func (this *Ads) UpdatedAdsStatus(id uint64, status_id uint32) int {

	var err error

	isGet := this.Get(id)
	if isGet == nil {
		return ERRCODE_ADS_NOTEXIST
	}

	if status_id == 1 || status_id == 2 {
		_, err = dao.DB.GetMysqlConn().Exec("UPDATE `ads` SET `states`=? WHERE `id`=?", status_id-1, id)
	} else if status_id == 3 || status_id == 4 {
		_, err = dao.DB.GetMysqlConn().Exec("UPDATE `ads` SET `is_del`=? WHERE `id`=?", status_id-3, id)
	}

	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}

// 法币交易列表 - (广告(买卖))
func (this *Ads) AdsList(TypeId, TokenId, Page, PageNum uint32) ([]AdsUserCurrencyCount, int64) {

	total, err := dao.DB.GetMysqlConn().Where("type_id=? AND token_id=?", TypeId, TokenId).Count(new(Ads))
	if err != nil {
		Log.Errorln(err.Error())
		return nil, 0
	}
	if total <= 0 {
		return nil, 0
	}

	limit := 0
	if Page > 0 {
		limit = int((Page - 1) * PageNum)
	}

	data := make([]AdsUserCurrencyCount, 0)
	err = dao.DB.GetMysqlConn().Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").Where("type_id=? AND token_id=?", TypeId, TokenId).Desc("updated_time").Limit(int(PageNum), limit).Find(&data)
	if err != nil {
		Log.Errorln(err.Error())
		return nil, 0
	}

	return data, total
}

// 个人法币交易列表 - (广告(买卖))
func (this *Ads) AdsUserList(Uid uint64, TypeId, Page, PageNum uint32) ([]Ads, int64) {

	total, err := dao.DB.GetMysqlConn().Where("uid=? AND type_id=?", Uid, TypeId).Count(new(Ads))
	if err != nil {
		Log.Errorln(err.Error())
		return nil, 0
	}
	if total <= 0 {
		return nil, 0
	}

	limit := 0
	if Page > 0 {
		limit = int((Page - 1) * PageNum)
	}

	data := make([]Ads, 0)
	err = dao.DB.GetMysqlConn().Where("uid=? AND type_id=?", Uid, TypeId).Desc("updated_time").Limit(int(PageNum), limit).Find(&data)
	if err != nil {
		Log.Errorln(err.Error())
		return nil, 0
	}

	return data, total
}
