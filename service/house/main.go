package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"service/house/handler"
	"service/house/model"
	house "service/house/proto/house"
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
		micro.Name("go.micro.srv.house"),
		micro.Version("latest"),
		micro.Address("192.168.17.129:12343"), // 指定固定端口 不同微服务之间要区分端口号
		micro.Registry(consulReg), //注册服务
	)

	//// Initialise service
	//service.Init()

	// Register Handler
	house.RegisterHouseHandler(service.Server(), new(handler.House))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
