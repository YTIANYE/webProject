package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
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
