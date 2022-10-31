package statuscontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sbasestarter/bizinters/userinters"
	"github.com/sgostarter/libeasygo/commerr"
)

func NewRedisStatusController(redisCli *redis.Client) (userinters.StatusController, error) {
	if redisCli == nil {
		return nil, commerr.ErrInvalidArgument
	}

	return &redisStatusControllerImpl{
		redisCli: redisCli,
	}, nil
}

type redisStatusControllerImpl struct {
	redisCli *redis.Client
}

func (impl *redisStatusControllerImpl) IsTokenBanned(ctx context.Context, id uint64) (banned bool, err error) {
	n, err := impl.redisCli.Exists(ctx, impl.tokenBanKey(id)).Result()
	if err != nil {
		return
	}

	banned = n > 0

	return
}

func (impl *redisStatusControllerImpl) BanToken(ctx context.Context, id uint64, expireAt int64) error {
	var expiration time.Duration
	if expireAt > 0 {
		expiration = time.Duration(expireAt - time.Now().Unix())
		if expiration < 0 {
			return nil
		}
	}

	return impl.redisCli.Set(ctx, impl.tokenBanKey(id), time.Now(), expiration).Err()
}

func (impl *redisStatusControllerImpl) tokenBanKey(id uint64) string {
	return fmt.Sprintf("banned:token:%d", id)
}
