package dao

import (
	cf "digicon/price_service/conf"
	. "digicon/price_service/log"
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
		Password: pass, // no password set
		DB:       11,    // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		Log.Fatalf("redis connect faild ")
	}
	Log.Infoln(pong)
	_, err = client.ZRangeWithScores("token:1/2", 0, 1).Result()
	if err != nil {
		Log.Fatalf("redis connect faild ")
	}

	return &RedisCli{
		rcon: client,
	}
}

func (s *Dao) GetRedisConn() *redis.Client {
	return s.redis.rcon
}
