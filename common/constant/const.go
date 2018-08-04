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
	APPLY_FOR_SECOND_NOT_ALREADY=8 //二级认证没有通过
	APPLY_FOR_FIRST_NOT_ALREADY =16 //一级认证状态未通过
)

//实名认证
const(
	FIRST_NOT_V=1 //没认证
	FIRST_ALREADY=2 //已经通过认证
	FIRST_NOT_ALREADY=3//没有通过认证
	FIRST_VERIFYING=4 //认证中
)

const(
	SECOND_NOT_V=1//没认证
	SECOND_ALREADY=2//已经通过认证
	SECOND_NOT_ALREADY=3//没有通过认证
	SECOND_VERIFYING=4//认证中
)
