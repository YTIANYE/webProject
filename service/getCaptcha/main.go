package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"webProject/service/getCaptcha/handler"
	getCaptcha "webProject/service/getCaptcha/proto/getCaptcha"
)

func main() {
	// 服务发现
	// 初始化consul
	consulReg := consul.NewRegistry()

	// New Service
	service := micro.NewService(
		micro.Address("192.168.17.129:12341"), //防止随即生成port  //端口号不超过65535
		// micro.Address("127.0.0.1:12341"),
		micro.Name("go.micro.srv.getCaptcha"), // 别名
		micro.Registry(consulReg),             //添加注册
		micro.Version("latest"),
	)

	// Register Handler
	getCaptcha.RegisterGetCaptchaHandler(service.Server(), new(handler.GetCaptcha))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
