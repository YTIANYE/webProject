package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// 创建第一种中间件

func CreateMid1(ctx *gin.Context) { //Test1 的类型就是Test2的返回类型
	fmt.Println("创建第一种中间件")
	fmt.Println("再创建第一种中间件")
}

// 创建第二种中间件/另一种格式的中间件

func CreateMid2() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("创建第二种中间件")
		fmt.Println("再创建第二种中间件")
	}
}

// 测试next()

func CreateMid1_next(ctx *gin.Context) { //Test1 的类型就是Test2的返回类型
	fmt.Println("创建第一种中间件")
	ctx.Next() //跳过当前中间件剩余内容， 去执行下一个中间件。 当所有操作执行完之后，以出栈的执行顺序返回，执行剩余代码。
	fmt.Println("再创建第一种中间件")
}

func CreateMid2_next() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("创建第二种中间件")
		context.Next()
		fmt.Println("再创建第二种中间件")
	}
}

// 测试return

func CreateMid2_return() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("创建第二种中间件")
		return // return   // 终止执行当前中间件的剩余内容，执行下一个中间件的内容
		context.Next()
		fmt.Println("再创建第二种中间件")
	}
}

// 测试Abort

func CreateMid2_Abort() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("创建第二种中间件")
		context.Abort() // 只执行当前中间件，操作完成后，执行上一个中间件
		fmt.Println("再创建第二种中间件")
	}
}

// 测试业务时间

func CreateMid1_Time(ctx *gin.Context) { //Test1 的类型就是Test2的返回类型
	fmt.Println("创建第一种中间件")
	t := time.Now()
	ctx.Next()
	fmt.Println(time.Now().Sub(t))
}

func Test() {
	router := gin.Default()

	// 使用中间件
	router.Use(CreateMid1) //注意两种调用方式
	router.Use(CreateMid2())

	router.GET("/test", func(context *gin.Context) {
		fmt.Println("使用中间件")
		context.Writer.WriteString("hello world!")
	})

	router.Run(":9999")
}

func TestNext() {
	router := gin.Default()

	// 使用中间件
	router.Use(CreateMid1_next) //注意两种调用方式
	router.Use(CreateMid2_next())

	router.GET("/test", func(context *gin.Context) {
		fmt.Println("使用中间件")
		context.Writer.WriteString("hello world!")
	})

	router.Run(":9999")
}

func TestReturn() {
	router := gin.Default()

	// 使用中间件
	router.Use(CreateMid1_next) //注意两种调用方式
	router.Use(CreateMid2_return())

	router.GET("/test", func(context *gin.Context) {
		fmt.Println("使用中间件")
		context.Writer.WriteString("hello world!")
	})

	router.Run(":9999")
}

func TestAbort() {
	router := gin.Default()

	// 使用中间件
	router.Use(CreateMid1_next) //注意两种调用方式
	router.Use(CreateMid2_Abort())

	router.GET("/test", func(context *gin.Context) {
		fmt.Println("使用中间件")
		context.Writer.WriteString("hello world!")
	})

	router.Run(":9999")
}

// 测试业务时间

func TestTime() {
	router := gin.Default()

	// 使用中间件

	// 测试第一种     6.902µs
	router.Use(CreateMid1_Time) //注意两种调用方式
	router.Use(CreateMid2_next())

	router.GET("/test", func(context *gin.Context) {
		fmt.Println("使用中间件")
		context.Writer.WriteString("hello world!")
	})

	/*	// 测试第二种		1.708µs
		router.Use(CreateMid1_Time) //注意两种调用方式
		router.Use(CreateMid2_Abort())

		router.GET("/test", func(context *gin.Context) {
			fmt.Println("使用中间件")
			context.Writer.WriteString("hello world!")
		})*/

	router.Run(":9999")
}

func main() {

	// Test()

	// TestNext()

	// TestReturn()

	TestAbort()

	// TestTime()

}
