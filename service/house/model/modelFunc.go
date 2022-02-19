package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"service/house/proto/house"
	"strconv"
)

// 创建全局连接池 redis 句柄
var RedisPool redis.Pool

// 创建函数 初始化连接池
func InitRedis() {
	// 连接 Redis 连接池
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "192.168.17.129:6379")
		},
	}
}

// 获取用户发布的房源信息

func GetUserHouse(userName string) ([]*house.Houses, error) {

	var houseInfos []*house.Houses

	// 判断用户是否存在
	var user User
	if err := GlobalConn.Where("name = ?", userName).Find(&user).Error; err != nil {//不要忘记加 ？
		fmt.Println("userName:", userName)
		fmt.Print("获取当前用户信息错误：", err)
		return nil, err
	}

	// 查询房源信息  一对多
	var houses []House
	GlobalConn.Model(&user).Related(&houses) //related函数可以是以主表关联从表,也可以是以从表关联主表

	for _, v := range houses {
		var houseInfo house.Houses
		houseInfo.Title = v.Title
		houseInfo.Address = v.Address
		houseInfo.Ctime = v.CreatedAt.Format("2022-01-02 12:04:15")
		houseInfo.HouseId = int32(v.ID)
		houseInfo.ImgUrl = "http://192.168.17.129:8888/" + v.Index_image_url
		houseInfo.OrderCount = int32(v.Order_count)
		houseInfo.Price = int32(v.Price)
		houseInfo.RoomCount = int32(v.Room_count)
		houseInfo.UserAvatar = "http://192.168.17.129:8888/" + user.Avatar_url

		// 获取地域信息
		var area Area
		GlobalConn.Where("id = ?", v.AreaId).Find(&area)
		houseInfo.AreaName = area.Name

		houseInfos = append(houseInfos, &houseInfo)
	}
	return houseInfos, nil
}

// 发布用户房屋信息

func AddHouse(req *house.PubReq) (int, error) {

	//根据userName获取userId
	var user User
	if err := GlobalConn.Where("name = ?", req.UserName).Find(&user).Error; err != nil {
		fmt.Println("查询当前用户失败", err)
		return 0, err
	}

	// 给house赋值
	var houseInfo House
	houseInfo.Address = req.Address
	price, _ := strconv.Atoi(req.Price)
	roomCount, _ := strconv.Atoi(req.RoomCount)
	houseInfo.Price = price
	houseInfo.Room_count = roomCount
	houseInfo.Unit = req.Unit
	houseInfo.Capacity, _ = strconv.Atoi(req.Capacity)
	houseInfo.Beds = req.Beds
	houseInfo.Deposit, _ = strconv.Atoi(req.Deposit)
	houseInfo.Min_days, _ = strconv.Atoi(req.MinDays)
	houseInfo.Max_days, _ = strconv.Atoi(req.MaxDays)
	houseInfo.Acreage, _ = strconv.Atoi(req.MaxDays)

	//sql中一对多插入,只是给外键赋值
	houseInfo.UserId = uint(user.ID)
	houseInfo.Title = req.Title
	areaId, _ := strconv.Atoi(req.AreaId)
	houseInfo.AreaId = uint(areaId)

	// 房屋设施
	for _, v := range req.Facility {
		id, _ := strconv.Atoi(v)
		var fac Facility
		if err := GlobalConn.Where("id = ?", id).First(&fac).Error; err != nil {
			fmt.Println("家具id错误", err)
			return 0, err
		}
		// 查询到数据
		houseInfo.Facilities = append(houseInfo.Facilities, &fac)
	}

	// 插入数据
	if err := GlobalConn.Create(&houseInfo).Error; err != nil {
		fmt.Println("房屋信息写入数据库错误", err)
		return 0, err
	}
	return int(houseInfo.ID), nil
}
