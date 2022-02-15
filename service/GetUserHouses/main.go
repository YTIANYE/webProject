package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"service/GetUserHouses/handler"
	"service/GetUserHouses/subscriber"

	GetUserHouses "service/GetUserHouses/proto/GetUserHouses"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.GetUserHouses"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetUserHouses.RegisterGetUserHousesHandler(service.Server(), new(handler.GetUserHouses))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetUserHouses", service.Server(), new(subscriber.GetUserHouses))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetUserHouses", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
