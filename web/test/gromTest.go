package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

/*
// 创建全局结构体
type Student struct {
	Id   int //[在grom中Id字段]成为默认的主键，自增长   ---   主键索引，查询速度快
	Name string
	Age  int
}*/

/*type Student struct {
	gorm.Model // golang中，匿名成员 --- 继承！   gorm自动维护
	Name       string
	Age        int
}*/

// 创建全局结构体
type Student struct {
	Id    int    //[在grom中Id字段]成为默认的主键，自增长   ---   主键索引，查询速度快
	Name  string `gorm:"size:100;default:'xiaoming'"` // string -- varcha  默认大小255.可以在创建表时，指定表的大小。默认值 xiaoming
	Age   int
	Class int       `gorm:"not null"`
	Join  time.Time `gorm:"type:timestamp"` // 创建 Student 表指定 timestamp类型。
}

// 创建全局连接池句柄
var GlobalConn *gorm.DB

func main() {

	// 连接数据库---获取连接池句柄   --- 格式： 用户名：密码@协议（IP：port）/数据库名
	conn, err := gorm.Open("mysql", "tian:password@tcp(127.0.0.1:3306)/test?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("grom.Open err:", err)
		return
	}
	// defer conn.Close()
	GlobalConn = conn //连接池句柄不能关闭

	GlobalConn.DB().SetMaxIdleConns(10)
	GlobalConn.DB().SetMaxOpenConns(100)

	GlobalConn.SingularTable(true) // 不要复数表名

	// 借助 grom 创建数据库表
	fmt.Println(GlobalConn.AutoMigrate(new(Student)).Error) // 执行创建过程打印结果

	//// 插入数据
	//InsertData()

	//// 查询数据
	//SearchData()

	//// 更新数据
	//UpdateDate()

	//// 删除数据
	//DeleteData()

}

func InsertData() {
	// 先创建数据   --- 创建对象
	var stu Student
	stu.Name = "zhangsan"
	stu.Age = 18

	// 插入（创建）数据
	fmt.Println(GlobalConn.Create(&stu).Error) //注意 &
}

func SearchData() {

	// 查询第一条
	//var stu Student
	//// GlobalConn.First(&stu) //查询第一条的全部信息
	//GlobalConn.Select("name, age").First(&stu) // 只是查询第一条 name age
	//
	//fmt.Println(stu)

	//// 查询多条
	//var stus []Student
	//GlobalConn.Find(&stus)
	//// GlobalConn.Select("name, age").Find(&stus) //Find()查询多条 name age
	//fmt.Println(stus)

	// WHERE 字句
	var stus []Student
	// GlobalConn.Select("name, age").Where("name = ?", "list").Find(&stus)//查询姓名为list 的 name age
	// GlobalConn.Select("name, age").Where("name = ?", "list").Where("age = ?", 28).Find(&stus)
	GlobalConn.Select("name, age").Where("name = ? and age = ?", "list", 28).Find(&stus)

	// GlobalConn.Find(&stus)// 无法查询软删除的数据
	// GlobalConn.Unscoped().Find(&stus) //可以查询软删除的数据
	fmt.Println(stus)

}

func UpdateDate() {

	////Model(new(Student): 指定更新 “student” 表
	////Where("name = ?", "zhaoliu")： 指定过滤条件。
	////Update("name", "lisi").Error)：指定 把 “list” 更新成 “zhaoliu”

	//fmt.Println(GlobalConn.Model(new(Student)).
	//	Where("name = ?", "list").
	//	Update("name", "zhaoliu").Error)

	//// 使用`map`更新多个属性，只会更新这些更改的字段
	//db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
	////// UPDATE users SET name='hello', age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

	fmt.Println(GlobalConn.Model(new(Student)).
		Where("name = ?", "zhaoliu").
		Updates(map[string]interface{}{"name": "wangwu", "age": 50}).Error) //interface{} 任意数据类型

}

func DeleteData() {
	// fmt.Println(GlobalConn.Where("name = ?", "zhangsan").Delete(new(Student)).Error) //删除后 deleted_at 值为删除时间,反复删除，时间也是第一次的时间
	fmt.Println(GlobalConn.Unscoped().Where("name = ?", "zhangsan").Delete(new(Student)).Error) // 物理删除
}
