package cache

import (
	"github.com/go-redis/redis"
	"exercise/conf"
)

var Client *redis.Client

func Init() error {
	Client = redis.NewClient(&redis.Options{
		Addr: conf.GetConf().Redis.Addr,
		Password: conf.GetConf().Redis.Password,
	})
	_, err := Client.Ping().Result()
	return err
}
