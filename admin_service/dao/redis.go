package dao

import (
	cf "digicon/admin_service/conf"
	. "digicon/admin_service/log"
	"github.com/go-redis/redis"
)

type RedisCli struct {
	rcon *redis.Client
}

func NewRedisCli() *RedisCli {
	addr := cf.Cfg.MustValue("redis", "addr")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		Log.Fatalf("redis connect faild ")
	}
	Log.Infoln(pong)

	return &RedisCli{
		rcon: client,
	}
}
