package dao

import (
	"github.com/go-redis/redis"
	cf "digicon/user_service/conf"
)

type RedisCli struct {
	rcon   *redis.Client
}

func NewRedisCli() *RedisCli {

	addr := cf.Cfg.MustValue("redis", "addr")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer client.Close()
	return &RedisCli{
		rcon: client,
	}
}
