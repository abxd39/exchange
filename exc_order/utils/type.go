package utils

//登录数据
type LoginData struct {
	Result bool
	Uid int64
	Token string
}

//拉取数据结果
type PullData struct {
	Result bool
	Price float64
	Amount float64}
