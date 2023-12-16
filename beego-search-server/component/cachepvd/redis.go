package cachepvd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	beeLogger "github.com/beego/bee/v2/logger"
	"github.com/redis/go-redis/v9"
)

type redisProvider struct {
	client    *redis.Client
	expiredIn int
}

func NewRedisProvider() *redisProvider {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if host == "" || port == "" || password == "" || err != nil {
		beeLogger.Log.Fatal(ErrProviderIsNotConfigured.Error())
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	expiredIn, err := strconv.Atoi(os.Getenv("REDIS_EXPIRED_IN"))
	if err != nil {
		expiredIn = 180
	}

	return &redisProvider{
		client:    client,
		expiredIn: expiredIn,
	}
}

func (provider *redisProvider) GetCacheData(ctx context.Context, key string) (string, error) {
	return provider.client.Get(ctx, key).Result()
}

func (provider *redisProvider) SetCacheData(ctx context.Context, key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return provider.client.SetEx(ctx, key, jsonData, time.Duration(provider.expiredIn)*time.Second).Err()
}
