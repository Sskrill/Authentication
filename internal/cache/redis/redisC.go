package redisC

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRedis struct {
	cl *redis.Client
}

func NewRedisCache(url, password string) *CacheRedis {
	cl := redis.NewClient(&redis.Options{Addr: url, Password: password})
	return &CacheRedis{cl: cl}
}

func (cR *CacheRedis) Save(key string, value []byte) error {
	ctx := context.Background()
	err := cR.cl.Set(ctx, key, value, time.Minute*1).Err()
	if err != nil {
		return err
	}
	return nil
}
func (cR *CacheRedis) Get(key string) ([]byte, error) {
	ctx := context.Background()
	val, err := cR.cl.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return val, nil
}
