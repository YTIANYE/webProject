package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"webProject/service/user/handler"
	"webProject/service/user/model"
	user "webProject/service/user/proto/user"
)

func main() {
	// 初始化 Redis 连接池
	model.InitRedis()

	// 初始化consul
	consulReg := consul.NewRegistry()

	// New Service   ---   指定consul
	service := micro.NewService(

		micro.Address("192.168.17.129:12341"), // 指定固定端口
		micro.Name("go.micro.srv.user"),
		micro.Registry(consulReg), //注册服务
		micro.Version("latest"),
	)

	// Register Handler
	user.RegisterUserHandler(service.Server(), new(handler.User))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
