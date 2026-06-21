package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.CacheConf
	Redis      struct {
		Host         string
		Type         string
		Pass         string
		PoolSize     int `json:",default=100"`
		MinIdleConns int `json:",default=10"`
	}
}
