package redis

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type RedisDb struct {
	RedisDbClient *redis.Client
}

func (d *RedisDb) Insert(key, value string) error {
	err := d.RedisDbClient.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	fmt.Printf("Inserted key: %s, value: %s\n", key, value)
	return nil
}

func (d *RedisDb) Update(key, newValue string) error {
	// Use the SET command to update the value of an existing key
	err := d.RedisDbClient.Set(key, newValue, 0).Err()
	if err != nil {
		return err
	}
	fmt.Printf("Updated key: %s, new value: %s\n", key, newValue)
	return nil
}

func (d *RedisDb) Find(key string) (*string, error) {
	// Use the SET command to update the value of an existing key
	val, err := d.RedisDbClient.Get(key).Result()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Updated key: %s, new value: %s\n", key, val)
	return &val, nil
}
