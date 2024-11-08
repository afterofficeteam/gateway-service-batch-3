package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	Host     string
	Port     int
	Password string
	DB       int
}

var ctx = context.Background()

func ConnectToRedis(conn RedisConnection) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conn.Host, conn.Port),
		Password: conn.Password,
		DB:       conn.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return rdb, nil
}
