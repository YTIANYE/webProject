package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"service/house/proto/house"
	"strconv"
	"time"
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
	if err := GlobalConn.Where("name = ?", userName).Find(&user).Error; err != nil { //不要忘记加 ？
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

// 房屋图片存入数据库

func SaveHouseImg(houseId, imgPath string) error {
	return GlobalConn.Model(new(House)).
		Where("id = ?", houseId).
		Update("index_image_url", imgPath).Error
}

// 获取房屋详细信息

func GetHouseDetail(houseId, userName string) (house.DetailData, error) {
	var respData house.DetailData     // 返回数据
	var houseDetail house.HouseDetail //

	// 1. 查询房屋信息
	var houseInfo House
	if err := GlobalConn.Where("id = ?", houseId).Find(&houseInfo).Error; err != nil {
		fmt.Println("查询房屋信息错误：", err)
		return respData, nil
	}
	// 查询房屋信息成功
	{
		houseDetail.Acreage = int32(houseInfo.Acreage)
		houseDetail.Address = houseInfo.Address
		houseDetail.Beds = houseInfo.Beds
		houseDetail.Capacity = int32(houseInfo.Capacity)
		houseDetail.Deposit = int32(houseInfo.Deposit)
		houseDetail.Hid = int32(houseInfo.ID)
		houseDetail.MaxDays = int32(houseInfo.Max_days)
		houseDetail.MinDays = int32(houseInfo.Min_days)
		houseDetail.Price = int32(houseInfo.Price)
		houseDetail.RoomCount = int32(houseInfo.Room_count)
		houseDetail.Title = houseInfo.Title
		houseDetail.Unit = houseInfo.Unit
		if houseInfo.Index_image_url != "" {
			houseDetail.ImgUrls = append(houseDetail.ImgUrls, "http://192.168.17.129:8888/"+houseInfo.Index_image_url)
		}
	}

	// 2. 查询评论
	var orders []OrderHouse
	if err := GlobalConn.Model(&houseInfo).Related(&orders).Error; err != nil {
		fmt.Println("查询评论信息错误")
		return respData, nil
	}
	// 查询评论信息成功
	for _, v := range orders {
		var commentTemp house.CommentData
		commentTemp.Comment = v.Comment
		commentTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		var tempUser User
		GlobalConn.Model(&v).Related(&tempUser)
		commentTemp.UserName = tempUser.Name

		houseDetail.Comments = append(houseDetail.Comments, &commentTemp)
	}

	// 3. 获取房屋家具信息
	var facs []Facility
	if err := GlobalConn.Model(&houseInfo).Related(&facs, "Facilities").Error; err != nil {
		fmt.Println("获取房屋家具信息错误：", err)
		return respData, nil
	}
	// 查询家具信息成功
	for _, v := range facs {
		houseDetail.Facilities = append(houseDetail.Facilities, int32(v.Id))
	}

	// 4. 获取幅图片
	var imgs []HouseImage
	if err := GlobalConn.Model(&houseInfo).Related(&imgs).Error; err != nil {
		fmt.Println("该房屋没有副图片")
	}
	// 查询房屋副图片成功
	for _, v := range imgs {
		if len(imgs) != 0 {
			houseDetail.ImgUrls = append(houseDetail.ImgUrls, "http://192.168.17.129:8888/"+v.Url)
		}
	}

	// 5. 获取房屋所有者信息
	var user User
	if err := GlobalConn.Model(&houseInfo).Related(&user).Error; err != nil {
		fmt.Println("查询房屋所有者错误：", err)
		return respData, nil
	}
	// 查询房屋所有者信息成功
	houseDetail.UserName = user.Name
	houseDetail.UserAvatar = user.Avatar_url
	houseDetail.UserId = int32(user.ID)

	// 6. 完成房屋具体信息查询
	respData.House = &houseDetail

	// 获取浏览人信息
	var nowUser User
	if err := GlobalConn.Where("name = ?", userName).Find(&nowUser).Error; err != nil {
		fmt.Println("查询当前浏览人信息错误", err)
		return respData, nil
	}
	respData.UserId = int32(nowUser.ID)
	return respData, nil
}

// 获取房屋信息
func GetIndexHouse() ([]*house.Houses, error) {
	var housesResp []*house.Houses
	var houses []House
	if err := GlobalConn.Limit(5).Find(&houses).Error; err != nil {
		fmt.Println("获取房屋信息失败", err)
		return nil, err
	}
	for _, v := range houses {
		var houseTemp house.Houses
		houseTemp.Address = v.Address
		//根据房屋信息获取地域信息
		var area Area
		var user User
		GlobalConn.Model(&v).Related(&area).Related(&user)
		houseTemp.AreaName = area.Name
		houseTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseTemp.HouseId = int32(v.ID)
		houseTemp.ImgUrl = "http://192.168.17.129:8888/" + v.Index_image_url
		houseTemp.OrderCount = int32(v.Order_count)
		houseTemp.Price = int32(v.Price)
		houseTemp.RoomCount = int32(v.Room_count)
		houseTemp.Title = v.Title
		houseTemp.UserAvatar = "http://192.168.17.129:8888/" + user.Avatar_url

		housesResp = append(housesResp, &houseTemp)
	}

	return housesResp, nil
}

// 搜索房屋
func SearhHouse(areaId, sd, ed, sk string) ([]*house.Houses, error) {
	var houseInfos []House

	//   minDays  <  (结束时间  -  开始时间) <  max_days
	//计算一个差值  先把string类型转为time类型
	sdTime, _ := time.Parse("2006-01-02", sd)
	edTime, _ := time.Parse("2006-01-02", ed)
	dur := edTime.Sub(sdTime)

	err := GlobalConn.Where("area_id = ?", areaId).
		Where("min_days < ?", dur.Hours()/24).
		Where("max_days > ?", dur.Hours()/24).
		Order("created_at desc").Find(&houseInfos).Error
	if err != nil {
		fmt.Println("搜索房屋失败",err)
		return nil,err
	}

	var housesResp []*house.Houses

	for _, v := range houseInfos{
		var houseTemp house.Houses
		houseTemp.Address = v.Address
		// 查询房屋对应的地域信息
		var area Area
		var user User
		GlobalConn.Model(&v).Related(&area).Related(&user)

		houseTemp.AreaName = area.Name
		houseTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseTemp.HouseId = int32(v.ID)
		houseTemp.ImgUrl = "http://192.168.17.129:8888/"+v.Index_image_url
		houseTemp.OrderCount = int32(v.Order_count)
		houseTemp.Price = int32(v.Price)
		houseTemp.RoomCount = int32(v.Room_count)
		houseTemp.Title = v.Title
		houseTemp.UserAvatar = "http://192.168.17.129:8888/"+user.Avatar_url

		housesResp = append(housesResp, &houseTemp)
	}

	return housesResp, nil
}
