package main

import (
	"flag"
	"fmt"
	"go-zero-ecommerce/service/product/rpc/config"
	"go-zero-ecommerce/service/product/rpc/server"
	"go-zero-ecommerce/service/product/rpc/svc"
	"go-zero-ecommerce/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/product.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		product.RegisterProductServer(grpcServer, server.NewProductServer(ctx))
		reflection.Register(grpcServer)
	})
	defer s.Stop()

	group := service.NewServiceGroup()
	group.Add(s)

	fmt.Printf("Starting product rpc server at %s...\n", c.ListenOn)
	group.Start()
}
