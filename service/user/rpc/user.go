package main

import (
	"flag"
	"fmt"
	"go-zero-ecommerce/service/user/rpc/config"
	"go-zero-ecommerce/service/user/rpc/server"
	"go-zero-ecommerce/service/user/rpc/svc"
	"go-zero-ecommerce/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		reflection.Register(grpcServer)
	})
	defer s.Stop()

	group := service.NewServiceGroup()
	group.Add(s)

	fmt.Printf("Starting user rpc server at %s...\n", c.ListenOn)
	group.Start()
}
