package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"webProject/web/controller"
	"webProject/web/model"
	// "webProject/web/model"
)

// 校验session 的 中间件

func LoginFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 初始化 Session对象
		s := sessions.Default(ctx)
		userName := s.Get("userName")

		if userName == nil {
			ctx.Abort() // 从这里返回，不必继续执行
		} else {
			ctx.Next() // 继续向下
		}
	}
}

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
	{ //注意 方法名和路径一定都要写正确
		r1.GET("session", controller.GetSession)
		r1.GET("/imagecode/:uuid", controller.GetImageCd)
		r1.GET("/smscode/:phone", controller.GetSmscd)
		r1.POST("/users", controller.PostRet)
		r1.GET("/areas", controller.GetArea)
		r1.POST("/sessions", controller.PostLogin)

		r1.Use(LoginFilter()) // 以后的路由都不需要校验 Session 了，直接获取数据即可

		r1.DELETE("/session", controller.DeleteSession)
		r1.GET("/user", controller.GetUserInfo)
		r1.PUT("/user/name", controller.PutUserInfo)
		r1.POST("/user/avatar", controller.PostAvatar)

		// 用户实名认证
		r1.GET("/user/auth", controller.GetUserInfo)
		r1.POST("user/auth", controller.PutUserAuth)

		// 房源信息相关
		// 查看用户已发布房源
		r1.GET("/user/houses", controller.GetUserHouses)
		// 发布房源信息
		r1.POST("/houses", controller.PostHouses)
		// 添加房源图片
		r1.POST("/houses/:id/images", controller.PostHousesImage)
		// 查看房源详细信息
		r1.GET("/houses/:id", controller.GetHouseDetailInfo)
		// 获取首页动画图片
		r1.GET("/house/index", controller.GetHouseIndex)
		// 搜索房屋
		r1.GET("/houses", controller.GetHouses)


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
