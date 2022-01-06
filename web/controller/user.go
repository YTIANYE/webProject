package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/consul"
	"image/png"
	"net/http"
	"webProject/web/proto/getCaptcha"
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
	resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{})
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
