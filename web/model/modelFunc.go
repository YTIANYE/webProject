package model

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 创建全局连接池 redis 句柄

var RedisPool redis.Pool

// 创建函数 初始化连接池

func InitRedis() {
	// 连接 Redis 连接池
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "192.168.17.129:6379")
		},
	}
}

// 校验图片验证码

func CheckImgCode(uuid, imgCode string) bool {
	// 连接redis  ---  从连接池中获取链接
	/*	conn, err := redis.Dial("tcp", "192.168.17.129:6379")
		if err != nil {
			fmt.Println("redis.Dial err:", err)
		}*/
	conn := RedisPool.Get()
	defer conn.Close()

	// 查询redis中的数据
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误 err:", err)
		return false
	}

	// 返回校验结果
	return code == imgCode
}

// 存储短信验证码

func SaveSmsCode(phone, code string) error {
	// 从 redis 连接池获取一条连接
	conn := RedisPool.Get()
	defer conn.Close() //相当于释放链接，不是关闭链接

	// 存储短信验证码 到redis中
	_, err := conn.Do("setex", phone+"_code", 60*3, code)
	return err
}

// 校验短信验证码

func CheckSmsCode(phone, code string) error {
	// 链接redis
	conn := RedisPool.Get()

	// 从 redis 中, 根据 key 获取 Value --- 短信验证码  码值
	smsCode, err := redis.String(conn.Do("get", phone+"_code"))
	if err != nil {
		fmt.Println("redis get phone_code err:", err)
		return err
	}
	// 验证码匹配  失败
	if smsCode != code {
		return errors.New("验证码匹配失败!")
	}
	// 匹配成功!
	return nil
}

// 注册用户信息,写 MySQL 数据库.

func RegisterUser(mobile, pwd string) error {
	var user User
	user.Name = mobile // 默认使用手机号作为用户名

	// 使用 md5 对 pwd 加密
	m5 := md5.New()                             // 初始md5对象
	m5.Write([]byte(pwd))                       // 将 pwd 写入缓冲区
	pwd_hash := hex.EncodeToString(m5.Sum(nil)) // 不使用额外的秘钥

	user.Password_hash = pwd_hash

	// 插入数据到MySQL
	return GlobalConn.Create(&user).Error
}

// 处理登录业务，根据手机号和密码获取用户名

func Login(mobile, pwd string) (string, error) {
	var user User

	// 对pwd做md5 hash
	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum(nil))

	// 查询
	// 相当于 select name from user where mobil = xxx and password_hash = xxx ;
	err := GlobalConn.Where("mobile = ?", mobile).
		Where("password_hash = ?", pwd_hash).
		Select("name").
		Find(&user).Error

	return user.Name, err
}

// 获取用户信息

func GetUserInfo(userName string) (User, error) {

	// 实现 SQL select * from user where name = userName;
	var user User
	err := GlobalConn.Where("name = ?", userName).First(&user).Error
	return user, err
}

// 更新用户名

func UpdateUserName(newName, oldName string) error {
	// update user set name = 'username01' where name = 旧用户名
	return GlobalConn.Model(new(User)).Where("name = ?", oldName).Update("name", newName).Error
}

// 根据用户名，更新用户头像

func UpdateAvatar(userName, avatar string) error {
	// update user set avatar_url = avatar, where name = username
	return GlobalConn.Model(new(User)).Where("name = ?", userName).Update("avatar_url", avatar).Error
}
