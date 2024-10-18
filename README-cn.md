# 启动

## 容器运行：
1.为节省时间本地编译二进制包：
windows下运行如下文件：
windows.bat
linux下，运行指令：go build -ldflags="-w -s" -o ./main ./main.go 

2.启动容器
执行：
docker-compose up --build

后端服务80080端口，postgresql数据库端口5432


## 本地调试：
1.启动postgresql数据库：
docker compose up db -d
2.启动后端：
go run ./main.go


## 文件说明
repository 数据层，由sqlc构造sql。
servicces 控制层
handler   api接口
tests     mock文件存放
util      读取配置，初始化表



## 耗时：48h
## 代码查看

[路由](.\handler\handler.go)

## 单元测试
[handler](.\handler\accounts_test.go)
[servicces](.\services\accounts_test.go)
[sqlc](.\repository\sqlc\account_test.go)







api：

curl --location --request POST 'http://127.0.0.1:8080/accounts/creat' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8080' \
--header 'Connection: keep-alive' \
--data-raw '{
    "owner":"xhy",
    "currency": "USD"
}'
