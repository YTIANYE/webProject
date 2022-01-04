package main

import (
	"github.com/gin-gonic/gin"
)

// 添加gin框架开发 3 步骤
func main() {
	// 初始化路由
	router := gin.Default()
	// 路由匹配
	router.GET("/", func(context *gin.Context) {
		context.Writer.WriteString("项目开始了......")
	})
	// 启动运行
	router.Run(":8080")

}
