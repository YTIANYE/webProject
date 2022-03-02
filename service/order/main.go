package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"service/order/handler"
	"service/order/model"
	"service/order/proto/order"
)

func main() {
	// 初始化 Redis 连接池
	model.InitRedis()

	// 初始化 MySQL 链接池
	model.InitDb()

	// 初始化consul
	consulReg := consul.NewRegistry()


	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.order"),
		micro.Version("latest"),
		micro.Address("192.168.17.129:12344"), // 指定固定端口 不同微服务之间要区分端口号
		micro.Registry(consulReg), //注册服务
	)

	// Initialise service
	service.Init()

	// Register Handler
	order.RegisterOrderHandler(service.Server(), new(handler.Order))


	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
