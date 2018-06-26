package dao

import (
	cf "digicon/user_service/conf"
	"digicon/user_service/log"
	"time"

	"github.com/go-redis/redis"
)

type RedisCli struct {
	rcon   *redis.Client
	KeyTtl time.Duration
	salt   string
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
		log.Log.Fatalf("redis connect faild ")
	}
	log.Log.Infoln(pong)

	ct, err := cf.Cfg.Int64("redis", "ttl")
	if err != nil {
		ct = 30
	}
	return &RedisCli{
		rcon:   client,
		salt:   "mjfdsap832-1##1!",
		KeyTtl: time.Duration(ct) * time.Second,
	}
}

func (rds *Dbs) GetRedisInstance() *redis.Client {
	return rds.redis.rcon
}
