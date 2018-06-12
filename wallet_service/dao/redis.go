package dao

import (
	//"github.com/go-redis/redis"
	// "github.com/golang/glog"
	"fmt"
	redis "github.com/garyburd/redigo/redis"
)

type RedisCli struct {
	redis redis.Conn
}

func NewRedisCli() *RedisCli {

	client, err := redis.Dial("tcp", "47.106.136.96:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return nil
	}

	defer client.Close()
	return &RedisCli{
		redis: client,
	}
}
