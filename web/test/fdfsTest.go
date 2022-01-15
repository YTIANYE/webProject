package main

import (
	"fmt"
	"github.com/tedcy/fdfs_client"
)

func main() {
	//  初始化客户端  ---   配置文件
	clt, err := fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")
	if err != nil {
		fmt.Println("初始化fdfs客户端错误：", err)
		return
	}

	// 上传文件   ---  尝试文件名上传   传入到storage
	resp, err := clt.UploadByFilename("头像.jpg")

	fmt.Println(resp, err)
}
