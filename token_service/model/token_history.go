package model

type TokenHistory struct {
	Id         int    `xorm:"comment('操作序号') INT(255)"`
	Uid        int    `xorm:"comment('用户id') INT(11)"`
	TokenId    int    `xorm:"comment(' 货币类型') INT(11)"`
	Num        string `xorm:"comment('提现数量') DECIMAL(20,4)"`
	Fee        string `xorm:"comment('手续费') DECIMAL(20,4)"`
	Address    string `xorm:"comment('提现地址') VARCHAR(255)"`
	RecordTime int    `xorm:"comment('提交时间') INT(11)"`
	CheckTime  int    `xorm:"comment('审核时间') INT(11)"`
	AdminId    int    `xorm:"comment('审核人') INT(11)"`
	Status     int    `xorm:"comment('状态0审核中，1拒绝，2成功') INT(11)"`
	Operator   int    `xorm:"comment('操作类型0充币，1提币') INT(11)"`
}
