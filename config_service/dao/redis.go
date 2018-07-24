package dao

import (
	cf "digicon/config_service/conf"
	log "github.com/sirupsen/logrus"
	"github.com/go-redis/redis"
)

type RedisCli struct {
	rcon *redis.Client
}

func NewRedisCli() *RedisCli {
	addr := cf.Cfg.MustValue("redis", "addr")
	pass := cf.Cfg.MustValue("redis", "pass")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("redis connect faild")
	}
	log.Infoln(pong)

	return &RedisCli{
		rcon: client,
	}
}
