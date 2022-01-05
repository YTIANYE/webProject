package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

	png.Encode(ctx.Writer, img)

	fmt.Println("str:", str)
	fmt.Println("uuid:", uuid)

}
