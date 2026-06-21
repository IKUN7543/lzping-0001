package main

import (
	"flag"
	"fmt"
	"go-zero-ecommerce/service/order/rpc/config"
	"go-zero-ecommerce/service/order/rpc/order"
	"go-zero-ecommerce/service/order/rpc/server"
	"go-zero-ecommerce/service/order/rpc/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	defer ctx.Close()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		order.RegisterOrderServer(grpcServer, server.NewOrderServer(ctx))
		reflection.Register(grpcServer)
	})
	defer s.Stop()

	group := service.NewServiceGroup()
	group.Add(s)

	fmt.Printf("Starting order rpc server at %s...\n", c.ListenOn)
	group.Start()
}
