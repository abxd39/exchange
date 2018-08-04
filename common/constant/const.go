package constant

const (
	ZOOKEEPER int = 0
	ETCD      int = 1
)

const (
	REGISTER int = 1
)

const (
	AUTH_NIL    = -1 //取消认证
	AUTH_EMAIL  = 2  //00000010 //邮箱
	AUTH_PHONE  = 1  //00000001 //电话
	AUTH_GOOGLE = 8  //00001000 //google
	AUTH_TWO    = 4  //0100 //二级
	AUTH_FIRST  = 16 //0001 0000 实名认证
)

//是否设置资金密码状态标识
const (
	AUTH_TRADEMARK = 1 //0001资金密码设置状态
	APPLY_FOR_FIRST=2//实名认证申请状态
	APPLY_FOR_SECOND=4//二级认证申请状态
)

// 划转
const (
	RDS_TOKEN_TO_CURRENCY_TODO = "token_to_currency_todo"
	RDS_TOKEN_TO_CURRENCY_DONE = "token_to_currency_done"
	RDS_CURRENCY_TO_TOKEN_TODO = "currency_to_token_todo"
	RDS_CURRENCY_TO_TOKEN_DONE = "currency_to_token_done"
)
