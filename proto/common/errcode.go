package errdefine

import "github.com/gin-gonic/gin"

var ErrCodeRet = "code"
var ErrCodeMessage = "msg"
var RetData = "data"

var message map[int32]string

const (
	//0-49 base
	ERRCODE_SUCCESS          = 0
	ERRCODE_UNKNOWN          = 1
	ERRCODE_PARAM            = 2
	ERRCODE_ACCOUNT_EXIST    = 4
	ERRCODE_ACCOUNT_NOTEXIST = 5
	ERRCODE_PWD              = 6
	ERRCODE_PWD_COMFIRM      = 7
	ERRCODE_SECURITY_KEY     = 8
	ERRCODE_SMS_CODE_DIFF    	 = 9
	ERRCODE_SMS_CODE_NIL    	 = 10

	//100-130 sms
	ERRCODE_SMS_COMMIT_QUICK  = 103
	ERRCODE_SMS_SYS_BUSY      = 104
	ERRCODE_SMS_TEL_FORMAT    = 107
	ERRCODE_SMS_MONEY_ENGOUGE = 109
	ERRCODE_SMS_PARAM         = 130
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
	ret[ErrCodeRet] = 0
	ret[ErrCodeMessage] = 0
	ret[RetData] = data

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
	message[ERRCODE_SMS_CODE_DIFF]="验证码错误"
	message[ERRCODE_SMS_CODE_NIL]="验证码未获取"



	message[ERRCODE_SMS_COMMIT_QUICK] = "提交过快"
	message[ERRCODE_SMS_SYS_BUSY] = "系统忙"
	message[ERRCODE_SMS_TEL_FORMAT] = "包含错误的手机号码"
	message[ERRCODE_SMS_MONEY_ENGOUGE] = "无发送额度"
	message[ERRCODE_SMS_PARAM] = "请求参数错误"
}
