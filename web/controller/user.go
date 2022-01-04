package controller

import (
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
	"image/color"
	"image/png"
	"net/http"
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
	// 生成图片验证码

	// 初始化对象
	cap := captcha.New()

	// 设置字体
	if err := cap.SetFont("web/conf/comic.ttf"); err != nil { //注意设置
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

	png.Encode(ctx.Writer, img)

	fmt.Println("str:", str)
	fmt.Println("uuid:", uuid)

}
