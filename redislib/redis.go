package redislib

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedis(dsn string) (cli *redis.Client, err error) {
	options, err := redis.ParseURL(dsn)
	if err != nil {
		return
	}

	cli = redis.NewClient(options)

	ctx, cf := context.WithTimeout(context.Background(), 3*time.Second)
	defer cf()

	err = cli.Ping(ctx).Err()
	if err != nil {
		return
	}

	return
}
