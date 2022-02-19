package handler

import (
	"context"
	"service/house/model"
	"service/house/utils"
	"strconv"

	"github.com/micro/go-micro/util/log"

	house "service/house/proto/house"
)

type House struct{}

// Call is a single request handler called via client.Call or the generated client code

// 获取用户的房屋信息
func (e *House) GetHouseInfo(ctx context.Context, req *house.InfoReq, resp *house.InfoResp) error {
	log.Log("Received House.GetHouseInfo request")

	//根据用户名获取所有的房屋数据
	houseInfos,err := model.GetUserHouse(req.UserName)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	// 返回成功信息
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var getData house.GetData
	getData.Houses = houseInfos
	resp.Data = &getData

	return nil
}

// 发布房屋信息
func (e *House) PubHouse(ctx context.Context, req *house.PubReq, resp *house.PubResp) error {
	log.Log("Received House.PubHouse request")

	//发布房屋
	houseId,err := model.AddHouse(req)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	// 返回成功信息
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var h house.HouseData
	h.HouseId = strconv.Itoa(houseId)
	resp.Data = &h
	return nil
}
