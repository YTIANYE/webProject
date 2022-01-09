package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// 创建全局结构体
type Student struct {
	Id   int //[在grom中Id字段]成为默认的主键，自增长   ---   主键索引，查询速度快
	Name string
	Age  int
}

func main() {

	// 连接数据库   --- 格式： 用户名：密码@协议（IP：port）/数据库名
	conn, err := gorm.Open("mysql", "tian:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		fmt.Println("grom.Open err:", err)
		return
	}
	defer conn.Close()

	conn.SingularTable(true) // 不要复数表名

	// 借助 grom 创建数据库表
	fmt.Println(conn.AutoMigrate(new(Student)).Error) // 执行创建过程打印结果

}
