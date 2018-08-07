package model

import (
	"digicon/common/errors"
	"digicon/common/model"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	"fmt"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

/*
const (
	MONEY_UKEY_TYPE_HASH     = 1
	MONEY_UKEY_TYPE_ENTRUST  = 2
	MONEY_UKEY_TYPE_REWARD   = 3
	MONEY_UKEY_TYPE_REGISTER = 4
	MONEY_UKEY_TYPE_TRADE    = 5
)
*/
type MoneyRecord struct {
	Id          uint64 `xorm:"pk autoincr BIGINT(20)"`
	Uid         uint64 `xorm:"comment('用户ID') unique(hash_index)  INT(11)"`
	TokenId     int    `xorm:"comment('代币ID') INT(11)"`
	Ukey        string `xorm:"comment('联合key') unique(hash_index) VARCHAR(128)"`
	Type        int    `xorm:"comment('流水类型1区块2委托') INT(11)"`
	Opt         int    `xorm:"comment('操作方向1加2减') unique(hash_index) TINYINT(4)"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)"`
	Balance     int64  `xorm:"comment('余额') BIGINT(20)"`
	CreatedTime int64  `xorm:"comment('操作时间')  created BIGINT(20)"`
}

type MoneyRecordWithToken struct {
	MoneyRecord `xorm:"extends"`
	TokenName   string
}

func (*MoneyRecord) TableName() string {
	return "money_record"
}

//流水列表
func (s *MoneyRecord) List(pageIndex, pageSize int, filter map[string]interface{}) (*model.ModelList, []*MoneyRecordWithToken, error) {
	query := DB.GetMysqlConn().Alias("mr").Join("LEFT", []string{new(OutCommonTokens).TableName(), "t"}, "mr.token_id=t.id").Where("1=1")

	//筛选
	if _, ok := filter["uid"]; ok {
		query.And("mr.uid=?", filter["uid"])
	}
	if _, ok := filter["transfer"]; ok { //划转流水
		query.And("mr.type IN (?,?)", proto.TOKEN_TYPE_OPERATOR_TRANSFER_TO_CURRENCY, proto.TOKEN_TYPE_OPERATOR_TRANSFER_FROM_CURRENCY)
	}
	if v, ok := filter["type"]; ok {
		query.And("mr.type=?", v)
	}

	//分页
	tmpQuery := *query
	total, err := tmpQuery.Count(s)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	offset, modelList := model.Paging(pageIndex, pageSize, int(total))

	var list []*MoneyRecordWithToken
	err = query.Select("mr.*,t.mark AS token_name").OrderBy("mr.id DESC").Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	modelList.Items = list

	return modelList, list, nil
}

//检查流水记录是否存在
func (s *MoneyRecord) CheckExist(ukey string, ty proto.TOKEN_TYPE_OPERATOR) (ok bool, err error) {
	ok, err = DB.GetMysqlConn().Where("ukey=? and type=?", ukey, ty).Exist(&MoneyRecord{})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

//新增一条流水
func InsertRecord(session *xorm.Session, p *MoneyRecord) (err error) {
	_, err = session.InsertOne(p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

// 检查币币划转到法币消息是否已处理，流水存在即表示已处理
func (s *MoneyRecord) IsTransferFromCurrencyHandled(transferId int64) (bool, *MoneyRecord, error) {
	has, err := DB.GetMysqlConn().Where(fmt.Sprintf("ukey='%d'", transferId)).And("type=11").Get(s)
	if err != nil {
		return false, nil, errors.NewSys(err)
	}

	return has, s, nil
}
