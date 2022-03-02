package handler

import (
	"context"
	"github.com/micro/go-micro/util/log"
	"service/order/model"
	order "service/order/proto/order"
	"service/order/utils"
	"strconv"
)

type Order struct{}

func (o Order) CreateOrder(ctx context.Context, request *order.CreateReq, response *order.CreateResp) error {
	orderId, err:= model.CreateOrder(request.HouseId, request.StartDate, request.EndDate, request.UserName)
	if err != nil{
		log.Log("创建订单信息失败：", err)
		response.Errno = utils.RECODE_DBWRITERR
		response.Errmsg = utils.RecodeText(utils.RECODE_DBWRITERR)
		return err
	}

	response.Errno = utils.RECODE_OK
	response.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var orderData order.OrderData
	orderData.OrderId = strconv.Itoa(orderId)

	response.Data = &orderData
	log.Log("response：", response)

	return nil
}

