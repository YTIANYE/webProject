package handler

import (
	"context"
	"service/order/proto/order"
)

type Order struct{}

func (o Order) Call(ctx context.Context, request *order.Request, response *order.Response) error {
	panic("implement me")
}

