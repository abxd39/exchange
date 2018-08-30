package errdefine

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var ERR_CODE_RET = "code"
var ERR_CODE_MESSAGE = "msg"
var RET_DATA = "data"

var message map[int32]string

const (
	//0-49 base
	ERRCODE_SUCCESS     = 0
	ERRCODE_UNKNOWN     = 1
	ERRCODE_PARAM       = 2
	ERRCODE_TOKENVERIFY = 3
	ERRCODE_GREE        = 4

	ERRCODE_NORMAL_ERROR = 100

	//200-
	ERRCODE_ACCOUNT_EXIST     = 202
	ERRCODE_ACCOUNT_NOTEXIST  = 203
	ERRCODE_PWD               = 204
	ERRCODE_PWD_COMFIRM       = 205
	ERRCODE_SECURITY_KEY      = 206
	ERRCODE_SMS_CODE_DIFF     = 207
	ERRCODE_SMS_CODE_NIL      = 208
	ERRCODE_SMS_COMMIT_QUICK  = 209
	ERRCODE_SMS_SYS_BUSY      = 210
	ERRCODE_SMS_MONEY_ENGOUGE = 211
	ERRCODE_SMS_PHONE_FORMAT  = 212
	ERRCODE_SMS_EMAIL_FORMAT  = 213

	ERRCODE_GOOGLE_CODE           = 214
	ERRCODE_GOOGLE_CODE_EXIST     = 215
	ERRCODE_GOOGLE_CODE_NOT_EXIST = 216

	ERRCODE_ACCOUNT_BANK_CARD_NUMBER_MISMATCH = 217
	ERRCODE_ARTICLE_NOT_EXIST                 = 218
	ERRCODE_OLDPWD                            = 219

	ERRCODE_PHONE_EXIST     = 220
	ERRCODE_PHONE_NOT_EXIST = 221
	ERRCODE_INVITE          = 222

	ERRCODE_EMAIL_EXIST   = 223
	ERRCODE_UOPLOA_FAILED = 224

	//300-

	ERRCODE_ADS_NOTEXIST       = 301
	ERRCODE_TOKENS_NOTEXIST    = 302
	ERRCODE_PAYS_NOTEXIST      = 303
	ERRCODE_ADS_TYPE_NOTEXIST  = 304
	ERRCODE_ORDER_NOTEXIST     = 305
	ERRCODE_ADS_EXISTS         = 306
	ERRCODE_ADS_SET_PRICE      = 307
	ERRCODE_ADS_MIN_LIMIT      = 308
	ERRCODE_ADS_NEED_TWO_LEVEL = 309
	ERRCODE_ADS_NUM_LESS       = 310

	//400-
	ERR_TOKEN_QUENE_NIL      = 401
	ERR_TOKEN_LESS           = 402
	ERR_TOKEN_REPEAT         = 403
	ERR_TOKEN_QUENE_CONF     = 404
	ERR_TOKEN_ENTRUST_STATES = 405
	ERR_TOKEN_ENTRUST_EXIST  = 406

	ERRCODE_ORDER_FREEZE        = 420
	ERRCODE_SELLER_LESS         = 421
	ERRCODE_USER_BALANCE        = 422
	ERRCODE_ORDER_ERROR         = 423
	ERRCODE_TRADE_ERROR         = 424
	ERRCODE_TRADE_ERROR_ADS_NUM = 425
	ERRCODE_TRADE_LOWER_PRICE   = 426
	ERRCODE_TRADE_LARGE_PRICE   = 427
	ERRCODE_TRADE_TO_SELF       = 428
	ERRCODE_TRADE_HAS_COMPLETED = 429
	ERRCODE_SAVE_ERROR = 430
	ERRCODE_SELECT_ERROR = 431
	ERRCODE_RPC_ERROR = 432
	ERRCODE_CREATE_ERROR = 433
	ERRCODE_TOKEN_INVALID = 434
	ERRCODE_PAY_PWD = 435
	ERRCODE_TOKEN_NOT_ENOUGH = 436
	ERRCODE_PARSE = 437
	ERRCODE_CNY_PRICE = 438
	ERRCODE_FORMAT = 439
	ERRCODE_FREEZE = 440
	ERRCODE_TIBI_ADDRESS = 441
)

func GetErrorMessage(code int32) string {
	return message[code]
}

func CheckErrorMessage(code int32) (ret string, ok bool) {
	ret, ok = message[code]
	if ok {
		return
	}
	return
}

func NewErrorMessage() gin.H {
	var ret = gin.H{}
	data := make(map[string]interface{}, 0)
	ret[ERR_CODE_RET] = 0
	ret[ERR_CODE_MESSAGE] = 0
	ret[RET_DATA] = data
	return ret
}

func init() {
	errors.New("{err:}")
	message = make(map[int32]string, 0)
	message[ERRCODE_SUCCESS] = "成功"
	message[ERRCODE_UNKNOWN] = "未知错误"
	message[ERRCODE_PARAM] = "参数错误"
	message[ERRCODE_TOKENVERIFY] = "令牌失效"
	message[ERRCODE_GREE] = "智能验证失败"

	message[ERRCODE_ACCOUNT_EXIST] = "账户已经存在"
	message[ERRCODE_ACCOUNT_NOTEXIST] = "账户不存在"
	message[ERRCODE_PWD] = "密码错误"
	message[ERRCODE_OLDPWD] = "旧密码不匹配"
	message[ERRCODE_PWD_COMFIRM] = "确认密码不一致"
	message[ERRCODE_SECURITY_KEY] = "安全码不一致"
	message[ERRCODE_SMS_CODE_DIFF] = "验证码错误"
	message[ERRCODE_SMS_CODE_NIL] = "验证码未获取"

	message[ERRCODE_SMS_COMMIT_QUICK] = "提交过快"
	message[ERRCODE_SMS_SYS_BUSY] = "系统忙"
	message[ERRCODE_SMS_MONEY_ENGOUGE] = "无发送额度"
	message[ERRCODE_SMS_PHONE_FORMAT] = "手机号格式错误"
	message[ERRCODE_SMS_EMAIL_FORMAT] = "邮箱格式错误"
	message[ERRCODE_GOOGLE_CODE] = "谷歌验证码错误"
	message[ERRCODE_GOOGLE_CODE_EXIST] = "谷歌验证码已经存在无法重复拉取"
	message[ERRCODE_GOOGLE_CODE_NOT_EXIST] = "谷歌验证码不存在无法解绑"
	message[ERRCODE_INVITE] = "邀请码不存在"

	message[ERRCODE_ADS_NOTEXIST] = "广告不存在"
	message[ERRCODE_TOKENS_NOTEXIST] = "货币类型不存在"
	message[ERRCODE_PAYS_NOTEXIST] = "支付方式不存在"
	message[ERRCODE_ADS_TYPE_NOTEXIST] = "广告类型不存在"
	message[ERRCODE_ORDER_NOTEXIST] = "订单不存在"
	message[ERRCODE_ADS_EXISTS] = "广告已存在"
	message[ERRCODE_ADS_SET_PRICE] = "当前广告的总价已小于最小价格的值"
	message[ERRCODE_ADS_MIN_LIMIT] = "限制最小价格要大于等于100"
	message[ERRCODE_ADS_NEED_TWO_LEVEL] = "对方设置了需要通过二级认证,请先进行认证"
	message[ERRCODE_ADS_NUM_LESS] = "上架时候币额不能为0"

	message[ERRCODE_SELLER_LESS] = "卖家余额不足"
	message[ERRCODE_USER_BALANCE] = "查询用户余额失败"
	message[ERRCODE_ORDER_ERROR] = "下单失败"
	message[ERRCODE_TRADE_ERROR_ADS_NUM] = "下单失败,购买的数量大于订单的数量!"
	message[ERRCODE_TRADE_ERROR] = "交易失败，请重试!"
	message[ERRCODE_ORDER_FREEZE] = "订单冻结"
	message[ERRCODE_TRADE_LOWER_PRICE] = "下单失败,买价小于允许的最小价格!"
	message[ERRCODE_TRADE_LARGE_PRICE] = "下单失败,买价大于允许的最大价格!"
	message[ERRCODE_TRADE_TO_SELF] = "不能下自己的单!"
	message[ERRCODE_TRADE_HAS_COMPLETED] = "订单已经完成!"

	message[ERR_TOKEN_QUENE_NIL] = "队列为空"
	message[ERR_TOKEN_LESS] = "币的余额不够"
	message[ERR_TOKEN_REPEAT] = "重复请求"
	message[ERR_TOKEN_QUENE_CONF] = "队列未配置"
	message[ERR_TOKEN_ENTRUST_STATES] = "委托状态错误"
	message[ERR_TOKEN_ENTRUST_EXIST] = "委托不存在"

	message[ERRCODE_ACCOUNT_BANK_CARD_NUMBER_MISMATCH] = "两次输入的银行卡号码不相同"
	message[ERRCODE_ARTICLE_NOT_EXIST] = "文章不存在"
	message[ERRCODE_PHONE_EXIST] = "电话号码已经存在"
	message[ERRCODE_PHONE_NOT_EXIST] = "电话号码不存在"
	message[ERRCODE_EMAIL_EXIST] = "邮箱已经存在"
	message[ERRCODE_UOPLOA_FAILED] = "上传图片到oss 失败"

	message[ERRCODE_SAVE_ERROR] = "保存数据失败"
	message[ERRCODE_SELECT_ERROR] = "查询失败"
	message[ERRCODE_CREATE_ERROR] = "创建失败"
	message[ERRCODE_TOKEN_INVALID] = "Token暂不可用"
	message[ERRCODE_PAY_PWD] = "支付密码错误"
	message[ERRCODE_TOKEN_NOT_ENOUGH] = "余额不足"
	message[ERRCODE_CNY_PRICE] = "获取价格出错"
	message[ERRCODE_FORMAT] = "格式化数据失败"
	message[ERRCODE_FREEZE] = "冻结失败"
	message[ERRCODE_TIBI_ADDRESS] = "提币地址未配置"
}

type PublicErrorType struct {
	ret  gin.H
	data map[string]interface{}
}

//创建统一错误返回格式
func NewPublciError() *PublicErrorType {
	s := new(PublicErrorType)
	s.init()
	return s
}

//初始化操作
func (s *PublicErrorType) init() {
	var ret = gin.H{}
	ret[ERR_CODE_RET] = 0
	ret[ERR_CODE_MESSAGE] = 0
	s.ret = ret
	s.data = make(map[string]interface{}, 0)
}

//设置错误代码，如果有自定义错误信息填写err_msg参数
func (s *PublicErrorType) SetErrCode(code int32, err_msg ...string) {
	s.ret[ERR_CODE_RET] = code

	if len(err_msg) > 0 {
		if err_msg[0] == "" {
			s.ret[ERR_CODE_MESSAGE] = GetErrorMessage(code)
		} else {
			s.ret[ERR_CODE_MESSAGE] = err_msg[0]
		}
	} else {
		s.ret[ERR_CODE_MESSAGE] = GetErrorMessage(code)
	}
}

//设置数据部分内容
func (s *PublicErrorType) SetDataSection(key string, value interface{}) {
	s.data[key] = value
}

//返回最终的数据
func (s *PublicErrorType) GetResult() gin.H {
	s.ret[RET_DATA] = s.data
	return s.ret
}

// 补充-设置数据部分内容
func (s *PublicErrorType) SetDataValue(value interface{}) {
	s.ret[RET_DATA] = value
}

// 补充-返回最终的数据
func (s *PublicErrorType) GetData() gin.H {
	return s.ret
}
