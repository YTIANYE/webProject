#! /bin/bash

# 启动 consul
consul agent -dev

# 启动redis
# sudo redis-server /etc/redis/redis.conf
sudo systemctl status redis-server
# redis-cli -h 192.168.17.129 -p 6379 --raw

# 启动mysql
# mysql -utian -ppassword

# 启动fdfs
sudo fdfs_trackerd /etc/fdfs/tracker.conf
sudo fdfs_storaged /etc/fdfs/storage.conf

# 启动nginx
sudo /usr/local/nginx/sbin/nginx




# 启动微服务

go run ~/go/src/webProject/service/getCaptcha/main.go

go run ~/go/src/webProject/service/user/main.go

# go run ~/go/src/webProject/web/main.go





:<< block
# 这一部分被多行注释掉了
block
