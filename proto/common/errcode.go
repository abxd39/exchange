package errdefine

import "github.com/gin-gonic/gin"

var ERR_CODE_RET = "code"
var ERR_CODE_MESSAGE = "msg"
var RET_DATA = "data"

var message map[int32]string

const (
	//0-49 base
	ERRCODE_SUCCESS = 0
	ERRCODE_UNKNOWN = 1
	ERRCODE_PARAM   = 2

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

	//300-

	ERRCODE_ADS_NOTEXIST      = 301
	ERRCODE_TOKENS_NOTEXIST   = 302
	ERRCODE_PAYS_NOTEXIST     = 303
	ERRCODE_ADS_TYPE_NOTEXIST = 304
	ERRCODE_ORDER_NOTEXIST    = 305
	ERRCODE_ADS_EXISTS        = 306

	//400-
	ERR_TOKEN_QUENE_NIL  = 401
	ERR_TOKEN_LESS       = 402
	ERR_TOKEN_REPEAT     = 403
	ERR_TOKEN_QUENE_CONF = 404

	ERRCODE_ORDER_FREEZE = 420
	ERRCODE_SELLER_LESS  = 421
	ERRCODE_USER_BALANCE = 422
	ERRCODE_ORDER_ERROR  = 423
	ERRCODE_TRADE_ERROR  = 424
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
	message = make(map[int32]string, 0)
	message[ERRCODE_SUCCESS] = "成功"
	message[ERRCODE_UNKNOWN] = "未知错误"
	message[ERRCODE_PARAM] = "参数错误"

	message[ERRCODE_ACCOUNT_EXIST] = "账户已经存在"
	message[ERRCODE_ACCOUNT_NOTEXIST] = "账户不存在"
	message[ERRCODE_PWD] = "密码错误"
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

	message[ERRCODE_ADS_NOTEXIST] = "广告不存在"
	message[ERRCODE_TOKENS_NOTEXIST] = "货币类型不存在"
	message[ERRCODE_PAYS_NOTEXIST] = "支付方式不存在"
	message[ERRCODE_ADS_TYPE_NOTEXIST] = "广告类型不存在"
	message[ERRCODE_ORDER_NOTEXIST] = "订单不存在"
	message[ERRCODE_ADS_EXISTS] = "广告已存在"

	message[ERRCODE_SELLER_LESS] = "卖家余额不足"
	message[ERRCODE_USER_BALANCE] = "查询用户余额失败"
	message[ERRCODE_ORDER_ERROR] = "下单失败"
	message[ERRCODE_TRADE_ERROR] = "交易失败，请重试!"

	message[ERR_TOKEN_QUENE_NIL] = "队列为空"
	message[ERR_TOKEN_LESS] = "币的余额不够"
	message[ERR_TOKEN_REPEAT] = "重复请求"
	message[ERR_TOKEN_QUENE_CONF] = "队列未配置"

	message[ERRCODE_ACCOUNT_BANK_CARD_NUMBER_MISMATCH] = "两次输入的银行卡号码不相同"
	message[ERRCODE_ARTICLE_NOT_EXIST] = "文章不存在"
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

	if code == ERRCODE_PARAM {
		if len(err_msg) > 0 {
			s.ret[ERR_CODE_MESSAGE] = err_msg[0]
		} else {
			s.ret[ERR_CODE_MESSAGE] = GetErrorMessage(code)
		}
	} else if code != ERRCODE_UNKNOWN {
		s.ret[ERR_CODE_MESSAGE] = GetErrorMessage(code)
	} else {
		if len(err_msg) > 0 {
			s.ret[ERR_CODE_MESSAGE] = err_msg[0]
		}
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
