package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/tedcy/fdfs_client"
	"image/png"
	"net/http"
	"path"
	"webProject/web/model"
	"webProject/web/proto/getCaptcha"
	houseMicro "webProject/web/proto/house"
	userMicro "webProject/web/proto/user"
	"webProject/web/utils"
)

// 获取 session 信息

/*func GetSession(ctx *gin.Context) {
	// 初始化一个错误返回的 map
	resp := make(map[string]string)
	resp["errno"] = utils.RECODE_SESSIONERR
	resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)

	ctx.JSON(http.StatusOK, resp)
}*/

func GetSession(ctx *gin.Context) {
	resp := make(map[string]interface{})

	// 获取 Session数据
	session := sessions.Default(ctx) // 初始化Session 对象
	userName := session.Get("userName")

	if userName == nil {
		// 用户没有登录   ---   没有存在MySQL中， 也没存在session中
		fmt.Println("用户不存在")
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	} else {
		// 用户存在
		fmt.Println("用户存在")
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		var nameData struct {
			Name string `json:"name"`
		}
		nameData.Name = userName.(string) // 类型断言
		resp["data"] = nameData
	}

	ctx.JSON(http.StatusOK, resp)
}

// 获取图片信息

func GetImageCd(ctx *gin.Context) {
	// 获取图片验证码的 uuid
	uuid := ctx.Param("uuid")

	// 指定 consul 的服务发现
	consulReg := consul.NewRegistry()
	consulService := micro.NewService(
		micro.Registry(consulReg),
	)

	// 初始化客户端
	microClient := getCaptcha.NewGetCaptchaService("go.micro.srv.getCaptcha", consulService.Client())

	// 调用远程函数
	resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{Uuid: uuid})
	if err != nil {
		fmt.Println("没有找到远程服务,err:", err)
		return
	}

	// 将得到的数据反序列化，得到图片数据
	var img captcha.Image
	json.Unmarshal(resp.Img, &img)

	// 将图片写出到 浏览器
	png.Encode(ctx.Writer, img)

	fmt.Println("uuid:", uuid)

}

// 获取短信验证码

func GetSmscd(ctx *gin.Context) {
	// 获取短信验证码的 uuid
	phone := ctx.Param("phone")
	// 拆分 GET 请求中的 URL   === 格式： 资源路径 ？key=value&key=value&key=value
	imgCode := ctx.Query("text")
	uuid := ctx.Query("id")
	fmt.Println("---out---:", "phone:", phone, "imgCode:", imgCode, "uuid", uuid)

	// 指定 consul 的服务发现
	consulReg := consul.NewRegistry()
	consulService := micro.NewService(
		micro.Registry(consulReg),
	)

	// 初始化客户端
	microClient := userMicro.NewUserService("go.micro.srv.user", consulService.Client())

	// 调用远程函数
	resp, err := microClient.SendSms(context.TODO(), &userMicro.Request{Phone: phone, Uuid: uuid, ImgCode: imgCode})
	if err != nil {
		fmt.Println("调用远程函数 SendMsm 失败：", err)
		return
	}

	// 发送校验结果给浏览器
	ctx.JSON(http.StatusOK, resp)

}

// 发送注册信息

func PostRet(ctx *gin.Context) {

	// 获取数据
	var regData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	ctx.Bind(&regData)

	fmt.Println("获取到的数据为:", regData)

	// 初始化consul
	microService := utils.InitMicro()

	// 初始化客户端
	microClient := userMicro.NewUserService("go.micro.srv.user", microService.Client())

	// 调用远程函数
	resp, err := microClient.Register(context.TODO(), &userMicro.RegReq{
		Mobile:   regData.Mobile,
		SmsCode:  regData.SmsCode,
		Password: regData.PassWord,
	})
	if err != nil {
		fmt.Println("注册用户, 找不到远程服务!", err)
		return
	}
	// 写给浏览器
	ctx.JSON(http.StatusOK, resp)
}

// 获取地域信息

func GetArea(ctx *gin.Context) {
	// 从缓存redis, 中获取数据
	conn := model.RedisPool.Get()
	// 当初使用“字节切片”，现在使用切片类型接收
	areaData, _ := redis.Bytes(conn.Do("get", "areaData"))

	var areas []model.Area
	if len(areaData) == 0 { // 没有从redis中获取数据
		fmt.Println("从mysql中获取数据")

		// 先从MySQL中获取数据.
		model.GlobalConn.Find(&areas)

		// 再把数据写入到 redis 中.  存储结构体序列化后的JSON串
		areaBuf, _ := json.Marshal(areas)
		conn.Do("set", "areaData", areaBuf)
	} else { // redis 中有数据
		fmt.Println("从redis中获取数据")

		json.Unmarshal(areaData, &areas) //反序列化

	}

	resp := make(map[string]interface{})

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = areas
	fmt.Println("resp:", resp)

	ctx.JSON(http.StatusOK, resp)
}

// 处理登陆业务

func PostLogin(ctx *gin.Context) {
	// 获取前端数据
	var loginData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
	}
	ctx.Bind(&loginData)

	resp := make(map[string]interface{})

	// 获取数据库数据，查询是否和输入的数据匹配
	userName, err := model.Login(loginData.Mobile, loginData.PassWord)
	if err == nil {
		// 登录成功
		resp["errno"] = utils.RECODE_OK
		resp["errmg"] = utils.RecodeText(utils.RECODE_OK)

		// 将 登录状态，保存在Session中
		s := sessions.Default(ctx)  //初始化session
		s.Set("userName", userName) // 将用户名设置到session中				// 注意不要写错啊，找了半天的bug,get总为nil
		s.Save()

	} else {
		// 登录失败
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmg"] = utils.RecodeText(utils.RECODE_LOGINERR)
	}

	ctx.JSON(http.StatusOK, resp)
}

// 退出登录

func DeleteSession(ctx *gin.Context) {

	resp := make(map[string]interface{})

	// 初始化 Session 对象
	session := sessions.Default(ctx)

	// 删除session数据
	session.Delete("userName") // 没有返回值

	// 必须使用Save()保存
	err := session.Save() //有返回值
	if err != nil {
		// 退出失败
		resp["errno"] = utils.RECODE_LOGOUTERR
		resp["errmg"] = utils.RecodeText(utils.RECODE_LOGOUTERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmg"] = utils.RecodeText(utils.RECODE_OK)
	}

	ctx.JSON(http.StatusOK, resp)
}

// 获取用户的基本信息  // 检查用户实名认证

func GetUserInfo(ctx *gin.Context) {

	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)

	// 获取Session,  得到当前用户的信息
	s := sessions.Default(ctx)
	userName := s.Get("userName")

	// 判断用户名是否存在
	if userName == nil { // 用户每登录，但是进入该页面，恶意进入
		fmt.Println("用户未登录")
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return //如果出错，报错， 直接退出程序
	}

	// 根据用户名获取用户信息   --- 查MySQL数据库 user 表
	user, err := model.GetUserInfo(userName.(string))
	if err != nil {
		fmt.Println("访问数据库出错")
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return
	}

	//
	resp["errno"] = utils.RECODE_OK
	resp["errmg"] = utils.RecodeText(utils.RECODE_OK)

	temp := make(map[string]interface{})
	temp["user_id"] = user.ID
	temp["name"] = user.Name
	temp["mobile"] = user.Mobile
	temp["real_name"] = user.Real_name
	temp["id_card"] = user.Id_card
	temp["avatar_url"] = "http://192.168.17.129:8888/" + user.Avatar_url

	resp["data"] = temp

}

// 更新用户名

func PutUserInfo(ctx *gin.Context) {
	// 从session获取当前用户名
	session := sessions.Default(ctx) // 初始化一个session对象
	userName := session.Get("userName")

	// 获取新用户名   ---   处理 Request Payload 数据类型.Bind()
	var nameData struct {
		Name string `json:"name"`
	}
	ctx.Bind(&nameData)

	// 更新用户名
	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)

	// 更新数据库中的 name
	err := model.UpdateUserName(nameData.Name, userName.(string))
	if err != nil {
		fmt.Println("用户名更新失败：", err)
		resp["errno"] = utils.RECODE_DBERR
		resp["errmg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}

	// 更新Session数据
	session.Set("userName", nameData.Name)

	err = session.Save()
	if err != nil {
		fmt.Println("更新session错误：", err)
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return
	}

	resp["errno"] = utils.RECODE_OK
	resp["errmg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = nameData
}

//  上传头像

func PostAvatar(ctx *gin.Context) {

	// 获取图片文件，得到静态文件对象
	file, _ := ctx.FormFile("avatar")

	// 上传文件到项目中
	/*	err := ctx.SaveUploadedFile(file, "web/img/"+file.Filename)
		if err != nil {
			fmt.Println("上传文件到项目中错误：", err)
		}*/

	// 上传头像到fastdfs中
	clt, _ := fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")

	// 打开文件读取内容
	f, _ := file.Open() // 只读打开

	buf := make([]byte, file.Size) // 按文件世纪大小，创建切片

	f.Read(buf) //读取文件内容， 保存至buf缓冲区

	// go语言根据文件名获取文件后缀
	fileExt := path.Ext(file.Filename) //传文件名， 获取文件后缀   ----  带有“.”

	// 按字节流上传图片内容
	remoteId, err := clt.UploadByBuffer(buf, fileExt[1:])
	if err != nil {
		fmt.Println("图片上传错误：", err)
	}

	// 获取session, 得到当前用户
	userName := sessions.Default(ctx).Get("userName")

	// 根据用户名，更新用户头像    ----   MySQL数据库
	model.UpdateAvatar(userName.(string), remoteId)

	resp := make(map[string]interface{})
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	temp := make(map[string]interface{})
	temp["avatar_url"] = "http://192.168.17.129:8888/" + remoteId
	resp["data"] = temp

	ctx.JSON(http.StatusOK, resp)

}

// 保存用户实名认证信息

func PutUserAuth(ctx *gin.Context) {
	// 从session获取当前用户名
	session := sessions.Default(ctx) // 初始化一个session对象
	userName := session.Get("userName")

	// 保存认证信息   ---   处理 Request Payload 数据类型.Bind()
	var authData struct {
		RealName string `json:"real_name"`
		IdCard   string `json:"id_card"`
	}
	ctx.Bind(&authData)

	// 初始化consul
	microService := utils.InitMicro()

	// 初始化客户端
	microClient := userMicro.NewUserService("go.micro.srv.user", microService.Client())

	// 调用远程函数
	resp, err := microClient.AuthUpdate(context.TODO(), &userMicro.AuthReq{
		UserName: userName.(string),
		RealName: authData.RealName,
		IdCard:   authData.IdCard,
	})
	if err != nil {
		fmt.Println("保存实名认证信息, 找不到远程服务!", err)
		return
	}

	// 更新信息
	ctx.JSON(http.StatusOK, resp)

}

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
