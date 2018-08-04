package dao

import (
	//"github.com/go-redis/redis"
	// "github.com/golang/glog"
	cf "digicon/currency_service/conf"
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
		DB:       0,    // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("redis connect faild ")
	}
	log.Infoln(pong)

	return &RedisCli{
		rcon: client,
	}
}

func (s *Dao) GetRedisConn() *redis.Client {
	return s.redis.rcon
}
