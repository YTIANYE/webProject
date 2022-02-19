package utils

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/registry/consul"
)

// 初始化consul

func InitMicro() micro.Service {

	// 初始化consul客户端
	consulReg := consul.NewRegistry()
	return micro.NewService(
		micro.Registry(consulReg),
	)
}

// 初始化micro客户端

func GetMicroClient() client.Client{
	consulReg := consul.NewRegistry()
	microService := micro.NewService(
		micro.Registry(consulReg),
	)
	return microService.Client()
}
