#!/bin/bash

# MySQL登录路径
login_path=local

# 数据库名和表名
dbname=dbtest
tbname=tbtest

# 创建数据库
mysql --login-path=$login_path -e "CREATE DATABASE IF NOT EXISTS $dbname;"

# 进入数据库
mysql --login-path=$login_path $dbname << EOF

# 创建表
CREATE TABLE IF NOT EXISTS $tbname (
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255)
);

EOF

# 循环插入数据
while true; do
    # 生成随机8位字符串
    name=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)

    # 插入数据
    mysql --login-path=$login_path $dbname -e "INSERT INTO $tbname (name) VALUES ('$name');"

    # 获取当前时间戳
    timestamp=$(date +"%Y-%m-%d %H:%M:%S")

    # 记录日志
    echo "Inserted data: created_at=$timestamp, name=$name"

    # 等待1秒
    sleep 1
done
