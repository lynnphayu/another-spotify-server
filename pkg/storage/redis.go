package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(connectionString string) (*Cache, error) {
	redisInstance := new(Cache)
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	redisInstance.client = client
	return redisInstance, nil
}

//set value in redis
func (redisInstance *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	return redisInstance.client.Set(context.TODO(), key, value, expiration).Err()
}

func (redisInstance *Cache) Get(key string) (interface{}, error) {
	return redisInstance.client.Get(context.TODO(), key).Result()
}

// clear cache with key from parameter
func (redisInstance *Cache) Clear(key string) error {
	return redisInstance.client.Del(context.TODO(), key).Err()
}
