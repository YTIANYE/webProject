package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/consul"
	"image/png"
	"math/rand"
	"net/http"
	"time"
	"webProject/web/model"
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

	// 创建一个容器， 存储返回的数据信息
	resp := make(map[string]string)

	// 校验图片验证码 是否正确
	result := model.CheckImgCode(uuid, imgCode)
	if result {
		// 发送短信
		client, _err := CreateClient(tea.String("LTAI5tNjSPYJNT4wyTbWzYLr"), tea.String("B3R0SpewZr83NCyDVuiPlPDQwk4hYj"))
		if _err != nil {
			fmt.Println("发送短信错误 _err", _err)
		}

		// 生成一个随机的六位数，作为验证码
		rand.Seed(time.Now().UnixNano()) //播种随机种子
		// 生成6位随机数
		smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
		smsCodeString := "{" + "\"code\":" + "\"" + smsCode + "\"" + "}"

		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			SignName:      tea.String("阿里云短信测试"),
			TemplateCode:  tea.String("SMS_154950909"),
			PhoneNumbers:  tea.String(phone),
			TemplateParam: tea.String(smsCodeString),
		}

		// 打印 API 的返回值
		_, _err = client.SendSms(sendSmsRequest)
		if _err != nil {
			resp["errno"] = utils.RECODE_SMSERR
			resp["errno"] = utils.RecodeText(utils.RECODE_SMSERR)
			fmt.Println("发送短信错误 _err", _err)
		} else {
			// 发送短信成功
			resp["errno"] = utils.RECODE_OK
			resp["errmg"] = utils.RecodeText(utils.RECODE_OK)
			fmt.Println("发送短信成功")

			// 将 电话号：短信验证码 存入 Redis
			err := model.SaveSmsCode(phone, smsCode)
			if err != nil {
				resp["errno"] = utils.RECODE_DBWRITERR
				resp["errmg"] = utils.RecodeText(utils.RECODE_DBWRITERR)
				fmt.Println("存储短信验证码到redis失败：", err)
			}
		}

	} else {
		// 校验失败 发送错误信息
		resp["errno"] = utils.RECODE_CHECKERR
		resp["errno"] = utils.RecodeText(utils.RECODE_CHECKERR)
		fmt.Println("校验失败， 发送短信错误")
	}

	// 发送成功或者失败的 结果
	ctx.JSON(http.StatusOK, resp)

}
