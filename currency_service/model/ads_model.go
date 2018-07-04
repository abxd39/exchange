package model

import (
	"digicon/currency_service/dao"
	. "digicon/currency_service/log"
	. "digicon/proto/common"
	"fmt"
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
	//Premium     int32  `xorm:"INT(10)" json:"premium"`          // 溢价
	Premium     int64  `xorm:"BIGINT(64)" json:"premium"`          // 溢价
	AcceptPrice uint64 `xorm:"BIGINT(20)" json:"accept_price"`  // 可接受最低[高]单价
	MinLimit    uint32 `xorm:"INT(10)" json:"min_limit"`        // 最小限额
	MaxLimit    uint32 `xorm:"INT(10)" json:"max_limit"`        // 最大限额
	IsTwolevel  uint32 `xorm:"TINYINT(1)" json:"is_twolevel"`   // 是否要通过二级认证:0不通过 1通过
	Pays        string `xorm:"VARBINARY(50)" json:"pays"`       // 支付方式:以 , 分隔: 1,2,3
	Remarks     string `xorm:"VARBINARY(512)" json:"remarks"`   // 交易备注
	Reply       string `xorm:"VARBINARY(512)" json:"reply"`     // 自动回复问候语
	IsUsd       uint32 `xorm:"TINYINT(1)" json:"is_usd"`        // 是否美元支付:0否 1是
	States      uint32 `xorm:"TINYINT(1)" json:"states"`        // 状态:0下架 1上架
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

	data := new(Ads)
	_, err := dao.DB.GetMysqlConn().Where("uid=? AND token_id=? AND type_id=?", this.Uid, this.TokenId, this.TypeId).Get(data)
	if err != nil {
		Log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	//if isdata && data.Id > 0 {                  /// 去掉 去重
	//	return ERRCODE_ADS_EXISTS
	//}

	_, err = dao.DB.GetMysqlConn().Insert(this)
	if err != nil {
		fmt.Println("ad ads error!,", err.Error())
		Log.Errorln(err.Error())
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

	// 1下架
	if isGet.IsDel == 0 && isGet.States == status_id-1 {
		return ERRCODE_SUCCESS
	}
	// 2上架
	if isGet.IsDel == 0 && isGet.States == status_id-1 {
		return ERRCODE_SUCCESS
	}
	// 3正常(不删除)
	if isGet.IsDel == 0 && isGet.IsDel == status_id-3 {
		return ERRCODE_SUCCESS
	}
	// 4删除
	if isGet.IsDel == 1 && isGet.IsDel == status_id-3 {
		return ERRCODE_SUCCESS
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

	return ERRCODE_UNKNOWN
}

// 法币交易列表 - (广告(买卖))
//func (this *Ads) AdsList(TypeId, TokenId, Page, PageNum uint32) ([]AdsUserCurrencyCount, int64) {
func (this *Ads) AdsList(TypeId, TokenId, Page, PageNum uint32) ([]Ads, int64) {
	total, err := dao.DB.GetMysqlConn().Where("type_id=? AND token_id=?", TypeId, TokenId).Count(new(Ads))
	fmt.Println("total:", total)
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

	//data := make([]AdsUserCurrencyCount, int(PageNum))
	data := []Ads{}
	err = dao.DB.GetMysqlConn().
		//Join("INNER", "user_currency", "ads.uid=user_currency.uid AND ads.token_id=user_currency.token_id").
		//Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("ads.type_id=? AND ads.token_id=?", TypeId, TokenId).
		Desc("updated_time").
		Limit(int(PageNum), limit).
		Find(&data)

	if err != nil {
		Log.Errorln(err.Error())
		return nil, 0
	}
	//fmt.Println(data)
	return data, total
}

// 个人法币交易列表 - (广告(买卖))
func (this *Ads) AdsUserList(Uid uint64, TypeId, Page, PageNum uint32) ([]AdsUserCurrencyCount, int64) {

	total, err := dao.DB.GetMysqlConn().Where("uid=? AND type_id=?", Uid, TypeId).Count(new(Ads))

	fmt.Println("total:", total)
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

	//data := make([]AdsUserCurrencyCount, int(PageNum))
	data := []AdsUserCurrencyCount{}
	err = dao.DB.GetMysqlConn().
		Join("INNER", "user_currency", "ads.uid=user_currency.uid AND ads.token_id=user_currency.token_id").
		Where("ads.uid=? AND ads.type_id=?", Uid, TypeId).
		Desc("updated_time").
		Limit(int(PageNum), limit).
		Find(&data)

	if err != nil {
		Log.Errorln(err.Error())
		return nil, 0
	}

	//fmt.Println("User data:", data)
	return data, total
}
