package handler

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/util/log"
	"service/order/model"
	order "service/order/proto/order"
	"service/order/utils"
	"strconv"
)

type Order struct{}


// 创建订单

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

// 查询用户订单

func (o Order) GetUserOrder(ctx context.Context, req *order.GetReq, resp *order.GetResp) error {
	respData, err := model.GetOrderInfo(req.GetRole(), req.GetUserName())
	if err != nil{
		fmt.Println("获取用户订单信息失败：", err)
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return err
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var getData order.GetData
	getData.Orders = respData
	resp.Data = &getData

	return nil
}

// 接受或拒绝订单

func (o Order) StateOrder(ctx context.Context, req *order.StateReq, resp *order.StateResp) error {
	err := model.StateOrder(req.Action, req.Reason, req.Id)
	if err != nil{
		fmt.Println("更新订单状态失败：", err )
		resp.Errno = utils.RECODE_DBWRITERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return err
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)

	return nil
}