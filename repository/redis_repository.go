package repository

import (
	"context"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

func GetOrSet[T any](
	ctx context.Context,
	rdb *redis.Client,
	key string,
	ttl time.Duration,
	fetchFunc func() (T, int64, error),
) (T, int64, error) {
	var data T
	var total int64

	var cache struct {
		Data  T     `json:"data"`
		Total int64 `json:"total"`
	}

	val, err := rdb.Get(ctx, key).Result()
	if err == nil {

		errUnmarshal := sonic.Unmarshal([]byte(val), &cache)
		if errUnmarshal == nil {
			return cache.Data, cache.Total, nil
		}
		log.Printf(errUnmarshal.Error())
	} else if err != redis.Nil {
		// log redis
	}

	data, total, err = fetchFunc()
	if err != nil {
		return data, total, err
	}

	cache.Data = data
	cache.Total = total
	marshaledData, errMarshal := sonic.Marshal(cache)

	go func() {
		if errMarshal == nil {
			rdb.Set(context.Background(), key, marshaledData, ttl)
		}
	}()

	return data, total, nil
}

func GetOrSetWithValidation[T any](
	ctx context.Context,
	rdb *redis.Client,
	key string,
	ttl time.Duration,
	fetchFunc func() (T, error),
	validateFunc func(T) bool,
) (T, error) {
	var data T
	var cacheHit bool

	val, err := rdb.Get(ctx, key).Result()
	if err == nil {
		if errUnmarshal := sonic.Unmarshal([]byte(val), &data); errUnmarshal == nil {
			cacheHit = true
		}
	}

	if cacheHit {
		isValid := validateFunc(data)
		if isValid {
			return data, nil
		}
	}

	data, err = fetchFunc()
	if err != nil {
		return data, err
	}

	marshaledData, errMarshal := sonic.Marshal(data)

	go func() {
		if errMarshal == nil {
			rdb.Set(context.Background(), key, marshaledData, ttl)
		}
	}()

	return data, nil
}

func InvalidateCache(ctx context.Context, rdb *redis.Client, key string) error {
	return rdb.Del(ctx, key).Err()
}
