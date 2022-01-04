package main

import (
	"github.com/gin-gonic/gin"
	"webProject/web/controller"
)

// 添加gin框架开发 3 步骤
// main 负责资源路径匹配
func main() {
	// 初始化路由
	router := gin.Default()

	// 路由匹配
	/*	router.GET("/", func(context *gin.Context) {
		context.Writer.WriteString("项目开始了......")
	})*/

	// router.Static("/", "web/view")                        //加载静态资源
	router.Static("/home", "web/view") //加载静态资源
	// 处理 Session
	router.GET("/api/v1.0/session", controller.GetSesion) //回调函数这里只放函数名就可以了
	//
	router.GET("/api/v1.0/imagecode/:uuid", controller.GetImageCd)

	// 启动运行
	router.Run(":8080")

}
