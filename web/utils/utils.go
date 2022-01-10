package utils

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/consul"
)

// 初始化micro客户端
func InitMicro() micro.Service {

	// 初始化consul客户端
	consulReg := consul.NewRegistry()
	return micro.NewService(
		micro.Registry(consulReg),
	)

}
