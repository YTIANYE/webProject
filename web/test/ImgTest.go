package main

import (
	"github.com/afocus/captcha"
	"image/color"
	"image/png"
	"net/http"
)

var cap *captcha.Captcha

func main() {
	// 初始化对象
	cap := captcha.New()

	// 设置字体
	if err := cap.SetFont("web/test/comic.ttf"); err != nil { //注意设置
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
	// img, str := cap.Create(4, captcha.NUM)

	// 将图片验证码展示到页面中
	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		img, str := cap.Create(4, captcha.NUM)
		png.Encode(w, img)
		println(str)
	})

	// 启动服务
	http.ListenAndServe(":8085", nil)
}
