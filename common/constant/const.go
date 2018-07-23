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
)
