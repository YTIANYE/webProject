package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func cookieTest() {
	router := gin.Default()

	router.GET("/test", func(context *gin.Context) {

		// 设置 Cookie
		// context.SetCookie("mytest", "chuanzhi", 60*60, "", "", true, true)
		// context.SetCookie("mytest", "chuanzhi", 60*60, "", "", false, true)
		// context.SetCookie("mytest", "chuanzhi", 0, "", "", false, true)

		// 获取Cookie
		cookieVal, _ := context.Cookie("mytest")
		fmt.Println("获取的Cookie 为：", cookieVal)

		context.Writer.WriteString("测试Cookie....")
	})

	router.Run(":9999")
}

func sessionTest() {

	router := gin.Default()

	// 初始化容器
	store, _ := redis.NewStore(10, "tcp", "192.168.17.129:6379", "", []byte("webProject"))

	// 设置临时session
	/*	store.Options(sessions.Options{
		MaxAge: 0,
	})*/

	// 使用容器
	router.Use(sessions.Sessions("mysession", store))
	router.GET("/test", func(context *gin.Context) {
		// 调用session设置session数据
		session := sessions.Default(context)

		/*		//设置session
				session.Set("key", "value")
				// 修改session时，需要Save，否则不生效
				session.Save()*/

		// 获取session
		value := session.Get("key")
		fmt.Println("获取 Session:", value.(string))

		context.Writer.WriteString("测试 Session ...")
	})

	router.Run(":9999")
}

func main() {
	//// 测试cookie
	//cookieTest()

	// 测试Session
	sessionTest()

}
