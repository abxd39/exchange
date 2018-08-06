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
	AUTH_TRADEMARK               = 1  //0001资金密码设置状态
	APPLY_FOR_FIRST              = 2  //实名认证申请状态
	APPLY_FOR_SECOND             = 4  //二级认证申请状态
	APPLY_FOR_SECOND_NOT_ALREADY = 8  //二级认证没有通过
	APPLY_FOR_FIRST_NOT_ALREADY  = 16 //一级认证状态未通过
)

//实名认证
const (
	FIRST_NOT_V       = 1 //没认证
	FIRST_ALREADY     = 2 //已经通过认证
	FIRST_NOT_ALREADY = 3 //没有通过认证
	FIRST_VERIFYING   = 4 //认证中
)

const (
	SECOND_NOT_V       = 1 //没认证
	SECOND_ALREADY     = 2 //已经通过认证
	SECOND_NOT_ALREADY = 3 //没有通过认证
	SECOND_VERIFYING   = 4 //认证中
)

// 划转
const (
	RDS_TOKEN_TO_CURRENCY_TODO = "token_to_currency_todo"
	RDS_TOKEN_TO_CURRENCY_DONE = "token_to_currency_done"
	RDS_CURRENCY_TO_TOKEN_TODO = "currency_to_token_todo"
	RDS_CURRENCY_TO_TOKEN_DONE = "currency_to_token_done"
)

// 验证类型   1 注册 2 忘记密码 3 修改手机号码 4重置谷歌验证码 5 重置资金密码 6 修改登录密码 7 设置银行卡支付 8 设置微信支付 9 设置支付宝支付 10 设置PayPal支付
//           11 绑定手机  12 绑定邮箱, 13 提币



//const (
//	SMS_REGISTER         = 1 //注册业务
//	SMS_FORGET           = 2
//	SMS_MODIFY_PHONE     = 3
//	SMS_SET_GOOGLE_CODE  = 4
//	SMS_RESET_TRADE_PWD  = 5
//	SMS_MODIFY_LOGIN_PWD = 6
//	SMS_BANK_PAY         = 7
//	SMS_WECHAT_PAY       = 8
//	SMS_AIL_PAY          = 9
//	SMS_PAYPAL_PAY       = 10
//	SMS_BIND_PHONE       = 11
//	SMS_BIND_EMAIL       = 12
//	SMS_CARRY_COIN       = 13
//	SMS_REAL_NAME        = 14
//	SMS_MAX              = 15
//)


//// 验证类型   1 注册 2 忘记密码 3 修改手机号码 4重置谷歌验证码 5 重置资金密码 6 修改登录密码 7 设置银行卡支付 8 设置微信支付 9 设置支付宝支付 10 设置PayPal支付
////           11 绑定手机  12 绑定邮箱, 13 提币
const (
	SMS_REGISTER         = 1 //注册业务
	SMS_FORGET           = 2
	SMS_MODIFY_PHONE     = 3
	SMS_SET_GOOGLE_CODE  = 4
	SMS_RESET_TRADE_PWD  = 5
	SMS_MODIFY_LOGIN_PWD = 6
	SMS_BANK_PAY         = 7
	SMS_WECHAT_PAY       = 8
	SMS_AIL_PAY          = 9
	SMS_PAYPAL_PAY       = 10
	SMS_BIND_PHONE       = 11
	SMS_BIND_EMAIL       = 12
	SMS_CARRY_COIN       = 13
	SMS_REAL_NAME        = 14
	SMS_MAX              = 15
)

// 各项目通信的redis dbs
const COMMON_REDIS_DB = 5
