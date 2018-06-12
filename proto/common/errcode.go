package errdefine

import "github.com/gin-gonic/gin"

var ErrCodeRet = "code"
var ErrCodeMessage = "msg"
var RetData = "data"

var message map[int32]string

const (
	ERRCODE_SUCCESS          = 0
	ERRCODE_UNKNOWN          = 1
	ERRCODE_PARAM            = 2
	ERRCODE_ACCOUNT_EXIST    = 4
	ERRCODE_ACCOUNT_NOTEXIST = 5
	ERRCODE_PWD              = 6
	ERRCODE_PWD_COMFIRM      = 7
	ERRCODE_SECURITY_KEY     = 8
)

func GetErrorMessage(code int32) string {
	return message[code]
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
}
