package main

import (
	"flag"
	"fmt"
	"go-zero-ecommerce/service/order/api/internal/config"
	"go-zero-ecommerce/service/order/api/internal/handler"
	"go-zero-ecommerce/service/order/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/order-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	server.AddRoutes(handler.Routes(ctx))

	fmt.Printf("Starting order api server at %s...\n", c.Host+":"+c.Port)
	server.Start()
}
