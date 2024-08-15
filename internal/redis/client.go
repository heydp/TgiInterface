package redis

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v7"
)

type RedisCreds struct {
	Host     string
	Port     int64
	UserName string
	Password string
	Database int
}

func (c *RedisCreds) GiveRedisClient() (*RedisDb, error) {
	addr := fmt.Sprintf("%v:%v", c.Host, c.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: c.Password,
		DB:       c.Database, // use default db
	})

	resp, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}
	log.Printf("connected to redis successfully, got the response - %v\n", resp)

	return &RedisDb{RedisDbClient: rdb}, nil
}
