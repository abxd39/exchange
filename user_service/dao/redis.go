package dao

import (
	"digicon/common/encryption"
	cf "digicon/user_service/conf"
	. "digicon/user_service/log"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisCli struct {
	rcon    *redis.Client
	key_ttl time.Duration
	salt    string
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

	ct, err := cf.Cfg.Int64("redis", "ttl")
	if err != nil {
		ct = 30
	}
	return &RedisCli{
		rcon:    client,
		salt:    "mjfdsap832-1##1!",
		key_ttl: time.Duration(ct) * time.Second,
	}
}

func GetUserTag(phone string) string {
	return fmt.Sprintf("%s:SecurityKey", phone)
}

func (s *Dao) GenSecurityKey(phone string) (security_key []byte, err error) {
	security_key = encryption.Gensha256(phone, time.Now().Unix(), s.redis.salt)
	err = s.redis.rcon.Set(GetUserTag(phone), security_key, s.redis.key_ttl).Err()
	if err != nil {
		Log.Errorln(err.Error())
		return
	}
	return
}
