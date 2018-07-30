package model

import (
	"digicon/currency_service/dao"
	log "github.com/sirupsen/logrus"
)

// 货币类型表
type CommonTokens struct {
	Id   uint32 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Name string `xorm:"VARBINARY(20)" json:"cn_name"` // 货币中文名
	Mark string `xorm:"VARBINARY(20)" json:"name"`    // 货币标识
}

// common 下获取货币类型
func (this *CommonTokens) Get(id uint32, mark string) *CommonTokens {

	data := new(CommonTokens)
	var isdata bool
	var err error
	if id > 0 {
		isdata, err = dao.DB.GetCommonMysqlConn().Table("tokens").Where("id=?", id).Get(data)
	} else {
		isdata, err = dao.DB.GetCommonMysqlConn().Table("tokens").Where("mark=?", mark).Get(data)
	}

	if err != nil {
		log.Errorln(err.Error())
		return nil
	}

	if isdata == false {
		return nil
	}

	return data
}

// 获取货币类型列表
func (this *CommonTokens) List() []CommonTokens {
	data := make([]CommonTokens, 0)
	//sql := "select id, name, mark from g_common.tokens"
	//err := dao.DB.GetCommonMysqlConn().SQL(sql).Find(&data)
	err := dao.DB.GetCommonMysqlConn().Table("tokens").Find(&data)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}

	return data
}

//* 根据token_id列表获取货币
func (this *CommonTokens) GetByTokenIds(ids []int) (data []CommonTokens){
	idsLen := len(ids)
	pageNum := 20
	if idsLen > pageNum{
		var page int
		if idsLen % pageNum == 0{
			page = idsLen / pageNum
		}else{
			page = idsLen / pageNum + 1
		}
		for i:= 0;i< page ; i++{
			var tmpdata []CommonTokens
			tmpids := ids[i * pageNum: (i+1)*pageNum]
			err := dao.DB.GetCommonMysqlConn().Table("tokens").In("id", tmpids).Find(&tmpdata)
			if err != nil {
				log.Error(err.Error())
			}else{
				data = append(data, tmpdata...)
			}
		}
	}else{
		err := dao.DB.GetCommonMysqlConn().Table("tokens").In("id", ids).Find(&data)
		if err != nil {
			log.Errorln(err.Error())
		}
	}
	return
}
