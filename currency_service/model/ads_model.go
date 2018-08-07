package model

import (
	"digicon/currency_service/dao"
	//log "github.com/sirupsen/logrus"
	. "digicon/proto/common"
	"fmt"
	log "github.com/sirupsen/logrus"
	"digicon/common/convert"
)

// 买卖(广告)表
type Ads struct {
	Id        uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid       uint64 `xorm:"INT(10)" json:"uid"`              // 用户ID
	TypeId    uint32 `xorm:"TINYINT(1)" json:"type_id"`       // 类型:1出售 2购买
	TokenId   uint32 `xorm:"INT(10)" json:"token_id"`         // 货币类型
	TokenName string `xorm:"VARBINARY(36)" json:"token_name"` // 货币名称
	Price     uint64 `xorm:"BIGINT(20)" json:"price"`         // 单价
	Num       uint64 `xorm:"BIGINT(20)" json:"num"`           // 数量
	//Premium     int32  `xorm:"INT(10)" json:"premium"`          // 溢价
	Premium     int64  `xorm:"BIGINT(64)" json:"premium"`      // 溢价
	AcceptPrice uint64 `xorm:"BIGINT(20)" json:"accept_price"` // 可接受最低[高]单价
	MinLimit    uint32 `xorm:"INT(10)" json:"min_limit"`       // 最小限额
	MaxLimit    uint32 `xorm:"INT(10)" json:"max_limit"`       // 最大限额
	IsTwolevel  uint32 `xorm:"TINYINT(1)" json:"is_twolevel"`  // 是否要通过二级认证:0不通过 1通过
	Pays        string `xorm:"VARBINARY(50)" json:"pays"`      // 支付方式:以 , 分隔: 1,2,3
	Remarks     string `xorm:"VARBINARY(512)" json:"remarks"`  // 交易备注
	Reply       string `xorm:"VARBINARY(512)" json:"reply"`    // 自动回复问候语
	IsUsd       uint32 `xorm:"TINYINT(1)" json:"is_usd"`       // 是否美元支付:0否 1是
	States      uint32 `xorm:"TINYINT(1)" json:"states"`       // 状态:0下架 1上架
	CreatedTime string `xorm:"DATETIME" json:"created_time"`   // 创建时间
	UpdatedTime string `xorm:"DATETIME" json:"updated_time"`   // 修改时间
	IsDel       uint32 `xorm:"TINYINT(1)" json:"is_del"`       // 是否删除:0不删除 1删除
}

func (this *Ads) Get(id uint64) *Ads {

	data := new(Ads)
	isdata, err := dao.DB.GetMysqlConn().Id(id).Get(data)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}

	if !isdata {
		return nil
	}

	return data
}

func (this *Ads) Add() int {
	engine := dao.DB.GetMysqlConn()

	////////  先判断是否有余额  /////////
	uCurrency := new(UserCurrency)
	_, err := engine.Where("uid = ? AND token_id =?", this.Uid, this.TokenId).Get(uCurrency)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	var sumLimit int64
	sumLimit, err = this.GetUserAdsLimit(this.Uid, this.TokenId)
	if err != nil {
		log.Errorln("get user ads sum limit error!", err.Error())
	}

	if this.MinLimit < 100 {
		log.Errorln("限制最小价格要大于等于100")
		return ERRCODE_ADS_MIN_LIMIT
	}

	curCnyPrice := convert.Int64MulInt64By8Bit(int64(this.Num), int64(this.Price))
	minLimit := this.MinLimit * 100000000
	if curCnyPrice < int64(minLimit) {
		log.Errorln("当前广告的总价已小于最小价格的值")
		return ERRCODE_ADS_SET_PRICE
	}


	if this.TypeId == 2 && (uCurrency.Balance - sumLimit) < int64(this.Num)  { // type_id=2  是发布出售单
		log.Errorln("add ads error, user currency balance lower this num!")
		return ERR_TOKEN_LESS
	}

	data := new(Ads)
	_, err = engine.Where("uid=? AND token_id=? AND type_id=?", this.Uid, this.TokenId, this.TypeId).Get(data)
	if err != nil {
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}
	//if isdata && data.Id > 0 {                  /// 去掉 去重
	//	return ERRCODE_ADS_EXISTS
	//}

	_, err = engine.Insert(this)
	if err != nil {
		fmt.Println("ad ads error!,", err.Error())
		log.Errorln(err.Error())
		return ERRCODE_UNKNOWN
	}

	return ERRCODE_SUCCESS
}

func (this *Ads) Update() int {

	isGet := this.Get(this.Id)
	if isGet == nil {
		return ERRCODE_ADS_NOTEXIST
	}
	_, err := dao.DB.GetMysqlConn().Id(this.Id).Cols("price",
		"num","premium", "accept_price", "min_limit", "max_limit",
		"is_twolevel", "pays", "remarks", "reply", "updated_time").Update(this)
	if err != nil {
		fmt.Println(err)
		log.Errorln(err.Error())
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
	//  上架时候检查币的余额是否足够
	if status_id == 2 {
		uCurrency, err  := new(UserCurrency).GetBalance(isGet.Uid, uint32(isGet.TokenId))
		if err != nil {
			log.Errorln("查询用户余额错误")
			return ERRCODE_USER_BALANCE
		}
		userBanalce := uCurrency.Balance
		if uint64(userBanalce) < isGet.Num {
			log.Errorln("币的余额不足上架失败")
			return  ERR_TOKEN_LESS
		}
	}
	//// 1下架
	//fmt.Println("111111111111111111")
	//if isGet.IsDel == 0 && isGet.States == status_id-1 {
	//	return ERRCODE_SUCCESS
	//}
	//// 2上架
	//fmt.Println("222222222222")
	//if isGet.IsDel == 0 && isGet.States == status_id-1 {
	//	return ERRCODE_SUCCESS
	//}
	//// 3正常(不删除)
	//fmt.Println("33333333333")
	//if isGet.IsDel == 0 && isGet.IsDel == status_id-3 {
	//	return ERRCODE_SUCCESS
	//}
	//// 4删除
	//fmt.Println("4444444444")
	//if isGet.IsDel == 1 && isGet.IsDel == status_id-3 {
	//	return ERRCODE_SUCCESS
	//}
	// 1下架
	//fmt.Println("111111111111111111")
	//if isGet.IsDel == 0 && isGet.States == status_id {
	//	return ERRCODE_SUCCESS
	//}
	//
	//// 3正常(不删除)
	//fmt.Println("33333333333")
	//if isGet.IsDel == 0 && isGet.IsDel == status_id-3 {
	//	return ERRCODE_SUCCESS
	//}
	//// 4删除
	//fmt.Println("4444444444")
	//if isGet.IsDel == 1 && isGet.IsDel == status_id-3 {
	//	return ERRCODE_SUCCESS
	//}

	if status_id == 1 || status_id == 2 {
		_, err = dao.DB.GetMysqlConn().Exec("UPDATE `ads` SET `states`=?,`is_del`=0  WHERE `id`=?", status_id-1, id)
	} else if status_id == 3 || status_id == 4 {
		_, err = dao.DB.GetMysqlConn().Exec("UPDATE `ads` SET `is_del`=? WHERE `id`=?", status_id-3, id)
	}

	if err != nil {
		log.Errorln(err.Error())
		fmt.Println(err.Error())
		return ERRCODE_UNKNOWN
	} else {
		return ERRCODE_SUCCESS
	}

}

// 法币交易列表 - (广告(买卖))
func (this *Ads) AdsList(TypeId, TokenId, Page, PageNum uint32) ([]AdsUserCurrencyCountList, int64) {
	//func (this *Ads) AdsList(TypeId, TokenId, Page, PageNum uint32) ([]Ads, int64) {
	total, err := dao.DB.GetMysqlConn().Where("type_id=? AND token_id=? AND states = 1 AND is_del = 0", TypeId, TokenId).Count(new(Ads))
	fmt.Println("total:", total)
	if err != nil {
		log.Errorln(err.Error())
		return nil, 0
	}
	if total <= 0 {
		return nil, 0
	}

	limit := 0
	if Page > 0 {
		limit = int((Page - 1) * PageNum)
	}

	data := []AdsUserCurrencyCountList{}
	//data := []Ads{}
	//sql := "SELECT * FROM `ads` INNER JOIN user_currency ON ads.uid=user_currency.uid AND ads.token_id=user_currency.token_id
	// LEFT JOIN user_currency_count ON ads.uid=user_currency_count.uid WHERE (ads.type_id=2 AND ads.token_id=1) ORDER BY `updated_time` DESC LIMIT 9 ;
	err = dao.DB.GetMysqlConn().Table("ads").
		//Join("INNER", "user_currency", "ads.uid=user_currency.uid AND ads.token_id=user_currency.token_id").
		Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid").
		Where("ads.type_id=? AND ads.token_id=?", TypeId, TokenId).
		Where("ads.states = 1").
		Where("ads.is_del = 0 ").
		Desc("updated_time").
		Limit(int(PageNum), limit).
		Find(&data)

	if err != nil {
		log.Errorln(err.Error())
		return nil, 0
	}
	//fmt.Println("total:", total)
	return data, total
}

// 个人法币交易列表 - (广告(买卖))
func (this *Ads) AdsUserList(Uid uint64, TypeId, Page, PageNum uint32) ([]AdsUserCurrencyCount, int64) {

	total, err := dao.DB.GetMysqlConn().Where("uid=? AND type_id=?", Uid, TypeId).Count(new(Ads))

	fmt.Println("total:", total)
	if err != nil {
		log.Errorln(err.Error())
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
	fmt.Println("uid:", Uid, " typeid:", TypeId)
	data := []AdsUserCurrencyCount{}
	err = dao.DB.GetMysqlConn().
		//Join("INNER", "user_currency", "ads.uid=user_currency.uid AND ads.token_id=user_currency.token_id").
		//Join("INNER", "user_currency","ads.uid=user_currency.uid ").
		Where("ads.uid=? AND ads.type_id=?", Uid, TypeId).
		Where("ads.is_del != 1").
		Desc("updated_time").
		Limit(int(PageNum), limit).
		Find(&data)

	if err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
		return nil, 0
	}

	//fmt.Println("User data:", data)
	return data, total
}



/*
	获取所有这个币种出售广告的总和额度
    type_id = 2 为出售单
*/
func (this *Ads) GetUserAdsLimit(uid uint64, tokenId uint32)( sumLimit int64, err error){
	ssAds := new(Ads)
	engine := dao.DB.GetMysqlConn()
	sumLimit, err = engine.Where("uid=? AND token_id=?  AND type_id=2  AND states=1  AND is_del = 0", uid, tokenId).SumInt(ssAds, "num")
	return
}

/*
	获取在线广告最高价格和最低价格
*/

func (this *Ads) GetOnlineAdsMaxMinPrice(tokenId uint32) ( MaxPrice, MinPrice int64, err error){
	//sql := "select MAX(price)  AS max_price,MIN(price) AS min_price from ads where token_id =? and states=1"
	sql := " SELECT MAX(order.price)  AS max_price,MIN(order.price) AS min_price FROM  `order` LEFT JOIN  `ads`  ON `ads`.`id` = `order`.`ad_id` WHERE order.token_id =? AND   `ads`.`states`=1 AND `ads`.`is_del`=0;"
	engine := dao.DB.GetMysqlConn()
	type PriceStruct struct {
		MaxPrice  int64   `xorm:"BIGINT(20)"  json:"max_price"`
		MinPrice  int64   `xorm:"BIGINT(20)"  json:"min_price"`
	}
	price := PriceStruct{}
	_, err  = engine.SQL(sql, tokenId).Get(&price)
	MaxPrice = price.MaxPrice
	MinPrice = price.MinPrice
	return
}




/*
	广告下架
*/
func AdsAutoDownline(id  uint64) {
	ads := new(Ads).Get(id)
	if ads.TypeId == 1 {
		log.Println("这是买价格要购买币的订单")
		return
	}
	curCnyPrice := convert.Int64MulInt64By8Bit(int64(ads.Num), int64(ads.Price))
	minLimit := ads.MinLimit * 100000000
	if curCnyPrice < int64(minLimit) {
		log.Errorln("当前广告的总价已小于最小价格的值, 广告自动下架")
		sql := "update ads set states=0 where id=?"
		engine := dao.DB.GetMysqlConn()
		_, err := engine.Exec(sql, id)
		if err != nil {
			log.Println(err.Error())
		}
	}

}

/*
	检查用户广告额度不足余额时自动下架广告
*/
func AutoDownlineUserAds(uid uint64, tokenId uint64) {
	log.Printf("auto check  user balance  uid:%v, tokenid: %v ", uid, tokenId)
	uCurrency, err  := new(UserCurrency).GetBalance(uid, uint32(tokenId))
	if err != nil {
		log.Errorln("查询用户余额错误")
		return
	}
	userBanalce := uCurrency.Balance
	log.Println(userBanalce)

	adsList := []Ads{}
	engine := dao.DB.GetMysqlConn()
	err = engine.Where("uid = ? and token_id = ? and  states=1 and  type_id=2  and is_del=0", uid, tokenId).Desc("created_time").Find(&adsList)
	if err != nil {
		log.Errorln("划转币的时候查询该币的广告错误!")
		return
	}

	var curSumAdsNum uint64
	var otherAdsIdList []uint64

	for _, adsi := range adsList{
		curSumAdsNum  = curSumAdsNum + adsi.Num
		if int64(curSumAdsNum) >= userBanalce{
			otherAdsIdList = append(otherAdsIdList, adsi.Id)
		}
		curCnyPrice := convert.Int64MulInt64By8Bit(int64(adsi.Num), int64(adsi.Price))
		minLimit := adsi.MinLimit * 100000000
		if curCnyPrice < int64(minLimit) {
			otherAdsIdList = append(otherAdsIdList, adsi.Id)
		}
	}
	fmt.Println("need to down ads:", otherAdsIdList)
	//downlineSql := "update ads set states=0  where id in ?"
	//_, err = engine.Exec(downlineSql, otherAdsIdList)
	if len(otherAdsIdList) <= 0 {
		log.Println("not need to downline ads...")
		return
	}
	ad := Ads{States:0}
	_, err = engine.In("id", otherAdsIdList).Cols("states").Update(&ad)
	if err != nil {
		log.Errorln("自动下架广告错误", err.Error())
	}

}

