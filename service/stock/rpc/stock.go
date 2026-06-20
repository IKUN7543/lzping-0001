package main

import (
	"flag"
	"fmt"
	"go-zero-ecommerce/service/stock/rpc/config"
	"go-zero-ecommerce/service/stock/rpc/server"
	"go-zero-ecommerce/service/stock/rpc/stock"
	svc2 "go-zero-ecommerce/service/stock/rpc/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/stock.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		stock.RegisterStockServer(grpcServer, server.NewStockServer(ctx))
		reflection.Register(grpcServer)
	})
	defer s.Stop()

	group := service.NewServiceGroup()
	group.Add(s)

	fmt.Printf("Starting stock rpc server at %s...\n", c.ListenOn)
	group.Start()
}
