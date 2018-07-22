package dao

import  (
	cf "digicon/config_service/conf"
	. "digicon/config_service/log"
	"github.com/go-redis/redis"
)

type RedisCli struct {
	rcon *redis.Client
}

func NewRedisCli() *RedisCli{
	addr := cf.Cfg.MustValue("redis", "addr")
	pass := cf.Cfg.MustValue("redis", "pass")
	client := redis.NewClient(&redis.Options{
		Addr:addr,
		Password:pass,
		DB: 0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		Log.Fatalf("redis connect faild")
	}
	Log.Infoln(pong)

	return &RedisCli{
		rcon:client,
	}
}