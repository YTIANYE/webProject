package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"webProject/service/user/handler"
	user "webProject/service/user/proto/user"
)

func main() {
	// 初始化consul
	consulReg := consul.NewRegistry()

	// New Service   ---   指定consul
	service := micro.NewService(

		micro.Address("192.168.17.129:12341"),
		micro.Name("go.micro.srv.user"),
		micro.Registry(consulReg),
		micro.Version("latest"),
	)

	// Register Handler
	user.RegisterUserHandler(service.Server(), new(handler.User))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
