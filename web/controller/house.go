package controller

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	houseMicro "webProject/web/proto/house"
	"webProject/web/utils"
)

// 获取用户已发布房源

func GetUserHouses(ctx *gin.Context) {
	// 获取Sesion, 得到用户名
	session := sessions.Default(ctx)
	userName := session.Get("userName")

	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//调用远程服务
	resp, _ := microClient.GetHouseInfo(context.TODO(), &houseMicro.InfoReq{UserName: userName.(string)})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

// 发布房源信息

type HouseStu struct {
	Acreage   string   `json:"acreage"`
	Address   string   `json:"address"`
	AreaId    string   `json:"area_id"`
	Beds      string   `json:"beds"`
	Capacity  string   `json:"capacity"`
	Deposit   string   `json:"deposit"`
	Facility  []string `json:"facility"`
	MaxDays   string   `json:"max_days"`
	MinDays   string   `json:"min_days"`
	Price     string   `json:"price"`
	RoomCount string   `json:"room_count"`
	Title     string   `json:"title"`
	Unit      string   `json:"unit"`
}

func PostHouses(ctx *gin.Context) {
	//校验数据   bind数据的时候不带自动转换   c.getInt()
	var house HouseStu
	err := ctx.Bind(&house)
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}

	// 获取session, 得到用户名
	s := sessions.Default(ctx)
	userName := s.Get("userName")
	fmt.Println("userName", userName)

	// 调用远程服务
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	resp, _ := microClient.PubHouse(context.TODO(), &houseMicro.PubReq{
		Acreage:   house.Acreage,
		Address:   house.Address,
		AreaId:    house.AreaId,
		Beds:      house.Beds,
		Capacity:  house.Capacity,
		Deposit:   house.Deposit,
		Facility:  house.Facility,
		MaxDays:   house.MaxDays,
		MinDays:   house.MinDays,
		Price:     house.Price,
		RoomCount: house.RoomCount,
		Title:     house.Title,
		Unit:      house.Unit,
		UserName:  userName.(string),
	})

	// 返回数据
	ctx.JSON(http.StatusOK, resp)
}
