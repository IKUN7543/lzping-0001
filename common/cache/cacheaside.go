package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheAside struct {
	rdb        *redis.Client
	expiration time.Duration
	prefix     string
}

func NewCacheAside(rdb *redis.Client, prefix string, expiration time.Duration) *CacheAside {
	return &CacheAside{
		rdb:        rdb,
		prefix:     prefix,
		expiration: expiration,
	}
}

func (c *CacheAside) key(id interface{}) string {
	return fmt.Sprintf("%s:%v", c.prefix, id)
}

func (c *CacheAside) Get(ctx context.Context, id interface{}, dest interface{}) error {
	data, err := c.rdb.Get(ctx, c.key(id)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *CacheAside) Set(ctx context.Context, id interface{}, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, c.key(id), data, c.expiration).Err()
}

func (c *CacheAside) Del(ctx context.Context, id interface{}) error {
	return c.rdb.Del(ctx, c.key(id)).Err()
}

func (c *CacheAside) LoadOrStore(ctx context.Context, id interface{}, dest interface{}, load func() (interface{}, error)) error {
	err := c.Get(ctx, id, dest)
	if err == nil {
		return nil
	}
	value, err := load()
	if err != nil {
		return err
	}
	if err := c.Set(ctx, id, value); err != nil {
		return err
	}
	data, _ := json.Marshal(value)
	return json.Unmarshal(data, dest)
}
