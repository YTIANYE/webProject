package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"webProject/web/controller"
	"webProject/web/model"
	// "webProject/web/model"
)

// 添加gin框架开发 3 步骤
// main 负责资源路径匹配
func main() {
	// 初始化 Redis 连接池
	model.InitRedis()

	// 初始化 MySQL 链接池
	model.InitDb()

	// 初始化路由
	router := gin.Default()

	// 初始化容器
	store, _ := redis.NewStore(10, "tcp", "192.168.17.129:6379", "", []byte("webProject"))

	// 使用容器
	router.Use(sessions.Sessions("mysession", store))

	//加载静态资源
	router.Static("/home", "web/view")
	// 添加路由分组
	r1 := router.Group("/api/v1.0")
	{
		r1.GET("session", controller.GetSession)
		r1.GET("/imagecode/:uuid", controller.GetImageCd)
		r1.GET("/smscode/:phone", controller.GetSmscd)
		r1.POST("/users", controller.PostRet)
		r1.GET("/areas", controller.GetArea)
		r1.POST("/sessions", controller.PostLogin)
		r1.DELETE("/session", controller.DeleteSession)
	}

	/*	// 处理 Session
		router.GET("/api/v1.0/session", controller.GetSesion) //回调函数这里只放函数名就可以了
		// 图片验证码
		router.GET("/api/v1.0/imagecode/:uuid", controller.GetImageCd)
		// 短信验证码
		router.GET("/api/v1.0/smscode/:phone", controller.GetSmscd)*/

	// 启动运行
	router.Run(":8080")

}
