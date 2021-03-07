package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"surf_be/internal/configuration"
	"time"
)

type Redis struct {
	Config     configuration.Config
	Connection *redis.Client
}

func (rd *Redis) Init() {
	rd.Connection = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", rd.Config.Server.DataBase.Redis.Host, rd.Config.Server.DataBase.Redis.Port),
		Password: "",
		DB:       0,
	})
}

func (rd *Redis) Get(ctx context.Context, key string) (interface{}, error) {
	val2, err := rd.Connection.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("key does not exist: %v", err)
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val2, nil
}

func (rd *Redis) Set(ctx context.Context, key string, value interface{}, ttlInSecond time.Duration) error {
	byteArray, err := json.Marshal(value)
	ok, err := rd.Connection.SetNX(ctx, key, byteArray, ttlInSecond*time.Second).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("redis value is existed")
	}
	return nil
}

func (rd *Redis) Delete(ctx context.Context, key string) error {
	err := rd.Connection.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
