package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Host         string
		Type         string
		Pass         string
		PoolSize     int `json:",default=100"`
		MinIdleConns int `json:",default=10"`
	}
	Kafka struct {
		Brokers string
		Topic   string
		Group   string
	}
	StockRpc   zrpc.RpcClientConf
	ProductRpc zrpc.RpcClientConf
}
