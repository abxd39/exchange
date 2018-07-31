package cron

import (
	"fmt"

	"digicon/user_service/model"
	"digicon/user_service/rpc/handler"

	"github.com/robfig/cron"
)

func InitCron() {
	c := cron.New()
	c.AddFunc("0 0 3 * * *", RegisterNoReward)
	c.Start()
}

// 注册未收到奖励，补上奖励
func RegisterNoReward() {
	list, _ := new(model.User).GetRegisterNoRewardUser()

	rpcServer := new(handler.RPCServer)
	for _, v := range list {
		fmt.Println("in", v)
		rpcServer.RegisterReward(v.Uid, v.InviteId)
	}
}
