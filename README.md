# 介绍

基于GOST开发的内网穿透管理平台，支持多用户、多节点，支持速率、连接数限制，中心化配置，通过网页修改配置，实时生效

**(此版本为开源版，去除了用户绑定QQ，客户端归属地区功能)**

[GOST开源地址](https://github.com/go-gost/gost)

如果想快速体验，可直接使用我提供的服务
[地址](https://gost.sian.one)

## License
This project is licensed under the [CC0-1.0 license](https://github.com/SianHH/gostc-open?tab=CC0-1.0-1-ov-file).

## 目录结构
```text
-- /
    -- server   // 后端项目代码
    -- web      // 前端项目代码
    -- client   // 节点和客户端项目代码
```

## 编译前端代码
```shell
cd web
npm i
npm run build
```
编译后的文件在web/dist目录中，压缩dist目录，将dist.zip复制到server/web/dist.zip

## 编译后端代码
```shell
cd server
go build -ldflags "-s -w" -o server main.go
```

## 编译节点/客户端(客户端和不开源版本通用)

编译方式一(当前平台)：
```shell
cd client
go build -ldflags "-s -w" -o gostc ./client/
```
编译后的文件为server/gostc

编译方式二(多平台编译，借助goreleaser工具)：
```shell
cd client
goreleaser release --snapshot --clean
```
编译后的文件在server/dist目录中

## 运行后台管理(服务端)
```shell
./server
```
默认端口8080，默认账号密码admin/admin

## 客户端和节点如何指定服务端地址
```shell
# 节点
./gostc -tls=false -addr 127.0.0.1:8080 -s -key ******

# 客户端
./gostc -tls=false -addr 127.0.0.1:8080 -key ******
```
如果后台管理服务不启用SSL，需要设置-tls=false
-addr指定的127.0.0.1:8080修改为实际的后台管理

