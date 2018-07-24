package dao

import (
	cf "digicon/price_service/conf"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type RedisCli struct {
	rcon *redis.Client
}

func NewRedisCli() *RedisCli {

	addr := cf.Cfg.MustValue("redis", "addr")
	pass := cf.Cfg.MustValue("redis", "pass")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass, // no password set
		DB:       11,   // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("redis connect faild ")
	}
	log.Infoln(pong)
	_, err = client.ZRangeWithScores("token:1/2", 0, 1).Result()
	if err != nil {
		log.Fatalf("redis connect faild ")
	}

	return &RedisCli{
		rcon: client,
	}
}

func (s *Dao) GetRedisConn() *redis.Client {
	return s.redis.rcon
}
