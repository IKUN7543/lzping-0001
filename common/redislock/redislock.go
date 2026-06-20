package redislock

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go-zero-ecommerce/common/errx"
	"time"
)

type RedisLock struct {
	rdb        *redis.Client
	key        string
	value      string
	expiration time.Duration
}

func NewRedisLock(rdb *redis.Client, key string, expiration time.Duration) *RedisLock {
	return &RedisLock{
		rdb:        rdb,
		key:        key,
		value:      uuid.New().String(),
		expiration: expiration,
	}
}

var lockScript = `
if redis.call('SETNX', KEYS[1], ARGV[1]) == 1 then
    return redis.call('PEXPIRE', KEYS[1], ARGV[2])
else
    return 0
end
`

func (l *RedisLock) Lock(ctx context.Context) error {
	expireMs := l.expiration.Milliseconds()
	result, err := l.rdb.Eval(ctx, lockScript, []string{l.key}, l.value, expireMs).Result()
	if err != nil {
		return err
	}
	if result.(int64) == 0 {
		return errx.ErrLockTimeout
	}
	return nil
}

var unlockScript = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call('DEL', KEYS[1])
else
    return 0
end
`

func (l *RedisLock) Unlock(ctx context.Context) error {
	result, err := l.rdb.Eval(ctx, unlockScript, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}
	if result.(int64) == 0 {
		return errx.ErrLockTimeout
	}
	return nil
}

func (l *RedisLock) TryLock(ctx context.Context, retryTimes int, retryInterval time.Duration) error {
	for i := 0; i < retryTimes; i++ {
		err := l.Lock(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(retryInterval)
	}
	return errx.ErrLockTimeout
}
