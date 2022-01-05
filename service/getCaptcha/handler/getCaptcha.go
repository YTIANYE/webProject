package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"

	"image/color"

	getCaptcha "webProject/service/getCaptcha/proto/getCaptcha"
)

type GetCaptcha struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetCaptcha) Call(ctx context.Context, req *getCaptcha.Request, rsp *getCaptcha.Response) error {

	// 生成图片验证码

	// 初始化对象
	cap := captcha.New()

	// 设置字体
	if err := cap.SetFont("service/getCaptcha/conf/comic.ttf"); err != nil { //注意设置路径
		panic(err.Error())
	}

	// 设置验证码大小
	cap.SetSize(128, 64)

	// 设置干扰强度
	cap.SetDisturbance(captcha.NORMAL)

	// 设置前景色
	cap.SetFrontColor(color.RGBA{161, 47, 47, 255})

	// 设置背景色
	cap.SetBkgColor(color.RGBA{137, 190, 178, 128}, color.RGBA{69, 137, 148, 128})

	// 生成字体
	img, str := cap.Create(4, captcha.NUM)
	fmt.Println("生成的验证码 str:", str)

	// 将生成的图片序列化
	imgBuf, _ := json.Marshal(img)

	// 将 imgBuf 使用参数 rsp 传出
	rsp.Img = imgBuf

	return nil
}
