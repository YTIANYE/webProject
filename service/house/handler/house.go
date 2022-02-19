package handler

import (
	"context"
	"service/house/model"
	"service/house/utils"

	"github.com/micro/go-micro/util/log"

	house "service/house/proto/house"
)

type House struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *House) GetHouseInfo(ctx context.Context, req *house.InfoReq, resp *house.InfoResp) error {
	log.Log("Received House.GetHouseInfo request")
	// 查询数据库


	// 返回结果

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
