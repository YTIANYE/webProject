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
	"image/png"
	"net/http"
	"webProject/web/model"
	"webProject/web/proto/getCaptcha"
	userMicro "webProject/web/proto/user"
	"webProject/web/utils"
)

// 获取 session 信息
func GetSesion(ctx *gin.Context) {
	// 初始化一个错误返回的 map
	resp := make(map[string]string)
	resp["errno"] = utils.RECODE_SESSIONERR
	resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)

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
		s := sessions.Default(ctx) //初始化session
		s.Set("useName", userName) // 将用户名设置到session中
		s.Save()

	} else {
		// 登录失败
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmg"] = utils.RecodeText(utils.RECODE_LOGINERR)
	}

	ctx.JSON(http.StatusOK, resp)
}
