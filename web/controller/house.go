package controller

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
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

// 添加房屋图片

func PostHousesImage(ctx *gin.Context) {

	// 获取数据
	houseId := ctx.Param("id")
	fileHeader, err := ctx.FormFile("house_image")

	// 校验数据
	if houseId == "" || err != nil{
		fmt.Println("传入参数不完整", err)
		return
	}

	// 三种校验 大小 类型 防止重名    fastdfs
	if fileHeader.Size > 50000000 {
		fmt.Println("文件过大！")
		return
	}

	fileExt := path.Ext(fileHeader.Filename)
	if fileExt != ".png" && fileExt != ".jpg" {
		fmt.Println("文件类型错误,请重新选择")
		return
	}

	// 获取文件字节切片
	file, _ := fileHeader.Open()
	buf := make([]byte, fileHeader.Size)
	file.Read(buf)

	// 远程调用
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	resp, _ := microClient.UploadHouseImg(context.TODO(), &houseMicro.ImgReq{
		HouseId: houseId,
		ImgData: buf,
		FileExt: fileExt,
	})
	ctx.JSON(http.StatusOK, resp)


}

// 查看房源详细信息
func GetHouseDetailInfo(ctx *gin.Context) {
	fmt.Println("开始查询房屋详细信息")

	// 获取数据
	houseId := ctx.Param("id")
	if houseId == ""{
		fmt.Println("获取房屋id错误")
		return
	}
	userName := sessions.Default(ctx).Get("userName")

	// 远程调用
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	resp, _ := microClient.GetHouseDetail(context.TODO(), &houseMicro.DetailReq{
		HouseId: houseId,
		UserName: userName.(string),
	})

	// 返回信息
	ctx.JSON(http.StatusOK, resp)
}
