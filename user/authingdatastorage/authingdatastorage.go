package authingdatastorage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sbasestarter/bizinters/userinters"
)

func NewRedisAuthingDataStorage(redisCli *redis.Client, groupKey string) userinters.AuthingDataStorage {
	if redisCli == nil {
		return nil
	}

	return &redisAuthingDataStorageImpl{
		redisCli: redisCli,
		groupKey: groupKey,
	}
}

type redisAuthingDataStorageImpl struct {
	redisCli *redis.Client
	groupKey string
}

func (impl *redisAuthingDataStorageImpl) Store(ctx context.Context, ad *userinters.AuthingData, expiration time.Duration) error {
	d, err := json.Marshal(ad)
	if err != nil {
		return err
	}

	return impl.redisCli.Set(ctx, impl.key(ad.UniqueID), d, expiration).Err()
}

func (impl *redisAuthingDataStorageImpl) Load(ctx context.Context, uniqueID uint64) (ad *userinters.AuthingData, err error) {
	d, err := impl.redisCli.Get(ctx, impl.key(uniqueID)).Bytes()
	if err != nil {
		return
	}

	ad = &userinters.AuthingData{}

	err = json.Unmarshal(d, &ad)

	return
}

func (impl *redisAuthingDataStorageImpl) Delete(ctx context.Context, uniqueID uint64) error {
	return impl.redisCli.Del(ctx, impl.key(uniqueID)).Err()
}

//
//
//

func (impl *redisAuthingDataStorageImpl) key(uniqueID uint64) string {
	key := fmt.Sprintf("authing_data:%d", uniqueID)

	if impl.groupKey != "" {
		key = impl.groupKey + ":" + key
	}

	return key
}
